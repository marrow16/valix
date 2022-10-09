package valix

const (
	ISO4217TestCurrencyCode        = "XTS"
	ISO4217TestCurrencyCodeNumeric = "963"
	ISO4217NoCurrencyCode          = "XXX"
	ISO4217NoCurrencyCodeNumeric   = "999"
)

var iSO4217CurrencyCodes = map[string]bool{
	"AED": true, // United Arab Emirates dirham
	"AFN": true, // Afghan afghani
	"ALL": true, // Albanian lek
	"AMD": true, // Armenian dram
	"ANG": true, // Netherlands Antillean guilder
	"AOA": true, // Angolan kwanza
	"ARS": true, // Argentine peso
	"AUD": true, // Australian dollar
	"AWG": true, // Aruban florin
	"AZN": true, // Azerbaijani manat
	"BAM": true, // Bosnia and Herzegovina convertible mark
	"BBD": true, // Barbados dollar
	"BDT": true, // Bangladeshi taka
	"BGN": true, // Bulgarian lev
	"BHD": true, // Bahraini dinar
	"BIF": true, // Burundian franc
	"BMD": true, // Bermudian dollar
	"BND": true, // Brunei dollar
	"BOB": true, // Boliviano
	"BOV": true, // Bolivian Mvdol (funds code)
	"BRL": true, // Brazilian real
	"BSD": true, // Bahamian dollar
	"BTN": true, // Bhutanese ngultrum
	"BWP": true, // Botswana pula
	"BYN": true, // Belarusian ruble
	"BZD": true, // Belize dollar
	"CAD": true, // Canadian dollar
	"CDF": true, // Congolese franc
	"CHE": true, // WIR euro (complementary currency)
	"CHF": true, // Swiss franc
	"CHW": true, // WIR franc (complementary currency)
	"CLF": true, // Unidad de Fomento (funds code)
	"CLP": true, // Chilean peso
	"CNY": true, // Renminbi
	"COP": true, // Colombian peso
	"COU": true, // Unidad de Valor Real (UVR) (funds code)
	"CRC": true, // Costa Rican colon
	"CUC": true, // Cuban convertible peso
	"CUP": true, // Cuban peso
	"CVE": true, // Cape Verdean escudo
	"CZK": true, // Czech koruna
	"DJF": true, // Djiboutian franc
	"DKK": true, // Danish krone
	"DOP": true, // Dominican peso
	"DZD": true, // Algerian dinar
	"EGP": true, // Egyptian pound
	"ERN": true, // Eritrean nakfa
	"ETB": true, // Ethiopian birr
	"EUR": true, // Euro
	"FJD": true, // Fiji dollar
	"FKP": true, // Falkland Islands pound
	"GBP": true, // Pound sterling
	"GEL": true, // Georgian lari
	"GHS": true, // Ghanaian cedi
	"GIP": true, // Gibraltar pound
	"GMD": true, // Gambian dalasi
	"GNF": true, // Guinean franc
	"GTQ": true, // Guatemalan quetzal
	"GYD": true, // Guyanese dollar
	"HKD": true, // Hong Kong dollar
	"HNL": true, // Honduran lempira
	"HRK": true, // Croatian kuna
	"HTG": true, // Haitian gourde
	"HUF": true, // Hungarian forint
	"IDR": true, // Indonesian rupiah
	"ILS": true, // Israeli new shekel
	"INR": true, // Indian rupee
	"IQD": true, // Iraqi dinar
	"IRR": true, // Iranian rial
	"ISK": true, // Icelandic króna (plural: krónur)
	"JMD": true, // Jamaican dollar
	"JOD": true, // Jordanian dinar
	"JPY": true, // Japanese yen
	"KES": true, // Kenyan shilling
	"KGS": true, // Kyrgyzstani som
	"KHR": true, // Cambodian riel
	"KMF": true, // Comoro franc
	"KPW": true, // North Korean won
	"KRW": true, // South Korean won
	"KWD": true, // Kuwaiti dinar
	"KYD": true, // Cayman Islands dollar
	"KZT": true, // Kazakhstani tenge
	"LAK": true, // Lao kip
	"LBP": true, // Lebanese pound
	"LKR": true, // Sri Lankan rupee
	"LRD": true, // Liberian dollar
	"LSL": true, // Lesotho loti
	"LYD": true, // Libyan dinar
	"MAD": true, // Moroccan dirham
	"MDL": true, // Moldovan leu
	"MGA": true, // Malagasy ariary
	"MKD": true, // Macedonian denar
	"MMK": true, // Myanmar kyat
	"MNT": true, // Mongolian tögrög
	"MOP": true, // Macanese pataca
	"MRU": true, // Mauritanian ouguiya
	"MUR": true, // Mauritian rupee
	"MVR": true, // Maldivian rufiyaa
	"MWK": true, // Malawian kwacha
	"MXN": true, // Mexican peso
	"MXV": true, // Mexican Unidad de Inversion (UDI) (funds code)
	"MYR": true, // Malaysian ringgit
	"MZN": true, // Mozambican metical
	"NAD": true, // Namibian dollar
	"NGN": true, // Nigerian naira
	"NIO": true, // Nicaraguan córdoba
	"NOK": true, // Norwegian krone
	"NPR": true, // Nepalese rupee
	"NZD": true, // New Zealand dollar
	"OMR": true, // Omani rial
	"PAB": true, // Panamanian balboa
	"PEN": true, // Peruvian sol
	"PGK": true, // Papua New Guinean kina
	"PHP": true, // Philippine peso
	"PKR": true, // Pakistani rupee
	"PLN": true, // Polish złoty
	"PYG": true, // Paraguayan guaraní
	"QAR": true, // Qatari riyal
	"RON": true, // Romanian leu
	"RSD": true, // Serbian dinar
	"RUB": true, // Russian ruble
	"RWF": true, // Rwandan franc
	"SAR": true, // Saudi riyal
	"SBD": true, // Solomon Islands dollar
	"SCR": true, // Seychelles rupee
	"SDG": true, // Sudanese pound
	"SEK": true, // Swedish krona (plural: kronor)
	"SGD": true, // Singapore dollar
	"SHP": true, // Saint Helena pound
	"SLE": true, // Sierra Leonean leone
	"SLL": true, // Sierra Leonean leone
	"SOS": true, // Somali shilling
	"SRD": true, // Surinamese dollar
	"SSP": true, // South Sudanese pound
	"STN": true, // São Tomé and Príncipe dobra
	"SVC": true, // Salvadoran colón
	"SYP": true, // Syrian pound
	"SZL": true, // Swazi lilangeni
	"THB": true, // Thai baht
	"TJS": true, // Tajikistani somoni
	"TMT": true, // Turkmenistan manat
	"TND": true, // Tunisian dinar
	"TOP": true, // Tongan paʻanga
	"TRY": true, // Turkish lira
	"TTD": true, // Trinidad and Tobago dollar
	"TWD": true, // New Taiwan dollar
	"TZS": true, // Tanzanian shilling
	"UAH": true, // Ukrainian hryvnia
	"UGX": true, // Ugandan shilling
	"USD": true, // United States dollar
	"USN": true, // United States dollar (next day) (funds code)
	"UYI": true, // Uruguay Peso en Unidades Indexadas (URUIURUI) (funds code)
	"UYU": true, // Uruguayan peso
	"UYW": true, // Unidad previsional
	"UZS": true, // Uzbekistan som
	"VED": true, // Venezuelan bolívar digital
	"VES": true, // Venezuelan bolívar soberano
	"VND": true, // Vietnamese đồng
	"VUV": true, // Vanuatu vatu
	"WST": true, // Samoan tala
	"XAF": true, // CFA franc BEAC
	"XAG": true, // Silver (one troy ounce)
	"XAU": true, // Gold (one troy ounce)
	"XBA": true, // European Composite Unit (EURCO) (bond market unit)
	"XBB": true, // European Monetary Unit (E.M.U.-6) (bond market unit)
	"XBC": true, // European Unit of Account 9 (E.U.A.-9) (bond market unit)
	"XBD": true, // European Unit of Account 17 (E.U.A.-17) (bond market unit)
	"XCD": true, // East Caribbean dollar
	"XDR": true, // Special drawing rights
	"XOF": true, // CFA franc BCEAO
	"XPD": true, // Palladium (one troy ounce)
	"XPF": true, // CFP franc (franc Pacifique)	French territories of the Pacific Ocean
	"XPT": true, // Platinum (one troy ounce)
	"XSU": true, // SUCRE Unified System for Regional Compensation (SUCRE)
	"XUA": true, // ADB Unit of Account	African Development Bank
	"YER": true, // Yemeni rial
	"ZAR": true, // South African rand
	"ZMW": true, // Zambian kwacha
	"ZWL": true, // Zimbabwean dollar
	//"XTS": true, // Code reserved for testing
	//"XXX": true, // No currency
}

var iSO4217CurrencyCodesNumeric = map[string]bool{
	"008": true, // ALL Albanian lek
	"012": true, // DZD Algerian dinar
	"032": true, // ARS Argentine peso
	"036": true, // AUD Australian dollar
	"044": true, // BSD Bahamian dollar
	"048": true, // BHD Bahraini dinar
	"050": true, // BDT Bangladeshi taka
	"051": true, // AMD Armenian dram
	"052": true, // BBD Barbados dollar
	"060": true, // BMD Bermudian dollar
	"064": true, // BTN Bhutanese ngultrum
	"068": true, // BOB Boliviano
	"072": true, // BWP Botswana pula
	"084": true, // BZD Belize dollar
	"090": true, // SBD Solomon Islands dollar
	"096": true, // BND Brunei dollar
	"104": true, // MMK Myanmar kyat
	"108": true, // BIF Burundian franc
	"116": true, // KHR Cambodian riel
	"124": true, // CAD Canadian dollar
	"132": true, // CVE Cape Verdean escudo
	"136": true, // KYD Cayman Islands dollar
	"144": true, // LKR Sri Lankan rupee
	"152": true, // CLP Chilean peso
	"156": true, // CNY Renminbi
	"170": true, // COP Colombian peso
	"174": true, // KMF Comoro franc
	"188": true, // CRC Costa Rican colon
	"191": true, // HRK Croatian kuna
	"192": true, // CUP Cuban peso
	"203": true, // CZK Czech koruna
	"208": true, // DKK Danish krone
	"214": true, // DOP Dominican peso
	"222": true, // SVC Salvadoran colón
	"230": true, // ETB Ethiopian birr
	"232": true, // ERN Eritrean nakfa
	"238": true, // FKP Falkland Islands pound
	"242": true, // FJD Fiji dollar
	"262": true, // DJF Djiboutian franc
	"270": true, // GMD Gambian dalasi
	"292": true, // GIP Gibraltar pound
	"320": true, // GTQ Guatemalan quetzal
	"324": true, // GNF Guinean franc
	"328": true, // GYD Guyanese dollar
	"332": true, // HTG Haitian gourde
	"340": true, // HNL Honduran lempira
	"344": true, // HKD Hong Kong dollar
	"348": true, // HUF Hungarian forint
	"352": true, // ISK Icelandic króna (plural: krónur)
	"356": true, // INR Indian rupee
	"360": true, // IDR Indonesian rupiah
	"364": true, // IRR Iranian rial
	"368": true, // IQD Iraqi dinar
	"376": true, // ILS Israeli new shekel
	"388": true, // JMD Jamaican dollar
	"392": true, // JPY Japanese yen
	"398": true, // KZT Kazakhstani tenge
	"400": true, // JOD Jordanian dinar
	"404": true, // KES Kenyan shilling
	"408": true, // KPW North Korean won
	"410": true, // KRW South Korean won
	"414": true, // KWD Kuwaiti dinar
	"417": true, // KGS Kyrgyzstani som
	"418": true, // LAK Lao kip
	"422": true, // LBP Lebanese pound
	"426": true, // LSL Lesotho loti
	"430": true, // LRD Liberian dollar
	"434": true, // LYD Libyan dinar
	"446": true, // MOP Macanese pataca
	"454": true, // MWK Malawian kwacha
	"458": true, // MYR Malaysian ringgit
	"462": true, // MVR Maldivian rufiyaa
	"480": true, // MUR Mauritian rupee
	"484": true, // MXN Mexican peso
	"496": true, // MNT Mongolian tögrög
	"498": true, // MDL Moldovan leu
	"504": true, // MAD Moroccan dirham
	"512": true, // OMR Omani rial
	"516": true, // NAD Namibian dollar
	"524": true, // NPR Nepalese rupee
	"532": true, // ANG Netherlands Antillean guilder
	"533": true, // AWG Aruban florin
	"548": true, // VUV Vanuatu vatu
	"554": true, // NZD New Zealand dollar
	"558": true, // NIO Nicaraguan córdoba
	"566": true, // NGN Nigerian naira
	"578": true, // NOK Norwegian krone
	"586": true, // PKR Pakistani rupee
	"590": true, // PAB Panamanian balboa
	"598": true, // PGK Papua New Guinean kina
	"600": true, // PYG Paraguayan guaraní
	"604": true, // PEN Peruvian sol
	"608": true, // PHP Philippine peso
	"634": true, // QAR Qatari riyal
	"643": true, // RUB Russian ruble
	"646": true, // RWF Rwandan franc
	"654": true, // SHP Saint Helena pound
	"682": true, // SAR Saudi riyal
	"690": true, // SCR Seychelles rupee
	"694": true, // SLL Sierra Leonean leone
	"702": true, // SGD Singapore dollar
	"704": true, // VND Vietnamese đồng
	"706": true, // SOS Somali shilling
	"710": true, // ZAR South African rand
	"728": true, // SSP South Sudanese pound
	"748": true, // SZL Swazi lilangeni
	"752": true, // SEK Swedish krona (plural: kronor)
	"756": true, // CHF Swiss franc
	"760": true, // SYP Syrian pound
	"764": true, // THB Thai baht
	"776": true, // TOP Tongan paʻanga
	"780": true, // TTD Trinidad and Tobago dollar
	"784": true, // AED United Arab Emirates dirham
	"788": true, // TND Tunisian dinar
	"800": true, // UGX Ugandan shilling
	"807": true, // MKD Macedonian denar
	"818": true, // EGP Egyptian pound
	"826": true, // GBP Pound sterling
	"834": true, // TZS Tanzanian shilling
	"840": true, // USD United States dollar
	"858": true, // UYU Uruguayan peso
	"860": true, // UZS Uzbekistan som
	"882": true, // WST Samoan tala
	"886": true, // YER Yemeni rial
	"901": true, // TWD New Taiwan dollar
	"925": true, // SLE Sierra Leonean leone
	"926": true, // VED Venezuelan bolívar digital
	"927": true, // UYW Unidad previsional
	"928": true, // VES Venezuelan bolívar soberano
	"929": true, // MRU Mauritanian ouguiya
	"930": true, // STN São Tomé and Príncipe dobra
	"931": true, // CUC Cuban convertible peso
	"932": true, // ZWL Zimbabwean dollar
	"933": true, // BYN Belarusian ruble
	"934": true, // TMT Turkmenistan manat
	"936": true, // GHS Ghanaian cedi
	"938": true, // SDG Sudanese pound
	"940": true, // UYI Uruguay Peso en Unidades Indexadas (URUIURUI) (funds code)
	"941": true, // RSD Serbian dinar
	"943": true, // MZN Mozambican metical
	"944": true, // AZN Azerbaijani manat
	"946": true, // RON Romanian leu
	"947": true, // CHE WIR euro (complementary currency)
	"948": true, // CHW WIR franc (complementary currency)
	"949": true, // TRY Turkish lira
	"950": true, // XAF CFA franc BEAC
	"951": true, // XCD East Caribbean dollar
	"952": true, // XOF CFA franc BCEAO
	"953": true, // XPF CFP franc (franc Pacifique)	French territories of the Pacific Ocean
	"955": true, // XBA European Composite Unit (EURCO) (bond market unit)
	"956": true, // XBB European Monetary Unit (E.M.U.-6) (bond market unit)
	"957": true, // XBC European Unit of Account 9 (E.U.A.-9) (bond market unit)
	"958": true, // XBD European Unit of Account 17 (E.U.A.-17) (bond market unit)
	"959": true, // XAU Gold (one troy ounce)
	"960": true, // XDR Special drawing rights
	"961": true, // XAG Silver (one troy ounce)
	"962": true, // XPT Platinum (one troy ounce)
	"964": true, // XPD Palladium (one troy ounce)
	"965": true, // XUA ADB Unit of Account	African Development Bank
	"967": true, // ZMW Zambian kwacha
	"968": true, // SRD Surinamese dollar
	"969": true, // MGA Malagasy ariary
	"970": true, // COU Unidad de Valor Real (UVR) (funds code)
	"971": true, // AFN Afghan afghani
	"972": true, // TJS Tajikistani somoni
	"973": true, // AOA Angolan kwanza
	"975": true, // BGN Bulgarian lev
	"976": true, // CDF Congolese franc
	"977": true, // BAM Bosnia and Herzegovina convertible mark
	"978": true, // EUR Euro
	"979": true, // MXV Mexican Unidad de Inversion (UDI) (funds code)
	"980": true, // UAH Ukrainian hryvnia
	"981": true, // GEL Georgian lari
	"984": true, // BOV Bolivian Mvdol (funds code)
	"985": true, // PLN Polish złoty
	"986": true, // BRL Brazilian real
	"990": true, // CLF Unidad de Fomento (funds code)
	"994": true, // XSU Unified System for Regional Compensation (SUCRE)
	"997": true, // USN United States dollar (next day) (funds code)
	//"963": true, // XTS Code reserved for testing
	//"999": true, // XXX No currency
}

