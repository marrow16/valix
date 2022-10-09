package valix

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Preset is the interface used for presets (any preset registered directly using RegisterPreset must implement this interface)
type Preset interface {
	Check(v string) bool
	GetRegexp() *regexp.Regexp
	GetPostChecker() PostPatternChecker
	GetMessage() string
}

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
	presetsRegistry.register(token, &patternPreset{
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

// RegisterPreset registers a preset pattern for use by the StringPresetPattern constraint
//
// * the `asConstraint` arg, if true, means the preset is also registered as a named constraint and is available for use
// as a `v8n` constraint tag
func RegisterPreset(token string, preset Preset, asConstraint bool) {
	presetsRegistry.register(token, preset)
	if asConstraint {
		constraintsRegistry.registerNamed(true, token, &StringPresetPattern{
			Preset: token,
		})
	}
}

func GetRegisteredPreset(token string) (Preset, bool) {
	return presetsRegistry.get(token)
}

const (
	numericPattern = "^[-+]?[0-9]*(?:\\.[0-9]+)?$"
	// allow number with scientific notation...
	numericWithScientific = "^(([-+]?[0-9]*(?:\\.[0-9]+)?)|(([+-]?\\d*)\\.?\\d+[eE][-+]?\\d+))$"
	// allow number with scientific notation - plus allow +-Inf or Nan...
	numericFull = "^(([-+]?[0-9]*(?:\\.[0-9]+)?)|(([+-]?\\d*)\\.?\\d+[eE][-+]?\\d+)|([-+]?[Ii][Nn][Ff])|([Nn][Aa][Nn]))$"
	// match anything
	anyMatchPattern = ".*"
	// cmyk...
	cmykNumPattern = "(?:100(?:\\.[0]+)?\\%)|(?:1(?:\\.[0]+)?)|(?:0(?:\\.[0-9]+)?)|(?:(?:[0-9]|[1-9][0-9])(?:\\.[0-9]+)?\\%)"
	cmykPattern    = "^cmyk\\(\\s*(" + cmykNumPattern + ")\\s*,\\s*(" + cmykNumPattern + ")\\s*,\\s*(" + cmykNumPattern + ")\\s*,\\s*(" + cmykNumPattern + ")\\s*\\)$"
	// rgb-icc...
	rgbIccPattern = "^rgb-icc\\((.*?)\\)"
)

var (
	spaceHyphenStripper = regexp.MustCompile(`[\s-]+`)
	matchAnything       = regexp.MustCompile(anyMatchPattern)
	cmykRegexp          = regexp.MustCompile(cmykPattern)
	rgbIccRegexp        = regexp.MustCompile(rgbIccPattern)
	rgbIccOkStrRegexp   = regexp.MustCompile("^[A-Za-z]([A-Za-z0-9_\\-]{1,31})?$")
	ean8Regexp          = regexp.MustCompile("^[0-9]{8}$")
	ean13Regexp         = regexp.MustCompile("^[0-9]{13}$")
	ean14Regexp         = regexp.MustCompile("^(([0-9]{14})|(\\(01\\)[0-9]{14}))$")
	ean18Regexp         = regexp.MustCompile("^(([0-9]{18})|(\\(00\\)[0-9]{18}))$")
	ean99Regexp         = regexp.MustCompile("^99[0-9]{11}$")
	isbn10Regexp        = regexp.MustCompile("^(?:[0-9]{9}X|[0-9]{10})$")
	isbn13Regexp        = regexp.MustCompile("^97[89][0-9]{10}$")
	issn8Regexp         = regexp.MustCompile("^(?:[0-9]{7}X|[0-9]{8})$")
	issn13Regexp        = regexp.MustCompile("^977[0-9]{10}$")
	upcARegexp          = regexp.MustCompile("^[0-9]{12}$")
	upcERegexp          = regexp.MustCompile("^0[0-9]{7}$")
)

type presetRegistry struct {
	namedPresets map[string]Preset
	sync         *sync.Mutex
}

var presetsRegistry presetRegistry

func init() {
	presetsRegistry = presetRegistry{
		namedPresets: getBuiltInPresets(),
		sync:         &sync.Mutex{},
	}
}

func (r *presetRegistry) register(token string, preset Preset) {
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

func (r *presetRegistry) get(token string) (Preset, bool) {
	defer r.sync.Unlock()
	r.sync.Lock()
	if p, ok := r.namedPresets[token]; ok {
		return p, true
	}
	return nil, false
}

type patternPreset struct {
	regex       *regexp.Regexp
	postChecker PostPatternChecker
	msg         string
}

func (pp *patternPreset) Check(v string) bool {
	result := pp.regex.MatchString(v)
	if result && pp.postChecker != nil {
		result = pp.postChecker.Check(v)
	}
	return result
}

func (pp *patternPreset) GetRegexp() *regexp.Regexp {
	return pp.regex
}

func (pp *patternPreset) GetPostChecker() PostPatternChecker {
	return pp.postChecker
}

func (pp *patternPreset) GetMessage() string {
	return pp.msg
}

func getBuiltInPresets() map[string]Preset {
	return map[string]Preset{
		PresetAlpha: &patternPreset{
			regex: regexp.MustCompile("^[a-zA-Z]+$"),
			msg:   msgPresetAlpha,
		},
		PresetAlphaNumeric: &patternPreset{
			regex: regexp.MustCompile("^[a-zA-Z0-9]+$"),
			msg:   msgPresetAlphaNumeric,
		},
		PresetCMYK: &patternPreset{
			regex:       cmykRegexp,
			postChecker: cmyk{400},
			msg:         msgPresetCMYK,
		},
		PresetCMYK300: &patternPreset{
			regex:       cmykRegexp,
			postChecker: cmyk{300},
			msg:         msgPresetCMYK300,
		},
		PresetBarcode: &patternPreset{
			regex:       regexp.MustCompile("^[0-9X()]{8,22}$"),
			postChecker: barcode{},
			msg:         msgPresetBarcode,
		},
		PresetBase64: &patternPreset{
			regex: regexp.MustCompile("^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})$"),
			msg:   msgPresetBase64,
		},
		PresetBase64URL: &patternPreset{
			regex: regexp.MustCompile("^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3}=|[A-Za-z0-9-_]{4})$"),
			msg:   msgPresetBase64URL,
		},
		PresetCard: &patternPreset{
			regex:       regexp.MustCompile("^(([0-9]{10,19})|([0-9]{4} [0-9]{4} [0-9]{2,4})|([0-9]{4} [0-9]{4} [0-9]{4} [0-9]{1,4})|([0-9]{4} [0-9]{4} [0-9]{4} [0-9]{4} [0-9]{1,3}))$"),
			postChecker: cardNumber{},
			msg:         msgValidCardNumber,
		},
		PresetE164: &patternPreset{
			regex: regexp.MustCompile("^\\+[1-9]?[0-9]{7,14}$"),
			msg:   msgPresetE164,
		},
		PresetEAN: &patternPreset{
			regex:       regexp.MustCompile("^(([0-9]{8,18})|(\\(01\\)[0-9]{14})|(\\(00\\)[0-9]{18}))$"),
			postChecker: &ean{},
			msg:         msgPresetEAN,
		},
		PresetEAN8: &patternPreset{
			regex:       ean8Regexp,
			postChecker: ean8,
			msg:         msgPresetEAN8,
		},
		PresetEAN13: &patternPreset{
			regex:       ean13Regexp,
			postChecker: ean13,
			msg:         msgPresetEAN13,
		},
		PresetDUN14: &patternPreset{
			regex:       regexp.MustCompile("^(([0-9]{14})|(\\(01\\)[0-9]{14}))$"),
			postChecker: ean14,
			msg:         msgPresetDUN14,
		},
		PresetEAN14: &patternPreset{
			regex:       ean14Regexp,
			postChecker: ean14,
			msg:         msgPresetEAN14,
		},
		PresetEAN18: &patternPreset{
			regex:       ean18Regexp,
			postChecker: ean18,
			msg:         msgPresetEAN18,
		},
		PresetEAN99: &patternPreset{
			regex:       ean99Regexp,
			postChecker: ean13,
			msg:         msgPresetEAN99,
		},
		PresetHexadecimal: &patternPreset{
			regex: regexp.MustCompile("^(0[xX])?[0-9a-fA-F]+$"),
			msg:   msgPresetHexadecimal,
		},
		PresetHsl: &patternPreset{
			regex: regexp.MustCompile("^hsl\\(\\s*(?:0|[1-9]\\d?|[12]\\d\\d|3[0-5]\\d|360)\\s*,\\s*(?:0|[1-9]\\d?|100)%\\s*,\\s*(?:0|[1-9]\\d?|100)%\\s*\\)$"),
			msg:   msgPresetHsl,
		},
		PresetHsla: &patternPreset{
			regex: regexp.MustCompile("^hsla\\(\\s*(?:0|[1-9]\\d?|[12]\\d\\d|3[0-5]\\d|360)\\s*,\\s*(?:0|[1-9]\\d?|100)%\\s*,\\s*(?:0|[1-9]\\d?|100)%\\s*,\\s*(?:0.[1-9]*|[01])\\s*\\)$"),
			msg:   msgPresetHsla,
		},
		PresetHtmlColor: &patternPreset{
			regex: regexp.MustCompile("^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$"),
			msg:   msgPresetHtmlColor,
		},
		PresetInteger: &patternPreset{
			regex: regexp.MustCompile("^[0-9]+$"),
			msg:   msgPresetInteger,
		},
		PresetISBN: &patternPreset{
			regex:       regexp.MustCompile("^((97[89][0-9]{10})|(?:[0-9]{9}X|[0-9]{10}))$"),
			postChecker: isbn{},
			msg:         msgPresetISBN,
		},
		PresetISBN10: &patternPreset{
			regex:       isbn10Regexp,
			postChecker: isbn10,
			msg:         msgPresetISBN10,
		},
		PresetISBN13: &patternPreset{
			regex:       isbn13Regexp,
			postChecker: isbn13,
			msg:         msgPresetISBN13,
		},
		PresetISSN: &patternPreset{
			regex:       regexp.MustCompile("^((977[0-9]{10})|(?:[0-9]{7}X|[0-9]{8}))$"),
			postChecker: issn{},
			msg:         msgPresetISSN,
		},
		PresetISSN8: &patternPreset{
			regex:       issn8Regexp,
			postChecker: issn8,
			msg:         msgPresetISSN8,
		},
		PresetISSN13: &patternPreset{
			regex:       issn13Regexp,
			postChecker: isbn13,
			msg:         msgPresetISSN13,
		},
		PresetNumeric: &patternPreset{
			regex: regexp.MustCompile(numericPattern),
			msg:   msgPresetNumeric,
		},
		PresetNumericE: &patternPreset{
			regex: regexp.MustCompile(numericWithScientific),
			msg:   msgPresetNumeric,
		},
		PresetNumericX: &patternPreset{
			regex: regexp.MustCompile(numericFull),
			msg:   msgPresetNumeric,
		},
		PresetPublication: &patternPreset{
			regex:       regexp.MustCompile("^((97[789][0-9]{10})|(?:[0-9]{9}X|[0-9]{10})|(?:[0-9]{7}X|[0-9]{8}))$"),
			postChecker: publication{},
			msg:         msgPresetPublication,
		},
		PresetRgb: &patternPreset{
			regex: regexp.MustCompile("^rgb\\(\\s*(?:(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])|(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%)\\s*\\)$"),
			msg:   msgPresetRgb,
		},
		PresetRgba: &patternPreset{
			regex: regexp.MustCompile("^rgba\\(\\s*(?:(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])|(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(?:0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%)\\s*,\\s*(?:0.[1-9]*|[01])\\s*\\)$"),
			msg:   msgPresetRgba,
		},
		PresetRgbIcc: &patternPreset{
			regex:       rgbIccRegexp,
			postChecker: rgbIcc{},
			msg:         msgPresetRgbIcc,
		},
		PresetULID: &patternPreset{
			regex: regexp.MustCompile("^[01234567][0123456789ABCDEFGHJKMNPQRSTVWXYZ]{25}$"),
			msg:   msgPresetULID,
		},
		PresetUPC: &patternPreset{
			regex:       regexp.MustCompile("^(([0-9]{12})|(0[0-9]{7}))$"),
			postChecker: upc{},
			msg:         msgPresetUPC,
		},
		PresetUPCA: &patternPreset{
			regex:       upcARegexp,
			postChecker: upcA,
			msg:         msgPresetUPCA,
		},
		PresetUPCE: &patternPreset{
			regex:       upcERegexp,
			postChecker: upcE,
			msg:         msgPresetUPCE,
		},
		PresetUuid: &patternPreset{
			regex: regexp.MustCompile("^([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})$"),
			msg:   msgValidUuid,
		},
		PresetUUID: &patternPreset{
			regex: regexp.MustCompile("^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})$"),
			msg:   msgValidUuid,
		},
		PresetUuid1: &patternPreset{
			regex: regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-1[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"),
			msg:   msgPresetUuid1,
		},
		PresetUUID1: &patternPreset{
			regex: regexp.MustCompile("^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-1[0-9A-Fa-f]{3}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}$"),
			msg:   msgPresetUuid1,
		},
		PresetUuid2: &patternPreset{
			regex: regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-2[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"),
			msg:   msgPresetUuid2,
		},
		PresetUUID2: &patternPreset{
			regex: regexp.MustCompile("^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-2[0-9A-Fa-f]{3}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}$"),
			msg:   msgPresetUuid2,
		},
		PresetUuid3: &patternPreset{
			regex: regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"),
			msg:   msgPresetUuid3,
		},
		PresetUUID3: &patternPreset{
			regex: regexp.MustCompile("^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-3[0-9A-Fa-f]{3}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}$"),
			msg:   msgPresetUuid3,
		},
		PresetUuid4: &patternPreset{
			regex: regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"),
			msg:   msgPresetUuid4,
		},
		PresetUUID4: &patternPreset{
			regex: regexp.MustCompile("^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-4[0-9A-Fa-f]{3}-[89abAB][0-9A-Fa-f]{3}-[0-9A-Fa-f]{12}$"),
			msg:   msgPresetUuid4,
		},
		PresetUuid5: &patternPreset{
			regex: regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"),
			msg:   msgPresetUuid5,
		},
		PresetUUID5: &patternPreset{
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
	stripPrefix   string
}

func (mc moduloCheck) Check(str string) bool {
	// assume our regex will have not matched with spaces, hyphens, braces etc. - but some may decide to allow them...
	useStr := spaceHyphenStripper.ReplaceAllString(str, "")
	if mc.stripPrefix != "" && strings.HasPrefix(useStr, mc.stripPrefix) && len(useStr)-len(mc.stripPrefix) == len(mc.weights)+1 {
		useStr = useStr[len(mc.stripPrefix):]
	}
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
	ean8 = moduloCheck{
		modulo:  10,
		weights: []int32{1, 3, 1, 3, 1, 3, 1},
	}
	ean13 = moduloCheck{
		modulo:  10,
		weights: []int32{1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3},
	}
	ean14 = moduloCheck{
		modulo:      10,
		weights:     []int32{1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1},
		stripPrefix: "(01)",
	}
	ean18 = moduloCheck{
		modulo:      10,
		weights:     []int32{1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1},
		stripPrefix: "(00)",
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

type barcode struct{}

func (ck barcode) Check(str string) bool {
	switch len(str) {
	case 8:
		return check8DigitBarcode(str)
	case 10:
		return check10DigitBarcode(str)
	case 12:
		return check12DigitBarcode(str)
	case 13:
		return check13DigitBarcode(str)
	case 14:
		return check14DigitBarcode(str)
	case 18:
		return check18DigitBarcode(str)
	case 22:
		return check22DigitBarcode(str)
	}
	return false
}

func check8DigitBarcode(str string) bool {
	result := false
	if ean8Regexp.MatchString(str) {
		result = ean8.Check(str) || issn8.Check(str) || upcE.Check(str)
	} else if issn8Regexp.MatchString(str) {
		result = issn8.Check(str)
	}
	return result
}

func check10DigitBarcode(str string) bool {
	return isbn10Regexp.MatchString(str) && isbn10.Check(str)
}

func check12DigitBarcode(str string) bool {
	return upcARegexp.MatchString(str) && upcA.Check(str)
}

func check13DigitBarcode(str string) bool {
	return ean13Regexp.MatchString(str) && ean13.Check(str)
}

func check14DigitBarcode(str string) bool {
	return ean14Regexp.MatchString(str) && ean14.Check(str)
}

func check18DigitBarcode(str string) bool {
	if strings.HasPrefix(str, "(") {
		return ean14Regexp.MatchString(str) && ean14.Check(str)
	}
	return ean18Regexp.MatchString(str) && ean18.Check(str)
}

func check22DigitBarcode(str string) bool {
	return ean18Regexp.MatchString(str) && ean18.Check(str)
}

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

type ean struct{}

func (ck ean) Check(str string) bool {
	switch len(str) {
	case 8:
		return ean8.Check(str)
	case 13:
		return ean13.Check(str)
	case 14:
		return ean14.Check(str)
	case 18:
		if strings.HasPrefix(str, "(") {
			return ean14.Check(str)
		}
		return ean18.Check(str)
	case 22:
		return ean18.Check(str)
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

type cmyk struct {
	max float64
}

func (ck cmyk) Check(str string) bool {
	parts := cmykRegexp.FindStringSubmatch(str)
	result := false
	if len(parts) == 5 {
		sum := 0.0
		result = true
		for _, part := range parts[1:] {
			if strings.HasSuffix(part, "%") {
				f, err := strconv.ParseFloat(part[:len(part)-1], 64)
				result = result && err == nil && f <= 100.0
				sum += f
			} else {
				f, err := strconv.ParseFloat(part, 64)
				result = result && err == nil && f <= 1.0
				sum += f * 100.0
			}
		}
		result = result && sum <= ck.max
	}
	return result
}

type rgbIcc struct{}

func (ck rgbIcc) Check(str string) bool {
	if subs := rgbIccRegexp.FindStringSubmatch(str); subs != nil && len(subs) == 2 {
		if args, ok := rgbIccArgsParse(strings.Trim(subs[1], " ")); ok && len(args) > 0 {
			tokenAt := 0
			token := args[0]
			if !isRgbIccToken(token) {
				// if the first isn't the token - then we expect the 4th to be and the first 3 to be rgb values
				if len(args) < 4 {
					return false
				} else if token = args[3]; !isRgbIccToken(token) {
					return false
				}
				tokenAt = 3
				for i := 0; i < 3; i++ {
					if val, err := strconv.ParseInt(args[i], 10, 0); err != nil || val < 0 || val > 255 {
						return false
					}
				}
			}
			return isValidRgbIccToken(token, args, tokenAt)
		}
	}
	return false
}

func isValidRgbIccToken(token string, args []string, tokenAt int) (result bool) {
	result = false
	switch token {
	case "#CMYK":
		result = isValidRgbIccCmykToken(args, tokenAt)
	case "#Grayscale":
		result = isValidRgbIccGrayscaleToken(args, tokenAt)
	case "#Separation":
		result = isValidRgbIccSeparationToken(args, tokenAt)
	case "#Registration":
		result = isValidRgbIccRegistrationToken(args, tokenAt)
	case "#SpotColor":
		result = isValidRgbIccSpotColorToken(args, tokenAt)
	}
	return
}

func isValidRgbIccCmykToken(args []string, tokenAt int) bool {
	// we expect there to be 4 more args with cmyk values...
	if len(args) != tokenAt+5 {
		return false
	}
	for i := 1; i < 5; i++ {
		if !isRgbIccValue(args[tokenAt+i]) {
			return false
		}
	}
	return true
}

func isValidRgbIccGrayscaleToken(args []string, tokenAt int) bool {
	// we expect there to be 1 more arg with a 0.0-1.0 value...
	if len(args) != tokenAt+2 {
		return false
	} else if !isRgbIccValue(args[tokenAt+1]) {
		return false
	}
	return true
}

func isValidRgbIccSeparationToken(args []string, tokenAt int) bool {
	// we expect there to be 1 more arg with a string value...
	if len(args) != tokenAt+2 {
		return false
	} else if !isRgbIccStringOk(args[tokenAt+1], true) {
		return false
	}
	return true
}

func isValidRgbIccRegistrationToken(args []string, tokenAt int) bool {
	// we expect there to be 0 or 1 more args - and if there's 1, it should be a 0.0-1.0 value...
	if len(args) == tokenAt+2 {
		if !isRgbIccValue(args[tokenAt+1]) {
			return false
		}
	} else if len(args) > tokenAt+2 {
		return false
	}
	return true
}

func isValidRgbIccSpotColorToken(args []string, tokenAt int) bool {
	// we expect there to be at least 2 more args after this - the 1st being 'stringy' and the 2nd being a 0.0-1.0 value...
	if len(args) < tokenAt+3 {
		return false
	} else if !isRgbIccStringOk(args[tokenAt+1], false) {
		return false
	} else if !isRgbIccValue(args[tokenAt+2]) {
		return false
	}
	if len(args) > tokenAt+3 {
		// followed by another #token - #CMYK, #Grayscale or #Registration..
		tokenAt = tokenAt + 3
		token := args[tokenAt]
		if !isRgbIccToken(token) {
			return false
		}
		followingTokenOk := false
		switch token {
		case "#CMYK":
			followingTokenOk = isValidRgbIccCmykToken(args, tokenAt)
		case "#Grayscale":
			followingTokenOk = isValidRgbIccGrayscaleToken(args, tokenAt)
		case "#Registration":
			followingTokenOk = isValidRgbIccRegistrationToken(args, tokenAt)
		}
		return followingTokenOk
	}
	return true
}

func isRgbIccToken(str string) bool {
	return str[0:1] == "#" &&
		(str == "#CMYK" || str == "#Grayscale" || str == "#Separation" || str == "#SpotColor" || str == "#Registration")
}

func isRgbIccStringOk(str string, quotedOnly bool) bool {
	if strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'") {
		return len(str) > 2
	} else if !quotedOnly {
		return rgbIccOkStrRegexp.MatchString(str)
	}
	return false
}

func isRgbIccValue(str string) bool {
	if val, err := strconv.ParseFloat(str, 64); err != nil || math.IsNaN(val) || math.IsInf(val, 0) ||
		val < 0.0 || val > 1.0 {
		return false
	}
	return true
}

func rgbIccArgsParse(str string) ([]string, bool) {
	if len(str) == 0 {
		return nil, false
	}
	result := make([]string, 0, len(str)/3)
	runes := []rune(str)
	lastArgAt := 0
	inQuote := false
	for i, r := range runes {
		switch r {
		case ',':
			if !inQuote {
				arg := strings.Trim(string(runes[lastArgAt:i]), " ")
				if !rgbIccArgCheck(arg) {
					return nil, false
				}
				result = append(result, arg)
				lastArgAt = i + 1
			}
		case '\'':
			inQuote = !inQuote
		}
	}
	if lastArgAt < len(runes) {
		arg := strings.Trim(string(runes[lastArgAt:]), " ")
		if !rgbIccArgCheck(arg) {
			return nil, false
		}
		result = append(result, arg)
	}
	return result, true
}

func rgbIccArgCheck(arg string) bool {
	if arg == "" {
		return false
	} else if strings.HasPrefix(arg, "'") && strings.HasSuffix(arg, "'") {
		return len(arg) > 2 && !strings.Contains(arg[1:len(arg)-1], "'")
	}
	return !strings.Contains(arg, " ")
}

const (
	msgPresetAlpha        = "Value must be only alphabet characters (A-Z, a-z)"
	msgPresetAlphaNumeric = "Value must be only alphanumeric characters (A-Z, a-z, 0-9)"
	msgPresetBarcode      = "Value must be a valid barcode"
	msgPresetISBN         = "Value must be a valid ISBN"
	msgPresetISBN10       = "Value must be a valid ISBN-10"
	msgPresetISBN13       = "Value must be a valid ISBN-13"
	msgPresetISSN         = "Value must be a valid ISSN"
	msgPresetISSN8        = "Value must be a valid ISSN-8"
	msgPresetISSN13       = "Value must be a valid ISSN-13"
	msgPresetEAN          = "Value must be a valid EAN code"
	msgPresetEAN8         = "Value must be a valid EAN-8 code"
	msgPresetEAN13        = "Value must be a valid EAN-13 code"
	msgPresetDUN14        = "Value must be a valid DUN-14 code"
	msgPresetEAN14        = "Value must be a valid EAN-14 code"
	msgPresetEAN18        = "Value must be a valid EAN-18 code"
	msgPresetEAN99        = "Value must be a valid EAN-99 code"
	msgPresetUPC          = "Value must be a valid UPC code (UPC-A or UPC-E)"
	msgPresetUPCA         = "Value must be a valid UPC-A code"
	msgPresetUPCE         = "Value must be a valid UPC-E code"
	msgPresetPublication  = "Value must be a valid ISBN or ISSN"
	msgPresetNumeric      = "Value must be a valid number string"
	msgPresetInteger      = "Value must be a valid integer string (characters 0-9)"
	msgPresetHexadecimal  = "Value must be a valid hexadecimal string"
	msgPresetCMYK         = "Value must be a valid cmyk() colour string"
	msgPresetCMYK300      = "Value must be a valid cmyk() colour string (maximum 300%)"
	msgPresetHtmlColor    = "Value must be a valid HTML colour string"
	msgPresetRgb          = "Value must be a valid rgb() colour string"
	msgPresetRgba         = "Value must be a valid rgba() colour string"
	msgPresetRgbIcc       = "Value must be a valid rgb-icc() colour string"
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
	// PresetAlpha checks for only alpha characters (use with StringPresetPattern)
	PresetAlpha = "alpha"
	// PresetAlphaNumeric checks for only alphanumeric characters (use with StringPresetPattern)
	PresetAlphaNumeric = "alphaNumeric"
	// PresetBarcode checks for a valid barcode (EAN, ISBN, ISSN, UPC) (use with StringPresetPattern)
	PresetBarcode = "barcode"
	// PresetBase64 checks for valid base64 encoded string (use with StringPresetPattern)
	PresetBase64 = "base64"
	// PresetBase64URL checks for valid base64 URL encoded string (use with StringPresetPattern)
	PresetBase64URL = "base64URL"
	// PresetCMYK checks for valid cmyk() color string (use with StringPresetPattern)
	PresetCMYK = "cmyk"
	// PresetCMYK300 checks for valid cmyk() color string (maximum 300%) (use with StringPresetPattern)
	PresetCMYK300 = "cmyk300"
	// PresetCard checks for valid card number (use with StringPresetPattern)
	PresetCard = "card"
	// PresetDUN14 checks for valid DUN-14 barcode (use with StringPresetPattern)
	PresetDUN14 = "DUN14"
	// PresetE164 checks for valid E.164 code (use with StringPresetPattern)
	PresetE164 = "e164"
	// PresetEAN checks for valid EAN barcode (EAN-8, 13, 14, 18 or 99) (use with StringPresetPattern)
	PresetEAN = "EAN"
	// PresetEAN13 checks for valid EAN-13 barcode (use with StringPresetPattern)
	PresetEAN13 = "EAN13"
	// PresetEAN14 checks for valid EAN-14 barcode (use with StringPresetPattern)
	PresetEAN14 = "EAN14"
	// PresetEAN18 checks for valid EAN-18 barcode (use with StringPresetPattern)
	PresetEAN18 = "EAN18"
	// PresetEAN8 checks for valid EAN-8 barcode (use with StringPresetPattern)
	PresetEAN8 = "EAN8"
	// PresetEAN99 checks for valid EAN-99 barcode (use with StringPresetPattern)
	PresetEAN99 = "EAN99"
	// PresetHexadecimal checks for valid hexadecimal string (use with StringPresetPattern)
	PresetHexadecimal = "hexadecimal"
	// PresetHsl checks for valid hsl() color string (use with StringPresetPattern)
	PresetHsl = "hsl"
	// PresetHsla checks for valid hsla() color string (use with StringPresetPattern)
	PresetHsla = "hsla"
	// PresetHtmlColor checks for valid HTML color string (use with StringPresetPattern)
	PresetHtmlColor = "htmlColor"
	// PresetISBN checks for valid ISBN barcode (ISBN-10 or 13) (use with StringPresetPattern)
	PresetISBN = "ISBN"
	// PresetISBN10 checks for valid ISBN-10 barcode (use with StringPresetPattern)
	PresetISBN10 = "ISBN10"
	// PresetISBN13 checks for valid ISBN-13 barcode (use with StringPresetPattern)
	PresetISBN13 = "ISBN13"
	// PresetISSN checks for valid ISSN barcode (ISSN-8 or 13) (use with StringPresetPattern)
	PresetISSN = "ISSN"
	// PresetISSN13 checks for valid ISSN-13 barcode (use with StringPresetPattern)
	PresetISSN13 = "ISSN13"
	// PresetISSN8 checks for valid ISSN-8 barcode (use with StringPresetPattern)
	PresetISSN8 = "ISSN8"
	// PresetInteger checks for valid integer string (chars 0-9) (use with StringPresetPattern)
	PresetInteger = "integer"
	// PresetNumeric checks for valid numeric string (use with StringPresetPattern)
	PresetNumeric = "numeric"
	// PresetNumericE checks for valid numeric string allowing scientific notation (use with StringPresetPattern)
	PresetNumericE = "numeric+e"
	// PresetNumericX checks for valid numeric string allowing scientific notation plus "Inf" and "NaN" (use with StringPresetPattern)
	PresetNumericX = "numeric+x"
	// PresetPublication checks for valid publication barcode (ISBN or ISSN) (use with StringPresetPattern)
	PresetPublication = "publication"
	// PresetRgb checks for valid rgb() color string (use with StringPresetPattern)
	PresetRgb = "rgb"
	// PresetRgbIcc checks for valid rgb-icc() color string (use with StringPresetPattern)
	PresetRgbIcc = "rgb-icc"
	// PresetRgba checks for valid rgba() color string (use with StringPresetPattern)
	PresetRgba = "rgba"
	// PresetULID checks for valid ULID code (use with StringPresetPattern)
	PresetULID = "ULID"
	// PresetUPC checks for valid UPC barcode (UPC-A or UPC-E) (use with StringPresetPattern)
	PresetUPC = "UPC"
	// PresetUPCA checks for valid UPC-A barcode (use with StringPresetPattern)
	PresetUPCA = "UPC-A"
	// PresetUPCE checks for valid UPC-E barcode (use with StringPresetPattern)
	PresetUPCE = "UPC-E"
	// PresetUUID checks for valid UUID (upper or lower hex chars) (use with StringPresetPattern)
	PresetUUID = "UUID"
	// PresetUUID1 checks for valid UUID Version 1 (upper or lower hex chars) (use with StringPresetPattern)
	PresetUUID1 = "UUID1"
	// PresetUUID2 checks for valid UUID Version 2 (upper or lower hex chars) (use with StringPresetPattern)
	PresetUUID2 = "UUID2"
	// PresetUUID3 checks for valid UUID Version 3 (upper or lower hex chars) (use with StringPresetPattern)
	PresetUUID3 = "UUID3"
	// PresetUUID4 checks for valid UUID Version 4 (upper or lower hex chars) (use with StringPresetPattern)
	PresetUUID4 = "UUID4"
	// PresetUUID5 checks for valid UUID Version 5 (upper or lower hex chars) (use with StringPresetPattern)
	PresetUUID5 = "UUID5"
	// PresetUuid checks for valid UUID (lower hex chars only) (use with StringPresetPattern)
	PresetUuid = "uuid"
	// PresetUuid1 checks for valid UUID Version 1 (lower hex chars only) (use with StringPresetPattern)
	PresetUuid1 = "uuid1"
	// PresetUuid2 checks for valid UUID Version 2 (lower hex chars only) (use with StringPresetPattern)
	PresetUuid2 = "uuid2"
	// PresetUuid3 checks for valid UUID Version 3 (lower hex chars only) (use with StringPresetPattern)
	PresetUuid3 = "uuid3"
	// PresetUuid4 checks for valid UUID Version 4 (lower hex chars only) (use with StringPresetPattern)
	PresetUuid4 = "uuid4"
	// PresetUuid5 checks for valid UUID Version 5 (lower hex chars only) (use with StringPresetPattern)
	PresetUuid5 = "uuid5"
)
