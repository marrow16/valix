package valix

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"regexp"
	"strconv"
	"testing"
)

const builtInPresetsCount = 50

func TestPresetsRegistryInitialized(t *testing.T) {
	require.Equal(t, builtInPresetsCount, len(presetsRegistry.namedPresets))
	require.Equal(t, builtInPresetsCount, len(getBuiltInPresets()))
}

func TestPresetsRegistryResets(t *testing.T) {
	defer func() {
		presetsRegistry.reset()
	}()
	presetsRegistry.reset()
	require.Equal(t, builtInPresetsCount, len(presetsRegistry.namedPresets))

	presetsRegistry.register("fooey", patternPreset{})
	require.Equal(t, builtInPresetsCount+1, len(presetsRegistry.namedPresets))
	// and didn't corrupt builtins...
	require.Equal(t, builtInPresetsCount, len(getBuiltInPresets()))
}

func TestPresetsRegistryInternalRegister(t *testing.T) {
	defer func() {
		presetsRegistry.reset()
	}()
	presetsRegistry.reset()
	_, ok := presetsRegistry.get("fooey")
	require.False(t, ok)

	presetsRegistry.register("fooey", patternPreset{})
	_, ok = presetsRegistry.get("fooey")
	require.True(t, ok)
	require.Equal(t, builtInPresetsCount+1, len(presetsRegistry.namedPresets))
}

func TestPresetsRegistryExternalRegister(t *testing.T) {
	const testName = "fooey"
	defer func() {
		presetsRegistry.reset()
		constraintsRegistry.reset()
	}()
	presetsRegistry.reset()
	constraintsRegistry.reset()
	require.False(t, constraintsRegistry.has(testName))

	RegisterPresetPattern(testName, nil, "", nil, false)
	_, ok := presetsRegistry.get(testName)
	require.True(t, ok)
	require.Equal(t, builtInPresetsCount+1, len(presetsRegistry.namedPresets))
	require.False(t, constraintsRegistry.has(testName))

	RegisterPresetPattern(testName, nil, "", nil, true)
	_, ok = presetsRegistry.get(testName)
	require.True(t, ok)
	require.Equal(t, builtInPresetsCount+1, len(presetsRegistry.namedPresets))
	require.True(t, constraintsRegistry.has(testName))
}

func TestRegisteredPresetUsedAsV8nTag(t *testing.T) {
	type FooTest struct {
		Foo string `json:"foo" v8n:"&fooey"`
	}
	const testName = "fooey"
	defer func() {
		presetsRegistry.reset()
		constraintsRegistry.reset()
	}()
	presetsRegistry.reset()
	constraintsRegistry.reset()

	_, err := ValidatorFor(FooTest{}, nil)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownConstraint, testName), err.Error())

	RegisterPresetPattern(testName, nil, "", nil, true)
	v, err := ValidatorFor(FooTest{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)
	pv := v.Properties["foo"]
	require.Equal(t, 1, len(pv.Constraints))
	constraint := pv.Constraints[0].(*StringPresetPattern)
	require.Equal(t, testName, constraint.Preset)
}

