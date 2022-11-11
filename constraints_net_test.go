package valix

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNetIsCIDR(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	testCases := map[string][]struct {
		constraint *NetIsCIDR
		expect     bool
	}{
		"": {
			{&NetIsCIDR{}, false},
		},
		"256.0.0.0/32": {
			{&NetIsCIDR{}, false},
		},
		"0.256.0.0/32": {
			{&NetIsCIDR{}, false},
		},
		"0.0.256.0/32": {
			{&NetIsCIDR{}, false},
		},
		"0.0.0.256/32": {
			{&NetIsCIDR{}, false},
		},
		"10.0.0.0": {
			{&NetIsCIDR{}, false},
		},
		"10.0.0.0/0": {
			{&NetIsCIDR{}, true},
		},
		"10.0.0.1/8": {
			{&NetIsCIDR{}, true},
		},
		"10.0.0.1/32": {
			{&NetIsCIDR{}, true},
		},
		"10.0.0.1/33": {
			{&NetIsCIDR{}, false},
		},

		"127.0.0.1/0": {
			{&NetIsCIDR{}, true},
			{&NetIsCIDR{DisallowPrivate: true}, true},
			{&NetIsCIDR{DisallowLoopback: true}, false},
			{&NetIsCIDR{V4Only: true}, true},
			{&NetIsCIDR{V4Only: true, DisallowLoopback: true}, false},
			{&NetIsCIDR{V6Only: true}, false},
			// setting both v4 and v6 always fails...
			{&NetIsCIDR{V4Only: true, V6Only: true}, false},
		},
		"::1/0": {
			{&NetIsCIDR{}, true},
			{&NetIsCIDR{DisallowPrivate: true}, true},
			{&NetIsCIDR{DisallowLoopback: true}, false},
			{&NetIsCIDR{V6Only: true}, true},
			{&NetIsCIDR{V6Only: true, DisallowLoopback: true}, false},
			{&NetIsCIDR{V4Only: true}, false},
			// setting both v4 and v6 always fails...
			{&NetIsCIDR{V4Only: true, V6Only: true}, false},
		},
		"10.255.0.0/0": {
			{&NetIsCIDR{}, true},
			{&NetIsCIDR{DisallowPrivate: true}, false},
		},
		"11.0.0.0/0": {
			{&NetIsCIDR{}, true},
			{&NetIsCIDR{DisallowPrivate: true}, true},
		},
		"fc00::/0": {
			{&NetIsCIDR{}, true},
			{&NetIsCIDR{DisallowPrivate: true}, false},
		},
		"fe00::/0": {
			{&NetIsCIDR{}, true},
			{&NetIsCIDR{DisallowPrivate: true}, true},
		},
		"192.0.2.1/0": {
			{&NetIsCIDR{}, true},
			{&NetIsCIDR{V4Only: true}, true},
			{&NetIsCIDR{V6Only: true}, false},
		},
		"2001:db8::68/64": {
			{&NetIsCIDR{}, true},
			{&NetIsCIDR{V6Only: true}, true},
			{&NetIsCIDR{V4Only: true}, false},
		},
		"2001:0db8:0000:0000:0000:0000:1234:5678/32": {
			{&NetIsCIDR{}, true},
			{&NetIsCIDR{V6Only: true}, true},
			{&NetIsCIDR{V4Only: true}, false},
		},
		"2001:db8:0:0:0:0:1234:5678/64": {
			{&NetIsCIDR{}, true},
			{&NetIsCIDR{V6Only: true}, true},
			{&NetIsCIDR{V4Only: true}, false},
		},
		"2001:0db8::1234:5678/64": {
			{&NetIsCIDR{}, true},
			{&NetIsCIDR{V6Only: true}, true},
			{&NetIsCIDR{V4Only: true}, false},
		},
	}
	for tv, tcs := range testCases {
		for i, tc := range tcs {
			c := tc.constraint
			t.Run(fmt.Sprintf("NetIsCIDR:\"%s\"[%d]V4:%v,V6:%v,!Lb:%v,!Pv:%v", tv, i+1, c.V4Only, c.V6Only, c.DisallowLoopback, c.DisallowPrivate), func(t *testing.T) {
				ok, msg := c.Check(tv, vcx)
				if tc.expect {
					require.True(t, ok)
				} else {
					require.False(t, ok)
					if c.V4Only {
						require.Equal(t, msgValidCIDRv4, msg)
					} else if c.V6Only {
						require.Equal(t, msgValidCIDRv6, msg)
					} else {
						require.Equal(t, msgValidCIDR, msg)
					}
				}
			})
		}
	}
}

