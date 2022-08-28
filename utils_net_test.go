package valix

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	testCases := map[string]struct {
		expectValid     bool
		tryWithOpts     bool
		expectOptsValid bool
		opts            domainOptions
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
			tryWithOpts: true,
			opts: domainOptions{
				allowTldOnly: true,
			},
		},
		"me@[127.0.0.1]": {
			expectValid:     false,
			expectOptsValid: true,
			opts: domainOptions{
				allowLocal:     true,
				allowIPAddress: true,
			},
		},
		"me@[123.456.789.0]": {
			expectValid:     false,
			expectOptsValid: false,
			tryWithOpts:     true,
			opts: domainOptions{
				allowIPAddress: true,
			},
		},
		"me@[123.45.67.89]": {
			expectValid:     false,
			expectOptsValid: true,
			opts: domainOptions{
				allowLocal:     true,
				allowIPAddress: true,
			},
		},
		"me@[::1]": {
			expectValid:     false,
			expectOptsValid: true,
			opts: domainOptions{
				allowLocal:     true,
				allowIPAddress: true,
				allowIPV6:      true,
			},
		},
		"me2@[::1]": {
			expectValid:     false,
			expectOptsValid: false,
			tryWithOpts:     true,
			opts: domainOptions{
				allowLocal:     true,
				allowIPAddress: true,
				allowIPV6:      false,
			},
		},
		"me@google.com": {
			expectValid: true,
		},
		"Me@Google.Com": {
			expectValid: true,
		},
		"me@localhost": {
			expectValid:     false,
			expectOptsValid: true,
			opts: domainOptions{
				allowLocal: true,
			},
		},
		"me@example.org": {
			expectValid:     false,
			expectOptsValid: true,
			opts: domainOptions{
				allowTestTlds: true,
			},
		},
		"me@audi": {
			expectValid:     false,
			expectOptsValid: true,
			opts: domainOptions{
				allowTldOnly:   true,
				allowBrandTlds: true,
			},
		},
		"me@sample.ಭಾರತ": {
			expectValid: true,
		},
		"me@sample.xn--2scrj9c": {
			expectValid: true,
		},
		"me@some.academy": {
			expectValid:     false,
			expectOptsValid: true,
			opts: domainOptions{
				allowGenericTlds: true,
			},
		},
		"me2@some.academy": {
			expectValid:     false,
			expectOptsValid: false,
			tryWithOpts:     true,
			opts: domainOptions{
				allowGenericTlds: true,
				excGenericTlds:   []string{"academy"},
			},
		},
	}
	defOpts := domainOptions{}
	for ea, tc := range testCases {
		t.Run(fmt.Sprintf("%s", ea), func(t *testing.T) {
			r := isValidEmail(ea, defOpts)
			require.Equal(t, tc.expectValid, r)
			if tc.tryWithOpts || tc.expectOptsValid {
				r = isValidEmail(ea, tc.opts)
				require.Equal(t, tc.expectOptsValid, r)
			}
		})
	}
}

func TestIsValidDomain(t *testing.T) {
	const longString = "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz0123456789A"
	require.Equal(t, 63, len(longString))
	longDomain := longString[:57] + "." + longString + "." + longString + "." + longString + ".com"
	require.Equal(t, 253, len(longDomain))

	testCases := map[string]struct {
		expectValid     bool
		tryWithOpts     bool
		expectOptsValid bool
		opts            domainOptions
	}{
		"apache.org":                    {expectValid: true},
		"www.google.com":                {expectValid: true},
		"test-domain.com":               {expectValid: true},
		"test---domain.com":             {expectValid: true},
		"test-d-o-m-ain.com":            {expectValid: true},
		"as.uk":                         {expectValid: true},
		"ApAchE.Org":                    {expectValid: true},
		"z.com":                         {expectValid: true},
		"i.have.an-example.domain.name": {},
		".org":                          {},
		" apache.org ":                  {},
		"apa che.org":                   {},
		"-testdomain.com":               {},
		"testdomain-.com":               {},
		"---c.com":                      {},
		"c--.com":                       {},
		"apache.rog":                    {},
		"http://www.apache.org":         {},
		"":                              {},
		" ":                             {},
		"a.ch":                          {expectValid: true},
		"9.ch":                          {expectValid: true},
		"az.ch":                         {expectValid: true},
		"09.ch":                         {expectValid: true},
		"9-1.ch":                        {expectValid: true},
		"91-.ch":                        {},
		"-.ch":                          {},
		"xn--d1abbgf6aiiy.xn--p1ai":     {expectValid: true},
		longString + ".com":             {expectValid: true},
		longString + "x.com":            {},
		longDomain:                      {expectValid: true},
		"x" + longDomain:                {},
		"example.com": {
			expectValid:     false,
			expectOptsValid: true,
			opts:            domainOptions{allowTestTlds: true},
		},
		"example": {
			expectValid:     false,
			expectOptsValid: true,
			opts: domainOptions{
				allowTestTlds: true,
				allowTldOnly:  true,
			},
		},
		"arpa": {
			expectValid:     false,
			expectOptsValid: true,
			opts: domainOptions{
				allowInfraTlds: true,
			},
		},
	}
	defOpts := domainOptions{}
	for d, tc := range testCases {
		t.Run(fmt.Sprintf("%s", d), func(t *testing.T) {
			r := isValidDomain(d, defOpts)
			require.Equal(t, tc.expectValid, r)
			if tc.tryWithOpts || tc.expectOptsValid {
				r = isValidDomain(d, tc.opts)
				require.Equal(t, tc.expectOptsValid, r)
			}
		})
	}
}

