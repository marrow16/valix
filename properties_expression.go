package valix

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type BooleanOperator int

const (
	And BooleanOperator = iota
	Or
)

// MustParseExpression is the same as ParseExpression - but panics if there is an error
func MustParseExpression(expr string) OthersExpr {
	result, err := ParseExpression(expr)
	if err != nil {
		panic(err)
	}
	return result
}

// ParseExpression parses a string expression representing the presence or non-presence of named
// properties in an object
//
// An example:
//   expr, err := valix.ParseExpression(`(foo && bar) || (foo && baz) || (bar && baz) && !(foo && bar && baz)`)
func ParseExpression(expr string) (OthersExpr, error) {
	tokens, groupsCount, err := parseExpressionTokens(expr)
	if err != nil {
		return nil, err
	}
	result := make([]Other, 0, len(tokens))
	nextNotted := false
	nextOperator := And
	var currentGroup *OtherGrouping = nil
	groupStack := make([]*OtherGrouping, 0, groupsCount)
	for _, token := range tokens {
		switch token.tokenType {
		case tokenTypeNot:
			nextNotted = !nextNotted
		case tokenTypeOperator:
			nextOperator = token.operator
		case tokenTypeName:
			newItem := &OtherProperty{
				Name: token.name,
				Not:  nextNotted,
				Op:   nextOperator,
			}
			nextNotted = false
			nextOperator = And
			if currentGroup != nil {
				currentGroup.Of = append(currentGroup.Of, newItem)
			} else {
				result = append(result, newItem)
			}
		case tokenTypeGroupStart:
			newGroup := &OtherGrouping{
				Of:  OthersExpr{},
				Not: nextNotted,
				Op:  nextOperator,
			}
			nextNotted = false
			nextOperator = And
			if currentGroup != nil {
				groupStack = append(groupStack, currentGroup)
				currentGroup.Of = append(currentGroup.Of, newGroup)
			} else {
				result = append(result, newGroup)
			}
			currentGroup = newGroup
		case tokenTypeGroupEnd:
			if l := len(groupStack); l > 0 {
				currentGroup = groupStack[l-1]
				groupStack = groupStack[0 : l-1]
			} else {
				currentGroup = nil
			}
		}
	}
	return result, nil
}