func TestNetIsHostname(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	testCases := map[string][]struct {
		constraint *NetIsHostname
		expect     bool
	}{
		"": {
			{&NetIsHostname{}, false},
		},
		" ": {
			{&NetIsHostname{}, false},
		},
		".": {
			{&NetIsHostname{}, false},
		},
		"http://some.com": {
			{&NetIsHostname{}, false},
		},
		"b a d.com": {
			{&NetIsHostname{}, false},
		},
		"bad..com": {
			{&NetIsHostname{}, false},
		},
		".bad.com": {
			{&NetIsHostname{}, false},
		},
		"some.com": {
			{&NetIsHostname{}, true},
		},
		"//some.com": {
			{&NetIsHostname{}, false},
		},
		"/some.com": {
			{&NetIsHostname{}, false},
		},
		"example.com": {
			{&NetIsHostname{}, false},
			{&NetIsHostname{
				AllowTestTlds: true,
			}, true},
		},
		"example.unknown": {
			{&NetIsHostname{}, false},
		},
		"com": {
			{&NetIsHostname{}, false},
			{&NetIsHostname{
				AllowTldOnly: true,
			}, false},
		},
		"127.0.0.1": {
			{&NetIsHostname{}, false},
			{&NetIsHostname{
				AllowIPAddress: true,
			}, false},
			{&NetIsHostname{
				AllowIPAddress: true,
				AllowLocal:     true,
			}, true},
		},
		"127.0.0.1:8080": {
			{&NetIsHostname{}, false},
		},
		"101.102.103.104": {
			{&NetIsHostname{}, false},
			{&NetIsHostname{
				AllowIPAddress: true,
			}, true},
		},
		"101.102.103.104:80": {
			{&NetIsHostname{}, false},
		},
		"::1": {
			{&NetIsHostname{}, false},
			{&NetIsHostname{
				AllowIPAddress: true,
			}, false},
			{&NetIsHostname{
				AllowIPAddress: true,
				AllowIPV6:      true,
			}, false},
			{&NetIsHostname{
				AllowIPAddress: true,
				AllowIPV6:      true,
				AllowLocal:     true,
			}, true},
		},
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334": {
			{&NetIsHostname{}, false},
			{&NetIsHostname{
				AllowIPAddress: true,
			}, false},
			{&NetIsHostname{
				AllowIPAddress: true,
				AllowIPV6:      true,
			}, true},
		},
		"localhost": {
			{&NetIsHostname{}, false},
			{&NetIsHostname{
				AllowLocal: true,
			}, true},
		},
		"arpa": {
			{&NetIsHostname{}, false},
			{&NetIsHostname{
				AllowInfraTlds: true,
			}, true},
		},
		"usr:pwd@some.com": {
			{&NetIsHostname{}, false},
		},
		"some.com?querystr": {
			{&NetIsHostname{}, false},
		},
		"some.com#bookmark": {
			{&NetIsHostname{}, false},
		},
		"some.company.uk": {
			{&NetIsHostname{}, true},
			{&NetIsHostname{
				ExcCountryCodeTlds: []string{"uk"},
			}, false},
		},
		"some.company.africa": {
			{&NetIsHostname{}, false},
			{&NetIsHostname{
				AllowGeographicTlds: true,
			}, true},
			{&NetIsHostname{
				AllowGeographicTlds: true,
				ExcCountryCodeTlds:  []string{"africa"},
			}, false},
		},
		"some.company.accountants": {
			{&NetIsHostname{}, false},
			{&NetIsHostname{
				AllowGenericTlds: true,
			}, true},
			{&NetIsHostname{
				AllowGenericTlds: true,
				ExcGenericTlds:   []string{"accountants"},
			}, false},
		},
		"audi": {
			{&NetIsHostname{}, false},
			{&NetIsHostname{
				AllowBrandTlds: true,
			}, false},
			{&NetIsHostname{
				AllowBrandTlds: true,
				AllowTldOnly:   true,
			}, true},
			{&NetIsHostname{
				AllowBrandTlds: true,
				AllowTldOnly:   true,
				ExcBrandTlds:   []string{"audi"},
			}, false},
		},
		"tt.audi": {
			{&NetIsHostname{}, false},
			{&NetIsHostname{
				AllowBrandTlds: true,
			}, true},
			{&NetIsHostname{
				AllowBrandTlds: true,
				ExcBrandTlds:   []string{"audi"},
			}, false},
		},
		"goxxx.compare": {
			{&NetIsHostname{
				AllowGenericTlds: true,
			}, true},
			{&NetIsHostname{
				CheckHost:        true,
				AllowGenericTlds: true,
			}, false},
		},
		"go.compare": {
			{&NetIsHostname{
				AllowGenericTlds: true,
			}, true},
			{&NetIsHostname{
				CheckHost:        true,
				AllowGenericTlds: true,
			}, true},
		},
	}
	for a, tcs := range testCases {
		for i, tc := range tcs {
			c := tc.constraint
			t.Run(fmt.Sprintf("\"%s\"[%d]", a, i+1), func(t *testing.T) {
				ok, _ := c.Check(a, vcx)
				require.Equal(t, tc.expect, ok)
			})
		}
	}
}

