package valix

import (
	"golang.org/x/net/idna"
	"net"
	"regexp"
	"strings"
)

const (
	maxDomainLength         = 253
	maxUsernameLength       = 64
	domainLabelPattern      = "[a-zA-Z0-9]{1}[a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9]?"
	topLabelPattern         = "[a-zA-Z]{1}[a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9]?"
	emailPattern            = "^(.+)@(\\S+)$"
	userSpecialCharsPattern = "\\x00-\\x1F\\x7F\\(\\)<>@,;:'\\\\\\\"\\.\\[\\]"
	userValidCharsPattern   = "(\\\\.)|[^\\s" + userSpecialCharsPattern + "]"
	userQuotedPattern       = "(\"(\\\\\"|[^\"])*\")"
	userWordPattern         = "((" + userValidCharsPattern + "|')+|" + userQuotedPattern + ")"
	ipDomainPattern         = "^\\[(.*)\\]$"
)

var (
	topLabelRegex    = regexp.MustCompile("^(" + topLabelPattern + ")$")
	domainLabelRegex = regexp.MustCompile("^(" + domainLabelPattern + ")$")
	ipDomainRegex    = regexp.MustCompile(ipDomainPattern)
	emailRegex       = regexp.MustCompile(emailPattern)
	userRegex        = regexp.MustCompile("^" + userWordPattern + "(\\." + userWordPattern + ")*$")
)

type domainOptions struct {
	allowIPAddress      bool
	allowIPV6           bool
	allowLocal          bool
	allowTldOnly        bool
	allowGeographicTlds bool
	allowGenericTlds    bool
	allowBrandTlds      bool
	allowInfraTlds      bool
	allowTestTlds       bool
	addCountryCodeTlds  []string
	excCountryCodeTlds  []string
	addGenericTlds      []string
	excGenericTlds      []string
	addBrandTlds        []string
	excBrandTlds        []string
	addLocalTlds        []string
	excLocalTlds        []string
}

func tldIn(tld string, sl []string) bool {
	for _, v := range sl {
		if v == tld {
			return true
		}
	}
	return false
}

func isValidEmail(addr string, opts domainOptions) bool {
	if strings.HasSuffix(addr, ".") {
		return false
	}
	parts := emailRegex.FindStringSubmatch(addr)
	if parts == nil {
		return false
	}
	if !isValidUser(parts[1]) {
		return false
	}
	return isValidEmailDomain(parts[2], opts)
}

func isValidUser(user string) bool {
	if len(user) > maxUsernameLength {
		return false
	}
	return userRegex.MatchString(user)
}

func isValidEmailDomain(d string, opts domainOptions) bool {
	if ipParts := ipDomainRegex.FindStringSubmatch(d); ipParts != nil {
		return opts.allowIPAddress && inetAddress(ipParts[1]).isValid(opts.allowLocal, opts.allowIPV6)
	}
	return domain(d).isValid(opts)
}

func isValidDomain(d string, opts domainOptions) bool {
	return domain(d).isValid(opts)
}

type domain string

func (d domain) isValid(opts domainOptions) (result bool) {
	result = false
	if da, err := d.toASCII(); err == nil {
		if len(da) > maxDomainLength {
			return
		}
		parts := strings.Split(strings.ToLower(da), ".")
		l := len(parts)
		top := parts[l-1]
		if strings.HasSuffix(top, "-") || !topLabelRegex.MatchString(top) {
			return
		}
		if l == 1 {
			result = tld(top).isValidInfrastructure(opts) ||
				tld(top).isValidLocal(opts) ||
				(opts.allowTldOnly && (tld(top).isValidBrand(opts) ||
					tld(top).isValidTest(opts)))
		} else {
			for i := 0; i < (l - 1); i++ {
				if strings.HasSuffix(parts[i], "-") || !domainLabelRegex.MatchString(parts[i]) {
					return
				}
			}
			if testTlds[tld(parts[l-2])] && originalTlds[tld(top)] {
				result = opts.allowTestTlds
			} else {
				result = tld(top).isValid(opts)
			}
		}
	}
	return
}

const (
	fullStop                     = "\u002E"
	ideographicFullStop          = "\u3002"
	fullwidthFullStop            = "\uFF0E"
	halfwidthIdeographicFullStop = "\uFF61"
)

func (d domain) toASCII() (string, error) {
	if d.isOnlyASCII() {
		return string(d), nil
	}
	dots := strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(string(d), halfwidthIdeographicFullStop, fullStop), fullwidthFullStop, fullStop), ideographicFullStop, fullStop)
	return idna.ToASCII(dots)
}

func (d domain) isOnlyASCII() bool {
	for _, ch := range string(d) {
		if ch > 127 {
			return false
		}
	}
	return true
}

type tld string

func (t tld) isValid(opts domainOptions) bool {
	return originalTlds[t] ||
		t.isValidLocal(opts) ||
		t.isValidInfrastructure(opts) ||
		t.isValidGeneric(opts) ||
		t.isValidBrand(opts) ||
		t.isValidCountryOrGeographic(opts) ||
		t.isValidTest(opts)
}

func (t tld) isValidLocal(opts domainOptions) bool {
	if opts.allowLocal {
		return (localTlds[t] || tldIn(string(t), opts.addLocalTlds)) && !tldIn(string(t), opts.excLocalTlds)
	}
	return false
}

func (t tld) isValidInfrastructure(opts domainOptions) bool {
	return opts.allowInfraTlds && infrastructureTlds[t]
}

func (t tld) isValidGeneric(opts domainOptions) bool {
	if opts.allowGenericTlds {
		return (binarySearch(genericTlds, string(t)) || tldIn(string(t), opts.addGenericTlds)) && !tldIn(string(t), opts.excGenericTlds)
	}
	return false
}

func (t tld) isValidBrand(opts domainOptions) bool {
	if opts.allowBrandTlds {
		return (binarySearch(brandTlds, string(t)) || tldIn(string(t), opts.addBrandTlds)) && !tldIn(string(t), opts.excBrandTlds)
	}
	return false
}

func (t tld) isValidCountryOrGeographic(opts domainOptions) bool {
	return ((binarySearch(countryCodeTlds, string(t)) ||
		tldIn(string(t), opts.addCountryCodeTlds)) ||
		(opts.allowGeographicTlds && binarySearch(geographicTlds, string(t)))) &&
		!tldIn(string(t), opts.excCountryCodeTlds)
}

func (t tld) isValidTest(opts domainOptions) bool {
	return opts.allowTestTlds && testTlds[t]
}

type inetAddress string

func (ia inetAddress) isValid(allowLocal bool, allowIPV6 bool) bool {
	ip := net.ParseIP(string(ia))
	if ip == nil {
		return false
	}
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
	} else if !allowIPV6 {
		return false
	}
	return allowLocal || !ip.IsLoopback()
}

var localTlds = map[tld]bool{
	"local":       true,
	"localdomain": true,
	"localhost":   true,
}

var testTlds = map[tld]bool{
	"example": true,
	"invalid": true,
	"onion":   true,
	"test":    true,
}

var infrastructureTlds = map[tld]bool{
	"arpa": true,
}

var originalTlds = map[tld]bool{
	"com": true,
	"edu": true,
	"gov": true,
	"int": true,
	"mil": true,
	"net": true,
	"org": true,
}