var iSO4217CurrencyCodesHistorical = map[string]bool{
	"ADP": true, // Andorran peseta
	"AFA": true, // Afghan afghani
	"ALK": true, // Old Albanian lek
	"AOK": true, // Angolan kwanza
	"AON": true, // Angolan novo kwanza
	"AOR": true, // Angolan kwanza reajustado
	"ARA": true, // Argentine austral
	"ARP": true, // Argentine peso argentino
	"ARY": true, // Argentine peso
	"ATS": true, // Austrian schilling
	"AYM": true, // Azerbaijani manat
	"AZM": true, // Azerbaijani manat
	"BAD": true, // Bosnia and Herzegovina dinar
	"BEC": true, // Belgian convertible franc (funds code)
	"BEF": true, // Belgian franc
	"BEL": true, // Belgian financial franc (funds code)
	"BGJ": true, // Bulgarian lev (first)
	"BGK": true, // Bulgarian lev (second)
	"BGL": true, // Bulgarian lev (third)
	"BOP": true, // Bolivian peso
	"BRB": true, // Brazilian cruzeiro
	"BRC": true, // Brazilian cruzado
	"BRE": true, // Brazilian cruzeiro
	"BRN": true, // Brazilian cruzado novo
	"BRR": true, // Brazilian cruzeiro real
	"BUK": true, // Burmese kyat
	"BYB": true, // Belarusian ruble
	"BYR": true, // Belarusian ruble
	"CHC": true, // WIR franc (for electronic currency)
	"CSD": true, // Serbian dinar
	"CSJ": true, // Czechoslovak koruna (second)
	"CSK": true, // Czechoslovak koruna
	"CYP": true, // Cypriot pound
	"DDM": true, // East German mark
	"DEM": true, // German mark
	"ECS": true, // Ecuadorian sucre
	"ECV": true, // Ecuador Unidad de Valor Constante (funds code)
	"EEK": true, // Estonian kroon
	"ESA": true, // Spanish peseta (account A)
	"ESB": true, // Spanish peseta (account B)
	"ESP": true, // Spanish peseta
	"FIM": true, // Finnish markka
	"FRF": true, // French franc
	"GEK": true, // Georgian kuponi
	"GHC": true, // Ghanaian cedi
	"GHP": true, // Ghanaian cedi
	"GNE": true, // Guinean syli
	"GNS": true, // Guinean syli
	"GQE": true, // Equatorial Guinean ekwele
	"GRD": true, // Greek drachma
	"GWE": true, // Guinean escudo
	"GWP": true, // Guinea-Bissau peso
	"HRD": true, // Croatian dinar
	"IEP": true, // Irish pound
	"ILP": true, // Israeli lira
	"ILR": true, // Israeli shekel
	"ISJ": true, // Icelandic króna
	"ITL": true, // Italian lira
	"LAJ": true, // Lao kip
	"LSM": true, // Lesotho loti
	"LTL": true, // Lithuanian litas
	"LTT": true, // Lithuanian talonas[56]
	"LUC": true, // Luxembourg convertible franc (funds code)
	"LUF": true, // Luxembourg franc
	"LUL": true, // Luxembourg financial franc (funds code)
	"LVL": true, // Latvian lats
	"LVR": true, // Latvian rublis
	"MGF": true, // Malagasy franc
	"MLF": true, // Malian franc
	"MRO": true, // Mauritanian ouguiya
	"MTL": true, // Maltese lira
	"MTP": true, // Maltese pound
	"MVQ": true, // Maldivian rupee
	"MXP": true, // Mexican peso
	"MZE": true, // Mozambican escudo
	"MZM": true, // Mozambican metical
	"NIC": true, // Nicaraguan córdoba
	"NLG": true, // Dutch guilder
	"PEH": true, // Peruvian old sol
	"PEI": true, // Peruvian inti
	"PES": true, // Peruvian sol
	"PLZ": true, // Polish zloty
	"PTE": true, // Portuguese escudo
	"RHD": true, // Rhodesian dollar
	"ROK": true, // Romanian leu (second)
	"ROL": true, // Romanian leu (third)
	"RUR": true, // Russian ruble
	"SDD": true, // Sudanese dinar
	"SDP": true, // Sudanese old pound
	"SIT": true, // Slovenian tolar
	"SKK": true, // Slovak koruna
	"SRG": true, // Surinamese guilder
	"STD": true, // São Tomé and Príncipe dobra
	"SUR": true, // Soviet Union ruble
	"TJR": true, // Tajikistani ruble
	"TMM": true, // Turkmenistani manat
	"TPE": true, // Portuguese Timorese escudo
	"TRL": true, // Turkish lira
	"UAK": true, // Ukrainian karbovanets
	"UGS": true, // Ugandan shilling
	"USS": true, // United States dollar (same day) (funds code)[59]
	"UYN": true, // Uruguay peso
	"UYP": true, // Uruguay new peso
	"VEB": true, // Venezuelan bolívar
	"VEF": true, // Venezuelan bolívar fuerte
	"VNC": true, // Old Vietnamese dong
	"XEU": true, // European Currency Unit
	"XFO": true, // Gold franc (special settlement currency)
	"XFU": true, // UIC franc (special settlement currency)
	"XRE": true, // RINET funds code[63]
	"YDD": true, // South Yemeni dinar
	"YUD": true, // Yugoslav dinar
	"YUM": true, // Yugoslav dinar
	"YUN": true, // Yugoslav dinar
	"ZAL": true, // South African financial rand (funds code)
	"ZMK": true, // Zambian kwacha
	"ZRN": true, // Zairean new zaire
	"ZRZ": true, // Zairean zaire
	"ZWC": true, // Rhodesian dollar
	"ZWD": true, // Zimbabwean dollar
	"ZWN": true, // Zimbabwean dollar
	"ZWR": true, // Zimbabwean dollar
}

var iSO4217CurrencyCodesNumericHistorical = map[string]bool{
	"004": true, // Afghan afghani
	"020": true, // Andorran peseta
	"024": true, // Angolan kwanza
	"031": true, // Azerbaijani manat
	"040": true, // Austrian schilling
	"056": true, // Belgian franc
	"070": true, // Bosnia and Herzegovina dinar
	"076": true, // Brazilian cruzeiro
	"100": true, // Bulgarian lev
	"112": true, // Belarusian ruble
	"180": true, // Zairean zaire
	"196": true, // Cypriot pound
	"200": true, // Czechoslovak koruna
	"218": true, // Ecuadorian sucre
	"226": true, // Equatorial Guinean ekwele
	"233": true, // Estonian kroon
	"246": true, // Finnish markka
	"250": true, // French franc
	"268": true, // Georgian kuponi
	"276": true, // German mark
	"278": true, // East German mark
	"288": true, // Ghanaian cedi
	"300": true, // Greek drachma
	"372": true, // Irish pound
	"380": true, // Italian lira
	"428": true, // Latvian lats / Latvian rublis
	"440": true, // Lithuanian litas / Lithuanian talonas
	"442": true, // Luxembourg franc
	"450": true, // Malagasy franc
	"466": true, // Malian franc
	"470": true, // Maltese lira / Maltese pound
	"478": true, // Mauritanian ouguiya
	"508": true, // Mozambican escudo / Mozambican metical
	"528": true, // Dutch guilder
	"616": true, // Polish zloty
	"620": true, // Portuguese escudo
	"624": true, // Guinean escudo / Guinea-Bissau peso
	"626": true, // Portuguese Timorese escudo
	"642": true, // Romanian leu
	"678": true, // São Tomé and Príncipe dobra
	"703": true, // Slovak koruna
	"705": true, // Slovenian tolar
	"716": true, // Rhodesian dollar / Zimbabwean dollar
	"720": true, // South Yemeni dinar
	"724": true, // Spanish peseta
	"736": true, // Sudanese dinar
	"740": true, // Surinamese guilder
	"762": true, // Tajikistani ruble
	"792": true, // Turkish lira
	"795": true, // Turkmenistani manat
	"804": true, // Ukrainian karbovanets
	"810": true, // Russian ruble / Soviet Union ruble
	"862": true, // Venezuelan bolívar
	"890": true, // Yugoslav dinar
	"891": true, // Yugoslav dinar / Serbian dinar
	"894": true, // Zambian kwacha
	"935": true, // Zimbabwean dollar
	"937": true, // Venezuelan bolívar fuerte
	"939": true, // Ghanaian cedi
	"942": true, // Zimbabwean dollar
	"945": true, // Azerbaijani manat
	"954": true, // European Currency Unit
	"974": true, // Belarusian ruble
	"982": true, // Angolan kwanza reajustado
	"983": true, // Ecuador Unidad de Valor Constante (funds code)
	"987": true, // Brazilian cruzeiro real
	"988": true, // Luxembourg financial franc (funds code)
	"989": true, // Luxembourg convertible franc (funds code)
	"991": true, // South African financial rand (funds code)
	"992": true, // Belgian financial franc (funds code)
	"993": true, // Belgian convertible franc (funds code)
	"995": true, // Spanish peseta (account B)
	"996": true, // Spanish peseta (account A)
	"998": true, // United States dollar (same day) (funds code)[59]
}

var unofficialCurrencyCodes = map[string]bool{
	"ADF": true, // Andorran franc
	"ARL": true, // Argentine peso ley
	"BDS": true, // (BBD) Barbados dollar
	"CNH": true, // Renminbi (offshore)
	"CNT": true, // Renminbi (offshore)
	"GGP": true, // Guernsey pound
	"IMP": true, // Isle of Man pound
	"JEP": true, // Jersey pound
	"KID": true, // Kiribati dollar
	"MAF": true, // Malian franc
	"MCF": true, // Monégasque franc
	"MKN": true, // Old Macedonian denar
	"NIS": true, // (ILS) Israeli new shekel
	"NTD": true, // (TWD) New Taiwan dollar
	"PRB": true, // Transnistrian ruble
	"RMB": true, // (CNY) Renminbi
	"SLS": true, // Somaliland shilling
	"SML": true, // San Marinese lira
	"STG": true, // (GBP) Sterling
	"TVD": true, // Tuvalu dollar
	"VAL": true, // Vatican lira
	"YUG": true, // Yugoslav dinar
	"YUO": true, // Yugoslav dinar
	"YUR": true, // Reformed Yugoslav dinar
	"ZWB": true, // Zimbabwean bonds
}

var cryptoCurrencyCodes = map[string]bool{
	"ADA":  true, // Ada	Currency on the Cardano platform
	"BCH":  true, // Bitcoin Cash
	"BNB":  true, // Binance	BNB
	"BSV":  true, // Bitcoin SV  (Bitcoin Satoshi Vision)
	"BTC":  true, // Bitcoin	BTC
	"DASH": true, // DASH
	"DOGE": true, // Dogecoin
	"EOS":  true, // EOS
	"ETH":  true, // Ethereum
	"LTC":  true, // Litecoin
	"VTC":  true, // Vertcoin
	"XBT":  true, // Bitcoin	BTC
	"XLM":  true, // Stellar Lumen
	"XMR":  true, // Monero
	"XNO":  true, // Nano
	"XRP":  true, // XRP
	"XTZ":  true, // Tez
	"ZEC":  true, // Zcash
}