func TestNetIsIP(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	testCases := map[string][]struct {
		constraint *NetIsIP
		expect     bool
	}{
		"": {
			{&NetIsIP{}, false},
		},
		"127.0.0.1:80": {
			{&NetIsIP{}, false},
		},
		"256.0.0.0": {
			{&NetIsIP{}, false},
		},
		"0.256.0.0": {
			{&NetIsIP{}, false},
		},
		"0.0.256.0": {
			{&NetIsIP{}, false},
		},
		"0.0.0.256": {
			{&NetIsIP{}, false},
		},
		"localhost": {
			{&NetIsIP{}, false},
			{&NetIsIP{AllowLocalhost: true}, true},
			{&NetIsIP{AllowLocalhost: true, Resolvable: true}, true},
			{&NetIsIP{AllowLocalhost: true, DisallowLoopback: true}, false},
			{&NetIsIP{AllowLocalhost: true, V6Only: true}, true},
			{&NetIsIP{AllowLocalhost: true, V6Only: true, DisallowLoopback: true}, false},
			{&NetIsIP{AllowLocalhost: true, DisallowPrivate: true}, true},
		},
		"127.0.0.1": {
			{&NetIsIP{}, true},
			{&NetIsIP{DisallowPrivate: true}, true},
			{&NetIsIP{DisallowPrivate: true, Resolvable: true}, true},
			{&NetIsIP{DisallowLoopback: true}, false},
			{&NetIsIP{V4Only: true}, true},
			{&NetIsIP{V4Only: true, DisallowLoopback: true}, false},
			{&NetIsIP{Resolvable: true}, true},
			{&NetIsIP{Resolvable: true, DisallowLoopback: true}, false},
			{&NetIsIP{V4Only: true, Resolvable: true}, true},
			{&NetIsIP{V4Only: true, Resolvable: true, DisallowLoopback: true}, false},
			{&NetIsIP{V6Only: true}, false},
			{&NetIsIP{V6Only: true, Resolvable: true}, false},
			{&NetIsIP{V6Only: true, Resolvable: true, DisallowLoopback: true}, false},
			// setting both v4 and v6 always fails...
			{&NetIsIP{V4Only: true, V6Only: true}, false},
		},
		"::1": {
			{&NetIsIP{}, true},
			{&NetIsIP{DisallowPrivate: true}, true},
			{&NetIsIP{DisallowPrivate: true, Resolvable: true}, true},
			{&NetIsIP{DisallowLoopback: true}, false},
			{&NetIsIP{V6Only: true}, true},
			{&NetIsIP{V6Only: true, DisallowLoopback: true}, false},
			{&NetIsIP{Resolvable: true}, true},
			{&NetIsIP{Resolvable: true, DisallowLoopback: true}, false},
			{&NetIsIP{V6Only: true, Resolvable: true}, true},
			{&NetIsIP{V6Only: true, Resolvable: true, DisallowLoopback: true}, false},
			{&NetIsIP{V4Only: true}, false},
			{&NetIsIP{V4Only: true, Resolvable: true}, false},
			{&NetIsIP{V4Only: true, Resolvable: true, DisallowLoopback: true}, false},
			// setting both v4 and v6 always fails...
			{&NetIsIP{V4Only: true, V6Only: true}, false},
		},
		"10.255.0.0": {
			{&NetIsIP{}, true},
			{&NetIsIP{DisallowPrivate: true}, false},
		},
		"11.0.0.0": {
			{&NetIsIP{}, true},
			{&NetIsIP{DisallowPrivate: true}, true},
		},
		"fc00::": {
			{&NetIsIP{}, true},
			{&NetIsIP{DisallowPrivate: true}, false},
		},
		"fe00::": {
			{&NetIsIP{}, true},
			{&NetIsIP{DisallowPrivate: true}, true},
		},
		"192.0.2.1": {
			{&NetIsIP{}, true},
			{&NetIsIP{V4Only: true}, true},
			{&NetIsIP{V6Only: true}, false},
		},
		"2001:db8::68": {
			{&NetIsIP{}, true},
			{&NetIsIP{V6Only: true}, true},
			{&NetIsIP{V4Only: true}, false},
		},
		"2001:0db8:0000:0000:0000:0000:1234:5678": {
			{&NetIsIP{}, true},
			{&NetIsIP{V6Only: true}, true},
			{&NetIsIP{V4Only: true}, false},
		},
		"2001:db8:0:0:0:0:1234:5678": {
			{&NetIsIP{}, true},
			{&NetIsIP{V6Only: true}, true},
			{&NetIsIP{V4Only: true}, false},
		},
		"2001:0db8::1234:5678": {
			{&NetIsIP{}, true},
			{&NetIsIP{V6Only: true}, true},
			{&NetIsIP{V4Only: true}, false},
		},
	}
	for tv, tcs := range testCases {
		for i, tc := range tcs {
			c := tc.constraint
			t.Run(fmt.Sprintf("NetIsIP:\"%s\"[%d]V4:%v,V6:%v,R:%v,!Lb:%v,!Pv:%v,ALh:%v", tv, i+1, c.V4Only, c.V6Only, c.Resolvable, c.DisallowLoopback, c.DisallowPrivate, c.AllowLocalhost), func(t *testing.T) {
				ok, msg := c.Check(tv, vcx)
				if tc.expect {
					require.True(t, ok)
				} else {
					require.False(t, ok)
					if c.V4Only {
						require.Equal(t, msgValidIPv4, msg)
					} else if c.V6Only {
						require.Equal(t, msgValidIPv6, msg)
					} else {
						require.Equal(t, msgValidIP, msg)
					}
				}
			})
		}
	}
}