func parseExpressionTokens(expr string) (tokens []*parsedToken, groupStarts int, err error) {
	err = nil
	groupStarts = 0
	if len(strings.Trim(expr, " \t\n")) == 0 {
		return
	}
	runes := []rune(expr)
	tokens = make([]*parsedToken, 0, len(runes))
	max := len(runes) - 1
	inQuote := false
	var quoteChar = ' '
	groupEnds := 0
	groupsOpen := 0
	currentToken := &parsedToken{}
	startNewToken := func(i int, t parseTokenType) *parsedToken {
		currentToken = &parsedToken{tokenType: t, start: i, end: i}
		tokens = append(tokens, currentToken)
		return currentToken
	}
	startNewToken(0, tokenTypeStart)
	for i := 0; i <= max; i++ {
		ch := runes[i]
		if inQuote {
			if ch == quoteChar {
				inQuote = false
				currentToken.end = i + 1
				startNewToken(i, tokenTypeWhitespace)
			}
		} else {
			currentToken.end = i
			switch ch {
			case '"', '\'':
				inQuote = true
				quoteChar = ch
				startNewToken(i, tokenTypeName)
			case ' ', '\t', '\n':
				if currentToken.tokenType != tokenTypeWhitespace && currentToken.tokenType != tokenTypeStart {
					startNewToken(i, tokenTypeWhitespace)
				}
			case '!':
				startNewToken(i, tokenTypeNot)
			case '&':
				if i < max && runes[i+1] == '&' {
					startNewToken(i, tokenTypeOperator)
					currentToken.operator = And
					i++
					currentToken.end = i + 1
				} else {
					err = fmt.Errorf("invalid operator character '%s' - expected '&&' (at position %d)", string(ch), i)
					return
				}
			case '|':
				if i < max && runes[i+1] == '|' {
					startNewToken(i, tokenTypeOperator)
					currentToken.operator = Or
					i++
					currentToken.end = i + 1
				} else {
					err = fmt.Errorf("invaid operator character '%s' - expected '||' (at position %d)", string(ch), i)
					return
				}
			case '(':
				startNewToken(i, tokenTypeGroupStart)
				groupStarts++
				groupsOpen++
			case ')':
				startNewToken(i, tokenTypeGroupEnd)
				groupEnds++
				groupsOpen--
				if groupsOpen < 0 {
					err = fmt.Errorf("unexpected group close character '%s' (at position %d)", string(ch), i)
					return
				}
			default:
				if ch < 32 || ch > 127 {
					err = fmt.Errorf("unexpected non-naming character code %v (at position %d) - use enclosing quotes if necessary", ch, i)
					return
				} else if !unicode.Is(allowedNameChars, ch) {
					err = fmt.Errorf("unexpected non-naming character '%s' (at position %d) - use enclosing quotes if necessary", string(ch), i)
					return
				}
				if currentToken.tokenType != tokenTypeName {
					startNewToken(i, tokenTypeName)
				}
			}
		}
	}
	if groupStarts != groupEnds {
		err = fmt.Errorf("unbalanced grouping parentheses (at position %d)", max)
		return
	}
	if inQuote {
		err = fmt.Errorf("unclosed quote (started at position %d)", currentToken.start)
		return
	}
	currentToken.end = max + 1
	if currentToken.tokenType == tokenTypeWhitespace {
		currentToken.tokenType = tokenTypeEnd
	} else {
		startNewToken(max, tokenTypeEnd)
	}
	// check all the tokens are correctly sequenced...
	max = len(tokens) - 1
	for i, token := range tokens {
		if i < max && !tokenAllowedFollowedBy[token.tokenType][tokens[i+1].tokenType] {
			err = fmt.Errorf("unexpected character '%s' (at position %d)", string(runes[tokens[i+1].start]), tokens[i+1].start)
			return
		}
		switch token.tokenType {
		case tokenTypeNot:
			found := false
			for j := i - 1; j >= 0; j-- {
				if tokens[j].tokenType == tokenTypeOperator || tokens[j].tokenType == tokenTypeGroupStart || tokens[j].tokenType == tokenTypeStart {
					found = true
					break
				} else if tokens[j].tokenType != tokenTypeWhitespace {
					break
				}
			}
			if !found {
				err = fmt.Errorf("unexpected not operator (at position %d)", token.start)
				return
			}
		case tokenTypeName:
			found := false
			for j := i - 1; j >= 0; j-- {
				if tokens[j].tokenType == tokenTypeOperator || tokens[j].tokenType == tokenTypeGroupStart || tokens[j].tokenType == tokenTypeStart {
					found = true
					break
				} else if tokens[j].tokenType != tokenTypeWhitespace && tokens[j].tokenType != tokenTypeNot {
					break
				}
			}
			if !found {
				err = fmt.Errorf("unexpected property name start (at position %d)", token.start)
				return
			}
			token.name = string(runes[token.start:token.end])
			if isQuotedStr(token.name, true) {
				token.name = token.name[1 : len(token.name)-1]
			}
		}
	}
	return
}

var allowedNameChars = &unicode.RangeTable{
	R16: []unicode.Range16{
		{'$', '$', 1},
		{'-', '9', 1},
		{'@', 'Z', 1},
		{'_', '_', 1},
		{'a', 'z', 1},
		{'~', '~', 1},
	},
}

type parseTokenType int

const (
	tokenTypeStart parseTokenType = iota
	tokenTypeWhitespace
	tokenTypeOperator
	tokenTypeNot
	tokenTypeName
	tokenTypeGroupStart
	tokenTypeGroupEnd
	tokenTypeEnd
)

var tokenAllowedFollowedBy = map[parseTokenType]map[parseTokenType]bool{
	tokenTypeStart: {
		tokenTypeGroupStart: true,
		tokenTypeNot:        true,
		tokenTypeName:       true,
	},
	tokenTypeWhitespace: {
		tokenTypeOperator:   true,
		tokenTypeNot:        true,
		tokenTypeName:       true,
		tokenTypeGroupStart: true,
		tokenTypeGroupEnd:   true,
	},
	tokenTypeOperator: {
		tokenTypeWhitespace: true,
		tokenTypeNot:        true,
		tokenTypeName:       true,
		tokenTypeGroupStart: true,
	},
	tokenTypeNot: {
		tokenTypeNot:        true,
		tokenTypeName:       true,
		tokenTypeGroupStart: true,
	},
	tokenTypeName: {
		tokenTypeWhitespace: true,
		tokenTypeOperator:   true,
		tokenTypeGroupEnd:   true,
		tokenTypeEnd:        true,
	},
	tokenTypeGroupStart: {
		tokenTypeWhitespace: true,
		tokenTypeNot:        true,
		tokenTypeGroupStart: true,
		tokenTypeName:       true,
	},
	tokenTypeGroupEnd: {
		tokenTypeWhitespace: true,
		tokenTypeGroupEnd:   true,
		tokenTypeOperator:   true,
		tokenTypeEnd:        true,
	},
	tokenTypeEnd: {},
}

