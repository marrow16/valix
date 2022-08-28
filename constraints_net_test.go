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