func TestNetIsMac(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	c := &NetIsMac{}
	testCases := map[string]bool{
		"":                        false,
		"f00":                     false,
		"00:00:5e:00:53:01":       true,
		"02:00:5e:10:00:00:00:01": true,
		"00:00:00:00:fe:80:00:00:00:00:00:00:02:00:5e:10:00:00:00:01": true,
		"00-00-5e-00-53-01":       true,
		"02-00-5e-10-00-00-00-01": true,
		"00-00-00-00-fe-80-00-00-00-00-00-00-02-00-5e-10-00-00-00-01": true,
		"0000.5e00.5301":      true,
		"0200.5e10.0000.0001": true,
		"0000.0000.fe80.0000.0000.0000.0200.5e10.0000.0001": true,
	}
	for k, expect := range testCases {
		t.Run(fmt.Sprintf("NetIsMac\"%s\"", k), func(t *testing.T) {
			ok, msg := c.Check(k, vcx)
			if expect {
				require.True(t, ok)
			} else {
				require.False(t, ok)
				require.Equal(t, msgValidMAC, msg)
			}
		})
	}
}

func TestNetIsTCP(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	testCases := map[string][]struct {
		constraint *NetIsTCP
		expect     bool
	}{
		"": {
			{&NetIsTCP{}, false},
		},
		"256.0.0.0:80": {
			{&NetIsTCP{}, false},
		},
		"0.256.0.0:80": {
			{&NetIsTCP{}, false},
		},
		"0.0.256.0:80": {
			{&NetIsTCP{}, false},
		},
		"0.0.0.256:80": {
			{&NetIsTCP{}, false},
		},
		"127.0.0.1:65536": {
			{&NetIsTCP{}, false},
		},
		"127.0.0.1:": {
			{&NetIsTCP{}, false},
		},
		"127.0.0.1:80": {
			{&NetIsTCP{}, true},
			{&NetIsTCP{DisallowPrivate: true}, true},
			{&NetIsTCP{DisallowLoopback: true}, false},
			{&NetIsTCP{V4Only: true}, true},
			{&NetIsTCP{V4Only: true, DisallowLoopback: true}, false},
			{&NetIsTCP{V6Only: true}, false},
			// setting both v4 and v6 always fails...
			{&NetIsTCP{V4Only: true, V6Only: true}, false},
		},
		"[::1]": {
			{&NetIsTCP{}, false},
		},
		"[::1]:": {
			{&NetIsTCP{}, false},
		},
		"[::1]:80": {
			{&NetIsTCP{}, true},
			{&NetIsTCP{DisallowPrivate: true}, true},
			{&NetIsTCP{DisallowLoopback: true}, false},
			{&NetIsTCP{V6Only: true}, true},
			{&NetIsTCP{V6Only: true, DisallowLoopback: true}, false},
			{&NetIsTCP{V4Only: true}, false},
			// setting both v4 and v6 always fails...
			{&NetIsTCP{V4Only: true, V6Only: true}, false},
		},
	}
	for tv, tcs := range testCases {
		for i, tc := range tcs {
			c := tc.constraint
			t.Run(fmt.Sprintf("NetIsTCP:\"%s\"[%d]V4:%v,V6:%v,!Lb:%v,!Pv:%v", tv, i+1, c.V4Only, c.V6Only, c.DisallowLoopback, c.DisallowPrivate), func(t *testing.T) {
				ok, msg := c.Check(tv, vcx)
				if tc.expect {
					require.True(t, ok)
				} else {
					require.False(t, ok)
					if c.V4Only {
						require.Equal(t, msgValidTCPv4, msg)
					} else if c.V6Only {
						require.Equal(t, msgValidTCPv6, msg)
					} else {
						require.Equal(t, msgValidTCP, msg)
					}
				}
			})
		}
	}
}

