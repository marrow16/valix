package valix

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringValidCardNumberConstraint(t *testing.T) {
	validator := buildFooValidator(JsonAny,
		&StringValidCardNumber{}, false)
	obj := map[string]interface{}{
		"foo": "",
	}
	testCardNumbers := map[string]bool{
		// valid VISA...
		"4902498374064506":    true,
		"4556494687321500":    true,
		"4969975508508776718": true,
		// valid MasterCard...
		"5385096580406173": true,
		"5321051022936318": true,
		"5316711656334695": true,
		// valid American Express (AMEX)...
		"344349836405221": true,
		"370579317661267": true,
		"379499960274519": true,
		// valid Discover...
		"6011816780702703":    true,
		"6011803416863109":    true,
		"6011658115940081575": true,
		// valid JCB...
		"3537599036021738":    true,
		"3535975622588649":    true,
		"3538751375142304859": true,
		// valid Diners Club (North America)...
		"5449933541584900": true,
		"5424932328289252": true,
		"5443511725058507": true,
		// valid Diners Club (Carte Blanche)...
		"30533389154463": true,
		"30225631860175": true,
		"30132515376577": true,
		// valid Diners Club (International)...
		"36023499511897": true,
		"36895205164933": true,
		"36614132300415": true,
		// valid Maestro...
		"5020799867464796": true,
		"5893155499331362": true,
		"6763838557832695": true,
		// valid Visa Electron...
		"4917958531215104": true,
		"4913912408530396": true,
		"4917079458677141": true,
		// valid InstaPayment...
		"6387294734923401": true,
		"6382441564848878": true,
		"6371830528023664": true,
		// valid all zeroes...
		"0000000000":          true,
		"00000000000":         true,
		"000000000000":        true,
		"0000000000000":       true,
		"00000000000000":      true,
		"000000000000000":     true,
		"0000000000000000":    true,
		"00000000000000000":   true,
		"000000000000000000":  true,
		"0000000000000000000": true,
		// invalid all zeroes...
		"000000000":            false,
		"00000000000000000000": false,
		// invalids...
		"1234567890123":        false, // too short
		"12345678901234567890": false, // too long
		"4902498374064505":     false,
		"4969975508508776717":  false,
		"5385096580406171":     false,
		"344349836405220":      false,
		"6011816780702709":     false,
		"6011658115940081576":  false,
		"3537599036021730":     false,
		"3538751375142304850":  false,
		"5449933541584901":     false,
		"30533389154466":       false,
		"36023499511898":       false,
		"5020799867464791":     false,
		"4917958531215105":     false,
		"6387294734923400":     false,
		// invalid bad digits...
		"12345678901234x": false,
		// with spaces (fails without AllowSpaces set)...
		"4902 4983 7406 4506":     false,
		"5385 0965 8040 6173":     false,
		"3443 4983 6405 221":      false,
		"6011 8167 8070 2703":     false,
		"3537 5990 3602 1738":     false,
		"5449 9335 4158 4900":     false,
		"3053 3389 1544 63":       false,
		"5020 7998 6746 4796":     false,
		"4917 9585 3121 5104":     false,
		"6387 2947 3492 3401":     false,
		"0000 0000 00":            false,
		"0000 0000 000":           false,
		"0000 0000 0000":          false,
		"0000 0000 0000 0":        false,
		"0000 0000 0000 00":       false,
		"0000 0000 0000 000":      false,
		"0000 0000 0000 0000":     false,
		"0000 0000 0000 0000 0":   false,
		"0000 0000 0000 0000 00":  false,
		"0000 0000 0000 0000 000": false,
	}
	for ccn, expect := range testCardNumbers {
		t.Run(fmt.Sprintf("CardNumber:\"%s\"", ccn), func(t *testing.T) {
			obj["foo"] = ccn
			ok, _ := validator.Validate(obj)
			require.Equal(t, expect, ok)
		})
	}
	// and again with spaces and AllowSpaces set...
	testCardNumbers = map[string]bool{
		"4902 4983 7406 4506":     true,
		"5385 0965 8040 6173":     true,
		"3443 4983 6405 221":      true,
		"6011 8167 8070 2703":     true,
		"3537 5990 3602 1738":     true,
		"5449 9335 4158 4900":     true,
		"3053 3389 1544 63":       true,
		"5020 7998 6746 4796":     true,
		"4917 9585 3121 5104":     true,
		"6387 2947 3492 3401":     true,
		"0000 0000 00":            true,
		"0000 0000 000":           true,
		"0000 0000 0000":          true,
		"0000 0000 0000 0":        true,
		"0000 0000 0000 00":       true,
		"0000 0000 0000 000":      true,
		"0000 0000 0000 0000":     true,
		"0000 0000 0000 0000 0":   true,
		"0000 0000 0000 0000 00":  true,
		"0000 0000 0000 0000 000": true,
		// spaces in wrong places...
		" 4902 4983 7406 4506": false,
		"490 24983 7406 4506":  false,
		"4902 4983 7406 4506 ": false,
	}
	validator = buildFooValidator(JsonAny,
		&StringValidCardNumber{AllowSpaces: true}, false)
	for ccn, expect := range testCardNumbers {
		t.Run(fmt.Sprintf("CardNumber:\"%s\"", ccn), func(t *testing.T) {
			obj["foo"] = ccn
			ok, _ := validator.Validate(obj)
			require.Equal(t, expect, ok)
		})
	}
}

