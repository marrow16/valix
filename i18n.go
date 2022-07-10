package valix

import (
	"encoding/json"
	"golang.org/x/text/language"
	"net/http"
	"strings"
)

var (
	DefaultI18nProvider I18n = &defaultI18nProvider{}
)

// DefaultLanguage is the default language used by the default I18nProvider
//
// Languages provided are "de", "en", "es", "fr" & "it"
var DefaultLanguage = "en"
var DefaultRegion = ""
var DefaultFallbackLanguages = map[string]string{}

// I18n interface for supporting i18n (internationalisation) in valix -
// used to provide I18nContext interfaces upon request by the Validators
type I18n interface {
	ContextFromRequest(r *http.Request) I18nContext
	DefaultContext() I18nContext
}

// I18nContext is the interface passed around during validation that provides translations of
// messages, message formats and individual word tokens
type I18nContext interface {
	TranslateMessage(msg string) string
	TranslateFormat(format string, a ...interface{}) string
	TranslateToken(token string) string
	Language() string
	Region() string
}

type defaultI18nContext struct {
	lang   string
	region string
}

func newDefaultI18nContext(lang string, region string) I18nContext {
	return &defaultI18nContext{
		lang:   defaultLanguage(lang),
		region: strings.ToUpper(region),
	}
}

func defaultLanguage(lang string) string {
	result := strings.ToLower(lang)
	if result != "en" && result != "de" && result != "es" && result != "fr" && result != "it" {
		if fb, ok := DefaultFallbackLanguages[result]; ok {
			result = strings.ToLower(fb)
		} else {
			result = strings.ToLower(DefaultLanguage)
		}
	} else {
		return result
	}
	if result != "en" && result != "de" && result != "es" && result != "fr" && result != "it" {
		result = "en"
	}
	return result
}

func (d *defaultI18nContext) TranslateMessage(msg string) string {
	return DefaultTranslator.TranslateMessage(d.lang, d.region, msg)
}

func (d *defaultI18nContext) TranslateFormat(format string, a ...interface{}) string {
	return DefaultTranslator.TranslateFormat(d.lang, d.region, format, a...)
}

func (d *defaultI18nContext) TranslateToken(token string) string {
	return DefaultTranslator.TranslateToken(d.lang, d.region, token)
}

func (d *defaultI18nContext) Language() string {
	return d.lang
}

func (d *defaultI18nContext) Region() string {
	return d.region
}

func (d *defaultI18nContext) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"tokens":   internalTokens,
		"messages": internalMessages,
		"formats":  internalFormats,
	}
	return json.Marshal(m)
}

type defaultI18nProvider struct{}

func (i *defaultI18nProvider) ContextFromRequest(r *http.Request) I18nContext {
	tags, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	useLang := DefaultLanguage
	useRegion := DefaultRegion
	if err == nil {
		for _, tag := range tags {
			baseLang, _ := tag.Base()
			if baseLang.String() == defaultLanguage(baseLang.String()) {
				useLang = baseLang.String()
				_, _, rawRgn := tag.Raw()
				if rawRgn.String() != "ZZ" {
					useRegion = rawRgn.String()
				}
				break
			}
		}
	}
	return newDefaultI18nContext(useLang, useRegion)
}

func (i *defaultI18nProvider) DefaultContext() I18nContext {
	return newDefaultI18nContext(DefaultLanguage, DefaultRegion)
}

var fallbackI18nProvider I18n = &defaultI18nProvider{}

func obtainI18nProvider() I18n {
	if DefaultI18nProvider != nil {
		return DefaultI18nProvider
	}
	return fallbackI18nProvider
}

var fallbackI18nContext I18nContext = &defaultI18nContext{lang: "en"}

// used by ValidatorContext and constraints to ensure they never try to use a nil I18nContext
func obtainI18nContext(i18nCtx I18nContext) I18nContext {
	if i18nCtx != nil {
		return i18nCtx
	}
	result := obtainI18nProvider().DefaultContext()
	if result != nil {
		return result
	}
	return fallbackI18nContext
}