type parsedToken struct {
	tokenType parseTokenType
	start     int
	end       int
	operator  BooleanOperator
	name      string
}

// OthersExpr is a list of expressions (OtherProperty, OtherGrouping) that can be evaluated against
// an object to determine the presence or non-presence of specific named properties
//
// This is used by the PropertyValidator.RequiredWith and PropertyValidator.UnwantedWith fields
type OthersExpr []Other

// Other is the interface for items in OthersExpr - and is implemented by OtherProperty, OtherGrouping and
// by OthersExpr itself
type Other interface {
	// Evaluate evaluates the presence or non-presence of named properties in a given object
	//
	// the currentObj arg is the object that the property should be checked for within
	//
	// the ancestryValues arg provides the ancestry of objects in case there is a need to traverse
	// upwards.  The first, index 0, item in the ancestryValues slice will be the parent of the
	// currentObj... the second item will be the grandparent of currentObj... etc.
	Evaluate(currentObj map[string]interface{}, ancestryValues []interface{}, vcx *ValidatorContext) bool
	// GetOperator returns the boolean operator (And / Or)
	GetOperator() BooleanOperator
	// String method provides a string representation of the expression
	String() string
}

// Evaluate implements Other.Evaluate
func (o *OthersExpr) Evaluate(currentObj map[string]interface{}, ancestryValues []interface{}, vcx *ValidatorContext) bool {
	if len(*o) == 0 {
		return true
	}
	result := false
	for i, item := range *o {
		itemResult := item.Evaluate(currentObj, ancestryValues, vcx)
		if i == 0 {
			result = itemResult
		} else if item.GetOperator() == Or {
			result = result || itemResult
		} else {
			result = result && itemResult
		}
	}
	return result
}

func (o *OthersExpr) String() string {
	var sb strings.Builder
	for i, other := range *o {
		if i > 0 {
			if other.GetOperator() == Or {
				sb.WriteString(" || ")
			} else {
				sb.WriteString(" && ")
			}
		}
		sb.WriteString(other.String())
	}
	return sb.String()
}

func (o *OthersExpr) AddProperty(name string) *OthersExpr {
	*o = append(*o, &OtherProperty{Name: name})
	return o
}

func (o *OthersExpr) AddNotProperty(name string) *OthersExpr {
	*o = append(*o, &OtherProperty{Name: name, Not: true})
	return o
}

func (o *OthersExpr) AddAndProperty(name string) *OthersExpr {
	*o = append(*o, &OtherProperty{Name: name, Op: And})
	return o
}

func (o *OthersExpr) AddAndNotProperty(name string) *OthersExpr {
	*o = append(*o, &OtherProperty{Name: name, Op: And, Not: true})
	return o
}

func (o *OthersExpr) AddOrProperty(name string) *OthersExpr {
	*o = append(*o, &OtherProperty{Name: name, Op: Or})
	return o
}

func (o *OthersExpr) AddOrNotProperty(name string) *OthersExpr {
	*o = append(*o, &OtherProperty{Name: name, Op: Or, Not: true})
	return o
}

func (o *OthersExpr) AddGroup(of *OthersExpr) *OthersExpr {
	*o = append(*o, &OtherGrouping{Of: *of, Op: And})
	return o
}

func (o *OthersExpr) AddNotGroup(of *OthersExpr) *OthersExpr {
	*o = append(*o, &OtherGrouping{Of: *of, Op: And, Not: true})
	return o
}

func (o *OthersExpr) AddAndGroup(of *OthersExpr) *OthersExpr {
	*o = append(*o, &OtherGrouping{Of: *of, Op: And})
	return o
}

func (o *OthersExpr) AddAndNotGroup(of *OthersExpr) *OthersExpr {
	*o = append(*o, &OtherGrouping{Of: *of, Op: And, Not: true})
	return o
}

func (o *OthersExpr) AddOrGroup(of *OthersExpr) *OthersExpr {
	*o = append(*o, &OtherGrouping{Of: *of, Op: Or})
	return o
}

func (o *OthersExpr) AddOrNotGroup(of *OthersExpr) *OthersExpr {
	*o = append(*o, &OtherGrouping{Of: *of, Op: Or, Not: true})
	return o
}

// GetOperator implements Other.GetOperator (always returns And)
func (o *OthersExpr) GetOperator() BooleanOperator {
	// collection always returns And (it's irrelevant because it's never used)
	return And
}