var iSO3166_2_CountryCodes = map[string]map[string]bool{
	"AD": { // Andorra
		"02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true,
	},
	"AE": { // United Arab Emirates
		"AZ": true, "AJ": true, "FU": true, "SH": true, "DU": true, "RK": true, "UQ": true,
	},
	"AF": { // Afghanistan
		"BDS": true,
		"BDG": true,
		"BGL": true,
		"BAL": true,
		"BAM": true,
		"DAY": true,
		"FRA": true,
		"FYB": true,
		"GHA": true,
		"GHO": true,
		"HEL": true,
		"HER": true,
		"JOW": true,
		"KAB": true,
		"KAN": true,
		"KAP": true,
		"KHO": true,
		"KNR": true,
		"KDZ": true,
		"LAG": true,
		"LOG": true,
		"NAN": true,
		"NIM": true,
		"NUR": true,
		"PKA": true,
		"PIA": true,
		"PAN": true,
		"PAR": true,
		"SAM": true,
		"SAR": true,
		"TAK": true,
		"URU": true,
		"WAR": true,
		"ZAB": true,
	},
	"AG": { // Antigua and Barbuda
		"03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "10": true, "11": true,
	},
	"AI": {}, // Anguilla
	"AL": { // Albania
		"01": true, "09": true, "02": true, "03": true, "04": true, "05": true,
		"06": true, "07": true, "08": true, "10": true, "11": true, "12": true,
	},
	"AM": { // Armenia
		"ER": true,
		"AG": true,
		"AR": true,
		"AV": true,
		"GR": true,
		"KT": true,
		"LO": true,
		"SH": true,
		"SU": true,
		"TV": true,
		"VD": true,
	},
	"AO": { // Angola
		"BGO": true,
		"BGU": true,
		"BIE": true,
		"CAB": true,
		"CNN": true,
		"HUA": true,
		"HUI": true,
		"CCU": true,
		"CNO": true,
		"CUS": true,
		"LUA": true,
		"LNO": true,
		"LSU": true,
		"MAL": true,
		"MOX": true,
		"NAM": true,
		"UIG": true,
		"ZAI": true,
	},
	"AQ": {}, // Antarctica
	"AR": { // Argentina
		"C": true, "B": true, "K": true, "H": true, "U": true, "X": true, "W": true, "E": true,
		"P": true, "Y": true, "L": true, "F": true, "M": true, "N": true, "Q": true, "R": true,
		"A": true, "J": true, "D": true, "Z": true, "S": true, "G": true, "V": true, "T": true,
	},
	"AS": {}, // American Samoa
	"AT": { // Austria
		"1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true, "8": true, "9": true,
	},
	"AU": { // Australia
		"NSW": true, "QLD": true, "SA": true, "TAS": true, "VIC": true, "WA": true, "ACT": true, "NT": true,
	},
	"AW": {}, // Aruba
	"AX": {}, // Åland Islands
	"AZ": { // Azerbaijan
		"NX":  true,
		"BA":  true,
		"GA":  true,
		"LA":  true,
		"MI":  true,
		"NA":  true,
		"NV":  true,
		"SA":  true,
		"SR":  true,
		"SM":  true,
		"XA":  true,
		"YE":  true,
		"ABS": true,
		"AGC": true,
		"AGM": true,
		"AGS": true,
		"AGA": true,
		"AGU": true,
		"AST": true,
		"BAB": true,
		"BAL": true,
		"BAR": true,
		"BEY": true,
		"BIL": true,
		"CAB": true,
		"CAL": true,
		"CUL": true,
		"DAS": true,
		"FUZ": true,
		"GAD": true,
		"GOR": true,
		"GOY": true,
		"GYG": true,
		"HAC": true,
		"IMI": true,
		"ISM": true,
		"KAL": true,
		"KAN": true,
		"KUR": true,
		"LAC": true,
		"LAN": true,
		"LER": true,
		"MAS": true,
		"NEF": true,
		"OGU": true,
		"ORD": true,
		"QAB": true,
		"QAX": true,
		"QAZ": true,
		"QOB": true,
		"QBA": true,
		"QBI": true,
		"QUS": true,
		"SAT": true,
		"SAB": true,
		"SBN": true,
		"SAD": true,
		"SAH": true,
		"SAK": true,
		"SAL": true,
		"SMI": true,
		"SKR": true,
		"SMX": true,
		"SAR": true,
		"SIY": true,
		"SUS": true,
		"TAR": true,
		"TOV": true,
		"UCA": true,
		"XAC": true,
		"XIZ": true,
		"XCI": true,
		"XVD": true,
		"YAR": true,
		"YEV": true,
		"ZAN": true,
		"ZAQ": true,
		"ZAR": true,
	},
	"BA": { // Bosnia and Herzegovina
		"BIH": true,
		"SRP": true,
		"BRC": true,
	},
	"BB": { // Barbados
		"01": true,
		"02": true,
		"03": true,
		"04": true,
		"05": true,
		"06": true,
		"07": true,
		"08": true,
		"09": true,
		"10": true,
		"11": true,
	},
	"BD": { // Bangladesh
		"A": true, "B": true, "C": true, "D": true, "H": true, "E": true, "F": true, "G": true, "05": true,
		"01": true, "02": true, "06": true, "07": true, "03": true, "04": true, "09": true, "45": true, "10": true, "12": true,
		"11": true, "08": true, "13": true, "14": true, "15": true, "16": true, "19": true, "18": true, "17": true, "20": true,
		"21": true, "22": true, "25": true, "23": true, "24": true, "29": true, "27": true, "26": true, "28": true, "30": true,
		"31": true, "32": true, "36": true, "37": true, "33": true, "39": true, "38": true, "35": true, "34": true, "48": true,
		"43": true, "40": true, "42": true, "44": true, "41": true, "46": true, "47": true, "49": true, "52": true, "51": true,
		"50": true, "53": true, "54": true, "56": true, "55": true, "58": true, "62": true, "57": true, "59": true, "61": true,
		"60": true, "63": true, "64": true,
	},
	"BE": { // Belgium
		"BRU": true,
		"VLG": true,
		"WAL": true,
		"VAN": true,
		"WBR": true,
		"WHT": true,
		"WLG": true,
		"VLI": true,
		"WLX": true,
		"WNA": true,
		"VOV": true,
		"VBR": true,
		"VWV": true,
	},
	"BF": { // Burkina Faso
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true,
		"BAL": true,
		"BAM": true,
		"BAN": true,
		"BAZ": true,
		"BGR": true,
		"BLG": true,
		"BLK": true,
		"COM": true,
		"GAN": true,
		"GNA": true,
		"GOU": true,
		"HOU": true,
		"IOB": true,
		"KAD": true,
		"KEN": true,
		"KMD": true,
		"KMP": true,
		"KOS": true,
		"KOP": true,
		"KOT": true,
		"KOW": true,
		"LER": true,
		"LOR": true,
		"MOU": true,
		"NAO": true,
		"NAM": true,
		"NAY": true,
		"NOU": true,
		"OUB": true,
		"OUD": true,
		"PAS": true,
		"PON": true,
		"SNG": true,
		"SMT": true,
		"SEN": true,
		"SIS": true,
		"SOM": true,
		"SOR": true,
		"TAP": true,
		"TUI": true,
		"YAG": true,
		"YAT": true,
		"ZIR": true,
		"ZON": true,
		"ZOU": true,
	},
	"BG": { // Bulgaria
		"01": true, "02": true, "08": true, "07": true, "26": true, "09": true, "10": true, "11": true, "12": true, "13": true,
		"14": true, "15": true, "16": true, "17": true, "18": true, "27": true, "19": true, "20": true, "21": true, "23": true,
		"22": true, "24": true, "25": true, "03": true, "04": true, "05": true, "06": true, "28": true,
	},
	"BH": { // Bahrain
		"13": true, "14": true, "15": true, "17": true,
	},
	"BI": { // Burundi
		"BB": true,
		"BM": true,
		"BL": true,
		"BR": true,
		"CA": true,
		"CI": true,
		"GI": true,
		"KR": true,
		"KY": true,
		"KI": true,
		"MA": true,
		"MU": true,
		"MY": true,
		"MW": true,
		"NG": true,
		"RM": true,
		"RT": true,
		"RY": true,
	},
	"BJ": { // Benin
		"AL": true,
		"AK": true,
		"AQ": true,
		"BO": true,
		"CO": true,
		"KO": true,
		"DO": true,
		"LI": true,
		"MO": true,
		"OU": true,
		"PL": true,
		"ZO": true,
	},
	"BL": {}, // Saint Barthélemy
	"BM": {}, // Bermuda
	"BN": { // Brunei Darussalam
		"BE": true,
		"BM": true,
		"TE": true,
		"TU": true,
	},
	"BO": { // Bolivia (Plurinational State of)
		"C": true, "H": true, "B": true, "L": true, "O": true, "N": true, "P": true, "S": true, "T": true,
	},
	"BQ": { // Bonaire, Sint Eustatius and Saba
		"BO": true, "SA": true, "SE": true,
	},
	"BR": { // Brazil
		"AC": true,
		"AL": true,
		"AP": true,
		"AM": true,
		"BA": true,
		"CE": true,
		"DF": true,
		"ES": true,
		"GO": true,
		"MA": true,
		"MT": true,
		"MS": true,
		"MG": true,
		"PA": true,
		"PB": true,
		"PR": true,
		"PE": true,
		"PI": true,
		"RJ": true,
		"RN": true,
		"RS": true,
		"RO": true,
		"RR": true,
		"SC": true,
		"SP": true,
		"SE": true,
		"TO": true,
	},
	"BS": { // Bahamas
		"AK": true,
		"BY": true,
		"BI": true,
		"BP": true,
		"CI": true,
		"CO": true,
		"CS": true,
		"CE": true,
		"FP": true,
		"CK": true,
		"EG": true,
		"EX": true,
		"GC": true,
		"HI": true,
		"HT": true,
		"IN": true,
		"LI": true,
		"MC": true,
		"MG": true,
		"MI": true,
		"NP": true,
		"NO": true,
		"NS": true,
		"NE": true,
		"RI": true,
		"RC": true,
		"SS": true,
		"SO": true,
		"SA": true,
		"SE": true,
		"SW": true,
		"WG": true,
	},
	"BT": { // Bhutan
		"33": true,
		"12": true,
		"22": true,
		"GA": true,
		"13": true,
		"44": true,
		"42": true,
		"11": true,
		"43": true,
		"23": true,
		"45": true,
		"14": true,
		"31": true,
		"15": true,
		"41": true,
		"TY": true,
		"32": true,
		"21": true,
		"24": true,
		"34": true,
	},
	"BV": {}, // Bouvet Island
	"BW": { // Botswana
		"CE": true,
		"CH": true,
		"FR": true,
		"GA": true,
		"GH": true,
		"JW": true,
		"KG": true,
		"KL": true,
		"KW": true,
		"LO": true,
		"NE": true,
		"NW": true,
		"SP": true,
		"SE": true,
		"SO": true,
		"ST": true,
	},
	"BY": { // Belarus
		"BR": true,
		"HO": true,
		"HR": true,
		"MA": true,
		"MI": true,
		"VI": true,
		"HM": true,
	},
	"BZ": { // Belize
		"BZ":  true,
		"CY":  true,
		"CZL": true,
		"OW":  true,
		"SC":  true,
		"TOL": true,
	},
	"CA": { // Canada
		"AB": true,
		"BC": true,
		"MB": true,
		"NB": true,
		"NL": true,
		"NS": true,
		"ON": true,
		"PE": true,
		"QC": true,
		"SK": true,
		"NT": true,
		"NU": true,
		"YT": true,
	},
	"CC": {}, // Cocos (Keeling) Islands
	"CD": { // Congo, Democratic Republic of the
		"KN": true,
		"BC": true,
		"EQ": true,
		"KE": true,
		"MA": true,
		"NK": true,
		"SK": true,
		"BU": true,
		"HK": true,
		"HL": true,
		"HU": true,
		"IT": true,
		"KC": true,
		"KG": true,
		"KL": true,
		"KS": true,
		"LO": true,
		"LU": true,
		"MN": true,
		"MO": true,
		"NU": true,
		"SA": true,
		"SU": true,
		"TA": true,
		"TO": true,
		"TU": true,
	},
	"CF": { // Central African Republic
		"BGF": true,
		"BB":  true,
		"BK":  true,
		"HK":  true,
		"HM":  true,
		"HS":  true,
		"KG":  true,
		"LB":  true,
		"MB":  true,
		"NM":  true,
		"MP":  true,
		"UK":  true,
		"AC":  true,
		"OP":  true,
		"VK":  true,
		"KB":  true,
		"SE":  true,
	},
	"CG": { // Congo
		"BZV": true,
		"11":  true,
		"8":   true,
		"15":  true,
		"5":   true,
		"2":   true,
		"7":   true,
		"9":   true,
		"14":  true,
		"16":  true,
		"12":  true,
		"13":  true,
	},
	"CH": { // Switzerland
		"AG": true,
		"AR": true,
		"AI": true,
		"BL": true,
		"BS": true,
		"BE": true,
		"FR": true,
		"GE": true,
		"GL": true,
		"GR": true,
		"JU": true,
		"LU": true,
		"NE": true,
		"NW": true,
		"OW": true,
		"SG": true,
		"SH": true,
		"SZ": true,
		"SO": true,
		"TG": true,
		"TI": true,
		"UR": true,
		"VS": true,
		"VD": true,
		"ZG": true,
		"ZH": true,
	},
	"CI": { // Côte d'Ivoire
		"AB": true,
		"BS": true,
		"CM": true,
		"DN": true,
		"GD": true,
		"LC": true,
		"LG": true,
		"MG": true,
		"SM": true,
		"SV": true,
		"VB": true,
		"WR": true,
		"YM": true,
		"ZZ": true,
	},
	"CK": {}, // Cook Islands
	"CL": { // Chile
		"AI": true,
		"AN": true,
		"AP": true,
		"AT": true,
		"BI": true,
		"CO": true,
		"AR": true,
		"LI": true,
		"LL": true,
		"LR": true,
		"MA": true,
		"ML": true,
		"NB": true,
		"RM": true,
		"TA": true,
		"VS": true,
	},
	"CM": { // Cameroon
		"AD": true,
		"CE": true,
		"ES": true,
		"EN": true,
		"LT": true,
		"NO": true,
		"NW": true,
		"SU": true,
		"SW": true,
		"OU": true,
	},
	"CN": { // China
		"AH": true,
		"BJ": true,
		"CQ": true,
		"FJ": true,
		"GD": true,
		"GS": true,
		"GX": true,
		"GZ": true,
		"HA": true,
		"HB": true,
		"HE": true,
		"HI": true,
		"HK": true,
		"HL": true,
		"HN": true,
		"JL": true,
		"JS": true,
		"JX": true,
		"LN": true,
		"MO": true,
		"NM": true,
		"NX": true,
		"QH": true,
		"SC": true,
		"SD": true,
		"SH": true,
		"SN": true,
		"SX": true,
		"TJ": true,
		"TW": true,
		"XJ": true,
		"XZ": true,
		"YN": true,
		"ZJ": true,
	},
	"CO": { // Colombia
		"DC":  true,
		"AMA": true,
		"ANT": true,
		"ARA": true,
		"ATL": true,
		"BOL": true,
		"BOY": true,
		"CAL": true,
		"CAQ": true,
		"CAS": true,
		"CAU": true,
		"CES": true,
		"COR": true,
		"CUN": true,
		"CHO": true,
		"GUA": true,
		"GUV": true,
		"HUI": true,
		"LAG": true,
		"MAG": true,
		"MET": true,
		"NAR": true,
		"NSA": true,
		"PUT": true,
		"QUI": true,
		"RIS": true,
		"SAP": true,
		"SAN": true,
		"SUC": true,
		"TOL": true,
		"VAC": true,
		"VAU": true,
		"VID": true,
	},
	"CR": { // Costa Rica
		"A": true, "C": true, "G": true, "H": true, "L": true, "P": true, "SJ": true,
	},
	"CU": { // Cuba
		"01": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true,
		"99": true,
	},
	"CV": { // Cabo Verde
		"B":  true,
		"S":  true,
		"BR": true,
		"BV": true,
		"CA": true,
		"CF": true,
		"CR": true,
		"MA": true,
		"MO": true,
		"PA": true,
		"PN": true,
		"PR": true,
		"RB": true,
		"RG": true,
		"RS": true,
		"SD": true,
		"SF": true,
		"SL": true,
		"SM": true,
		"SO": true,
		"SS": true,
		"SV": true,
		"TA": true,
		"TS": true,
	},
	"CW": {}, // Curaçao
	"CX": {}, // Christmas Island
	"CY": { // Cyprus
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true,
	},
	"CZ": { // Czechia
		"10":  true,
		"20":  true,
		"31":  true,
		"32":  true,
		"41":  true,
		"42":  true,
		"51":  true,
		"52":  true,
		"53":  true,
		"63":  true,
		"64":  true,
		"71":  true,
		"72":  true,
		"80":  true,
		"20A": true,
		"20B": true,
		"20C": true,
		"201": true,
		"202": true,
		"203": true,
		"204": true,
		"205": true,
		"206": true,
		"207": true,
		"208": true,
		"209": true,
		"311": true,
		"312": true,
		"313": true,
		"314": true,
		"315": true,
		"316": true,
		"317": true,
		"321": true,
		"322": true,
		"323": true,
		"324": true,
		"325": true,
		"326": true,
		"327": true,
		"411": true,
		"412": true,
		"413": true,
		"421": true,
		"422": true,
		"423": true,
		"424": true,
		"425": true,
		"426": true,
		"427": true,
		"511": true,
		"512": true,
		"513": true,
		"514": true,
		"521": true,
		"522": true,
		"523": true,
		"524": true,
		"525": true,
		"531": true,
		"532": true,
		"533": true,
		"534": true,
		"631": true,
		"632": true,
		"633": true,
		"634": true,
		"635": true,
		"641": true,
		"642": true,
		"643": true,
		"644": true,
		"645": true,
		"646": true,
		"647": true,
		"711": true,
		"712": true,
		"713": true,
		"714": true,
		"715": true,
		"721": true,
		"722": true,
		"723": true,
		"724": true,
		"801": true,
		"802": true,
		"803": true,
		"804": true,
		"805": true,
		"806": true,
	},
	"DE": { // Germany
		"BW": true,
		"BY": true,
		"BE": true,
		"BB": true,
		"HB": true,
		"HH": true,
		"HE": true,
		"MV": true,
		"NI": true,
		"NW": true,
		"RP": true,
		"SL": true,
		"SN": true,
		"ST": true,
		"SH": true,
		"TH": true,
	},
	"DJ": { // Djibouti
		"AS": true,
		"AR": true,
		"DI": true,
		"OB": true,
		"TA": true,
		"DJ": true,
	},
	"DK": { // Denmark
		"81": true, "82": true, "83": true, "84": true, "85": true,
	},
	"DM": { // Dominica
		"02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true, "11": true,
	},
	"DO": { // Dominican Republic
		"33": true, "34": true, "35": true, "36": true, "37": true, "38": true, "39": true, "40": true, "41": true, "42": true,
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true,
		"11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true, "20": true,
		"21": true, "22": true, "23": true, "24": true, "25": true, "26": true, "27": true, "28": true, "29": true, "30": true,
		"31": true, "32": true,
	},
	"DZ": { // Algeria
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true,
		"11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true, "20": true,
		"21": true, "22": true, "23": true, "24": true, "25": true, "26": true, "27": true, "28": true, "29": true, "30": true,
		"31": true, "32": true, "33": true, "34": true, "35": true, "36": true, "37": true, "38": true, "39": true, "40": true,
		"41": true, "42": true, "43": true, "44": true, "45": true, "46": true, "47": true, "48": true,
	},
	"EC": { // Ecuador
		"A":  true,
		"B":  true,
		"C":  true,
		"D":  true,
		"E":  true,
		"F":  true,
		"G":  true,
		"H":  true,
		"I":  true,
		"L":  true,
		"M":  true,
		"N":  true,
		"O":  true,
		"P":  true,
		"R":  true,
		"S":  true,
		"SD": true,
		"SE": true,
		"T":  true,
		"U":  true,
		"W":  true,
		"X":  true,
		"Y":  true,
		"Z":  true,
	},
	"EE": { // Estonia
		"37":  true,
		"39":  true,
		"45":  true,
		"50":  true,
		"52":  true,
		"56":  true,
		"60":  true,
		"64":  true,
		"68":  true,
		"71":  true,
		"74":  true,
		"79":  true,
		"81":  true,
		"84":  true,
		"87":  true,
		"130": true,
		"141": true,
		"142": true,
		"171": true,
		"184": true,
		"191": true,
		"198": true,
		"205": true,
		"214": true,
		"245": true,
		"247": true,
		"251": true,
		"255": true,
		"272": true,
		"283": true,
		"284": true,
		"291": true,
		"293": true,
		"296": true,
		"303": true,
		"305": true,
		"317": true,
		"321": true,
		"338": true,
		"353": true,
		"424": true,
		"430": true,
		"431": true,
		"432": true,
		"441": true,
		"442": true,
		"446": true,
		"478": true,
		"480": true,
		"486": true,
		"503": true,
		"511": true,
		"514": true,
		"528": true,
		"557": true,
		"567": true,
		"586": true,
		"615": true,
		"618": true,
		"622": true,
		"624": true,
		"638": true,
		"651": true,
		"653": true,
		"661": true,
		"663": true,
		"668": true,
		"689": true,
		"698": true,
		"708": true,
		"712": true,
		"714": true,
		"719": true,
		"726": true,
		"732": true,
		"735": true,
		"784": true,
		"792": true,
		"793": true,
		"796": true,
		"803": true,
		"809": true,
		"824": true,
		"834": true,
		"855": true,
		"890": true,
		"897": true,
		"899": true,
		"901": true,
		"903": true,
		"907": true,
		"917": true,
		"919": true,
		"928": true,
	},
	"EG": { // Egypt
		"ALX": true,
		"ASN": true,
		"AST": true,
		"BA":  true,
		"BH":  true,
		"BNS": true,
		"C":   true,
		"DK":  true,
		"DT":  true,
		"FYM": true,
		"GH":  true,
		"GZ":  true,
		"IS":  true,
		"JS":  true,
		"KB":  true,
		"KFS": true,
		"KN":  true,
		"LX":  true,
		"MN":  true,
		"MNF": true,
		"MT":  true,
		"PTS": true,
		"SHG": true,
		"SHR": true,
		"SIN": true,
		"SUZ": true,
		"WAD": true,
	},
	"EH": {}, // Western Sahara
	"ER": { // Eritrea
		"AN": true,
		"DK": true,
		"DU": true,
		"GB": true,
		"MA": true,
		"SK": true,
	},
	"ES": { // Spain
		"AN": true,
		"AR": true,
		"AS": true,
		"CB": true,
		"CE": true,
		"CL": true,
		"CM": true,
		"CN": true,
		"CT": true,
		"EX": true,
		"GA": true,
		"IB": true,
		"MC": true,
		"MD": true,
		"ML": true,
		"NC": true,
		"PV": true,
		"RI": true,
		"VC": true,
		"A":  true,
		"AB": true,
		"AL": true,
		"AV": true,
		"B":  true,
		"BA": true,
		"BI": true,
		"BU": true,
		"C":  true,
		"CA": true,
		"CC": true,
		"CO": true,
		"CR": true,
		"CS": true,
		"CU": true,
		"GC": true,
		"GI": true,
		"GR": true,
		"GU": true,
		"H":  true,
		"HU": true,
		"J":  true,
		"L":  true,
		"LE": true,
		"LO": true,
		"LU": true,
		"M":  true,
		"MA": true,
		"MU": true,
		"NA": true,
		"O":  true,
		"OR": true,
		"P":  true,
		"PM": true,
		"PO": true,
		"S":  true,
		"SA": true,
		"SE": true,
		"SG": true,
		"SO": true,
		"SS": true,
		"T":  true,
		"TE": true,
		"TF": true,
		"TO": true,
		"V":  true,
		"VA": true,
		"VI": true,
		"Z":  true,
		"ZA": true,
	},
	"ET": { // Ethiopia
		"AA": true,
		"AF": true,
		"AM": true,
		"BE": true,
		"DD": true,
		"GA": true,
		"HA": true,
		"OR": true,
		"SI": true,
		"SN": true,
		"SO": true,
		"TI": true,
	},
	"FI": { // Finland
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true,
		"11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true,
	},
	"FJ": { // Fiji
		"C": true, "E": true, "N": true, "R": true, "W": true,
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true,
		"11": true, "12": true, "13": true, "14": true,
	},
	"FK": {}, // Falkland Islands (Malvinas)
	"FM": { // Micronesia (Federated States of)
		"KSA": true,
		"PNI": true,
		"TRK": true,
		"YAP": true,
	},
	"FO": {}, // Faroe Islands
	"FR": { // France
		"20R": true,
		"ARA": true,
		"BFC": true,
		"BRE": true,
		"CVL": true,
		"GES": true,
		"HDF": true,
		"IDF": true,
		"NAQ": true,
		"NOR": true,
		"OCC": true,
		"PAC": true,
		"PDL": true,
		"01":  true, "02": true, "03": true, "04": true, "05": true, "06": true, "6AE": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true,
		"2A": true, "2B": true,
		"21": true, "22": true, "23": true, "24": true, "25": true, "26": true, "27": true, "28": true, "29": true,
		"30": true, "31": true, "32": true, "33": true, "34": true, "35": true, "36": true, "37": true, "38": true, "39": true,
		"40": true, "41": true, "42": true, "43": true, "44": true, "45": true, "46": true, "47": true, "48": true, "49": true,
		"50": true, "51": true, "52": true, "53": true, "54": true, "55": true, "56": true, "57": true, "58": true, "59": true,
		"60": true, "61": true, "62": true, "63": true, "64": true, "65": true, "66": true, "69": true, "69M": true,
		"70": true, "71": true, "72": true, "73": true, "74": true, "75C": true, "76": true, "77": true, "78": true, "79": true,
		"80": true, "81": true, "82": true, "83": true, "84": true, "85": true, "86": true, "87": true, "88": true, "89": true,
		"90": true, "91": true, "92": true, "93": true, "94": true, "95": true,
		"971": true, "972": true, "973": true, "974": true, "976": true,
		"67": true, "68": true,
		"BL": true,
		"CP": true,
		"MF": true,
		"NC": true,
		"PF": true,
		"PM": true,
		"TF": true,
		"WF": true,
	},
	"GA": { // Gabon
		"1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true, "8": true, "9": true,
	},
	"GB": { // United Kingdom of Great Britain and Northern Ireland
		"ENG": true,
		"NIR": true,
		"SCT": true,
		"WLS": true,
		"EAW": true,
		"GBN": true,
		"UKM": true,
		"ABC": true,
		"ABD": true,
		"ABE": true,
		"AGB": true,
		"AGY": true,
		"AND": true,
		"ANN": true,
		"ANS": true,
		"BAS": true,
		"BBD": true,
		"BCP": true,
		"BDF": true,
		"BDG": true,
		"BEN": true,
		"BEX": true,
		"BFS": true,
		"BGE": true,
		"BGW": true,
		"BIR": true,
		"BKM": true,
		"BNE": true,
		"BNH": true,
		"BNS": true,
		"BOL": true,
		"BPL": true,
		"BRC": true,
		"BRD": true,
		"BRY": true,
		"BST": true,
		"BUR": true,
		"CAM": true,
		"CAY": true,
		"CBF": true,
		"CCG": true,
		"CGN": true,
		"CHE": true,
		"CHW": true,
		"CLD": true,
		"CLK": true,
		"CMA": true,
		"CMD": true,
		"CMN": true,
		"CON": true,
		"COV": true,
		"CRF": true,
		"CRY": true,
		"CWY": true,
		"DAL": true,
		"DBY": true,
		"DEN": true,
		"DER": true,
		"DEV": true,
		"DGY": true,
		"DNC": true,
		"DND": true,
		"DOR": true,
		"DRS": true,
		"DUD": true,
		"DUR": true,
		"EAL": true,
		"EAY": true,
		"EDH": true,
		"EDU": true,
		"ELN": true,
		"ELS": true,
		"ENF": true,
		"ERW": true,
		"ERY": true,
		"ESS": true,
		"ESX": true,
		"FAL": true,
		"FIF": true,
		"FLN": true,
		"FMO": true,
		"GAT": true,
		"GLG": true,
		"GLS": true,
		"GRE": true,
		"GWN": true,
		"HAL": true,
		"HAM": true,
		"HAV": true,
		"HCK": true,
		"HEF": true,
		"HIL": true,
		"HLD": true,
		"HMF": true,
		"HNS": true,
		"HPL": true,
		"HRT": true,
		"HRW": true,
		"HRY": true,
		"IOS": true,
		"IOW": true,
		"ISL": true,
		"IVC": true,
		"KEC": true,
		"KEN": true,
		"KHL": true,
		"KIR": true,
		"KTT": true,
		"KWL": true,
		"LAN": true,
		"LBC": true,
		"LBH": true,
		"LCE": true,
		"LDS": true,
		"LEC": true,
		"LEW": true,
		"LIN": true,
		"LIV": true,
		"LND": true,
		"LUT": true,
		"MAN": true,
		"MDB": true,
		"MDW": true,
		"MEA": true,
		"MIK": true,
		"MLN": true,
		"MON": true,
		"MRT": true,
		"MRY": true,
		"MTY": true,
		"MUL": true,
		"NAY": true,
		"NBL": true,
		"NEL": true,
		"NET": true,
		"NFK": true,
		"NGM": true,
		"NLK": true,
		"NLN": true,
		"NMD": true,
		"NSM": true,
		"NTH": true,
		"NTL": true,
		"NTT": true,
		"NTY": true,
		"NWM": true,
		"NWP": true,
		"NYK": true,
		"OLD": true,
		"ORK": true,
		"OXF": true,
		"PEM": true,
		"PKN": true,
		"PLY": true,
		"POR": true,
		"POW": true,
		"PTE": true,
		"RCC": true,
		"RCH": true,
		"RCT": true,
		"RDB": true,
		"RDG": true,
		"RFW": true,
		"RIC": true,
		"ROT": true,
		"RUT": true,
		"SAW": true,
		"SAY": true,
		"SCB": true,
		"SFK": true,
		"SFT": true,
		"SGC": true,
		"SHF": true,
		"SHN": true,
		"SHR": true,
		"SKP": true,
		"SLF": true,
		"SLG": true,
		"SLK": true,
		"SND": true,
		"SOL": true,
		"SOM": true,
		"SOS": true,
		"SRY": true,
		"STE": true,
		"STG": true,
		"STH": true,
		"STN": true,
		"STS": true,
		"STT": true,
		"STY": true,
		"SWA": true,
		"SWD": true,
		"SWK": true,
		"TAM": true,
		"TFW": true,
		"THR": true,
		"TOB": true,
		"TOF": true,
		"TRF": true,
		"TWH": true,
		"VGL": true,
		"WAR": true,
		"WBK": true,
		"WDU": true,
		"WFT": true,
		"WGN": true,
		"WIL": true,
		"WKF": true,
		"WLL": true,
		"WLN": true,
		"WLV": true,
		"WND": true,
		"WNM": true,
		"WOK": true,
		"WOR": true,
		"WRL": true,
		"WRT": true,
		"WRX": true,
		"WSM": true,
		"WSX": true,
		"YOR": true,
		"ZET": true,
	},
	"GD": { // Grenada
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "10": true,
	},
	"GE": { // Georgia
		"AB": true,
		"AJ": true,
		"GU": true,
		"IM": true,
		"KA": true,
		"KK": true,
		"MM": true,
		"RL": true,
		"SJ": true,
		"SK": true,
		"SZ": true,
		"TB": true,
	},
	"GF": {}, // French Guiana
	"GG": {}, // Guernsey
	"GH": { // Ghana
		"AA": true,
		"AF": true,
		"AH": true,
		"BE": true,
		"BO": true,
		"CP": true,
		"EP": true,
		"NE": true,
		"NP": true,
		"OT": true,
		"SV": true,
		"TV": true,
		"UE": true,
		"UW": true,
		"WN": true,
		"WP": true,
	},
	"GI": {}, // Gibraltar
	"GL": { // Greenland
		"AV": true,
		"KU": true,
		"QE": true,
		"QT": true,
		"SM": true,
	},
	"GM": { // Gambia
		"B": true,
		"L": true,
		"M": true,
		"N": true,
		"U": true,
		"W": true,
	},
	"GN": { // Guinea
		"B":  true,
		"C":  true,
		"D":  true,
		"F":  true,
		"K":  true,
		"L":  true,
		"M":  true,
		"N":  true,
		"BE": true,
		"BF": true,
		"BK": true,
		"CO": true,
		"DB": true,
		"DI": true,
		"DL": true,
		"DU": true,
		"FA": true,
		"FO": true,
		"FR": true,
		"GA": true,
		"GU": true,
		"KA": true,
		"KB": true,
		"KD": true,
		"KE": true,
		"KN": true,
		"KO": true,
		"KS": true,
		"LA": true,
		"LE": true,
		"LO": true,
		"MC": true,
		"MD": true,
		"ML": true,
		"MM": true,
		"NZ": true,
		"PI": true,
		"SI": true,
		"TE": true,
		"TO": true,
		"YO": true,
	},
	"GP": {}, // Guadeloupe
	"GQ": { // Equatorial Guinea
		"C":  true,
		"I":  true,
		"AN": true,
		"BN": true,
		"BS": true,
		"CS": true,
		"DJ": true,
		"KN": true,
		"LI": true,
		"WN": true,
	},
	"GR": { // Greece
		"A":  true,
		"B":  true,
		"C":  true,
		"D":  true,
		"E":  true,
		"F":  true,
		"G":  true,
		"H":  true,
		"I":  true,
		"J":  true,
		"K":  true,
		"L":  true,
		"M":  true,
		"69": true,
	},
	"GS": {}, // South Georgia and the South Sandwich Islands
	"GT": { // Guatemala
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true,
		"11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true, "20": true,
		"21": true, "22": true,
	},
	"GU": {}, // Guam
	"GW": { // Guinea-Bissau
		"L":  true,
		"N":  true,
		"S":  true,
		"BA": true,
		"BL": true,
		"BM": true,
		"BS": true,
		"CA": true,
		"GA": true,
		"OI": true,
		"QU": true,
		"TO": true,
	},
	"GY": { // Guyana
		"BA": true,
		"CU": true,
		"DE": true,
		"EB": true,
		"ES": true,
		"MA": true,
		"PM": true,
		"PT": true,
		"UD": true,
		"UT": true,
	},
	"HK": {}, // Hong Kong
	"HM": {}, // Heard Island and McDonald Islands
	"HN": { // Honduras
		"AT": true,
		"CH": true,
		"CL": true,
		"CM": true,
		"CP": true,
		"CR": true,
		"EP": true,
		"FM": true,
		"GD": true,
		"IB": true,
		"IN": true,
		"LE": true,
		"LP": true,
		"OC": true,
		"OL": true,
		"SB": true,
		"VA": true,
		"YO": true,
	},
	"HR": { // Croatia
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true,
		"11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true, "20": true,
		"21": true,
	},
	"HT": { // Haiti
		"AR": true,
		"CE": true,
		"GA": true,
		"ND": true,
		"NE": true,
		"NI": true,
		"NO": true,
		"OU": true,
		"SD": true,
		"SE": true,
	},
	"HU": { // Hungary
		"BA": true,
		"BC": true,
		"BE": true,
		"BK": true,
		"BU": true,
		"BZ": true,
		"CS": true,
		"DE": true,
		"DU": true,
		"EG": true,
		"ER": true,
		"FE": true,
		"GS": true,
		"GY": true,
		"HB": true,
		"HE": true,
		"HV": true,
		"JN": true,
		"KE": true,
		"KM": true,
		"KV": true,
		"MI": true,
		"NK": true,
		"NO": true,
		"NY": true,
		"PE": true,
		"PS": true,
		"SD": true,
		"SF": true,
		"SH": true,
		"SK": true,
		"SN": true,
		"SO": true,
		"SS": true,
		"ST": true,
		"SZ": true,
		"TB": true,
		"TO": true,
		"VA": true,
		"VE": true,
		"VM": true,
		"ZA": true,
		"ZE": true,
	},
	"ID": { // Indonesia
		"JW": true,
		"KA": true,
		"ML": true,
		"NU": true,
		"PP": true,
		"SL": true,
		"SM": true,
		"AC": true,
		"BA": true,
		"BB": true,
		"BE": true,
		"BT": true,
		"GO": true,
		"JA": true,
		"JB": true,
		"JI": true,
		"JK": true,
		"JT": true,
		"KB": true,
		"KI": true,
		"KR": true,
		"KS": true,
		"KT": true,
		"KU": true,
		"LA": true,
		"MA": true,
		"MU": true,
		"NB": true,
		"NT": true,
		"PA": true,
		"PB": true,
		"RI": true,
		"SA": true,
		"SB": true,
		"SG": true,
		"SN": true,
		"SR": true,
		"SS": true,
		"ST": true,
		"SU": true,
		"YO": true,
	},
	"IE": { // Ireland
		"C":  true,
		"L":  true,
		"M":  true,
		"U":  true,
		"CE": true,
		"CN": true,
		"CO": true,
		"CW": true,
		"D":  true,
		"DL": true,
		"G":  true,
		"KE": true,
		"KK": true,
		"KY": true,
		"LD": true,
		"LH": true,
		"LK": true,
		"LM": true,
		"LS": true,
		"MH": true,
		"MN": true,
		"MO": true,
		"OY": true,
		"RN": true,
		"SO": true,
		"TA": true,
		"WD": true,
		"WH": true,
		"WW": true,
		"WX": true,
	},
	"IL": { // Israel
		"D":  true,
		"HA": true,
		"JM": true,
		"M":  true,
		"TA": true,
		"Z":  true,
	},
	"IM": {}, // Isle of Man
	"IN": { // India
		"AN": true,
		"AP": true,
		"AR": true,
		"AS": true,
		"BR": true,
		"CH": true,
		"CT": true,
		"DH": true,
		"DL": true,
		"GA": true,
		"GJ": true,
		"HP": true,
		"HR": true,
		"JH": true,
		"JK": true,
		"KA": true,
		"KL": true,
		"LA": true,
		"LD": true,
		"MH": true,
		"ML": true,
		"MN": true,
		"MP": true,
		"MZ": true,
		"NL": true,
		"OR": true,
		"PB": true,
		"PY": true,
		"RJ": true,
		"SK": true,
		"TG": true,
		"TN": true,
		"TR": true,
		"UP": true,
		"UT": true,
		"WB": true,
	},
	"IO": {}, // British Indian Ocean Territory
	"IQ": { // Iraq
		"AN": true,
		"AR": true,
		"BA": true,
		"BB": true,
		"BG": true,
		"DA": true,
		"DI": true,
		"DQ": true,
		"KA": true,
		"KI": true,
		"KR": true,
		"MA": true,
		"MU": true,
		"NA": true,
		"NI": true,
		"QA": true,
		"SD": true,
		"SU": true,
		"WA": true,
	},
	"IR": { // Iran (Islamic Republic of)
		"00": true, "01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true,
		"20": true, "21": true, "22": true, "23": true, "24": true, "25": true, "26": true, "27": true, "28": true, "29": true,
		"30": true,
	},
	"IS": { // Iceland
		"1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true, "8": true,
		"AKH": true,
		"AKN": true,
		"AKU": true,
		"ARN": true,
		"ASA": true,
		"BLA": true,
		"BLO": true,
		"BOG": true,
		"BOL": true,
		"DAB": true,
		"DAV": true,
		"EOM": true,
		"EYF": true,
		"FJD": true,
		"FJL": true,
		"FLA": true,
		"FLR": true,
		"GAR": true,
		"GOG": true,
		"GRN": true,
		"GRU": true,
		"GRY": true,
		"HAF": true,
		"HEL": true,
		"HRG": true,
		"HRU": true,
		"HUT": true,
		"HUV": true,
		"HVA": true,
		"HVE": true,
		"ISA": true,
		"KAL": true,
		"KJO": true,
		"KOP": true,
		"LAN": true,
		"MOS": true,
		"MUL": true,
		"MYR": true,
		"NOR": true,
		"RGE": true,
		"RGY": true,
		"RHH": true,
		"RKN": true,
		"RKV": true,
		"SBH": true,
		"SBT": true,
		"SDN": true,
		"SDV": true,
		"SEL": true,
		"SFA": true,
		"SHF": true,
		"SKF": true,
		"SKG": true,
		"SKO": true,
		"SKU": true,
		"SNF": true,
		"SOG": true,
		"SOL": true,
		"SSF": true,
		"SSS": true,
		"STR": true,
		"STY": true,
		"SVG": true,
		"TAL": true,
		"THG": true,
		"TJO": true,
		"VEM": true,
		"VER": true,
		"VOP": true,
	},
	"IT": { // Italy
		"21": true,
		"25": true,
		"34": true,
		"42": true,
		"45": true,
		"52": true,
		"55": true,
		"57": true,
		"62": true,
		"65": true,
		"67": true,
		"72": true,
		"75": true,
		"77": true,
		"78": true,
		"23": true,
		"32": true,
		"36": true,
		"82": true,
		"88": true,
		"AL": true,
		"AN": true,
		"AP": true,
		"AQ": true,
		"AR": true,
		"AT": true,
		"AV": true,
		"BG": true,
		"BI": true,
		"BL": true,
		"BN": true,
		"BR": true,
		"BS": true,
		"BT": true,
		"CB": true,
		"CE": true,
		"CH": true,
		"CN": true,
		"CO": true,
		"CR": true,
		"CS": true,
		"CZ": true,
		"FC": true,
		"FE": true,
		"FG": true,
		"FM": true,
		"FR": true,
		"GR": true,
		"IM": true,
		"IS": true,
		"KR": true,
		"LC": true,
		"LE": true,
		"LI": true,
		"LO": true,
		"LT": true,
		"LU": true,
		"MB": true,
		"MC": true,
		"MN": true,
		"MO": true,
		"MS": true,
		"MT": true,
		"NO": true,
		"NU": true,
		"OR": true,
		"PC": true,
		"PD": true,
		"PE": true,
		"PG": true,
		"PI": true,
		"PO": true,
		"PR": true,
		"PT": true,
		"PU": true,
		"PV": true,
		"PZ": true,
		"RA": true,
		"RE": true,
		"RI": true,
		"RN": true,
		"RO": true,
		"SA": true,
		"SI": true,
		"SO": true,
		"SP": true,
		"SS": true,
		"SU": true,
		"SV": true,
		"TA": true,
		"TE": true,
		"TR": true,
		"TV": true,
		"VA": true,
		"VB": true,
		"VC": true,
		"VI": true,
		"VR": true,
		"VT": true,
		"VV": true,
		"BZ": true,
		"TN": true,
		"GO": true,
		"PN": true,
		"TS": true,
		"UD": true,
		"AG": true,
		"CL": true,
		"EN": true,
		"RG": true,
		"SR": true,
		"TP": true,
		"BA": true,
		"BO": true,
		"CA": true,
		"CT": true,
		"FI": true,
		"GE": true,
		"ME": true,
		"MI": true,
		"NA": true,
		"PA": true,
		"RC": true,
		"RM": true,
		"TO": true,
		"VE": true,
	},
	"JE": {}, // Jersey
	"JM": { // Jamaica
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true,
		"11": true, "12": true, "13": true, "14": true,
	},
	"JO": { // Jordan
		"AJ": true,
		"AM": true,
		"AQ": true,
		"AT": true,
		"AZ": true,
		"BA": true,
		"IR": true,
		"JA": true,
		"KA": true,
		"MA": true,
		"MD": true,
		"MN": true,
	},
	"JP": { // Japan
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true,
		"11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true, "20": true,
		"21": true, "22": true, "23": true, "24": true, "25": true, "26": true, "27": true, "28": true, "29": true, "30": true,
		"31": true, "32": true, "33": true, "34": true, "35": true, "36": true, "37": true, "38": true, "39": true, "40": true,
		"41": true, "42": true, "43": true, "44": true, "45": true, "46": true, "47": true,
	},
	"KE": { // Kenya
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true,
		"11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true, "20": true,
		"21": true, "22": true, "23": true, "24": true, "25": true, "26": true, "27": true, "28": true, "29": true, "30": true,
		"31": true, "32": true, "33": true, "34": true, "35": true, "36": true, "37": true, "38": true, "39": true, "40": true,
		"41": true, "42": true, "43": true, "44": true, "45": true, "46": true, "47": true,
	},
	"KG": { // Kyrgyzstan
		"B":  true,
		"C":  true,
		"GB": true,
		"GO": true,
		"J":  true,
		"N":  true,
		"O":  true,
		"T":  true,
		"Y":  true,
	},
	"KH": { // Cambodia
		"1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true, "8": true, "9": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true,
		"20": true, "21": true, "22": true, "23": true, "24": true, "25": true,
	},
	"KI": { // Kiribati
		"G": true, "L": true, "P": true,
	},
	"KM": { // Comoros
		"A": true, "G": true, "M": true,
	},
	"KN": { // Saint Kitts and Nevis
		"K":  true,
		"N":  true,
		"01": true,
		"02": true,
		"03": true,
		"04": true,
		"05": true,
		"06": true,
		"07": true,
		"08": true,
		"09": true,
		"10": true,
		"11": true,
		"12": true,
		"13": true,
		"15": true,
	},
	"KP": { // Korea (Democratic People's Republic of)
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "13": true, "14": true,
	},
	"KR": { // Korea, Republic of
		"11": true,
		"26": true,
		"27": true,
		"28": true,
		"29": true,
		"30": true,
		"31": true,
		"41": true,
		"42": true,
		"43": true,
		"44": true,
		"45": true,
		"46": true,
		"47": true,
		"48": true,
		"49": true,
		"50": true,
	},
	"KW": { // Kuwait
		"AH": true,
		"FA": true,
		"HA": true,
		"JA": true,
		"KU": true,
		"MU": true,
	},
	"KY": {}, // Cayman Islands
	"KZ": { // Kazakhstan
		"AKM": true,
		"AKT": true,
		"ALA": true,
		"ALM": true,
		"AST": true,
		"ATY": true,
		"KAR": true,
		"KUS": true,
		"KZY": true,
		"MAN": true,
		"PAV": true,
		"SEV": true,
		"SHY": true,
		"VOS": true,
		"YUZ": true,
		"ZAP": true,
		"ZHA": true,
	},
	"LA": { // Lao People's Democratic Republic
		"AT": true,
		"BK": true,
		"BL": true,
		"CH": true,
		"HO": true,
		"KH": true,
		"LM": true,
		"LP": true,
		"OU": true,
		"PH": true,
		"SL": true,
		"SV": true,
		"VI": true,
		"VT": true,
		"XA": true,
		"XE": true,
		"XI": true,
		"XS": true,
	},
	"LB": { // Lebanon
		"AK": true,
		"AS": true,
		"BA": true,
		"BH": true,
		"BI": true,
		"JA": true,
		"JL": true,
		"NA": true,
	},
	"LC": { // Saint Lucia
		"01": true, "02": true, "03": true, "05": true, "06": true, "07": true, "08": true, "10": true, "11": true, "12": true,
	},
	"LI": { // Liechtenstein
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true, "11": true,
	},
	"LK": { // Sri Lanka
		"1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true, "8": true, "9": true,
		"11": true, "12": true, "13": true,
		"21": true, "22": true, "23": true,
		"31": true, "32": true, "33": true,
		"41": true, "42": true, "43": true, "44": true, "45": true,
		"51": true, "52": true, "53": true,
		"61": true, "62": true,
		"71": true, "72": true,
		"81": true, "82": true,
		"91": true, "92": true,
	},
	"LR": { // Liberia
		"BG": true,
		"BM": true,
		"CM": true,
		"GB": true,
		"GG": true,
		"GK": true,
		"GP": true,
		"LO": true,
		"MG": true,
		"MO": true,
		"MY": true,
		"NI": true,
		"RG": true,
		"RI": true,
		"SI": true,
	},
	"LS": { // Lesotho
		"A": true, "B": true, "C": true, "D": true, "E": true, "F": true, "G": true, "H": true, "J": true, "K": true,
	},
	"LT": { // Lithuania
		"AL": true,
		"KL": true,
		"KU": true,
		"MR": true,
		"PN": true,
		"SA": true,
		"TA": true,
		"TE": true,
		"UT": true,
		"VL": true,
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true,
		"20": true, "21": true, "22": true, "23": true, "24": true, "25": true, "26": true, "27": true, "28": true, "29": true,
		"30": true, "31": true, "32": true, "33": true, "34": true, "35": true, "36": true, "37": true, "38": true, "39": true,
		"40": true, "41": true, "42": true, "43": true, "44": true, "45": true, "46": true, "47": true, "48": true, "49": true,
		"50": true, "51": true, "52": true, "53": true, "54": true, "55": true, "56": true, "57": true, "58": true, "59": true,
		"60": true,
	},
	"LU": { // Luxembourg
		"CA": true,
		"CL": true,
		"DI": true,
		"EC": true,
		"ES": true,
		"GR": true,
		"LU": true,
		"ME": true,
		"RD": true,
		"RM": true,
		"VD": true,
		"WI": true,
	},
	"LV": { // Latvia
		"002": true,
		"007": true,
		"011": true,
		"015": true,
		"016": true,
		"022": true,
		"026": true,
		"033": true,
		"041": true,
		"042": true,
		"047": true,
		"050": true,
		"052": true,
		"054": true,
		"056": true,
		"058": true,
		"059": true,
		"062": true,
		"067": true,
		"068": true,
		"073": true,
		"077": true,
		"080": true,
		"087": true,
		"088": true,
		"089": true,
		"091": true,
		"094": true,
		"097": true,
		"099": true,
		"101": true,
		"102": true,
		"106": true,
		"111": true,
		"112": true,
		"113": true,
		"DGV": true,
		"JEL": true,
		"JUR": true,
		"LPX": true,
		"REZ": true,
		"RIX": true,
		"VEN": true,
	},
	"LY": { // Libya
		"BA": true,
		"BU": true,
		"DR": true,
		"GT": true,
		"JA": true,
		"JG": true,
		"JI": true,
		"JU": true,
		"KF": true,
		"MB": true,
		"MI": true,
		"MJ": true,
		"MQ": true,
		"NL": true,
		"NQ": true,
		"SB": true,
		"SR": true,
		"TB": true,
		"WA": true,
		"WD": true,
		"WS": true,
		"ZA": true,
	},
	"MA": { // Morocco
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true, "11": true, "12": true,
		"AGD": true,
		"AOU": true,
		"ASZ": true,
		"AZI": true,
		"BEM": true,
		"BER": true,
		"BES": true,
		"BOD": true,
		"BOM": true,
		"BRR": true,
		"CAS": true,
		"CHE": true,
		"CHI": true,
		"CHT": true,
		"DRI": true,
		"ERR": true,
		"ESI": true,
		"ESM": true,
		"FAH": true,
		"FES": true,
		"FIG": true,
		"FQH": true,
		"GUE": true,
		"GUF": true,
		"HAJ": true,
		"HAO": true,
		"HOC": true,
		"IFR": true,
		"INE": true,
		"JDI": true,
		"JRA": true,
		"KEN": true,
		"KES": true,
		"KHE": true,
		"KHN": true,
		"KHO": true,
		"LAA": true,
		"LAR": true,
		"MAR": true,
		"MDF": true,
		"MED": true,
		"MEK": true,
		"MID": true,
		"MOH": true,
		"MOU": true,
		"NAD": true,
		"NOU": true,
		"OUA": true,
		"OUD": true,
		"OUJ": true,
		"OUZ": true,
		"RAB": true,
		"REH": true,
		"SAF": true,
		"SAL": true,
		"SEF": true,
		"SET": true,
		"SIB": true,
		"SIF": true,
		"SIK": true,
		"SIL": true,
		"SKH": true,
		"TAF": true,
		"TAI": true,
		"TAO": true,
		"TAR": true,
		"TAT": true,
		"TAZ": true,
		"TET": true,
		"TIN": true,
		"TIZ": true,
		"TNG": true,
		"TNT": true,
		"YUS": true,
		"ZAG": true,
	},
	"MC": { // Monaco
		"CL": true,
		"CO": true,
		"FO": true,
		"GA": true,
		"JE": true,
		"LA": true,
		"MA": true,
		"MC": true,
		"MG": true,
		"MO": true,
		"MU": true,
		"PH": true,
		"SD": true,
		"SO": true,
		"SP": true,
		"SR": true,
		"VR": true,
	},
	"MD": { // Moldova, Republic of
		"AN": true,
		"BA": true,
		"BD": true,
		"BR": true,
		"BS": true,
		"CA": true,
		"CL": true,
		"CM": true,
		"CR": true,
		"CS": true,
		"CT": true,
		"CU": true,
		"DO": true,
		"DR": true,
		"DU": true,
		"ED": true,
		"FA": true,
		"FL": true,
		"GA": true,
		"GL": true,
		"HI": true,
		"IA": true,
		"LE": true,
		"NI": true,
		"OC": true,
		"OR": true,
		"RE": true,
		"RI": true,
		"SD": true,
		"SI": true,
		"SN": true,
		"SO": true,
		"ST": true,
		"SV": true,
		"TA": true,
		"TE": true,
		"UN": true,
	},
	"ME": { // Montenegro
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true,
		"11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true, "20": true,
		"21": true, "22": true, "23": true, "24": true,
	},
	"MF": {}, // Saint Martin (French part)
	"MG": { // Madagascar
		"A": true, "D": true, "F": true, "M": true, "T": true, "U": true,
	},
	"MH": { // Marshall Islands
		"L":   true,
		"T":   true,
		"ALK": true,
		"ALL": true,
		"ARN": true,
		"AUR": true,
		"EBO": true,
		"ENI": true,
		"JAB": true,
		"JAL": true,
		"KIL": true,
		"KWA": true,
		"LAE": true,
		"LIB": true,
		"LIK": true,
		"MAJ": true,
		"MAL": true,
		"MEJ": true,
		"MIL": true,
		"NMK": true,
		"NMU": true,
		"RON": true,
		"UJA": true,
		"UTI": true,
		"WTH": true,
		"WTJ": true,
	},
	"MK": { // North Macedonia
		"101": true,
		"102": true,
		"103": true,
		"104": true,
		"105": true,
		"106": true,
		"107": true,
		"108": true,
		"109": true,
		"201": true,
		"202": true,
		"203": true,
		"204": true,
		"205": true,
		"206": true,
		"207": true,
		"208": true,
		"209": true,
		"210": true,
		"211": true,
		"301": true,
		"303": true,
		"304": true,
		"307": true,
		"308": true,
		"310": true,
		"311": true,
		"312": true,
		"313": true,
		"401": true,
		"402": true,
		"403": true,
		"404": true,
		"405": true,
		"406": true,
		"407": true,
		"408": true,
		"409": true,
		"410": true,
		"501": true,
		"502": true,
		"503": true,
		"504": true,
		"505": true,
		"506": true,
		"507": true,
		"508": true,
		"509": true,
		"601": true,
		"602": true,
		"603": true,
		"604": true,
		"605": true,
		"606": true,
		"607": true,
		"608": true,
		"609": true,
		"701": true,
		"702": true,
		"703": true,
		"704": true,
		"705": true,
		"706": true,
		"801": true,
		"802": true,
		"803": true,
		"804": true,
		"805": true,
		"806": true,
		"807": true,
		"808": true,
		"809": true,
		"810": true,
		"811": true,
		"812": true,
		"813": true,
		"814": true,
		"815": true,
		"816": true,
		"817": true,
	},
	"ML": { // Mali
		"1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true, "8": true, "9": true, "10": true,
		"BKO": true,
	},
	"MM": { // Myanmar
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true,
		"11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true,
	},
	"MN": { // Mongolia
		"1":   true,
		"035": true,
		"037": true,
		"039": true,
		"041": true,
		"043": true,
		"046": true,
		"047": true,
		"049": true,
		"051": true,
		"053": true,
		"055": true,
		"057": true,
		"059": true,
		"061": true,
		"063": true,
		"064": true,
		"065": true,
		"067": true,
		"069": true,
		"071": true,
		"073": true,
	},
	"MO": {}, // Macao
	"MP": {}, // Northern Mariana Islands
	"MQ": {}, // Martinique
	"MR": { // Mauritania
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true,
	},
	"MS": {}, // Montserrat
	"MT": { // Malta
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true,
		"20": true, "21": true, "22": true, "23": true, "24": true, "25": true, "26": true, "27": true, "28": true, "29": true,
		"30": true, "31": true, "32": true, "33": true, "34": true, "35": true, "36": true, "37": true, "38": true, "39": true,
		"40": true, "41": true, "42": true, "43": true, "44": true, "45": true, "46": true, "47": true, "48": true, "49": true,
		"50": true, "51": true, "52": true, "53": true, "54": true, "55": true, "56": true, "57": true, "58": true, "59": true,
		"60": true, "61": true, "62": true, "63": true, "64": true, "65": true, "66": true, "67": true, "68": true,
	},
	"MU": { // Mauritius
		"AG": true,
		"BL": true,
		"CC": true,
		"FL": true,
		"GP": true,
		"MO": true,
		"PA": true,
		"PL": true,
		"PW": true,
		"RO": true,
		"RR": true,
		"SA": true,
	},
	"MV": { // Maldives
		"01":  true,
		"MLE": true,
		"00":  true, "02": true, "03": true, "04": true, "05": true, "07": true, "08": true,
		"12": true, "13": true, "14": true, "17": true,
		"20": true, "23": true, "24": true, "25": true, "26": true, "27": true, "28": true, "29": true,
	},
	"MW": { // Malawi
		"C":  true,
		"N":  true,
		"S":  true,
		"BA": true,
		"BL": true,
		"CK": true,
		"CR": true,
		"CT": true,
		"DE": true,
		"DO": true,
		"KR": true,
		"KS": true,
		"LI": true,
		"LK": true,
		"MC": true,
		"MG": true,
		"MH": true,
		"MU": true,
		"MW": true,
		"MZ": true,
		"NB": true,
		"NE": true,
		"NI": true,
		"NK": true,
		"NS": true,
		"NU": true,
		"PH": true,
		"RU": true,
		"SA": true,
		"TH": true,
		"ZO": true,
	},
	"MX": { // Mexico
		"AGU": true,
		"BCN": true,
		"BCS": true,
		"CAM": true,
		"CHH": true,
		"CHP": true,
		"CMX": true,
		"COA": true,
		"COL": true,
		"DUR": true,
		"GRO": true,
		"GUA": true,
		"HID": true,
		"JAL": true,
		"MEX": true,
		"MIC": true,
		"MOR": true,
		"NAY": true,
		"NLE": true,
		"OAX": true,
		"PUE": true,
		"QUE": true,
		"ROO": true,
		"SIN": true,
		"SLP": true,
		"SON": true,
		"TAB": true,
		"TAM": true,
		"TLA": true,
		"VER": true,
		"YUC": true,
		"ZAC": true,
	},
	"MY": { // Malaysia
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true,
	},
	"MZ": { // Mozambique
		"A":   true,
		"B":   true,
		"G":   true,
		"I":   true,
		"L":   true,
		"MPM": true,
		"N":   true,
		"P":   true,
		"Q":   true,
		"S":   true,
		"T":   true,
	},
	"NA": { // Namibia
		"CA": true,
		"ER": true,
		"HA": true,
		"KA": true,
		"KE": true,
		"KH": true,
		"KU": true,
		"KW": true,
		"OD": true,
		"OH": true,
		"ON": true,
		"OS": true,
		"OT": true,
		"OW": true,
	},
	"NC": {}, // New Caledonia
	"NE": { // Niger
		"1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true, "8": true,
	},
	"NF": {}, // Norfolk Island
	"NG": { // Nigeria
		"AB": true,
		"AD": true,
		"AK": true,
		"AN": true,
		"BA": true,
		"BE": true,
		"BO": true,
		"BY": true,
		"CR": true,
		"DE": true,
		"EB": true,
		"ED": true,
		"EK": true,
		"EN": true,
		"FC": true,
		"GO": true,
		"IM": true,
		"JI": true,
		"KD": true,
		"KE": true,
		"KN": true,
		"KO": true,
		"KT": true,
		"KW": true,
		"LA": true,
		"NA": true,
		"NI": true,
		"OG": true,
		"ON": true,
		"OS": true,
		"OY": true,
		"PL": true,
		"RI": true,
		"SO": true,
		"TA": true,
		"YO": true,
		"ZA": true,
	},
	"NI": { // Nicaragua
		"AN": true,
		"AS": true,
		"BO": true,
		"CA": true,
		"CI": true,
		"CO": true,
		"ES": true,
		"GR": true,
		"JI": true,
		"LE": true,
		"MD": true,
		"MN": true,
		"MS": true,
		"MT": true,
		"NS": true,
		"RI": true,
		"SJ": true,
	},
	"NL": { // Netherlands
		"DR":  true,
		"FL":  true,
		"FR":  true,
		"GE":  true,
		"GR":  true,
		"LI":  true,
		"NB":  true,
		"NH":  true,
		"OV":  true,
		"UT":  true,
		"ZE":  true,
		"ZH":  true,
		"AW":  true,
		"BQ1": true,
		"BQ2": true,
		"BQ3": true,
		"CW":  true,
		"SX":  true,
	},
	"NO": { // Norway
		"03": true,
		"11": true,
		"15": true,
		"18": true,
		"21": true,
		"22": true,
		"30": true,
		"34": true,
		"38": true,
		"42": true,
		"46": true,
		"50": true,
		"54": true,
	},
	"NP": { // Nepal
		"1":  true,
		"2":  true,
		"3":  true,
		"4":  true,
		"5":  true,
		"P1": true,
		"P2": true,
		"P3": true,
		"P4": true,
		"P5": true,
		"P6": true,
		"P7": true,
		"BA": true,
		"BH": true,
		"DH": true,
		"GA": true,
		"JA": true,
		"KA": true,
		"KO": true,
		"LU": true,
		"MA": true,
		"ME": true,
		"NA": true,
		"RA": true,
		"SA": true,
		"SE": true,
	},
	"NR": { // Nauru
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true, "14": true,
	},
	"NU": {}, // Niue
	"NZ": { // New Zealand
		"AUK": true,
		"BOP": true,
		"CAN": true,
		"CIT": true,
		"GIS": true,
		"HKB": true,
		"MBH": true,
		"MWT": true,
		"NSN": true,
		"NTL": true,
		"OTA": true,
		"STL": true,
		"TAS": true,
		"TKI": true,
		"WGN": true,
		"WKO": true,
		"WTC": true,
	},
	"OM": { // Oman
		"BJ": true,
		"BS": true,
		"BU": true,
		"DA": true,
		"MA": true,
		"MU": true,
		"SJ": true,
		"SS": true,
		"WU": true,
		"ZA": true,
		"ZU": true,
	},
	"PA": { // Panama
		"1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true, "8": true, "9": true, "10": true,
		"EM": true,
		"KY": true,
		"NB": true,
		"NT": true,
	},
	"PE": { // Peru
		"AMA": true,
		"ANC": true,
		"APU": true,
		"ARE": true,
		"AYA": true,
		"CAJ": true,
		"CAL": true,
		"CUS": true,
		"HUC": true,
		"HUV": true,
		"ICA": true,
		"JUN": true,
		"LAL": true,
		"LAM": true,
		"LIM": true,
		"LMA": true,
		"LOR": true,
		"MDD": true,
		"MOQ": true,
		"PAS": true,
		"PIU": true,
		"PUN": true,
		"SAM": true,
		"TAC": true,
		"TUM": true,
		"UCA": true,
	},
	"PF": {}, // French Polynesia
	"PG": { // Papua New Guinea
		"CPK": true,
		"CPM": true,
		"EBR": true,
		"EHG": true,
		"EPW": true,
		"ESW": true,
		"GPK": true,
		"HLA": true,
		"JWK": true,
		"MBA": true,
		"MPL": true,
		"MPM": true,
		"MRL": true,
		"NCD": true,
		"NIK": true,
		"NPP": true,
		"NSB": true,
		"SAN": true,
		"SHM": true,
		"WBK": true,
		"WHM": true,
		"WPD": true,
	},
	"PH": { // Philippines
		"00": true, "01": true, "02": true, "03": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true,
		"40": true, "41": true,
		"ABR": true,
		"AGN": true,
		"AGS": true,
		"AKL": true,
		"ALB": true,
		"ANT": true,
		"APA": true,
		"AUR": true,
		"BAN": true,
		"BAS": true,
		"BEN": true,
		"BIL": true,
		"BOH": true,
		"BTG": true,
		"BTN": true,
		"BUK": true,
		"BUL": true,
		"CAG": true,
		"CAM": true,
		"CAN": true,
		"CAP": true,
		"CAS": true,
		"CAT": true,
		"CAV": true,
		"CEB": true,
		"COM": true,
		"DAO": true,
		"DAS": true,
		"DAV": true,
		"DIN": true,
		"DVO": true,
		"EAS": true,
		"GUI": true,
		"IFU": true,
		"ILI": true,
		"ILN": true,
		"ILS": true,
		"ISA": true,
		"KAL": true,
		"LAG": true,
		"LAN": true,
		"LAS": true,
		"LEY": true,
		"LUN": true,
		"MAD": true,
		"MAG": true,
		"MAS": true,
		"MDC": true,
		"MDR": true,
		"MOU": true,
		"MSC": true,
		"MSR": true,
		"NCO": true,
		"NEC": true,
		"NER": true,
		"NSA": true,
		"NUE": true,
		"NUV": true,
		"PAM": true,
		"PAN": true,
		"PLW": true,
		"QUE": true,
		"QUI": true,
		"RIZ": true,
		"ROM": true,
		"SAR": true,
		"SCO": true,
		"SIG": true,
		"SLE": true,
		"SLU": true,
		"SOR": true,
		"SUK": true,
		"SUN": true,
		"SUR": true,
		"TAR": true,
		"TAW": true,
		"WSA": true,
		"ZAN": true,
		"ZAS": true,
		"ZMB": true,
		"ZSI": true,
	},
	"PK": { // Pakistan
		"BA": true,
		"GB": true,
		"IS": true,
		"JK": true,
		"KP": true,
		"PB": true,
		"SD": true,
	},
	"PL": { // Poland
		"02": true,
		"04": true,
		"06": true,
		"08": true,
		"10": true,
		"12": true,
		"14": true,
		"16": true,
		"18": true,
		"20": true,
		"22": true,
		"24": true,
		"26": true,
		"28": true,
		"30": true,
		"32": true,
	},
	"PM": {}, // Saint Pierre and Miquelon
	"PN": {}, // Pitcairn
	"PR": {}, // Puerto Rico
	"PS": { // Palestine, State of
		"BTH": true,
		"DEB": true,
		"GZA": true,
		"HBN": true,
		"JEM": true,
		"JEN": true,
		"JRH": true,
		"KYS": true,
		"NBS": true,
		"NGZ": true,
		"QQA": true,
		"RBH": true,
		"RFH": true,
		"SLT": true,
		"TBS": true,
		"TKM": true,
	},
	"PT": { // Portugal
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true,
		"20": true,
		"30": true,
	},
	"PW": { // Palau
		"002": true,
		"004": true,
		"010": true,
		"050": true,
		"100": true,
		"150": true,
		"212": true,
		"214": true,
		"218": true,
		"222": true,
		"224": true,
		"226": true,
		"227": true,
		"228": true,
		"350": true,
		"370": true,
	},
	"PY": { // Paraguay
		"1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true, "8": true, "9": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "19": true,
		"ASU": true,
	},
	"QA": { // Qatar
		"DA": true,
		"KH": true,
		"MS": true,
		"RA": true,
		"SH": true,
		"US": true,
		"WA": true,
		"ZA": true,
	},
	"RE": {}, // Réunion
	"RO": { // Romania
		"AB": true,
		"AG": true,
		"AR": true,
		"B":  true,
		"BC": true,
		"BH": true,
		"BN": true,
		"BR": true,
		"BT": true,
		"BV": true,
		"BZ": true,
		"CJ": true,
		"CL": true,
		"CS": true,
		"CT": true,
		"CV": true,
		"DB": true,
		"DJ": true,
		"GJ": true,
		"GL": true,
		"GR": true,
		"HD": true,
		"HR": true,
		"IF": true,
		"IL": true,
		"IS": true,
		"MH": true,
		"MM": true,
		"MS": true,
		"NT": true,
		"OT": true,
		"PH": true,
		"SB": true,
		"SJ": true,
		"SM": true,
		"SV": true,
		"TL": true,
		"TM": true,
		"TR": true,
		"VL": true,
		"VN": true,
		"VS": true,
	},
	"RS": { // Serbia
		"KM": true,
		"VO": true,
		"00": true,
		"01": true,
		"02": true,
		"03": true,
		"04": true,
		"05": true,
		"06": true,
		"07": true,
		"08": true,
		"09": true,
		"10": true,
		"11": true,
		"12": true,
		"13": true,
		"14": true,
		"15": true,
		"16": true,
		"17": true,
		"18": true,
		"19": true,
		"20": true,
		"21": true,
		"22": true,
		"23": true,
		"24": true,
		"25": true,
		"26": true,
		"27": true,
		"28": true,
		"29": true,
	},
	"RU": { // Russian Federation
		"AD":  true,
		"AL":  true,
		"ALT": true,
		"AMU": true,
		"ARK": true,
		"AST": true,
		"BA":  true,
		"BEL": true,
		"BRY": true,
		"BU":  true,
		"CE":  true,
		"CHE": true,
		"CHU": true,
		"CU":  true,
		"DA":  true,
		"IN":  true,
		"IRK": true,
		"IVA": true,
		"KAM": true,
		"KB":  true,
		"KC":  true,
		"KDA": true,
		"KEM": true,
		"KGD": true,
		"KGN": true,
		"KHA": true,
		"KHM": true,
		"KIR": true,
		"KK":  true,
		"KL":  true,
		"KLU": true,
		"KO":  true,
		"KOS": true,
		"KR":  true,
		"KRS": true,
		"KYA": true,
		"LEN": true,
		"LIP": true,
		"MAG": true,
		"ME":  true,
		"MO":  true,
		"MOS": true,
		"MOW": true,
		"MUR": true,
		"NEN": true,
		"NGR": true,
		"NIZ": true,
		"NVS": true,
		"OMS": true,
		"ORE": true,
		"ORL": true,
		"PER": true,
		"PNZ": true,
		"PRI": true,
		"PSK": true,
		"ROS": true,
		"RYA": true,
		"SA":  true,
		"SAK": true,
		"SAM": true,
		"SAR": true,
		"SE":  true,
		"SMO": true,
		"SPE": true,
		"STA": true,
		"SVE": true,
		"TA":  true,
		"TAM": true,
		"TOM": true,
		"TUL": true,
		"TVE": true,
		"TY":  true,
		"TYU": true,
		"UD":  true,
		"ULY": true,
		"VGG": true,
		"VLA": true,
		"VLG": true,
		"VOR": true,
		"YAN": true,
		"YAR": true,
		"YEV": true,
		"ZAB": true,
	},
	"RW": { // Rwanda
		"01": true, "02": true, "03": true, "04": true, "05": true,
	},
	"SA": { // Saudi Arabia
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "14": true,
	},
	"SB": { // Solomon Islands
		"CE": true,
		"CH": true,
		"CT": true,
		"GU": true,
		"IS": true,
		"MK": true,
		"ML": true,
		"RB": true,
		"TE": true,
		"WE": true,
	},
	"SC": { // Seychelles
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true,
		"20": true, "21": true, "22": true, "23": true, "24": true, "25": true, "26": true, "27": true,
	},
	"SD": { // Sudan
		"DC": true,
		"DE": true,
		"DN": true,
		"DS": true,
		"DW": true,
		"GD": true,
		"GK": true,
		"GZ": true,
		"KA": true,
		"KH": true,
		"KN": true,
		"KS": true,
		"NB": true,
		"NO": true,
		"NR": true,
		"NW": true,
		"RS": true,
		"SI": true,
	},
	"SE": { // Sweden
		"AB": true,
		"AC": true,
		"BD": true,
		"C":  true,
		"D":  true,
		"E":  true,
		"F":  true,
		"G":  true,
		"H":  true,
		"I":  true,
		"K":  true,
		"M":  true,
		"N":  true,
		"O":  true,
		"S":  true,
		"T":  true,
		"U":  true,
		"W":  true,
		"X":  true,
		"Y":  true,
		"Z":  true,
	},
	"SG": { // Singapore
		"01": true, "02": true, "03": true, "04": true, "05": true,
	},
	"SH": { // Saint Helena, Ascension and Tristan da Cunha
		"AC": true,
		"HL": true,
		"TA": true,
	},
	"SI": { // Slovenia
		"001": true,
		"002": true,
		"003": true,
		"004": true,
		"005": true,
		"006": true,
		"007": true,
		"008": true,
		"009": true,
		"010": true,
		"011": true,
		"012": true,
		"013": true,
		"014": true,
		"015": true,
		"016": true,
		"017": true,
		"018": true,
		"019": true,
		"020": true,
		"021": true,
		"022": true,
		"023": true,
		"024": true,
		"025": true,
		"026": true,
		"027": true,
		"028": true,
		"029": true,
		"030": true,
		"031": true,
		"032": true,
		"033": true,
		"034": true,
		"035": true,
		"036": true,
		"037": true,
		"038": true,
		"039": true,
		"040": true,
		"041": true,
		"042": true,
		"043": true,
		"044": true,
		"045": true,
		"046": true,
		"047": true,
		"048": true,
		"049": true,
		"050": true,
		"051": true,
		"052": true,
		"053": true,
		"054": true,
		"055": true,
		"056": true,
		"057": true,
		"058": true,
		"059": true,
		"060": true,
		"061": true,
		"062": true,
		"063": true,
		"064": true,
		"065": true,
		"066": true,
		"067": true,
		"068": true,
		"069": true,
		"070": true,
		"071": true,
		"072": true,
		"073": true,
		"074": true,
		"075": true,
		"076": true,
		"077": true,
		"078": true,
		"079": true,
		"080": true,
		"081": true,
		"082": true,
		"083": true,
		"084": true,
		"085": true,
		"086": true,
		"087": true,
		"088": true,
		"089": true,
		"090": true,
		"091": true,
		"092": true,
		"093": true,
		"094": true,
		"095": true,
		"096": true,
		"097": true,
		"098": true,
		"099": true,
		"100": true,
		"101": true,
		"102": true,
		"103": true,
		"104": true,
		"105": true,
		"106": true,
		"107": true,
		"108": true,
		"109": true,
		"110": true,
		"111": true,
		"112": true,
		"113": true,
		"114": true,
		"115": true,
		"116": true,
		"117": true,
		"118": true,
		"119": true,
		"120": true,
		"121": true,
		"122": true,
		"123": true,
		"124": true,
		"125": true,
		"126": true,
		"127": true,
		"128": true,
		"129": true,
		"130": true,
		"131": true,
		"132": true,
		"133": true,
		"134": true,
		"135": true,
		"136": true,
		"137": true,
		"138": true,
		"139": true,
		"140": true,
		"141": true,
		"142": true,
		"143": true,
		"144": true,
		"146": true,
		"147": true,
		"148": true,
		"149": true,
		"150": true,
		"151": true,
		"152": true,
		"153": true,
		"154": true,
		"155": true,
		"156": true,
		"157": true,
		"158": true,
		"159": true,
		"160": true,
		"161": true,
		"162": true,
		"163": true,
		"164": true,
		"165": true,
		"166": true,
		"167": true,
		"168": true,
		"169": true,
		"170": true,
		"171": true,
		"172": true,
		"173": true,
		"174": true,
		"175": true,
		"176": true,
		"177": true,
		"178": true,
		"179": true,
		"180": true,
		"181": true,
		"182": true,
		"183": true,
		"184": true,
		"185": true,
		"186": true,
		"187": true,
		"188": true,
		"189": true,
		"190": true,
		"191": true,
		"192": true,
		"193": true,
		"194": true,
		"195": true,
		"196": true,
		"197": true,
		"198": true,
		"199": true,
		"200": true,
		"201": true,
		"202": true,
		"203": true,
		"204": true,
		"205": true,
		"206": true,
		"207": true,
		"208": true,
		"209": true,
		"210": true,
		"211": true,
		"212": true,
		"213": true,
	},
	"SJ": {}, // Svalbard and Jan Mayen
	"SK": { // Slovakia
		"BC": true,
		"BL": true,
		"KI": true,
		"NI": true,
		"PV": true,
		"TA": true,
		"TC": true,
		"ZI": true,
	},
	"SL": { // Sierra Leone
		"E":  true,
		"N":  true,
		"NW": true,
		"S":  true,
		"W":  true,
	},
	"SM": { // San Marino
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
	},
	"SN": { // Senegal
		"DB": true,
		"DK": true,
		"FK": true,
		"KA": true,
		"KD": true,
		"KE": true,
		"KL": true,
		"LG": true,
		"MT": true,
		"SE": true,
		"SL": true,
		"TC": true,
		"TH": true,
		"ZG": true,
	},
	"SO": { // Somalia
		"AW": true,
		"BK": true,
		"BN": true,
		"BR": true,
		"BY": true,
		"GA": true,
		"GE": true,
		"HI": true,
		"JD": true,
		"JH": true,
		"MU": true,
		"NU": true,
		"SA": true,
		"SD": true,
		"SH": true,
		"SO": true,
		"TO": true,
		"WO": true,
	},
	"SR": { // Suriname
		"BR": true,
		"CM": true,
		"CR": true,
		"MA": true,
		"NI": true,
		"PM": true,
		"PR": true,
		"SA": true,
		"SI": true,
		"WA": true,
	},
	"SS": { // South Sudan
		"BN": true,
		"BW": true,
		"EC": true,
		"EE": true,
		"EW": true,
		"JG": true,
		"LK": true,
		"NU": true,
		"UY": true,
		"WR": true,
	},
	"ST": { // Sao Tome and Principe
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true,
		"P": true,
	},
	"SV": { // El Salvador
		"AH": true,
		"CA": true,
		"CH": true,
		"CU": true,
		"LI": true,
		"MO": true,
		"PA": true,
		"SA": true,
		"SM": true,
		"SO": true,
		"SS": true,
		"SV": true,
		"UN": true,
		"US": true,
	},
	"SX": {}, // Sint Maarten (Dutch part)
	"SY": { // Syrian Arab Republic
		"DI": true,
		"DR": true,
		"DY": true,
		"HA": true,
		"HI": true,
		"HL": true,
		"HM": true,
		"ID": true,
		"LA": true,
		"QU": true,
		"RA": true,
		"RD": true,
		"SU": true,
		"TA": true,
	},
	"SZ": { // Eswatini
		"HH": true,
		"LU": true,
		"MA": true,
		"SH": true,
	},
	"TC": {}, // Turks and Caicos Islands
	"TD": { // Chad
		"BA": true,
		"BG": true,
		"BO": true,
		"CB": true,
		"EE": true,
		"EO": true,
		"GR": true,
		"HL": true,
		"KA": true,
		"LC": true,
		"LO": true,
		"LR": true,
		"MA": true,
		"MC": true,
		"ME": true,
		"MO": true,
		"ND": true,
		"OD": true,
		"SA": true,
		"SI": true,
		"TA": true,
		"TI": true,
		"WF": true,
	},
	"TF": {}, // French Southern Territories
	"TG": { // Togo
		"C": true,
		"K": true,
		"M": true,
		"P": true,
		"S": true,
	},
	"TH": { // Thailand
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true,
		"20": true, "21": true, "22": true, "23": true, "24": true, "25": true, "26": true, "27": true,
		"30": true, "31": true, "32": true, "33": true, "34": true, "35": true, "36": true, "37": true, "38": true, "39": true,
		"40": true, "41": true, "42": true, "43": true, "44": true, "45": true, "46": true, "47": true, "48": true, "49": true,
		"50": true, "51": true, "52": true, "53": true, "54": true, "55": true, "56": true, "57": true, "58": true,
		"60": true, "61": true, "62": true, "63": true, "64": true, "65": true, "66": true, "67": true,
		"70": true, "71": true, "72": true, "73": true, "74": true, "75": true, "76": true, "77": true,
		"80": true, "81": true, "82": true, "83": true, "84": true, "85": true, "86": true,
		"90": true, "91": true, "92": true, "93": true, "94": true, "95": true, "96": true,
		"S": true,
	},
	"TJ": { // Tajikistan
		"DU": true,
		"GB": true,
		"KT": true,
		"RA": true,
		"SU": true,
	},
	"TK": {}, // Tokelau
	"TL": { // Timor-Leste
		"AL": true,
		"AN": true,
		"BA": true,
		"BO": true,
		"CO": true,
		"DI": true,
		"ER": true,
		"LA": true,
		"LI": true,
		"MF": true,
		"MT": true,
		"OE": true,
		"VI": true,
	},
	"TM": { // Turkmenistan
		"A": true, "B": true, "D": true, "L": true, "M": true, "S": true,
	},
	"TN": { // Tunisia
		"11": true, "12": true, "13": true, "14": true,
		"21": true, "22": true, "23": true,
		"31": true, "32": true, "33": true, "34": true,
		"41": true, "42": true, "43": true,
		"51": true, "52": true, "53": true,
		"61": true,
		"71": true, "72": true, "73": true,
		"81": true, "82": true, "83": true,
	},
	"TO": { // Tonga
		"01": true, "02": true, "03": true, "04": true, "05": true,
	},
	"TR": { // Türkiye
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true,
		"20": true, "21": true, "22": true, "23": true, "24": true, "25": true, "26": true, "27": true, "28": true, "29": true,
		"30": true, "31": true, "32": true, "33": true, "34": true, "35": true, "36": true, "37": true, "38": true, "39": true,
		"40": true, "41": true, "42": true, "43": true, "44": true, "45": true, "46": true, "47": true, "48": true, "49": true,
		"50": true, "51": true, "52": true, "53": true, "54": true, "55": true, "56": true, "57": true, "58": true, "59": true,
		"60": true, "61": true, "62": true, "63": true, "64": true, "65": true, "66": true, "67": true, "68": true, "69": true,
		"70": true, "71": true, "72": true, "73": true, "74": true, "75": true, "76": true, "77": true, "78": true, "79": true,
		"80": true, "81": true,
	},
	"TT": { // Trinidad and Tobago
		"ARI": true,
		"CHA": true,
		"CTT": true,
		"DMN": true,
		"MRC": true,
		"PED": true,
		"POS": true,
		"PRT": true,
		"PTF": true,
		"SFO": true,
		"SGE": true,
		"SIP": true,
		"SJL": true,
		"TOB": true,
		"TUP": true,
	},
	"TV": { // Tuvalu
		"FUN": true,
		"NIT": true,
		"NKF": true,
		"NKL": true,
		"NMA": true,
		"NMG": true,
		"NUI": true,
		"VAI": true,
	},
	"TW": { // Taiwan, Province of China
		"CHA": true,
		"CYI": true,
		"CYQ": true,
		"HSQ": true,
		"HSZ": true,
		"HUA": true,
		"ILA": true,
		"KEE": true,
		"KHH": true,
		"KIN": true,
		"LIE": true,
		"MIA": true,
		"NAN": true,
		"NWT": true,
		"PEN": true,
		"PIF": true,
		"TAO": true,
		"TNN": true,
		"TPE": true,
		"TTT": true,
		"TXG": true,
		"YUN": true,
	},
	"TZ": { // Tanzania, United Republic of
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true,
		"10": true, "11": true, "12": true, "13": true, "14": true, "15": true, "16": true, "17": true, "18": true, "19": true,
		"20": true, "21": true, "22": true, "23": true, "24": true, "25": true, "26": true, "27": true, "28": true, "29": true,
		"30": true, "31": true,
	},
	"UA": { // Ukraine
		"05": true,
		"07": true,
		"09": true,
		"12": true,
		"14": true,
		"18": true,
		"21": true,
		"23": true,
		"26": true,
		"30": true,
		"32": true,
		"35": true,
		"40": true,
		"43": true,
		"46": true,
		"48": true,
		"51": true,
		"53": true,
		"56": true,
		"59": true,
		"61": true,
		"63": true,
		"65": true,
		"68": true,
		"71": true,
		"74": true,
		"77": true,
	},
	"UG": { // Uganda
		"C":   true,
		"E":   true,
		"N":   true,
		"W":   true,
		"101": true,
		"102": true,
		"103": true,
		"104": true,
		"105": true,
		"106": true,
		"107": true,
		"108": true,
		"109": true,
		"110": true,
		"111": true,
		"112": true,
		"113": true,
		"114": true,
		"115": true,
		"116": true,
		"117": true,
		"118": true,
		"119": true,
		"120": true,
		"121": true,
		"122": true,
		"123": true,
		"124": true,
		"125": true,
		"126": true,
		"201": true,
		"202": true,
		"203": true,
		"204": true,
		"205": true,
		"206": true,
		"207": true,
		"208": true,
		"209": true,
		"210": true,
		"211": true,
		"212": true,
		"213": true,
		"214": true,
		"215": true,
		"216": true,
		"217": true,
		"218": true,
		"219": true,
		"220": true,
		"221": true,
		"222": true,
		"223": true,
		"224": true,
		"225": true,
		"226": true,
		"227": true,
		"228": true,
		"229": true,
		"230": true,
		"231": true,
		"232": true,
		"233": true,
		"234": true,
		"235": true,
		"236": true,
		"237": true,
		"301": true,
		"302": true,
		"303": true,
		"304": true,
		"305": true,
		"306": true,
		"307": true,
		"308": true,
		"309": true,
		"310": true,
		"311": true,
		"312": true,
		"313": true,
		"314": true,
		"315": true,
		"316": true,
		"317": true,
		"318": true,
		"319": true,
		"320": true,
		"321": true,
		"322": true,
		"323": true,
		"324": true,
		"325": true,
		"326": true,
		"327": true,
		"328": true,
		"329": true,
		"330": true,
		"331": true,
		"332": true,
		"333": true,
		"334": true,
		"335": true,
		"336": true,
		"337": true,
		"401": true,
		"402": true,
		"403": true,
		"404": true,
		"405": true,
		"406": true,
		"407": true,
		"408": true,
		"409": true,
		"410": true,
		"411": true,
		"412": true,
		"413": true,
		"414": true,
		"415": true,
		"416": true,
		"417": true,
		"418": true,
		"419": true,
		"420": true,
		"421": true,
		"422": true,
		"423": true,
		"424": true,
		"425": true,
		"426": true,
		"427": true,
		"428": true,
		"429": true,
		"430": true,
		"431": true,
		"432": true,
		"433": true,
		"434": true,
		"435": true,
	},
	"UM": { // United States Minor Outlying Islands
		"67": true,
		"71": true,
		"76": true,
		"79": true,
		"81": true,
		"84": true,
		"86": true,
		"89": true,
		"95": true,
	},
	"US": { // United States of America
		"AK": true,
		"AL": true,
		"AR": true,
		"AS": true,
		"AZ": true,
		"CA": true,
		"CO": true,
		"CT": true,
		"DC": true,
		"DE": true,
		"FL": true,
		"GA": true,
		"GU": true,
		"HI": true,
		"IA": true,
		"ID": true,
		"IL": true,
		"IN": true,
		"KS": true,
		"KY": true,
		"LA": true,
		"MA": true,
		"MD": true,
		"ME": true,
		"MI": true,
		"MN": true,
		"MO": true,
		"MP": true,
		"MS": true,
		"MT": true,
		"NC": true,
		"ND": true,
		"NE": true,
		"NH": true,
		"NJ": true,
		"NM": true,
		"NV": true,
		"NY": true,
		"OH": true,
		"OK": true,
		"OR": true,
		"PA": true,
		"PR": true,
		"RI": true,
		"SC": true,
		"SD": true,
		"TN": true,
		"TX": true,
		"UM": true,
		"UT": true,
		"VA": true,
		"VI": true,
		"VT": true,
		"WA": true,
		"WI": true,
		"WV": true,
		"WY": true,
	},
	"UY": { // Uruguay
		"AR": true,
		"CA": true,
		"CL": true,
		"CO": true,
		"DU": true,
		"FD": true,
		"FS": true,
		"LA": true,
		"MA": true,
		"MO": true,
		"PA": true,
		"RN": true,
		"RO": true,
		"RV": true,
		"SA": true,
		"SJ": true,
		"SO": true,
		"TA": true,
		"TT": true,
	},
	"UZ": { // Uzbekistan
		"AN": true,
		"BU": true,
		"FA": true,
		"JI": true,
		"NG": true,
		"NW": true,
		"QA": true,
		"QR": true,
		"SA": true,
		"SI": true,
		"SU": true,
		"TK": true,
		"TO": true,
		"XO": true,
	},
	"VA": {}, // Holy City
	"VC": { // Saint Vincent and the Grenadines
		"01": true,
		"02": true,
		"03": true,
		"04": true,
		"05": true,
		"06": true,
	},
	"VE": { // Venezuela (Bolivarian Republic of)
		"A": true,
		"B": true,
		"C": true,
		"D": true,
		"E": true,
		"F": true,
		"G": true,
		"H": true,
		"I": true,
		"J": true,
		"K": true,
		"L": true,
		"M": true,
		"N": true,
		"O": true,
		"P": true,
		"R": true,
		"S": true,
		"T": true,
		"U": true,
		"V": true,
		"W": true,
		"X": true,
		"Y": true,
		"Z": true,
	},
	"VG": {}, // Virgin Islands (British)
	"VI": {}, // Virgin Islands (U.S.)
	"VN": { // Viet Nam
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "09": true,
		"13": true, "14": true, "18": true,
		"20": true, "21": true, "22": true, "23": true, "24": true, "25": true, "26": true, "27": true, "28": true, "29": true,
		"30": true, "31": true, "32": true, "33": true, "34": true, "35": true, "36": true, "37": true, "39": true,
		"40": true, "41": true, "43": true, "44": true, "45": true, "46": true, "47": true, "49": true,
		"50": true, "51": true, "52": true, "53": true, "54": true, "55": true, "56": true, "57": true, "58": true, "59": true,
		"61": true, "63": true, "66": true, "67": true, "68": true, "69": true,
		"70": true, "71": true, "72": true, "73": true,
		"CT": true,
		"DN": true,
		"HN": true,
		"HP": true,
		"SG": true,
	},
	"VU": { // Vanuatu
		"MAP": true,
		"PAM": true,
		"SAM": true,
		"SEE": true,
		"TAE": true,
		"TOB": true,
	},
	"WF": { // Wallis and Futuna
		"AL": true,
		"SG": true,
		"UV": true,
	},
	"WS": { // Samoa
		"AA": true,
		"AL": true,
		"AT": true,
		"FA": true,
		"GE": true,
		"GI": true,
		"PA": true,
		"SA": true,
		"TU": true,
		"VF": true,
		"VS": true,
	},
	"YE": { // Yemen
		"AB": true,
		"AD": true,
		"AM": true,
		"BA": true,
		"DA": true,
		"DH": true,
		"HD": true,
		"HJ": true,
		"HU": true,
		"IB": true,
		"JA": true,
		"LA": true,
		"MA": true,
		"MR": true,
		"MW": true,
		"RA": true,
		"SA": true,
		"SD": true,
		"SH": true,
		"SN": true,
		"SU": true,
		"TA": true,
	},
	"YT": {}, // Mayotte
	"ZA": { // South Africa
		"EC":  true,
		"FS":  true,
		"GP":  true,
		"KZN": true,
		"LP":  true,
		"MP":  true,
		"NC":  true,
		"NW":  true,
		"WC":  true,
	},
	"ZM": { // Zambia
		"01": true, "02": true, "03": true, "04": true, "05": true, "06": true, "07": true, "08": true, "09": true, "10": true,
	},
	"ZW": { // Zimbabwe
		"BU": true,
		"HA": true,
		"MA": true,
		"MC": true,
		"ME": true,
		"MI": true,
		"MN": true,
		"MS": true,
		"MV": true,
		"MW": true,
	},
}