func TestStringValidCountryCode(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	c := &StringValidCountryCode{}
	for cc := range ISO3166_2_CountryCodes {
		ok, _ := c.Check(cc, vcx)
		require.True(t, ok)
	}
}

func TestStringValidCountryCodeAssignmentTypes(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	c := &StringValidCountryCode{}
	for l1, subs := range iso3166_1_CountryCodesMatrix {
		for l2, assignmentType := range subs {
			cc := string([]byte{l1, l2})
			switch assignmentType {
			case ccAs, ccRa:
				ok, _ := c.Check(cc, vcx)
				require.True(t, ok)
			case ccER:
				c.Allow3166_1_ExceptionallyReserved = false
				ok, _ := c.Check(cc, vcx)
				require.False(t, ok)
				c.Allow3166_1_ExceptionallyReserved = true
				ok, _ = c.Check(cc, vcx)
				require.True(t, ok)
			case ccTR:
				c.Allow3166_1_TransitionallyReserved = false
				ok, _ := c.Check(cc, vcx)
				require.False(t, ok)
				c.Allow3166_1_TransitionallyReserved = true
				ok, _ = c.Check(cc, vcx)
				require.True(t, ok)
			case ccIR:
				c.Allow3166_1_IndeterminatelyReserved = false
				ok, _ := c.Check(cc, vcx)
				require.False(t, ok)
				c.Allow3166_1_IndeterminatelyReserved = true
				ok, _ = c.Check(cc, vcx)
				require.True(t, ok)
			case ccUA:
				c.AllowUserAssigned = false
				ok, _ := c.Check(cc, vcx)
				require.False(t, ok)
				c.AllowUserAssigned = true
				ok, _ = c.Check(cc, vcx)
				require.True(t, ok)
			case ccDl:
				c.Allow3166_1_Deleted = false
				ok, _ := c.Check(cc, vcx)
				require.False(t, ok)
				c.Allow3166_1_Deleted = true
				ok, _ = c.Check(cc, vcx)
				require.True(t, ok)
			default:
				ok, _ := c.Check(cc, vcx)
				require.False(t, ok)
			}
		}
	}
}

func TestStringValidCountryCodeAndRegion(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	c := &StringValidCountryCode{}
	gb := ISO3166_2_CountryCodes["GB"]
	for r := range gb {
		c.Allow3166_2 = false
		ok, _ := c.Check("GB-"+r, vcx)
		require.False(t, ok)
		c.Allow3166_2 = true
		ok, _ = c.Check("GB-"+r, vcx)
		require.True(t, ok)
	}
	c.Allow3166_2 = false
	ok, _ := c.Check("GB-ZZZ", vcx)
	require.False(t, ok)
	c.Allow3166_2 = true
	ok, _ = c.Check("GB-ZZZ", vcx)
	require.False(t, ok)
}