// used by defaultI18nContext.MarshalJSON - to allow listing of translation reference
// also as a kind of double-entry bookkeeping to ensure everything has been translated into TranslationsMessages
var internalMessages = map[string]string{
	msgUnableToDecode:                 msgUnableToDecode,
	msgNotJsonNull:                    msgNotJsonNull,
	msgNotJsonArray:                   msgNotJsonArray,
	msgNotJsonObject:                  msgNotJsonObject,
	msgExpectedJsonArray:              msgExpectedJsonArray,
	msgExpectedJsonObject:             msgExpectedJsonObject,
	msgErrorReading:                   msgErrorReading,
	msgErrorUnmarshall:                msgErrorUnmarshall,
	msgRequestBodyEmpty:               msgRequestBodyEmpty,
	msgUnableToDecodeRequest:          msgUnableToDecodeRequest,
	msgRequestBodyNotJsonNull:         msgRequestBodyNotJsonNull,
	msgRequestBodyNotJsonArray:        msgRequestBodyNotJsonArray,
	msgRequestBodyNotJsonObject:       msgRequestBodyNotJsonObject,
	msgRequestBodyExpectedJsonArray:   msgRequestBodyExpectedJsonArray,
	msgRequestBodyExpectedJsonObject:  msgRequestBodyExpectedJsonObject,
	msgArrayElementMustBeObject:       msgArrayElementMustBeObject,
	msgArrayElementMustNotBeNull:      msgArrayElementMustNotBeNull,
	msgMissingProperty:                msgMissingProperty,
	msgUnwantedProperty:               msgUnwantedProperty,
	msgUnknownProperty:                msgUnknownProperty,
	msgInvalidProperty:                msgInvalidProperty,
	msgInvalidPropertyName:            msgInvalidPropertyName,
	msgPropertyValueMustBeObject:      msgPropertyValueMustBeObject,
	msgPropertyRequiredWhen:           msgPropertyRequiredWhen,
	msgPropertyUnwantedWhen:           msgPropertyUnwantedWhen,
	msgValueCannotBeNull:              msgValueCannotBeNull,
	msgValueMustBeObject:              msgValueMustBeObject,
	msgValueMustBeArray:               msgValueMustBeArray,
	msgValueMustBeObjectOrArray:       msgValueMustBeObjectOrArray,
	msgPropertyObjectValidatorError:   msgPropertyObjectValidatorError,
	msgNotEmptyString:                 msgNotEmptyString,
	msgNotBlankString:                 msgNotBlankString,
	msgNoControlChars:                 msgNoControlChars,
	msgValidPattern:                   msgValidPattern,
	msgInvalidCharacters:              msgInvalidCharacters,
	msgStringLowercase:                msgStringLowercase,
	msgStringUppercase:                msgStringUppercase,
	msgUnicodeNormalization:           msgUnicodeNormalization,
	msgUnicodeNormalizationNFC:        msgUnicodeNormalizationNFC,
	msgUnicodeNormalizationNFKC:       msgUnicodeNormalizationNFKC,
	msgUnicodeNormalizationNFD:        msgUnicodeNormalizationNFD,
	msgUnicodeNormalizationNFKD:       msgUnicodeNormalizationNFKD,
	msgPositive:                       msgPositive,
	msgPositiveOrZero:                 msgPositiveOrZero,
	msgNegative:                       msgNegative,
	msgNegativeOrZero:                 msgNegativeOrZero,
	msgArrayUnique:                    msgArrayUnique,
	msgValidUuid:                      msgValidUuid,
	msgValidCardNumber:                msgValidCardNumber,
	msgValidEmail:                     msgValidEmail,
	msgFailure:                        msgFailure,
	msgValidISODate:                   msgValidISODate,
	msgValidISODatetimeFormatFull:     msgValidISODatetimeFormatFull,
	msgValidISODatetimeFormatNoOffs:   msgValidISODatetimeFormatNoOffs,
	msgValidISODatetimeFormatNoMillis: msgValidISODatetimeFormatNoMillis,
	msgValidISODatetimeFormatMin:      msgValidISODatetimeFormatMin,
	msgDatetimeFuture:                 msgDatetimeFuture,
	msgDatetimeFutureOrPresent:        msgDatetimeFutureOrPresent,
	msgDatetimePast:                   msgDatetimePast,
	msgDatetimePastOrPresent:          msgDatetimePastOrPresent,
	msgPresetISBN:                     msgPresetISBN,
	msgPresetISBN10:                   msgPresetISBN10,
	msgPresetISBN13:                   msgPresetISBN13,
	msgPresetISSN:                     msgPresetISSN,
	msgPresetEAN:                      msgPresetEAN,
	msgPresetEAN8:                     msgPresetEAN8,
	msgPresetEAN13:                    msgPresetEAN13,
	msgPresetDUN14:                    msgPresetDUN14,
	msgPresetEAN14:                    msgPresetEAN14,
	msgPresetEAN18:                    msgPresetEAN18,
	msgPresetEAN99:                    msgPresetEAN99,
	msgPresetUPC:                      msgPresetUPC,
	msgPresetUPCA:                     msgPresetUPCA,
	msgPresetUPCE:                     msgPresetUPCE,
	msgPresetPublication:              msgPresetPublication,
	msgPresetAlpha:                    msgPresetAlpha,
	msgPresetAlphaNumeric:             msgPresetAlphaNumeric,
	msgPresetBarcode:                  msgPresetBarcode,
	msgPresetNumeric:                  msgPresetNumeric,
	msgPresetInteger:                  msgPresetInteger,
	msgPresetHexadecimal:              msgPresetHexadecimal,
	msgPresetCMYK:                     msgPresetCMYK,
	msgPresetCMYK300:                  msgPresetCMYK300,
	msgPresetHtmlColor:                msgPresetHtmlColor,
	msgPresetRgb:                      msgPresetRgb,
	msgPresetRgba:                     msgPresetRgba,
	msgPresetRgbIcc:                   msgPresetRgbIcc,
	msgPresetHsl:                      msgPresetHsl,
	msgPresetHsla:                     msgPresetHsla,
	msgPresetE164:                     msgPresetE164,
	msgPresetBase64:                   msgPresetBase64,
	msgPresetBase64URL:                msgPresetBase64URL,
	msgPresetUuid1:                    msgPresetUuid1,
	msgPresetUuid2:                    msgPresetUuid2,
	msgPresetUuid3:                    msgPresetUuid3,
	msgPresetUuid4:                    msgPresetUuid4,
	msgPresetUuid5:                    msgPresetUuid5,
	msgPresetULID:                     msgPresetULID,
}