// countryCodeTlds must be sorted (for binary search to work)
var countryCodeTlds = []string{
	"ac",                     // Ascension Island (United Kingdom)
	"ad",                     // Andorra
	"ae",                     // United Arab Emirates
	"af",                     // Afghanistan
	"ag",                     // Antigua and Barbuda
	"ai",                     // Anguilla (United Kingdom)
	"al",                     // Albania
	"am",                     // Armenia
	"ao",                     // Angola
	"aq",                     // Antarctica
	"ar",                     // Argentina
	"as",                     // American Samoa (United States)
	"at",                     // Austria
	"au",                     // Australia
	"aw",                     // Aruba (Kingdom of the Netherlands)
	"ax",                     // Åland (Finland)
	"az",                     // Azerbaijan
	"ba",                     // Bosnia and Herzegovina
	"bb",                     // Barbados
	"bd",                     // Bangladesh
	"be",                     // Belgium
	"bf",                     // Burkina Faso
	"bg",                     // Bulgaria
	"bh",                     // Bahrain
	"bi",                     // Burundi
	"bj",                     // Benin
	"bm",                     // Bermuda (United Kingdom)
	"bn",                     // Brunei
	"bo",                     // Bolivia
	"br",                     // Brazil
	"bs",                     // Bahamas
	"bt",                     // Bhutan
	"bw",                     // Botswana
	"by",                     // Belarus
	"bz",                     // Belize
	"ca",                     // Canada
	"cc",                     // Cocos (Keeling) Islands
	"cd",                     // Democratic Republic of the Congo
	"cf",                     // Central African Republic
	"cg",                     // Republic of the Congo
	"ch",                     // Switzerland
	"ci",                     // Ivory Coast
	"ck",                     // Cook Islands
	"cl",                     // Chile
	"cm",                     // Cameroon
	"cn",                     // People's Republic of China
	"co",                     // Colombia
	"cr",                     // Costa Rica
	"cu",                     // Cuba
	"cv",                     // Cape Verde
	"cw",                     // Curaçao (Kingdom of the Netherlands)
	"cx",                     // Christmas Island
	"cy",                     // Cyprus
	"cz",                     // Czech Republic
	"de",                     // Germany
	"dj",                     // Djibouti
	"dk",                     // Denmark
	"dm",                     // Dominica
	"do",                     // Dominican Republic
	"dz",                     // Algeria
	"ec",                     // Ecuador
	"ee",                     // Estonia
	"eg",                     // Egypt
	"er",                     // Eritrea
	"es",                     // Spain
	"et",                     // Ethiopia
	"eu",                     // European Union
	"fi",                     // Finland
	"fj",                     // Fiji
	"fk",                     // Falkland Islands (United Kingdom)
	"fm",                     // Federated States of Micronesia
	"fo",                     // Faroe Islands (Kingdom of Denmark)
	"fr",                     // France
	"ga",                     // Gabon
	"gd",                     // Grenada
	"ge",                     // Georgia
	"gf",                     // French Guiana (France)
	"gg",                     // Guernsey
	"gh",                     // Ghana
	"gi",                     // Gibraltar (United Kingdom)
	"gl",                     // Greenland (Kingdom of Denmark)
	"gm",                     // The Gambia
	"gn",                     // Guinea
	"gp",                     // Guadeloupe (France)
	"gq",                     // Equatorial Guinea
	"gr",                     // Greece
	"gs",                     // South Georgia and the South Sandwich Islands (United Kingdom)
	"gt",                     // Guatemala
	"gu",                     // Guam (United States)
	"gw",                     // Guinea-Bissau
	"gy",                     // Guyana
	"hk",                     // Hong Kong
	"hm",                     // Heard Island and McDonald Islands
	"hn",                     // Honduras
	"hr",                     // Croatia
	"ht",                     // Haiti
	"hu",                     // Hungary
	"id",                     // Indonesia
	"ie",                     // Ireland
	"il",                     // Israel
	"im",                     // Isle of Man
	"in",                     // India
	"io",                     // British Indian Ocean Territory (United Kingdom)
	"iq",                     // Iraq
	"ir",                     // Iran
	"is",                     // Iceland
	"it",                     // Italy
	"je",                     // Jersey
	"jm",                     // Jamaica
	"jo",                     // Jordan
	"jp",                     // Japan
	"ke",                     // Kenya
	"kg",                     // Kyrgyzstan
	"kh",                     // Cambodia
	"ki",                     // Kiribati
	"km",                     // Comoros
	"kn",                     // Saint Kitts and Nevis
	"kp",                     // North Korea
	"kr",                     // South Korea
	"kw",                     // Kuwait
	"ky",                     // Cayman Islands (United Kingdom)
	"kz",                     // Kazakhstan
	"la",                     // Laos
	"lb",                     // Lebanon
	"lc",                     // Saint Lucia
	"li",                     // Liechtenstein
	"lk",                     // Sri Lanka
	"lr",                     // Liberia
	"ls",                     // Lesotho
	"lt",                     // Lithuania
	"lu",                     // Luxembourg
	"lv",                     // Latvia
	"ly",                     // Libya
	"ma",                     // Morocco
	"mc",                     // Monaco
	"md",                     // Moldova
	"me",                     // Montenegro
	"mg",                     // Madagascar
	"mh",                     // Marshall Islands
	"mk",                     // North Macedonia
	"ml",                     // Mali
	"mm",                     // Myanmar
	"mn",                     // Mongolia
	"mo",                     // Macau
	"mp",                     // Northern Mariana Islands (United States)
	"mq",                     // Martinique (France)
	"mr",                     // Mauritania
	"ms",                     // Montserrat (United Kingdom)
	"mt",                     // Malta
	"mu",                     // Mauritius
	"mv",                     // Maldives
	"mw",                     // Malawi
	"mx",                     // Mexico
	"my",                     // Malaysia
	"mz",                     // Mozambique
	"na",                     // Namibia
	"nc",                     // New Caledonia (France)
	"ne",                     // Niger
	"nf",                     // Norfolk Island
	"ng",                     // Nigeria
	"ni",                     // Nicaragua
	"nl",                     // Netherlands
	"no",                     // Norway
	"np",                     // Nepal
	"nr",                     // Nauru
	"nu",                     // Niue
	"nz",                     // New Zealand
	"om",                     // Oman
	"pa",                     // Panama
	"pe",                     // Peru
	"pf",                     // French Polynesia (France)
	"pg",                     // Papua New Guinea
	"ph",                     // Philippines
	"pk",                     // Pakistan
	"pl",                     // Poland
	"pm",                     // Saint-Pierre and Miquelon (France)
	"pn",                     // Pitcairn Islands (United Kingdom)
	"pr",                     // Puerto Rico (United States)
	"ps",                     // Palestine[51]
	"pt",                     // Portugal
	"pw",                     // Palau
	"py",                     // Paraguay
	"qa",                     // Qatar
	"re",                     // Réunion (France)
	"ro",                     // Romania
	"rs",                     // Serbia
	"ru",                     // Russia
	"rw",                     // Rwanda
	"sa",                     // Saudi Arabia
	"sb",                     // Solomon Islands
	"sc",                     // Seychelles
	"sd",                     // Sudan
	"se",                     // Sweden
	"sg",                     // Singapore
	"sh",                     // Saint Helena, Ascension and Tristan da Cunha (United Kingdom)
	"si",                     // Slovenia
	"sk",                     // Slovakia
	"sl",                     // Sierra Leone
	"sm",                     // San Marino
	"sn",                     // Senegal
	"so",                     // Somalia
	"sr",                     // Suriname
	"ss",                     // South Sudan
	"st",                     // São Tomé and Príncipe
	"su",                     // Soviet Union
	"sv",                     // El Salvador
	"sx",                     // Sint Maarten (Kingdom of the Netherlands)
	"sy",                     // Syria
	"sz",                     // Eswatini
	"tc",                     // Turks and Caicos Islands (United Kingdom)
	"td",                     // Chad
	"tf",                     // French Southern and Antarctic Lands
	"tg",                     // Togo
	"th",                     // Thailand
	"tj",                     // Tajikistan
	"tk",                     // Tokelau
	"tl",                     // East Timor
	"tm",                     // Turkmenistan
	"tn",                     // Tunisia
	"to",                     // Tonga
	"tr",                     // Turkey
	"tt",                     // Trinidad and Tobago
	"tv",                     // Tuvalu
	"tw",                     // Taiwan
	"tz",                     // Tanzania
	"ua",                     // Ukraine
	"ug",                     // Uganda
	"uk",                     // United Kingdom
	"us",                     // United States of America
	"uy",                     // Uruguay
	"uz",                     // Uzbekistan
	"va",                     // Vatican City
	"vc",                     // Saint Vincent and the Grenadines
	"ve",                     // Venezuela
	"vg",                     // British Virgin Islands (United Kingdom)
	"vi",                     // United States Virgin Islands (United States)
	"vn",                     // Vietnam
	"vu",                     // Vanuatu
	"wf",                     // Wallis and Futuna
	"ws",                     // Samoa
	"xn--2scrj9c",            // ಭಾರತ (India)
	"xn--3e0b707e",           // 한국 (South Korea)
	"xn--3hcrj9c",            // ଭାରତ (India)
	"xn--45br5cyl",           // ভাৰত (India)
	"xn--45brj9c",            // ভারত (India)
	"xn--4dbrk0ce",           // ישראל (Israel)
	"xn--54b7fta0cc",         // বাংলা (Bangladesh)
	"xn--80ao21a",            // қаз (Kazakhstan)
	"xn--90a3ac",             // срб (Serbia)
	"xn--90ae",               // бг (Bulgaria)
	"xn--90ais",              // бел (Belarus)
	"xn--clchc0ea0b2g2a9gcd", // சிங்கப்பூர் (Singapore)
	"xn--d1alf",              // мкд (North Macedonia)
	"xn--e1a4c",              // ею (European Union)
	"xn--fiqs8s",             // 中国 (China)
	"xn--fiqz9s",             // 中國 (China)
	"xn--fpcrj9c3d",          // భారత్ (India)
	"xn--fzc2c9e2c",          // ලංකා (Sri Lanka)
	"xn--gecrj9c",            // ભારત (India)
	"xn--h2breg3eve",         // भारतम् (India)
	"xn--h2brj9c",            // भारत (India)
	"xn--h2brj9c8c",          // भारोत (India)
	"xn--j1amh",              // укр (Ukraine)
	"xn--j6w193g",            // 香港 (Hong Kong)
	"xn--kprw13d",            // 台湾 (Taiwan)
	"xn--kpry57d",            // 台灣 (Taiwan)
	"xn--l1acc",              // мон (Mongolia)
	"xn--lgbbat1ad8j",        // الجزائر (Algeria)
	"xn--mgb9awbf",           // عمان (Oman)
	"xn--mgba3a4f16a",        // ایران (Iran)
	"xn--mgbaam7a8h",         // امارات (United Arab Emirates)
	"xn--mgbah1a3hjkrd",      // موريتانيا (Mauritania)
	"xn--mgbai9azgqp6j",      // پاکستان (Pakistan)
	"xn--mgbayh7gpa",         // الاردن (Jordan)
	"xn--mgbbh1a",            // بارت (India)
	"xn--mgbbh1a71e",         // بھارت (India)
	"xn--mgbc0a9azcg",        // المغرب (Morocco)
	"xn--mgbcpq6gpa1a",       // البحرين (Bahrain)
	"xn--mgberp4a5d4ar",      // السعودية (Saudi Arabia)
	"xn--mgbgu82a",           // ڀارت (India)
	"xn--mgbpl2fh",           // سودان (Sudan)
	"xn--mgbx4cd0ab",         // مليسيا (Malaysia)
	"xn--mix891f",            // 澳門 (Macao)
	"xn--node",               // გე (Georgia)
	"xn--o3cw4h",             // ไทย (Thailand)
	"xn--ogbpf8fl",           // سورية (Syria)
	"xn--p1ai",               // рф (Russia)
	"xn--pgbs0dh",            // تونس (Tunisia)
	"xn--q7ce6a",             // ລາວ (Laos)
	"xn--qxa6a",              // ευ (European Union)
	"xn--qxam",               // ελ (Greece)
	"xn--rvc1e0am3e",         // ഭാരതം (India)
	"xn--s9brj9c",            // ਭਾਰਤ (India)
	"xn--wgbh1c",             // مصر (Egypt)
	"xn--wgbl6a",             // قطر (Qatar)
	"xn--xkc2al3hye2a",       // இலங்கை (Sri Lanka)
	"xn--xkc2dl3a5ee0h",      // இந்தியா (India)
	"xn--y9a3aq",             // հայ (Armenia)
	"xn--yfro4i67o",          // 新加坡 (Singapore)
	"xn--ygbi2ammx",          // فلسطين (Palestinian Authority)
	"ye",                     // Yemen
	"yt",                     // Mayotte
	"za",                     // South Africa
	"zm",                     // Zambia
	"zw",                     // Zimbabwe
}