func TestBarcodesWithModuloChecked(t *testing.T) {
	testCases := []struct {
		code        string
		preset      string
		expectRegex bool
		expectCheck bool
	}{
		{
			"123601057072",
			presetTokenUPCA,
			true,
			true,
		},
		{
			"123601057072",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"123601057072",
			presetTokenUPC,
			true,
			true,
		},
		{
			"01234565",
			presetTokenUPCE,
			true,
			true,
		},
		{
			"01234565",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"01234565",
			presetTokenUPC,
			true,
			true,
		},
		{
			"96385074",
			presetTokenEAN8,
			true,
			true,
		},
		{
			"96385074",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"96385074",
			presetTokenEAN,
			true,
			true,
		},
		{
			"996385074",
			presetTokenEAN,
			true,
			false,
		},
		{
			"96385079",
			presetTokenEAN8,
			true,
			false,
		},
		{
			"96385079",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"96385079",
			presetTokenEAN,
			true,
			false,
		},
		{
			"5000237128237",
			presetTokenEAN13,
			true,
			true,
		},
		{
			"5000237128237",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"5000237128237",
			presetTokenEAN,
			true,
			true,
		},
		{
			"5000237128239",
			presetTokenEAN13,
			true,
			false,
		},
		{
			"5000237128239",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"5000237128239",
			presetTokenEAN,
			true,
			false,
		},
		{
			"40700719670720",
			presetTokenEAN14,
			true,
			true,
		},
		{
			"40700719670729",
			presetTokenEAN14,
			true,
			false,
		},
		{
			"40700719670720",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"40700719670729",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"40700719670720",
			presetTokenEAN,
			true,
			true,
		},
		{
			"(01)40700719670720",
			presetTokenEAN14,
			true,
			true,
		},
		{
			"(01)40700719670720",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"(01)40700719670720",
			presetTokenEAN,
			true,
			true,
		},
		{
			"(01)40700719670729",
			presetTokenEAN14,
			true,
			false,
		},
		{
			"(01)40700719670729",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"(01)40700719670729",
			presetTokenEAN,
			true,
			false,
		},
		{
			"40700719670729",
			presetTokenEAN14,
			true,
			false,
		},
		{
			"40700719670729",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"40700719670729",
			presetTokenEAN,
			true,
			false,
		},
		{
			"407000000719670720",
			presetTokenEAN18,
			true,
			true,
		},
		{
			"407000000719670720",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"407000000719670720",
			presetTokenEAN,
			true,
			true,
		},
		{
			"407000000719670729",
			presetTokenEAN18,
			true,
			false,
		},
		{
			"407000000719670729",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"407000000719670729",
			presetTokenEAN,
			true,
			false,
		},
		{
			"(00)407000000719670720",
			presetTokenEAN18,
			true,
			true,
		},
		{
			"(00)407000000719670720",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"(00)407000000719670720",
			presetTokenEAN,
			true,
			true,
		},
		{
			"(00)407000000719670729",
			presetTokenEAN18,
			true,
			false,
		},
		{
			"(00)407000000719670729",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"(00)407000000719670729",
			presetTokenEAN,
			true,
			false,
		},
		{
			"5000237128237",
			presetTokenEAN99,
			false,
			false,
		},
		{
			"5000237128237",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"5000237128239",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"5000237128237",
			presetTokenEAN,
			true,
			true,
		},
		{
			"9911171967072",
			presetTokenEAN99,
			true,
			true,
		},
		{
			"9911171967072",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"9911171967072",
			presetTokenEAN,
			true,
			true,
		},
		{
			"9911171967079",
			presetTokenEAN99,
			true,
			false,
		},
		{
			"9911171967079",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"9911171967079",
			presetTokenEAN,
			true,
			false,
		},
		{
			"9771364743155",
			presetTokenISSN,
			true,
			true,
		},
		{
			"9771364743155",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"9770262407244",
			presetTokenISSN,
			true,
			true,
		},
		{
			"9770262407244",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"9770262407249",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"9770036873121",
			presetTokenISSN,
			true,
			true,
		},
		{
			"9772054638003",
			presetTokenISSN,
			true,
			true,
		},
		{
			"9771144875007",
			presetTokenISSN,
			true,
			true,
		},
		{
			"9771144875007",
			presetTokenPublication,
			true,
			true,
		},
		{
			"1144875X",
			presetTokenISSN,
			true,
			true,
		},
		{
			"1144875X",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"11448759",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"1144875X",
			presetTokenISSN8,
			true,
			true,
		},
		{
			"1144875X",
			presetTokenISSN13,
			false,
			false,
		},
		{
			"1144875X",
			presetTokenPublication,
			true,
			true,
		},
		{
			"020161622X",
			presetTokenISBN10,
			true,
			true,
		},
		{
			"020161622X",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"0201616229",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"020161622X",
			presetTokenISBN,
			true,
			true,
		},
		{
			"020161622X",
			presetTokenPublication,
			true,
			true,
		},
		{
			"9780201616224",
			presetTokenISBN13,
			true,
			true,
		},
		{
			"9780201616224",
			presetTokenISBN,
			true,
			true,
		},
		{
			"9780201616224",
			presetTokenPublication,
			true,
			true,
		},
		{
			"0805300910",
			presetTokenISBN10,
			true,
			true,
		},
		{
			"0805300910",
			presetTokenISBN,
			true,
			true,
		},
		{
			"9780201733860",
			presetTokenISBN13,
			true,
			true,
		},
		{
			"9780201733860",
			presetTokenISBN,
			true,
			true,
		},
		{
			"9780201733860",
			presetTokenISBN10,
			false,
			false,
		},
		{
			"0201733862",
			presetTokenISBN10,
			true,
			true,
		},
		{
			"0201733862",
			presetTokenISBN,
			true,
			true,
		},
		{
			"0201733862",
			presetTokenISBN13,
			false,
			false,
		},
		{
			"0805300911",
			presetTokenISBN10,
			true,
			false,
		},
		{
			"0805300911",
			presetTokenISBN,
			true,
			false,
		},
		{
			"0",
			presetTokenBarcode,
			false,
			false,
		},
		{
			"00",
			presetTokenBarcode,
			false,
			false,
		},
		{
			"000",
			presetTokenBarcode,
			false,
			false,
		},
		{
			"0000",
			presetTokenBarcode,
			false,
			false,
		},
		{
			"00000",
			presetTokenBarcode,
			false,
			false,
		},
		{
			"000000",
			presetTokenBarcode,
			false,
			false,
		},
		{
			"0000000",
			presetTokenBarcode,
			false,
			false,
		},
		{
			"00000000",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"0000000X",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"000000000",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"0000000000",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"00000000000",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"000000000000",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"0000000000000",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"00000000000000",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"000000000000000",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"0000000000000000",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"00000000000000000",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"000000000000000000",
			presetTokenBarcode,
			true,
			true,
		},
		{
			"0000000000000000000",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"00000000000000000000",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"000000000000000000000",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"0000000000000000000000",
			presetTokenBarcode,
			true,
			false,
		},
		{
			"00000000000000000000000",
			presetTokenBarcode,
			false,
			false,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("[%d]%s:\"%s\",regexp:%t,modck:%t", i+1, tc.preset, tc.code, tc.expectRegex, tc.expectCheck), func(t *testing.T) {
			pp, _ := presetsRegistry.get(tc.preset)
			rx := pp.regex
			matches := rx.MatchString(tc.code)
			require.Equal(t, tc.expectRegex, matches)
			if matches && pp.postChecker != nil {
				require.Equal(t, tc.expectCheck, pp.postChecker.Check(tc.code))
			}
		})
	}
}