func TestNetIsTld(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	testCases := map[string]struct {
		expect         bool
		tryWithOpts    bool
		expectWithOpts bool
		opts           NetIsTld
	}{
		"":    {},
		" ":   {},
		"com": {expect: true},
		"uk":  {expect: true},
		"abudhabi": {
			expect:         false,
			expectWithOpts: true,
			opts: NetIsTld{
				AllowGeographicTlds: true,
			},
		},
		"africa": {
			expect:      false,
			tryWithOpts: true,
			opts: NetIsTld{
				AllowGeographicTlds: true,
				ExcCountryCodeTlds:  []string{"africa"},
			},
		},
		"academy": {
			expect:         false,
			expectWithOpts: true,
			opts: NetIsTld{
				AllowGenericTlds: true,
			},
		},
		"accountant": {
			expect:      false,
			tryWithOpts: true,
			opts: NetIsTld{
				AllowGenericTlds: true,
				ExcGenericTlds:   []string{"accountant"},
			},
		},
		"audi": {
			expect:         false,
			expectWithOpts: true,
			opts: NetIsTld{
				AllowBrandTlds: true,
			},
		},
		"aws": {
			expect:      false,
			tryWithOpts: true,
			opts: NetIsTld{
				AllowBrandTlds: true,
				ExcBrandTlds:   []string{"aws"},
			},
		},
	}
	c := NetIsTld{}
	for k, tc := range testCases {
		t.Run(k, func(t *testing.T) {
			ok, _ := c.Check(k, vcx)
			require.Equal(t, tc.expect, ok)
			if tc.tryWithOpts || tc.expectWithOpts {
				ok, _ = tc.opts.Check(k, vcx)
				require.Equal(t, tc.expectWithOpts, ok)
			}
		})
	}
}