func TestStringValidCountryCodeAndRegion_ObsoleteISO3166_2_Codes(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	c := &StringValidCountryCode{Allow3166_2: true}
	for cc := range ISO3166_2_ObsoleteCodes {
		c.Allow3166_2_Obsoletes = false
		ok, _ := c.Check(cc, vcx)
		require.False(t, ok)
		c.Allow3166_2_Obsoletes = true
		ok, _ = c.Check(cc, vcx)
		require.True(t, ok)
	}

}

func TestTestStringValidCountryCodeNumericOnly(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	c := &StringValidCountryCode{}
	for v := range ISO3166_1_NumericCodes {
		iv, _ := strconv.Atoi(v)
		c.NumericOnly = false
		c.Allow3166_1_Numeric = false
		ok, _ := c.Check(v, vcx)
		require.False(t, ok)
		ok, _ = c.Check(iv, vcx)
		require.False(t, ok)

		c.NumericOnly = true
		ok, _ = c.Check(v, vcx)
		require.True(t, ok)
		ok, _ = c.Check(iv, vcx)
		require.True(t, ok)

		c.NumericOnly = false
		c.Allow3166_1_Numeric = true
		ok, _ = c.Check(v, vcx)
		require.True(t, ok)
		ok, _ = c.Check(iv, vcx)
		require.True(t, ok)
	}

	c.NumericOnly = true
	c.AllowUserAssigned = false
	ok, _ := c.Check(899, vcx)
	require.False(t, ok)
	ok, _ = c.Check("899", vcx)
	require.False(t, ok)
	ok, _ = c.Check(1000, vcx)
	require.False(t, ok)
	ok, _ = c.Check("1000", vcx)
	require.False(t, ok)

	for i := 900; i < 1000; i++ {
		istr := fmt.Sprintf("%03d", i)

		c.NumericOnly = true
		c.AllowUserAssigned = false
		ok, _ := c.Check(i, vcx)
		require.False(t, ok)
		ok, _ = c.Check(istr, vcx)
		require.False(t, ok)

		c.AllowUserAssigned = true
		ok, _ = c.Check(i, vcx)
		require.True(t, ok)
		ok, _ = c.Check(istr, vcx)
		require.True(t, ok)

		c.NumericOnly = false
		c.Allow3166_1_Numeric = true
		c.AllowUserAssigned = false
		ok, _ = c.Check(i, vcx)
		require.False(t, ok)
		ok, _ = c.Check(istr, vcx)
		require.False(t, ok)

		c.AllowUserAssigned = true
		ok, _ = c.Check(i, vcx)
		require.True(t, ok)
		ok, _ = c.Check(istr, vcx)
		require.True(t, ok)
	}
}

func TestStringValidCountryCodeBad(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	c := &StringValidCountryCode{}
	ok, _ := c.Check("ZZ", vcx)
	require.False(t, ok)

	ok, _ = c.Check("-", vcx)
	require.False(t, ok)
	ok, _ = c.Check("ZZ-ZZ", vcx)
	require.False(t, ok)
	ok, _ = c.Check("ZZ-ZZ-ZZ", vcx)
	require.False(t, ok)

	c.Allow3166_2 = true
	ok, _ = c.Check("-", vcx)
	require.False(t, ok)
	ok, _ = c.Check("ZZ-ZZ", vcx)
	require.False(t, ok)
	ok, _ = c.Check("ZZ-ZZ-ZZ", vcx)
	require.False(t, ok)
}

