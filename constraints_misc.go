package valix

import (
	"fmt"
	"golang.org/x/text/language"
	"net"
	"net/mail"
	"strconv"
	"strings"
)

// StringValidCardNumber constraint checks that a string contains a valid card number according
// to Luhn Algorithm and checking that card number is 10 to 19 digits
type StringValidCardNumber struct {
	// if set to true, AllowSpaces accepts space separators in the card number (but must appear between each 4 digits)
	AllowSpaces bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringValidCardNumber) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringValidCardNumber) checkString(str string, vcx *ValidatorContext) bool {
	const digitMinChar = '0'
	const digitMaxChar = '9'
	buffer := []byte(str)
	l := len(buffer)
	if c.AllowSpaces {
		stripBuffer := make([]byte, 0, len(str))
		for i, by := range buffer {
			if by == ' ' {
				if (i+1)%5 != 0 || i+1 == l {
					return false
				}
			} else {
				stripBuffer = append(stripBuffer, by)
			}
		}
		buffer = stripBuffer
		l = len(buffer)
	}
	if l < 10 || l > 19 {
		return false
	}
	var checkSum uint8 = 0
	doubling := false
	for i := l - 1; i >= 0; i-- {
		ch := buffer[i]
		if (ch < digitMinChar) || (ch > digitMaxChar) {
			return false
		}
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

// GetMessage implements the Constraint.GetMessage
func (c *StringValidCardNumber) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidCardNumber)
}

// StringValidCountryCode constraint checks that a string is a valid ISO-3166 (3166-1 / 3166-2) country code
type StringValidCountryCode struct {
	Allow3166_2                         bool
	Allow3166_2_Obsoletes               bool
	AllowUserAssigned                   bool
	Allow3166_1_ExceptionallyReserved   bool
	Allow3166_1_IndeterminatelyReserved bool
	Allow3166_1_TransitionallyReserved  bool
	Allow3166_1_Deleted                 bool
	Allow3166_1_Numeric                 bool
	// overrides all other flags (with the exception of AllowUserAssigned) and allows only ISO-3166-1 numeric codes
	NumericOnly bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringValidCountryCode) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if c.NumericOnly {
		if c.checkNumericOnly(v) {
			return true, ""
		}
	} else {
		if c.checkAll(v) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

func (c *StringValidCountryCode) checkAll(v interface{}) bool {
	if str, ok := v.(string); ok {
		return c.checkStringOnly(str)
	} else if c.Allow3166_1_Numeric {
		if iv, ok, isNumber := coerceToInt(v); ok && isNumber {
			return c.checkNumeric(iv)
		}
	}
	return false
}

func (c *StringValidCountryCode) checkStringOnly(str string) bool {
	if strings.Contains(str, "-") {
		return c.isOk31662(str)
	} else if _, ok := iSO3166_2_CountryCodes[str]; ok {
		return true
	} else if c.Allow3166_1_Numeric && iSO3166_1_NumericCodes[str] {
		return true
	} else if len(str) == 2 {
		return c.isOk31661(str)
	} else if c.Allow3166_1_Numeric && c.AllowUserAssigned {
		if iv, err := strconv.Atoi(str); err == nil && iv >= 900 && iv <= 999 {
			return true
		}
	}
	return false
}

func (c *StringValidCountryCode) isOk31661(str string) bool {
	if subs, ok := iso3166_1_CountryCodesMatrix[str[0]]; ok {
		if assignment, ok := subs[str[1]]; ok {
			if assignment == ccAs || assignment == ccRa ||
				(c.AllowUserAssigned && assignment == ccUA) ||
				(c.Allow3166_1_ExceptionallyReserved && assignment == ccER) ||
				(c.Allow3166_1_IndeterminatelyReserved && assignment == ccIR) ||
				(c.Allow3166_1_TransitionallyReserved && assignment == ccTR) ||
				(c.Allow3166_1_Deleted && assignment == ccDl) {
				return true
			}
		}
	}
	return false
}

func (c *StringValidCountryCode) isOk31662(str string) bool {
	if c.Allow3166_2 {
		if parts := strings.Split(str, "-"); len(parts) == 2 {
			if rs, ok := iSO3166_2_CountryCodes[parts[0]]; ok {
				if rs[parts[1]] {
					return true
				}
			}
			if c.Allow3166_2_Obsoletes && iSO3166_2_ObsoleteCodes[str] {
				return true
			}
		}
	}
	return false
}

func (c *StringValidCountryCode) checkNumericOnly(v interface{}) bool {
	if str, ok := v.(string); ok {
		if iSO3166_1_NumericCodes[str] {
			return true
		} else if c.AllowUserAssigned {
			if iv, err := strconv.Atoi(str); err == nil && iv >= 900 && iv <= 999 {
				return true
			}
		}
	} else if iv, ok, isNumber := coerceToInt(v); ok && isNumber {
		return c.checkNumeric(iv)
	}
	return false
}

func (c *StringValidCountryCode) checkNumeric(iv int64) bool {
	if c.AllowUserAssigned && iv >= 900 && iv <= 999 {
		return true
	}
	ic := fmt.Sprintf("%03d", iv)
	if iSO3166_1_NumericCodes[ic] {
		return true
	}
	return false
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidCountryCode) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidCountryCode)
}