type OtherProperty struct {
	// Name is the name of the property whose presence or non-presence is to be checked
	Name string
	// Not is whether presence is NOTed (!) - i.e. if Not is set to true, then the non-presence is checked
	Not bool
	// Op is the boolean operator (And / Or) applied to the previous resultant
	Op BooleanOperator
	// internals...
	cachedName     string // check against Name to see if it's changed
	normalizedName string
	pathed         bool // whether the name is a pathing at all
	upPath         int
	downPath       []string
}

func NewOtherProperty(name string) *OtherProperty {
	return &OtherProperty{
		Name: name,
		Op:   And,
	}
}

func (p *OtherProperty) NOTed() *OtherProperty {
	p.Not = !p.Not
	return p
}

func (p *OtherProperty) ANDed() *OtherProperty {
	p.Op = And
	return p
}

func (p *OtherProperty) ORed() *OtherProperty {
	p.Op = Or
	return p
}

// Evaluate implements Other.Evaluate
func (p *OtherProperty) Evaluate(currentObj map[string]interface{}, ancestryValues []interface{}, vcx *ValidatorContext) bool {
	if strings.HasPrefix(p.Name, "~") {
		// names starting with ~ (tilde) are condition checks...
		return vcx.IsCondition(p.Name[1:])
	}
	p.checkChanged()
	r := false
	if p.pathed {
		r = p.walkPath(currentObj, ancestryValues)
	} else if p.normalizedName == "" {
		r = currentObj != nil
	} else {
		_, r = currentObj[p.normalizedName]
	}
	return !p.Not == r
}

func (p *OtherProperty) checkChanged() {
	if p.Name == p.cachedName {
		// no changes
		return
	}
	// reset...
	p.cachedName = p.Name
	p.normalizedName = p.Name
	p.pathed = false
	p.upPath = 0
	p.downPath = nil
	// short circuit out of any obvious non-pathing dots...
	if p.normalizedName == "." || p.normalizedName == "" {
		// empty or just dot is ignored
		p.normalizedName = ""
		return
	} else if strings.HasPrefix(p.normalizedName, ".") && !strings.Contains(p.normalizedName[1:], ".") {
		// starts with just one dot...
		p.normalizedName = p.normalizedName[1:]
		return
	} else if !strings.Contains(p.normalizedName, ".") {
		return
	}
	// now count and splice the dots...
	pastFirstDots := false
	prevCh := ' '
	lastDot := -1
	for i, ch := range p.cachedName {
		if ch == '.' && prevCh != '\\' {
			if i == 1 && prevCh == '/' {
				p.upPath = -1
			} else if pastFirstDots {
				if pathPart := p.cachedName[lastDot+1 : i]; pathPart != "" {
					p.downPath = append(p.downPath, unescapeDots(pathPart))
				}
			} else if p.upPath >= 0 {
				p.upPath++
			}
			lastDot = i
		} else {
			pastFirstDots = true
		}
		prevCh = ch
	}
	if lastDot < len(p.cachedName)-1 {
		p.downPath = append(p.downPath, unescapeDots(p.cachedName[lastDot+1:]))
	}
	p.pathed = p.upPath != 0 || len(p.downPath) > 0
}

var escapedDotsRegexp = regexp.MustCompile(`\\\.`)

func unescapeDots(str string) string {
	return escapedDotsRegexp.ReplaceAllString(str, ".")
}

func (p *OtherProperty) walkPath(currentObj map[string]interface{}, ancestryValues []interface{}) bool {
	result := false
	if p.upPath == -1 {
		// from root...
		if len(ancestryValues) > 0 {
			lastAncestor := ancestryValues[len(ancestryValues)-1]
			if m, ok := lastAncestor.(map[string]interface{}); ok {
				result = walkPathFrom(m, p.downPath, 0)
			}
		}
	} else if p.upPath != 0 {
		if len(ancestryValues) > 0 && (p.upPath-1) < len(ancestryValues) {
			ancestor := ancestryValues[p.upPath-1]
			if m, ok := ancestor.(map[string]interface{}); ok {
				result = walkPathFrom(m, p.downPath, 0)
			}
		}
	} else {
		result = walkPathFrom(currentObj, p.downPath, 0)
	}
	return result
}

func walkPathFrom(currentObj map[string]interface{}, downPath []string, downIndex int) bool {
	ptyName := downPath[downIndex]
	if v, ok := currentObj[ptyName]; ok {
		if downIndex+1 == len(downPath) {
			return ok
		} else if next, ok := v.(map[string]interface{}); ok {
			return walkPathFrom(next, downPath, downIndex+1)
		}
	}
	return false
}