func TestStringValidCurrencyCode(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	c := &StringValidCurrencyCode{}
	// test currency code...
	ok, _ := c.Check(ISO4217TestCurrencyCode, vcx)
	require.False(t, ok)
	ok, _ = c.Check(ISO4217TestCurrencyCodeNumeric, vcx)
	require.False(t, ok)
	c.AllowTestCode = true
	ok, _ = c.Check(ISO4217TestCurrencyCode, vcx)
	require.True(t, ok)
	ok, _ = c.Check(ISO4217TestCurrencyCodeNumeric, vcx)
	require.False(t, ok)
	c.AllowNumeric = true
	ok, _ = c.Check(ISO4217TestCurrencyCodeNumeric, vcx)
	require.True(t, ok)
	// no currency code...
	c.AllowNumeric = false
	ok, _ = c.Check(ISO4217NoCurrencyCode, vcx)
	require.False(t, ok)
	ok, _ = c.Check(ISO4217NoCurrencyCodeNumeric, vcx)
	require.False(t, ok)
	c.AllowNoCode = true
	ok, _ = c.Check(ISO4217NoCurrencyCode, vcx)
	require.True(t, ok)
	ok, _ = c.Check(ISO4217NoCurrencyCodeNumeric, vcx)
	require.False(t, ok)
	c.AllowNumeric = true
	ok, _ = c.Check(ISO4217NoCurrencyCodeNumeric, vcx)
	require.True(t, ok)
	// regular currency codes...
	for code := range ISO4217CurrencyCodes {
		ok, _ = c.Check(code, vcx)
		require.True(t, ok)
	}
	ok, msg := c.Check("!!!", vcx)
	require.False(t, ok)
	require.Equal(t, msgValidCurrencyCode, msg)
	// regular numerics...
	for code := range ISO4217CurrencyCodesNumeric {
		c.AllowNumeric = false
		ok, _ = c.Check(code, vcx)
		require.False(t, ok)
		icode, _ := strconv.Atoi(code)
		ok, _ = c.Check(icode, vcx)
		require.False(t, ok)

		c.AllowNumeric = true
		ok, _ = c.Check(code, vcx)
		require.True(t, ok)
		ok, _ = c.Check(icode, vcx)
		require.True(t, ok)
	}
	// historical...
	for code := range ISO4217CurrencyCodesHistorical {
		c.AllowHistorical = false
		ok, _ = c.Check(code, vcx)
		require.False(t, ok)
		c.AllowHistorical = true
		ok, _ = c.Check(code, vcx)
		require.True(t, ok)
	}
	// historical numerics...
	c.AllowNumeric = true
	for code := range ISO4217CurrencyCodesNumericHistorical {
		c.AllowHistorical = false
		ok, _ = c.Check(code, vcx)
		require.False(t, ok)
		c.AllowHistorical = true
		ok, _ = c.Check(code, vcx)
		require.True(t, ok)
	}
	// unofficial...
	for code := range UnofficialCurrencyCodes {
		c.AllowUnofficial = false
		ok, _ = c.Check(code, vcx)
		require.False(t, ok)
		c.AllowUnofficial = true
		ok, _ = c.Check(code, vcx)
		require.True(t, ok)
	}
	// cryptos...
	for code := range CryptoCurrencyCodes {
		c.AllowCrypto = false
		ok, _ = c.Check(code, vcx)
		require.False(t, ok)
		c.AllowCrypto = true
		ok, _ = c.Check(code, vcx)
		require.True(t, ok)
	}
}

