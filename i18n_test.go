package valix

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
)

func TestFallbackI18nContext(t *testing.T) {
	i18ctx := DefaultI18nProvider.DefaultContext()
	require.Equal(t, fallbackI18nContext, i18ctx)
}

func TestEnsureNonNilI18nContext(t *testing.T) {
	defer func() {
		DefaultI18nProvider = fallbackI18nProvider
	}()

	i18ctx := obtainI18nContext(nil)
	require.NotNil(t, i18ctx)

	DefaultI18nProvider = nil
	i18ctx = obtainI18nContext(nil)
	require.NotNil(t, i18ctx)
}

func TestCanTranslateToken(t *testing.T) {
	tcx := &defaultI18nContext{lang: "en"}
	s := tcx.TranslateToken("xxxx")
	require.Equal(t, "xxxx", s)

	s = tcx.TranslateToken("century")
	require.Equal(t, "century", s)
	s = tcx.TranslateToken("century...")
	require.Equal(t, "centuries", s)
}

func TestCanTranslateMessage(t *testing.T) {
	tcx := &defaultI18nContext{lang: "en"}
	s := tcx.TranslateMessage("xxxx")
	require.Equal(t, "xxxx", s)

	s = tcx.TranslateMessage(msgUnableToDecode)
	require.Equal(t, msgUnableToDecode, s)
}

func TestCanTranslateFormat(t *testing.T) {
	tcx := &defaultI18nContext{lang: "en"}
	s := tcx.TranslateFormat("xxxx %d", 1)
	require.Equal(t, "xxxx 1", s)

	s = tcx.TranslateFormat(fmtMsgValueExpectedType, "foo")
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, "foo"), s)
}

func TestCanMarshalDefaultI18nContext(t *testing.T) {
	_, err := json.Marshal(&defaultI18nContext{lang: "en"})
	require.Nil(t, err)
}

func TestAllTokensTranslated(t *testing.T) {
	require.Equal(t, len(internalTokens), len(defaultInternalTranslator.Tokens))
	for k := range internalTokens {
		_, ok := defaultInternalTranslator.Tokens[k]
		require.True(t, ok)
	}
}

func TestAllMessagesTranslated(t *testing.T) {
	require.Equal(t, len(internalMessages), len(defaultInternalTranslator.Messages))
	for k := range internalMessages {
		_, ok := defaultInternalTranslator.Messages[k]
		require.True(t, ok)
	}
}

func TestAllFormatMessagesTranslated(t *testing.T) {
	require.Equal(t, len(internalFormats), len(defaultInternalTranslator.Formats))
	for k := range internalFormats {
		_, ok := defaultInternalTranslator.Formats[k]
		require.True(t, ok)
	}
}

func TestDefaultLanguage(t *testing.T) {
	defer func() {
		DefaultLanguage = "en"
	}()

	lang := defaultLanguage("")
	require.Equal(t, "en", lang)

	lang = defaultLanguage("EN")
	require.Equal(t, "en", lang)

	lang = defaultLanguage("fr")
	require.Equal(t, "fr", lang)

	lang = defaultLanguage("xx")
	require.Equal(t, "en", lang)

	DefaultLanguage = "xx"
	lang = defaultLanguage("xx")
	require.Equal(t, "en", lang)
	lang = defaultLanguage("en")
	require.Equal(t, "en", lang)
	lang = defaultLanguage("")
	require.Equal(t, "en", lang)
	lang = defaultLanguage("fr")
	require.Equal(t, "fr", lang)
}

func TestDefaultFallbackLanguage(t *testing.T) {
	defer func() {
		DefaultFallbackLanguages = map[string]string{}
	}()
	i18ctx := newDefaultI18nContext("mt", "")
	require.Equal(t, "en", i18ctx.Language())

	DefaultFallbackLanguages["mt"] = "IT"
	i18ctx = newDefaultI18nContext("mt", "")
	require.Equal(t, "it", i18ctx.Language())
}

