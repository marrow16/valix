package valix

import (
	"net"
)

const (
	msgValidMAC    = "String value must be a valid MAC address"
	msgValidCIDR   = "String value must be a valid CIDR address"
	msgValidCIDRv4 = "String value must be a valid CIDR (version 4) address"
	msgValidCIDRv6 = "String value must be a valid CIDR (version 6) address"
	msgValidIP     = "String must be a valid IP address"
	msgValidIPv4   = "String must be a valid IP (version 4) address"
	msgValidIPv6   = "String must be a valid IP (version 6) address"
	msgValidTCP    = "String must be a valid TCP address"
	msgValidTCPv4  = "String must be a valid TCP (version 4) address"
	msgValidTCPv6  = "String must be a valid TCP (version 6) address"
	msgValidUDP    = "String must be a valid UDP address"
	msgValidUDPv4  = "String must be a valid UDP (version 4) address"
	msgValidUDPv6  = "String must be a valid UDP (version 6) address"
	msgValidTld    = "String must be a valid TLD"
)

// NetIsCIDR constraint to check that string value is a valid CIDR (v4 or v6) address
//
// NB. Setting both V4Only and V6Only to true will cause this constraint to always fail!
type NetIsCIDR struct {
	// if set, allows only CIDR v4
	V4Only bool
	// if set, allows only CIDR v6
	V6Only bool
	// if set, disallows loopback addresses
	DisallowLoopback bool
	// if set, disallows private addresses
	DisallowPrivate bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *NetIsCIDR) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		pass := false
		if !(c.V4Only && c.V6Only) { // can only pass if v4 and v6 are not both set
			if ip, _, err := net.ParseCIDR(str); err == nil {
				ipv4 := ip.To4()
				pass = (!c.V4Only && !c.V6Only) || (c.V4Only && ipv4 != nil) || (c.V6Only && ipv4 == nil)
				if pass && c.DisallowPrivate {
					pass = !ip.IsPrivate()
				}
				if pass && c.DisallowLoopback {
					pass = !((ipv4 != nil && ipv4[0] == 127) || (ipv4 == nil && ip.Equal(net.IPv6loopback)))
				}
			}
		}
		if !pass {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *NetIsCIDR) GetMessage(tcx I18nContext) string {
	if c.V4Only {
		return defaultMessage(tcx, c.Message, msgValidCIDRv4)
	} else if c.V6Only {
		return defaultMessage(tcx, c.Message, msgValidCIDRv6)
	}
	return defaultMessage(tcx, c.Message, msgValidCIDR)
}

// NetIsIP constraint to check that string value is a valid IP (v4 or v6) address
//
// NB. Setting both V4Only and V6Only to true will cause this constraint to always fail!
type NetIsIP struct {
	// if set, allows only IP v4
	V4Only bool
	// if set, allows only IP v6
	V6Only bool
	// if set, checks that the address is resolvable
	Resolvable bool
	// if set, disallows loopback addresses
	DisallowLoopback bool
	// if set, disallows private addresses
	DisallowPrivate bool
	// if set, allows value of "localhost" to be seen as valid
	AllowLocalhost bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *NetIsIP) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		pass := false
		if !(c.V4Only && c.V6Only) { // can only pass if v4 and v6 are not both set
			var ip net.IP = nil
			if c.Resolvable {
				if ipa, err := net.ResolveIPAddr(useNetwork("ip", c.V4Only, c.V6Only), str); err == nil {
					ip = ipa.IP
				}
			} else {
				if c.AllowLocalhost && str == "localhost" {
					if c.V6Only {
						str = "::1"
					} else {
						str = "127.0.0.1"
					}
				}
				ip = net.ParseIP(str)
			}
			if ip != nil {
				ipv4 := ip.To4()
				if ipv4 != nil {
					pass = !c.V6Only
				} else {
					pass = !c.V4Only
				}
				if pass && c.DisallowPrivate {
					pass = !ip.IsPrivate()
				}
				if pass && c.DisallowLoopback {
					pass = !((ipv4 != nil && ipv4[0] == 127) || (ipv4 == nil && ip.Equal(net.IPv6loopback)))
				}
			}
		}
		if !pass {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *NetIsIP) GetMessage(tcx I18nContext) string {
	if c.V4Only {
		return defaultMessage(tcx, c.Message, msgValidIPv4)
	} else if c.V6Only {
		return defaultMessage(tcx, c.Message, msgValidIPv6)
	}
	return defaultMessage(tcx, c.Message, msgValidIP)
}