// geographicTlds must be sorted (for binary search to work)
var geographicTlds = []string{
	"abudhabi",   // Abu Dhabi
	"africa",     // Africa
	"alsace",     //  Alsace
	"amsterdam",  //  Amsterdam, Netherlands
	"arab",       // League of Arab States
	"asia",       // Asia-Pacific region
	"bar",        // ZastavaBar.png Bar, Montenegro
	"barcelona",  // Barcelona
	"bayern",     //  Bavaria
	"bcn",        // Barcelona
	"berlin",     //  Berlin
	"boston",     //  Boston, Massachusetts
	"brussels",   //  Brussels, Belgium
	"budapest",   //  Budapest, Hungary
	"bzh",        //  Brittany; Breton language and culture
	"capetown",   // Cape Town, South Africa
	"cat",        //  Catalonia; Catalan language and culture
	"cologne",    //  Cologne
	"corsica",    //  Corsica
	"cymru",      //  Wales, United Kingdom
	"doha",       // Doha
	"dubai",      // Dubai
	"durban",     // Durban, South Africa
	"eus",        // Basque, Spain and France
	"frl",        //  Friesland, Netherlands
	"gal",        //  Galicia
	"gent",       // Ghent, Belgium
	"hamburg",    //  Hamburg
	"helsinki",   // Helsinki, Finland
	"irish",      //  Ireland; global Irish community
	"ist",        // İstanbul, Turkey
	"istanbul",   // İstanbul, Turkey
	"joburg",     // Johannesburg, South Africa
	"kiwi",       // New Zealand New Zealanders
	"koeln",      //  Cologne
	"krd",        //  Kurdistan
	"kyoto",      // Kyoto, Japan
	"lat",        // Latin America
	"london",     // London, United Kingdom
	"madrid",     //  Madrid
	"melbourne",  //  Melbourne, Australia
	"miami",      //  Miami, Florida
	"moscow",     //  Moscow, Russia
	"nagoya",     // Nagoya, Japan
	"nrw",        //  North Rhine-Westphalia
	"nyc",        //  New York City, New York
	"okinawa",    // Okinawa, Japan
	"osaka",      // Osaka, Japan
	"paris",      //  Paris
	"quebec",     //  Quebec, Canada
	"rio",        // Rio de Janeiro, Brazil
	"ruhr",       // Ruhr
	"ryukyu",     // Ryukyu Islands, Japan
	"saarland",   //  Saarland
	"scot",       //  Scotland, United Kingdom
	"stockholm",  //  Stockholm, Sweden
	"swiss",      //   Switzerland
	"sydney",     //  Sydney, Australia
	"taipei",     // Taipei, Taiwan
	"tatar",      // Tatar peoples and places
	"tirol",      //  Tyrol, Austria
	"tokyo",      // Tokyo, Japan
	"vegas",      //  Las Vegas, Nevada
	"vlaanderen", //  Flanders, Belgium
	"wales",      //  Wales, United Kingdom
	"wien",       //  Vienna, Austria
	"xn--1qqw23a",
	"xn--80adxhks",
	"xn--mgbca7dzdo",
	"xn--xhq521b",
	"yokohama", // Yokohama, Japan
	"zuerich",  //  Zurich, Switzerland
}