func TestRegionTranslations(t *testing.T) {
	tcx := &defaultI18nContext{
		lang:   "en",
		region: "US",
	}
	msg := tcx.TranslateMessage(msgPresetRgb)
	require.Equal(t, "Value must be a valid rgb() color string", msg)
	msg = tcx.TranslateMessage(msgPresetPublication)
	require.Equal(t, msgPresetPublication, msg)

	tcx.region = ""
	msg = tcx.TranslateMessage(msgPresetRgb)
	require.Equal(t, msgPresetRgb, msg)

	tcx.region = "US"
	DefaultTranslator.AddMessageLanguageTranslation("en", "test", "TEST ENGLISH",
		RegionalVariantTranslation{"US", "TEST ENGLISH US"})
	msg = tcx.TranslateMessage("test")
	require.Equal(t, "TEST ENGLISH US", msg)
	tcx.region = "GB"
	msg = tcx.TranslateMessage("test")
	require.Equal(t, "TEST ENGLISH", msg)
	DefaultTranslator.AddMessageRegionTranslation("en", "GB", "test", "TEST ENGLISH GB")
	msg = tcx.TranslateMessage("test")
	require.Equal(t, "TEST ENGLISH GB", msg)

	tcx.region = "US"
	DefaultTranslator.AddTokenLanguageTranslation("en", "test", "TEST ENGLISH",
		RegionalVariantTranslation{"US", "TEST ENGLISH US"})
	msg = tcx.TranslateToken("test")
	require.Equal(t, "TEST ENGLISH US", msg)
	tcx.region = "GB"
	msg = tcx.TranslateToken("test")
	require.Equal(t, "TEST ENGLISH", msg)
	DefaultTranslator.AddTokenRegionTranslation("en", "GB", "test", "TEST ENGLISH GB")
	msg = tcx.TranslateToken("test")
	require.Equal(t, "TEST ENGLISH GB", msg)

	tcx.region = "US"
	DefaultTranslator.AddFormatLanguageTranslation("en", "TEST %d", "TEST ENGLISH %d",
		RegionalVariantTranslation{"US", "%d TEST ENGLISH US"})
	msg = tcx.TranslateFormat("TEST %d", 1)
	require.Equal(t, "1 TEST ENGLISH US", msg)
	tcx.region = "GB"
	msg = tcx.TranslateFormat("TEST %d", 2)
	require.Equal(t, "TEST ENGLISH 2", msg)
	DefaultTranslator.AddFormatRegionTranslation("en", "GB", "TEST %d", "TEST ENGLISH GB %d")
	msg = tcx.TranslateFormat("TEST %d", 3)
	require.Equal(t, "TEST ENGLISH GB 3", msg)
}

func TestTranslatorTokenAdditions(t *testing.T) {
	const testToken = "TEST_TEST"
	msg := DefaultTranslator.TranslateToken("en", "", testToken)
	require.Equal(t, testToken, msg)
	l := len(defaultInternalTranslator.Tokens)

	DefaultTranslator.AddTokenLanguageTranslation("en", testToken, "TEST IN ENGLISH", RegionalVariantTranslation{"gb", ""})
	require.Equal(t, l+1, len(defaultInternalTranslator.Tokens))
	msg = DefaultTranslator.TranslateToken("en", "", testToken)
	require.Equal(t, "TEST IN ENGLISH", msg)
	msg = DefaultTranslator.TranslateToken("en", "GB", testToken)
	require.Equal(t, "TEST IN ENGLISH", msg)

	const testToken2 = "TEST_TEST_2"
	DefaultTranslator.AddTokenRegionTranslation("en", "GB", testToken2, "ANOTHER TEST IN ENGLISH")
	require.Equal(t, l+2, len(defaultInternalTranslator.Tokens))
	msg = DefaultTranslator.TranslateToken("en", "", testToken2)
	require.Equal(t, "ANOTHER TEST IN ENGLISH", msg)
	msg = DefaultTranslator.TranslateToken("en", "GB", testToken2)
	require.Equal(t, "ANOTHER TEST IN ENGLISH", msg)
}

func TestTranslatorMessageAdditions(t *testing.T) {
	const testMessage = "TEST_TEST"
	msg := DefaultTranslator.TranslateMessage("en", "", testMessage)
	require.Equal(t, testMessage, msg)
	l := len(defaultInternalTranslator.Messages)

	DefaultTranslator.AddMessageLanguageTranslation("en", testMessage, "TEST IN ENGLISH", RegionalVariantTranslation{"gb", ""})
	require.Equal(t, l+1, len(defaultInternalTranslator.Messages))
	msg = DefaultTranslator.TranslateMessage("en", "", testMessage)
	require.Equal(t, "TEST IN ENGLISH", msg)
	msg = DefaultTranslator.TranslateMessage("en", "GB", testMessage)
	require.Equal(t, "TEST IN ENGLISH", msg)

	const testMessage2 = "TEST_TEST_2"
	DefaultTranslator.AddMessageRegionTranslation("en", "GB", testMessage2, "ANOTHER TEST IN ENGLISH")
	require.Equal(t, l+2, len(defaultInternalTranslator.Messages))
	msg = DefaultTranslator.TranslateMessage("en", "", testMessage2)
	require.Equal(t, "ANOTHER TEST IN ENGLISH", msg)
	msg = DefaultTranslator.TranslateMessage("en", "GB", testMessage2)
	require.Equal(t, "ANOTHER TEST IN ENGLISH", msg)
}