func TestStringValidCurrencyCodeNumericOnly(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	c := &StringValidCurrencyCode{NumericOnly: true}
	// alphas all fail...
	for code := range ISO4217CurrencyCodes {
		ok, _ := c.Check(code, vcx)
		require.False(t, ok)
	}
	// numerics...
	for code := range ISO4217CurrencyCodesNumeric {
		ok, _ := c.Check(code, vcx)
		require.True(t, ok)

		icode, _ := strconv.Atoi(code)
		ok, _ = c.Check(icode, vcx)
		require.True(t, ok)
	}

	ok, _ := c.Check(ISO4217TestCurrencyCodeNumeric, vcx)
	require.False(t, ok)
	icode, _ := strconv.Atoi(ISO4217TestCurrencyCodeNumeric)
	ok, _ = c.Check(icode, vcx)
	require.False(t, ok)
	c.AllowTestCode = true
	ok, _ = c.Check(ISO4217TestCurrencyCodeNumeric, vcx)
	require.True(t, ok)
	ok, _ = c.Check(icode, vcx)
	require.True(t, ok)

	ok, _ = c.Check(ISO4217NoCurrencyCodeNumeric, vcx)
	require.False(t, ok)
	icode, _ = strconv.Atoi(ISO4217NoCurrencyCodeNumeric)
	ok, _ = c.Check(icode, vcx)
	require.False(t, ok)
	c.AllowNoCode = true
	ok, _ = c.Check(ISO4217NoCurrencyCodeNumeric, vcx)
	require.True(t, ok)
	ok, _ = c.Check(icode, vcx)
	require.True(t, ok)
}

func TestStringValidEmail(t *testing.T) {
	c := &StringValidEmail{}
	vcx := newValidatorContext(nil, nil, false, nil)
	testCases := map[string]struct {
		expectValid            bool
		expectValidWithOptions bool
		tryWith                *StringValidEmail
	}{
		"":                      {},
		"me@com.":               {},
		"me@123":                {},
		"me@abc.-bad-.com":      {},
		"me and you@google.com": {},
		"me.too.long.34567890123456789012345678901234567890123456789012345@google.com": {},
		"me@too.long.a1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234": {},
		"me@too.long.after.expand.சிங்கப்பூர்.சிங்கப்பூர்.சிங்கப்பூர்.சிங்கப்பூர்.சிங்கப்பூர்.சிங்கப்பூர்.சிங்கப்பூர்.சிங்கப்பூர்.சிங்கப்பூர்.சிங்கப்பூர்.சிங்கப்பூர்":                                                                                                      {},
		"me@com": {
			tryWith: &StringValidEmail{
				AllowTldOnly: true,
			},
		},
		"me@[127.0.0.1]": {
			expectValid:            false,
			expectValidWithOptions: true,
			tryWith: &StringValidEmail{
				AllowLocal:     true,
				AllowIPAddress: true,
			},
		},
		"me@[123.456.789.0]": {
			expectValid:            false,
			expectValidWithOptions: false,
			tryWith: &StringValidEmail{
				AllowIPAddress: true,
			},
		},
		"me@[123.45.67.89]": {
			expectValid:            false,
			expectValidWithOptions: true,
			tryWith: &StringValidEmail{
				AllowLocal:     true,
				AllowIPAddress: true,
			},
		},
		"me@[::1]": {
			expectValid:            false,
			expectValidWithOptions: true,
			tryWith: &StringValidEmail{
				AllowLocal:     true,
				AllowIPAddress: true,
				AllowIPV6:      true,
			},
		},
		"me2@[::1]": {
			expectValid:            false,
			expectValidWithOptions: false,
			tryWith: &StringValidEmail{
				AllowLocal:     true,
				AllowIPAddress: true,
				AllowIPV6:      false,
			},
		},
		"me@google.com": {
			expectValid: true,
		},
		"me2@google.com": {
			expectValid:            true,
			expectValidWithOptions: true,
			tryWith: &StringValidEmail{
				DisallowRFC5322: true,
			},
		},
		"Me@Google.Com": {
			expectValid: true,
		},
		"\"marrow\" <me@google.com>": {
			expectValid:            true,
			expectValidWithOptions: false,
			tryWith: &StringValidEmail{
				DisallowRFC5322: true,
			},
		},
		"me@localhost": {
			expectValid:            false,
			expectValidWithOptions: true,
			tryWith: &StringValidEmail{
				AllowLocal: true,
			},
		},
		"me2@localhost": {
			expectValid:            false,
			expectValidWithOptions: false,
			tryWith: &StringValidEmail{
				AllowLocal:   true,
				ExcLocalTlds: []string{"localhost"},
			},
		},
		"me@example.org": {
			expectValid:            false,
			expectValidWithOptions: true,
			tryWith: &StringValidEmail{
				AllowTestTlds: true,
			},
		},
		"me@audi": {
			expectValid:            false,
			expectValidWithOptions: true,
			tryWith: &StringValidEmail{
				AllowTldOnly:   true,
				AllowBrandTlds: true,
			},
		},
		"me@sample.ಭಾರತ": {
			expectValid: true,
		},
		"me@sample.xn--2scrj9c": {
			expectValid: true,
		},
		"me@some.academy": {
			expectValid:            false,
			expectValidWithOptions: true,
			tryWith: &StringValidEmail{
				AllowGenericTlds: true,
			},
		},
		"me2@some.academy": {
			expectValid:            false,
			expectValidWithOptions: false,
			tryWith: &StringValidEmail{
				AllowGenericTlds: true,
				ExcGenericTlds:   []string{"academy"},
			},
		},
		"Bilbo <bilbo@example.com>": {
			expectValid:            false,
			expectValidWithOptions: true,
			tryWith: &StringValidEmail{
				AllowTestTlds: true,
			},
		},
		"me@some-company.africa": {
			expectValid:            false,
			expectValidWithOptions: true,
			tryWith: &StringValidEmail{
				AllowGeographicTlds: true,
			},
		},
		"me2@some-company.africa": {
			expectValid:            false,
			expectValidWithOptions: false,
			tryWith: &StringValidEmail{
				AllowGeographicTlds: true,
				ExcCountryCodeTlds:  []string{"africa"},
			},
		},
		"me@my.audi": {
			expectValid:            false,
			expectValidWithOptions: true,
			tryWith: &StringValidEmail{
				AllowBrandTlds: true,
			},
		},
		"me2@my.audi": {
			expectValid:            false,
			expectValidWithOptions: false,
			tryWith: &StringValidEmail{
				AllowBrandTlds: true,
				ExcBrandTlds:   []string{"audi"},
			},
		},
		"me@arpa": {
			expectValid:            false,
			expectValidWithOptions: true,
			tryWith: &StringValidEmail{
				AllowInfraTlds: true,
			},
		},
	}
	for addr, tc := range testCases {
		t.Run(fmt.Sprintf("Email:\"%s\"", addr), func(t *testing.T) {
			ok, _ := c.Check(addr, vcx)
			require.Equal(t, tc.expectValid, ok)
			if tc.tryWith != nil {
				ok, _ := tc.tryWith.Check(addr, vcx)
				require.Equal(t, tc.expectValidWithOptions, ok)
			}
		})
	}

	// with non-string...
	ok, _ := c.Check(true, vcx)
	require.True(t, ok)
}