// StringValidCurrencyCode constraint checks that a string is a valid ISO-4217 currency code
type StringValidCurrencyCode struct {
	AllowNumeric    bool
	AllowHistorical bool
	AllowUnofficial bool
	AllowCrypto     bool
	// AllowTestCode when set to true, allows test currency codes (i.e. "XTS" or numeric "963")
	AllowTestCode bool
	// AllowNoCode when set to true, allows no code (i.e. "XXX" or numeric "999")
	AllowNoCode bool
	// set to true to only allow ISO-4217 numeric currency codes
	NumericOnly bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringValidCurrencyCode) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if c.NumericOnly {
		if c.checkNumericOnly(v) {
			return true, ""
		}
	} else if c.checkEither(v) {
		return true, ""
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

func (c *StringValidCurrencyCode) checkEither(v interface{}) bool {
	if str, ok := v.(string); ok {
		if c.isOkTest(str) || c.isOkNoCode(str) || c.isOkRegular(str) ||
			c.isOkHistorical(str) || c.isOkOther(str) {
			return true
		}
	} else if c.AllowNumeric {
		if iv, ok, isNumber := coerceToInt(v); ok && isNumber {
			return c.checkNumeric(iv)
		}
	}
	return false
}

func (c *StringValidCurrencyCode) isOkRegular(str string) bool {
	return iSO4217CurrencyCodes[str] || (c.AllowNumeric && iSO4217CurrencyCodesNumeric[str])
}

func (c *StringValidCurrencyCode) isOkTest(str string) bool {
	return c.AllowTestCode && (str == ISO4217TestCurrencyCode || (c.AllowNumeric && str == ISO4217TestCurrencyCodeNumeric))
}

func (c *StringValidCurrencyCode) isOkNoCode(str string) bool {
	return c.AllowNoCode && (str == ISO4217NoCurrencyCode || (c.AllowNumeric && str == ISO4217NoCurrencyCodeNumeric))
}

func (c *StringValidCurrencyCode) isOkHistorical(str string) bool {
	return c.AllowHistorical && (iSO4217CurrencyCodesHistorical[str] || (c.AllowNumeric && iSO4217CurrencyCodesNumericHistorical[str]))
}

func (c *StringValidCurrencyCode) isOkOther(str string) bool {
	return (c.AllowUnofficial && unofficialCurrencyCodes[str]) || (c.AllowCrypto && cryptoCurrencyCodes[str])
}

func (c *StringValidCurrencyCode) checkNumericOnly(v interface{}) bool {
	if str, ok := v.(string); ok {
		if (c.AllowTestCode && str == ISO4217TestCurrencyCodeNumeric) ||
			(c.AllowNoCode && str == ISO4217NoCurrencyCodeNumeric) ||
			iSO4217CurrencyCodesNumeric[str] ||
			(c.AllowHistorical && iSO4217CurrencyCodesNumericHistorical[str]) {
			return true
		}
	} else if iv, ok, isNumber := coerceToInt(v); ok && isNumber {
		return c.checkNumeric(iv)
	}
	return false
}

func (c *StringValidCurrencyCode) checkNumeric(iv int64) bool {
	ic := fmt.Sprintf("%03d", iv)
	if (c.AllowTestCode && ic == ISO4217TestCurrencyCodeNumeric) ||
		(c.AllowNoCode && ic == ISO4217NoCurrencyCodeNumeric) ||
		iSO4217CurrencyCodesNumeric[ic] ||
		(c.AllowHistorical && iSO4217CurrencyCodesNumericHistorical[ic]) {
		return true
	}
	return false
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidCurrencyCode) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidCurrencyCode)
}