func TestTranslatorFormatAdditions(t *testing.T) {
	const testFormat = "TEST_TEST %d"
	msg := DefaultTranslator.TranslateFormat("en", "", testFormat, 10)
	require.Equal(t, fmt.Sprintf(testFormat, 10), msg)
	l := len(defaultInternalTranslator.Formats)

	DefaultTranslator.AddFormatLanguageTranslation("en", testFormat, "TEST IN ENGLISH %d", RegionalVariantTranslation{"gb", ""})
	require.Equal(t, l+1, len(defaultInternalTranslator.Formats))
	msg = DefaultTranslator.TranslateFormat("en", "", testFormat, 11)
	require.Equal(t, "TEST IN ENGLISH 11", msg)
	msg = DefaultTranslator.TranslateFormat("en", "GB", testFormat, 12)
	require.Equal(t, "TEST IN ENGLISH 12", msg)

	const testFormat2 = "TEST_TEST_2 %d"
	DefaultTranslator.AddFormatRegionTranslation("en", "GB", testFormat2, "ANOTHER TEST IN ENGLISH %d")
	require.Equal(t, l+2, len(defaultInternalTranslator.Formats))
	msg = DefaultTranslator.TranslateFormat("en", "", testFormat2, 13)
	require.Equal(t, "ANOTHER TEST IN ENGLISH 13", msg)
	msg = DefaultTranslator.TranslateFormat("en", "GB", testFormat2, 14)
	require.Equal(t, "ANOTHER TEST IN ENGLISH 14", msg)
}

func TestRequestLanguage(t *testing.T) {
	testCases := []struct {
		header string
		lang   string
		region string
	}{
		{
			"nn;q=0.3, en;q=0.8, en-US,",
			"en",
			"US",
		},
		{
			"en-gb, en-US, en,",
			"en",
			"GB",
		},
		{
			"fr, it, es, de, en,",
			"fr",
			"",
		},
		{
			"fr;q=0.5, fr-CA, it, es, de, en,",
			"fr",
			"CA",
		},
		{
			"gsw, en;q=0.7, en-US;q=0.8",
			"en",
			"US",
		},
		{
			"gsw, nl, da",
			"en",
			"",
		},
		{
			"invalid",
			"en",
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("[%d]Accept-Language=\"%s\"", i+1, tc.header), func(t *testing.T) {
			r, _ := http.NewRequest("GET", "example.com", strings.NewReader("Hello"))
			r.Header.Set("Accept-Language", tc.header)
			tcx := DefaultI18nProvider.ContextFromRequest(r)
			require.NotNil(t, tcx)
			require.Equal(t, tc.lang, tcx.Language())
			require.Equal(t, tc.region, tcx.Region())
		})
	}
}

func TestValidationMessagesTranslated(t *testing.T) {
	foundLang := ""
	foundRegion := ""
	validator := &Validator{
		Constraints: Constraints{
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (passed bool, message string) {
				foundLang = vcx.Language()
				foundRegion = vcx.Region()
				return true, ""
			}, ""),
		},
		Properties: Properties{
			"foo": {
				NotNull: true,
				Constraints: Constraints{
					&StringPresetPattern{Preset: PresetRgb},
				},
			},
		},
	}
	r, _ := http.NewRequest("GET", "example.com", strings.NewReader(`{"foo": null}`))
	r.Header.Set("Accept-Language", "fr")
	ok, violations, _ := validator.RequestValidate(r)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "La valeur ne peut pas être nulle", violations[0].Message)
	require.Equal(t, "fr", foundLang)
	require.Equal(t, "", foundRegion)

	r, _ = http.NewRequest("GET", "example.com", strings.NewReader(`null`))
	r.Header.Set("Accept-Language", "es")
	ok, violations, _ = validator.RequestValidate(r)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "El cuerpo de la solicitud no debe ser JSON nulo", violations[0].Message)

	r, _ = http.NewRequest("GET", "example.com", strings.NewReader(`{"foo": "not an rgb string!"}`))
	r.Header.Set("Accept-Language", "en")
	ok, violations, _ = validator.RequestValidate(r)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Value must be a valid rgb() colour string", violations[0].Message)

	r, _ = http.NewRequest("GET", "example.com", strings.NewReader(`{"foo": "not an rgb string!"}`))
	r.Header.Set("Accept-Language", "en-us")
	ok, violations, _ = validator.RequestValidate(r)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Value must be a valid rgb() color string", violations[0].Message)
}