// genericTlds must be sorted (for binary search to work)
var genericTlds = []string{
	"abogado",
	"academy",
	"accountant",
	"accountants",
	"actor",
	"ads",
	"adult",
	"aero",
	"africa",
	"agency",
	"airforce",
	"analytics",
	"anquan",
	"apartments",
	"app",
	"archi",
	"army",
	"art",
	"associates",
	"attorney",
	"auction",
	"audio",
	"author",
	"auto",
	"autos",
	"baby",
	"band",
	"bank",
	"bar",
	"bargains",
	"baseball",
	"basketball",
	"beauty",
	"beer",
	"best",
	"bestbuy",
	"bet",
	"bible",
	"bid",
	"bike",
	"bingo",
	"bio",
	"biz",
	"black",
	"blackfriday",
	"blockbuster",
	"blog",
	"blue",
	"boats",
	"boo",
	"book",
	"booking", // Booking.com
	"boston",
	"bot",
	"boutique",
	"box",
	"broadway",
	"broker",
	"build",
	"builders",
	"business",
	"buy",
	"buzz",
	"cab",
	"cafe",
	"call",
	"cam",
	"camera",
	"camp",
	"cancerresearch",
	"capital",
	"car",
	"cards",
	"care",
	"career",
	"careers",
	"cars",
	"casa",
	"case",
	"cash",
	"casino",
	"catering",
	"catholic",
	"center",
	"cfd",
	"channel",
	"charity",
	"chat",
	"cheap",
	"christmas",
	"church",
	"circle",
	"city",
	"cityeats",
	"claims",
	"cleaning",
	"click",
	"clinic",
	"clothing",
	"cloud",
	"club",
	"coach",
	"codes",
	"coffee",
	"college",
	"community",
	"company",
	"compare",
	"computer",
	"condos",
	"construction",
	"consulting",
	"contact",
	"contractors",
	"cooking",
	"cool",
	"coop",
	"country",
	"coupon",
	"coupons",
	"courses",
	"cpa",
	"credit",
	"creditcard",
	"creditunion", // CUNA Performance Resources
	"cricket",
	"cruise",
	"cruises",
	"cyou",
	"dad",
	"dance",
	"data",
	"date",
	"dating",
	"day",
	"dds",
	"deal",
	"deals",
	"degree",
	"delivery",
	"democrat",
	"dental",
	"dentist",
	"desi",
	"design",
	"dev",
	"diamonds",
	"diet",
	"digital",
	"direct",
	"directory",
	"discount",
	"diy",
	"docs",
	"doctor",
	"dog",
	"domains",
	"dot",
	"download",
	"drive",
	"earth",
	"eat",
	"eco",
	"education",
	"email",
	"energy",
	"engineer",
	"engineering",
	"enterprises",
	"equipment",
	"esq",
	"estate",
	"events",
	"exchange",
	"expert",
	"exposed",
	"express",
	"fail",
	"faith",
	"family",
	"fan",
	"fans",
	"farm",
	"fashion",
	"fast",
	"feedback",
	"film",
	"final",
	"finance",
	"financial",
	"fire",
	"fish",
	"fishing",
	"fit",
	"fitness",
	"flights",
	"florist",
	"flowers",
	"fly",
	"foo",
	"food",
	"foodnetwork",
	"football",
	"forsale",
	"forum",
	"foundation",
	"free",
	"frontdoor",
	"fun",
	"fund",
	"furniture",
	"futbol",
	"fyi",
	"gallery",
	"game",
	"games",
	"garden",
	"gay",
	"gdn",
	"gift",
	"gifts",
	"gives",
	"giving", // Giving Limited
	"glass",
	"global",
	"gmbh",
	"gold",
	"golf",
	"gop",
	"got",
	"graphics",
	"gratis",
	"green",
	"gripe",
	"grocery",
	"group",
	"guide",
	"guitars",
	"guru",
	"hair",
	"hangout",
	"haus",
	"health",
	"healthcare",
	"help",
	"here",
	"hiphop",
	"hiv",
	"hockey",
	"holdings",
	"holiday",
	"homegoods",
	"homes",
	"homesense",
	"horse",
	"hospital",
	"host",
	"hosting",
	"hot",
	"hoteles",
	"hotels",
	"house",
	"how",
	"ice",
	"icu",
	"imamat",
	"immo",
	"immobilien",
	"inc",
	"industries",
	"info",
	"ing",
	"ink",
	"institute",
	"insurance",
	"insure",
	"international",
	"investments",
	"irish",
	"ismaili",
	"java", // Oracle (Java)
	"jetzt",
	"jewelry",
	"jobs",
	"jot",
	"joy",
	"juegos",
	"kaufen",
	"kids",
	"kim",
	"kitchen",
	"kosher",
	"kpn",
	"land",
	"lat",
	"latino",
	"law",
	"lawyer",
	"lds", // The Church of Jesus Christ of Latter-day Saints (LDS Church)[186]
	"lease",
	"legal",
	"lgbt",
	"life",
	"lifeinsurance",
	"lifestyle", // Lifestyle Domain Holdings, Inc.
	"lighting",
	"like",
	"limited",
	"limo",
	"link",
	"live",
	"living",
	"llc",
	"llp",
	"loan",
	"loans",
	"locker",
	"lol",
	"lotto",
	"love",
	"ltd",
	"ltda",
	"luxe",
	"luxury",
	"maison",
	"makeup",
	"management",
	"map",
	"market",
	"marketing",
	"markets",
	"mba",
	"med",
	"media",
	"meet",
	"meme",
	"memorial",
	"men",
	"menu",
	"mint",
	"mls",
	"mobi",
	"mobile",
	"moda",
	"moe",
	"moi",
	"mom",
	"money",
	"monster",
	"mormon", // The Church of Jesus Christ of Latter-day Saints (LDS Church)[186]
	"mortgage",
	"motorcycles",
	"mov",
	"movie",
	"movistar", // Telefónica (Movistar)
	"museum",
	"music",
	"name",
	"navy",
	"network",
	"new",
	"news",
	"ngo",
	"ninja",
	"now",
	"nowruz",
	"observer",
	"one",
	"ong",
	"onl",
	"online",
	"ooo",
	"open",
	"org",
	"organic",
	"origins",
	"ott",
	"page",
	"pars",
	"partners",
	"parts",
	"party",
	"passagens",
	"pay",
	"pet",
	"pharmacy",
	"phd",
	"phone",
	"photo",
	"photography",
	"photos",
	"physio",
	"pics",
	"pictures",
	"pid",
	"pin",
	"pink",
	"pizza",
	"place",
	"play", // Google (Google Play)
	"plumbing",
	"plus",
	"poker",
	"politie", // Politie Nederland
	"porn",
	"post",
	"press",
	"prime",
	"pro",
	"productions",
	"prof",
	"promo",
	"properties",
	"property",
	"protection",
	"pub",
	"qpon",
	"quebec",
	"racing",
	"radio",
	"read",
	"realestate",
	"realtor",
	"realty",
	"recipes",
	"red",
	"rehab",
	"reise",
	"reisen",
	"reit",
	"rent",
	"rentals",
	"repair",
	"report",
	"republican",
	"rest",
	"restaurant",
	"review",
	"reviews",
	"rich",
	"rip",
	"rocks",
	"rodeo",
	"room",
	"rsvp",
	"rugby",
	"run",
	"safe",
	"safety", // Safety Registry Services, LLC.
	"sale",
	"salon",
	"sarl",
	"save",
	"scholarships",
	"school",
	"schule",
	"science",
	"search",
	"secure",
	"security",
	"select",
	"services",
	"sex",
	"sexy",
	"shia",
	"shiksha",
	"shoes",
	"shop",
	"shopping",
	"show",
	"showtime",
	"silk",
	"singles",
	"site",
	"ski",
	"skin",
	"sling",
	"smile",
	"soccer",
	"social",
	"software",
	"solar",
	"solutions",
	"song",
	"soy",
	"spa",
	"space",
	"sport",
	"spot",
	"spreadbetting",
	"srl",
	"stockholm",
	"storage",
	"store",
	"stream",
	"studio",
	"study",
	"style",
	"sucks",
	"supplies",
	"supply",
	"support",
	"surf",
	"surgery",
	"systems",
	"talk",
	"tattoo",
	"tax",
	"taxi",
	"team",
	"tech",
	"technology",
	"tel",
	"tennis",
	"theater",
	"theatre",
	"tickets",
	"tienda",
	"tips",
	"tires",
	"today",
	"tools",
	"top",
	"tours",
	"town",
	"toys",
	"trade",
	"trading",
	"training",
	"travel",
	"travelersinsurance",
	"trust",
	"tube",
	"tunes",
	"university",
	"uno",
	"vacations",
	"ventures",
	"versicherung",
	"vet",
	"viajes",
	"video",
	"villas",
	"vin",
	"vip",
	"vision",
	"vivo", // Telefónica Brasil
	"vodka",
	"vote",
	"voting",
	"voto",
	"voyage",
	"vuelos",
	"wang",
	"wanggou",
	"watch",
	"watches",
	"weather",
	"webcam",
	"website",
	"wed",
	"wedding",
	"weibo",
	"whoswho",
	"wiki",
	"win",
	"wine",
	"winners",
	"work",
	"works",
	"world",
	"wow",
	"wtf",
	"xihuan",
	"xin",
	"xn--11b4c3d",
	"xn--1ck2e1b",
	"xn--30rr7y",
	"xn--3bst00m",
	"xn--3ds443g",
	"xn--3pxu8k",
	"xn--42c2d9a",
	"xn--45q11c",
	"xn--4gbrim",
	"xn--55qw42g",
	"xn--55qx5d",
	"xn--5tzm5g",
	"xn--6frz82g",
	"xn--6qq986b3xl",
	"xn--80aqecdr1a",
	"xn--80asehdb",
	"xn--80aswg",
	"xn--9dbq2a",
	"xn--9et52u",
	"xn--9krt00a",
	"xn--bck1b9a5dre4c",
	"xn--c1avg",
	"xn--c2br7g",
	"xn--cck2b3b",
	"xn--czr694b",
	"xn--czrs0t",
	"xn--czru2d",
	"xn--d1acj3b",
	"xn--efvy88h",
	"xn--fct429k",
	"xn--fhbei",
	"xn--fiq228c5hs",
	"xn--fjq720a",
	"xn--g2xx48c",
	"xn--gk3at1e",
	"xn--hxt814e",
	"xn--i1b6b1a6a2e",
	"xn--imr513n",
	"xn--io0a7i",
	"xn--j1aef",
	"xn--jvr189m",
	"xn--kpu716f", // 手表
	"xn--kput3i",
	"xn--mgbab2bd",
	"xn--mgbt3dhd",
	"xn--mk1bu44c",
	"xn--mxtq1m",
	"xn--ngbc5azd",
	"xn--nqv7f",
	"xn--nqv7fs00ema",
	"xn--nyqy26a",
	"xn--p1acf",
	"xn--pbt977c", // 珠宝
	"xn--pssy2u",
	"xn--q9jyb4c",
	"xn--rhqv96g",
	"xn--rovu88b",
	"xn--ses554g",
	"xn--t60b56a",
	"xn--tckwe",
	"xn--tiq49xqyj",
	"xn--unup4y",
	"xn--vermgensberater-ctb",
	"xn--vermgensberatung-pwb",
	"xn--vhquv",
	"xn--vuq861b",
	"xn--zfr164b",
	"xxx",
	"xyz",
	"yachts",
	"yoga",
	"you",
	"yun",
	"zero",
	"zip",
	"zone",
}