// used by defaultI18nContext.MarshalJSON - to allow listing of translation reference
// also as a kind of double-entry bookkeeping to ensure everything has been translated into TranslationsFormats
var internalFormats = map[string]string{
	// property validator...
	fmtMsgValueExpectedType: fmtMsgValueExpectedType,
	// constraints...
	fmtMsgUnknownPresetPattern:   fmtMsgUnknownPresetPattern,
	fmtMsgValidToken:             fmtMsgValidToken,
	fmtMsgStringMinLen:           fmtMsgStringMinLen,
	fmtMsgStringMinLenExc:        fmtMsgStringMinLenExc,
	fmtMsgStringMaxLen:           fmtMsgStringMaxLen,
	fmtMsgStringMaxLenExc:        fmtMsgStringMaxLenExc,
	fmtMsgStringExactLen:         fmtMsgStringExactLen,
	fmtMsgStringMinMaxLen:        fmtMsgStringMinMaxLen,
	fmtMsgMinLen:                 fmtMsgMinLen,
	fmtMsgMinLenExc:              fmtMsgMinLenExc,
	fmtMsgExactLen:               fmtMsgExactLen,
	fmtMsgMinMax:                 fmtMsgMinMax,
	fmtMsgGt:                     fmtMsgGt,
	fmtMsgGte:                    fmtMsgGte,
	fmtMsgLt:                     fmtMsgLt,
	fmtMsgLte:                    fmtMsgLte,
	fmtMsgRange:                  fmtMsgRange,
	fmtMsgMultipleOf:             fmtMsgMultipleOf,
	fmtMsgArrayElementType:       fmtMsgArrayElementType,
	fmtMsgArrayElementTypeOrNull: fmtMsgArrayElementTypeOrNull,
	fmtMsgUuidMinVersion:         fmtMsgUuidMinVersion,
	fmtMsgUuidCorrectVer:         fmtMsgUuidCorrectVer,
	fmtMsgEqualsOther:            fmtMsgEqualsOther,
	fmtMsgNotEqualsOther:         fmtMsgNotEqualsOther,
	fmtMsgGtOther:                fmtMsgGtOther,
	fmtMsgGteOther:               fmtMsgGteOther,
	fmtMsgLtOther:                fmtMsgLtOther,
	fmtMsgLteOther:               fmtMsgLteOther,
	fmtMsgDtGt:                   fmtMsgDtGt,
	fmtMsgDtGte:                  fmtMsgDtGte,
	fmtMsgDtLt:                   fmtMsgDtLt,
	fmtMsgDtLte:                  fmtMsgDtLte,
	// datetime tolerances...
	fmtMsgDtToleranceFixedSame:                      fmtMsgDtToleranceFixedSame,
	fmtMsgDtToleranceFixedNotMoreThanAfterPlural:    fmtMsgDtToleranceFixedNotMoreThanAfterPlural,
	fmtMsgDtToleranceFixedNotMoreThanAfterSingular:  fmtMsgDtToleranceFixedNotMoreThanAfterSingular,
	fmtMsgDtToleranceFixedNotMoreThanBeforePlural:   fmtMsgDtToleranceFixedNotMoreThanBeforePlural,
	fmtMsgDtToleranceFixedNotMoreThanBeforeSingular: fmtMsgDtToleranceFixedNotMoreThanBeforeSingular,
	fmtMsgDtToleranceNowSame:                        fmtMsgDtToleranceNowSame,
	fmtMsgDtToleranceNowNotMoreThanAfterPlural:      fmtMsgDtToleranceNowNotMoreThanAfterPlural,
	fmtMsgDtToleranceNowNotMoreThanAfterSingular:    fmtMsgDtToleranceNowNotMoreThanAfterSingular,
	fmtMsgDtToleranceNowNotMoreThanBeforePlural:     fmtMsgDtToleranceNowNotMoreThanBeforePlural,
	fmtMsgDtToleranceNowNotMoreThanBeforeSingular:   fmtMsgDtToleranceNowNotMoreThanBeforeSingular,
	fmtMsgDtToleranceOtherSame:                      fmtMsgDtToleranceOtherSame,
	fmtMsgDtToleranceOtherNotMoreThanAfterPlural:    fmtMsgDtToleranceOtherNotMoreThanAfterPlural,
	fmtMsgDtToleranceOtherNotMoreThanAfterSingular:  fmtMsgDtToleranceOtherNotMoreThanAfterSingular,
	fmtMsgDtToleranceOtherNotMoreThanBeforePlural:   fmtMsgDtToleranceOtherNotMoreThanBeforePlural,
	fmtMsgDtToleranceOtherNotMoreThanBeforeSingular: fmtMsgDtToleranceOtherNotMoreThanBeforeSingular,
	// constraint set...
	fmtMsgConstraintSetDefaultAllOf: fmtMsgConstraintSetDefaultAllOf,
	fmtMsgConstraintSetDefaultOneOf: fmtMsgConstraintSetDefaultOneOf,
}