func TestModuloCheck(t *testing.T) {
	moduloCheck := moduloCheck{
		modulo:        11,
		weights:       []int32{10, 9, 8, 7, 6, 5, 4, 3, 2},
		overflowDigit: "X",
	}
	result := moduloCheck.Check("020161622X")
	require.True(t, result)

	result = moduloCheck.Check("0201616229") // bad check digit
	require.False(t, result)

	result = moduloCheck.Check("02016162290") // too long
	require.False(t, result)

	result = moduloCheck.Check("020161622") // too short
	require.False(t, result)
}

func TestIsbn_Check(t *testing.T) {
	moduloCheck := isbn{}

	result := moduloCheck.Check("020161622X")
	require.True(t, result)

	result = moduloCheck.Check("9780201616224")
	require.True(t, result)

	result = moduloCheck.Check("0201616229") // bad check digit
	require.False(t, result)

	result = moduloCheck.Check("9780201616229") // bad check digit
	require.False(t, result)

	result = moduloCheck.Check("978020161622") // bad length
	require.False(t, result)
}

func TestIssn_Check(t *testing.T) {
	moduloCheck := issn{}

	result := moduloCheck.Check("1144875X")
	require.True(t, result)

	result = moduloCheck.Check("9771364743155")
	require.True(t, result)

	result = moduloCheck.Check("11448759") // bad check digit
	require.False(t, result)

	result = moduloCheck.Check("9771364743159") // bad check digit
	require.False(t, result)

	result = moduloCheck.Check("977136474315") // bad length
	require.False(t, result)
}

func TestPublication_Check(t *testing.T) {
	moduloCheck := publication{}

	result := moduloCheck.Check("1144875X")
	require.True(t, result)

	result = moduloCheck.Check("020161622X")
	require.True(t, result)

	result = moduloCheck.Check("9771364743155")
	require.True(t, result)

	result = moduloCheck.Check("11448759") // bad check digit
	require.False(t, result)

	result = moduloCheck.Check("0201616229") // bad check digit
	require.False(t, result)

	result = moduloCheck.Check("9771364743159") // bad check digit
	require.False(t, result)

	result = moduloCheck.Check("977136474315") // bad length
	require.False(t, result)
}

