package valix

import (
	"regexp"
	"sync"
)

// RegisterPresetPattern registers a preset pattern for use by the StringPresetPattern constraint
//
// * the `token` arg specifies the token for the preset (as used by the StringPresetPattern.Preset field)
//
// * the `rx` arg specifies the *regexp.Regexp that must be matched (if nil, a match anything regexp is used)
//
// * the `message` arg specifies the message for the preset
//
// * the `postCheck` is any post pattern checking that needs to be satisfied
//
// * the `asConstraint` arg, if true, means the preset is also registered as a named constraint and is available for use
// as a `v8n` constraint tag
func RegisterPresetPattern(token string, rx *regexp.Regexp, message string, postCheck PostPatternChecker, asConstraint bool) {
	useRx := rx
	if rx == nil {
		useRx = matchAnything
	}
	presetsRegistry.register(token, patternPreset{
		regex:       useRx,
		postChecker: postCheck,
		msg:         message,
	})
	if asConstraint {
		constraintsRegistry.registerNamed(true, token, &StringPresetPattern{
			Preset: token,
		})
	}
}

const (
	numericPattern = "^[-+]?[0-9]*(?:\\.[0-9]+)?$"
	// allow number with scientific notation...
	numericWithScientific = "^(([-+]?[0-9]*(?:\\.[0-9]+)?)|(([+-]?\\d*)\\.?\\d+[eE][-+]?\\d+))$"
	// allow number with scientific notation - plus allow +-Inf or Nan...
	numericFull = "^(([-+]?[0-9]*(?:\\.[0-9]+)?)|(([+-]?\\d*)\\.?\\d+[eE][-+]?\\d+)|([-+]?[Ii][Nn][Ff])|([Nn][Aa][Nn]))$"
	// match anything
	anyMatchPattern = ".*"
)

var (
	spaceHyphenStripper = regexp.MustCompile(`[\s-]+`)
	matchAnything       = regexp.MustCompile(anyMatchPattern)
)

type presetRegistry struct {
	namedPresets map[string]patternPreset
	sync         *sync.Mutex
}

var presetsRegistry presetRegistry

func init() {
	presetsRegistry = presetRegistry{
		namedPresets: getBuiltInPresets(),
		sync:         &sync.Mutex{},
	}
}

func (r *presetRegistry) register(token string, preset patternPreset) {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.namedPresets[token] = preset
}

// reset for testing
func (r *presetRegistry) reset() {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.namedPresets = getBuiltInPresets()
}

func (r *presetRegistry) get(token string) (*patternPreset, bool) {
	defer r.sync.Unlock()
	r.sync.Lock()
	if p, ok := r.namedPresets[token]; ok {
		return &p, true
	}
	return nil, false
}

type patternPreset struct {
	regex       *regexp.Regexp
	postChecker PostPatternChecker
	msg         string
}

func (pp *patternPreset) check(v string) bool {
	result := pp.regex.MatchString(v)
	if result && pp.postChecker != nil {
		result = pp.postChecker.Check(v)
	}
	return result
}