// used by defaultI18nContext.MarshalJSON - to allow listing of translation reference
// also as a kind of double-entry bookkeeping to ensure everything has been translated into TranslationsTokens
// "..." at the end indicates pluralisation
var internalTokens = map[string]string{
	jsonTypeTokenString:  jsonTypeTokenString,
	jsonTypeTokenNumber:  jsonTypeTokenNumber,
	jsonTypeTokenInteger: jsonTypeTokenInteger,
	jsonTypeTokenBoolean: jsonTypeTokenBoolean,
	jsonTypeTokenObject:  jsonTypeTokenObject,
	jsonTypeTokenArray:   jsonTypeTokenArray,
	jsonTypeTokenAny:     jsonTypeTokenAny,
	tokenInclusive:       tokenInclusive,
	tokenExclusive:       tokenExclusive,
	"millennium":         "millennium",
	"millennium...":      "millennia",
	"century":            "century",
	"century...":         "centuries",
	"decade":             "decade",
	"decade...":          "decades",
	"year":               "year",
	"year...":            "years",
	"month":              "month",
	"month...":           "months",
	"week":               "week",
	"week...":            "weeks",
	"day":                "day",
	"day...":             "days",
	"hour":               "hour",
	"hour...":            "hours",
	"minute":             "minute",
	"minute...":          "minutes",
	"second":             "second",
	"second...":          "seconds",
	"millisecond":        "millisecond",
	"millisecond...":     "milliseconds",
	"microsecond":        "microsecond",
	"microsecond...":     "microseconds",
	"nanosecond":         "nanosecond",
	"nanosecond...":      "nanoseconds",
}