// brandTlds must be sorted (for binary search to work)
var brandTlds = []string{
	"aaa",    // American Automobile Association
	"aarp",   // AARP
	"abarth", // Fiat Chrysler Automobiles (Abarth)
	"abb",    // ABB
	"abbott", // Abbott Laboratories
	"abbvie", // AbbVie
	"abc",    // American Broadcasting Company[146]
	"able",
	"accenture", // Accenture
	"aco",       // ACO Severin Ahlmann GmbH & Co. KG
	"active",
	"adac",
	"aeg",             // Aktiebolaget Electrolux
	"aetna",           // Aetna
	"afl",             // Australian Football League
	"agakhan",         // Aga Khan Foundation
	"aig",             // American International Group
	"aigo",            // Aigo
	"airbus",          // Airbus
	"airtel",          // Bharti Airtel
	"akdn",            // Aga Khan Foundation
	"alfaromeo",       // Fiat Chrysler Automobiles (Alfa Romeo)
	"alibaba",         // Alibaba Group
	"alipay",          // Alibaba Group
	"allfinanz",       // Allfinanz Deutsche Vermögensberatung Aktiengesellschaft
	"allstate",        // Allstate
	"ally",            // Ally Financial
	"alstom",          // Alstom
	"amazon",          // Amazon
	"americanexpress", // American Express
	"americanfamily",
	"amex", // American Express
	"amfam",
	"amica",     // Amica Mutual Insurance
	"android",   // Google (Android)
	"anz",       // Australia & New Zealand Banking Group
	"aol",       // AOL
	"apple",     // Apple
	"aquarelle", // Aquarelle.com Group
	"aramco",    // Aramco
	"arte",
	"asda",
	"athleta",
	"audi", // Audi
	"audible",
	"auspost", // Australia Post
	"avianca",
	"aws",   // Amazon Web Services
	"axa",   // Axa
	"azure", // Microsoft
	"baidu", // Baidu
	"banamex",
	"bananarepublic", // Gap (Banana Republic)
	"barclays",       // Barclays
	"barefoot",
	"bauhaus", // Werkhaus GmbH
	"bbc",     // BBC
	"bbt",     // BB&T
	"bbva",    // Banco Bilbao Vizcaya Argentaria
	"bcg",     // Boston Consulting Group
	"beats",
	"bentley",    // Bentley
	"bharti",     // Bharti Enterprises
	"bing",       // Microsoft (Bing)
	"blanco",     // BLANCO GmbH + Co KG
	"bloomberg",  // Bloomberg
	"bms",        // Bristol-Myers Squibb
	"bmw",        // BMW
	"bnl",        // Banca Nazionale del Lavoro
	"bnpparibas", // BNP Paribas
	"boehringer", // Boehringer Ingelheim
	"bofa",
	"bom",
	"bond", // Bond University[8]
	"boots",
	"bosch",       // Robert Bosch GmbH
	"bostik",      // Bostik
	"bradesco",    // Bradesco
	"bridgestone", // Bridgestone
	"brother",     // Brother Industries
	"bugatti",     // Bugatti
	"cal",         // Google (Google Calendar)
	"calvinklein", // PVH
	"canon",       // Canon
	"capitalone",  // Capital One
	"caravan",     // Caravan International
	"cartier",     // Richemont DNS
	"cba",         // Commonwealth Bank
	"cbn",         // Christian Broadcasting Network
	"cbre",        // CBRE Group
	"cbs",         // CBS
	"ceo",
	"cern",    // CERN
	"cfa",     // CFA Institute
	"chanel",  // Chanel
	"chase",   // Chase Bank
	"chintai", // Chintai Corporation
	"chloe",
	"chrome",   // Google (Google Chrome)
	"chrysler", // Chrysler
	"cipriani", // Hotel Cipriani Srl
	"cisco",    // Cisco Systems
	"citadel",  // Citadel Domain
	"citi",     // Citigroup
	"citic",    // CITIC Group
	"clinique",
	"clubmed",  // Club Med
	"comcast",  // Comcast
	"commbank", // Commonwealth Bank
	"comsec",
	"cookingchannel",
	"crown",      // Crown Equipment Corporation
	"crs",        // Federated Co-operatives (Co-Operative Retailing System)
	"csc",        // DXC Technology
	"cuisinella", // Société Alsacienne de Meubles (Cuisinella)
	"dabur",      // Dabur
	"datsun",     // Nissan (Datsun)
	"dclk",
	"dealer",   // Dealer Dot Com
	"dell",     // Dell
	"deloitte", // Deloitte
	"delta",    // Delta Air Lines
	"dhl",      // Deutsche Post
	"discover", // Discover Financial Services
	"dish",     // Dish Network
	"dnp",      // Dai Nippon Printing
	"dodge",    // Chrysler (Dodge)
	"dtv",
	"duck",
	"dunlop", // Goodyear Tire and Rubber Company (Dunlop Tyres)
	"dupont", // DuPont
	"dvag",   // Deutsche Vermögensberatung
	"dvr",
	"edeka",  // Edeka
	"emerck", // Merck Group
	"epost",
	"epson",      // Seiko Epson
	"ericsson",   // Ericsson
	"erni",       // ERNI Group Holding AG
	"esurance",   // Esurance
	"etisalat",   // Etisalat
	"eurovision", // European Broadcasting Union (Eurovision)
	"everbank",   // EverBank
	"extraspace", // Extra Space Storage
	"fage",       // Fage
	"fairwinds",  // FairWinds Partners
	"farmers",    // Farmers Insurance Exchange
	"fedex",      // FedEx
	"ferrari",    // Fiat Chrysler Automobiles (Ferrari)
	"ferrero",    // Ferrero
	"fiat",       // Fiat Chrysler Automobiles
	"fidelity",   // Fidelity Investments
	"fido",
	"firestone", // Bridgestone
	"firmdale",  // Firmdale Holdings
	"flickr",    // Yahoo! (Flickr)
	"flir",      // FLIR Systems
	"flsmidth",  // FLSmidth
	"ford",      // Ford
	"forex",     // Dotforex Registry Ltd
	"fox",       // Fox Broadcasting Company
	"fresenius", // Fresenius Immobilien-Verwaltungs-GmbH
	"frogans",   // OP3FT
	"frontier",  // Frontier Communications
	"ftr",
	"fujitsu",   // Fujitsu
	"fujixerox", // Fuji Xerox
	"gallo",     // E & J Gallo Winery
	"gallup",    // Gallup
	"gap",       // Gap
	"gbiz",      // Google
	"gea",       // GEA Group
	"genting",   // Genting Group
	"george",
	"ggee",
	"gle",       // Google
	"globo",     // Grupo Globo
	"gmail",     // Google (Gmail)
	"gmo",       // GMO Internet
	"gmx",       // 1&1 Mail & Media (GMX, Global Message Exchange)
	"godaddy",   // Go Daddy
	"goldpoint", // Yodobashi Camera
	"goo",
	"goodyear", // Goodyear Tire and Rubber Company
	"goog",     // Google
	"google",   // Google
	"grainger", // Grainger Registry Services, LLC
	"guardian", // Guardian Life Insurance Company of America
	"gucci",    // Gucci
	"guge",
	"hbo",      // HBO
	"hdfc",     // Housing Development Finance Corporation
	"hdfcbank", // HDFC Bank
	"hermes",   // Hermès
	"hgtv",
	"hisamitsu", // Hisamitsu Pharmaceutical
	"hitachi",   // Hitachi
	"hkt",       // Hong Kong Telecom
	"homedepot",
	"honda",     // Honda
	"honeywell", // Honeywell
	"hotmail",   // Microsoft (Hotmail)
	"hsbc",      // HSBC
	"htc",
	"hughes",  // Hughes Network Systems
	"hyatt",   // Hyatt
	"hyundai", // Hyundai Motor Company
	"ibm",     // IBM
	"icbc",
	"ieee", // Institute of Electrical & Electronics Engineers
	"ifm",  // ifm electronic gmbh
	"iinet",
	"ikano",    // Ikano
	"imdb",     // Amazon (IMDb)
	"infiniti", // Nissan (Infiniti)
	"intel",    // Intel
	"intuit",   // Intuit
	"ipiranga", // Ipiranga
	"iselect",  // iSelect
	"itau",     // Itaú Unibanco
	"itv",      // ITV
	"iveco",    // CNH Industrial (Iveco)
	"iwc",
	"jaguar", // Jaguar Land Rover
	"jcb",    // JCB
	"jcp",    // JCP Media
	"jeep",   // Chrysler (Jeep)
	"jio",
	"jlc",
	"jll",
	"jmp",
	"jnj",
	"jpmorgan", // JPMorgan Chase
	"jprs",
	"juniper",         // Juniper Networks
	"kddi",            // KDDI
	"kerryhotels",     // Kerry Trading Co. Limited
	"kerrylogistics",  // Kerry Trading Co. Limited
	"kerryproperties", // Kerry Trading Co. Limited
	"kfh",             // Kuwait Finance House
	"kia",             // Kia Motors Corporation
	"kinder",          // Ferrero (Kinder Surprise)
	"kindle",          // Amazon (Amazon Kindle)
	"komatsu",         // Komatsu
	"kpmg",            // KPMG
	"kred",            // KredTLD
	"kuokgroup",       // Kerry Trading Co. Limited
	"lacaixa",         // Caixa d’Estalvis i Pensions de Barcelona
	"ladbrokes",       // Ladbrokes
	"lamborghini",     // Lamborghini
	"lamer",
	"lancaster", // Lancaster
	"lancia",    // Fiat Chrysler Automobiles (Lancia)
	"lancome",   // L'Oréal
	"landrover", // Jaguar Land Rover
	"lanxess",   // Lanxess
	"lasalle",   // JLL
	"latrobe",   // La Trobe University
	"leclerc",   // E.Leclerc
	"lefrak",
	"lego",    // Lego Group
	"lexus",   // Toyota (Lexus)
	"liaison", // Liaison Technologies
	"lidl",    // Lidl
	"lilly",   // Eli Lilly & Company
	"lincoln", // Ford (Lincoln)
	"linde",   // Linde
	"lipsy",   // Lipsy Ltd
	"lixil",   // Lixil Group
	"locus",   // Locus Analytics
	"loft",
	"lotte",        // Lotte Holdings
	"lpl",          // LPL Financial
	"lplfinancial", // LPL Financial
	"lundbeck",     // Lundbeck
	"lupin",        // Lupin Limited
	"macys",        // Macy's
	"maif",         // Mutuelle Assurance Instituteur France (MAIF)
	"man",          // MAN
	"mango",        // Mango
	"marriott",     // Marriott International
	"marshalls",
	"maserati", // Fiat Chrysler Automobiles (Maserati)
	"mattel",   // Mattel
	"mckinsey", // McKinsey & Company
	"meo",
	"metlife",    // MetLife
	"microsoft",  // Microsoft
	"mini",       // BMW (Mini)
	"mit",        // Massachusetts Institute of Technology
	"mitsubishi", // Mitsubishi Corporation
	"mlb",        // MLB Advanced Media
	"mma",        // MMA IARD
	"mobily",
	"monash", // Monash University
	"montblanc",
	"mopar",
	"moto",   // Motorola
	"msd",    // MSD Registry Holdings, Inc.
	"mtn",    // MTN
	"mtr",    // MTR Corporation
	"mutual", // Northwestern Mutual
	"mutuelle",
	"nab",
	"nadex",      // Nadex
	"nationwide", // Nationwide Mutual Insurance Company
	"natura",     // Natura & Co
	"nba",        // National Basketball Association
	"nec",        // NEC
	"netbank",
	"netflix",    // Netflix
	"neustar",    // Neustar
	"newholland", // CNH Industrial (New Holland Agriculture, New Holland Construction)
	"next",
	"nextdirect",
	"nexus",              // Google (Google Nexus)
	"nfl",                // National Football League
	"nhk",                // NHK
	"nico",               // Dwango (Niconico)
	"nike",               // Nike
	"nikon",              // Nikon
	"nissan",             // Nissan
	"nissay",             // Nippon Life
	"nokia",              // Nokia
	"northwesternmutual", // Northwestern Mutual
	"norton",             // NortonLifeLock
	"nowtv",
	"nra", // National Rifle Association of America
	"ntt", // Nippon Telegraph & Telephone
	"obi", // OBI Group Holding SE & Co. KGaA
	"off",
	"office", // Microsoft (Microsoft Office)
	"olayan",
	"olayangroup",
	"oldnavy",
	"ollo",
	"omega",  // Swatch Group (Omega)
	"oracle", // Oracle
	"orange", // Orange
	"orientexpress",
	"otsuka", // Otsuka Pharmaceutical
	"ovh",    // OVH
	"pamperedchef",
	"panasonic",   // Panasonic
	"pccw",        // PCCW
	"pfizer",      // Pfizer
	"philips",     // Philips
	"piaget",      // Piaget
	"pictet",      // Pictet
	"ping",        // Ping
	"pioneer",     // Pioneer Corporation
	"playstation", // Sony (PlayStation)
	"pnc",
	"pohl", // Deutsche Vermögensberatung
	"pramerica",
	"praxi",       // Praxi
	"prod",        // Google (products)
	"progressive", // Progressive Corporation
	"pru",         // Prudential Financial
	"prudential",  // Prudential Financial
	"pwc",         // PwC
	"quest",       // Quest Software
	"qvc",         // QVC
	"redstone",    // Redstone Haute Couture Co
	"redumbrella",
	"reliance", // Reliance Industries
	"ren",
	"rexroth", // Robert Bosch GmbH
	"richardli",
	"ricoh", // Ricoh
	"ril",
	"rmit",   // Royal Melbourne Institute of Technology
	"rocher", // Ferrero (Ferrero Rocher)
	"rogers", // Rogers Communications
	"rwe",    // RWE
	"sakura", // SAKURA Internet Inc.
	"samsclub",
	"samsung",         // Samsung SDS
	"sandvik",         // Sandvik
	"sandvikcoromant", // Sandvik Coromant
	"sanofi",          // Sanofi
	"sanvik",
	"sap", // SAP
	"sapo",
	"sas",
	"saxo",       // Saxo Bank
	"sbi",        // State Bank of India
	"sbs",        // Special Broadcasting Service[8]
	"sca",        // Svenska Cellulosa
	"scb",        // Siam Commercial Bank
	"schaeffler", // Schaeffler Technologies
	"schmidt",    // Société Alsacienne de Meubles (Cuisines Schmidt)
	"schwarz",    // Schwarz Gruppe
	"scjohnson",  // SC Johnson & Son
	"scor",       // SCOR
	"seat",       // SEAT
	"seek",       // Seek
	"sener",      // Sener Ingeniería y Sistemas, S.A.
	"ses",        // SES
	"seven",      // Seven West Media
	"sew",        // SEW Eurodrive
	"sfr",        // SFR
	"shangrila",  // Shangri-La Hotels & Resorts
	"sharp",      // Sharp Corporation
	"shaw",       // Shaw Cablesystems G.P.
	"shell",      // Shell
	"shouji",
	"shriram",    // Shriram Capital
	"sina",       // Sina Corp
	"sky",        // Sky Group
	"skype",      // Microsoft (Skype)
	"smart",      // Smart Communications
	"sncf",       // SNCF
	"softbank",   // SoftBank
	"sohu",       // Sohu
	"sony",       // Sony
	"spiegel",    // Der Spiegel
	"stada",      // Stada Arzneimittel
	"staples",    // Staples
	"star",       // Star India
	"starhub",    // StarHub
	"statebank",  // State Bank of India
	"statefarm",  // State Farm
	"statoil",    // Statoil
	"stc",        // Saudi Telecom Company
	"stcgroup",   // Saudi Telecom Company
	"suzuki",     // Suzuki
	"swatch",     // Swatch Group
	"swiftcover", // Swiftcover
	"symantec",   // NortonLifeLock
	"tab",
	"taobao",     // Alibaba Group (Taobao)
	"target",     // Target
	"tatamotors", // Tata Motors
	"tci",
	"tdk",        // TDK
	"telecity",   // TelecityGroup
	"telefonica", // Telefónica
	"temasek",    // Temasek Holdings
	"teva",       // Teva Pharmaceuticals
	"thd",
	"tiaa",
	"tiffany", // Tiffany & Co
	"tjmaxx",
	"tjx", // TJX Companies
	"tkmaxx",
	"tmall",
	"toray",         // Toray Industries
	"toshiba",       // Toshiba
	"total",         // Total
	"toyota",        // Toyota
	"travelchannel", // Lifestyle Domain Holdings, Inc. (Travel Channel)
	"travelers",     // Travelers TLD, LLC (The Travelers Companies)
	"trv",
	"tui", // TUI
	"tushu",
	"tvs", // T V Sundram Iyengar & Sons Private Limited
	"ubank",
	"ubs", // UBS
	"uconnect",
	"unicom", // China Unicom
	"uol",    // Universo Online
	"ups",    // United Parcel Service
	"vana",
	"vanguard",   // Vanguard Group
	"verisign",   // VeriSign
	"vig",        // Vienna Insurance Group
	"viking",     // Viking River Cruises
	"virgin",     // Virgin Group
	"visa",       // Visa
	"vista",      // Vistaprint
	"vistaprint", // Vistaprint
	"viva",
	"volkswagen", // Volkswagen Group of America
	"volvo",      // Volvo
	"walmart",    // Walmart
	"walter",     // Sandvik
	"warman",
	"weatherchannel", // IBM (The Weather Company)
	"weber",          // Saint-Gobain
	"weir",           // Weir Group
	"williamhill",    // William Hill
	"windows",        // Microsoft (Microsoft Windows)
	"wme",            // Endeavor
	"wolterskluwer",  // Wolters Kluwer
	"woodside",       // Woodside Petroleum
	"wtc",            // World Trade Centers Association
	"xbox",           // Microsoft Xbox
	"xerox",          // Fuji Xerox
	"xfinity",        // Comcast
	"xn--5su34j936bgsg",
	"xn--8y0a063a",
	"xn--b4w605ferd",
	"xn--cg4bki",
	"xn--eckvdtc9d",
	"xn--estv75g", // 工行
	"xn--fiq64b",
	"xn--flw351e",
	"xn--fzys8d69uvgm",
	"xn--jlq61u9w7b",
	"xn--kcrx77d1x4a",
	"xn--mgba3a3ejt",
	"xn--mgba7c0bbn0a",
	"xn--mgbb9fbpob",
	"xn--mgbi4ecexp",
	"xn--ngbe9e0a",
	"xn--qcka1pmc",
	"xn--w4r85el8fhu5dnra",
	"xn--w4rs40l",
	"xperia",
	"yahoo",     // Yahoo!
	"yamaxun",   // Amazon (Amazon China)
	"yandex",    // Yandex
	"yodobashi", // Yodobashi Camera
	"youtube",   // Google (YouTube)
	"zappos",    // Amazon (Zappos)
	"zara",      // Inditex (Zara)
	"zippo",     // Zippo
}