func TestUpc_Check(t *testing.T) {
	moduloCheck := upc{}

	result := moduloCheck.Check("01234565")
	require.True(t, result)

	result = moduloCheck.Check("123601057072")
	require.True(t, result)

	result = moduloCheck.Check("01234569") // bad check digit
	require.False(t, result)

	result = moduloCheck.Check("123601057079") // bad check digit
	require.False(t, result)

	result = moduloCheck.Check("12360105707") // bad length
	require.False(t, result)
}

func TestNumericWithScientific(t *testing.T) {
	// check that numeric scientific pattern matches everything that can be parsed using strconv.ParseFloat() - and vice versa
	regx := regexp.MustCompile(numericWithScientific)
	testCases := []struct {
		str string
		ok  bool
	}{
		{
			"1E0",
			true,
		},
		{
			"-1.234E10",
			true,
		},
		{
			"-1.23e-12",
			true,
		},
		{
			"-1.23e+12",
			true,
		},
		{
			"+10E0",
			true,
		},
		{
			"-10E0",
			true,
		},
		{
			"10E0",
			true,
		},
		{
			"2.3e4.5",
			false,
		},
		{
			"1",
			true,
		},
		{
			".5",
			true,
		},
		{
			"-.5",
			true,
		},
		{
			"+.5",
			true,
		},
		{
			"-.5e2",
			true,
		},
		{
			"+.5e2",
			true,
		},
		{
			"+.5e-2",
			true,
		},
		{
			"+.5E+2",
			true,
		},
		{
			"e2",
			false,
		},
		{
			"-e2",
			false,
		},
		{
			"NAN",
			false,
		},
		{
			"NaN",
			false,
		},
		{
			"INF",
			false,
		},
		{
			"inf",
			false,
		},
		{
			"iNf",
			false,
		},
		{
			"+Inf",
			false,
		},
		{
			"-Inf",
			false,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("[%d]\"%s\"", i+1, tc.str), func(t *testing.T) {
			f, err := strconv.ParseFloat(tc.str, 64)
			if tc.ok {
				require.Nil(t, err)
				require.True(t, regx.MatchString(tc.str))
			} else {
				if math.IsNaN(f) || math.IsInf(f, 0) {
					require.Nil(t, err)
				} else {
					require.NotNil(t, err)
				}
				require.False(t, regx.MatchString(tc.str))
			}
		})
	}
}

func TestNumericFullPattern(t *testing.T) {
	// check that numeric full pattern matches everything that can be parsed using strconv.ParseFloat() - and vice versa
	regx := regexp.MustCompile(numericFull)
	testCases := []struct {
		str string
		ok  bool
	}{
		{
			"1E0",
			true,
		},
		{
			"-1.234E10",
			true,
		},
		{
			"-1.23e-12",
			true,
		},
		{
			"-1.23e+12",
			true,
		},
		{
			"+10E0",
			true,
		},
		{
			"-10E0",
			true,
		},
		{
			"10E0",
			true,
		},
		{
			"2.3e4.5",
			false,
		},
		{
			"1",
			true,
		},
		{
			".5",
			true,
		},
		{
			"-.5",
			true,
		},
		{
			"+.5",
			true,
		},
		{
			"-.5e2",
			true,
		},
		{
			"+.5e2",
			true,
		},
		{
			"+.5e-2",
			true,
		},
		{
			"+.5E+2",
			true,
		},
		{
			"e2",
			false,
		},
		{
			"-e2",
			false,
		},
		{
			"NAN",
			true,
		},
		{
			"NaN",
			true,
		},
		{
			"INF",
			true,
		},
		{
			"inf",
			true,
		},
		{
			"iNf",
			true,
		},
		{
			"+Inf",
			true,
		},
		{
			"-Inf",
			true,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("[%d]\"%s\"", i+1, tc.str), func(t *testing.T) {
			_, err := strconv.ParseFloat(tc.str, 64)
			if tc.ok {
				require.Nil(t, err)
				require.True(t, regx.MatchString(tc.str))
			} else {
				require.NotNil(t, err)
				require.False(t, regx.MatchString(tc.str))
			}
		})
	}
}