var iSO3166_1_NumericCodes = map[string]bool{
	"004": true, // Afghanistan
	"008": true, // Albania
	"010": true, // Antarctica
	"012": true, // Algeria
	"016": true, // American Samoa
	"020": true, // Andorra
	"024": true, // Angola
	"028": true, // Antigua and Barbuda
	"031": true, // Azerbaijan	Until 1991 part of the USSR
	"032": true, // Argentina
	"036": true, // Australia
	"040": true, // Austria
	"044": true, // Bahamas
	"048": true, // Bahrain
	"050": true, // Bangladesh
	"051": true, // Armenia	Until 1991 part of the USSR
	"052": true, // Barbados
	"056": true, // Belgium
	"060": true, // Bermuda
	"064": true, // Bhutan
	"068": true, // Bolivia (Plurinational State of)
	"070": true, // Bosnia and Herzegovina	Until 1992 part of Yugoslavia
	"072": true, // Botswana
	"074": true, // Bouvet Island
	"076": true, // Brazil
	"084": true, // Belize	Formerly British Honduras
	"086": true, // British Indian Ocean Territory
	"090": true, // Solomon Islands	Formerly British Solomon Islands
	"092": true, // Virgin Islands (British)
	"096": true, // Brunei Darussalam
	"100": true, // Bulgaria
	"104": true, // Myanmar	Formerly Burma
	"108": true, // Burundi
	"112": true, // Belarus	Formerly Byelorussian SSR
	"116": true, // Cambodia
	"120": true, // Cameroon
	"124": true, // Canada
	"132": true, // Cabo Verde	Formerly Cape Verde
	"136": true, // Cayman Islands
	"140": true, // Central African Republic
	"144": true, // Sri Lanka	Formerly Ceylon
	"148": true, // Chad
	"152": true, // Chile
	"156": true, // China
	"158": true, // Taiwan, Province of China
	"162": true, // Christmas Island
	"166": true, // Cocos (Keeling) Islands
	"170": true, // Colombia
	"174": true, // Comoros
	"175": true, // Mayotte	Until 1975 part of Comoros, own ISO code since 1993
	"178": true, // Congo
	"180": true, // Congo, Democratic Republic of the
	"184": true, // Cook Islands
	"188": true, // Costa Rica
	"191": true, // Croatia	Until 1992 part of Yugoslavia
	"192": true, // Cuba
	"196": true, // Cyprus
	"203": true, // Czechia	Until 1993 part of Czechoslovakia
	"204": true, // Benin	Formerly Dahomey
	"208": true, // Denmark
	"212": true, // Dominica
	"214": true, // Dominican Republic
	"218": true, // Ecuador
	"222": true, // El Salvador
	"226": true, // Equatorial Guinea
	"231": true, // Ethiopia
	"232": true, // Eritrea	Until 1993 part of Ethiopia
	"233": true, // Estonia	Until 1991 part of the USSR
	"234": true, // Faroe Islands	Previously spelled Faeroe Islands
	"238": true, // Falkland Islands (Malvinas)
	"239": true, // South Georgia and the South Sandwich Islands	Until 1993 part of the Falkland Islands
	"242": true, // Fiji
	"246": true, // Finland
	"248": true, // Åland Islands	Until 2004 part of Finland
	"250": true, // France
	"254": true, // French Guiana
	"258": true, // French Polynesia
	"260": true, // French Southern Territories
	"262": true, // Djibouti	Formerly French Territory of the Afars and the Issas
	"266": true, // Gabon
	"268": true, // Georgia	Until 1991 part of the USSR
	"270": true, // Gambia
	"275": true, // Palestine, State of	Replaced the Gaza Strip, which was assigned code 274 by the United Nations Statistics Division
	"276": true, // Germany	Unified country since 1990
	"288": true, // Ghana
	"292": true, // Gibraltar
	"296": true, // Kiribati	Formerly Gilbert and Ellice Islands
	"300": true, // Greece
	"304": true, // Greenland
	"308": true, // Grenada
	"312": true, // Guadeloupe
	"316": true, // Guam
	"320": true, // Guatemala
	"324": true, // Guinea
	"328": true, // Guyana
	"332": true, // Haiti
	"334": true, // Heard Island and McDonald Islands
	"336": true, // Holy See
	"340": true, // Honduras
	"344": true, // Hong Kong
	"348": true, // Hungary
	"352": true, // Iceland
	"356": true, // India
	"360": true, // Indonesia
	"364": true, // Iran (Islamic Republic of)
	"368": true, // Iraq
	"372": true, // Ireland
	"376": true, // Israel
	"380": true, // Italy
	"384": true, // Côte d'Ivoire	Formerly Ivory Coast
	"388": true, // Jamaica
	"392": true, // Japan
	"398": true, // Kazakhstan	Until 1991 part of the USSR
	"400": true, // Jordan
	"404": true, // Kenya
	"408": true, // Korea (Democratic People's Republic of)
	"410": true, // Korea, Republic of
	"414": true, // Kuwait
	"417": true, // Kyrgyzstan	Until 1991 part of the USSR
	"418": true, // Lao People's Democratic Republic
	"422": true, // Lebanon
	"426": true, // Lesotho
	"428": true, // Latvia	Until 1991 part of the USSR
	"430": true, // Liberia
	"434": true, // Libya
	"438": true, // Liechtenstein
	"440": true, // Lithuania	Until 1991 part of the USSR
	"442": true, // Luxembourg
	"446": true, // Macao
	"450": true, // Madagascar
	"454": true, // Malawi
	"458": true, // Malaysia
	"462": true, // Maldives
	"466": true, // Mali
	"470": true, // Malta
	"474": true, // Martinique
	"478": true, // Mauritania
	"480": true, // Mauritius
	"484": true, // Mexico
	"492": true, // Monaco
	"496": true, // Mongolia
	"498": true, // Moldova, Republic of	Until 1991 part of the USSR
	"499": true, // Montenegro	Until 2006 part of Yugoslavia/Serbia and Montenegro
	"500": true, // Montserrat
	"504": true, // Morocco
	"508": true, // Mozambique
	"512": true, // Oman	Formerly Muscat and Oman
	"516": true, // Namibia
	"520": true, // Nauru
	"524": true, // Nepal
	"528": true, // Netherlands
	"531": true, // Curaçao	Until 2010 part of the Netherlands Antilles
	"533": true, // Aruba	Until 1986 part of the Netherlands Antilles
	"534": true, // Sint Maarten (Dutch part)	Until 2010 part of the Netherlands Antilles
	"535": true, // Bonaire, Sint Eustatius and Saba	Until 2010 part of the Netherlands Antilles
	"540": true, // New Caledonia
	"548": true, // Vanuatu	Formerly New Hebrides
	"554": true, // New Zealand
	"558": true, // Nicaragua
	"562": true, // Niger
	"566": true, // Nigeria
	"570": true, // Niue
	"574": true, // Norfolk Island
	"578": true, // Norway
	"580": true, // Northern Mariana Islands	Until 1986 part of Pacific Islands (Trust Territory)
	"581": true, // United States Minor Outlying Islands	Merger of uninhabited U.S. islands on the Pacific Ocean in 1986
	"583": true, // Micronesia (Federated States of)	Until 1986 part of Pacific Islands (Trust Territory)
	"584": true, // Marshall Islands	Until 1986 part of Pacific Islands (Trust Territory)
	"585": true, // Palau	Until 1986 part of Pacific Islands (Trust Territory)
	"586": true, // Pakistan
	"591": true, // Panama
	"598": true, // Papua New Guinea
	"600": true, // Paraguay
	"604": true, // Peru
	"608": true, // Philippines
	"612": true, // Pitcairn
	"616": true, // Poland
	"620": true, // Portugal
	"624": true, // Guinea-Bissau	Formerly Portuguese Guinea
	"626": true, // Timor-Leste	Formerly Portuguese Timor and East Timor
	"630": true, // Puerto Rico
	"634": true, // Qatar
	"638": true, // Réunion
	"642": true, // Romania
	"643": true, // Russian Federation	Until 1991 part of the USSR
	"646": true, // Rwanda
	"652": true, // Saint Barthélemy	Until 2007 part of Guadeloupe
	"654": true, // Saint Helena, Ascension and Tristan da Cunha
	"659": true, // Saint Kitts and Nevis	Until 1985 part of Saint Kitts-Nevis-Anguilla
	"660": true, // Anguilla	Until 1985 part of Saint Kitts-Nevis-Anguilla
	"662": true, // Saint Lucia
	"663": true, // Saint Martin (French part)	Until 2007 part of Guadeloupe
	"666": true, // Saint Pierre and Miquelon
	"670": true, // Saint Vincent and the Grenadines
	"674": true, // San Marino
	"678": true, // Sao Tome and Principe
	"682": true, // Saudi Arabia
	"686": true, // Senegal
	"688": true, // Serbia	Until 2006 part of Yugoslavia/Serbia and Montenegro
	"690": true, // Seychelles
	"694": true, // Sierra Leone
	"702": true, // Singapore
	"703": true, // Slovakia	Until 1993 part of Czechoslovakia
	"704": true, // Viet Nam	Official name Socialist Republic of Viet Nam
	"705": true, // Slovenia	Until 1992 part of Yugoslavia
	"706": true, // Somalia
	"710": true, // South Africa
	"716": true, // Zimbabwe	Formerly Southern Rhodesia
	"724": true, // Spain
	"728": true, // South Sudan	Until 2011 part of Sudan
	"729": true, // Sudan
	"732": true, // Western Sahara	Formerly Spanish Sahara
	"740": true, // Suriname
	"744": true, // Svalbard and Jan Mayen
	"748": true, // Eswatini	Formerly Swaziland
	"752": true, // Sweden
	"756": true, // Switzerland
	"760": true, // Syrian Arab Republic
	"762": true, // Tajikistan	Until 1991 part of the USSR
	"764": true, // Thailand
	"768": true, // Togo
	"772": true, // Tokelau
	"776": true, // Tonga
	"780": true, // Trinidad and Tobago
	"784": true, // United Arab Emirates	Formerly Trucial States
	"788": true, // Tunisia
	"792": true, // Türkiye
	"795": true, // Turkmenistan	Until 1991 part of the USSR
	"796": true, // Turks and Caicos Islands
	"798": true, // Tuvalu
	"800": true, // Uganda
	"804": true, // Ukraine
	"807": true, // North Macedonia	Until 1993 part of Yugoslavia
	"818": true, // Egypt	Formerly United Arab Republic
	"826": true, // United Kingdom of Great Britain and Northern Ireland
	"831": true, // Guernsey	Until 2006 part of the United Kingdom
	"832": true, // Jersey	Until 2006 part of the United Kingdom
	"833": true, // Isle of Man	Until 2006 part of the United Kingdom
	"834": true, // Tanzania, United Republic of
	"840": true, // United States of America
	"850": true, // Virgin Islands (U.S.)
	"854": true, // Burkina Faso	Formerly Upper Volta
	"858": true, // Uruguay
	"860": true, // Uzbekistan	Until 1991 part of the USSR
	"862": true, // Venezuela (Bolivarian Republic of)
	"876": true, // Wallis and Futuna
	"882": true, // Samoa	Formerly Western Samoa
	"887": true, // Yemen	Unified country since 1990
	"894": true, // Zambia
}