// proposedBrandTlds must be sorted (for binary search to work)
var proposedBrandTlds = []string{
	"acer",
	"afamilycompany",
	"africamagic",
	"alcon",
	"amp",
	"barclaycard", // Barclays
	"barclaysbank",
	"bbb",
	"bestbuy",
	"blockbuster",
	"canalplus",
	"case",
	"changiairport",
	"cimb",
	"deutschepost",
	"digikey",
	"dnb",
	"docomo",
	"dstv",
	"duns",
	"dwg",
	"emerson",
	"gecompany",
	"glade",
	"goodhands",
	"gotv",
	"indians",
	"jpmorganchase",
	"konami",
	"kyknet",
	"livestrong",
	"mcd",
	"mcdonalds",
	"merck",
	"merckmsd",
	"mih",
	"mnet",
	"mozaic",
	"mrmuscle",
	"multichoice",
	"mzansimagic",
	"naspers",
	"panerai",
	"payu",
	"piperlime",
	"pitney",
	"qtel",
	"raid",
	"ram",
	"rockwool",
	"schwarzgroup",
	"shopyourway",
	"srt",
	"supersport",
	"tata",
	"terra",
	"theguardian",
	"tradershotels",
	"transunion",
	"travelguard",
	"unicorn",
	"vanish",
	"webjet",
	"wilmar",
	"xn--4gq48lf9j",           // 一号店
	"xn--55qx5d8y0buji4b870u", // 通用电气公司
	"xn--cckwcxetd",           // アマゾン
	"xn--hxt035czzpffl",       // 盛贸饭店
	"xn--mgbaakc7dvf",
	"xn--mgbv6cfpo",
	"xn--pgb3ceoj",
}