func TestNetIsUDP(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	testCases := map[string][]struct {
		constraint *NetIsUDP
		expect     bool
	}{
		"": {
			{&NetIsUDP{}, false},
		},
		"256.0.0.0:80": {
			{&NetIsUDP{}, false},
		},
		"0.256.0.0:80": {
			{&NetIsUDP{}, false},
		},
		"0.0.256.0:80": {
			{&NetIsUDP{}, false},
		},
		"0.0.0.256:80": {
			{&NetIsUDP{}, false},
		},
		"127.0.0.1:65536": {
			{&NetIsUDP{}, false},
		},
		"127.0.0.1:80": {
			{&NetIsUDP{}, true},
			{&NetIsUDP{DisallowPrivate: true}, true},
			{&NetIsUDP{DisallowLoopback: true}, false},
			{&NetIsUDP{V4Only: true}, true},
			{&NetIsUDP{V4Only: true, DisallowLoopback: true}, false},
			{&NetIsUDP{V6Only: true}, false},
			// setting both v4 and v6 always fails...
			{&NetIsUDP{V4Only: true, V6Only: true}, false},
		},
		"[::1]": {
			{&NetIsUDP{}, false},
		},
		"[::1]:80": {
			{&NetIsUDP{}, true},
			{&NetIsUDP{DisallowPrivate: true}, true},
			{&NetIsUDP{DisallowLoopback: true}, false},
			{&NetIsUDP{V6Only: true}, true},
			{&NetIsUDP{V6Only: true, DisallowLoopback: true}, false},
			{&NetIsUDP{V4Only: true}, false},
			// setting both v4 and v6 always fails...
			{&NetIsUDP{V4Only: true, V6Only: true}, false},
		},
	}
	for tv, tcs := range testCases {
		for i, tc := range tcs {
			c := tc.constraint
			t.Run(fmt.Sprintf("NetIsUDP:\"%s\"[%d]V4:%v,V6:%v,!Lb:%v,!Pv:%v", tv, i+1, c.V4Only, c.V6Only, c.DisallowLoopback, c.DisallowPrivate), func(t *testing.T) {
				ok, msg := c.Check(tv, vcx)
				if tc.expect {
					require.True(t, ok)
				} else {
					require.False(t, ok)
					if c.V4Only {
						require.Equal(t, msgValidUDPv4, msg)
					} else if c.V6Only {
						require.Equal(t, msgValidUDPv6, msg)
					} else {
						require.Equal(t, msgValidUDP, msg)
					}
				}
			})
		}
	}
}