var iSO3166_2_ObsoleteCodes = map[string]bool{
	"AL-BR": true, "AL-BU": true, "AL-DI": true, "AL-DL": true, "AL-DR": true, "AL-DV": true, "AL-EL": true, "AL-ER": true,
	"AL-FR": true, "AL-GJ": true, "AL-GR": true, "AL-HA": true, "AL-KA": true, "AL-KB": true, "AL-KC": true, "AL-KO": true,
	"AL-KR": true, "AL-KU": true, "AL-LB": true, "AL-LE": true, "AL-LU": true, "AL-MK": true, "AL-MM": true, "AL-MR": true,
	"AL-MT": true, "AL-PG": true, "AL-PQ": true, "AL-PR": true, "AL-PU": true, "AL-SH": true, "AL-SK": true, "AL-SR": true,
	"AL-TE": true, "AL-TP": true, "AL-TR": true, "AL-VL": true,
	"BA-05": true, "BA-07": true, "BA-10": true, "BA-09": true, "BA-02": true,
	"BA-06": true, "BA-03": true, "BA-01": true, "BA-08": true, "BA-04": true,
	"BH-16": true, "BR-FN": true,
	"CD-BN": true, "CD-KA": true, "CD-KW": true, "CD-OR": true,
	"CI-01": true, "CI-02": true, "CI-03": true, "CI-04": true, "CI-05": true, "CI-06": true, "CI-07": true, "CI-08": true, "CI-09": true,
	"CI-10": true, "CI-11": true, "CI-12": true, "CI-13": true, "CI-14": true, "CI-15": true, "CI-16": true, "CI-17": true, "CI-18": true, "CI-19": true,
	"CN-11": true, "CN-12": true, "CN-13": true, "CN-14": true, "CN-15": true,
	"CN-21": true, "CN-22": true, "CN-23": true,
	"CN-31": true, "CN-32": true, "CN-33": true, "CN-34": true, "CN-35": true, "CN-36": true, "CN-37": true,
	"CN-41": true, "CN-42": true, "CN-43": true, "CN-44": true, "CN-45": true, "CN-46": true,
	"CN-50": true, "CN-51": true, "CN-52": true, "CN-53": true, "CN-54": true,
	"CN-61": true, "CN-62": true, "CN-63": true, "CN-64": true, "CN-65": true,
	"CN-71": true,
	"CN-91": true, "CN-92": true,
	"CU-02":  true,
	"CZ-101": true, "CZ-102": true, "CZ-103": true, "CZ-104": true, "CZ-105": true, "CZ-106": true, "CZ-107": true, "CZ-108": true, "CZ-109": true,
	"CZ-110": true, "CZ-111": true, "CZ-112": true, "CZ-113": true, "CZ-114": true, "CZ-115": true, "CZ-116": true, "CZ-117": true, "CZ-118": true, "CZ-119": true,
	"CZ-120": true, "CZ-121": true, "CZ-122": true,
	"DM-01": true,
	"EE-44": true, "EE-49": true, "EE-51": true, "EE-57": true, "EE-59": true, "EE-65": true, "EE-67": true, "EE-70": true,
	"EE-78": true, "EE-82": true, "EE-86": true,
	"EG-HU": true, "EG-SU": true,
	"FR-75": true, "FR-COR": true, "FR-GF": true, "FR-GP": true, "FR-GUA": true, "FR-LRE": true, "FR-MAY": true, "FR-MQ": true, "FR-RE": true, "FR-YT": true,
	"GB-BMH": true, "GB-POL": true,
	"GH-BA": true,
	"GL-QA": true,
	"GR-01": true, "GR-03": true, "GR-04": true, "GR-05": true, "GR-06": true, "GR-07": true,
	"GR-11": true, "GR-12": true, "GR-13": true, "GR-14": true, "GR-15": true, "GR-16": true, "GR-17": true,
	"GR-21": true, "GR-22": true, "GR-23": true, "GR-24": true,
	"GR-31": true, "GR-32": true, "GR-33": true, "GR-34": true,
	"GR-41": true, "GR-42": true, "GR-43": true, "GR-44": true,
	"GR-51": true, "GR-52": true, "GR-53": true, "GR-54": true, "GR-55": true, "GR-56": true, "GR-57": true, "GR-58": true, "GR-59": true,
	"GR-61": true, "GR-62": true, "GR-63": true, "GR-64": true,
	"GR-71": true, "GR-72": true, "GR-73": true,
	"GR-81": true, "GR-82": true, "GR-83": true, "GR-84": true, "GR-85": true,
	"GR-91": true, "GR-92": true, "GR-93": true, "GR-94": true,
	"GR-A1": true,
	"GT-AV": true, "GT-BV": true, "GT-CM": true, "GT-CQ": true, "GT-ES": true, "GT-GU": true, "GT-HU": true, "GT-IZ": true,
	"GT-JA": true, "GT-JU": true, "GT-PE": true, "GT-PR": true, "GT-QC": true, "GT-QZ": true, "GT-RE": true, "GT-SA": true,
	"GT-SM": true, "GT-SO": true, "GT-SR": true, "GT-SU": true, "GT-TO": true, "GT-ZA": true,
	"ID-IJ": true,
	"IN-DD": true, "IN-DN": true,
	"IQ-SW": true, "IQ-TS": true,
	"IR-31": true,
	"IS-0":  true,
	"IT-AO": true, "IT-CI": true, "IT-OG": true, "IT-OT": true, "IT-VS": true,
	"KE-110": true, "KE-200": true, "KE-300": true, "KE-400": true, "KE-500": true, "KE-700": true, "KE-800": true,
	"LU-D": true, "LU-G": true, "LU-L": true,
	"LV-001": true, "LV-003": true, "LV-004": true, "LV-005": true, "LV-006": true, "LV-008": true, "LV-009": true,
	"LV-010": true, "LV-012": true, "LV-013": true, "LV-014": true, "LV-017": true, "LV-018": true, "LV-019": true,
	"LV-020": true, "LV-021": true, "LV-023": true, "LV-024": true, "LV-025": true, "LV-027": true, "LV-028": true, "LV-029": true,
	"LV-030": true, "LV-031": true, "LV-032": true, "LV-034": true, "LV-035": true, "LV-036": true, "LV-037": true, "LV-038": true, "LV-039": true,
	"LV-040": true, "LV-043": true, "LV-044": true, "LV-045": true, "LV-046": true, "LV-048": true, "LV-049": true,
	"LV-051": true, "LV-053": true, "LV-055": true, "LV-057": true,
	"LV-060": true, "LV-061": true, "LV-063": true, "LV-064": true, "LV-065": true, "LV-066": true, "LV-069": true,
	"LV-070": true, "LV-071": true, "LV-072": true, "LV-074": true, "LV-075": true, "LV-076": true, "LV-078": true, "LV-079": true,
	"LV-081": true, "LV-082": true, "LV-083": true, "LV-084": true, "LV-085": true, "LV-086": true,
	"LV-090": true, "LV-092": true, "LV-093": true, "LV-095": true, "LV-096": true, "LV-098": true,
	"LV-100": true, "LV-103": true, "LV-104": true, "LV-105": true, "LV-107": true, "LV-108": true, "LV-109": true,
	"LV-110": true,
	"LV-JKB": true, "LV-VMR": true,
	"LY-JB": true,
	"MA-13": true, "MA-14": true, "MA-15": true, "MA-16": true,
	"MA-MMD": true, "MA-MMN": true, "MA-SYB": true,
	"MH-WTN": true,
	"MK-01":  true, "MK-02": true, "MK-03": true, "MK-04": true, "MK-05": true, "MK-06": true, "MK-07": true, "MK-08": true, "MK-09": true,
	"MK-10": true, "MK-11": true, "MK-12": true, "MK-13": true, "MK-14": true, "MK-15": true, "MK-16": true, "MK-17": true, "MK-18": true, "MK-19": true,
	"MK-20": true, "MK-21": true, "MK-22": true, "MK-23": true, "MK-24": true, "MK-25": true, "MK-26": true, "MK-27": true, "MK-28": true, "MK-29": true,
	"MK-30": true, "MK-31": true, "MK-32": true, "MK-33": true, "MK-34": true, "MK-35": true, "MK-36": true, "MK-37": true, "MK-38": true, "MK-39": true,
	"MK-40": true, "MK-41": true, "MK-42": true, "MK-43": true, "MK-44": true, "MK-45": true, "MK-46": true, "MK-47": true, "MK-48": true, "MK-49": true,
	"MK-50": true, "MK-51": true, "MK-52": true, "MK-53": true, "MK-54": true, "MK-55": true, "MK-56": true, "MK-57": true, "MK-58": true, "MK-59": true,
	"MK-60": true, "MK-61": true, "MK-62": true, "MK-63": true, "MK-64": true, "MK-65": true, "MK-66": true, "MK-67": true, "MK-68": true, "MK-69": true,
	"MK-70": true, "MK-71": true, "MK-72": true, "MK-73": true, "MK-74": true, "MK-75": true, "MK-76": true, "MK-77": true, "MK-78": true, "MK-79": true,
	"MK-80": true, "MK-81": true, "MK-82": true, "MK-83": true, "MK-84": true,
	"ML-BK0": true,
	"MR-NKC": true,
	"MU-BR":  true, "MU-CU": true, "MU-PU": true, "MU-QB": true, "MU-RP": true, "MU-VP": true,
	"MV-CE": true, "MV-NC": true, "MV-NO": true, "MV-SC": true, "MV-SU": true, "MV-UN": true, "MV-US": true,
	"MX-DIF": true,
	"NA-OK":  true,
	"NO-01":  true, "NO-02": true, "NO-04": true, "NO-05": true, "NO-06": true, "NO-07": true, "NO-08": true, "NO-09": true,
	"NO-10": true, "NO-12": true, "NO-14": true, "NO-16": true, "NO-17": true, "NO-19": true, "NO-20": true,
	"NZ-N": true, "NZ-S": true,
	"OM-BA": true, "OM-SH": true,
	"PK-TA": true,
	"PL-DS": true, "PL-KP": true, "PL-LB": true, "PL-LD": true, "PL-LU": true, "PL-MA": true, "PL-MZ": true, "PL-OP": true,
	"PL-PD": true, "PL-PK": true, "PL-PM": true, "PL-SK": true, "PL-SL": true, "PL-WN": true, "PL-WP": true, "PL-ZP": true,
	"SS-EE8": true,
	"ST-S":   true,
	"TD-EN":  true,
	"TT-ETO": true, "TT-RCM": true, "TT-WTO": true,
	"TW-KHQ": true, "TW-TNQ": true, "TW-TPQ": true, "TW-TXQ": true,
	"VN-15": true,
	"YE-MU": true,
	"ZA-ZN": true,
}