func TestStringValidLanguageCode(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	c := &StringValidLanguageCode{}

	ok, _ := c.Check("en", vcx)
	require.True(t, ok)
	ok, _ = c.Check("en_GB", vcx)
	require.True(t, ok)
	ok, _ = c.Check("en-US", vcx)
	require.True(t, ok)

	ok, _ = c.Check("en-", vcx)
	require.False(t, ok)
	ok, _ = c.Check("en_", vcx)
	require.False(t, ok)
	ok, _ = c.Check("en-.", vcx)
	require.False(t, ok)
}

func TestStringValidUuid(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringValidUuid{MinVersion: 4}, false)
	obj := jsonObject(`{
		"foo": "not a uuid"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgUuidMinVersion, 4), violations[0].Message)

	obj["foo"] = "db15398d-328f-4d16-be2f-f38e8f2d0a79"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "db15398d-328f-3d16-be2f-f38e8f2d0a79"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgUuidMinVersion, 4), violations[0].Message)

	validator = buildFooValidator(JsonString,
		&StringValidUuid{SpecificVersion: 4}, false)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgUuidCorrectVer, 4), violations[0].Message)

	validator = buildFooValidator(JsonString,
		&StringValidUuid{}, false)
	obj = jsonObject(`{
		"foo": "not a uuid"
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(msgValidUuid), violations[0].Message)
}
