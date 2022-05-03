package valix

import (
	"fmt"
	"strings"
)

// DefaultTranslator is the default translator used by valix I18n support
//
// Replace it with your own where necessary
var DefaultTranslator Translator = defaultInternalTranslator

// Translations created using Google Translate - apologies for any mistakes
// (correct them, if needed, in these variables)

type Translator interface {
	TranslateToken(lang string, region string, token string) string
	TranslateMessage(lang string, region string, message string) string
	TranslateFormat(lang string, region string, format string, a ...interface{}) string
	AddTokenLanguageTranslation(lang string, token string, translation string, regionals ...RegionalVariantTranslation)
	AddMessageLanguageTranslation(lang string, message string, translation string, regionals ...RegionalVariantTranslation)
	AddFormatLanguageTranslation(lang string, format string, translation string, regionals ...RegionalVariantTranslation)
	AddTokenRegionTranslation(lang string, region string, token string, translation string)
	AddMessageRegionTranslation(lang string, region string, message string, translation string)
	AddFormatRegionTranslation(lang string, region string, format string, translation string)
}

// RegionalVariantTranslation for use with Add...LanguageTranslation methods of Translator
type RegionalVariantTranslation struct {
	// is the region (must not be "" or will not be added)
	Region string
	// the regional translation (if "" uses the parent language translation)
	Translation string
}

type internalTranslator struct {
	Tokens   map[string]map[string]string `json:"tokens"`
	Messages map[string]map[string]string `json:"messages"`
	Formats  map[string]map[string]string `json:"formats"`
}

func lookupTranslation(trs map[string]map[string]string, str string, lang string, rgn string) string {
	result := str
	if ts, ok := trs[str]; ok {
		if rgn != "" {
			if trr, ok := ts[lang+"-"+rgn]; ok {
				result = trr
			} else if tr, ok := ts[lang]; ok {
				result = tr
			}
		} else if tr, ok := ts[lang]; ok {
			result = tr
		}
	}
	return result
}

func (t *internalTranslator) TranslateToken(lang string, region string, token string) string {
	return lookupTranslation(t.Tokens, token, lang, region)
}

func (t *internalTranslator) TranslateMessage(lang string, region string, message string) string {
	return lookupTranslation(t.Messages, message, lang, region)
}

func (t *internalTranslator) TranslateFormat(lang string, region string, format string, a ...interface{}) string {
	return fmt.Sprintf(lookupTranslation(t.Formats, format, lang, region), a...)
}

func (t *internalTranslator) AddTokenLanguageTranslation(lang string, token string, translation string, regionals ...RegionalVariantTranslation) {
	tr, present := t.Tokens[token]
	if !present {
		tr = map[string]string{
			strings.ToLower(lang): translation,
		}
		t.Tokens[token] = tr
	}
	for _, regional := range regionals {
		if regional.Region != "" {
			if regional.Translation != "" {
				tr[strings.ToLower(lang)+"-"+strings.ToUpper(regional.Region)] = regional.Translation
			} else {
				tr[strings.ToLower(lang)+"-"+strings.ToUpper(regional.Region)] = translation
			}
		}
	}
}

func (t *internalTranslator) AddTokenRegionTranslation(lang string, region string, token string, translation string) {
	if region != "" {
		if tr, ok := t.Tokens[token]; ok {
			tr[strings.ToLower(lang)+"-"+strings.ToUpper(region)] = translation
		} else {
			newTr := map[string]string{
				strings.ToLower(lang): translation,
			}
			if region != "" {
				newTr[strings.ToLower(lang)+"-"+strings.ToUpper(region)] = translation
			}
			t.Tokens[token] = newTr
		}
	}
}

func (t *internalTranslator) AddMessageLanguageTranslation(lang string, message string, translation string, regionals ...RegionalVariantTranslation) {
	tr, present := t.Messages[message]
	if !present {
		tr = map[string]string{
			strings.ToLower(lang): translation,
		}
		t.Messages[message] = tr
	}
	for _, regional := range regionals {
		if regional.Region != "" {
			if regional.Translation != "" {
				tr[strings.ToLower(lang)+"-"+strings.ToUpper(regional.Region)] = regional.Translation
			} else {
				tr[strings.ToLower(lang)+"-"+strings.ToUpper(regional.Region)] = translation
			}
		}
	}
}

func (t *internalTranslator) AddMessageRegionTranslation(lang string, region string, message string, translation string) {
	if region != "" {
		if tr, ok := t.Messages[message]; ok {
			tr[strings.ToLower(lang)+"-"+strings.ToUpper(region)] = translation
		} else {
			newTr := map[string]string{
				strings.ToLower(lang): translation,
			}
			if region != "" {
				newTr[strings.ToLower(lang)+"-"+strings.ToUpper(region)] = translation
			}
			t.Messages[message] = newTr
		}
	}
}

func (t *internalTranslator) AddFormatLanguageTranslation(lang string, format string, translation string, regionals ...RegionalVariantTranslation) {
	tr, present := t.Formats[format]
	if !present {
		tr = map[string]string{
			strings.ToLower(lang): translation,
		}
		t.Formats[format] = tr
	}
	for _, regional := range regionals {
		if regional.Region != "" {
			if regional.Translation != "" {
				tr[strings.ToLower(lang)+"-"+strings.ToUpper(regional.Region)] = regional.Translation
			} else {
				tr[strings.ToLower(lang)+"-"+strings.ToUpper(regional.Region)] = translation
			}
		}
	}
}

func (t *internalTranslator) AddFormatRegionTranslation(lang string, region string, format string, translation string) {
	if region != "" {
		if tr, ok := t.Formats[format]; ok {
			tr[strings.ToLower(lang)+"-"+strings.ToUpper(region)] = translation
		} else {
			newTr := map[string]string{
				strings.ToLower(lang): translation,
			}
			if region != "" {
				newTr[strings.ToLower(lang)+"-"+strings.ToUpper(region)] = translation
			}
			t.Formats[format] = newTr
		}
	}
}

