package valix

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestISO3166_1_CountryCodesMatrix(t *testing.T) {
	expectedCounts := map[assignmentType3166_1]int{
		ccAs: 242,
		ccRa: 7,
		ccUA: 43,
		ccER: 13,
		ccIR: 30,
		ccTR: 7,
		ccDl: 14,
		ccUn: 320,
	}
	actualCounts := map[assignmentType3166_1]int{
		ccAs: 0,
		ccRa: 0,
		ccUA: 0,
		ccER: 0,
		ccIR: 0,
		ccTR: 0,
		ccDl: 0,
		ccUn: 0,
	}
	totalCount := 0
	assignedCount := 0
	for k1, v1 := range iso3166_1_CountryCodesMatrix {
		for k2, v2 := range v1 {
			totalCount++
			actualCounts[v2] = actualCounts[v2] + 1
			if v2 == ccAs || v2 == ccRa {
				assignedCount++
				cc := string([]byte{k1, k2})
				_, exists := ISO3166_2_CountryCodes[cc]
				require.True(t, exists)
			}
		}
	}
	assert.Equal(t, 676, totalCount)
	assert.Equal(t, len(ISO3166_2_CountryCodes), assignedCount)
	for k, v := range expectedCounts {
		t.Run(fmt.Sprintf("CC3166[%d]=%d", k, v), func(t *testing.T) {
			actual := actualCounts[k]
			require.Equal(t, v, actual)
		})
	}
}

func TestISO3166_2(t *testing.T) {
	for k1, v1 := range iso3166_1_CountryCodesMatrix {
		for k2, v2 := range v1 {
			if v2 == ccAs || v2 == ccRa {
				cc := []byte{k1, k2}
				_, ok := ISO3166_2_CountryCodes[string(cc)]
				require.True(t, ok)
			}
		}
	}
}

func TestISO3166_2_CheckCounts(t *testing.T) {
	for k, cnt := range iso3166_2_checks {
		t.Run(fmt.Sprintf("Regions[%s]=%d", k, cnt), func(t *testing.T) {
			rs, ok := ISO3166_2_CountryCodes[k]
			require.True(t, ok)
			require.Equal(t, cnt, len(rs))
		})
	}
}