// GetOperator implements Other.GetOperator
func (p *OtherProperty) GetOperator() BooleanOperator {
	return p.Op
}

func (p *OtherProperty) String() string {
	str := ternary(p.Not).string("!", "")
	needsQuotes := false
	for _, ch := range p.Name {
		if !unicode.Is(allowedNameChars, ch) {
			needsQuotes = true
			break
		}
	}
	if needsQuotes {
		useQuote := ternary(strings.Contains(p.Name, "'")).string("\"", "'")
		str = useQuote + p.Name + useQuote
	} else {
		str = str + p.Name
	}
	return str
}

type OtherGrouping struct {
	// Of is the items within the grouping
	Of OthersExpr
	// Not is whether the grouping is NOTed (!)
	Not bool
	// Op is the boolean operator (And / Or) applied to the previous resultant
	Op BooleanOperator
}

// Evaluate implements Other.Evaluate
func (g *OtherGrouping) Evaluate(currentObj map[string]interface{}, ancestryValues []interface{}, vcx *ValidatorContext) bool {
	return !g.Not == g.Of.Evaluate(currentObj, ancestryValues, vcx)
}

// GetOperator implements Other.GetOperator (always returns And)
func (g *OtherGrouping) GetOperator() BooleanOperator {
	return g.Op
}

func (g *OtherGrouping) String() string {
	return ternary(g.Not).string("!", "") + "(" + g.Of.String() + ")"
}

func NewOtherGrouping(items ...interface{}) *OtherGrouping {
	result := &OtherGrouping{}
	var prevOp BooleanOperator = -1
	for i, item := range items {
		switch v := item.(type) {
		case BooleanOperator:
			prevOp = v
			if i == 0 {
				result.Op = prevOp
				prevOp = -1
			}
		case string:
			if prevOp != -1 {
				result.Of = append(result.Of, &OtherProperty{Name: v, Op: prevOp})
			} else {
				result.Of = append(result.Of, &OtherProperty{Name: v, Op: And})
			}
			prevOp = -1
		case OtherProperty:
			result.Of = append(result.Of, &v)
		case *OtherProperty:
			result.Of = append(result.Of, v)
		case OtherGrouping:
			result.Of = append(result.Of, &v)
		case *OtherGrouping:
			result.Of = append(result.Of, v)
		default:
			panic("Illegal argument")
		}
	}
	return result
}

func (g *OtherGrouping) NOTed() *OtherGrouping {
	g.Not = !g.Not
	return g
}

func (g *OtherGrouping) ANDed() *OtherGrouping {
	g.Op = And
	return g
}

func (g *OtherGrouping) ORed() *OtherGrouping {
	g.Op = Or
	return g
}

func (g *OtherGrouping) AddProperty(name string) *OtherGrouping {
	g.Of.AddProperty(name)
	return g
}

func (g *OtherGrouping) AddNotProperty(name string) *OtherGrouping {
	g.Of.AddNotProperty(name)
	return g
}

func (g *OtherGrouping) AddAndProperty(name string) *OtherGrouping {
	g.Of.AddAndProperty(name)
	return g
}

func (g *OtherGrouping) AddAndNotProperty(name string) *OtherGrouping {
	g.Of.AddAndNotProperty(name)
	return g
}

func (g *OtherGrouping) AddOrProperty(name string) *OtherGrouping {
	g.Of.AddOrProperty(name)
	return g
}

func (g *OtherGrouping) AddOrNotProperty(name string) *OtherGrouping {
	g.Of.AddOrNotProperty(name)
	return g
}

func (g *OtherGrouping) AddGroup(of *OthersExpr) *OtherGrouping {
	g.Of.AddGroup(of)
	return g
}

func (g *OtherGrouping) AddNotGroup(of *OthersExpr) *OtherGrouping {
	g.Of.AddNotGroup(of)
	return g
}

func (g *OtherGrouping) AddAndGroup(of *OthersExpr) *OtherGrouping {
	g.Of.AddAndGroup(of)
	return g
}

func (g *OtherGrouping) AddAndNotGroup(of *OthersExpr) *OtherGrouping {
	g.Of.AddAndNotGroup(of)
	return g
}

func (g *OtherGrouping) AddOrGroup(of *OthersExpr) *OtherGrouping {
	g.Of.AddOrGroup(of)
	return g
}

func (g *OtherGrouping) AddOrNotGroup(of *OthersExpr) *OtherGrouping {
	g.Of.AddOrNotGroup(of)
	return g
}