// terminatedBrandTlds must be sorted (for binary search to work)
var terminatedBrandTlds = []string{
	"axis",
	"caseih", // termination in process
	"ceb",
	"delmonte",
	"doosan", // retrired
	"gcc",    // rejected
	"mtpc",
	"onyourside",
	"rightathome",
}

var withdrawnBrandTlds = []string{
	"allfinanzberater",  // 	WITHDRAWN
	"allfinanzberatung", // 	WITHDRAWN
	"and",               // 	WITHDRAWN
	"ansons",            // 	WITHDRAWN
	"anthem",            // 	WITHDRAWN
	"astrium",           // 	WITHDRAWN
	"avery",             // 	WITHDRAWN
	"beknown",           // 	WITHDRAWN
	"bloomingdales",     // 	WITHDRAWN
	"buick",             // 	WITHDRAWN
	"bway",
	"cadillac",      // 	WITHDRAWN
	"caremore",      // 	WITHDRAWN
	"chartis",       // 	WITHDRAWN
	"chatr",         // 	WITHDRAWN
	"chesapeake",    // 	WITHDRAWN
	"chevrolet",     // 	WITHDRAWN
	"chevy",         // 	WITHDRAWN
	"chk",           // 	WITHDRAWN
	"cialis",        // 	WITHDRAWN
	"fls",           // 	WITHDRAWN. Application submitted by Thomsen Trampedach.
	"garnier",       // 	WITHDRAWN
	"glean",         // 	WITHDRAWN
	"globalx",       // 	WITHDRAWN
	"gmc",           // 	WITHDRAWN
	"gree",          // 	WITHDRAWN; Community Application
	"guardianlife",  // 	WITHDRAWN
	"guardianmedia", // 	WITHDRAWN
	"heinz",         // 	WITHDRAWN
	"hilton",        // 	WITHDRAWN
	"infosys",       // 	WITHDRAWN
	"infy",          // 	WITHDRAWN
	"justforu",
	"kerastase",
	"kiehls",     // 	WITHDRAWN
	"kone",       // 	WITHDRAWN
	"ksb",        // 	WITHDRAWN
	"loreal",     // 	WITHDRAWN
	"matrix",     // 	WITHDRAWN
	"maybelline", // 	WITHDRAWN
	"mii",
	"mitek", // 	WITHDRAWN
	"mrporter",
	"netaporter",
	"northlandinsurance", // 	WITHDRAWN
	"olympus",            // 	WITHDRAWN
	"patagonia",          // 	GAC Early Warning: Australia. Objected to by Independent Objector. Deemed Geographic by GAC, WITHDRAWN.
	"polo",               // 	WITHDRAWN
	"redken",             // 	WITHDRAWN
	"safeway",
	"saphhire",     // 	WITHDRAWN
	"skolkovo",     // 	WITHDRAWN
	"skydrive",     // 	WITHDRAWN
	"svr",          // 	WITHDRAWN
	"swiss",        // 	WITHDRAWN
	"thehartford",  // 	WITHDRAWN
	"transformers", // 	WITHDRAWN
	"ultrabook",
	"vons",
	"xn--dkwm73cwpn",     // 欧莱雅
	"xn--hxt035cmppuel",  // 盛貿飯店
	"xn--j6w470d71issc",  // 香港電訊
	"xn--kcrx7bb75ajk3b", // 普利司通
}

// unknownTlds must be sorted (for binary search to work)
var unknownTlds = []string{
	"bq", // Caribbean Netherlands ( Bonaire,  Saba, and  Sint Eustatius)
	"bv",
	"eh", // Western Sahara
	"entertainment",
	"gb",
	"sj",
	"xn--cckwcxetd",
	"xn--gckr3f0f",
	"xn--jlq480n2rg",
	"xn--mgb2ddes", // اليمن (Yemen)
	"xn--mgbaakc7dvf",
	"xn--mgbtx2b", // عراق (Iraq)
	"xn--mix082f", // 澳门 (Macao)
	"xn--ngbrx",
	"xn--otu796d",
}

func binarySearch(sl []string, key string) bool {
	low := 0
	high := len(sl) - 1
	for low <= high {
		mid := (low + high) >> 1
		midVal := sl[mid]
		if cmp := strings.Compare(midVal, key); cmp == 0 {
			return true
		} else if cmp < 0 {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return false
}
