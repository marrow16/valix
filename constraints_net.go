package valix

import (
	"net"
	"net/url"
	"strings"
)

const (
	msgValidMAC      = "String value must be a valid MAC address"
	msgValidCIDR     = "String value must be a valid CIDR address"
	msgValidCIDRv4   = "String value must be a valid CIDR (version 4) address"
	msgValidCIDRv6   = "String value must be a valid CIDR (version 6) address"
	msgValidIP       = "String must be a valid IP address"
	msgValidIPv4     = "String must be a valid IP (version 4) address"
	msgValidIPv6     = "String must be a valid IP (version 6) address"
	msgValidTCP      = "String must be a valid TCP address"
	msgValidTCPv4    = "String must be a valid TCP (version 4) address"
	msgValidTCPv6    = "String must be a valid TCP (version 6) address"
	msgValidUDP      = "String must be a valid UDP address"
	msgValidUDPv4    = "String must be a valid UDP (version 4) address"
	msgValidUDPv6    = "String must be a valid UDP (version 6) address"
	msgValidTld      = "String must be a valid TLD"
	msgValidHostname = "String must be a valid hostname"
	msgValidURI      = "String must be a valid URI"
	msgValidURL      = "String must be a valid URL"
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
	if str, ok := v.(string); ok && !(c.V4Only && c.V6Only) {
		pass := false
		if ip, _, err := net.ParseCIDR(str); err == nil {
			pass = isValidIp(ip, c.V4Only, c.V6Only, c.DisallowPrivate, c.DisallowLoopback)
		}
		if pass {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

func isValidIp(ip net.IP, v4Only, v6Only, disallowPrivate, disallowLoopback bool) bool {
	ipv4 := ip.To4()
	return ((!v4Only && !v6Only) || (v4Only && ipv4 != nil) || (v6Only && ipv4 == nil)) &&
		(!disallowPrivate || !ip.IsPrivate()) &&
		(!disallowLoopback || !((ipv4 != nil && ipv4[0] == 127) || (ipv4 == nil && ip.Equal(net.IPv6loopback))))
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

// NetIsHostname constraint to check that string value is a valid hostname
type NetIsHostname struct {
	// AllowIPAddress when set, allows IP address hostnames
	AllowIPAddress bool
	// AllowIPV6 when set, allows IP v6 address hostnames
	AllowIPV6 bool
	// AllowLocal when set, allows hostnames 'local' (e.g. "localhost", "local", "localdomain", "127.0.0.1", "::1")
	AllowLocal bool
	// AllowTldOnly when set, allows hostnames with only Tld specified (e.g. "audi")
	AllowTldOnly bool
	// AllowGeographicTlds when set, allows hostnames with geographic Tlds (e.g. "some-company.africa")
	AllowGeographicTlds bool
	// AllowGenericTlds when set, allows hostnames with generic Tlds (e.g. "some.academy")
	AllowGenericTlds bool
	// AllowBrandTlds when set, allows hostnames with brand Tlds (e.g. "my.audi")
	AllowBrandTlds bool
	// AllowInfraTlds when set, allows hostnames with infrastructure Tlds (e.g. "arpa")
	AllowInfraTlds bool
	// AllowTestTlds when set, allows hostnames with test Tlds and test domains (e.g. "example.com", "test.com")
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
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *NetIsHostname) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, true, c.Stop)
}

func (c *NetIsHostname) checkString(str string, vcx *ValidatorContext) bool {
	return len(str) > 0 && isValidDomain(str, domainOptions{
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
	})
}

// GetMessage implements the Constraint.GetMessage
func (c *NetIsHostname) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidHostname)
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
	if str, ok := v.(string); ok && !(c.V4Only && c.V6Only) {
		pass := false
		var ip net.IP = nil
		if c.Resolvable {
			if ipa, err := net.ResolveIPAddr(useNetwork("ip", c.V4Only, c.V6Only), str); err == nil {
				ip = ipa.IP
			}
		} else if c.AllowLocalhost && str == "localhost" {
			lhip := ternary(c.V6Only).string("::1", "127.0.0.1")
			ip = net.ParseIP(lhip)
		} else {
			ip = net.ParseIP(str)
		}
		if ip != nil {
			pass = isValidIp(ip, c.V4Only, c.V6Only, c.DisallowPrivate, c.DisallowLoopback)
		}
		if pass {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
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
	return checkStringConstraint(v, vcx, c, true, c.Stop)
}

func (c *NetIsMac) checkString(str string, vcx *ValidatorContext) bool {
	if _, err := net.ParseMAC(str); err == nil {
		return true
	}
	return false
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
	if str, ok := v.(string); ok && !(c.V4Only && c.V6Only) && str != "" {
		pass := false
		if ta, err := net.ResolveTCPAddr(useNetwork("tcp", c.V4Only, c.V6Only), str); err == nil {
			pass = isValidIp(ta.IP, c.V4Only, c.V6Only, c.DisallowPrivate, c.DisallowLoopback)
		}
		if pass {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
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
	return checkStringConstraint(v, vcx, c, true, c.Stop)
}

func (c *NetIsTld) checkString(str string, vcx *ValidatorContext) bool {
	return tld(str).isValid(domainOptions{
		allowGeographicTlds: c.AllowGeographicTlds,
		allowGenericTlds:    c.AllowGenericTlds,
		allowBrandTlds:      c.AllowBrandTlds,
		addCountryCodeTlds:  c.AddCountryCodeTlds,
		excCountryCodeTlds:  c.ExcCountryCodeTlds,
		addGenericTlds:      c.AddGenericTlds,
		excGenericTlds:      c.ExcGenericTlds,
		addBrandTlds:        c.AddBrandTlds,
		excBrandTlds:        c.ExcBrandTlds,
	})
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
	if str, ok := v.(string); ok && !(c.V4Only && c.V6Only) && str != "" {
		pass := false
		if ta, err := net.ResolveUDPAddr(useNetwork("udp", c.V4Only, c.V6Only), str); err == nil {
			pass = isValidIp(ta.IP, c.V4Only, c.V6Only, c.DisallowPrivate, c.DisallowLoopback)
		}
		if pass {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
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

// NetIsURI constraint to check that string value is a valid URN
type NetIsURI struct {
	// if set, the host is also checked (see also AllowIPAddress and others)
	CheckHost bool `v8n:"default"`
	// AllowIPAddress when set, allows IP address hostnames
	AllowIPAddress bool
	// AllowIPV6 when set, allows IP v6 address hostnames
	AllowIPV6 bool
	// AllowLocal when set, allows hostnames 'local' (e.g. "localhost", "local", "localdomain", "127.0.0.1", "::1")
	AllowLocal bool
	// AllowTldOnly when set, allows hostnames with only Tld specified (e.g. "audi")
	AllowTldOnly bool
	// AllowGeographicTlds when set, allows hostnames with geographic Tlds (e.g. "some-company.africa")
	AllowGeographicTlds bool
	// AllowGenericTlds when set, allows hostnames with generic Tlds (e.g. "some.academy")
	AllowGenericTlds bool
	// AllowBrandTlds when set, allows hostnames with brand Tlds (e.g. "my.audi")
	AllowBrandTlds bool
	// AllowInfraTlds when set, allows hostnames with infrastructure Tlds (e.g. "arpa")
	AllowInfraTlds bool
	// AllowTestTlds when set, allows hostnames with test Tlds and test domains (e.g. "example.com", "test.com")
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
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *NetIsURI) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, true, c.Stop)
}

func (c *NetIsURI) checkString(str string, vcx *ValidatorContext) bool {
	if hat := strings.Index(str, "#"); hat != -1 {
		str = str[:hat]
	}
	if len(str) > 0 {
		if u, err := url.ParseRequestURI(str); err == nil {
			if c.CheckHost {
				if isValidDomain(u.Hostname(), domainOptions{
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
					return true
				}
			} else {
				return true
			}
		}
	}
	return false
}

// GetMessage implements the Constraint.GetMessage
func (c *NetIsURI) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidURI)
}

// NetIsURL constraint to check that string value is a valid URL
type NetIsURL struct {
	// if set, the host is also checked (see also AllowIPAddress and others)
	CheckHost bool `v8n:"default"`
	// AllowIPAddress when set, allows IP address hostnames
	AllowIPAddress bool
	// AllowIPV6 when set, allows IP v6 address hostnames
	AllowIPV6 bool
	// AllowLocal when set, allows hostnames 'local' (e.g. "localhost", "local", "localdomain", "127.0.0.1", "::1")
	AllowLocal bool
	// AllowTldOnly when set, allows hostnames with only Tld specified (e.g. "audi")
	AllowTldOnly bool
	// AllowGeographicTlds when set, allows hostnames with geographic Tlds (e.g. "some-company.africa")
	AllowGeographicTlds bool
	// AllowGenericTlds when set, allows hostnames with generic Tlds (e.g. "some.academy")
	AllowGenericTlds bool
	// AllowBrandTlds when set, allows hostnames with brand Tlds (e.g. "my.audi")
	AllowBrandTlds bool
	// AllowInfraTlds when set, allows hostnames with infrastructure Tlds (e.g. "arpa")
	AllowInfraTlds bool
	// AllowTestTlds when set, allows hostnames with test Tlds and test domains (e.g. "example.com", "test.com")
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
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *NetIsURL) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, true, c.Stop)
}

func (c *NetIsURL) checkString(str string, vcx *ValidatorContext) bool {
	if hat := strings.Index(str, "#"); hat != -1 {
		str = str[:hat]
	}
	if len(str) > 0 {
		if u, err := url.ParseRequestURI(str); err == nil && u.Scheme != "" && u.Host != "" {
			if c.CheckHost {
				if isValidDomain(u.Hostname(), domainOptions{
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
					return true
				}
			} else {
				return true
			}
		}
	}
	return false
}

// GetMessage implements the Constraint.GetMessage
func (c *NetIsURL) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidURL)
}