type assignmentType3166_1 int

const (
	ccAs assignmentType3166_1 = iota // assigned
	ccRa                             // re-assigned
	ccUA                             // user-assigned
	ccER                             // exceptionally reserved
	ccIR                             // indeterminately reserved
	ccTR                             // transitionally reserved
	ccDl                             // deleted
	ccUn                             // Unassigned
)

type countryCodeMatrix map[byte]map[byte]assignmentType3166_1

var iso3166_1_CountryCodesMatrix = countryCodeMatrix{
	'A': {
		'A': ccUA, 'B': ccUn, 'C': ccER, 'D': ccAs, 'E': ccAs, 'F': ccAs, 'G': ccAs, 'H': ccUn, 'I': ccRa, 'J': ccUn, 'K': ccUn, 'L': ccAs, 'M': ccAs,
		'N': ccTR, 'O': ccAs, 'P': ccIR, 'Q': ccAs, 'R': ccAs, 'S': ccAs, 'T': ccAs, 'U': ccAs, 'V': ccUn, 'W': ccAs, 'X': ccAs, 'Y': ccUn, 'Z': ccAs,
	},
	'B': {
		'A': ccAs, 'B': ccAs, 'C': ccUn, 'D': ccAs, 'E': ccAs, 'F': ccAs, 'G': ccAs, 'H': ccAs, 'I': ccAs, 'J': ccAs, 'K': ccUn, 'L': ccAs, 'M': ccAs,
		'N': ccAs, 'O': ccAs, 'P': ccUn, 'Q': ccRa, 'R': ccAs, 'S': ccAs, 'T': ccAs, 'U': ccTR, 'V': ccAs, 'W': ccAs, 'X': ccIR, 'Y': ccAs, 'Z': ccAs,
	},
	'C': {
		'A': ccAs, 'B': ccUn, 'C': ccAs, 'D': ccAs, 'E': ccUn, 'F': ccAs, 'G': ccAs, 'H': ccAs, 'I': ccAs, 'J': ccUn, 'K': ccAs, 'L': ccAs, 'M': ccAs,
		'N': ccAs, 'O': ccAs, 'P': ccER, 'Q': ccER, 'R': ccAs, 'S': ccTR, 'T': ccDl, 'U': ccAs, 'V': ccAs, 'W': ccAs, 'X': ccAs, 'Y': ccAs, 'Z': ccAs,
	},
	'D': {
		'A': ccUn, 'B': ccUn, 'C': ccUn, 'D': ccDl, 'E': ccAs, 'F': ccUn, 'G': ccER, 'H': ccUn, 'I': ccUn, 'J': ccAs, 'K': ccAs, 'L': ccUn, 'M': ccAs,
		'N': ccUn, 'O': ccAs, 'P': ccUn, 'Q': ccUn, 'R': ccUn, 'S': ccUn, 'T': ccUn, 'U': ccUn, 'V': ccUn, 'W': ccUn, 'X': ccUn, 'Y': ccIR, 'Z': ccAs,
	},
	'E': {
		'A': ccER, 'B': ccUn, 'C': ccAs, 'D': ccUn, 'E': ccAs, 'F': ccIR, 'G': ccAs, 'H': ccAs, 'I': ccUn, 'J': ccUn, 'K': ccUn, 'L': ccUn, 'M': ccIR,
		'N': ccUn, 'O': ccUn, 'P': ccIR, 'Q': ccUn, 'R': ccAs, 'S': ccAs, 'T': ccAs, 'U': ccER, 'V': ccIR, 'W': ccIR, 'X': ccUn, 'Y': ccUn, 'Z': ccER,
	},
	'F': {
		'A': ccUn, 'B': ccUn, 'C': ccUn, 'D': ccUn, 'E': ccUn, 'F': ccUn, 'G': ccUn, 'H': ccUn, 'I': ccAs, 'J': ccAs, 'K': ccAs, 'L': ccIR, 'M': ccAs,
		'N': ccUn, 'O': ccAs, 'P': ccUn, 'Q': ccDl, 'R': ccAs, 'S': ccUn, 'T': ccUn, 'U': ccUn, 'V': ccUn, 'W': ccUn, 'X': ccER, 'Y': ccUn, 'Z': ccUn,
	},
	'G': {
		'A': ccAs, 'B': ccAs, 'C': ccIR, 'D': ccAs, 'E': ccRa, 'F': ccAs, 'G': ccAs, 'H': ccAs, 'I': ccAs, 'J': ccUn, 'K': ccUn, 'L': ccAs, 'M': ccAs,
		'N': ccAs, 'O': ccUn, 'P': ccAs, 'Q': ccAs, 'R': ccAs, 'S': ccAs, 'T': ccAs, 'U': ccAs, 'V': ccUn, 'W': ccAs, 'X': ccUn, 'Y': ccAs, 'Z': ccUn,
	},
	'H': {
		'A': ccUn, 'B': ccUn, 'C': ccUn, 'D': ccUn, 'E': ccUn, 'F': ccUn, 'G': ccUn, 'H': ccUn, 'I': ccUn, 'J': ccUn, 'K': ccAs, 'L': ccUn, 'M': ccAs,
		'N': ccAs, 'O': ccUn, 'P': ccUn, 'Q': ccUn, 'R': ccAs, 'S': ccUn, 'T': ccAs, 'U': ccAs, 'V': ccDl, 'W': ccUn, 'X': ccUn, 'Y': ccUn, 'Z': ccUn,
	},
	'I': {
		'A': ccUn, 'B': ccIR, 'C': ccER, 'D': ccAs, 'E': ccAs, 'F': ccUn, 'G': ccUn, 'H': ccUn, 'I': ccUn, 'J': ccUn, 'K': ccUn, 'L': ccAs, 'M': ccAs,
		'N': ccAs, 'O': ccAs, 'P': ccUn, 'Q': ccAs, 'R': ccAs, 'S': ccAs, 'T': ccAs, 'U': ccUn, 'V': ccUn, 'W': ccUn, 'X': ccUn, 'Y': ccUn, 'Z': ccUn,
	},
	'J': {
		'A': ccIR, 'B': ccUn, 'C': ccUn, 'D': ccUn, 'E': ccAs, 'F': ccUn, 'G': ccUn, 'H': ccUn, 'I': ccUn, 'J': ccUn, 'K': ccUn, 'L': ccUn, 'M': ccAs,
		'N': ccUn, 'O': ccAs, 'P': ccAs, 'Q': ccUn, 'R': ccUn, 'S': ccUn, 'T': ccDl, 'U': ccUn, 'V': ccUn, 'W': ccUn, 'X': ccUn, 'Y': ccUn, 'Z': ccUn,
	},
	'K': {
		'A': ccUn, 'B': ccUn, 'C': ccUn, 'D': ccUn, 'E': ccAs, 'F': ccUn, 'G': ccAs, 'H': ccAs, 'I': ccAs, 'J': ccUn, 'K': ccUn, 'L': ccUn, 'M': ccAs,
		'N': ccAs, 'O': ccUn, 'P': ccAs, 'Q': ccUn, 'R': ccAs, 'S': ccUn, 'T': ccUn, 'U': ccUn, 'V': ccUn, 'W': ccAs, 'X': ccUn, 'Y': ccAs, 'Z': ccAs,
	},
	'L': {
		'A': ccAs, 'B': ccAs, 'C': ccAs, 'D': ccUn, 'E': ccUn, 'F': ccIR, 'G': ccUn, 'H': ccUn, 'I': ccAs, 'J': ccUn, 'K': ccAs, 'L': ccUn, 'M': ccUn,
		'N': ccUn, 'O': ccUn, 'P': ccUn, 'Q': ccUn, 'R': ccAs, 'S': ccAs, 'T': ccRa, 'U': ccAs, 'V': ccAs, 'W': ccUn, 'X': ccUn, 'Y': ccAs, 'Z': ccUn,
	},
	'M': {
		'A': ccAs, 'B': ccUn, 'C': ccAs, 'D': ccAs, 'E': ccRa, 'F': ccAs, 'G': ccAs, 'H': ccAs, 'I': ccDl, 'J': ccUn, 'K': ccAs, 'L': ccAs, 'M': ccAs,
		'N': ccAs, 'O': ccAs, 'P': ccAs, 'Q': ccAs, 'R': ccAs, 'S': ccAs, 'T': ccAs, 'U': ccAs, 'V': ccAs, 'W': ccAs, 'X': ccAs, 'Y': ccAs, 'Z': ccAs,
	},
	'N': {
		'A': ccAs, 'B': ccUn, 'C': ccAs, 'D': ccUn, 'E': ccAs, 'F': ccAs, 'G': ccAs, 'H': ccDl, 'I': ccAs, 'J': ccUn, 'K': ccUn, 'L': ccAs, 'M': ccUn,
		'N': ccUn, 'O': ccAs, 'P': ccAs, 'Q': ccDl, 'R': ccAs, 'S': ccUn, 'T': ccTR, 'U': ccAs, 'V': ccUn, 'W': ccUn, 'X': ccUn, 'Y': ccUn, 'Z': ccAs,
	},
	'O': {
		'A': ccIR, 'B': ccUn, 'C': ccUn, 'D': ccUn, 'E': ccUn, 'F': ccUn, 'G': ccUn, 'H': ccUn, 'I': ccUn, 'J': ccUn, 'K': ccUn, 'L': ccUn, 'M': ccAs,
		'N': ccUn, 'O': ccUA, 'P': ccUn, 'Q': ccUn, 'R': ccUn, 'S': ccUn, 'T': ccUn, 'U': ccUn, 'V': ccUn, 'W': ccUn, 'X': ccUn, 'Y': ccUn, 'Z': ccUn,
	},
	'P': {
		'A': ccAs, 'B': ccUn, 'C': ccDl, 'D': ccUn, 'E': ccAs, 'F': ccAs, 'G': ccAs, 'H': ccAs, 'I': ccIR, 'J': ccUn, 'K': ccAs, 'L': ccAs, 'M': ccAs,
		'N': ccAs, 'O': ccUn, 'P': ccUn, 'Q': ccUn, 'R': ccAs, 'S': ccAs, 'T': ccAs, 'U': ccDl, 'V': ccUn, 'W': ccAs, 'X': ccUn, 'Y': ccAs, 'Z': ccDl,
	},
	'Q': {
		'A': ccAs, 'B': ccUn, 'C': ccUn, 'D': ccUn, 'E': ccUn, 'F': ccUn, 'G': ccUn, 'H': ccUn, 'I': ccUn, 'J': ccUn, 'K': ccUn, 'L': ccUn, 'M': ccUA,
		'N': ccUA, 'O': ccUA, 'P': ccUA, 'Q': ccUA, 'R': ccUA, 'S': ccUA, 'T': ccUA, 'U': ccUA, 'V': ccUA, 'W': ccUA, 'X': ccUA, 'Y': ccUA, 'Z': ccUA,
	},
	'R': {
		'A': ccIR, 'B': ccIR, 'C': ccIR, 'D': ccUn, 'E': ccAs, 'F': ccUn, 'G': ccUn, 'H': ccIR, 'I': ccIR, 'J': ccUn, 'K': ccUn, 'L': ccIR, 'M': ccIR,
		'N': ccIR, 'O': ccAs, 'P': ccIR, 'Q': ccUn, 'R': ccUn, 'S': ccAs, 'T': ccUn, 'U': ccRa, 'V': ccUn, 'W': ccAs, 'X': ccUn, 'Y': ccUn, 'Z': ccUn,
	},
	'S': {
		'A': ccAs, 'B': ccAs, 'C': ccAs, 'D': ccAs, 'E': ccAs, 'F': ccIR, 'G': ccAs, 'H': ccAs, 'I': ccAs, 'J': ccAs, 'K': ccRa, 'L': ccAs, 'M': ccAs,
		'N': ccAs, 'O': ccAs, 'P': ccUn, 'Q': ccUn, 'R': ccAs, 'S': ccAs, 'T': ccAs, 'U': ccER, 'V': ccAs, 'W': ccUn, 'X': ccAs, 'Y': ccAs, 'Z': ccAs,
	},
	'T': {
		'A': ccER, 'B': ccUn, 'C': ccAs, 'D': ccAs, 'E': ccUn, 'F': ccAs, 'G': ccAs, 'H': ccAs, 'I': ccUn, 'J': ccAs, 'K': ccAs, 'L': ccAs, 'M': ccAs,
		'N': ccAs, 'O': ccAs, 'P': ccTR, 'Q': ccUn, 'R': ccAs, 'S': ccUn, 'T': ccAs, 'U': ccUn, 'V': ccAs, 'W': ccAs, 'X': ccUn, 'Y': ccUn, 'Z': ccAs,
	},
	'U': {
		'A': ccAs, 'B': ccUn, 'C': ccUn, 'D': ccUn, 'E': ccUn, 'F': ccUn, 'G': ccAs, 'H': ccUn, 'I': ccUn, 'J': ccUn, 'K': ccER, 'L': ccUn, 'M': ccAs,
		'N': ccER, 'O': ccUn, 'P': ccUn, 'Q': ccUn, 'R': ccUn, 'S': ccAs, 'T': ccUn, 'U': ccUn, 'V': ccUn, 'W': ccUn, 'X': ccUn, 'Y': ccAs, 'Z': ccAs,
	},
	'V': {
		'A': ccAs, 'B': ccUn, 'C': ccAs, 'D': ccDl, 'E': ccAs, 'F': ccUn, 'G': ccAs, 'H': ccUn, 'I': ccAs, 'J': ccUn, 'K': ccUn, 'L': ccUn, 'M': ccUn,
		'N': ccAs, 'O': ccUn, 'P': ccUn, 'Q': ccUn, 'R': ccUn, 'S': ccUn, 'T': ccUn, 'U': ccAs, 'V': ccUn, 'W': ccUn, 'X': ccUn, 'Y': ccUn, 'Z': ccUn,
	},
	'W': {
		'A': ccUn, 'B': ccUn, 'C': ccUn, 'D': ccUn, 'E': ccUn, 'F': ccAs, 'G': ccIR, 'H': ccUn, 'I': ccUn, 'J': ccUn, 'K': ccDl, 'L': ccIR, 'M': ccUn,
		'N': ccUn, 'O': ccIR, 'P': ccUn, 'Q': ccUn, 'R': ccUn, 'S': ccAs, 'T': ccUn, 'U': ccUn, 'V': ccIR, 'W': ccUn, 'X': ccUn, 'Y': ccUn, 'Z': ccUn,
	},
	'X': {
		'A': ccUA, 'B': ccUA, 'C': ccUA, 'D': ccUA, 'E': ccUA, 'F': ccUA, 'G': ccUA, 'H': ccUA, 'I': ccUA, 'J': ccUA, 'K': ccUA, 'L': ccUA, 'M': ccUA,
		'N': ccUA, 'O': ccUA, 'P': ccUA, 'Q': ccUA, 'R': ccUA, 'S': ccUA, 'T': ccUA, 'U': ccUA, 'V': ccUA, 'W': ccUA, 'X': ccUA, 'Y': ccUA, 'Z': ccUA,
	},
	'Y': {
		'A': ccUn, 'B': ccUn, 'C': ccUn, 'D': ccDl, 'E': ccAs, 'F': ccUn, 'G': ccUn, 'H': ccUn, 'I': ccUn, 'J': ccUn, 'K': ccUn, 'L': ccUn, 'M': ccUn,
		'N': ccUn, 'O': ccUn, 'P': ccUn, 'Q': ccUn, 'R': ccUn, 'S': ccUn, 'T': ccAs, 'U': ccTR, 'V': ccIR, 'W': ccUn, 'X': ccUn, 'Y': ccUn, 'Z': ccUn,
	},
	'Z': {
		'A': ccAs, 'B': ccUn, 'C': ccUn, 'D': ccUn, 'E': ccUn, 'F': ccUn, 'G': ccUn, 'H': ccUn, 'I': ccUn, 'J': ccUn, 'K': ccUn, 'L': ccUn, 'M': ccAs,
		'N': ccUn, 'O': ccUn, 'P': ccUn, 'Q': ccUn, 'R': ccTR, 'S': ccUn, 'T': ccUn, 'U': ccUn, 'V': ccUn, 'W': ccAs, 'X': ccUn, 'Y': ccUn, 'Z': ccUA,
	},
}