func TestCardNumberPreset(t *testing.T) {
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
		"000000000000":        true,
		"0000000000000":       true,
		"00000000000000":      true,
		"000000000000000":     true,
		"0000000000000000":    true,
		"00000000000000000":   true,
		"000000000000000000":  true,
		"0000000000000000000": true,
		// invalid all zeroes...
		"00000000000":          false,
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
		"12345678901234x":         false,
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
	pp, _ := presetsRegistry.get("card")
	rx := pp.regex
	ck := pp.postChecker
	for ccn, expect := range testCardNumbers {
		t.Run(fmt.Sprintf("CardNumber:\"%s\"", ccn), func(t *testing.T) {
			if expect {
				require.True(t, rx.MatchString(ccn))
				require.True(t, ck.Check(ccn))
			} else if rx.MatchString(ccn) {
				require.False(t, ck.Check(ccn))
			}
		})
	}
}

func TestCmyks(t *testing.T) {
	checkCmyk := presetsRegistry.namedPresets[presetTokenCMYK]
	checkCmyk300 := presetsRegistry.namedPresets[presetTokenCMYK300]
	tcs := map[string]struct {
		expectMatch bool
		expect300Ok bool
	}{
		"cmyk(0,0,0,0)": {
			true,
			true,
		},
		"cmyk(0%,0%,0%,0%)": {
			true,
			true,
		},
		"cmyk(1,1,1,1)": {
			true,
			false,
		},
		"cmyk(1.0, 1.0, 1.0, 1.0)": {
			true,
			false,
		},
		"cmyk(1.1, 1.0, 1.0, 1.0)": {
			false,
			false,
		},
		"cmyk(100%, 100.0%, 100.0000%, 100.00000%)": {
			true,
			false,
		},
		"cmyk(100%, 100.0%, 100.0000%, 100.00009%)": {
			false,
			false,
		},
		"cmyk(100%, 100%, 100%, 0.1)": {
			true,
			false,
		},
		"cmyk(100%, 100%, 100%, 0.00001%)": {
			true,
			false,
		},
	}
	for str, tc := range tcs {
		t.Run(fmt.Sprintf("CMYK\"%s\"", str), func(t *testing.T) {
			require.Equal(t, tc.expectMatch, checkCmyk.regex.MatchString(str))
			require.Equal(t, tc.expectMatch, checkCmyk300.regex.MatchString(str))
			if tc.expectMatch {
				require.Equal(t, tc.expect300Ok, checkCmyk300.check(str))
				require.True(t, checkCmyk.check(str))
			}
		})
	}
}