var iso3166_2_checks = map[string]int{
	"AD": 7, // Andorra
	// 7 parishes
	"AE": 7, // United Arab Emirates
	// 7 emirates
	"AF": 34, // Afghanistan
	// 34 provinces
	"AG": 6 + 2, // Antigua and Barbuda
	// 6 parishes, 2 dependencies
	"AI": 0,  // Anguilla
	"AL": 12, // Albania
	// 12 counties
	"AM": 1 + 10, // Armenia
	// 1 city, 10 regions
	"AO": 18, // Angola
	// 18 provinces
	"AQ": 0,      // Antarctica
	"AR": 1 + 23, // Argentina
	// 1 city, 23 provinces
	"AS": 0, // American Samoa
	"AT": 9, // Austria
	// 9 states
	"AU": 6 + 2, // Australia
	// 6 states, 2 territories
	"AW": 0,           // Aruba
	"AX": 0,           // Åland Islands
	"AZ": 1 + 11 + 66, // Azerbaijan
	// 1 autonomous republic, 11 municipalities, 66 rayons
	"BA": 2 + 1, // Bosnia and Herzegovina
	// 2 entities, 1 district with special status
	"BB": 11, // Barbados
	// 11 parishes
	"BD": 8 + 64, // Bangladesh
	// 8 divisions, 64 districts
	"BE": 3 + 10, // Belgium
	// 3 regions, 10 provinces
	"BF": 13 + 45, // Burkina Faso
	// 13 regions, 45 provinces
	"BG": 28, // Bulgaria
	// 28 regions
	"BH": 4, // Bahrain
	// 4 governorates
	"BI": 18, // Burundi
	// 18 provinces
	"BJ": 12, // Benin
	// 12 departments
	"BL": 0, // Saint Barthélemy
	"BM": 0, // Bermuda
	"BN": 4, // Brunei Darussalam
	// 4 districts
	"BO": 9, // Bolivia (Plurinational State of)
	// 9 departments
	"BQ": 3, // Bonaire, Sint Eustatius and Saba
	// 3 special municipalities
	"BR": 1 + 26, // Brazil
	// 1 federal district, 26 states
	"BS": 31 + 1, // Bahamas
	// 31 districts, 1 island
	"BT": 20, // Bhutan
	// 20 districts
	"BV": 0,          // Bouvet Island
	"BW": 10 + 4 + 2, // Botswana
	// 10 districts, 4 towns, 2 cities
	"BY": 6 + 1, // Belarus
	// 6 oblasts, 1 city
	"BZ": 6, // Belize
	// 6 districts
	"CA": 10 + 3, // Canada
	// 10 provinces, 3 territories
	"CC": 0,      // Cocos (Keeling) Islands
	"CD": 1 + 25, // Congo, Democratic Republic of the
	// 1 city, 25 provinces
	"CF": 1 + 14 + 2, // Central African Republic
	// 1 commune, 14 prefectures, 2 economic prefectures
	"CG": 12, // Congo
	// 12 departments
	"CH": 26, // Switzerland
	// 26 cantons
	"CI": 12 + 2, // Côte d'Ivoire
	// 12 districts, 2 autonomous districts
	"CK": 0,  // Cook Islands
	"CL": 16, // Chile
	// 16 regions
	"CM": 10, // Cameroon
	// 10 regions
	"CN": 4 + 23 + 5 + 2, // China
	// 4 municipalities, 23 provinces, 5 autonomous regions, 2 special administrative regions
	"CO": 1 + 32, // Colombia
	// 1 capital district, 32 departments
	"CR": 7, // Costa Rica
	// 7 provinces
	"CU": 15 + 1, // Cuba
	// 15 provinces, 1 special municipality
	"CV": 2 + 22, // Cabo Verde
	// 2 geographical regions, 22 municipalities
	"CW": 0, // Curaçao
	"CX": 0, // Christmas Island
	"CY": 6, // Cyprus
	// 6 districts
	"CZ": 13 + 1 + 76, // Czechia
	// 13 regions and 1 capital city, 76 districts
	"DE": 16, // Germany
	// 16 states
	"DJ": 5 + 1, // Djibouti
	// 5 regions, 1 city
	"DK": 5, // Denmark
	// 5 regions
	"DM": 10, // Dominica
	// 10 parishes
	"DO": 10 + 1 + 31, // Dominican Republic
	// 10 regions, 1 district, 31 provinces
	"DZ": 48, // Algeria
	// 48 provinces
	"EC": 24, // Ecuador
	// 24 provinces
	"EE": 15 + 64 + 15, // Estonia
	// 15 counties, 64 rural municipalities, 15 urban municipalities
	"EG": 27, // Egypt
	// 27 governorates
	"EH": 0, // Western Sahara
	"ER": 6, // Eritrea
	// 6 regions
	"ES": 17 + 2 + 50, // Spain
	// 17 autonomous communities, 2 autonomous cities in North Africa, 50 provinces
	"ET": 2 + 10, // Ethiopia
	// 2 administrations, 10 regional states
	"FI": 19, // Finland
	// 19 regions
	"FJ": 4 + 1 + 14, // Fiji
	// 4 divisions, 1 dependency, 14 provinces
	"FK": 0, // Falkland Islands (Malvinas)
	"FM": 4, // Micronesia (Federated States of)
	// 4 states
	"FO": 0,   // Faroe Islands
	"FR": 124, // France
	"GA": 9,   // Gabon
	// 9 provinces
	"GB": 223,   // United Kingdom of Great Britain and Northern Ireland
	"GD": 6 + 1, // Grenada
	// 6 parishes, 1 dependency
	"GE": 2 + 1 + 9, // Georgia
	// 2 autonomous republics, 1 city, 9 regions
	"GF": 0,  // French Guiana
	"GG": 0,  // Guernsey
	"GH": 16, // Ghana
	// 16 regions
	"GI": 0, // Gibraltar
	"GL": 5, // Greenland
	// 5 municipalities
	"GM": 1 + 5, // Gambia
	// 1 city, 5 divisions
	"GN": 7 + 1 + 33, // Guinea
	// 7 administrative regions, 1 governorate, 33 prefectures
	"GP": 0,     // Guadeloupe
	"GQ": 2 + 8, // Equatorial Guinea
	// 2 regions, 8 provinces
	"GR": 13 + 1, // Greece
	// 13 administrative regions, 1 self-governed part
	"GS": 0,  // South Georgia and the South Sandwich Islands
	"GT": 22, // Guatemala
	// 22 departments
	"GU": 0,         // Guam
	"GW": 3 + 1 + 8, // Guinea-Bissau
	// 3 provinces, 1 autonomous sector, 8 regions
	"GY": 10, // Guyana
	// 10 regions
	"HK": 0,  // Hong Kong
	"HM": 0,  // Heard Island and McDonald Islands
	"HN": 18, // Honduras
	// 18 departments
	"HR": 1 + 20, // Croatia
	// 1 city, 20 counties
	"HT": 10, // Haiti
	// 10 departments
	"HU": 1 + 19 + 23, // Hungary
	// 1 capital city, 19 counties, 23 cities of county right
	"ID": 7 + 32 + 1 + 1, // Indonesia
	// 7 geographical units, 32 provinces, 1 capital district, 1 special region
	"IE": 4 + 26, // Ireland
	// 4 provinces, 26 counties
	"IL": 6, // Israel
	// 6 districts
	"IM": 0,      // Isle of Man
	"IN": 28 + 8, // India
	// 28 states, 8 union territories
	"IO": 0,      // British Indian Ocean Territory
	"IQ": 18 + 1, // Iraq
	// 18 governorates, 1 region
	"IR": 31, // Iran (Islamic Republic of)
	// 31 provinces
	"IS": 8 + 69, // Iceland
	// 8 regions, 69 municipalities
	"IT": 15 + 80 + 5 + 2 + 6 + 14 + 4, // Italy
	// 15 regions
	// 80 provinces, 5 autonomous regions, 2 autonomous provinces
	// 6 free municipal consortiums, 14 metropolitan cities, 4 decentralized regional entities
	"JE": 0,  // Jersey
	"JM": 14, // Jamaica
	// 14 parishes
	"JO": 12, // Jordan
	// 12 governorates
	"JP": 47, // Japan
	// 47 prefectures
	"KE": 47, // Kenya
	// 47 counties
	"KG": 2 + 7, // Kyrgyzstan
	// 2 cities, 7 regions
	"KH": 1 + 24, // Cambodia
	// 1 autonomous municipality, 24 provinces
	"KI": 3, // Kiribati
	// 3 groups of islands
	"KM": 3, // Comoros
	// 3 islands
	"KN": 2 + 14, // Saint Kitts and Nevis
	// 2 states, 14 parishes
	"KP": 1 + 1 + 1 + 9, // Korea (Democratic People's Republic of)
	// 1 capital city, 1 metropolitan city, 1 special city, 9 provinces
	"KR": 6 + 1 + 1 + 8 + 1, // Korea, Republic of
	// 6 metropolitan cities, 1 special city, 1 special self-governing city, 8 provinces, 1 special self-governing province
	"KW": 6, // Kuwait
	// 6 governorates
	"KY": 0,      // Cayman Islands
	"KZ": 3 + 14, // Kazakhstan
	// 3 cities, 14 regions
	"LA": 1 + 17, // Lao People's Democratic Republic
	// 1 prefecture, 17 provinces
	"LB": 8, // Lebanon
	// 8 governorates
	"LC": 10, // Saint Lucia
	// 10 districts
	"LI": 11, // Liechtenstein
	// 11 communes
	"LK": 9 + 25, // Sri Lanka
	// 9 provinces, 25 districts
	"LR": 15, // Liberia
	// 15 counties
	"LS": 10, // Lesotho
	// 10 districts
	"LT": 10 + 9 + 7 + 44, // Lithuania
	// 10 counties, 9 municipalities, 7 city municipalities, 44 district municipalities
	"LU": 12, // Luxembourg
	// 12 cantons
	"LV": 36 + 7, // Latvia
	// 36 municipalities, 7 state cities
	"LY": 22, // Libya
	// 22 popularates
	"MA": 12 + 62 + 13, // Morocco
	// 12 regions, 62 provinces, 13 prefectures
	"MC": 17, // Monaco
	// 17 quarters
	"MD": 1 + 3 + 32 + 1, // Moldova, Republic of
	// 1 autonomous territorial unit, 3 cities, 32 districts, 1 territorial unit
	"ME": 24, // Montenegro
	// 24 municipalities
	"MF": 0, // Saint Martin (French part)
	"MG": 6, // Madagascar
	// 6 provinces
	"MH": 2 + 24, // Marshall Islands
	// 2 chains of islands, 24 municipalities
	"MK": 80, // North Macedonia
	// 80 municipalities
	"ML": 1 + 10, // Mali
	// 1 district, 10 regions
	"MM": 7 + 7 + 1, // Myanmar
	// 7 regions, 7 states, 1 union territory
	"MN": 1 + 21, // Mongolia
	// 1 capital city, 21 provinces
	"MO": 0,  // Macao
	"MP": 0,  // Northern Mariana Islands
	"MQ": 0,  // Martinique
	"MR": 15, // Mauritania
	// 15 regions
	"MS": 0,  // Montserrat
	"MT": 68, // Malta
	// 68 local councils
	"MU": 3 + 9, // Mauritius
	// 3 dependencies, 9 districts
	"MV": 19 + 2, // Maldives
	// 19 administrative atolls, 2 cities
	"MW": 3 + 28, // Malawi
	// 3 regions, 28 districts
	"MX": 31 + 1, // Mexico
	// 31 states, 1 federal district
	"MY": 3 + 13, // Malaysia
	// 3 federal territories, 13 states
	"MZ": 1 + 10, // Mozambique
	// 1 city, 10 provinces
	"NA": 14, // Namibia
	// 14 regions
	"NC": 0,     // New Caledonia
	"NE": 1 + 7, // Niger
	// 1 urban community, 7 departments
	"NF": 0,      // Norfolk Island
	"NG": 1 + 36, // Nigeria
	// 1 capital territory, 36 states
	"NI": 15 + 2, // Nicaragua
	// 15 departments, 2 autonomous regions
	"NL": 12 + 3 + 3, // Netherlands[note 1]
	// 12 provinces, 3 countries, 3 special municipalities
	"NO": 11 + 2, // Norway
	// 11 counties, 2 arctic regions
	"NP": 5 + 7 + 14, // Nepal
	// 5 development regions, 7 provinces, 14 zones
	"NR": 14, // Nauru
	// 14 districts
	"NU": 0,      // Niue
	"NZ": 16 + 1, // New Zealand
	// 16 regions, 1 special island authority
	"OM": 11, // Oman
	// 11 governorates
	"PA": 10 + 4, // Panama
	// 10 provinces, 4 indigenous regions
	"PE": 25 + 1, // Peru
	// 25 regions, 1 municipality
	"PF": 0,          // French Polynesia
	"PG": 1 + 20 + 1, // Papua New Guinea
	// 1 district, 20 provinces, 1 autonomous region
	"PH": 17 + 81, // Philippines
	// 17 regions, 81 provinces
	"PK": 4 + 2 + 1, // Pakistan
	// 4 provinces, 2 autonomous territories, 1 federal territory
	"PL": 16, // Poland
	// 16 voivodships
	"PM": 0,  // Saint Pierre and Miquelon
	"PN": 0,  // Pitcairn
	"PR": 0,  // Puerto Rico
	"PS": 16, // Palestine, State of
	// 16 governorates
	"PT": 18 + 2, // Portugal
	// 18 districts, 2 autonomous regions
	"PW": 16, // Palau
	// 16 states
	"PY": 1 + 17, // Paraguay
	// 1 capital, 17 departments
	"QA": 8, // Qatar
	// 8 municipalities
	"RE": 0,      // Réunion
	"RO": 41 + 1, // Romania
	// 41 departments, 1 municipality
	"RS": 2 + 1 + 29, // Serbia
	// 2 autonomous provinces, 1 city, 29 districts
	"RU": 21 + 9 + 46 + 2 + 1 + 4, // Russian Federation
	// 21 republics, 9 administrative territories
	// 46 administrative regions, 2 autonomous cities, 1 autonomous region, 4 autonomous districts
	"RW": 1 + 4, // Rwanda
	// 1 town council, 4 provinces
	"SA": 13, // Saudi Arabia
	// 13 regions
	"SB": 1 + 9, // Solomon Islands
	// 1 capital territory, 9 provinces
	"SC": 27, // Seychelles
	// 27 districts
	"SD": 18, // Sudan
	// 18 states
	"SE": 21, // Sweden
	// 21 counties
	"SG": 5, // Singapore
	// 5 districts
	"SH": 3, // Saint Helena, Ascension and Tristan da Cunha
	// 3 geographical entities
	"SI": 212, // Slovenia
	// 212 municipalities
	"SJ": 0, // Svalbard and Jan Mayen
	"SK": 8, // Slovakia
	// 8 regions
	"SL": 1 + 4, // Sierra Leone
	// 1 area, 4 provinces
	"SM": 9, // San Marino
	// 9 municipalities
	"SN": 14, // Senegal
	// 14 regions
	"SO": 18, // Somalia
	// 18 regions
	"SR": 10, // Suriname
	// 10 districts
	"SS": 10, // South Sudan
	// 10 states
	"ST": 1 + 6, // Sao Tome and Principe
	// 1 autonomous region, 6 districts
	"SV": 14, // El Salvador
	// 14 departments
	"SX": 0,  // Sint Maarten (Dutch part)
	"SY": 14, // Syrian Arab Republic
	// 14 provinces
	"SZ": 4, // Eswatini
	// 4 regions
	"TC": 0,  // Turks and Caicos Islands
	"TD": 23, // Chad
	// 23 provinces
	"TF": 0, // French Southern Territories
	"TG": 5, // Togo
	// 5 regions
	"TH": 1 + 1 + 76, // Thailand
	// 1 metropolitan administration, 1 special administrative city, 76 provinces
	"TJ": 1 + 2 + 1 + 1, // Tajikistan
	// 1 autonomous region, 2 regions, 1 capital territory, 1 district under republic administration
	"TK": 0,      // Tokelau
	"TL": 12 + 1, // Timor-Leste
	// 12 municipalities, 1 special administrative region
	"TM": 5 + 1, // Turkmenistan
	// 5 regions, 1 city
	"TN": 24, // Tunisia
	// 24 governorates
	"TO": 5, // Tonga
	// 5 divisions
	"TR": 81, // Türkiye
	// 81 provinces
	"TT": 9 + 3 + 2 + 1, // Trinidad and Tobago
	// 9 regions, 3 boroughs, 2 cities, 1 ward
	"TV": 1 + 7, // Tuvalu
	// 1 town council, 7 island councils
	"TW": 13 + 3 + 6, // Taiwan, Province of China[note 2]
	// 13 counties, 3 cities, 6 special municipalities
	"TZ": 31, // Tanzania, United Republic of
	// 31 regions
	"UA": 24 + 1 + 2, // Ukraine
	// 24 regions, 1 republic, 2 cities
	"UG": 4 + 134 + 1, // Uganda
	// 4 geographical regions, 134 districts, 1 city
	"UM": 9, // United States Minor Outlying Islands
	// 9 islands, groups of islands
	"US": 50 + 1 + 6, // United States of America
	// 50 states, 1 district, 6 outlying areas
	"UY": 19, // Uruguay
	// 19 departments
	"UZ": 1 + 12 + 1, // Uzbekistan
	// 1 city, 12 regions, 1 republic
	"VA": 0, // Holy City
	"VC": 6, // Saint Vincent and the Grenadines
	// 6 parishes
	"VE": 1 + 1 + 23, // Venezuela (Bolivarian Republic of)
	// 1 federal dependency, 1 federal district, 23 states
	"VG": 0,      // Virgin Islands (British)
	"VI": 0,      // Virgin Islands (U.S.)
	"VN": 58 + 5, // Viet Nam
	// 58 provinces, 5 municipalities
	"VU": 6, // Vanuatu
	// 6 provinces
	"WF": 3, // Wallis and Futuna
	// 3 administrative precincts
	"WS": 11, // Samoa
	// 11 districts
	"YE": 1 + 21, // Yemen
	// 1 municipality, 21 governorates
	"YT": 0, // Mayotte
	"ZA": 9, // South Africa
	// 9 provinces
	"ZM": 10, // Zambia
	// 10 provinces
	"ZW": 10, // Zimbabwe
	// 10 provinces
}