// NetIsMac constraint to check that string value is a valid MAC address
type NetIsMac struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *NetIsMac) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if _, err := net.ParseMAC(str); err != nil {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *NetIsMac) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidMAC)
}

// NetIsTCP constraint to check that string value is a valid resolvable TCP (v4 or v6) address
//
// NB. Setting both V4Only and V6Only to true will cause this constraint to always fail!
type NetIsTCP struct {
	// if set, allows only TCP v4
	V4Only bool
	// if set, allows only TCP v6
	V6Only bool
	// if set, disallows loopback addresses
	DisallowLoopback bool
	// if set, disallows private addresses
	DisallowPrivate bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *NetIsTCP) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		pass := false
		if !(c.V4Only && c.V6Only) && str != "" { // can only pass if v4 and v6 are not both set
			if ta, err := net.ResolveTCPAddr(useNetwork("tcp", c.V4Only, c.V6Only), str); err == nil {
				ip := ta.IP
				ipv4 := ip.To4()
				pass = (!c.V4Only && !c.V6Only) || (c.V4Only && ipv4 != nil) || (c.V6Only && ipv4 == nil)
				if pass && c.DisallowPrivate {
					pass = !ip.IsPrivate()
				}
				if pass && c.DisallowLoopback {
					pass = !((ipv4 != nil && ipv4[0] == 127) || (ipv4 == nil && ip.Equal(net.IPv6loopback)))
				}
			}
		}
		if !pass {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *NetIsTCP) GetMessage(tcx I18nContext) string {
	if c.V4Only {
		return defaultMessage(tcx, c.Message, msgValidTCPv4)
	} else if c.V6Only {
		return defaultMessage(tcx, c.Message, msgValidTCPv6)
	}
	return defaultMessage(tcx, c.Message, msgValidTCP)
}

// NetIsTld constraint to check that string value is a valid Tld (top level domain)
type NetIsTld struct {
	AllowGeographicTlds bool
	AllowGenericTlds    bool
	AllowBrandTlds      bool
	AddCountryCodeTlds  []string
	ExcCountryCodeTlds  []string
	AddGenericTlds      []string
	ExcGenericTlds      []string
	AddBrandTlds        []string
	ExcBrandTlds        []string
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *NetIsTld) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if !tld(str).isValid(domainOptions{
			allowGeographicTlds: c.AllowGeographicTlds,
			allowGenericTlds:    c.AllowGenericTlds,
			allowBrandTlds:      c.AllowBrandTlds,
			addCountryCodeTlds:  c.AddCountryCodeTlds,
			excCountryCodeTlds:  c.ExcCountryCodeTlds,
			addGenericTlds:      c.AddGenericTlds,
			excGenericTlds:      c.ExcGenericTlds,
			addBrandTlds:        c.AddBrandTlds,
			excBrandTlds:        c.ExcBrandTlds,
		}) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *NetIsTld) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidTld)
}

// NetIsUDP constraint to check that string value is a valid resolvable UDP (v4 or v6) address
//
// NB. Setting both V4Only and V6Only to true will cause this constraint to always fail!
type NetIsUDP struct {
	// if set, allows only TCP v4
	V4Only bool
	// if set, allows only TCP v6
	V6Only bool
	// if set, disallows loopback addresses
	DisallowLoopback bool
	// if set, disallows private addresses
	DisallowPrivate bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *NetIsUDP) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		pass := false
		if !(c.V4Only && c.V6Only) && str != "" { // can only pass if v4 and v6 are not both set
			if ta, err := net.ResolveUDPAddr(useNetwork("udp", c.V4Only, c.V6Only), str); err == nil {
				ip := ta.IP
				ipv4 := ip.To4()
				pass = (!c.V4Only && !c.V6Only) || (c.V4Only && ipv4 != nil) || (c.V6Only && ipv4 == nil)
				if pass && c.DisallowPrivate {
					pass = !ip.IsPrivate()
				}
				if pass && c.DisallowLoopback {
					pass = !((ipv4 != nil && ipv4[0] == 127) || (ipv4 == nil && ip.Equal(net.IPv6loopback)))
				}
			}
		}
		if !pass {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *NetIsUDP) GetMessage(tcx I18nContext) string {
	if c.V4Only {
		return defaultMessage(tcx, c.Message, msgValidUDPv4)
	} else if c.V6Only {
		return defaultMessage(tcx, c.Message, msgValidUDPv6)
	}
	return defaultMessage(tcx, c.Message, msgValidUDP)
}

func useNetwork(addr string, v4, v6 bool) string {
	if v4 {
		return addr + "4"
	} else if v6 {
		return addr + "6"
	}
	return addr
}
