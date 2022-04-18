package valix

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"regexp"
	"strconv"
	"testing"
)

const builtInPresetsCount = 40

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

func TestCodesWithModuloChecked(t *testing.T) {
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
			presetTokenUPC,
			true,
			true,
		},
		{
			"5000237128237",
			presetTokenEAN13,
			true,
			true,
		},
		{
			"9771364743155",
			presetTokenISSN,
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