func TestNetIsURI(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	testCases := map[string][]struct {
		constraint *NetIsURI
		expect     bool
	}{
		"": {
			{&NetIsURI{}, false},
		},
		" ": {
			{&NetIsURI{}, false},
		},
		".": {
			{&NetIsURI{}, false},
		},
		"http://": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
		},
		"http://b a d.com": {
			{&NetIsURI{}, false},
			{&NetIsURI{CheckHost: true}, false},
		},
		"http://bad..com": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
		},
		"http://.bad.com": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
		},
		"http://bad.com.": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
		},
		"some.com": {
			{&NetIsURI{}, false},
			{&NetIsURI{CheckHost: true}, false},
		},
		"//some.com": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
		},
		"/some.com": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
		},
		"http://example.com": {
			{&NetIsURI{}, true},
			{&NetIsURI{
				CheckHost: true,
			}, false},
			{&NetIsURI{
				CheckHost:     true,
				AllowTestTlds: true,
			}, true},
		},
		"http://example.unknown": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
		},
		"http://com": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
			{&NetIsURI{
				CheckHost:    true,
				AllowTldOnly: true,
			}, false},
		},
		"http://127.0.0.1": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowIPAddress: true,
			}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowIPAddress: true,
				AllowLocal:     true,
			}, true},
		},
		"http://127.0.0.1:8080": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowIPAddress: true,
			}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowIPAddress: true,
				AllowLocal:     true,
			}, true},
		},
		"http://101.102.103.104": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowIPAddress: true,
			}, true},
		},
		"http://101.102.103.104:80": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowIPAddress: true,
			}, true},
		},
		"https://[::1]": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowIPAddress: true,
			}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowIPAddress: true,
				AllowIPV6:      true,
			}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowIPAddress: true,
				AllowIPV6:      true,
				AllowLocal:     true,
			}, true},
		},
		"https://[::1]:8080": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowIPAddress: true,
			}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowIPAddress: true,
				AllowIPV6:      true,
			}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowIPAddress: true,
				AllowIPV6:      true,
				AllowLocal:     true,
			}, true},
		},
		"https://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:17000": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowIPAddress: true,
			}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowIPAddress: true,
				AllowIPV6:      true,
			}, true},
		},
		"http://localhost": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
			{&NetIsURI{
				CheckHost:  true,
				AllowLocal: true,
			}, true},
		},
		"http://arpa": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowInfraTlds: true,
			}, true},
		},
		"http://usr:pwd@some.com": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, true},
		},
		"http://some.com?querystr": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, true},
		},
		"http://some.com#bookmark": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, true},
		},
		"http://some.company.uk": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, true},
			{&NetIsURI{
				CheckHost:          true,
				ExcCountryCodeTlds: []string{"uk"},
			}, false},
		},
		"http://some.company.africa": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
			{&NetIsURI{
				CheckHost:           true,
				AllowGeographicTlds: true,
			}, true},
			{&NetIsURI{
				CheckHost:           true,
				AllowGeographicTlds: true,
				ExcCountryCodeTlds:  []string{"africa"},
			}, false},
		},
		"http://some.company.accountants": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
			{&NetIsURI{
				CheckHost:        true,
				AllowGenericTlds: true,
			}, true},
			{&NetIsURI{
				CheckHost:        true,
				AllowGenericTlds: true,
				ExcGenericTlds:   []string{"accountants"},
			}, false},
		},
		"http://audi": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowBrandTlds: true,
			}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowBrandTlds: true,
				AllowTldOnly:   true,
			}, true},
		},
		"http://tt.audi": {
			{&NetIsURI{}, true},
			{&NetIsURI{CheckHost: true}, false},
			{&NetIsURI{
				CheckHost:      true,
				AllowBrandTlds: true,
			}, true},
			{&NetIsURI{
				CheckHost:      true,
				AllowBrandTlds: true,
				ExcBrandTlds:   []string{"audi"},
			}, false},
		},
	}
	for a, tcs := range testCases {
		for i, tc := range tcs {
			c := tc.constraint
			t.Run(fmt.Sprintf("\"%s\"[%d]", a, i+1), func(t *testing.T) {
				ok, _ := c.Check(a, vcx)
				require.Equal(t, tc.expect, ok)
			})
		}
	}
}

