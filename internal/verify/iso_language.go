// Copyright © 2025 Ping Identity Corporation

// Package verify provides validation utilities and constants for the PingOne Terraform provider.
package verify

import (
	"fmt"
	"slices"
	"strings"
)

// reservedLanguageCodes contains language codes that are reserved by PingOne and cannot be customized
var reservedLanguageCodes = []string{
	"cs",
	"de",
	"en",
	"es",
	"fr-CA",
	"fr",
	"hu",
	"it",
	"ja",
	"ko",
	"nl",
	"pl",
	"pt",
	"ru",
	"th",
	"tr",
	"zh",
}

// IsoCountry represents an ISO language/country code with its descriptive name
type IsoCountry struct {
	// Code is the ISO language or locale code (e.g., "en", "en-US")
	Code string
	// Name is the human-readable name for the language or locale
	Name string
}

var isoList = []IsoCountry{
	{
		Code: "aa",
		Name: "Afar",
	},
	{
		Code: "ab",
		Name: "Abkhazian",
	},
	{
		Code: "ae",
		Name: "Avestan",
	},
	{
		Code: "af",
		Name: "Afrikaans",
	},
	{
		Code: "af-ZA",
		Name: "Afrikaans (South Africa)",
	},
	{
		Code: "ak",
		Name: "Akan",
	},
	{
		Code: "am",
		Name: "Amharic",
	},
	{
		Code: "an",
		Name: "Aragonese",
	},
	{
		Code: "ar",
		Name: "Arabic",
	},
	{
		Code: "ar-AE",
		Name: "Arabic (U.A.E.)",
	},
	{
		Code: "ar-BH",
		Name: "Arabic (Bahrain)",
	},
	{
		Code: "ar-DZ",
		Name: "Arabic (Algeria)",
	},
	{
		Code: "ar-EG",
		Name: "Arabic (Egypt)",
	},
	{
		Code: "ar-IQ",
		Name: "Arabic (Iraq)",
	},
	{
		Code: "ar-JO",
		Name: "Arabic (Jordan)",
	},
	{
		Code: "ar-KW",
		Name: "Arabic (Kuwait)",
	},
	{
		Code: "ar-LB",
		Name: "Arabic (Lebanon)",
	},
	{
		Code: "ar-LY",
		Name: "Arabic (Libya)",
	},
	{
		Code: "ar-MA",
		Name: "Arabic (Morocco)",
	},
	{
		Code: "ar-OM",
		Name: "Arabic (Oman)",
	},
	{
		Code: "ar-QA",
		Name: "Arabic (Qatar)",
	},
	{
		Code: "ar-SA",
		Name: "Arabic (Saudi Arabia)",
	},
	{
		Code: "ar-SY",
		Name: "Arabic (Syria)",
	},
	{
		Code: "ar-TN",
		Name: "Arabic (Tunisia)",
	},
	{
		Code: "ar-YE",
		Name: "Arabic (Yemen)",
	},
	{
		Code: "as",
		Name: "Assamese",
	},
	{
		Code: "av",
		Name: "Avaric",
	},
	{
		Code: "ay",
		Name: "Aymara",
	},
	{
		Code: "az",
		Name: "Azeri (Latin)",
	},
	{
		Code: "az-AZ",
		Name: "Azeri (Latin) (Azerbaijan)",
	},
	{
		Code: "ba",
		Name: "Bashkir",
	},
	{
		Code: "be",
		Name: "Belarusian",
	},
	{
		Code: "be-BY",
		Name: "Belarusian (Belarus)",
	},
	{
		Code: "bg",
		Name: "Bulgarian",
	},
	{
		Code: "bg-BG",
		Name: "Bulgarian (Bulgaria)",
	},
	{
		Code: "bi",
		Name: "Bislama",
	},
	{
		Code: "bm",
		Name: "Bambara",
	},
	{
		Code: "bn",
		Name: "Bengali",
	},
	{
		Code: "bo",
		Name: "Tibetan",
	},
	{
		Code: "br",
		Name: "Breton",
	},
	{
		Code: "bs",
		Name: "Bosnian",
	},
	{
		Code: "bs-BA",
		Name: "Bosnian (Bosnia and Herzegovina)",
	},
	{
		Code: "ca",
		Name: "Catalan",
	},
	{
		Code: "ca-ES",
		Name: "Catalan (Spain)",
	},
	{
		Code: "ce",
		Name: "Chechen",
	},
	{
		Code: "ch",
		Name: "Chamorro",
	},
	{
		Code: "cmn-CN",
		Name: "Chinese",
	},
	{
		Code: "cmn-TW",
		Name: "Chinese (Taiwan)",
	},
	{
		Code: "co",
		Name: "Corsican",
	},
	{
		Code: "cr",
		Name: "Cree",
	},
	{
		Code: "cs",
		Name: "Czech",
	},
	{
		Code: "cs-CZ",
		Name: "Czech (Czech Republic)",
	},
	{
		Code: "cu",
		Name: "Church Slavonic",
	},
	{
		Code: "cv",
		Name: "Chuvash",
	},
	{
		Code: "cy",
		Name: "Welsh",
	},
	{
		Code: "cy-GB",
		Name: "Welsh (United Kingdom)",
	},
	{
		Code: "da",
		Name: "Danish",
	},
	{
		Code: "da-DK",
		Name: "Danish (Denmark)",
	},
	{
		Code: "de",
		Name: "German",
	},
	{
		Code: "de-AT",
		Name: "German (Austria)",
	},
	{
		Code: "de-CH",
		Name: "German (Switzerland)",
	},
	{
		Code: "de-DE",
		Name: "German (Germany)",
	},
	{
		Code: "de-LI",
		Name: "German (Liechtenstein)",
	},
	{
		Code: "de-LU",
		Name: "German (Luxembourg)",
	},
	{
		Code: "dv",
		Name: "Divehi",
	},
	{
		Code: "dv-MV",
		Name: "Divehi (Maldives)",
	},
	{
		Code: "dz",
		Name: "Dzongkha",
	},
	{
		Code: "ee",
		Name: "Ewe",
	},
	{
		Code: "el",
		Name: "Greek",
	},
	{
		Code: "el-GR",
		Name: "Greek (Greece)",
	},
	{
		Code: "en",
		Name: "English",
	},
	{
		Code: "en-AU",
		Name: "English (Australia)",
	},
	{
		Code: "en-BZ",
		Name: "English (Belize)",
	},
	{
		Code: "en-CA",
		Name: "English (Canada)",
	},
	{
		Code: "en-CB",
		Name: "English (Caribbean)",
	},
	{
		Code: "en-GB",
		Name: "English (United Kingdom)",
	},
	{
		Code: "en-GB-WLS",
		Name: "English (Welsh)",
	},
	{
		Code: "en-IE",
		Name: "English (Ireland)",
	},
	{
		Code: "en-IN",
		Name: "English (Indian)",
	},
	{
		Code: "en-JM",
		Name: "English (Jamaica)",
	},
	{
		Code: "en-NZ",
		Name: "English (New Zealand)",
	},
	{
		Code: "en-PH",
		Name: "English (Republic of the Philippines)",
	},
	{
		Code: "en-TT",
		Name: "English (Trinidad and Tobago)",
	},
	{
		Code: "en-US",
		Name: "English (United States)",
	},
	{
		Code: "en-ZA",
		Name: "English (South Africa)",
	},
	{
		Code: "en-ZW",
		Name: "English (Zimbabwe)",
	},
	{
		Code: "eo",
		Name: "Esperanto",
	},
	{
		Code: "es",
		Name: "Spanish",
	},
	{
		Code: "es-AR",
		Name: "Spanish (Argentina)",
	},
	{
		Code: "es-BO",
		Name: "Spanish (Bolivia)",
	},
	{
		Code: "es-CL",
		Name: "Spanish (Chile)",
	},
	{
		Code: "es-CO",
		Name: "Spanish (Colombia)",
	},
	{
		Code: "es-CR",
		Name: "Spanish (Costa Rica)",
	},
	{
		Code: "es-DO",
		Name: "Spanish (Dominican Republic)",
	},
	{
		Code: "es-EC",
		Name: "Spanish (Ecuador)",
	},
	{
		Code: "es-ES",
		Name: "Spanish (Spain)",
	},
	{
		Code: "es-GT",
		Name: "Spanish (Guatemala)",
	},
	{
		Code: "es-HN",
		Name: "Spanish (Honduras)",
	},
	{
		Code: "es-MX",
		Name: "Spanish (Mexico)",
	},
	{
		Code: "es-NI",
		Name: "Spanish (Nicaragua)",
	},
	{
		Code: "es-PA",
		Name: "Spanish (Panama)",
	},
	{
		Code: "es-PE",
		Name: "Spanish (Peru)",
	},
	{
		Code: "es-PR",
		Name: "Spanish (Puerto Rico)",
	},
	{
		Code: "es-PY",
		Name: "Spanish (Paraguay)",
	},
	{
		Code: "es-SV",
		Name: "Spanish (El Salvador)",
	},
	{
		Code: "es-US",
		Name: "Spanish (United States)",
	},
	{
		Code: "es-UY",
		Name: "Spanish (Uruguay)",
	},
	{
		Code: "es-VE",
		Name: "Spanish (Venezuela)",
	},
	{
		Code: "et",
		Name: "Estonian",
	},
	{
		Code: "et-EE",
		Name: "Estonian (Estonia)",
	},
	{
		Code: "eu",
		Name: "Basque",
	},
	{
		Code: "eu-ES",
		Name: "Basque (Spain)",
	},
	{
		Code: "fa",
		Name: "Farsi",
	},
	{
		Code: "fa-IR",
		Name: "Farsi (Iran)",
	},
	{
		Code: "ff",
		Name: "Fulah",
	},
	{
		Code: "fi",
		Name: "Finnish",
	},
	{
		Code: "fi-FI",
		Name: "Finnish (Finland)",
	},
	{
		Code: "fj",
		Name: "Fijian",
	},
	{
		Code: "fo",
		Name: "Faroese",
	},
	{
		Code: "fo-FO",
		Name: "Faroese (Faroe Islands)",
	},
	{
		Code: "fr",
		Name: "French",
	},
	{
		Code: "fr-BE",
		Name: "French (Belgium)",
	},
	{
		Code: "fr-CA",
		Name: "French (Canada)",
	},
	{
		Code: "fr-CH",
		Name: "French (Switzerland)",
	},
	{
		Code: "fr-FR",
		Name: "French (France)",
	},
	{
		Code: "fr-LU",
		Name: "French (Luxembourg)",
	},
	{
		Code: "fr-MC",
		Name: "French (Principality of Monaco)",
	},
	{
		Code: "fy",
		Name: "Western Frisian",
	},
	{
		Code: "ga",
		Name: "Irish",
	},
	{
		Code: "gd",
		Name: "Gaelic",
	},
	{
		Code: "gl",
		Name: "Galician",
	},
	{
		Code: "gl-ES",
		Name: "Galician (Spain)",
	},
	{
		Code: "gn",
		Name: "Guarani",
	},
	{
		Code: "gu",
		Name: "Gujarati",
	},
	{
		Code: "gu-IN",
		Name: "Gujarati (India)",
	},
	{
		Code: "gv",
		Name: "Manx",
	},
	{
		Code: "ha",
		Name: "Hausa",
	},
	{
		Code: "he",
		Name: "Hebrew",
	},
	{
		Code: "he-IL",
		Name: "Hebrew (Israel)",
	},
	{
		Code: "hi",
		Name: "Hindi",
	},
	{
		Code: "hi-IN",
		Name: "Hindi (India)",
	},
	{
		Code: "ho",
		Name: "Hiri Motu",
	},
	{
		Code: "hr",
		Name: "Croatian",
	},
	{
		Code: "hr-BA",
		Name: "Croatian (Bosnia and Herzegovina)",
	},
	{
		Code: "hr-HR",
		Name: "Croatian (Croatia)",
	},
	{
		Code: "ht",
		Name: "Haitian",
	},
	{
		Code: "hu",
		Name: "Hungarian",
	},
	{
		Code: "hu-HU",
		Name: "Hungarian (Hungary)",
	},
	{
		Code: "hy",
		Name: "Armenian",
	},
	{
		Code: "hy-AM",
		Name: "Armenian (Armenia)",
	},
	{
		Code: "hz",
		Name: "Herero",
	},
	{
		Code: "ia",
		Name: "Interlingua (International Auxiliary Language Association)",
	},
	{
		Code: "id",
		Name: "Indonesian",
	},
	{
		Code: "id-ID",
		Name: "Indonesian (Indonesia)",
	},
	{
		Code: "ie",
		Name: "Interlingue",
	},
	{
		Code: "ig",
		Name: "Igbo",
	},
	{
		Code: "ii",
		Name: "Sichuan Yi",
	},
	{
		Code: "ik",
		Name: "Inupiaq",
	},
	{
		Code: "io",
		Name: "Ido",
	},
	{
		Code: "is",
		Name: "Icelandic",
	},
	{
		Code: "is-IS",
		Name: "Icelandic (Iceland)",
	},
	{
		Code: "it",
		Name: "Italian",
	},
	{
		Code: "it-CH",
		Name: "Italian (Switzerland)",
	},
	{
		Code: "it-IT",
		Name: "Italian (Italy)",
	},
	{
		Code: "iu",
		Name: "Inuktitut",
	},
	{
		Code: "ja",
		Name: "Japanese",
	},
	{
		Code: "ja-JP",
		Name: "Japanese (Japan)",
	},
	{
		Code: "jv",
		Name: "Javanese",
	},
	{
		Code: "ka",
		Name: "Georgian",
	},
	{
		Code: "ka-GE",
		Name: "Georgian (Georgia)",
	},
	{
		Code: "kg",
		Name: "Kongo",
	},
	{
		Code: "ki",
		Name: "Kikuyu",
	},
	{
		Code: "kj",
		Name: "Kuanyama",
	},
	{
		Code: "kk",
		Name: "Kazakh",
	},
	{
		Code: "kk-KZ",
		Name: "Kazakh (Kazakhstan)",
	},
	{
		Code: "kl",
		Name: "Kalaallisut",
	},
	{
		Code: "km",
		Name: "Central Khmer",
	},
	{
		Code: "kn",
		Name: "Kannada",
	},
	{
		Code: "kn-IN",
		Name: "Kannada (India)",
	},
	{
		Code: "ko",
		Name: "Korean",
	},
	{
		Code: "ko-KR",
		Name: "Korean (Korea)",
	},
	{
		Code: "kok",
		Name: "Konkani",
	},
	{
		Code: "kok-IN",
		Name: "Konkani (India)",
	},
	{
		Code: "kr",
		Name: "Kanuri",
	},
	{
		Code: "ks",
		Name: "Kashmiri",
	},
	{
		Code: "ku",
		Name: "Kurdish",
	},
	{
		Code: "kv",
		Name: "Komi",
	},
	{
		Code: "kw",
		Name: "Cornish",
	},
	{
		Code: "ky",
		Name: "Kyrgyz",
	},
	{
		Code: "ky-KG",
		Name: "Kyrgyz (Kyrgyzstan)",
	},
	{
		Code: "la",
		Name: "Latin",
	},
	{
		Code: "lb",
		Name: "Luxembourgish",
	},
	{
		Code: "lg",
		Name: "Ganda",
	},
	{
		Code: "li",
		Name: "Limburgan",
	},
	{
		Code: "ln",
		Name: "Lingala",
	},
	{
		Code: "lo",
		Name: "Lao",
	},
	{
		Code: "lt",
		Name: "Lithuanian",
	},
	{
		Code: "lt-LT",
		Name: "Lithuanian (Lithuania)",
	},
	{
		Code: "lu",
		Name: "Luba-Katanga",
	},
	{
		Code: "lv",
		Name: "Latvian",
	},
	{
		Code: "lv-LV",
		Name: "Latvian (Latvia)",
	},
	{
		Code: "mg",
		Name: "Malagasy",
	},
	{
		Code: "mh",
		Name: "Marshallese",
	},
	{
		Code: "mi",
		Name: "Maori",
	},
	{
		Code: "mi-NZ",
		Name: "Maori (New Zealand)",
	},
	{
		Code: "mk",
		Name: "FYRO Macedonian",
	},
	{
		Code: "mk-MK",
		Name: "FYRO Macedonian (Former Yugoslav Republic of Macedonia)",
	},
	{
		Code: "ml",
		Name: "Malayalam",
	},
	{
		Code: "mn",
		Name: "Mongolian",
	},
	{
		Code: "mn-MN",
		Name: "Mongolian (Mongolia)",
	},
	{
		Code: "mr",
		Name: "Marathi",
	},
	{
		Code: "mr-IN",
		Name: "Marathi (India)",
	},
	{
		Code: "ms",
		Name: "Malay",
	},
	{
		Code: "ms-BN",
		Name: "Malay (Brunei Darussalam)",
	},
	{
		Code: "ms-MY",
		Name: "Malay (Malaysia)",
	},
	{
		Code: "mt",
		Name: "Maltese",
	},
	{
		Code: "mt-MT",
		Name: "Maltese (Malta)",
	},
	{
		Code: "my",
		Name: "Burmese",
	},
	{
		Code: "na",
		Name: "Nauru",
	},
	{
		Code: "nb",
		Name: "Norwegian (Bokmål)",
	},
	{
		Code: "nb-NO",
		Name: "Norwegian (Bokmål) (Norway)",
	},
	{
		Code: "nd",
		Name: "North Ndebele",
	},
	{
		Code: "ne",
		Name: "Nepali",
	},
	{
		Code: "ng",
		Name: "Ndonga",
	},
	{
		Code: "nl",
		Name: "Dutch",
	},
	{
		Code: "nl-BE",
		Name: "Dutch (Belgium)",
	},
	{
		Code: "nl-NL",
		Name: "Dutch (Netherlands)",
	},
	{
		Code: "nn",
		Name: "Norwegian Nynorsk",
	},
	{
		Code: "nn-NO",
		Name: "Norwegian (Nynorsk) (Norway)",
	},
	{
		Code: "no",
		Name: "Norwegian",
	},
	{
		Code: "nr",
		Name: "South Ndebele",
	},
	{
		Code: "ns",
		Name: "Northern Sotho",
	},
	{
		Code: "ns-ZA",
		Name: "Northern Sotho (South Africa)",
	},
	{
		Code: "nv",
		Name: "Navajo",
	},
	{
		Code: "ny",
		Name: "Chichewa",
	},
	{
		Code: "oc",
		Name: "Occitan",
	},
	{
		Code: "oj",
		Name: "Ojibwa",
	},
	{
		Code: "om",
		Name: "Oromo",
	},
	{
		Code: "or",
		Name: "Oriya",
	},
	{
		Code: "os",
		Name: "Ossetian",
	},
	{
		Code: "pa",
		Name: "Punjabi",
	},
	{
		Code: "pa-IN",
		Name: "Punjabi (India)",
	},
	{
		Code: "pi",
		Name: "Pali",
	},
	{
		Code: "pl",
		Name: "Polish",
	},
	{
		Code: "pl-PL",
		Name: "Polish (Poland)",
	},
	{
		Code: "ps",
		Name: "Pashto",
	},
	{
		Code: "ps-AR",
		Name: "Pashto (Afghanistan)",
	},
	{
		Code: "pt",
		Name: "Portuguese",
	},
	{
		Code: "pt-BR",
		Name: "Portuguese (Brazil)",
	},
	{
		Code: "pt-PT",
		Name: "Portuguese (Portugal)",
	},
	{
		Code: "qu",
		Name: "Quechua",
	},
	{
		Code: "qu-BO",
		Name: "Quechua (Bolivia)",
	},
	{
		Code: "qu-EC",
		Name: "Quechua (Ecuador)",
	},
	{
		Code: "qu-PE",
		Name: "Quechua (Peru)",
	},
	{
		Code: "rm",
		Name: "Romansh",
	},
	{
		Code: "rn",
		Name: "Rundi",
	},
	{
		Code: "ro",
		Name: "Romanian",
	},
	{
		Code: "ro-RO",
		Name: "Romanian (Romania)",
	},
	{
		Code: "ru",
		Name: "Russian",
	},
	{
		Code: "ru-RU",
		Name: "Russian (Russia)",
	},
	{
		Code: "rw",
		Name: "Kinyarwanda",
	},
	{
		Code: "sa",
		Name: "Sanskrit",
	},
	{
		Code: "sa-IN",
		Name: "Sanskrit (India)",
	},
	{
		Code: "sc",
		Name: "Sardinian",
	},
	{
		Code: "sd",
		Name: "Sindhi",
	},
	{
		Code: "se",
		Name: "Sami (Northern)",
	},
	{
		Code: "se-FI",
		Name: "Sami (Northern) (Finland)",
	},
	{
		Code: "se-NO",
		Name: "Sami (Northern) (Norway)",
	},
	{
		Code: "se-SE",
		Name: "Sami (Northern) (Sweden)",
	},
	{
		Code: "sg",
		Name: "Sango",
	},
	{
		Code: "si",
		Name: "Sinhala",
	},
	{
		Code: "sk",
		Name: "Slovak",
	},
	{
		Code: "sk-SK",
		Name: "Slovak (Slovakia)",
	},
	{
		Code: "sl",
		Name: "Slovenian",
	},
	{
		Code: "sl-SI",
		Name: "Slovenian (Slovenia)",
	},
	{
		Code: "sm",
		Name: "Samoan",
	},
	{
		Code: "sn",
		Name: "Shona",
	},
	{
		Code: "so",
		Name: "Somali",
	},
	{
		Code: "sq",
		Name: "Albanian",
	},
	{
		Code: "sq-AL",
		Name: "Albanian (Albania)",
	},
	{
		Code: "sr",
		Name: "Serbian",
	},
	{
		Code: "sr-BA",
		Name: "Serbian (Latin) (Bosnia and Herzegovina)",
	},
	{
		Code: "sr-SP",
		Name: "Serbian (Latin) (Serbia and Montenegro)",
	},
	{
		Code: "ss",
		Name: "Swati",
	},
	{
		Code: "st",
		Name: "Southern Sotho",
	},
	{
		Code: "su",
		Name: "Sundanese",
	},
	{
		Code: "sv",
		Name: "Swedish",
	},
	{
		Code: "sv-FI",
		Name: "Swedish (Finland)",
	},
	{
		Code: "sv-SE",
		Name: "Swedish (Sweden)",
	},
	{
		Code: "sw",
		Name: "Swahili",
	},
	{
		Code: "sw-KE",
		Name: "Swahili (Kenya)",
	},
	{
		Code: "syr",
		Name: "Syriac",
	},
	{
		Code: "syr-SY",
		Name: "Syriac (Syria)",
	},
	{
		Code: "ta",
		Name: "Tamil",
	},
	{
		Code: "ta-IN",
		Name: "Tamil (India)",
	},
	{
		Code: "te",
		Name: "Telugu",
	},
	{
		Code: "te-IN",
		Name: "Telugu (India)",
	},
	{
		Code: "tg",
		Name: "Tajik",
	},
	{
		Code: "th",
		Name: "Thai",
	},
	{
		Code: "th-TH",
		Name: "Thai (Thailand)",
	},
	{
		Code: "ti",
		Name: "Tigrinya",
	},
	{
		Code: "tk",
		Name: "Turkmen",
	},
	{
		Code: "tl",
		Name: "Tagalog",
	},
	{
		Code: "tl-PH",
		Name: "Tagalog (Philippines)",
	},
	{
		Code: "tn",
		Name: "Tswana",
	},
	{
		Code: "tn-ZA",
		Name: "Tswana (South Africa)",
	},
	{
		Code: "to",
		Name: "Tonga (Tonga Islands)",
	},
	{
		Code: "tr",
		Name: "Turkish",
	},
	{
		Code: "tr-TR",
		Name: "Turkish (Turkey)",
	},
	{
		Code: "ts",
		Name: "Tsonga",
	},
	{
		Code: "tt",
		Name: "Tatar",
	},
	{
		Code: "tt-RU",
		Name: "Tatar (Russia)",
	},
	{
		Code: "tw",
		Name: "Twi",
	},
	{
		Code: "ty",
		Name: "Tahitian",
	},
	{
		Code: "ug",
		Name: "Uighur",
	},
	{
		Code: "uk",
		Name: "Ukrainian",
	},
	{
		Code: "uk-UA",
		Name: "Ukrainian (Ukraine)",
	},
	{
		Code: "ur",
		Name: "Urdu",
	},
	{
		Code: "ur-PK",
		Name: "Urdu (Islamic Republic of Pakistan)",
	},
	{
		Code: "uz",
		Name: "Uzbek (Latin)",
	},
	{
		Code: "uz-UZ",
		Name: "Uzbek (Latin) (Uzbekistan)",
	},
	{
		Code: "ve",
		Name: "Venda",
	},
	{
		Code: "vi",
		Name: "Vietnamese",
	},
	{
		Code: "vi-VN",
		Name: "Vietnamese (Viet Nam)",
	},
	{
		Code: "vo",
		Name: "Volapük",
	},
	{
		Code: "wa",
		Name: "Walloon",
	},
	{
		Code: "wo",
		Name: "Wolof",
	},
	{
		Code: "xh",
		Name: "Xhosa",
	},
	{
		Code: "xh-ZA",
		Name: "Xhosa (South Africa)",
	},
	{
		Code: "yi",
		Name: "Yiddish",
	},
	{
		Code: "yo",
		Name: "Yoruba",
	},
	{
		Code: "yue-CN",
		Name: "Chinese (Yue)",
	},
	{
		Code: "za",
		Name: "Zhuang",
	},
	{
		Code: "zh",
		Name: "Chinese",
	},
	{
		Code: "zh-CN",
		Name: "Chinese (Simplified)",
	},
	{
		Code: "zh-HK",
		Name: "Chinese (Hong Kong)",
	},
	{
		Code: "zh-MO",
		Name: "Chinese (Macau)",
	},
	{
		Code: "zh-SG",
		Name: "Chinese (Singapore)",
	},
	{
		Code: "zh-TW",
		Name: "Chinese (Taiwan)",
	},
	{
		Code: "zu",
		Name: "Zulu",
	},
	{
		Code: "zu-ZA",
		Name: "Zulu (South Africa)",
	},
}

