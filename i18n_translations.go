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

const (
	langEn = "en"
	langFr = "fr"
	langEs = "es"
	langIt = "it"
	langDe = "de"
)

var defaultInternalTranslator = &internalTranslator{
	Tokens: map[string]map[string]string{
		jsonTypeTokenAny: {
			langEn: jsonTypeTokenAny,
			langFr: "tout",
			langEs: "alguno",
			langIt: "qualsiasi",
			langDe: "beliebig",
		},
		jsonTypeTokenArray: {
			langEn: jsonTypeTokenArray,
			langFr: "tableau",
			langEs: "matriz",
			langIt: "array",
			langDe: "Array",
		},
		jsonTypeTokenBoolean: {
			langEn: jsonTypeTokenBoolean,
			langFr: "booléen",
			langEs: "booleano",
			langIt: "booleano",
			langDe: "boolesch",
		},
		jsonTypeTokenInteger: {
			langEn: jsonTypeTokenInteger,
			langFr: "entier",
			langEs: "entero",
			langIt: "intero",
			langDe: "Ganzzahl",
		},
		jsonTypeTokenNumber: {
			langEn: jsonTypeTokenNumber,
			langFr: "nombre",
			langEs: "número",
			langIt: "numero",
			langDe: "Nummer",
		},
		jsonTypeTokenObject: {
			langEn: jsonTypeTokenObject,
			langFr: "objet",
			langEs: "objeto",
			langIt: "oggetto",
			langDe: "Objekt",
		},
		jsonTypeTokenString: {
			langEn: jsonTypeTokenString,
			langFr: "chaîne",
			langEs: "cadena",
			langIt: "stringa",
			langDe: "Zeichenfolge",
		},
		tokenExclusive: {
			langEn: tokenExclusive,
			langFr: "exclusif",
			langEs: "exclusivo",
			langIt: "esclusivo",
			langDe: "exklusiv",
		},
		tokenInclusive: {
			langEn: tokenInclusive,
			langFr: "inclusif",
			langEs: "inclusivo",
			langIt: "comprensivo",
			langDe: "inklusive",
		},
		"century": {
			langEn: "century",
			langFr: "siècle",
			langEs: "siglo",
			langIt: "secolo",
			langDe: "Jahrhundert",
		},
		"century...": {
			langEn: "centuries",
			langFr: "siècles",
			langEs: "siglos",
			langIt: "secoli",
			langDe: "Jahrhunderte",
		},
		"day": {
			langEn: "day",
			langFr: "jour",
			langEs: "día",
			langIt: "giorno",
			langDe: "Tag",
		},
		"day...": {
			langEn: "days",
			langFr: "jours",
			langEs: "días",
			langIt: "giorni",
			langDe: "Tage",
		},
		"decade": {
			langEn: "decade",
			langFr: "décennie",
			langEs: "década",
			langIt: "decennio",
			langDe: "Jahrzehnt",
		},
		"decade...": {
			langEn: "decades",
			langFr: "décennies",
			langEs: "décadas",
			langIt: "decenni",
			langDe: "Jahrzehnte",
		},
		"hour": {
			langEn: "hour",
			langFr: "heure",
			langEs: "hora",
			langIt: "ora",
			langDe: "Stunde",
		},
		"hour...": {
			langEn: "hours",
			langFr: "heures",
			langEs: "horas",
			langIt: "ore",
			langDe: "Stunden",
		},
		"microsecond": {
			langEn: "microsecond",
			langFr: "microsecondes",
			langEs: "microsegundos",
			langIt: "microsecondo",
			langDe: "Mikrosekunde",
		},
		"microsecond...": {
			langEn: "microseconds",
			langFr: "microsecondes",
			langEs: "microsegundos",
			langIt: "microsecondi",
			langDe: "Mikrosekunden",
		},
		"millennium": {
			langEn: "millennium",
			langFr: "millénaire",
			langEs: "milenio",
			langIt: "millennio",
			langDe: "Jahrtausend",
		},
		"millennium...": {
			langEn: "millennia",
			langFr: "millénaires",
			langEs: "milenios",
			langIt: "millenni",
			langDe: "Jahrtausende",
		},
		"millisecond": {
			langEn: "millisecond",
			langFr: "milliseconde",
			langEs: "milisegundos",
			langIt: "millisecondo",
			langDe: "Millisekunde",
		},
		"millisecond...": {
			langEn: "milliseconds",
			langFr: "millisecondes",
			langEs: "milisegundos",
			langIt: "millisecondi",
			langDe: "Millisekunden",
		},
		"minute": {
			langEn: "minute",
			langFr: "minute",
			langEs: "minuto",
			langIt: "minuto",
			langDe: "Minute",
		},
		"minute...": {
			langEn: "minutes",
			langFr: "minutes",
			langEs: "minutos",
			langIt: "minuti",
			langDe: "Minuten",
		},
		"month": {
			langEn: "month",
			langFr: "mois",
			langEs: "mes",
			langIt: "mese",
			langDe: "Monat",
		},
		"month...": {
			langEn: "months",
			langFr: "mois",
			langEs: "meses",
			langIt: "mesi",
			langDe: "Monate",
		},
		"nanosecond": {
			langEn: "nanosecond",
			langFr: "nanoseconde",
			langEs: "nanosegundo",
			langIt: "nanosecondo",
			langDe: "Nanosekunde",
		},
		"nanosecond...": {
			langEn: "nanoseconds",
			langFr: "nanosecondes",
			langEs: "nanosegundos",
			langIt: "nanosecondi",
			langDe: "Nanosekunden",
		},
		"second": {
			langEn: "second",
			langFr: "seconde",
			langEs: "segundo",
			langIt: "secondo",
			langDe: "Sekunde",
		},
		"second...": {
			langEn: "seconds",
			langFr: "secondes",
			langEs: "segundos",
			langIt: "secondi",
			langDe: "Sekunden",
		},
		"week": {
			langEn: "week",
			langFr: "semaine",
			langEs: "semana",
			langIt: "settimana",
			langDe: "Woche",
		},
		"week...": {
			langEn: "weeks",
			langFr: "semaines",
			langEs: "semanas",
			langIt: "settimane",
			langDe: "Wochen",
		},
		"year": {
			langEn: "year",
			langFr: "année",
			langEs: "año",
			langIt: "anno",
			langDe: "Jahr",
		},
		"year...": {
			langEn: "years",
			langFr: "ans",
			langEs: "años",
			langIt: "anni",
			langDe: "Jahre",
		},
	},
	Messages: map[string]map[string]string{
		msgArrayElementMustBeObject: {
			langEn: msgArrayElementMustBeObject,
			langFr: "L'élément du tableau JSON doit être un objet",
			langEs: "El elemento de matriz JSON debe ser un objeto",
			langIt: "L'elemento dell'array JSON deve essere un oggetto",
			langDe: "JSON-Array-Element muss ein Objekt sein",
		},
		msgArrayElementMustNotBeNull: {
			langEn: msgArrayElementMustNotBeNull,
			langFr: "L'élément de tableau JSON ne doit pas être nul",
			langEs: "El elemento de matriz JSON no debe ser nulo",
			langIt: "L'elemento dell'array JSON non deve essere nullo",
			langDe: "JSON-Array-Element darf nicht null sein",
		},
		msgArrayUnique: {
			langEn: msgArrayUnique,
			langFr: "Les éléments du tableau doivent être uniques",
			langEs: "Los elementos del arreglo deben ser únicos",
			langIt: "Gli elementi dell'array devono essere univoci",
			langDe: "Array-Elemente müssen eindeutig sein",
		},
		msgDatetimeFuture: {
			langEn: msgDatetimeFuture,
			langFr: "La valeur doit être une date/heure valide dans le futur",
			langEs: "El valor debe ser una fecha/hora válida en el futuro",
			langIt: "Il valore deve essere una data/ora valida nel futuro",
			langDe: "Wert muss ein gültiges Datum/Zeit in der Zukunft sein",
		},
		msgDatetimeFutureOrPresent: {
			langEn: msgDatetimeFutureOrPresent,
			langFr: "La valeur doit être une date/heure valide dans le futur ou le présent",
			langEs: "El valor debe ser una fecha/hora válida en el futuro o presente",
			langIt: "Il valore deve essere una data/ora valida futura o presente",
			langDe: "Wert muss ein gültiges Datum/Zeit in der Zukunft oder Gegenwart sein",
		},
		msgDatetimePast: {
			langEn: msgDatetimePast,
			langFr: "La valeur doit être une date/heure valide dans le passé",
			langEs: "El valor debe ser una fecha/hora válida en el pasado",
			langIt: "Il valore deve essere una data/ora valida nel passato",
			langDe: "Wert muss ein gültiges Datum/Zeit in der Vergangenheit sein",
		},
		msgDatetimePastOrPresent: {
			langEn: msgDatetimePastOrPresent,
			langFr: "La valeur doit être une date/heure valide dans le passé ou le présent",
			langEs: "El valor debe ser una fecha/hora válida en el pasado o presente",
			langIt: "Il valore deve essere una data/ora valida nel passato o nel presente",
			langDe: "Wert muss ein gültiges Datum/Zeit in der Vergangenheit oder Gegenwart sein",
		},
		msgErrorReading: {
			langEn: msgErrorReading,
			langFr: "Erreur inattendue lors de la lecture du lecteur",
			langEs: "Error inesperado al leer el lector",
			langIt: "Errore imprevisto durante la lettura del lettore",
			langDe: "Unerwarteter Fehler beim Lesen des Lesegeräts",
		},
		msgErrorUnmarshall: {
			langEn: msgErrorUnmarshall,
			langFr: "Erreur inattendue lors du démarshalling",
			langEs: "Error inesperado durante la desorganización",
			langIt: "Errore imprevisto durante l'annullamento del marshalling",
			langDe: "Unerwarteter Fehler beim Unmarshalling",
		},
		msgExpectedJsonArray: {
			langEn: msgExpectedJsonArray,
			langFr: "JSON devrait être un tableau JSON",
			langEs: "Se esperaba que JSON fuera una matriz JSON",
			langIt: "JSON dovrebbe essere un array JSON",
			langDe: "JSON soll JSON-Array sein",
		},
		msgExpectedJsonObject: {
			langEn: msgExpectedJsonObject,
			langFr: "JSON devrait être un objet JSON",
			langEs: "Se esperaba que JSON fuera un objeto JSON",
			langIt: "JSON dovrebbe essere un oggetto JSON",
			langDe: "JSON soll JSON-Objekt sein",
		},
		msgFailure: {
			langEn: msgFailure,
			langFr: "Échec de la validation",
			langEs: "Validación fallida",
			langIt: "Convalida non riuscita",
			langDe: "Validierung fehlgeschlagen",
		},
		msgInvalidCharacters: {
			langEn: msgInvalidCharacters,
			langFr: "La valeur de la chaîne ne doit pas contenir de caractères invalides",
			langEs: "El valor de la cadena no debe tener caracteres inválidos",
			langIt: "Il valore della stringa non deve contenere caratteri non validi",
			langDe: "String-Wert darf keine ungültigen Zeichen enthalten",
		},
		msgInvalidProperty: {
			langEn: msgInvalidProperty,
			langFr: "Propriété invalide",
			langEs: "Propiedad no válida",
			langIt: "Proprietà non valida",
			langDe: "Ungültige Eigenschaft",
		},
		msgInvalidPropertyName: {
			langEn: msgInvalidPropertyName,
			langFr: "Nom de propriété invalide",
			langEs: "Nombre de propiedad inválido",
			langIt: "Nome proprietà non valido",
			langDe: "Ungültiger Eigenschaftsname",
		},
		msgMissingProperty: {
			langEn: msgMissingProperty,
			langFr: "Propriété manquante",
			langEs: "Propiedad faltante",
			langIt: "Proprietà mancante",
			langDe: "Fehlende Eigenschaft",
		},
		msgNegative: {
			langEn: msgNegative,
			langFr: "La valeur doit être négative",
			langEs: "El valor debe ser negativo",
			langIt: "Il valore deve essere negativo",
			langDe: "Wert muss negativ sein",
		},
		msgNegativeOrZero: {
			langEn: msgNegativeOrZero,
			langFr: "La valeur doit être négative ou nulle",
			langEs: "El valor debe ser negativo o cero",
			langIt: "Il valore deve essere negativo o zero",
			langDe: "Wert muss negativ oder Null sein",
		},
		msgNoControlChars: {
			langEn: msgNoControlChars,
			langFr: "La valeur de la chaîne ne doit pas contenir de caractères de contrôle",
			langEs: "El valor de la cadena no debe contener caracteres de control",
			langIt: "Il valore della stringa non deve contenere caratteri di controllo",
			langDe: "Stringwert darf keine Steuerzeichen enthalten",
		},
		msgNotBlankString: {
			langEn: msgNotBlankString,
			langFr: "La valeur de la chaîne ne doit pas être une chaîne vide",
			langEs: "El valor de la cadena no debe ser una cadena en blanco",
			langIt: "Il valore della stringa non deve essere una stringa vuota",
			langDe: "String-Wert darf kein leerer String sein",
		},
		msgNotEmpty: {
			langEn: msgNotEmpty,
			langFr: "La valeur ne doit pas être vide",
			langEs: "El valor no debe estar vacío",
			langIt: "Il valore non deve essere vuoto",
			langDe: "Wert darf nicht leer sein",
		},
		msgNotEmptyString: {
			langEn: msgNotEmptyString,
			langFr: "La valeur de la chaîne ne doit pas être une chaîne vide",
			langEs: "El valor de la cadena no debe ser una cadena vacía",
			langIt: "Il valore della stringa non deve essere una stringa vuota",
			langDe: "Stringwert darf kein leerer String sein",
		},
		msgNotJsonArray: {
			langEn: msgNotJsonArray,
			langFr: "JSON ne doit pas être un tableau JSON",
			langEs: "JSON no debe ser una matriz JSON",
			langIt: "JSON non deve essere un array JSON",
			langDe: "JSON darf kein JSON-Array sein",
		},
		msgNotJsonNull: {
			langEn: msgNotJsonNull,
			langFr: "JSON ne doit pas être JSON null",
			langEs: "JSON no debe ser JSON nulo",
			langIt: "JSON non deve essere JSON null",
			langDe: "JSON darf nicht JSON null sein",
		},
		msgNotJsonObject: {
			langEn: msgNotJsonObject,
			langFr: "JSON ne doit pas être un objet JSON",
			langEs: "JSON no debe ser un objeto JSON",
			langIt: "JSON non deve essere un oggetto JSON",
			langDe: "JSON darf kein JSON-Objekt sein",
		},
		msgPositive: {
			langEn: msgPositive,
			langFr: "La valeur doit être positive",
			langEs: "El valor debe ser positivo",
			langIt: "Il valore deve essere positivo",
			langDe: "Wert muss positiv sein",
		},
		msgPositiveOrZero: {
			langEn: msgPositiveOrZero,
			langFr: "La valeur doit être positive ou nulle",
			langEs: "El valor debe ser positivo o cero",
			langIt: "Il valore deve essere positivo o zero",
			langDe: "Wert muss positiv oder Null sein",
		},
		msgPropertyObjectValidatorError: {
			langEn: msgPropertyObjectValidatorError,
			langFr: "Erreur du validateur d'objet - n'autorise pas l'objet ou le tableau!",
			langEs: "Error del validador de objetos: ¡no permite el objeto o la matriz!",
			langIt: "Errore del validatore di oggetti - non consente l'oggetto o l'array!",
			langDe: "Objekt-Validator-Fehler - erlaubt kein Objekt oder Array!",
		},
		msgPropertyValueMustBeObject: {
			langEn: msgPropertyValueMustBeObject,
			langFr: "La valeur de la propriété doit être un objet",
			langEs: "El valor de la propiedad debe ser un objeto",
			langIt: "Il valore della proprietà deve essere un oggetto",
			langDe: "Eigenschaftswert muss ein Objekt sein",
		},
		msgPropertyRequiredWhen: {
			langEn: msgPropertyRequiredWhen,
			langFr: "La propriété est requise selon certains critères",
			langEs: "Se requiere propiedad bajo ciertos criterios",
			langIt: "L'immobile è richiesto secondo determinati criteri",
			langDe: "Eigentum wird unter bestimmten Kriterien benötigt",
		},
		msgPropertyUnwantedWhen: {
			langEn: msgPropertyUnwantedWhen,
			langFr: "La propriété ne doit pas être présente dans certaines conditions",
			langEs: "La propiedad no debe estar presente bajo ciertas condiciones",
			langIt: "L'immobile non deve essere presente in determinate condizioni",
			langDe: "Die Immobilie darf unter bestimmten Voraussetzungen nicht vorhanden sein",
		},
		msgRequestBodyEmpty: {
			langEn: msgRequestBodyEmpty,
			langFr: "Le corps de la requête est vide",
			langEs: "El cuerpo de la solicitud está vacío",
			langIt: "Il corpo della richiesta è vuoto",
			langDe: "Anfragetext ist leer",
		},
		msgRequestBodyExpectedJsonArray: {
			langEn: msgRequestBodyExpectedJsonArray,
			langFr: "Le corps de la requête devrait être un tableau JSON",
			langEs: "Se espera que el cuerpo de la solicitud sea una matriz JSON",
			langIt: "Il corpo della richiesta dovrebbe essere un array JSON",
			langDe: "Anforderungstext soll JSON-Array sein",
		},
		msgRequestBodyExpectedJsonObject: {
			langEn: msgRequestBodyExpectedJsonObject,
			langFr: "Le corps de la requête devrait être un objet JSON",
			langEs: "Se espera que el cuerpo de la solicitud sea un objeto JSON",
			langIt: "Il corpo della richiesta dovrebbe essere un oggetto JSON",
			langDe: "Anfragetext soll JSON-Objekt sein",
		},
		msgRequestBodyNotJsonArray: {
			langEn: msgRequestBodyNotJsonArray,
			langFr: "Le corps de la requête ne doit pas être un tableau JSON",
			langEs: "El cuerpo de la solicitud no debe ser una matriz JSON",
			langIt: "Il corpo della richiesta non deve essere un array JSON",
			langDe: "Anfragetext darf kein JSON-Array sein",
		},
		msgRequestBodyNotJsonNull: {
			langEn: msgRequestBodyNotJsonNull,
			langFr: "Le corps de la requête ne doit pas être nul en JSON",
			langEs: "El cuerpo de la solicitud no debe ser JSON nulo",
			langIt: "Il corpo della richiesta non deve essere JSON null",
			langDe: "Anfragetext darf nicht JSON null sein",
		},
		msgRequestBodyNotJsonObject: {
			langEn: msgRequestBodyNotJsonObject,
			langFr: "Le corps de la requête ne doit pas être un objet JSON",
			langEs: "El cuerpo de la solicitud no debe ser un objeto JSON",
			langIt: "Il corpo della richiesta non deve essere un oggetto JSON",
			langDe: "Anfragetext darf kein JSON-Objekt sein",
		},
		msgStringValidJson: {
			langEn: msgStringValidJson,
			langFr: "La valeur de la chaîne doit être un JSON valide",
			langEs: "El valor de cadena debe ser JSON válido",
			langIt: "Il valore della stringa deve essere JSON valido",
			langDe: "Zeichenfolgenwert muss gültiges JSON sein",
		},
		msgStringLowercase: {
			langEn: msgStringLowercase,
			langFr: "La valeur de la chaîne ne doit contenir que des lettres minuscules",
			langEs: "El valor de la cadena debe contener solo letras minúsculas",
			langIt: "Il valore della stringa deve contenere solo lettere minuscole",
			langDe: "Stringwert darf nur Kleinbuchstaben enthalten",
		},
		msgStringUppercase: {
			langEn: msgStringUppercase,
			langFr: "La valeur de la chaîne ne doit contenir que des lettres majuscules",
			langEs: "El valor de la cadena debe contener solo letras mayúsculas",
			langIt: "Il valore della stringa deve contenere solo lettere maiuscole",
			langDe: "String-Wert darf nur Großbuchstaben enthalten",
		},
		msgUnableToDecode: {
			langEn: msgUnableToDecode,
			langFr: "Impossible de décoder en JSON",
			langEs: "No se puede decodificar como JSON",
			langIt: "Impossibile decodificare come JSON",
			langDe: "Als JSON kann nicht dekodiert werden",
		},
		msgUnableToDecodeRequest: {
			langEn: msgUnableToDecodeRequest,
			langFr: "Impossible de décoder le corps de la requête en JSON",
			langEs: "No se puede decodificar el cuerpo de la solicitud como JSON",
			langIt: "Impossibile decodificare il corpo della richiesta come JSON",
			langDe: "Anforderungstext konnte nicht als JSON entschlüsselt werden",
		},
		msgUnicodeNormalization: {
			langEn: msgUnicodeNormalization,
			langFr: "La valeur de la chaîne doit être une forme de normalisation correcte",
			langEs: "El valor de la cadena debe ser la forma de normalización correcta",
			langIt: "Il valore della stringa deve essere un modulo di normalizzazione corretto",
			langDe: "String-Wert muss korrekte Normalisierungsform sein",
		},
		msgUnicodeNormalizationNFC: {
			langEn: msgUnicodeNormalizationNFC,
			langFr: "La valeur de la chaîne doit être la forme de normalisation correcte NFC",
			langEs: "El valor de la cadena debe ser la normalización correcta de NFC",
			langIt: "Il valore della stringa deve essere la normalizzazione corretta da NFC",
			langDe: "Stringwert muss korrekte Normalisierung von NFC sein",
		},
		msgUnicodeNormalizationNFD: {
			langEn: msgUnicodeNormalizationNFD,
			langFr: "La valeur de la chaîne doit être la forme de normalisation correcte NFD",
			langEs: "El valor de la cadena debe ser el formulario de normalización correcto NFD",
			langIt: "Il valore della stringa deve essere la normalizzazione corretta dal modulo NFD",
			langDe: "String-Wert muss korrekte Normalisierung von NFD sein",
		},
		msgUnicodeNormalizationNFKC: {
			langEn: msgUnicodeNormalizationNFKC,
			langFr: "La valeur de la chaîne doit être la forme de normalisation correcte NFKC",
			langEs: "El valor de la cadena debe ser el formulario de normalización correcto NFKC",
			langIt: "Il valore della stringa deve essere la normalizzazione corretta da NFKC",
			langDe: "String-Wert muss korrekte Normalisierung von NFKC sein",
		},
		msgUnicodeNormalizationNFKD: {
			langEn: msgUnicodeNormalizationNFKD,
			langFr: "La valeur de la chaîne doit être la forme de normalisation correcte NFKD",
			langEs: "El valor de la cadena debe ser el formulario de normalización correcto NFKD",
			langIt: "Il valore della stringa deve essere la normalizzazione corretta da NFKD",
			langDe: "String-Wert muss korrekte Normalisierung von NFKD sein",
		},
		msgUnknownProperty: {
			langEn: msgUnknownProperty,
			langFr: "Propriété inconnue",
			langEs: "Propiedad desconocida",
			langIt: "Proprietà sconosciuta",
			langDe: "Unbekanntes Eigentum",
		},
		msgOnlyProperty: {
			langEn: msgOnlyProperty,
			langFr: "La propriété ne peut pas être présente avec d'autres propriétés",
			langEs: "La propiedad no puede estar presente con otras propiedades",
			langIt: "L'immobile non può essere presente con altri immobili",
			langDe: "Eigenschaft kann nicht mit anderen Eigenschaften vorhanden sein",
		},
		msgUnwantedProperty: {
			langEn: msgUnwantedProperty,
			langFr: "La propriété ne doit pas être présente",
			langEs: "La propiedad no debe estar presente",
			langIt: "L'immobile non deve essere presente",
			langDe: "Eigenschaft darf nicht vorhanden sein",
		},
		msgValueCannotBeNull: {
			langEn: msgValueCannotBeNull,
			langFr: "La valeur ne peut pas être nulle",
			langEs: "El valor no puede ser nulo",
			langIt: "Il valore non può essere nullo",
			langDe: "Wert darf nicht null sein",
		},
		msgValueMustBeArray: {
			langEn: msgValueMustBeArray,
			langFr: "La valeur doit être un tableau",
			langEs: "El valor debe ser una matriz",
			langIt: "Il valore deve essere un array",
			langDe: "Wert muss ein Array sein",
		},
		msgValueMustBeObject: {
			langEn: msgValueMustBeObject,
			langFr: "La valeur doit être un objet",
			langEs: "El valor debe ser un objeto",
			langIt: "Il valore deve essere un oggetto",
			langDe: "Wert muss ein Objekt sein",
		},
		msgValueMustBeObjectOrArray: {
			langEn: msgValueMustBeObjectOrArray,
			langFr: "La valeur doit être un objet ou un tableau",
			langEs: "El valor debe ser un objeto o matriz",
			langIt: "Il valore deve essere un oggetto o un array",
			langDe: "Wert muss ein Objekt oder Array sein",
		},
		msgValidCardNumber: {
			langEn: msgValidCardNumber,
			langFr: "La valeur doit être un numéro de carte valide",
			langEs: "El valor debe ser un número de tarjeta válido",
			langIt: "Il valore deve essere un numero di carta valido",
			langDe: "Wert muss eine gültige Kartennummer sein",
		},
		msgValidCountryCode: {
			langEn: msgValidCountryCode,
			langFr: "La valeur doit être un code de pays ISO-3166 valide",
			langEs: "El valor debe ser un código de país ISO-3166 válido",
			langIt: "Il valore deve essere un codice paese ISO-3166 valido",
			langDe: "Wert muss ein gültiger ISO-3166-Ländercode sein",
		},
		msgValidCurrencyCode: {
			langEn: msgValidCurrencyCode,
			langFr: "La valeur doit être un code de devise ISO-4217 valide",
			langEs: "El valor debe ser un código de moneda ISO-4217 válido",
			langIt: "Il valore deve essere un codice valuta ISO-4217 valido",
			langDe: "Wert muss ein gültiger ISO-4217-Währungscode sein",
		},
		msgValidEmail: {
			langEn: msgValidEmail,
			langFr: "La valeur doit être une adresse e-mail",
			langEs: "El valor debe ser una dirección de correo electrónico",
			langIt: "Il valore deve essere un indirizzo email",
			langDe: "Wert muss eine E-Mail-Adresse sein",
		},
		msgValidLanguageCode: {
			langEn: msgValidLanguageCode,
			langFr: "La valeur doit être un code de langue valide",
			langEs: "El valor debe ser un código de idioma válido",
			langIt: "Il valore deve essere un codice lingua valido",
			langDe: "Wert muss ein gültiger Sprachcode sein",
		},
		msgValidISODate: {
			langEn: msgValidISODate,
			langFr: "La valeur doit être une chaîne de date valide (format : AAAA-MM-JJ)",
			langEs: "El valor debe ser una cadena de fecha válida (formato: AAAA-MM-DD)",
			langIt: "Il valore deve essere una stringa di data valida (formato: AAAA-MM-GG)",
			langDe: "Wert muss eine gültige Datumszeichenfolge sein (Format: JJJJ-MM-TT)",
		},
		msgValidISODatetimeFormatFull: {
			langEn: msgValidISODatetimeFormatFull,
			langFr: "La valeur doit être une chaîne de date/heure valide (format : AAAA-MM-JJThh : mm:ss.sss [Z|+-hh:mm])",
			langEs: "El valor debe ser una cadena de fecha/hora válida (formato: AAAA-MM-DDThh: mm:ss.sss [Z|+- hh:mm ])",
			langIt: "Il valore deve essere una stringa di data/ora valida (formato: AAAA-MM-GGThh: mm:ss.sss [Z|+- hh:mm ])",
			langDe: "Wert muss ein gültiger Datums-/Uhrzeit-String sein (Format: YYYY-MM-DDThh: mm:ss.sss [Z|+- hh:mm ])",
		},
		msgValidISODatetimeFormatMin: {
			langEn: msgValidISODatetimeFormatMin,
			langFr: "La valeur doit être une chaîne date/heure valide (format : AAAA-MM-JJThh: mm:ss)",
			langEs: "El valor debe ser una cadena de fecha/hora válida (formato: AAAA-MM-DDThh: mm:ss)",
			langIt: "Il valore deve essere una stringa di data/ora valida (formato: AAAA-MM-GGThh: mm:ss)",
			langDe: "Wert muss ein gültiger Datums-/Uhrzeit-String sein (Format: YYYY-MM-DDThh: mm:ss )",
		},
		msgValidISODatetimeFormatNoOffs: {
			langEn: msgValidISODatetimeFormatNoOffs,
			langFr: "La valeur doit être une chaîne date/heure valide (format : AAAA-MM-JJThh: mm:ss.sss)",
			langEs: "El valor debe ser una cadena de fecha/hora válida (formato: AAAA-MM-DDThh: mm:ss.sss)",
			langIt: "Il valore deve essere una stringa di data/ora valida (formato: AAAA-MM-GGThh: mm:ss.sss)",
			langDe: "Wert muss ein gültiger Datums-/Uhrzeit-String sein (Format: YYYY-MM-DDThh: mm:ss.sss )",
		},
		msgValidISODatetimeFormatNoMillis: {
			langEn: msgValidISODatetimeFormatNoMillis,
			langFr: "La valeur doit être une chaîne de date/heure valide (format : AAAA-MM-JJThh : mm:ss [Z|+-hh:mm])",
			langEs: "El valor debe ser una cadena de fecha/hora válida (formato: AAAA-MM-DDThh: mm:ss [Z|+- hh:mm ])",
			langIt: "Il valore deve essere una stringa di data/ora valida (formato: AAAA-MM-GGThh: mm:ss [Z|+- hh:mm ])",
			langDe: "Wert muss ein gültiger Datums-/Uhrzeit-String sein (Format: YYYY-MM-DDThh: mm:ss [Z|+- hh:mm ])",
		},
		msgValidTimezone: {
			langEn: msgValidTimezone,
			langFr: "La valeur doit être un fuseau horaire valide",
			langEs: "El valor debe ser una zona horaria válida",
			langIt: "Il valore deve essere un fuso orario valido",
			langDe: "Wert muss eine gültige Zeitzone sein",
		},
		msgValidPattern: {
			langEn: msgValidPattern,
			langFr: "La valeur de la chaîne doit avoir un modèle valide",
			langEs: "El valor de la cadena debe tener un patrón válido",
			langIt: "Il valore della stringa deve avere un modello valido",
			langDe: "String-Wert muss gültiges Muster haben",
		},
		msgValidUuid: {
			langEn: msgValidUuid,
			langFr: "La valeur doit être un UUID valide",
			langEs: "El valor debe ser un UUID válido",
			langIt: "Il valore deve essere un UUID valido",
			langDe: "Wert muss eine gültige UUID sein",
		},
		msgPresetAlpha: {
			langEn: msgPresetAlpha,
			langFr: "La valeur ne doit être que des caractères alphabétiques (A-Z, a-z)",
			langEs: "El valor debe ser solo caracteres alfabéticos (A-Z, a-z)",
			langIt: "Il valore deve essere solo caratteri alfabetici (A-Z, a-z)",
			langDe: "Wert darf nur aus Buchstaben bestehen (A-Z, a-z)",
		},
		msgPresetAlphaNumeric: {
			langEn: msgPresetAlphaNumeric,
			langFr: "La valeur ne doit être que des caractères alphanumériques (A-Z, a-z, 0-9)",
			langEs: "El valor debe ser solo caracteres alfanuméricos (A-Z, a-z, 0-9)",
			langIt: "Il valore deve essere solo caratteri alfanumerici (A-Z, a-z, 0-9)",
			langDe: "Wert darf nur aus alphanumerischen Zeichen bestehen (A-Z, a-z, 0-9)",
		},
		msgPresetBarcode: {
			langEn: "Value must be a valid barcode",
			langFr: "La valeur doit être un code-barres valide",
			langEs: "El valor debe ser un código de barras válido",
			langIt: "Il valore deve essere un codice a barre valido",
			langDe: "Wert muss ein gültiger Strichcode sein",
		},
		msgPresetBase64: {
			langEn: msgPresetBase64,
			langFr: "La valeur doit être une chaîne valide encodée en base64",
			langEs: "El valor debe ser una cadena codificada en base64 válida",
			langIt: "Il valore deve essere una stringa codificata base64 valida",
			langDe: "Wert muss eine gültige base64-codierte Zeichenfolge sein",
		},
		msgPresetBase64URL: {
			langEn: msgPresetBase64URL,
			langFr: "La valeur doit être une chaîne encodée URL base64 valide",
			langEs: "El valor debe ser una cadena codificada en URL base64 válida",
			langIt: "Il valore deve essere una stringa codificata URL base64 valida",
			langDe: "Wert muss eine gültige Base64-URL-codierte Zeichenfolge sein",
		},
		msgPresetCMYK: {
			langEn:  msgPresetCMYK,
			"en-US": "Value must be a valid cmyk() color string",
			langFr:  "La valeur doit être une chaîne de couleur cmyk() valide",
			langEs:  "El valor debe ser una cadena de color cmyk() válida",
			langIt:  "Il valore deve essere una stringa di colore cmyk() valida",
			langDe:  "Wert muss eine gültige cmyk()-Farbzeichenfolge sein",
		},
		msgPresetCMYK300: {
			langEn:  msgPresetCMYK300,
			"en-US": "Value must be a valid cmyk() color string (maximum 300%)",
			langFr:  "La valeur doit être une chaîne de couleur cmyk() valide (maximum 300 %)",
			langEs:  "El valor debe ser una cadena de color cmyk() válida (máximo 300 %)",
			langIt:  "Il valore deve essere una stringa di colore cmyk() valida (massimo 300%)",
			langDe:  "Wert muss eine gültige cmyk()-Farbzeichenfolge sein (maximal 300 %)",
		},
		msgPresetE164: {
			langEn: msgPresetE164,
			langFr: "La valeur doit être un code E.164 valide",
			langEs: "El valor debe ser un código E.164 válido",
			langIt: "Il valore deve essere un codice E.164 valido",
			langDe: "Wert muss ein gültiger E.164-Code sein",
		},
		msgPresetEAN: {
			langEn: msgPresetEAN,
			langFr: "La valeur doit être un code EAN valide",
			langEs: "El valor debe ser un código EAN válido",
			langIt: "Il valore deve essere un codice EAN valido",
			langDe: "Wert muss ein gültiger EAN-Code sein",
		},
		msgPresetEAN8: {
			langEn: msgPresetEAN8,
			langFr: "La valeur doit être un code EAN-8 valide",
			langEs: "El valor debe ser un código EAN-8 válido",
			langIt: "Il valore deve essere un codice EAN-8 valido",
			langDe: "Wert muss ein gültiger EAN-8-Code sein",
		},
		msgPresetEAN13: {
			langEn: msgPresetEAN13,
			langFr: "La valeur doit être un code EAN-13 valide",
			langEs: "El valor debe ser un código EAN-13 válido",
			langIt: "Il valore deve essere un codice EAN-13 valido",
			langDe: "Wert muss ein gültiger EAN-13-Code sein",
		},
		msgPresetDUN14: {
			langEn: msgPresetDUN14,
			langFr: "La valeur doit être un code DUN-14 valide",
			langEs: "El valor debe ser un código DUN-14 válido",
			langIt: "Il valore deve essere un codice DUN-14 valido",
			langDe: "Wert muss ein gültiger DUN-14-Code sein",
		},
		msgPresetEAN14: {
			langEn: msgPresetEAN14,
			langFr: "La valeur doit être un code EAN-14 valide",
			langEs: "El valor debe ser un código EAN-14 válido",
			langIt: "Il valore deve essere un codice EAN-14 valido",
			langDe: "Wert muss ein gültiger EAN-14-Code sein",
		},
		msgPresetEAN18: {
			langEn: msgPresetEAN18,
			langFr: "La valeur doit être un code EAN-18 valide",
			langEs: "El valor debe ser un código EAN-18 válido",
			langIt: "Il valore deve essere un codice EAN-18 valido",
			langDe: "Wert muss ein gültiger EAN-18-Code sein",
		},
		msgPresetEAN99: {
			langEn: msgPresetEAN99,
			langFr: "La valeur doit être un code EAN-99 valide",
			langEs: "El valor debe ser un código EAN-99 válido",
			langIt: "Il valore deve essere un codice EAN-99 valido",
			langDe: "Wert muss ein gültiger EAN-99-Code sein",
		},
		msgPresetHexadecimal: {
			langEn: msgPresetHexadecimal,
			langFr: "La valeur doit être une chaîne hexadécimale valide",
			langEs: "El valor debe ser una cadena hexadecimal válida",
			langIt: "Il valore deve essere una stringa esadecimale valida",
			langDe: "Wert muss eine gültige hexadezimale Zeichenfolge sein",
		},
		msgPresetHsl: {
			langEn:  msgPresetHsl,
			"en-US": "Value must be a valid hsl() color string",
			langFr:  "La valeur doit être une chaîne de couleur hsl() valide",
			langEs:  "El valor debe ser una cadena de color hsl() válida",
			langIt:  "Il valore deve essere una stringa di colore hsl() valida",
			langDe:  "Wert muss eine gültige hsl() Farbzeichenfolge sein",
		},
		msgPresetHsla: {
			langEn:  msgPresetHsla,
			"en-US": "Value must be a valid hsla() color string",
			langFr:  "La valeur doit être une chaîne de couleur hsla() valide",
			langEs:  "El valor debe ser una cadena de color hsla() válida",
			langIt:  "Il valore deve essere una stringa di colore hsla() valida",
			langDe:  "Wert muss eine gültige hsla() Farbzeichenfolge sein",
		},
		msgPresetHtmlColor: {
			langEn:  msgPresetHtmlColor,
			"en-US": "Value must be a valid HTML color string",
			langFr:  "La valeur doit être une chaîne de couleur HTML valide",
			langEs:  "El valor debe ser una cadena de color HTML válida",
			langIt: "Il valore deve essere una 	stringa di colori HTML valida",
			langDe: "Wert muss ein gültiger HTML -Farbstring sein",
		},
		msgPresetInteger: {
			langEn: msgPresetInteger,
			langFr: "La valeur doit être une chaîne d'entiers valide (caractères 0-9)",
			langEs: "El valor debe ser una cadena entera válida (caracteres 0-9)",
			langIt: "Il valore deve essere una stringa intera valida (caratteri 0-9)",
			langDe: "Wert muss eine gültige Ganzzahl sein (Zeichen 0-9)",
		},
		msgPresetISBN: {
			langEn: msgPresetISBN,
			langFr: "La valeur doit être un ISBN valide",
			langEs: "El valor debe ser un ISBN válido",
			langIt: "Il valore deve essere un ISBN valido",
			langDe: "Wert muss eine gültige ISBN sein",
		},
		msgPresetISBN10: {
			langEn: msgPresetISBN10,
			langFr: "La valeur doit être un ISBN-10 valide",
			langEs: "El valor debe ser un ISBN-10 válido",
			langIt: "Il valore deve essere un ISBN-10 valido",
			langDe: "Wert muss eine gültige ISBN-10 sein",
		},
		msgPresetISBN13: {
			langEn: msgPresetISBN13,
			langFr: "La valeur doit être un ISBN-13 valide",
			langEs: "El valor debe ser un ISBN-13 válido",
			langIt: "Il valore deve essere un ISBN-13 valido",
			langDe: "Wert muss eine gültige ISBN-13 sein",
		},
		msgPresetISSN: {
			langEn: msgPresetISSN,
			langFr: "La valeur doit être un ISSN valide",
			langEs: "El valor debe ser un ISSN válido",
			langIt: "Il valore deve essere un ISSN valido",
			langDe: "Wert muss eine gültige ISSN sein",
		},
		msgPresetNumeric: {
			langEn: msgPresetNumeric,
			langFr: "La valeur doit être une chaîne numérique valide",
			langEs: "El valor debe ser una cadena de números válida",
			langIt: "Il valore deve essere una stringa numerica valida",
			langDe: "Wert muss eine gültige Zahlenfolge sein",
		},
		msgPresetPublication: {
			langEn: msgPresetPublication,
			langFr: "La valeur doit être un ISBN ou un ISSN valide",
			langEs: "El valor debe ser un ISBN o ISSN válido",
			langIt: "Il valore deve essere un ISBN o ISSN valido",
			langDe: "Wert muss eine gültige ISBN oder ISSN sein",
		},
		msgPresetRgb: {
			langEn:  msgPresetRgb,
			"en-US": "Value must be a valid rgb() color string",
			langFr:  "La valeur doit être une chaîne de couleur rgb() valide",
			langEs:  "El valor debe ser una cadena de color rgb() válida",
			langIt:  "Il valore deve essere una stringa di colore rgb() valida",
			langDe:  "Wert muss eine gültige rgb() Farbzeichenfolge sein",
		},
		msgPresetRgba: {
			langEn:  msgPresetRgba,
			"en-US": "Value must be a valid rgba() color string",
			langFr:  "La valeur doit être une chaîne de couleur rgba() valide",
			langEs:  "El valor debe ser una cadena de color rgba() válida",
			langIt:  "Il valore deve essere una stringa di colore rgba() valida",
			langDe:  "Wert muss eine gültige rgba() Farbzeichenfolge sein",
		},
		msgPresetRgbIcc: {
			langEn:  msgPresetRgbIcc,
			"en-US": "Value must be a valid rgb-icc() color string",
			langFr:  "La valeur doit être une chaîne de couleur rgb-icc() valide",
			langEs:  "El valor debe ser una cadena de color rgb-icc() válida",
			langIt:  "Il valore deve essere una stringa di colore rgb-icc() valida",
			langDe:  "Wert muss eine gültige rgb-icc() Farbzeichenfolge sein",
		},
		msgPresetULID: {
			langEn: msgPresetULID,
			langFr: "La valeur doit être un ULID valide",
			langEs: "El valor debe ser un ULID válido",
			langIt: "Il valore deve essere un ULID valido",
			langDe: "Wert muss eine gültige ULID sein",
		},
		msgPresetUPC: {
			langEn: msgPresetUPC,
			langFr: "La valeur doit être un code UPC valide (UPC-A ou UPC-E)",
			langEs: "El valor debe ser un código UPC válido (UPC-A o UPC-E)",
			langIt: "Il valore deve essere un codice UPC valido (UPC-A o UPC-E)",
			langDe: "Wert muss ein gültiger UPC-Code sein (UPC-A oder UPC-E)",
		},
		msgPresetUPCA: {
			langEn: msgPresetUPCA,
			langFr: "La valeur doit être un code UPC-A valide",
			langEs: "El valor debe ser un código UPC-A válido",
			langIt: "Il valore deve essere un codice UPC-A valido",
			langDe: "Wert muss ein gültiger UPC-A-Code sein",
		},
		msgPresetUPCE: {
			langEn: msgPresetUPCE,
			langFr: "La valeur doit être un code UPC-E valide",
			langEs: "El valor debe ser un código UPC-E válido",
			langIt: "Il valore deve essere un codice UPC-E valido",
			langDe: "Wert muss ein gültiger UPC-E-Code sein",
		},
		msgPresetUuid1: {
			langEn: msgPresetUuid1,
			langFr: "La valeur doit être un UUID valide (Version 1)",
			langEs: "El valor debe ser un UUID válido (Versión 1)",
			langIt: "Il valore deve essere un UUID valido (versione 1)",
			langDe: "Wert muss eine gültige UUID sein (Version 1)",
		},
		msgPresetUuid2: {
			langEn: msgPresetUuid2,
			langFr: "La valeur doit être un UUID valide (Version 2)",
			langEs: "El valor debe ser un UUID válido (Versión 2)",
			langIt: "Il valore deve essere un UUID valido (versione 2)",
			langDe: "Wert muss eine gültige UUID sein (Version 2)",
		},
		msgPresetUuid3: {
			langEn: msgPresetUuid3,
			langFr: "La valeur doit être un UUID valide (Version 3)",
			langEs: "El valor debe ser un UUID válido (Versión 3)",
			langIt: "Il valore deve essere un UUID valido (versione 3)",
			langDe: "Wert muss eine gültige UUID sein (Version 3)",
		},
		msgPresetUuid4: {
			langEn: msgPresetUuid4,
			langFr: "La valeur doit être un UUID valide (Version 4)",
			langEs: "El valor debe ser un UUID válido (Versión 4)",
			langIt: "Il valore deve essere un UUID valido (versione 4)",
			langDe: "Wert muss eine gültige UUID sein (Version 4)",
		},
		msgPresetUuid5: {
			langEn: msgPresetUuid5,
			langFr: "La valeur doit être un UUID valide (Version 5)",
			langEs: "El valor debe ser un UUID válido (Versión 5)",
			langIt: "Il valore deve essere un UUID valido (versione 5)",
			langDe: "Wert muss eine gültige UUID sein (Version 5)",
		},
		msgValidMAC: {
			langEn: msgValidMAC,
			langFr: "La valeur de la chaîne doit être une adresse MAC valide",
			langEs: "El valor de cadena debe ser una dirección MAC válida",
			langIt: "Il valore della stringa deve essere un indirizzo MAC valido",
			langDe: "String-Wert muss eine gültige MAC-Adresse sein",
		},
		msgValidCIDR: {
			langEn: msgValidCIDR,
			langFr: "La valeur de chaîne doit être une adresse CIDR valide",
			langEs: "El valor de cadena debe ser una dirección CIDR válida",
			langIt: "Il valore della stringa deve essere un indirizzo CIDR valido",
			langDe: "String-Wert muss eine gültige CIDR-Adresse sein",
		},
		msgValidCIDRv4: {
			langEn: msgValidCIDRv4,
			langFr: "La valeur de chaîne doit être une adresse CIDR (version 4) valide",
			langEs: "El valor de cadena debe ser una dirección CIDR (versión 4) válida",
			langIt: "Il valore della stringa deve essere un indirizzo CIDR (versione 4) valido",
			langDe: "String-Wert muss eine gültige CIDR-Adresse (Version 4) sein",
		},
		msgValidCIDRv6: {
			langEn: msgValidCIDRv6,
			langFr: "La valeur de chaîne doit être une adresse CIDR (version 6) valide",
			langEs: "El valor de cadena debe ser una dirección CIDR (versión 6) válida",
			langIt: "Il valore della stringa deve essere un indirizzo CIDR (versione 6) valido",
			langDe: "String-Wert muss eine gültige CIDR-Adresse (Version 6) sein",
		},
		msgValidIP: {
			langEn: msgValidIP,
			langFr: "La chaîne doit être une adresse IP valide",
			langEs: "La cadena debe ser una dirección IP válida",
			langIt: "La stringa deve essere un indirizzo IP valido",
			langDe: "String muss eine gültige IP-Adresse sein",
		},
		msgValidIPv4: {
			langEn: msgValidIPv4,
			langFr: "La chaîne doit être une adresse IP valide (version 4)",
			langEs: "La cadena debe ser una dirección IP válida (versión 4)",
			langIt: "La stringa deve essere un indirizzo IP (versione 4) valido",
			langDe: "String muss eine gültige IP-Adresse (Version 4) sein",
		},
		msgValidIPv6: {
			langEn: msgValidIPv6,
			langFr: "La chaîne doit être une adresse IP valide (version 6)",
			langEs: "La cadena debe ser una dirección IP válida (versión 6)",
			langIt: "La stringa deve essere un indirizzo IP (versione 6) valido",
			langDe: "String muss eine gültige IP-Adresse (Version 6) sein",
		},
		msgValidTCP: {
			langEn: msgValidTCP,
			langFr: "La chaîne doit être une adresse TCP valide",
			langEs: "La cadena debe ser una dirección TCP válida",
			langIt: "La stringa deve essere un indirizzo TCP valido",
			langDe: "String muss eine gültige TCP-Adresse sein",
		},
		msgValidTCPv4: {
			langEn: msgValidTCPv4,
			langFr: "La chaîne doit être une adresse TCP (version 4) valide",
			langEs: "La cadena debe ser una dirección TCP válida (versión 4)",
			langIt: "La stringa deve essere un indirizzo TCP (versione 4) valido",
			langDe: "String muss eine gültige TCP-Adresse (Version 4) sein",
		},
		msgValidTCPv6: {
			langEn: msgValidTCPv6,
			langFr: "La chaîne doit être une adresse TCP (version 6) valide",
			langEs: "La cadena debe ser una dirección TCP válida (versión 6)",
			langIt: "La stringa deve essere un indirizzo TCP (versione 6) valido",
			langDe: "String muss eine gültige TCP-Adresse (Version 6) sein",
		},
		msgValidUDP: {
			langEn: msgValidUDP,
			langFr: "La chaîne doit être une adresse UDP valide",
			langEs: "La cadena debe ser una dirección UDP válida",
			langIt: "La stringa deve essere un indirizzo UDP valido",
			langDe: "String muss eine gültige UDP-Adresse sein",
		},
		msgValidUDPv4: {
			langEn: msgValidUDPv4,
			langFr: "La chaîne doit être une adresse UDP (version 4) valide",
			langEs: "La cadena debe ser una dirección UDP (versión 4) válida",
			langIt: "La stringa deve essere un indirizzo UDP (versione 4) valido",
			langDe: "String muss eine gültige UDP-Adresse (Version 4) sein",
		},
		msgValidUDPv6: {
			langEn: msgValidUDPv6,
			langFr: "La chaîne doit être une adresse UDP (version 6) valide",
			langEs: "La cadena debe ser una dirección UDP (versión 6) válida",
			langIt: "La stringa deve essere un indirizzo UDP (versione 6) valido",
			langDe: "String muss eine gültige UDP-Adresse (Version 6) sein",
		},
		msgValidTld: {
			langEn: msgValidTld,
			langFr: "La chaîne doit être un TLD valide",
			langEs: "La cadena debe ser un TLD válido",
			langIt: "La stringa deve essere un TLD valido",
			langDe: "String muss eine gültige TLD sein",
		},
		msgValidHostname: {
			langEn: msgValidHostname,
			langFr: "La chaîne doit être un nom d'hôte valide",
			langEs: "La cadena debe ser un nombre de host válido",
			langIt: "La stringa deve essere un nome host valido",
			langDe: "String muss ein gültiger Hostname sein",
		},
		msgValidURI: {
			langEn: msgValidURI,
			langFr: "La chaîne doit être un URI valide",
			langEs: "La cadena debe ser un URI válido",
			langIt: "La stringa deve essere un URI valido",
			langDe: "String muss eine gültige URI sein",
		},
		msgValidURL: {
			langEn: msgValidURL,
			langFr: "La chaîne doit être un URL valide",
			langEs: "La cadena debe ser un URL válido",
			langIt: "La stringa deve essere un URL valido",
			langDe: "String muss eine gültige URL sein",
		},
		msgQueryParamMultiNotAllowed: {
			langEn: msgQueryParamMultiNotAllowed,
			langFr: "Le paramètre de requête ne peut pas être spécifié plus d'une fois",
			langEs: "El parámetro de consulta no se puede especificar más de una vez",
			langIt: "Il parametro di query non può essere specificato più di una volta",
			langDe: "Abfrageparameter dürfen nicht mehrfach angegeben werden",
		},
	},
	Formats: map[string]map[string]string{
		fmtMsgArrayElementType: {
			langEn: fmtMsgArrayElementType,
			langFr: "Les éléments du tableau doivent être de type %[1]s",
			langEs: "Los elementos del arreglo deben ser del tipo %[1]s",
			langIt: "Gli elementi dell'array devono essere di tipo %[1]s",
			langDe: "Array-Elemente müssen vom Typ %[1]s sein",
		},
		fmtMsgArrayElementTypeOrNull: {
			langEn: fmtMsgArrayElementTypeOrNull,
			langFr: "Les éléments du tableau doivent être de type %[1]s ou nuls",
			langEs: "Los elementos del arreglo deben ser del tipo %[1]s o nulo",
			langIt: "Gli elementi dell'array devono essere di tipo %[1]s o null",
			langDe: "Array-Elemente müssen vom Typ %[1]s oder null sein",
		},
		fmtMsgConstraintSetDefaultAllOf: {
			langEn: fmtMsgConstraintSetDefaultAllOf,
			langFr: "L'ensemble de contraintes doit réussir toutes les %[1]d validations non divulguées",
			langEs: "El conjunto de restricciones debe pasar todas las %[1]d validaciones no reveladas",
			langIt: "Il set di vincoli deve superare tutte le %[1]d convalide non divulgate",
			langDe: "Einschränkungssatz muss alle %[1]d nicht offengelegten Validierungen bestehen",
		},
		fmtMsgConstraintSetDefaultOneOf: {
			langEn: fmtMsgConstraintSetDefaultOneOf,
			langFr: "L'ensemble de contraintes doit réussir l'une des %[1]d validations non divulguées",
			langEs: "El conjunto de restricciones debe pasar una de %[1]d validaciones no reveladas",
			langIt: "Il set di vincoli deve superare una delle %[1]d convalide non divulgate",
			langDe: "Einschränkungssatz muss eine von %[1]d nicht offengelegten Validierungen bestehen",
		},
		fmtMsgDtGt: {
			langEn: fmtMsgDtGt,
			langFr: "La valeur doit être après '%[1]s'",
			langEs: "El valor debe estar después de '%[1]s'",
			langIt: "Il valore deve essere successivo a '%[1]s'",
			langDe: "Wert muss nach '%[1]s' liegen",
		},
		fmtMsgDtGte: {
			langEn: fmtMsgDtGte,
			langFr: "La valeur doit être supérieure ou égale à '%[1]s'",
			langEs: "El valor debe ser posterior o igual a '%[1]s'",
			langIt: "Il valore deve essere successivo o uguale a '%[1]s'",
			langDe: "Wert muss nach oder gleich '%[1]s' sein",
		},
		fmtMsgDtLt: {
			langEn: fmtMsgDtLt,
			langFr: "La valeur doit être avant '%[1]s'",
			langEs: "El valor debe estar antes de '%[1]s'",
			langIt: "Il valore deve essere prima di '%[1]s'",
			langDe: "Wert muss vor '%[1]s' liegen",
		},
		fmtMsgDtLte: {
			langEn: fmtMsgDtLte,
			langFr: "La valeur doit être avant ou égale à '%[1]s'",
			langEs: "El valor debe ser anterior o igual a '%[1]s'",
			langIt: "Il valore deve essere prima o uguale a '%[1]s'",
			langDe: "Wert muss vor oder gleich '%[1]s' sein",
		},
		fmtMsgDtToleranceFixedMaxAfter: {
			langEn: fmtMsgDtToleranceFixedMaxAfter,
			langFr: "La valeur ne doit pas être supérieure à %[1]d %[2]s après %[3]s",
			langEs: "El valor no debe ser mayor que %[1]d %[2]s después de %[3]s",
			langIt: "Il valore non deve essere superiore a %[1]d %[2]s dopo %[3]s",
			langDe: "Wert darf nach %[3]s nicht größer als %[1]d %[2]s sein",
		},
		fmtMsgDtToleranceFixedMaxBefore: {
			langEn: fmtMsgDtToleranceFixedMaxBefore,
			langFr: "La valeur ne doit pas être supérieure à %[1]d %[2]s avant %[3]s",
			langEs: "El valor no debe ser mayor que %[1]d %[2]s antes de %[3]s",
			langIt: "Il valore non deve essere superiore a %[1]d %[2]s prima di %[3]s",
			langDe: "Wert darf nicht größer als %[1]d %[2]s vor %[3]s sein",
		},
		fmtMsgDtToleranceFixedMinAfter: {
			langEn: fmtMsgDtToleranceFixedMinAfter,
			langFr: "La valeur doit être au moins %[1]d %[2]s après %[3]s",
			langEs: "El valor debe ser al menos %[1]d %[2]s después de %[3]s",
			langIt: "Il valore deve essere almeno %[1]d %[2]s dopo %[3]s",
			langDe: "Wert muss mindestens %[1]d %[2]s nach %[3]s betragen",
		},
		fmtMsgDtToleranceFixedMinBefore: {
			langEn: fmtMsgDtToleranceFixedMinBefore,
			langFr: "La valeur doit être au moins %[1]d %[2]s avant %[3]s",
			langEs: "El valor debe ser al menos %[1]d %[2]s antes de %[3]s",
			langIt: "Il valore deve essere almeno %[1]d %[2]s prima di %[3]s",
			langDe: "Wert muss mindestens %[1]d %[2]s vor %[3]s betragen",
		},
		fmtMsgDtToleranceFixedSame: {
			langEn: fmtMsgDtToleranceFixedSame,
			langFr: "La valeur doit être la même %[1]s que %[2]s",
			langEs: "El valor debe ser el mismo %[1]s que %[2]s",
			langIt: "Il valore deve essere lo stesso %[1]s di %[2]s",
			langDe: "Wert muss gleich %[1]s wie %[2]s sein",
		},
		fmtMsgDtToleranceNowSame: {
			langEn: fmtMsgDtToleranceNowSame,
			langFr: "La valeur doit être la même %[1]s qu'actuellement",
			langEs: "El valor debe ser el mismo %[1]s que ahora",
			langIt: "Il valore deve essere lo stesso %[1]s di adesso",
			langDe: "Wert muss gleich %[1]s sein wie jetzt",
		},
		fmtMsgDtToleranceNowMaxAfter: {
			langEn: fmtMsgDtToleranceNowMaxAfter,
			langFr: "La valeur ne doit pas dépasser %[1]d %[2]s après maintenant",
			langEs: "El valor no debe ser mayor que %[1]d %[2]s después de ahora",
			langIt: "Il valore non deve essere superiore a %[1]d %[2]s dopo ora",
			langDe: "Wert darf nach jetzt nicht mehr als %[1]d %[2]s betragen",
		},
		fmtMsgDtToleranceNowMaxBefore: {
			langEn: fmtMsgDtToleranceNowMaxBefore,
			langFr: "La valeur ne doit pas dépasser %[1]d %[2]s avant maintenant",
			langEs: "El valor no debe ser superior a %[1]d %[2]s antes de ahora",
			langIt: "Il valore non deve essere superiore a %[1]d %[2]s prima di ora",
			langDe: "Wert darf bis jetzt nicht größer als %[1]d %[2]s sein",
		},
		fmtMsgDtToleranceNowMinAfter: {
			langEn: fmtMsgDtToleranceNowMinAfter,
			langFr: "La valeur doit être au moins %[1]d %[2]s après maintenant",
			langEs: "El valor debe ser al menos %[1]d %[2]s después de ahora",
			langIt: "Il valore deve essere almeno %[1]d %[2]s dopo ora",
			langDe: "Wert muss nach jetzt mindestens %[1]d %[2]s betragen",
		},
		fmtMsgDtToleranceNowMinBefore: {
			langEn: fmtMsgDtToleranceNowMinBefore,
			langFr: "La valeur doit être au moins %[1]d %[2]s avant maintenant",
			langEs: "El valor debe ser al menos %[1]d %[2]s antes de ahora",
			langIt: "Il valore deve essere almeno %[1]d %[2]s prima di ora",
			langDe: "Wert muss vorher mindestens %[1]d %[2]s betragen",
		},
		fmtMsgDtToleranceOtherSame: {
			langEn: fmtMsgDtToleranceOtherSame,
			langFr: "La valeur doit être la même %[1]s que la valeur de la propriété '%[2]s'",
			langEs: "El valor debe ser el mismo %[1]s que el valor de la propiedad '%[2]s'",
			langIt: "Il valore deve essere lo stesso %[1]s del valore della proprietà '%[2]s'",
			langDe: "Wert muss gleich %[1]s sein wie Wert der Eigenschaft '%[2]s'",
		},
		fmtMsgDtToleranceOtherMaxAfter: {
			langEn: fmtMsgDtToleranceOtherMaxAfter,
			langFr: "La valeur ne doit pas être supérieure à %[1]d %[2]s après la valeur de la propriété '%[3]s'",
			langEs: "El valor no debe ser mayor que %[1]d %[2]s después del valor de la propiedad '%[3]s'",
			langIt: "Il valore non deve essere superiore a %[1]d %[2]s dopo il valore della proprietà '%[3]s'",
			langDe: "Wert darf nicht größer als %[1]d %[2]s nach dem Wert der Eigenschaft '%[3]s' sein",
		},
		fmtMsgDtToleranceOtherMaxBefore: {
			langEn: fmtMsgDtToleranceOtherMaxBefore,
			langFr: "La valeur ne doit pas être supérieure à %[1]d %[2]s avant la valeur de la propriété '%[3]s'",
			langEs: "El valor no debe ser mayor que %[1]d %[2]s antes del valor de la propiedad '%[3]s'",
			langIt: "Il valore non deve essere superiore a %[1]d %[2]s prima del valore della proprietà '%[3]s'",
			langDe: "Wert darf nicht größer als %[1]d %[2]s vor dem Wert der Eigenschaft '%[3]s' sein",
		},
		fmtMsgDtToleranceOtherMinAfter: {
			langEn: fmtMsgDtToleranceOtherMinAfter,
			langFr: "La valeur doit être au moins %[1]d %[2]s après la valeur de la propriété '%[3]s'",
			langEs: "El valor debe ser al menos %[1]d %[2]s después del valor de la propiedad '%[3]s'",
			langIt: "Il valore deve essere almeno %[1]d %[2]s dopo il valore della proprietà '%[3]s'",
			langDe: "Wert muss mindestens %[1]d %[2]s nach dem Wert der Eigenschaft '%[3]s' betragen",
		},
		fmtMsgDtToleranceOtherMinBefore: {
			langEn: fmtMsgDtToleranceOtherMinBefore,
			langFr: "La valeur doit être au moins %[1]d %[2]s avant la valeur de la propriété '%[3]s'",
			langEs: "El valor debe ser al menos %[1]d %[2]s antes del valor de la propiedad '%[3]s'",
			langIt: "Il valore deve essere almeno %[1]d %[2]s prima del valore della proprietà '%[3]s'",
			langDe: "Wert muss mindestens %[1]d %[2]s vor dem Wert der Eigenschaft '%[3]s' liegen",
		},
		fmtMsgEqualsOther: {
			langEn: fmtMsgEqualsOther,
			langFr: "La valeur doit être égale à la valeur de la propriété '%[1]s'",
			langEs: "El valor debe ser igual al valor de la propiedad '%[1]s'",
			langIt: "Il valore deve essere uguale al valore della proprietà '%[1]s'",
			langDe: "Wert muss gleich dem Wert der Eigenschaft '%[1]s' sein",
		},
		fmtMsgExactLen: {
			langEn: fmtMsgExactLen,
			langFr: "La longueur de la valeur doit être %[1]d",
			langEs: "La longitud del valor debe ser %[1]d",
			langIt: "La lunghezza del valore deve essere %[1]d",
			langDe: "Wertlänge muss %[1]d sein",
		},
		fmtMsgGt: {
			langEn: fmtMsgGt,
			langFr: "La valeur doit être supérieure à %[1]v",
			langEs: "El valor debe ser mayor que %[1]v",
			langIt: "Il valore deve essere maggiore di %[1]v",
			langDe: "Wert muss größer als %[1]v sein",
		},
		fmtMsgGte: {
			langEn: fmtMsgGte,
			langFr: "La valeur doit être supérieure ou égale à %[1]v",
			langEs: "El valor debe ser mayor o igual que %[1]v",
			langIt: "Il valore deve essere maggiore o uguale a %[1]v",
			langDe: "Wert muss größer oder gleich %[1]v sein",
		},
		fmtMsgGtOther: {
			langEn: fmtMsgGtOther,
			langFr: "La valeur doit être supérieure à la valeur de la propriété '%[1]s'",
			langEs: "El valor debe ser mayor que el valor de la propiedad '%[1]s'",
			langIt: "Il valore deve essere maggiore del valore della proprietà '%[1]s'",
			langDe: "Wert muss größer sein als Wert der Eigenschaft '%[1]s'",
		},
		fmtMsgGteOther: {
			langEn: fmtMsgGteOther,
			langFr: "La valeur doit être supérieure ou égale à la valeur de la propriété '%[1]s'",
			langEs: "El valor debe ser mayor o igual que el valor de la propiedad '%[1]s'",
			langIt: "Il valore deve essere maggiore o uguale al valore della proprietà '%[1]s'",
			langDe: "Wert muss größer oder gleich dem Wert der Eigenschaft '%[1]s' sein",
		},
		fmtMsgLt: {
			langEn: fmtMsgLt,
			langFr: "La valeur doit être inférieure à %[1]v",
			langEs: "El valor debe ser menor que %[1]v",
			langIt: "Il valore deve essere inferiore a %[1]v",
			langDe: "Wert muss kleiner als %[1]v sein",
		},
		fmtMsgLte: {
			langEn: fmtMsgLte,
			langFr: "La valeur doit être inférieure ou égale à %[1]v",
			langEs: "El valor debe ser menor o igual que %[1]v",
			langIt: "Il valore deve essere inferiore o uguale a %[1]v",
			langDe: "Wert muss kleiner oder gleich %[1]v sein",
		},
		fmtMsgLtOther: {
			langEn: fmtMsgLtOther,
			langFr: "La valeur doit être inférieure à la valeur de la propriété '%[1]s'",
			langEs: "El valor debe ser menor que el valor de la propiedad '%[1]s'",
			langIt: "Il valore deve essere inferiore al valore della proprietà '%[1]s'",
			langDe: "Wert muss kleiner sein als Wert der Eigenschaft '%[1]s'",
		},
		fmtMsgLteOther: {
			langEn: fmtMsgLteOther,
			langFr: "La valeur doit être inférieure ou égale à la valeur de la propriété '%[1]s'",
			langEs: "El valor debe ser menor o igual que el valor de la propiedad '%[1]s'",
			langIt: "Il valore deve essere inferiore o uguale al valore della proprietà '%[1]s'",
			langDe: "Wert muss kleiner oder gleich dem Wert der Eigenschaft '%[1]s' sein",
		},
		fmtMsgMinLen: {
			langEn: fmtMsgMinLen,
			langFr: "La longueur de la valeur doit être d'au moins %[1]d",
			langEs: "La longitud del valor debe ser al menos %[1]d",
			langIt: "La lunghezza del valore deve essere almeno %[1]d",
			langDe: "Wertlänge muss mindestens %[1]d betragen",
		},
		fmtMsgMinLenExc: {
			langEn: fmtMsgMinLenExc,
			langFr: "La longueur de la valeur doit être supérieure à %[1]d",
			langEs: "La longitud del valor debe ser mayor que %[1]d",
			langIt: "La lunghezza del valore deve essere maggiore di %[1]d",
			langDe: "Wertlänge muss größer sein als %[1]d",
		},
		fmtMsgMinMax: {
			langEn: fmtMsgMinMax,
			langFr: "La longueur de la valeur doit être comprise entre %[1]d (%[2]s) et %[3]d (%[4]s)",
			langEs: "La longitud del valor debe estar entre %[1]d (%[2]s) y %[3]d (%[4]s)",
			langIt: "La lunghezza del valore deve essere compresa tra %[1]d (%[2]s) e %[3]d (%[4]s)",
			langDe: "Wertlänge muss zwischen %[1]d (%[2]s) und %[3]d (%[4]s) liegen",
		},
		fmtMsgMultipleOf: {
			langEn: fmtMsgMultipleOf,
			langFr: "La valeur doit être un multiple de %[1]d",
			langEs: "El valor debe ser un múltiplo de %[1]d",
			langIt: "Il valore deve essere un multiplo di %[1]d",
			langDe: "Wert muss ein Vielfaches von %[1]d sein",
		},
		fmtMsgNotEqualsOther: {
			langEn: fmtMsgNotEqualsOther,
			langFr: "La valeur ne doit pas être égale à la valeur de la propriété '%[1]s'",
			langEs: "El valor no debe ser igual al valor de la propiedad '%[1]s'",
			langIt: "Il valore non deve essere uguale al valore della proprietà '%[1]s'",
			langDe: "Wert darf nicht gleich dem Wert der Eigenschaft '%[1]s' sein",
		},
		fmtMsgRange: {
			langEn: fmtMsgRange,
			langFr: "La valeur doit être comprise entre %[1]v (%[2]s) et %[3]v (%[4]s)",
			langEs: "El valor debe estar entre %[1]v (%[2]s) y %[3]v (%[4]s)",
			langIt: "Il valore deve essere compreso tra %[1]v (%[2]s) e %[3]v (%[4]s)",
			langDe: "Wert muss zwischen %[1]v (%[2]s) und %[3]v (%[4]s) liegen",
		},
		fmtMsgStringExactLen: {
			langEn: fmtMsgStringExactLen,
			langFr: "La longueur de la valeur de la chaîne doit être de %[1]d caractères",
			langEs: "La longitud del valor de la cadena debe ser %[1]d caracteres",
			langIt: "La lunghezza del valore della stringa deve essere %[1]d caratteri",
			langDe: "String-Wert muss %[1]d Zeichen lang sein",
		},
		fmtMsgStringMaxLen: {
			langEn: fmtMsgStringMaxLen,
			langFr: "La longueur de la valeur de la chaîne ne doit pas dépasser %[1]d caractères",
			langEs: "La longitud del valor de la cadena no debe exceder %[1]d caracteres",
			langIt: "La lunghezza del valore della stringa non deve superare %[1]d caratteri",
			langDe: "Stringwertlänge darf %[1]d Zeichen nicht überschreiten",
		},
		fmtMsgStringMaxLenExc: {
			langEn: fmtMsgStringMaxLenExc,
			langFr: "La longueur de la valeur de la chaîne doit être inférieure à %[1]d caractères",
			langEs: "La longitud del valor de la cadena debe ser inferior a %[1]d caracteres",
			langIt: "La lunghezza del valore della stringa deve essere inferiore a %[1]d caratteri",
			langDe: "String-Wert muss weniger als %[1]d Zeichen lang sein",
		},
		fmtMsgStringMinLen: {
			langEn: fmtMsgStringMinLen,
			langFr: "La longueur de la valeur de la chaîne doit être d'au moins %[1]d caractères",
			langEs: "La longitud del valor de la cadena debe ser de al menos %[1]d caracteres",
			langIt: "La lunghezza del valore della stringa deve essere di almeno %[1]d caratteri",
			langDe: "String-Wert muss mindestens %[1]d Zeichen lang sein",
		},
		fmtMsgStringMinLenExc: {
			langEn: fmtMsgStringMinLenExc,
			langFr: "La longueur de la valeur de la chaîne doit être supérieure à %[1]d caractères",
			langEs: "La longitud del valor de la cadena debe ser mayor que %[1]d caracteres",
			langIt: "La lunghezza del valore della stringa deve essere maggiore di %[1]d caratteri",
			langDe: "Stringwertlänge muss größer als %[1]d Zeichen sein",
		},
		fmtMsgStringMinMaxLen: {
			langEn: fmtMsgStringMinMaxLen,
			langFr: "La longueur de la valeur de chaîne doit être comprise entre %[1]d (%[2]s) et %[3]d (%[4]s)",
			langEs: "La longitud del valor de la cadena debe estar entre %[1]d (%[2]s) y %[3]d (%[4]s)",
			langIt: "La lunghezza del valore della stringa deve essere compresa tra %[1]d (%[2]s) e %[3]d (%[4]s)",
			langDe: "Stringwertlänge muss zwischen %[1]d (%[2]s) und %[3]d (%[4]s) liegen",
		},
		fmtMsgUnknownPresetPattern: {
			langEn: fmtMsgUnknownPresetPattern,
			langFr: "Modèle prédéfini inconnu '%[1]s'",
			langEs: "Patrón predeterminado desconocido '%[1]s'",
			langIt: "Modello predefinito sconosciuto '%[1]s'",
			langDe: "Unbekanntes voreingestelltes Muster '%[1]s'",
		},
		fmtMsgUuidCorrectVer: {
			langEn: fmtMsgUuidCorrectVer,
			langFr: "La valeur doit être un UUID valide (version %[1]d)",
			langEs: "El valor debe ser un UUID válido (versión %[1]d)",
			langIt: "Il valore deve essere un UUID valido (versione %[1]d)",
			langDe: "Wert muss eine gültige UUID sein (Version %[1]d)",
		},
		fmtMsgUuidMinVersion: {
			langEn: fmtMsgUuidMinVersion,
			langFr: "La valeur doit être un UUID valide (version minimale %[1]d)",
			langEs: "El valor debe ser un UUID válido (versión mínima %[1]d)",
			langIt: "Il valore deve essere un UUID valido (versione minima %[1]d)",
			langDe: "Wert muss eine gültige UUID sein (Mindestversion %[1]d)",
		},
		fmtMsgValidToken: {
			langEn: fmtMsgValidToken,
			langFr: "La valeur de la chaîne doit être un jeton valide - %[1]s",
			langEs: "El valor de la cadena debe ser un token válido - %[1]s",
			langIt: "Il valore della stringa deve essere un token valido - %[1]s",
			langDe: "String-Wert muss gültiges Token sein - %[1]s",
		},
		fmtMsgValueExpectedType: {
			langEn: fmtMsgValueExpectedType,
			langFr: "Valeur supposée être de type %[1]s",
			langEs: "Se espera que el valor sea del tipo %[1]s",
			langIt: "Valore previsto di tipo %[1]s",
			langDe: "Wert sollte vom Typ %[1]s sein",
		},
		fmtMsgStringContains: {
			langEn: fmtMsgStringContains,
			langFr: "La chaîne doit contenir %[1]s",
			langEs: "La cadena debe contener %[1]s",
			langIt: "La stringa deve contenere %[1]s",
			langDe: "Zeichenfolge muss %[1]s enthalten",
		},
		fmtMsgStringNotContains: {
			langEn: fmtMsgStringNotContains,
			langFr: "La chaîne ne doit pas contenir %[1]s",
			langEs: "La cadena no debe contener %[1]s",
			langIt: "La stringa non deve contenere %[1]s",
			langDe: "Zeichenfolge darf %[1]s nicht enthalten",
		},
		fmtMsgStringStartsWith: {
			langEn: fmtMsgStringStartsWith,
			langFr: "La valeur de chaîne doit commencer par %[1]s",
			langEs: "El valor de cadena debe comenzar con %[1]s",
			langIt: "Il valore della stringa deve iniziare con %[1]s",
			langDe: "Zeichenfolgenwert muss mit %[1]s beginnen",
		},
		fmtMsgStringNotStartsWith: {
			langEn: fmtMsgStringNotStartsWith,
			langFr: "La valeur de chaîne ne doit pas commencer par %[1]s",
			langEs: "El valor de cadena no debe comenzar con %[1]s",
			langIt: "Il valore della stringa non deve iniziare con %[1]s",
			langDe: "Zeichenfolgenwert darf nicht mit %[1]s beginnen",
		},
		fmtMsgStringEndsWith: {
			langEn: fmtMsgStringEndsWith,
			langFr: "La valeur de chaîne doit se terminer par %[1]s",
			langEs: "El valor de cadena debe terminar con %[1]s",
			langIt: "Il valore della stringa deve terminare con %[1]s",
			langDe: "Zeichenfolgenwert muss mit %[1]s enden",
		},
		fmtMsgStringNotEndsWith: {
			langEn: fmtMsgStringNotEndsWith,
			langFr: "La valeur de chaîne ne doit pas se terminer par %[1]s",
			langEs: "El valor de cadena no debe terminar con %[1]s",
			langIt: "Il valore della stringa non deve terminare con %[1]s",
			langDe: "Zeichenfolgenwert darf nicht mit %[1]s enden",
		},
		fmtMsgQueryParamType: {
			langEn: fmtMsgQueryParamType,
			langFr: "Le paramètre de requête doit être de type %[1]s",
			langEs: "El parámetro de consulta debe ser del tipo %[1]s",
			langIt: "Il parametro della query deve essere di tipo %[1]s",
			langDe: "Der Abfrageparameter muss vom Typ %[1]s sein",
		},
	},
}