func TestRgbIccs(t *testing.T) {
	checkRgbIcc := presetsRegistry.namedPresets[presetTokenRgbIcc]
	tcs := map[string]struct {
		expectMatch bool
		expectCheck bool
	}{
		"xxx": {
			false,
			false,
		},
		"rgb-icc(": {
			false,
			false,
		},
		"rgb-icc())": {
			true,
			false,
		},
		"rgb-icc(''')": {
			true,
			false,
		},
		"rgb-icc()": {
			true,
			false,
		},
		"rgb-icc(,)": {
			true,
			false,
		},
		"rgb-icc(255,255,255, #CMYK, 1, 0.1, 1, 1.0)": {
			true,
			true,
		},
		"rgb-icc(255,255, #CMYK, 1, 0.1, 1, 1.0)": {
			true,
			false,
		},
		"rgb-icc(256, 0, 0, #CMYK, 1, 1, 1, 1)": {
			true,
			false,
		},
		"rgb-icc(#CMYK, 1, 1, 1, 0)": {
			true,
			true,
		},
		"rgb-icc(#CMYK, 1.0, 0, 0, 0)": {
			true,
			true,
		},
		"rgb-icc(#CMYK, 1.5, 0, 0, 0)": {
			true,
			false,
		},
		"rgb-icc(#CMYK, 2, 0, 0, 0)": {
			true,
			false,
		},
		"rgb-icc(#CMYK, 0, 0, 0)": {
			true,
			false,
		},
		"rgb-icc(#Grayscale, 0.5)": {
			true,
			true,
		},
		"rgb-icc(#Grayscale)": {
			true,
			false,
		},
		"rgb-icc(#Grayscale, 1.1)": {
			true,
			false,
		},
		"rgb-icc(#Grayscale, -1)": {
			true,
			false,
		},
		"rgb-icc(#Grayscale, NaN)": {
			true,
			false,
		},
		"rgb-icc(#Grayscale, Inf)": {
			true,
			false,
		},
		"rgb-icc(#Grayscale, xxx)": {
			true,
			false,
		},
		"rgb-icc(#Grayscale, 'xxx')": {
			true,
			false,
		},
		"rgb-icc(128,128, 128,#Grayscale, 0.5)": {
			true,
			true,
		},
		"rgb-icc(128, 128, 128, #Grayscale, 0.5)": {
			true,
			true,
		},
		"rgb-icc(#Separation, 'Name')": {
			true,
			true,
		},
		"rgb-icc(#Separation, 'All')": {
			true,
			true,
		},
		"rgb-icc(#Separation)": {
			true,
			false,
		},
		"rgb-icc(#Separation, 1)": {
			true,
			false,
		},
		"rgb-icc(#Registration)": {
			true,
			true,
		},
		"rgb-icc(#Registration, 0.5)": {
			true,
			true,
		},
		"rgb-icc(#Registration, 1)": {
			true,
			true,
		},
		"rgb-icc(#Registration, 0)": {
			true,
			true,
		},
		"rgb-icc(#Registration, 0, 1)": {
			true,
			false,
		},
		"rgb-icc(#Registration, -0.5)": {
			true,
			false,
		},
		"rgb-icc(#Registration, 1.1)": {
			true,
			false,
		},
		"rgb-icc(#Registration, NaN)": {
			true,
			false,
		},
		"rgb-icc(#Registration, Inf)": {
			true,
			false,
		},
		"rgb-icc(128, 128, 128, #Registration, 0.5)": {
			true,
			true,
		},
		"rgb-icc(256, 128, 128, #Registration, 0.5)": {
			true,
			false,
		},
		"rgb-icc(1, 2, #Registration)": {
			true,
			false,
		},
		"rgb-icc(1, 2, #Registration, 0.5)": {
			true,
			false,
		},
		"rgb-icc(128, 128, 128, #Registration)": {
			true,
			true,
		},
		"rgb-icc(255,255,0, #SpotColor,'PANTONE Orange 021 C',0.33)": {
			true,
			true,
		},
		"rgb-icc(255,255,0, #SpotColor,'PANTONE 169 M',0.5, #CMYK,0,0.2,0.2,0)": {
			true,
			true,
		},
		"rgb-icc(255,255,0, #SpotColor,'PANTONE 169 M',0.5, #CMYK,0,0.2,0.2)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor,0.33, #Grayscale,0.5)": {
			true,
			true,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor,0.33, #Grayscale)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor,0.33, #Grayscale, 1.1)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor,0.33, #Grayscale, -0.5)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor,0.33, #Grayscale, NaN)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor,0.33, #Grayscale, Inf)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor,0.33, #Grayscale, xxx)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor,0.33, #Registration)": {
			true,
			true,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor,0.33, #Registration, 0.5)": {
			true,
			true,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor,0.33, #Registration, NaN)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor,0.33, #Registration, 0.5, 'what?')": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor,0.33, #SpotColor,'pp')": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor,)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor, NaN)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor, +Inf)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor, xxx)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,123, 1)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,'', 1)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor, 1.1)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor, 1.0, #Unknown)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor, 1.0, #CMYK, 1, 1, 0, x)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor, 1.0, #CMYK, 1, 1, 0, -1)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor, 1.0, #CMYK, 1, 1, 0, 1.1)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor, 1.0, #CMYK, 1, 1, 0, NaN)": {
			true,
			false,
		},
		"rgb-icc(255,255,0, #SpotColor,MyColor, 1.0, #CMYK, 1, 1, 0, inf)": {
			true,
			false,
		},
	}
	for str, tc := range tcs {
		t.Run(fmt.Sprintf("RGB-ICC\"%s\"", str), func(t *testing.T) {
			require.Equal(t, tc.expectMatch, checkRgbIcc.regex.MatchString(str))
			if tc.expectMatch {
				require.Equal(t, tc.expectCheck, checkRgbIcc.check(str))
			}
		})
	}
}