// StringValidEmail constraint checks that a string contains a valid email address (does not
// verify the email address!)
//
// NB. Uses mail.ParseAddress to check valid email address
type StringValidEmail struct {
	// DisallowRFC5322 when set, disallows email addresses in RFC5322 format (i.e "Barry Gibbs <bg@example.com>")
	DisallowRFC5322 bool
	// CheckExchange when set, checks the MX (mail exchange) for the email address
	CheckExchange bool
	// AllowIPAddress when set, allows email addresses with IP (e.g. "me@[123.0.1.2]")
	AllowIPAddress bool
	// AllowIPV6 when set, allows email addresses with IP v6 (e.g. "me@[2001:db8::68]")
	AllowIPV6 bool
	// AllowLocal when set, allows email addresses with 'local' (e.g. "me@localhost", "me@local", "me@localdomain", "me@[127.0.0.1]", "me@[::1]")
	AllowLocal bool
	// AllowTldOnly when set, allows email addresses with only Tld specified (e.g. "me@audi")
	AllowTldOnly bool
	// AllowGeographicTlds when set, allows email addresses with geographic Tlds (e.g. "me@some-company.africa")
	AllowGeographicTlds bool
	// AllowGenericTlds when set, allows email addresses with generic Tlds (e.g. "me@some.academy")
	AllowGenericTlds bool
	// AllowBrandTlds when set, allows email addresses with brand Tlds (e.g. "me@my.audi")
	AllowBrandTlds bool
	// AllowInfraTlds when set, allows email addresses with infrastructure Tlds (e.g. "me@arpa")
	AllowInfraTlds bool
	// AllowTestTlds when set, allows email addresses with test Tlds and test domains (e.g. "me@example.com", "me@test.com")
	AllowTestTlds bool
	// AddCountryCodeTlds is an optional slice of additional country (and geographic) Tlds to allow
	AddCountryCodeTlds []string
	// ExcCountryCodeTlds is an optional slice of country (and geographic) Tlds to disallow
	ExcCountryCodeTlds []string
	// AddGenericTlds is an optional slice of additional generic Tlds to allow (only checked if AllowGenericTlds is also set)
	AddGenericTlds []string
	// ExcGenericTlds is an optional slice of generic Tlds to disallow (only relevant if AllowGenericTlds is also set)
	ExcGenericTlds []string
	// AddBrandTlds is an optional slice of additional brand Tlds to allow (only checked if AllowBrandTlds is also set)
	AddBrandTlds []string
	// ExcBrandTlds is an optional slice of brand Tlds to disallow (only relevant if AllowBrandTlds is also set)
	ExcBrandTlds []string
	// AddLocalTlds is an optional slice of additional local Tlds to allow (only checked if AllowLocal is also set)
	AddLocalTlds []string
	// ExcLocalTlds is an optional slice of local Tlds to disallow (only relevant if AllowLocal is also set)
	ExcLocalTlds []string
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringValidEmail) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringValidEmail) checkString(str string, vcx *ValidatorContext) bool {
	result := false
	if c.DisallowRFC5322 {
		result = c.checkEmailAddress(str)
	} else if a, err := mail.ParseAddress(str); err != nil {
		// fails to parse addresses with IPv6 - so try directly...
		if c.AllowIPAddress {
			result = c.checkEmailAddress(str)
		}
	} else {
		result = c.checkEmailAddress(a.Address)
	}
	return result
}

func (c *StringValidEmail) checkEmailAddress(addr string) bool {
	if isValidEmail(addr, domainOptions{
		allowIPAddress:      c.AllowIPAddress,
		allowIPV6:           c.AllowIPV6,
		allowLocal:          c.AllowLocal,
		allowTldOnly:        c.AllowTldOnly,
		allowGeographicTlds: c.AllowGeographicTlds,
		allowGenericTlds:    c.AllowGenericTlds,
		allowBrandTlds:      c.AllowBrandTlds,
		allowInfraTlds:      c.AllowInfraTlds,
		allowTestTlds:       c.AllowTestTlds,
		addCountryCodeTlds:  c.AddCountryCodeTlds,
		excCountryCodeTlds:  c.ExcCountryCodeTlds,
		addGenericTlds:      c.AddGenericTlds,
		excGenericTlds:      c.ExcGenericTlds,
		addBrandTlds:        c.AddBrandTlds,
		excBrandTlds:        c.ExcBrandTlds,
		addLocalTlds:        c.AddLocalTlds,
		excLocalTlds:        c.ExcLocalTlds,
	}) {
		if c.CheckExchange {
			if aAt := strings.LastIndex(addr, "@"); aAt != -1 {
				dom := addr[aAt+1:]
				if mxs, err := net.LookupMX(dom); err == nil && len(mxs) > 0 {
					return true
				}
			}
		} else {
			return true
		}
	}
	return false
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidEmail) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidEmail)
}

// StringValidLanguageCode constraint checks that a string is a valid BCP-47 language code
//
// NB. Uses language.Parse to check valid language code
type StringValidLanguageCode struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringValidLanguageCode) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringValidLanguageCode) checkString(str string, vcx *ValidatorContext) bool {
	if _, err := language.Parse(str); err != nil {
		return false
	}
	return true
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidLanguageCode) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidLanguageCode)
}

// StringValidUuid constraint to check that a string value is a valid UUID
type StringValidUuid struct {
	// the minimum UUID version (optional - if zero this is not checked)
	MinVersion uint8
	// the specific UUID version (optional - if zero this is not checked)
	SpecificVersion uint8
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringValidUuid) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringValidUuid) checkString(str string, vcx *ValidatorContext) bool {
	if !uuidRegexp.MatchString(str) {
		return false
	}
	var version = str[14] - 48
	if (c.MinVersion > 0 && version < c.MinVersion) || (c.SpecificVersion > 0 && version != c.SpecificVersion) {
		return false
	}
	return true
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidUuid) GetMessage(tcx I18nContext) string {
	if c.SpecificVersion > 0 {
		return defaultMessage(tcx, c.Message, fmtMsgUuidCorrectVer, c.SpecificVersion)
	} else if c.MinVersion > 0 {
		return defaultMessage(tcx, c.Message, fmtMsgUuidMinVersion, c.MinVersion)
	}
	return defaultMessage(tcx, c.Message, msgValidUuid)
}