func getBuiltInPresets() map[string]patternPreset {
	return map[string]patternPreset{
		presetTokenAlpha: {
			regex: regexp.MustCompile("^[a-zA-Z]+$"),
			msg:   msgPresetAlpha,
		},
		presetTokenAlphaNumeric: {
			regex: regexp.MustCompile("^[a-zA-Z0-9]+$"),
			msg:   msgPresetAlphaNumeric,
		},
		presetTokenBase64: {
			regex: regexp.MustCompile("^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})$"),
			msg:   msgPresetBase64,
		},
		presetTokenBase64URL: {
			regex: regexp.MustCompile("^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3}=|[A-Za-z0-9-_]{4})$"),
			msg:   msgPresetBase64URL,
		},
		presetTokenCard: {
			regex:       regexp.MustCompile("^(([0-9]{12,19})|([0-9]{4} [0-9]{4} [0-9]{4})|([0-9]{4} [0-9]{4} [0-9]{4} [0-9]{1,4})|([0-9]{4} [0-9]{4} [0-9]{4} [0-9]{4} [0-9]{1,3}))$"),
			postChecker: cardNumber{},
			msg:         msgValidCardNumber,
		},
		presetTokenE164: {
			regex: regexp.MustCompile("^\\+[1-9]?[0-9]{7,14}$"),
			msg:   msgPresetE164,
		},
		presetTokenEAN13: {
			regex:       regexp.MustCompile("^[0-9]{13}$"),
			postChecker: ean13,
			msg:         msgPresetEAN13,
		},
		presetTokenHexadecimal: {
			regex: regexp.MustCompile("^(0[xX])?[0-9a-fA-F]+$"),
			msg:   msgPresetHexadecimal,
		},
		presetTokenHsl: {
			regex: regexp.MustCompile("^hsl\\(\\s*(?:0|[1-9]\\d?|[12]\\d\\d|3[0-5]\\d|360)\\s*,\\s*(?:0|[1-9]\\d?|100)%\\s*,\\s*(?:0|[1-9]\\d?|100)%\\s*\\)$"),
			msg:   msgPresetHsl,
		},
		presetTokenHsla: {
			regex: regexp.MustCompile("^hsla\\(\\s*(?:0|[1-9]\\d?|[12]\\d\\d|3[0-5]\\d|360)\\s*,\\s*(?:0|[1-9]\\d?|100)%\\s*,\\s*(?:0|[1-9]\\d?|100)%\\s*,\\s*(?:0.[1-9]*|[01])\\s*\\)$"),
			msg:   msgPresetHsla,
		},
		presetTokenHtmlColor: {
			regex: regexp.MustCompile("^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$"),
			msg:   msgPresetHtmlColor,
		},
		presetTokenInteger: {
			regex: regexp.MustCompile("^[0-9]+$"),
			msg:   msgPresetInteger,
		},
		presetTokenISBN: {
			regex:       regexp.MustCompile("^((97[89][0-9]{10})|(?:[0-9]{9}X|[0-9]{10}))$"),
			postChecker: isbn{},
			msg:         msgPresetISBN,
		},
		presetTokenISBN10: {
			regex:       regexp.MustCompile("^(?:[0-9]{9}X|[0-9]{10})$"),
			postChecker: isbn10,
			msg:         msgPresetISBN10,
		},
		presetTokenISBN13: {
			regex:       regexp.MustCompile("^97[89][0-9]{10}$"),
			postChecker: isbn13,
			msg:         msgPresetISBN13,
		},
		presetTokenISSN: {
			regex:       regexp.MustCompile("^((977[0-9]{10})|(?:[0-9]{7}X|[0-9]{8}))$"),
			postChecker: issn{},
			msg:         msgPresetISSN,
		},
		presetTokenISSN8: {
			regex:       regexp.MustCompile("^(?:[0-9]{7}X|[0-9]{8})$"),
			postChecker: issn8,
			msg:         msgPresetISSN8,
		},
		presetTokenISSN13: {
			regex:       regexp.MustCompile("^977[0-9]{10}$"),
			postChecker: isbn13,
			msg:         msgPresetISSN13,
		},
		presetTokenNumeric: {
			regex: regexp.MustCompile(numericPattern),
			msg:   msgPresetNumeric,
		},
		presetTokenNumericE: {
			regex: regexp.MustCompile(numericWithScientific),
			msg:   msgPresetNumeric,
		},
		presetTokenNumericX: {
			regex: regexp.MustCompile(numericFull),
			msg:   msgPresetNumeric,
		},
		presetTokenPublication: {
			regex:       regexp.MustCompile("^((97[789][0-9]{10})|(?:[0-9]{9}X|[0-9]{10})|(?:[0-9]{7}X|[0-9]{8}))$"),
			postChecker: publication{},
			msg:         msgPresetPublication,
		},
		presetTokenRgb: {
			regex: regexp.MustCompile("^rgb\\(\\s*(?:(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])|(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%)\\s*\\)$"),
			msg:   msgPresetRgb,
		},
		presetTokenRgba: {
			regex: regexp.MustCompile("^rgba\\(\\s*(?:(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])|(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%)\\s*,\\s*(?:0.[1-9]*|[01])\\s*\\)$"),
			msg:   msgPresetRgba,
		},
		presetTokenULID: {
			regex: regexp.MustCompile("^[01234567][0123456789ABCDEFGHJKMNPQRSTVWXYZ]{25}$"),
			msg:   msgPresetULID,
		},
		presetTokenUPC: {
			regex:       regexp.MustCompile("^(([0-9]{12})|(0[0-9]{7}))$"),
			postChecker: upc{},
			msg:         msgPresetUPC,
		},
		presetTokenUPCA: {
			regex:       regexp.MustCompile("^[0-9]{12}$"),
			postChecker: upcA,
			msg:         msgPresetUPCA,
		},
		presetTokenUPCE: {
			regex:       regexp.MustCompile("^0[0-9]{7}$"),
			postChecker: upcE,
			msg:         msgPresetUPCE,
		},
		presetTokenUuid: {
			regex: regexp.MustCompile("^([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})$"),
			msg:   msgValidUuid,
		},
		presetTokenUUID: {
			regex: regexp.MustCompile("^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})$"),
			msg:   msgValidUuid,
		},
		presetTokenUuid1: {
			regex: regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-1[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"),
			msg:   msgPresetUuid1,
		},
		presetTokenUUID1: {
			regex: regexp.MustCompile("^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-1[0-9A-Fa-f]{3}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}$"),
			msg:   msgPresetUuid1,
		},
		presetTokenUuid2: {
			regex: regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-2[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"),
			msg:   msgPresetUuid2,
		},
		presetTokenUUID2: {
			regex: regexp.MustCompile("^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-2[0-9A-Fa-f]{3}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}$"),
			msg:   msgPresetUuid2,
		},
		presetTokenUuid3: {
			regex: regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"),
			msg:   msgPresetUuid3,
		},
		presetTokenUUID3: {
			regex: regexp.MustCompile("^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-3[0-9A-Fa-f]{3}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}$"),
			msg:   msgPresetUuid3,
		},
		presetTokenUuid4: {
			regex: regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"),
			msg:   msgPresetUuid4,
		},
		presetTokenUUID4: {
			regex: regexp.MustCompile("^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-4[0-9A-Fa-f]{3}-[89abAB][0-9A-Fa-f]{3}-[0-9A-Fa-f]{12}$"),
			msg:   msgPresetUuid4,
		},
		presetTokenUuid5: {
			regex: regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"),
			msg:   msgPresetUuid5,
		},
		presetTokenUUID5: {
			regex: regexp.MustCompile("^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-5[0-9A-Fa-f]{3}-[89abAB][0-9A-Fa-f]{3}-[0-9A-Fa-f]{12}$"),
			msg:   msgPresetUuid5,
		},
	}
}

type PostPatternChecker interface {
	Check(str string) bool
}

type moduloCheck struct {
	modulo        int32
	weights       []int32
	overflowDigit string
}

func (mc moduloCheck) Check(str string) bool {
	// assume our regex will have not matched with spaces and hyphens - but some may decide to allow them...
	useStr := spaceHyphenStripper.ReplaceAllString(str, "")
	if len(useStr) == len(mc.weights)+1 {
		var sum int32 = 0
		for i, ch := range useStr[0:len(mc.weights)] {
			sum += (ch - 48) * mc.weights[i]
		}
		remainder := sum % mc.modulo
		if remainder == 0 {
			return useStr[len(useStr)-1:] == "0"
		} else {
			ck := mc.modulo - remainder
			ckDigit := ""
			if ck >= 10 {
				ckDigit = mc.overflowDigit
			} else {
				ckDigit = string(ck + 48)
			}
			return useStr[len(useStr)-1:] == ckDigit
		}
	}
	return false
}

var (
	issn8 = moduloCheck{
		modulo:        11,
		weights:       []int32{8, 7, 6, 5, 4, 3, 2},
		overflowDigit: "X",
	}
	isbn10 = moduloCheck{
		modulo:        11,
		weights:       []int32{10, 9, 8, 7, 6, 5, 4, 3, 2},
		overflowDigit: "X",
	}
	isbn13 = moduloCheck{
		modulo:  10,
		weights: []int32{1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3},
	}
	ean13 = moduloCheck{
		modulo:  10,
		weights: []int32{1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3},
	}
	upcA = moduloCheck{
		modulo:  10,
		weights: []int32{3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3},
	}
	upcE = moduloCheck{
		modulo:  10,
		weights: []int32{3, 1, 3, 1, 3, 1, 3},
	}
)

type isbn struct{}

func (ck isbn) Check(str string) bool {
	switch len(str) {
	case 10:
		return isbn10.Check(str)
	case 13:
		return isbn13.Check(str)
	}
	return false
}

type issn struct{}

func (ck issn) Check(str string) bool {
	switch len(str) {
	case 8:
		return issn8.Check(str)
	case 13:
		return isbn13.Check(str)
	}
	return false
}

type upc struct{}

func (ck upc) Check(str string) bool {
	switch len(str) {
	case 8:
		return upcE.Check(str)
	case 12:
		return upcA.Check(str)
	}
	return false
}

type publication struct{}

func (ck publication) Check(str string) bool {
	switch len(str) {
	case 8:
		return issn8.Check(str)
	case 10:
		return isbn10.Check(str)
	case 13:
		return isbn13.Check(str)
	}
	return false
}

type cardNumber struct{}

func (ck cardNumber) Check(str string) bool {
	const digitMinChar = '0'
	buffer := []byte(spaceHyphenStripper.ReplaceAllString(str, ""))
	var checkSum uint8 = 0
	doubling := false
	for i := len(buffer) - 1; i >= 0; i-- {
		ch := buffer[i]
		digit := ch - digitMinChar
		if doubling && digit > 4 {
			digit = (digit * 2) - 9
		} else if doubling {
			digit = digit * 2
		}
		checkSum = checkSum + digit
		doubling = !doubling
	}
	return checkSum%10 == 0
}

const (
	msgPresetAlpha        = "Value must be only alphabet characters (A-Z, a-z)"
	msgPresetAlphaNumeric = "Value must be only alphanumeric characters (A-Z, a-z, 0-9)"
	msgPresetISBN         = "Value must be a valid ISBN"
	msgPresetISBN10       = "Value must be a valid ISBN-10"
	msgPresetISBN13       = "Value must be a valid ISBN-13"
	msgPresetISSN         = "Value must be a valid ISSN"
	msgPresetISSN8        = "Value must be a valid ISSN-8"
	msgPresetISSN13       = "Value must be a valid ISSN-13"
	msgPresetEAN13        = "Value must be a valid EAN-13 code"
	msgPresetUPC          = "Value must be a valid UPC code (UPC-A or UPC-E)"
	msgPresetUPCA         = "Value must be a valid UPC-A code"
	msgPresetUPCE         = "Value must be a valid UPC-E code"
	msgPresetPublication  = "Value must be a valid ISBN or ISSN"
	msgPresetNumeric      = "Value must be a valid number string"
	msgPresetInteger      = "Value must be a valid integer string (characters 0-9)"
	msgPresetHexadecimal  = "Value must be a valid hexadecimal string"
	msgPresetHtmlColor    = "Value must be a valid HTML colour string"
	msgPresetRgb          = "Value must be a valid rgb() colour string"
	msgPresetRgba         = "Value must be a valid rgba() colour string"
	msgPresetHsl          = "Value must be a valid hsl() colour string"
	msgPresetHsla         = "Value must be a valid hsla() colour string"
	msgPresetE164         = "Value must be a valid E.164 code"
	msgPresetBase64       = "Value must be a valid base64 encoded string"
	msgPresetBase64URL    = "Value must be a valid base64 URL encoded string"
	msgPresetUuid1        = "Value must be a valid UUID (Version 1)"
	msgPresetUuid2        = "Value must be a valid UUID (Version 2)"
	msgPresetUuid3        = "Value must be a valid UUID (Version 3)"
	msgPresetUuid4        = "Value must be a valid UUID (Version 4)"
	msgPresetUuid5        = "Value must be a valid UUID (Version 5)"
	msgPresetULID         = "Value must be a valid ULID"
)

const (
	presetTokenAlpha        = "alpha"
	presetTokenAlphaNumeric = "alphaNumeric"
	presetTokenBase64       = "base64"
	presetTokenBase64URL    = "base64URL"
	presetTokenCard         = "card"
	presetTokenE164         = "e164"
	presetTokenEAN13        = "EAN13"
	presetTokenHexadecimal  = "hexadecimal"
	presetTokenHsl          = "hsl"
	presetTokenHsla         = "hsla"
	presetTokenHtmlColor    = "htmlColor"
	presetTokenInteger      = "integer"
	presetTokenISBN         = "ISBN"
	presetTokenISBN10       = "ISBN10"
	presetTokenISBN13       = "ISBN13"
	presetTokenISSN         = "ISSN"
	presetTokenISSN8        = "ISSN8"
	presetTokenISSN13       = "ISSN13"
	presetTokenNumeric      = "numeric"
	presetTokenNumericE     = "numeric+e"
	presetTokenNumericX     = "numeric+x"
	presetTokenPublication  = "publication"
	presetTokenRgb          = "rgb"
	presetTokenRgba         = "rgba"
	presetTokenULID         = "ULID"
	presetTokenUPC          = "UPC"
	presetTokenUPCA         = "UPC-A"
	presetTokenUPCE         = "UPC-E"
	presetTokenUuid         = "uuid"
	presetTokenUUID         = "UUID"
	presetTokenUuid1        = "uuid1"
	presetTokenUUID1        = "UUID1"
	presetTokenUuid2        = "uuid2"
	presetTokenUUID2        = "UUID2"
	presetTokenUuid3        = "uuid3"
	presetTokenUUID3        = "UUID3"
	presetTokenUuid4        = "uuid4"
	presetTokenUUID4        = "UUID4"
	presetTokenUuid5        = "uuid5"
	presetTokenUUID5        = "UUID5"
)