var defaultInternalTranslator = &internalTranslator{
	Tokens: map[string]map[string]string{
		jsonTypeTokenAny: {
			"en": jsonTypeTokenAny,
			"fr": "tout",
			"es": "alguno",
			"it": "qualsiasi",
			"de": "beliebig",
		},
		jsonTypeTokenArray: {
			"en": jsonTypeTokenArray,
			"fr": "tableau",
			"es": "matriz",
			"it": "array",
			"de": "Array",
		},
		jsonTypeTokenBoolean: {
			"en": jsonTypeTokenBoolean,
			"fr": "booléen",
			"es": "booleano",
			"it": "booleano",
			"de": "boolesch",
		},
		jsonTypeTokenInteger: {
			"en": jsonTypeTokenInteger,
			"fr": "entier",
			"es": "entero",
			"it": "intero",
			"de": "Ganzzahl",
		},
		jsonTypeTokenNumber: {
			"en": jsonTypeTokenNumber,
			"fr": "nombre",
			"es": "número",
			"it": "numero",
			"de": "Nummer",
		},
		jsonTypeTokenObject: {
			"en": jsonTypeTokenObject,
			"fr": "objet",
			"es": "objeto",
			"it": "oggetto",
			"de": "Objekt",
		},
		jsonTypeTokenString: {
			"en": jsonTypeTokenString,
			"fr": "chaîne",
			"es": "cadena",
			"it": "stringa",
			"de": "Zeichenfolge",
		},
		tokenExclusive: {
			"en": tokenExclusive,
			"fr": "exclusif",
			"es": "exclusivo",
			"it": "esclusivo",
			"de": "exklusiv",
		},
		tokenInclusive: {
			"en": tokenInclusive,
			"fr": "inclusif",
			"es": "inclusivo",
			"it": "comprensivo",
			"de": "inklusive",
		},
		"century": {
			"en": "century",
			"fr": "siècle",
			"es": "siglo",
			"it": "secolo",
			"de": "Jahrhundert",
		},
		"century...": {
			"en": "centuries",
			"fr": "siècles",
			"es": "siglos",
			"it": "secoli",
			"de": "Jahrhunderte",
		},
		"day": {
			"en": "day",
			"fr": "jour",
			"es": "día",
			"it": "giorno",
			"de": "Tag",
		},
		"day...": {
			"en": "days",
			"fr": "jours",
			"es": "días",
			"it": "giorni",
			"de": "Tage",
		},
		"decade": {
			"en": "decade",
			"fr": "décennie",
			"es": "década",
			"it": "decennio",
			"de": "Jahrzehnt",
		},
		"decade...": {
			"en": "decades",
			"fr": "décennies",
			"es": "décadas",
			"it": "decenni",
			"de": "Jahrzehnte",
		},
		"hour": {
			"en": "hour",
			"fr": "heure",
			"es": "hora",
			"it": "ora",
			"de": "Stunde",
		},
		"hour...": {
			"en": "hours",
			"fr": "heures",
			"es": "horas",
			"it": "ore",
			"de": "Stunden",
		},
		"microsecond": {
			"en": "microsecond",
			"fr": "microsecondes",
			"es": "microsegundos",
			"it": "microsecondo",
			"de": "Mikrosekunde",
		},
		"microsecond...": {
			"en": "microseconds",
			"fr": "microsecondes",
			"es": "microsegundos",
			"it": "microsecondi",
			"de": "Mikrosekunden",
		},
		"millennium": {
			"en": "millennium",
			"fr": "millénaire",
			"es": "milenio",
			"it": "millennio",
			"de": "Jahrtausend",
		},
		"millennium...": {
			"en": "millennia",
			"fr": "millénaires",
			"es": "milenios",
			"it": "millenni",
			"de": "Jahrtausende",
		},
		"millisecond": {
			"en": "millisecond",
			"fr": "milliseconde",
			"es": "milisegundos",
			"it": "millisecondo",
			"de": "Millisekunde",
		},
		"millisecond...": {
			"en": "milliseconds",
			"fr": "millisecondes",
			"es": "milisegundos",
			"it": "millisecondi",
			"de": "Millisekunden",
		},
		"minute": {
			"en": "minute",
			"fr": "minute",
			"es": "minuto",
			"it": "minuto",
			"de": "Minute",
		},
		"minute...": {
			"en": "minutes",
			"fr": "minutes",
			"es": "minutos",
			"it": "minuti",
			"de": "Minuten",
		},
		"month": {
			"en": "month",
			"fr": "mois",
			"es": "mes",
			"it": "mese",
			"de": "Monat",
		},
		"month...": {
			"en": "months",
			"fr": "mois",
			"es": "meses",
			"it": "mesi",
			"de": "Monate",
		},
		"nanosecond": {
			"en": "nanosecond",
			"fr": "nanoseconde",
			"es": "nanosegundo",
			"it": "nanosecondo",
			"de": "Nanosekunde",
		},
		"nanosecond...": {
			"en": "nanoseconds",
			"fr": "nanosecondes",
			"es": "nanosegundos",
			"it": "nanosecondi",
			"de": "Nanosekunden",
		},
		"second": {
			"en": "second",
			"fr": "seconde",
			"es": "segundo",
			"it": "secondo",
			"de": "Sekunde",
		},
		"second...": {
			"en": "seconds",
			"fr": "secondes",
			"es": "segundos",
			"it": "secondi",
			"de": "Sekunden",
		},
		"week": {
			"en": "week",
			"fr": "semaine",
			"es": "semana",
			"it": "settimana",
			"de": "Woche",
		},
		"week...": {
			"en": "weeks",
			"fr": "semaines",
			"es": "semanas",
			"it": "settimane",
			"de": "Wochen",
		},
		"year": {
			"en": "year",
			"fr": "année",
			"es": "año",
			"it": "anno",
			"de": "Jahr",
		},
		"year...": {
			"en": "years",
			"fr": "ans",
			"es": "años",
			"it": "anni",
			"de": "Jahre",
		},
	},
	Messages: map[string]map[string]string{
		msgArrayElementMustBeObject: {
			"en": msgArrayElementMustBeObject,
			"fr": "L'élément du tableau JSON doit être un objet",
			"es": "El elemento de matriz JSON debe ser un objeto",
			"it": "L'elemento dell'array JSON deve essere un oggetto",
			"de": "JSON-Array-Element muss ein Objekt sein",
		},
		msgArrayUnique: {
			"en": msgArrayUnique,
			"fr": "Les éléments du tableau doivent être uniques",
			"es": "Los elementos del arreglo deben ser únicos",
			"it": "Gli elementi dell'array devono essere univoci",
			"de": "Array-Elemente müssen eindeutig sein",
		},
		msgDatetimeFuture: {
			"en": msgDatetimeFuture,
			"fr": "La valeur doit être une date/heure valide dans le futur",
			"es": "El valor debe ser una fecha/hora válida en el futuro",
			"it": "Il valore deve essere una data/ora valida nel futuro",
			"de": "Wert muss ein gültiges Datum/Zeit in der Zukunft sein",
		},
		msgDatetimeFutureOrPresent: {
			"en": msgDatetimeFutureOrPresent,
			"fr": "La valeur doit être une date/heure valide dans le futur ou le présent",
			"es": "El valor debe ser una fecha/hora válida en el futuro o presente",
			"it": "Il valore deve essere una data/ora valida futura o presente",
			"de": "Wert muss ein gültiges Datum/Zeit in der Zukunft oder Gegenwart sein",
		},
		msgDatetimePast: {
			"en": msgDatetimePast,
			"fr": "La valeur doit être une date/heure valide dans le passé",
			"es": "El valor debe ser una fecha/hora válida en el pasado",
			"it": "Il valore deve essere una data/ora valida nel passato",
			"de": "Wert muss ein gültiges Datum/Zeit in der Vergangenheit sein",
		},
		msgDatetimePastOrPresent: {
			"en": msgDatetimePastOrPresent,
			"fr": "La valeur doit être une date/heure valide dans le passé ou le présent",
			"es": "El valor debe ser una fecha/hora válida en el pasado o presente",
			"it": "Il valore deve essere una data/ora valida nel passato o nel presente",
			"de": "Wert muss ein gültiges Datum/Zeit in der Vergangenheit oder Gegenwart sein",
		},
		msgErrorReading: {
			"en": msgErrorReading,
			"fr": "Erreur inattendue lors de la lecture du lecteur",
			"es": "Error inesperado al leer el lector",
			"it": "Errore imprevisto durante la lettura del lettore",
			"de": "Unerwarteter Fehler beim Lesen des Lesegeräts",
		},
		msgErrorUnmarshall: {
			"en": msgErrorUnmarshall,
			"fr": "Erreur inattendue lors du démarshalling",
			"es": "Error inesperado durante la desorganización",
			"it": "Errore imprevisto durante l'annullamento del marshalling",
			"de": "Unerwarteter Fehler beim Unmarshalling",
		},
		msgExpectedJsonArray: {
			"en": msgExpectedJsonArray,
			"fr": "JSON devrait être un tableau JSON",
			"es": "Se esperaba que JSON fuera una matriz JSON",
			"it": "JSON dovrebbe essere un array JSON",
			"de": "JSON soll JSON-Array sein",
		},
		msgExpectedJsonObject: {
			"en": msgExpectedJsonObject,
			"fr": "JSON devrait être un objet JSON",
			"es": "Se esperaba que JSON fuera un objeto JSON",
			"it": "JSON dovrebbe essere un oggetto JSON",
			"de": "JSON soll JSON-Objekt sein",
		},
		msgFailure: {
			"en": msgFailure,
			"fr": "Échec de la validation",
			"es": "Validación fallida",
			"it": "Convalida non riuscita",
			"de": "Validierung fehlgeschlagen",
		},
		msgInvalidCharacters: {
			"en": msgInvalidCharacters,
			"fr": "La valeur de la chaîne ne doit pas contenir de caractères invalides",
			"es": "El valor de la cadena no debe tener caracteres inválidos",
			"it": "Il valore della stringa non deve contenere caratteri non validi",
			"de": "String-Wert darf keine ungültigen Zeichen enthalten",
		},
		msgInvalidProperty: {
			"en": msgInvalidProperty,
			"fr": "Propriété invalide",
			"es": "Propiedad no válida",
			"it": "Proprietà non valida",
			"de": "Ungültige Eigenschaft",
		},
		msgInvalidPropertyName: {
			"en": msgInvalidPropertyName,
			"fr": "Nom de propriété invalide",
			"es": "Nombre de propiedad inválido",
			"it": "Nome proprietà non valido",
			"de": "Ungültiger Eigenschaftsname",
		},
		msgMissingProperty: {
			"en": msgMissingProperty,
			"fr": "Propriété manquante",
			"es": "Propiedad faltante",
			"it": "Proprietà mancante",
			"de": "Fehlende Eigenschaft",
		},
		msgNegative: {
			"en": msgNegative,
			"fr": "La valeur doit être négative",
			"es": "El valor debe ser negativo",
			"it": "Il valore deve essere negativo",
			"de": "Wert muss negativ sein",
		},
		msgNegativeOrZero: {
			"en": msgNegativeOrZero,
			"fr": "La valeur doit être négative ou nulle",
			"es": "El valor debe ser negativo o cero",
			"it": "Il valore deve essere negativo o zero",
			"de": "Wert muss negativ oder Null sein",
		},
		msgNoControlChars: {
			"en": msgNoControlChars,
			"fr": "La valeur de la chaîne ne doit pas contenir de caractères de contrôle",
			"es": "El valor de la cadena no debe contener caracteres de control",
			"it": "Il valore della stringa non deve contenere caratteri di controllo",
			"de": "Stringwert darf keine Steuerzeichen enthalten",
		},
		msgNotBlankString: {
			"en": msgNotBlankString,
			"fr": "La valeur de la chaîne ne doit pas être une chaîne vide",
			"es": "El valor de la cadena no debe ser una cadena en blanco",
			"it": "Il valore della stringa non deve essere una stringa vuota",
			"de": "String-Wert darf kein leerer String sein",
		},
		msgNotEmptyString: {
			"en": msgNotEmptyString,
			"fr": "La valeur de la chaîne ne doit pas être une chaîne vide",
			"es": "El valor de la cadena no debe ser una cadena vacía",
			"it": "Il valore della stringa non deve essere una stringa vuota",
			"de": "Stringwert darf kein leerer String sein",
		},
		msgNotJsonArray: {
			"en": msgNotJsonArray,
			"fr": "JSON ne doit pas être un tableau JSON",
			"es": "JSON no debe ser una matriz JSON",
			"it": "JSON non deve essere un array JSON",
			"de": "JSON darf kein JSON-Array sein",
		},
		msgNotJsonNull: {
			"en": msgNotJsonNull,
			"fr": "JSON ne doit pas être JSON null",
			"es": "JSON no debe ser JSON nulo",
			"it": "JSON non deve essere JSON null",
			"de": "JSON darf nicht JSON null sein",
		},
		msgNotJsonObject: {
			"en": msgNotJsonObject,
			"fr": "JSON ne doit pas être un objet JSON",
			"es": "JSON no debe ser un objeto JSON",
			"it": "JSON non deve essere un oggetto JSON",
			"de": "JSON darf kein JSON-Objekt sein",
		},
		msgPositive: {
			"en": msgPositive,
			"fr": "La valeur doit être positive",
			"es": "El valor debe ser positivo",
			"it": "Il valore deve essere positivo",
			"de": "Wert muss positiv sein",
		},
		msgPositiveOrZero: {
			"en": msgPositiveOrZero,
			"fr": "La valeur doit être positive ou nulle",
			"es": "El valor debe ser positivo o cero",
			"it": "Il valore deve essere positivo o zero",
			"de": "Wert muss positiv oder Null sein",
		},
		msgPropertyObjectValidatorError: {
			"en": msgPropertyObjectValidatorError,
			"fr": "Erreur du validateur d'objet - n'autorise pas l'objet ou le tableau!",
			"es": "Error del validador de objetos: ¡no permite el objeto o la matriz!",
			"it": "Errore del validatore di oggetti - non consente l'oggetto o l'array!",
			"de": "Objekt-Validator-Fehler - erlaubt kein Objekt oder Array!",
		},
		msgPropertyValueMustBeObject: {
			"en": msgPropertyValueMustBeObject,
			"fr": "La valeur de la propriété doit être un objet",
			"es": "El valor de la propiedad debe ser un objeto",
			"it": "Il valore della proprietà deve essere un oggetto",
			"de": "Eigenschaftswert muss ein Objekt sein",
		},
		msgPropertyRequiredWhen: {
			"en": msgPropertyRequiredWhen,
			"fr": "La propriété est requise selon certains critères",
			"es": "Se requiere propiedad bajo ciertos criterios",
			"it": "L'immobile è richiesto secondo determinati criteri",
			"de": "Eigentum wird unter bestimmten Kriterien benötigt",
		},
		msgPropertyUnwantedWhen: {
			"en": msgPropertyUnwantedWhen,
			"fr": "La propriété ne doit pas être présente dans certaines conditions",
			"es": "La propiedad no debe estar presente bajo ciertas condiciones",
			"it": "L'immobile non deve essere presente in determinate condizioni",
			"de": "Die Immobilie darf unter bestimmten Voraussetzungen nicht vorhanden sein",
		},
		msgRequestBodyEmpty: {
			"en": msgRequestBodyEmpty,
			"fr": "Le corps de la requête est vide",
			"es": "El cuerpo de la solicitud está vacío",
			"it": "Il corpo della richiesta è vuoto",
			"de": "Anfragetext ist leer",
		},
		msgRequestBodyExpectedJsonArray: {
			"en": msgRequestBodyExpectedJsonArray,
			"fr": "Le corps de la requête devrait être un tableau JSON",
			"es": "Se espera que el cuerpo de la solicitud sea una matriz JSON",
			"it": "Il corpo della richiesta dovrebbe essere un array JSON",
			"de": "Anforderungstext soll JSON-Array sein",
		},
		msgRequestBodyExpectedJsonObject: {
			"en": msgRequestBodyExpectedJsonObject,
			"fr": "Le corps de la requête devrait être un objet JSON",
			"es": "Se espera que el cuerpo de la solicitud sea un objeto JSON",
			"it": "Il corpo della richiesta dovrebbe essere un oggetto JSON",
			"de": "Anfragetext soll JSON-Objekt sein",
		},
		msgRequestBodyNotJsonArray: {
			"en": msgRequestBodyNotJsonArray,
			"fr": "Le corps de la requête ne doit pas être un tableau JSON",
			"es": "El cuerpo de la solicitud no debe ser una matriz JSON",
			"it": "Il corpo della richiesta non deve essere un array JSON",
			"de": "Anfragetext darf kein JSON-Array sein",
		},
		msgRequestBodyNotJsonNull: {
			"en": msgRequestBodyNotJsonNull,
			"fr": "Le corps de la requête ne doit pas être nul en JSON",
			"es": "El cuerpo de la solicitud no debe ser JSON nulo",
			"it": "Il corpo della richiesta non deve essere JSON null",
			"de": "Anfragetext darf nicht JSON null sein",
		},
		msgRequestBodyNotJsonObject: {
			"en": msgRequestBodyNotJsonObject,
			"fr": "Le corps de la requête ne doit pas être un objet JSON",
			"es": "El cuerpo de la solicitud no debe ser un objeto JSON",
			"it": "Il corpo della richiesta non deve essere un oggetto JSON",
			"de": "Anfragetext darf kein JSON-Objekt sein",
		},
		msgStringLowercase: {
			"en": msgStringLowercase,
			"fr": "La valeur de la chaîne ne doit contenir que des lettres minuscules",
			"es": "El valor de la cadena debe contener solo letras minúsculas",
			"it": "Il valore della stringa deve contenere solo lettere minuscole",
			"de": "Stringwert darf nur Kleinbuchstaben enthalten",
		},
		msgStringUppercase: {
			"en": msgStringUppercase,
			"fr": "La valeur de la chaîne ne doit contenir que des lettres majuscules",
			"es": "El valor de la cadena debe contener solo letras mayúsculas",
			"it": "Il valore della stringa deve contenere solo lettere maiuscole",
			"de": "String-Wert darf nur Großbuchstaben enthalten",
		},
		msgUnableToDecode: {
			"en": msgUnableToDecode,
			"fr": "Impossible de décoder en JSON",
			"es": "No se puede decodificar como JSON",
			"it": "Impossibile decodificare come JSON",
			"de": "Als JSON kann nicht dekodiert werden",
		},
		msgUnableToDecodeRequest: {
			"en": msgUnableToDecodeRequest,
			"fr": "Impossible de décoder le corps de la requête en JSON",
			"es": "No se puede decodificar el cuerpo de la solicitud como JSON",
			"it": "Impossibile decodificare il corpo della richiesta come JSON",
			"de": "Anforderungstext konnte nicht als JSON entschlüsselt werden",
		},
		msgUnicodeNormalization: {
			"en": msgUnicodeNormalization,
			"fr": "La valeur de la chaîne doit être une forme de normalisation correcte",
			"es": "El valor de la cadena debe ser la forma de normalización correcta",
			"it": "Il valore della stringa deve essere un modulo di normalizzazione corretto",
			"de": "String-Wert muss korrekte Normalisierungsform sein",
		},
		msgUnicodeNormalizationNFC: {
			"en": msgUnicodeNormalizationNFC,
			"fr": "La valeur de la chaîne doit être la forme de normalisation correcte NFC",
			"es": "El valor de la cadena debe ser la normalización correcta de NFC",
			"it": "Il valore della stringa deve essere la normalizzazione corretta da NFC",
			"de": "Stringwert muss korrekte Normalisierung von NFC sein",
		},
		msgUnicodeNormalizationNFD: {
			"en": msgUnicodeNormalizationNFD,
			"fr": "La valeur de la chaîne doit être la forme de normalisation correcte NFD",
			"es": "El valor de la cadena debe ser el formulario de normalización correcto NFD",
			"it": "Il valore della stringa deve essere la normalizzazione corretta dal modulo NFD",
			"de": "String-Wert muss korrekte Normalisierung von NFD sein",
		},
		msgUnicodeNormalizationNFKC: {
			"en": msgUnicodeNormalizationNFKC,
			"fr": "La valeur de la chaîne doit être la forme de normalisation correcte NFKC",
			"es": "El valor de la cadena debe ser el formulario de normalización correcto NFKC",
			"it": "Il valore della stringa deve essere la normalizzazione corretta da NFKC",
			"de": "String-Wert muss korrekte Normalisierung von NFKC sein",
		},
		msgUnicodeNormalizationNFKD: {
			"en": msgUnicodeNormalizationNFKD,
			"fr": "La valeur de la chaîne doit être la forme de normalisation correcte NFKD",
			"es": "El valor de la cadena debe ser el formulario de normalización correcto NFKD",
			"it": "Il valore della stringa deve essere la normalizzazione corretta da NFKD",
			"de": "String-Wert muss korrekte Normalisierung von NFKD sein",
		},
		msgUnknownProperty: {
			"en": msgUnknownProperty,
			"fr": "Propriété inconnue",
			"es": "Propiedad desconocida",
			"it": "Proprietà sconosciuta",
			"de": "Unbekanntes Eigentum",
		},
		msgUnwantedProperty: {
			"en": msgUnwantedProperty,
			"fr": "La propriété ne doit pas être présente",
			"es": "La propiedad no debe estar presente",
			"it": "L'immobile non deve essere presente",
			"de": "Eigenschaft darf nicht vorhanden sein",
		},
		msgValueCannotBeNull: {
			"en": msgValueCannotBeNull,
			"fr": "La valeur ne peut pas être nulle",
			"es": "El valor no puede ser nulo",
			"it": "Il valore non può essere nullo",
			"de": "Wert darf nicht null sein",
		},
		msgValueMustBeArray: {
			"en": msgValueMustBeArray,
			"fr": "La valeur doit être un tableau",
			"es": "El valor debe ser una matriz",
			"it": "Il valore deve essere un array",
			"de": "Wert muss ein Array sein",
		},
		msgValueMustBeObject: {
			"en": msgValueMustBeObject,
			"fr": "La valeur doit être un objet",
			"es": "El valor debe ser un objeto",
			"it": "Il valore deve essere un oggetto",
			"de": "Wert muss ein Objekt sein",
		},
		msgValueMustBeObjectOrArray: {
			"en": msgValueMustBeObjectOrArray,
			"fr": "La valeur doit être un objet ou un tableau",
			"es": "El valor debe ser un objeto o matriz",
			"it": "Il valore deve essere un oggetto o un array",
			"de": "Wert muss ein Objekt oder Array sein",
		},
		msgValidCardNumber: {
			"en": msgValidCardNumber,
			"fr": "La valeur doit être un numéro de carte valide",
			"es": "El valor debe ser un número de tarjeta válido",
			"it": "Il valore deve essere un numero di carta valido",
			"de": "Wert muss eine gültige Kartennummer sein",
		},
		msgValidEmail: {
			"en": msgValidEmail,
			"fr": "La valeur doit être une adresse e-mail",
			"es": "El valor debe ser una dirección de correo electrónico",
			"it": "Il valore deve essere un indirizzo email",
			"de": "Wert muss eine E-Mail-Adresse sein",
		},
		msgValidISODate: {
			"en": msgValidISODate,
			"fr": "La valeur doit être une chaîne de date valide (format : AAAA-MM-JJ)",
			"es": "El valor debe ser una cadena de fecha válida (formato: AAAA-MM-DD)",
			"it": "Il valore deve essere una stringa di data valida (formato: AAAA-MM-GG)",
			"de": "Wert muss eine gültige Datumszeichenfolge sein (Format: JJJJ-MM-TT)",
		},
		msgValidISODatetimeFormatFull: {
			"en": msgValidISODatetimeFormatFull,
			"fr": "La valeur doit être une chaîne de date/heure valide (format : AAAA-MM-JJThh : mm:ss.sss [Z|+-hh:mm])",
			"es": "El valor debe ser una cadena de fecha/hora válida (formato: AAAA-MM-DDThh: mm:ss.sss [Z|+- hh:mm ])",
			"it": "Il valore deve essere una stringa di data/ora valida (formato: AAAA-MM-GGThh: mm:ss.sss [Z|+- hh:mm ])",
			"de": "Wert muss ein gültiger Datums-/Uhrzeit-String sein (Format: YYYY-MM-DDThh: mm:ss.sss [Z|+- hh:mm ])",
		},
		msgValidISODatetimeFormatMin: {
			"en": msgValidISODatetimeFormatMin,
			"fr": "La valeur doit être une chaîne date/heure valide (format : AAAA-MM-JJThh: mm:ss)",
			"es": "El valor debe ser una cadena de fecha/hora válida (formato: AAAA-MM-DDThh: mm:ss)",
			"it": "Il valore deve essere una stringa di data/ora valida (formato: AAAA-MM-GGThh: mm:ss)",
			"de": "Wert muss ein gültiger Datums-/Uhrzeit-String sein (Format: YYYY-MM-DDThh: mm:ss )",
		},
		msgValidISODatetimeFormatNoOffs: {
			"en": msgValidISODatetimeFormatNoOffs,
			"fr": "La valeur doit être une chaîne date/heure valide (format : AAAA-MM-JJThh: mm:ss.sss)",
			"es": "El valor debe ser una cadena de fecha/hora válida (formato: AAAA-MM-DDThh: mm:ss.sss)",
			"it": "Il valore deve essere una stringa di data/ora valida (formato: AAAA-MM-GGThh: mm:ss.sss)",
			"de": "Wert muss ein gültiger Datums-/Uhrzeit-String sein (Format: YYYY-MM-DDThh: mm:ss.sss )",
		},
		msgValidISODatetimeFormatNoMillis: {
			"en": msgValidISODatetimeFormatNoMillis,
			"fr": "La valeur doit être une chaîne de date/heure valide (format : AAAA-MM-JJThh : mm:ss [Z|+-hh:mm])",
			"es": "El valor debe ser una cadena de fecha/hora válida (formato: AAAA-MM-DDThh: mm:ss [Z|+- hh:mm ])",
			"it": "Il valore deve essere una stringa di data/ora valida (formato: AAAA-MM-GGThh: mm:ss [Z|+- hh:mm ])",
			"de": "Wert muss ein gültiger Datums-/Uhrzeit-String sein (Format: YYYY-MM-DDThh: mm:ss [Z|+- hh:mm ])",
		},
		msgValidPattern: {
			"en": msgValidPattern,
			"fr": "La valeur de la chaîne doit avoir un modèle valide",
			"es": "El valor de la cadena debe tener un patrón válido",
			"it": "Il valore della stringa deve avere un modello valido",
			"de": "String-Wert muss gültiges Muster haben",
		},
		msgValidUuid: {
			"en": msgValidUuid,
			"fr": "La valeur doit être un UUID valide",
			"es": "El valor debe ser un UUID válido",
			"it": "Il valore deve essere un UUID valido",
			"de": "Wert muss eine gültige UUID sein",
		},
		msgPresetAlpha: {
			"en": msgPresetAlpha,
			"fr": "La valeur ne doit être que des caractères alphabétiques (A-Z, a-z)",
			"es": "El valor debe ser solo caracteres alfabéticos (A-Z, a-z)",
			"it": "Il valore deve essere solo caratteri alfabetici (A-Z, a-z)",
			"de": "Wert darf nur aus Buchstaben bestehen (A-Z, a-z)",
		},
		msgPresetAlphaNumeric: {
			"en": msgPresetAlphaNumeric,
			"fr": "La valeur ne doit être que des caractères alphanumériques (A-Z, a-z, 0-9)",
			"es": "El valor debe ser solo caracteres alfanuméricos (A-Z, a-z, 0-9)",
			"it": "Il valore deve essere solo caratteri alfanumerici (A-Z, a-z, 0-9)",
			"de": "Wert darf nur aus alphanumerischen Zeichen bestehen (A-Z, a-z, 0-9)",
		},
		msgPresetBarcode: {
			"en": "Value must be a valid barcode",
			"fr": "La valeur doit être un code-barres valide",
			"es": "El valor debe ser un código de barras válido",
			"it": "Il valore deve essere un codice a barre valido",
			"de": "Wert muss ein gültiger Strichcode sein",
		},
		msgPresetBase64: {
			"en": msgPresetBase64,
			"fr": "La valeur doit être une chaîne valide encodée en base64",
			"es": "El valor debe ser una cadena codificada en base64 válida",
			"it": "Il valore deve essere una stringa codificata base64 valida",
			"de": "Wert muss eine gültige base64-codierte Zeichenfolge sein",
		},
		msgPresetBase64URL: {
			"en": msgPresetBase64URL,
			"fr": "La valeur doit être une chaîne encodée URL base64 valide",
			"es": "El valor debe ser una cadena codificada en URL base64 válida",
			"it": "Il valore deve essere una stringa codificata URL base64 valida",
			"de": "Wert muss eine gültige Base64-URL-codierte Zeichenfolge sein",
		},
		msgPresetCMYK: {
			"en":    msgPresetCMYK,
			"en-US": "Value must be a valid cmyk() color string",
			"fr":    "La valeur doit être une chaîne de couleur cmyk() valide",
			"es":    "El valor debe ser una cadena de color cmyk() válida",
			"it":    "Il valore deve essere una stringa di colore cmyk() valida",
			"de":    "Der Wert muss eine gültige cmyk()-Farbzeichenfolge sein",
		},
		msgPresetCMYK300: {
			"en":    msgPresetCMYK300,
			"en-US": "Value must be a valid cmyk() color string (maximum 300%)",
			"fr":    "La valeur doit être une chaîne de couleur cmyk() valide (maximum 300 %)",
			"es":    "El valor debe ser una cadena de color cmyk() válida (máximo 300 %)",
			"it":    "Il valore deve essere una stringa di colore cmyk() valida (massimo 300%)",
			"de":    "Der Wert muss eine gültige cmyk()-Farbzeichenfolge sein (maximal 300 %)",
		},
		msgPresetE164: {
			"en": msgPresetE164,
			"fr": "La valeur doit être un code E.164 valide",
			"es": "El valor debe ser un código E.164 válido",
			"it": "Il valore deve essere un codice E.164 valido",
			"de": "Wert muss ein gültiger E.164-Code sein",
		},
		msgPresetEAN: {
			"en": msgPresetEAN,
			"fr": "La valeur doit être un code EAN valide",
			"es": "El valor debe ser un código EAN válido",
			"it": "Il valore deve essere un codice EAN valido",
			"de": "Wert muss ein gültiger EAN-Code sein",
		},
		msgPresetEAN8: {
			"en": msgPresetEAN8,
			"fr": "La valeur doit être un code EAN-8 valide",
			"es": "El valor debe ser un código EAN-8 válido",
			"it": "Il valore deve essere un codice EAN-8 valido",
			"de": "Wert muss ein gültiger EAN-8-Code sein",
		},
		msgPresetEAN13: {
			"en": msgPresetEAN13,
			"fr": "La valeur doit être un code EAN-13 valide",
			"es": "El valor debe ser un código EAN-13 válido",
			"it": "Il valore deve essere un codice EAN-13 valido",
			"de": "Wert muss ein gültiger EAN-13-Code sein",
		},
		msgPresetDUN14: {
			"en": msgPresetDUN14,
			"fr": "La valeur doit être un code DUN-14 valide",
			"es": "El valor debe ser un código DUN-14 válido",
			"it": "Il valore deve essere un codice DUN-14 valido",
			"de": "Wert muss ein gültiger DUN-14-Code sein",
		},
		msgPresetEAN14: {
			"en": msgPresetEAN14,
			"fr": "La valeur doit être un code EAN-14 valide",
			"es": "El valor debe ser un código EAN-14 válido",
			"it": "Il valore deve essere un codice EAN-14 valido",
			"de": "Wert muss ein gültiger EAN-14-Code sein",
		},
		msgPresetEAN18: {
			"en": msgPresetEAN18,
			"fr": "La valeur doit être un code EAN-18 valide",
			"es": "El valor debe ser un código EAN-18 válido",
			"it": "Il valore deve essere un codice EAN-18 valido",
			"de": "Wert muss ein gültiger EAN-18-Code sein",
		},
		msgPresetEAN99: {
			"en": msgPresetEAN99,
			"fr": "La valeur doit être un code EAN-99 valide",
			"es": "El valor debe ser un código EAN-99 válido",
			"it": "Il valore deve essere un codice EAN-99 valido",
			"de": "Wert muss ein gültiger EAN-99-Code sein",
		},
		msgPresetHexadecimal: {
			"en": msgPresetHexadecimal,
			"fr": "La valeur doit être une chaîne hexadécimale valide",
			"es": "El valor debe ser una cadena hexadecimal válida",
			"it": "Il valore deve essere una stringa esadecimale valida",
			"de": "Wert muss eine gültige hexadezimale Zeichenfolge sein",
		},
		msgPresetHsl: {
			"en":    msgPresetHsl,
			"en-US": "Value must be a valid hsl() color string",
			"fr":    "La valeur doit être une chaîne de couleur hsl() valide",
			"es":    "El valor debe ser una cadena de color hsl() válida",
			"it":    "Il valore deve essere una stringa di colore hsl() valida",
			"de":    "Wert muss eine gültige hsl() Farbzeichenfolge sein",
		},
		msgPresetHsla: {
			"en":    msgPresetHsla,
			"en-US": "Value must be a valid hsla() color string",
			"fr":    "La valeur doit être une chaîne de couleur hsla() valide",
			"es":    "El valor debe ser una cadena de color hsla() válida",
			"it":    "Il valore deve essere una stringa di colore hsla() valida",
			"de":    "Wert muss eine gültige hsla() Farbzeichenfolge sein",
		},
		msgPresetHtmlColor: {
			"en":    msgPresetHtmlColor,
			"en-US": "Value must be a valid HTML color string",
			"fr":    "La valeur doit être une chaîne de couleur HTML valide",
			"es":    "El valor debe ser una cadena de color HTML válida",
			"it": "Il valore deve essere una 	stringa di colori HTML valida",
			"de": "Wert muss ein gültiger HTML -Farbstring sein",
		},
		msgPresetInteger: {
			"en": msgPresetInteger,
			"fr": "La valeur doit être une chaîne d'entiers valide (caractères 0-9)",
			"es": "El valor debe ser una cadena entera válida (caracteres 0-9)",
			"it": "Il valore deve essere una stringa intera valida (caratteri 0-9)",
			"de": "Wert muss eine gültige Ganzzahl sein (Zeichen 0-9)",
		},
		msgPresetISBN: {
			"en": msgPresetISBN,
			"fr": "La valeur doit être un ISBN valide",
			"es": "El valor debe ser un ISBN válido",
			"it": "Il valore deve essere un ISBN valido",
			"de": "Wert muss eine gültige ISBN sein",
		},
		msgPresetISBN10: {
			"en": msgPresetISBN10,
			"fr": "La valeur doit être un ISBN-10 valide",
			"es": "El valor debe ser un ISBN-10 válido",
			"it": "Il valore deve essere un ISBN-10 valido",
			"de": "Wert muss eine gültige ISBN-10 sein",
		},
		msgPresetISBN13: {
			"en": msgPresetISBN13,
			"fr": "La valeur doit être un ISBN-13 valide",
			"es": "El valor debe ser un ISBN-13 válido",
			"it": "Il valore deve essere un ISBN-13 valido",
			"de": "Wert muss eine gültige ISBN-13 sein",
		},
		msgPresetISSN: {
			"en": msgPresetISSN,
			"fr": "La valeur doit être un ISSN valide",
			"es": "El valor debe ser un ISSN válido",
			"it": "Il valore deve essere un ISSN valido",
			"de": "Wert muss eine gültige ISSN sein",
		},
		msgPresetNumeric: {
			"en": msgPresetNumeric,
			"fr": "La valeur doit être une chaîne numérique valide",
			"es": "El valor debe ser una cadena de números válida",
			"it": "Il valore deve essere una stringa numerica valida",
			"de": "Wert muss eine gültige Zahlenfolge sein",
		},
		msgPresetPublication: {
			"en": msgPresetPublication,
			"fr": "La valeur doit être un ISBN ou un ISSN valide",
			"es": "El valor debe ser un ISBN o ISSN válido",
			"it": "Il valore deve essere un ISBN o ISSN valido",
			"de": "Wert muss eine gültige ISBN oder ISSN sein",
		},
		msgPresetRgb: {
			"en":    msgPresetRgb,
			"en-US": "Value must be a valid rgb() color string",
			"fr":    "La valeur doit être une chaîne de couleur rgb() valide",
			"es":    "El valor debe ser una cadena de color rgb() válida",
			"it":    "Il valore deve essere una stringa di colore rgb() valida",
			"de":    "Wert muss eine gültige rgb() Farbzeichenfolge sein",
		},
		msgPresetRgba: {
			"en":    msgPresetRgba,
			"en-US": "Value must be a valid rgba() color string",
			"fr":    "La valeur doit être une chaîne de couleur rgba() valide",
			"es":    "El valor debe ser una cadena de color rgba() válida",
			"it":    "Il valore deve essere una stringa di colore rgba() valida",
			"de":    "Wert muss eine gültige rgba() Farbzeichenfolge sein",
		},
		msgPresetRgbIcc: {
			"en":    msgPresetRgbIcc,
			"en-US": "Value must be a valid rgb-icc() color string",
			"fr":    "La valeur doit être une chaîne de couleur rgb-icc() valide",
			"es":    "El valor debe ser una cadena de color rgb-icc() válida",
			"it":    "Il valore deve essere una stringa di colore rgb-icc() valida",
			"de":    "Wert muss eine gültige rgb-icc() Farbzeichenfolge sein",
		},
		msgPresetULID: {
			"en": msgPresetULID,
			"fr": "La valeur doit être un ULID valide",
			"es": "El valor debe ser un ULID válido",
			"it": "Il valore deve essere un ULID valido",
			"de": "Wert muss eine gültige ULID sein",
		},
		msgPresetUPC: {
			"en": msgPresetUPC,
			"fr": "La valeur doit être un code UPC valide (UPC-A ou UPC-E)",
			"es": "El valor debe ser un código UPC válido (UPC-A o UPC-E)",
			"it": "Il valore deve essere un codice UPC valido (UPC-A o UPC-E)",
			"de": "Wert muss ein gültiger UPC-Code sein (UPC-A oder UPC-E)",
		},
		msgPresetUPCA: {
			"en": msgPresetUPCA,
			"fr": "La valeur doit être un code UPC-A valide",
			"es": "El valor debe ser un código UPC-A válido",
			"it": "Il valore deve essere un codice UPC-A valido",
			"de": "Wert muss ein gültiger UPC-A-Code sein",
		},
		msgPresetUPCE: {
			"en": msgPresetUPCE,
			"fr": "La valeur doit être un code UPC-E valide",
			"es": "El valor debe ser un código UPC-E válido",
			"it": "Il valore deve essere un codice UPC-E valido",
			"de": "Wert muss ein gültiger UPC-E-Code sein",
		},
		msgPresetUuid1: {
			"en": msgPresetUuid1,
			"fr": "La valeur doit être un UUID valide (Version 1)",
			"es": "El valor debe ser un UUID válido (Versión 1)",
			"it": "Il valore deve essere un UUID valido (versione 1)",
			"de": "Wert muss eine gültige UUID sein (Version 1)",
		},
		msgPresetUuid2: {
			"en": msgPresetUuid2,
			"fr": "La valeur doit être un UUID valide (Version 2)",
			"es": "El valor debe ser un UUID válido (Versión 2)",
			"it": "Il valore deve essere un UUID valido (versione 2)",
			"de": "Wert muss eine gültige UUID sein (Version 2)",
		},
		msgPresetUuid3: {
			"en": msgPresetUuid3,
			"fr": "La valeur doit être un UUID valide (Version 3)",
			"es": "El valor debe ser un UUID válido (Versión 3)",
			"it": "Il valore deve essere un UUID valido (versione 3)",
			"de": "Wert muss eine gültige UUID sein (Version 3)",
		},
		msgPresetUuid4: {
			"en": msgPresetUuid4,
			"fr": "La valeur doit être un UUID valide (Version 4)",
			"es": "El valor debe ser un UUID válido (Versión 4)",
			"it": "Il valore deve essere un UUID valido (versione 4)",
			"de": "Wert muss eine gültige UUID sein (Version 4)",
		},
		msgPresetUuid5: {
			"en": msgPresetUuid5,
			"fr": "La valeur doit être un UUID valide (Version 5)",
			"es": "El valor debe ser un UUID válido (Versión 5)",
			"it": "Il valore deve essere un UUID valido (versione 5)",
			"de": "Wert muss eine gültige UUID sein (Version 5)",
		},
	},
	Formats: map[string]map[string]string{
		fmtMsgArrayElementType: {
			"en": fmtMsgArrayElementType,
			"fr": "Les éléments du tableau doivent être de type %[1]s",
			"es": "Los elementos del arreglo deben ser del tipo %[1]s",
			"it": "Gli elementi dell'array devono essere di tipo %[1]s",
			"de": "Array-Elemente müssen vom Typ %[1]s sein",
		},
		fmtMsgArrayElementTypeOrNull: {
			"en": fmtMsgArrayElementTypeOrNull,
			"fr": "Les éléments du tableau doivent être de type %[1]s ou nuls",
			"es": "Los elementos del arreglo deben ser del tipo %[1]s o nulo",
			"it": "Gli elementi dell'array devono essere di tipo %[1]s o null",
			"de": "Array-Elemente müssen vom Typ %[1]s oder null sein",
		},
		fmtMsgConstraintSetDefaultAllOf: {
			"en": fmtMsgConstraintSetDefaultAllOf,
			"fr": "L'ensemble de contraintes doit réussir toutes les %[1]d validations non divulguées",
			"es": "El conjunto de restricciones debe pasar todas las %[1]d validaciones no reveladas",
			"it": "Il set di vincoli deve superare tutte le %[1]d convalide non divulgate",
			"de": "Einschränkungssatz muss alle %[1]d nicht offengelegten Validierungen bestehen",
		},
		fmtMsgConstraintSetDefaultOneOf: {
			"en": fmtMsgConstraintSetDefaultOneOf,
			"fr": "L'ensemble de contraintes doit réussir l'une des %[1]d validations non divulguées",
			"es": "El conjunto de restricciones debe pasar una de %[1]d validaciones no reveladas",
			"it": "Il set di vincoli deve superare una delle %[1]d convalide non divulgate",
			"de": "Einschränkungssatz muss eine von %[1]d nicht offengelegten Validierungen bestehen",
		},
		fmtMsgDtGt: {
			"en": fmtMsgDtGt,
			"fr": "La valeur doit être après '%[1]s'",
			"es": "El valor debe estar después de '%[1]s'",
			"it": "Il valore deve essere successivo a '%[1]s'",
			"de": "Wert muss nach '%[1]s' liegen",
		},
		fmtMsgDtGte: {
			"en": fmtMsgDtGte,
			"fr": "La valeur doit être supérieure ou égale à '%[1]s'",
			"es": "El valor debe ser posterior o igual a '%[1]s'",
			"it": "Il valore deve essere successivo o uguale a '%[1]s'",
			"de": "Wert muss nach oder gleich '%[1]s' sein",
		},
		fmtMsgDtLt: {
			"en": fmtMsgDtLt,
			"fr": "La valeur doit être avant '%[1]s'",
			"es": "El valor debe estar antes de '%[1]s'",
			"it": "Il valore deve essere prima di '%[1]s'",
			"de": "Wert muss vor '%[1]s' liegen",
		},
		fmtMsgDtLte: {
			"en": fmtMsgDtLte,
			"fr": "La valeur doit être avant ou égale à '%[1]s'",
			"es": "El valor debe ser anterior o igual a '%[1]s'",
			"it": "Il valore deve essere prima o uguale a '%[1]s'",
			"de": "Wert muss vor oder gleich '%[1]s' sein",
		},
		fmtMsgDtToleranceFixedNotMoreThanAfterPlural: {
			"en": fmtMsgDtToleranceFixedNotMoreThanAfterPlural,
			"fr": "La valeur ne doit pas être supérieure à %[1]d %[2]s après %[3]s",
			"es": "El valor no debe ser mayor que %[1]d %[2]s después de %[3]s",
			"it": "Il valore non deve essere superiore a %[1]d %[2]s dopo %[3]s",
			"de": "Wert darf nach %[3]s nicht größer als %[1]d %[2]s sein",
		},
		fmtMsgDtToleranceFixedNotMoreThanAfterSingular: {
			"en": fmtMsgDtToleranceFixedNotMoreThanAfterSingular,
			"fr": "La valeur ne doit pas être supérieure à %[1]s après %[2]s",
			"es": "El valor no debe ser mayor que %[1]s después de %[2]s",
			"it": "Il valore non deve essere superiore a %[1]s dopo %[2]s",
			"de": "Wert darf nach %[2]s nicht größer als %[1]s sein",
		},
		fmtMsgDtToleranceFixedNotMoreThanBeforePlural: {
			"en": fmtMsgDtToleranceFixedNotMoreThanBeforePlural,
			"fr": "La valeur ne doit pas être supérieure à %[1]d %[2]s avant %[3]s",
			"es": "El valor no debe ser mayor que %[1]d %[2]s antes de %[3]s",
			"it": "Il valore non deve essere superiore a %[1]d %[2]s prima di %[3]s",
			"de": "Wert darf nicht größer als %[1]d %[2]s vor %[3]s sein",
		},
		fmtMsgDtToleranceFixedNotMoreThanBeforeSingular: {
			"en": fmtMsgDtToleranceFixedNotMoreThanBeforeSingular,
			"fr": "La valeur ne doit pas être supérieure à %[1]s avant %[2]s",
			"es": "El valor no debe ser mayor que %[1]s antes de %[2]s",
			"it": "Il valore non deve essere superiore a %[1]s prima di %[2]s",
			"de": "Wert darf nicht länger als %[1]s vor %[2]s sein",
		},
		fmtMsgDtToleranceFixedSame: {
			"en": fmtMsgDtToleranceFixedSame,
			"fr": "La valeur doit être la même %[1]s que %[2]s",
			"es": "El valor debe ser el mismo %[1]s que %[2]s",
			"it": "Il valore deve essere lo stesso %[1]s di %[2]s",
			"de": "Wert muss gleich %[1]s wie %[2]s sein",
		},
		fmtMsgDtToleranceNowSame: {
			"en": fmtMsgDtToleranceNowSame,
			"fr": "La valeur doit être la même %[1]s qu'actuellement",
			"es": "El valor debe ser el mismo %[1]s que ahora",
			"it": "Il valore deve essere lo stesso %[1]s di adesso",
			"de": "Wert muss gleich %[1]s sein wie jetzt",
		},
		fmtMsgDtToleranceNowNotMoreThanAfterPlural: {
			"en": fmtMsgDtToleranceNowNotMoreThanAfterPlural,
			"fr": "La valeur ne doit pas dépasser %[1]d %[2]s après maintenant",
			"es": "El valor no debe ser mayor que %[1]d %[2]s después de ahora",
			"it": "Il valore non deve essere superiore a %[1]d %[2]s dopo ora",
			"de": "Wert darf nach jetzt nicht mehr als %[1]d %[2]s betragen",
		},
		fmtMsgDtToleranceNowNotMoreThanAfterSingular: {
			"en": fmtMsgDtToleranceNowNotMoreThanAfterSingular,
			"fr": "La valeur ne doit pas être supérieure à %[1]s après maintenant",
			"es": "El valor no debe ser mayor a %[1]s después de ahora",
			"it": "Il valore non deve essere superiore a %[1]s dopo ora",
			"de": "Wert darf nicht mehr als %[1]s nach jetzt sein",
		},
		fmtMsgDtToleranceNowNotMoreThanBeforePlural: {
			"en": fmtMsgDtToleranceNowNotMoreThanBeforePlural,
			"fr": "La valeur ne doit pas dépasser %[1]d %[2]s avant maintenant",
			"es": "El valor no debe ser superior a %[1]d %[2]s antes de ahora",
			"it": "Il valore non deve essere superiore a %[1]d %[2]s prima di ora",
			"de": "Wert darf bis jetzt nicht größer als %[1]d %[2]s sein",
		},
		fmtMsgDtToleranceNowNotMoreThanBeforeSingular: {
			"en": fmtMsgDtToleranceNowNotMoreThanBeforeSingular,
			"fr": "La valeur ne doit pas dépasser un %[1]s avant maintenant",
			"es": "El valor no debe ser superior a %[1]s antes de ahora",
			"it": "Il valore non deve essere superiore a %[1]s prima di ora",
			"de": "Wert darf bis jetzt nicht größer als %[1]s sein",
		},
		fmtMsgDtToleranceOtherSame: {
			"en": fmtMsgDtToleranceOtherSame,
			"fr": "La valeur doit être la même %[1]s que la valeur de la propriété '%[2]s'",
			"es": "El valor debe ser el mismo %[1]s que el valor de la propiedad '%[2]s'",
			"it": "Il valore deve essere lo stesso %[1]s del valore della proprietà '%[2]s'",
			"de": "Wert muss gleich %[1]s sein wie Wert der Eigenschaft '%[2]s'",
		},
		fmtMsgDtToleranceOtherNotMoreThanAfterSingular: {
			"en": fmtMsgDtToleranceOtherNotMoreThanAfterSingular,
			"fr": "La valeur ne doit pas être supérieure à %[1]s après la valeur de la propriété '%[2]s'",
			"es": "El valor no debe ser mayor que %[1]s después del valor de la propiedad '%[2]s'",
			"it": "Il valore non deve essere superiore a %[1]s dopo il valore della proprietà '%[2]s'",
			"de": "Der Wert darf nicht mehr als ein %[1]s nach dem Wert der Eigenschaft '%[2]s' liegen",
		},
		fmtMsgDtToleranceOtherNotMoreThanAfterPlural: {
			"en": fmtMsgDtToleranceOtherNotMoreThanAfterPlural,
			"fr": "La valeur ne doit pas être supérieure à %[1]d %[2]s après la valeur de la propriété '%[3]s'",
			"es": "El valor no debe ser mayor que %[1]d %[2]s después del valor de la propiedad '%[3]s'",
			"it": "Il valore non deve essere superiore a %[1]d %[2]s dopo il valore della proprietà '%[3]s'",
			"de": "Der Wert darf nicht größer als %[1]d %[2]s nach dem Wert der Eigenschaft '%[3]s' sein",
		},
		fmtMsgDtToleranceOtherNotMoreThanBeforePlural: {
			"en": fmtMsgDtToleranceOtherNotMoreThanBeforePlural,
			"fr": "La valeur ne doit pas être supérieure à %[1]d %[2]s avant la valeur de la propriété '%[3]s'",
			"es": "El valor no debe ser mayor que %[1]d %[2]s antes del valor de la propiedad '%[3]s'",
			"it": "Il valore non deve essere superiore a %[1]d %[2]s prima del valore della proprietà '%[3]s'",
			"de": "Der Wert darf nicht größer als %[1]d %[2]s vor dem Wert der Eigenschaft '%[3]s' sein",
		},
		fmtMsgDtToleranceOtherNotMoreThanBeforeSingular: {
			"en": fmtMsgDtToleranceOtherNotMoreThanBeforeSingular,
			"fr": "La valeur ne doit pas être supérieure à %[1]s avant la valeur de la propriété '%[2]s'",
			"es": "El valor no debe ser mayor que %[1]s antes del valor de la propiedad '%[2]s'",
			"it": "Il valore non deve essere superiore a %[1]s prima del valore della proprietà '%[2]s'",
			"de": "Der Wert darf nicht mehr als ein %[1]s vor dem Wert der Eigenschaft '%[2]s' liegen",
		},
		fmtMsgEqualsOther: {
			"en": fmtMsgEqualsOther,
			"fr": "La valeur doit être égale à la valeur de la propriété '%[1]s'",
			"es": "El valor debe ser igual al valor de la propiedad '%[1]s'",
			"it": "Il valore deve essere uguale al valore della proprietà '%[1]s'",
			"de": "Wert muss gleich dem Wert der Eigenschaft '%[1]s' sein",
		},
		fmtMsgExactLen: {
			"en": fmtMsgExactLen,
			"fr": "La longueur de la valeur doit être %[1]d",
			"es": "La longitud del valor debe ser %[1]d",
			"it": "La lunghezza del valore deve essere %[1]d",
			"de": "Wertlänge muss %[1]d sein",
		},
		fmtMsgGt: {
			"en": fmtMsgGt,
			"fr": "La valeur doit être supérieure à %[1]v",
			"es": "El valor debe ser mayor que %[1]v",
			"it": "Il valore deve essere maggiore di %[1]v",
			"de": "Wert muss größer als %[1]v sein",
		},
		fmtMsgGte: {
			"en": fmtMsgGte,
			"fr": "La valeur doit être supérieure ou égale à %[1]v",
			"es": "El valor debe ser mayor o igual que %[1]v",
			"it": "Il valore deve essere maggiore o uguale a %[1]v",
			"de": "Wert muss größer oder gleich %[1]v sein",
		},
		fmtMsgGtOther: {
			"en": fmtMsgGtOther,
			"fr": "La valeur doit être supérieure à la valeur de la propriété '%[1]s'",
			"es": "El valor debe ser mayor que el valor de la propiedad '%[1]s'",
			"it": "Il valore deve essere maggiore del valore della proprietà '%[1]s'",
			"de": "Wert muss größer sein als Wert der Eigenschaft '%[1]s'",
		},
		fmtMsgGteOther: {
			"en": fmtMsgGteOther,
			"fr": "La valeur doit être supérieure ou égale à la valeur de la propriété '%[1]s'",
			"es": "El valor debe ser mayor o igual que el valor de la propiedad '%[1]s'",
			"it": "Il valore deve essere maggiore o uguale al valore della proprietà '%[1]s'",
			"de": "Wert muss größer oder gleich dem Wert der Eigenschaft '%[1]s' sein",
		},
		fmtMsgLt: {
			"en": fmtMsgLt,
			"fr": "La valeur doit être inférieure à %[1]v",
			"es": "El valor debe ser menor que %[1]v",
			"it": "Il valore deve essere inferiore a %[1]v",
			"de": "Wert muss kleiner als %[1]v sein",
		},
		fmtMsgLte: {
			"en": fmtMsgLte,
			"fr": "La valeur doit être inférieure ou égale à %[1]v",
			"es": "El valor debe ser menor o igual que %[1]v",
			"it": "Il valore deve essere inferiore o uguale a %[1]v",
			"de": "Wert muss kleiner oder gleich %[1]v sein",
		},
		fmtMsgLtOther: {
			"en": fmtMsgLtOther,
			"fr": "La valeur doit être inférieure à la valeur de la propriété '%[1]s'",
			"es": "El valor debe ser menor que el valor de la propiedad '%[1]s'",
			"it": "Il valore deve essere inferiore al valore della proprietà '%[1]s'",
			"de": "Wert muss kleiner sein als Wert der Eigenschaft '%[1]s'",
		},
		fmtMsgLteOther: {
			"en": fmtMsgLteOther,
			"fr": "La valeur doit être inférieure ou égale à la valeur de la propriété '%[1]s'",
			"es": "El valor debe ser menor o igual que el valor de la propiedad '%[1]s'",
			"it": "Il valore deve essere inferiore o uguale al valore della proprietà '%[1]s'",
			"de": "Wert muss kleiner oder gleich dem Wert der Eigenschaft '%[1]s' sein",
		},
		fmtMsgMinLen: {
			"en": fmtMsgMinLen,
			"fr": "La longueur de la valeur doit être d'au moins %[1]d",
			"es": "La longitud del valor debe ser al menos %[1]d",
			"it": "La lunghezza del valore deve essere almeno %[1]d",
			"de": "Wertlänge muss mindestens %[1]d betragen",
		},
		fmtMsgMinLenExc: {
			"en": fmtMsgMinLenExc,
			"fr": "La longueur de la valeur doit être supérieure à %[1]d",
			"es": "La longitud del valor debe ser mayor que %[1]d",
			"it": "La lunghezza del valore deve essere maggiore di %[1]d",
			"de": "Wertlänge muss größer sein als %[1]d",
		},
		fmtMsgMinMax: {
			"en": fmtMsgMinMax,
			"fr": "La longueur de la valeur doit être comprise entre %[1]d (%[2]s) et %[3]d (%[4]s)",
			"es": "La longitud del valor debe estar entre %[1]d (%[2]s) y %[3]d (%[4]s)",
			"it": "La lunghezza del valore deve essere compresa tra %[1]d (%[2]s) e %[3]d (%[4]s)",
			"de": "Wertlänge muss zwischen %[1]d (%[2]s) und %[3]d (%[4]s) liegen",
		},
		fmtMsgMultipleOf: {
			"en": fmtMsgMultipleOf,
			"fr": "La valeur doit être un multiple de %[1]d",
			"es": "El valor debe ser un múltiplo de %[1]d",
			"it": "Il valore deve essere un multiplo di %[1]d",
			"de": "Wert muss ein Vielfaches von %[1]d sein",
		},
		fmtMsgNotEqualsOther: {
			"en": fmtMsgNotEqualsOther,
			"fr": "La valeur ne doit pas être égale à la valeur de la propriété '%[1]s'",
			"es": "El valor no debe ser igual al valor de la propiedad '%[1]s'",
			"it": "Il valore non deve essere uguale al valore della proprietà '%[1]s'",
			"de": "Wert darf nicht gleich dem Wert der Eigenschaft '%[1]s' sein",
		},
		fmtMsgRange: {
			"en": fmtMsgRange,
			"fr": "La valeur doit être comprise entre %[1]v (%[2]s) et %[3]v (%[4]s)",
			"es": "El valor debe estar entre %[1]v (%[2]s) y %[3]v (%[4]s)",
			"it": "Il valore deve essere compreso tra %[1]v (%[2]s) e %[3]v (%[4]s)",
			"de": "Wert muss zwischen %[1]v (%[2]s) und %[3]v (%[4]s) liegen",
		},
		fmtMsgStringExactLen: {
			"en": fmtMsgStringExactLen,
			"fr": "La longueur de la valeur de la chaîne doit être de %[1]d caractères",
			"es": "La longitud del valor de la cadena debe ser %[1]d caracteres",
			"it": "La lunghezza del valore della stringa deve essere %[1]d caratteri",
			"de": "String-Wert muss %[1]d Zeichen lang sein",
		},
		fmtMsgStringMaxLen: {
			"en": fmtMsgStringMaxLen,
			"fr": "La longueur de la valeur de la chaîne ne doit pas dépasser %[1]d caractères",
			"es": "La longitud del valor de la cadena no debe exceder %[1]d caracteres",
			"it": "La lunghezza del valore della stringa non deve superare %[1]d caratteri",
			"de": "Stringwertlänge darf %[1]d Zeichen nicht überschreiten",
		},
		fmtMsgStringMaxLenExc: {
			"en": fmtMsgStringMaxLenExc,
			"fr": "La longueur de la valeur de la chaîne doit être inférieure à %[1]d caractères",
			"es": "La longitud del valor de la cadena debe ser inferior a %[1]d caracteres",
			"it": "La lunghezza del valore della stringa deve essere inferiore a %[1]d caratteri",
			"de": "String-Wert muss weniger als %[1]d Zeichen lang sein",
		},
		fmtMsgStringMinLen: {
			"en": fmtMsgStringMinLen,
			"fr": "La longueur de la valeur de la chaîne doit être d'au moins %[1]d caractères",
			"es": "La longitud del valor de la cadena debe ser de al menos %[1]d caracteres",
			"it": "La lunghezza del valore della stringa deve essere di almeno %[1]d caratteri",
			"de": "String-Wert muss mindestens %[1]d Zeichen lang sein",
		},
		fmtMsgStringMinLenExc: {
			"en": fmtMsgStringMinLenExc,
			"fr": "La longueur de la valeur de la chaîne doit être supérieure à %[1]d caractères",
			"es": "La longitud del valor de la cadena debe ser mayor que %[1]d caracteres",
			"it": "La lunghezza del valore della stringa deve essere maggiore di %[1]d caratteri",
			"de": "Stringwertlänge muss größer als %[1]d Zeichen sein",
		},
		fmtMsgStringMinMaxLen: {
			"en": fmtMsgStringMinMaxLen,
			"fr": "La longueur de la valeur de chaîne doit être comprise entre %[1]d (%[2]s) et %[3]d (%[4]s)",
			"es": "La longitud del valor de la cadena debe estar entre %[1]d (%[2]s) y %[3]d (%[4]s)",
			"it": "La lunghezza del valore della stringa deve essere compresa tra %[1]d (%[2]s) e %[3]d (%[4]s)",
			"de": "Stringwertlänge muss zwischen %[1]d (%[2]s) und %[3]d (%[4]s) liegen",
		},
		fmtMsgUnknownPresetPattern: {
			"en": fmtMsgUnknownPresetPattern,
			"fr": "Modèle prédéfini inconnu '%[1]s'",
			"es": "Patrón predeterminado desconocido '%[1]s'",
			"it": "Modello predefinito sconosciuto '%[1]s'",
			"de": "Unbekanntes voreingestelltes Muster '%[1]s'",
		},
		fmtMsgUuidCorrectVer: {
			"en": fmtMsgUuidCorrectVer,
			"fr": "La valeur doit être un UUID valide (version %[1]d)",
			"es": "El valor debe ser un UUID válido (versión %[1]d)",
			"it": "Il valore deve essere un UUID valido (versione %[1]d)",
			"de": "Wert muss eine gültige UUID sein (Version %[1]d)",
		},
		fmtMsgUuidMinVersion: {
			"en": fmtMsgUuidMinVersion,
			"fr": "La valeur doit être un UUID valide (version minimale %[1]d)",
			"es": "El valor debe ser un UUID válido (versión mínima %[1]d)",
			"it": "Il valore deve essere un UUID valido (versione minima %[1]d)",
			"de": "Wert muss eine gültige UUID sein (Mindestversion %[1]d)",
		},
		fmtMsgValidToken: {
			"en": fmtMsgValidToken,
			"fr": "La valeur de la chaîne doit être un jeton valide - %[1]s",
			"es": "El valor de la cadena debe ser un token válido - %[1]s",
			"it": "Il valore della stringa deve essere un token valido - %[1]s",
			"de": "String-Wert muss gültiges Token sein - %[1]s",
		},
		fmtMsgValueExpectedType: {
			"en": fmtMsgValueExpectedType,
			"fr": "Valeur supposée être de type %[1]s",
			"es": "Se espera que el valor sea del tipo %[1]s",
			"it": "Valore previsto di tipo %[1]s",
			"de": "Wert sollte vom Typ %[1]s sein",
		},
	},
}