func TestDomainToASCIIDots(t *testing.T) {
	otherDots := map[string]string{
		"b" + fullStop:                           "b.",
		"b" + ideographicFullStop:                "b.",
		"b" + fullwidthFullStop:                  "b.",
		"b" + halfwidthIdeographicFullStop:       "b.",
		"a" + fullStop + "b":                     "a.b",
		"a" + ideographicFullStop + "b":          "a.b",
		"a" + fullwidthFullStop + "b":            "a.b",
		"a" + halfwidthIdeographicFullStop + "b": "a.b",
		fullStop:                                 ".",
		ideographicFullStop:                      ".",
		fullwidthFullStop:                        ".",
		halfwidthIdeographicFullStop:             ".",
	}
	for k, v := range otherDots {
		t.Run(k, func(t *testing.T) {
			s, err := domain(k).toASCII()
			require.Nil(t, err)
			require.Equal(t, v, s)
		})
	}
}

func TestTldsSortedForBinarySearch(t *testing.T) {
	for i := 1; i < len(countryCodeTlds); i++ {
		require.Equal(t, -1, strings.Compare(countryCodeTlds[i-1], countryCodeTlds[i]), "country: %s < %s", countryCodeTlds[i-1], countryCodeTlds[i])
	}
	for i := 1; i < len(geographicTlds); i++ {
		require.Equal(t, -1, strings.Compare(geographicTlds[i-1], geographicTlds[i]), "geographic: %s < %s", geographicTlds[i-1], geographicTlds[i])
	}
	for i := 1; i < len(genericTlds); i++ {
		require.Equal(t, -1, strings.Compare(genericTlds[i-1], genericTlds[i]), "generic: %s < %s", genericTlds[i-1], genericTlds[i])
	}
	for i := 1; i < len(brandTlds); i++ {
		require.Equal(t, -1, strings.Compare(brandTlds[i-1], brandTlds[i]), "brand: %s < %s", brandTlds[i-1], brandTlds[i])
	}
	for i := 1; i < len(terminatedBrandTlds); i++ {
		require.Equal(t, -1, strings.Compare(terminatedBrandTlds[i-1], terminatedBrandTlds[i]), "terminated brand: %s < %s", terminatedBrandTlds[i-1], terminatedBrandTlds[i])
	}
	for i := 1; i < len(withdrawnBrandTlds); i++ {
		require.Equal(t, -1, strings.Compare(withdrawnBrandTlds[i-1], withdrawnBrandTlds[i]), "withdrawn brand: %s < %s", withdrawnBrandTlds[i-1], withdrawnBrandTlds[i])
	}
	for i := 1; i < len(unknownTlds); i++ {
		require.Equal(t, -1, strings.Compare(unknownTlds[i-1], unknownTlds[i]), "unknown: %s < %s", unknownTlds[i-1], unknownTlds[i])
	}
}

func TestBrandGenericDupes(t *testing.T) {
	for _, s := range brandTlds {
		require.False(t, binarySearch(genericTlds, s))
		require.False(t, binarySearch(withdrawnBrandTlds, s))
		require.False(t, binarySearch(terminatedBrandTlds, s))
		require.False(t, binarySearch(unknownTlds, s))
	}
	for _, s := range genericTlds {
		require.False(t, binarySearch(brandTlds, s))
		require.False(t, binarySearch(withdrawnBrandTlds, s))
		require.False(t, binarySearch(terminatedBrandTlds, s))
		require.False(t, binarySearch(unknownTlds, s))
	}
}