// FullIsoList returns a slice of all supported ISO language and locale codes.
// This includes both standard language codes (e.g., "en") and region-specific codes (e.g., "en-US").
// The returned slice contains all language codes that are valid for use in PingOne,
// including both reserved and customizable language codes.
func FullIsoList() []string {
	returnVar := make([]string, len(isoList))
	for i, c := range isoList {
		returnVar[i] = c.Code
	}

	return returnVar
}

// FullIsoListString returns a formatted string containing all supported ISO language and locale codes.
// The codes are sorted alphabetically and formatted with backticks for documentation purposes.
// This function is useful for generating documentation or validation error messages
// that need to display the complete list of supported language codes.
func FullIsoListString() string {

	slices.Sort(FullIsoList())

	v := make([]string, len(FullIsoList()))
	for i, c := range FullIsoList() {
		v[i] = fmt.Sprintf("`%s`", c)
	}
	return strings.Join(v, ", ")

}

// IsoList returns a slice of ISO language and locale codes that are available for customization.
// This excludes reserved language codes that are managed by PingOne and cannot be customized.
// The returned slice contains only language codes that can be used for custom language configurations.
func IsoList() []string {

	v := make([]string, 0)
	for _, c := range FullIsoList() {
		if !slices.Contains(reservedLanguageCodes, c) {
			v = append(v, c)
		}
	}

	return v
}

// ReservedIsoList returns a slice of language codes that are reserved by PingOne.
// These language codes have predefined translations and configurations that cannot be customized.
// This function is useful for validation to ensure users don't attempt to override reserved languages.
func ReservedIsoList() []string {
	return reservedLanguageCodes
}

// IsoReservedListString returns a formatted string containing all reserved language codes.
// The codes are sorted alphabetically and formatted with backticks for documentation purposes.
// This function is useful for generating documentation or validation error messages
// that need to display the list of reserved language codes that cannot be customized.
func IsoReservedListString() string {

	slices.Sort(reservedLanguageCodes)

	v := make([]string, len(reservedLanguageCodes))
	for i, c := range reservedLanguageCodes {
		v[i] = fmt.Sprintf("`%s`", c)
	}
	return strings.Join(v, ", ")

}