func TestNetIsURL(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	testCases := map[string][]struct {
		constraint *NetIsURL
		expect     bool
	}{
		"": {
			{&NetIsURL{}, false},
		},
		" ": {
			{&NetIsURL{}, false},
		},
		".": {
			{&NetIsURL{}, false},
		},
		"http://": {
			{&NetIsURL{}, false},
			{&NetIsURL{CheckHost: true}, false},
		},
		"http://b a d.com": {
			{&NetIsURL{}, false},
			{&NetIsURL{CheckHost: true}, false},
		},
		"http://bad..com": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
		},
		"http://.bad.com": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
		},
		"http://bad.com.": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
		},
		"some.com": {
			{&NetIsURL{}, false},
			{&NetIsURL{CheckHost: true}, false},
		},
		"//some.com": {
			{&NetIsURL{}, false},
			{&NetIsURL{CheckHost: true}, false},
		},
		"/some.com": {
			{&NetIsURL{}, false},
			{&NetIsURL{CheckHost: true}, false},
		},
		"http://example.com": {
			{&NetIsURL{}, true},
			{&NetIsURL{
				CheckHost: true,
			}, false},
			{&NetIsURL{
				CheckHost:     true,
				AllowTestTlds: true,
			}, true},
		},
		"http://example.unknown": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
		},
		"http://com": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
			{&NetIsURL{
				CheckHost:    true,
				AllowTldOnly: true,
			}, false},
		},
		"http://127.0.0.1": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowIPAddress: true,
			}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowIPAddress: true,
				AllowLocal:     true,
			}, true},
		},
		"http://127.0.0.1:8080": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowIPAddress: true,
			}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowIPAddress: true,
				AllowLocal:     true,
			}, true},
		},
		"http://101.102.103.104": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowIPAddress: true,
			}, true},
		},
		"http://101.102.103.104:80": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowIPAddress: true,
			}, true},
		},
		"https://[::1]": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowIPAddress: true,
			}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowIPAddress: true,
				AllowIPV6:      true,
			}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowIPAddress: true,
				AllowIPV6:      true,
				AllowLocal:     true,
			}, true},
		},
		"https://[::1]:8080": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowIPAddress: true,
			}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowIPAddress: true,
				AllowIPV6:      true,
			}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowIPAddress: true,
				AllowIPV6:      true,
				AllowLocal:     true,
			}, true},
		},
		"https://[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:17000": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowIPAddress: true,
			}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowIPAddress: true,
				AllowIPV6:      true,
			}, true},
		},
		"http://localhost": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
			{&NetIsURL{
				CheckHost:  true,
				AllowLocal: true,
			}, true},
		},
		"http://arpa": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowInfraTlds: true,
			}, true},
		},
		"http://usr:pwd@some.com": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, true},
		},
		"http://some.com?querystr": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, true},
		},
		"http://some.com#bookmark": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, true},
		},
		"http://some.company.uk": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, true},
			{&NetIsURL{
				CheckHost:          true,
				ExcCountryCodeTlds: []string{"uk"},
			}, false},
		},
		"http://some.company.africa": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
			{&NetIsURL{
				CheckHost:           true,
				AllowGeographicTlds: true,
			}, true},
			{&NetIsURL{
				CheckHost:           true,
				AllowGeographicTlds: true,
				ExcCountryCodeTlds:  []string{"africa"},
			}, false},
		},
		"http://some.company.accountants": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
			{&NetIsURL{
				CheckHost:        true,
				AllowGenericTlds: true,
			}, true},
			{&NetIsURL{
				CheckHost:        true,
				AllowGenericTlds: true,
				ExcGenericTlds:   []string{"accountants"},
			}, false},
		},
		"http://audi": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowBrandTlds: true,
			}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowBrandTlds: true,
				AllowTldOnly:   true,
			}, true},
		},
		"http://tt.audi": {
			{&NetIsURL{}, true},
			{&NetIsURL{CheckHost: true}, false},
			{&NetIsURL{
				CheckHost:      true,
				AllowBrandTlds: true,
			}, true},
			{&NetIsURL{
				CheckHost:      true,
				AllowBrandTlds: true,
				ExcBrandTlds:   []string{"audi"},
			}, false},
		},
	}
	for a, tcs := range testCases {
		for i, tc := range tcs {
			c := tc.constraint
			t.Run(fmt.Sprintf("\"%s\"[%d]", a, i+1), func(t *testing.T) {
				ok, _ := c.Check(a, vcx)
				require.Equal(t, tc.expect, ok)
			})
		}
	}
}
