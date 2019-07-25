package localisation

import "strings"

// Reader reads files of a particular file type and converts it to the `languages` type
type Reader interface {
	Read(path string, languagesToInclude []string) (Languages, error)
	DetectLanguages(path string) []string
}

// ConvertToLanguage converts a slice of keys, values and comments into a `Language` object
func ConvertToLanguage(x LanguageEntries, lang string) Language {
	ld := map[string]entry{}

	keys := []string{}

	for i := range x {
		e := x[i]
		ld[e.Key] = entry{e.Value, e.Comment}
		keys = append(keys, e.Key)
	}

	return Language{
		dict:         ld,
		keys:         keys,
		languageCode: strings.ToUpper(lang),
	}
}

// ConvertToSpreadsheet converts a `Languages` object to rows and columns
func ConvertToSpreadsheet(langs Languages) [][]string {
	headerRow := []string{"Key", langs.Base.languageCode}

	for _, l := range langs.Translations {
		headerRow = append(headerRow, l.languageCode)
	}

	headerRow = append(headerRow, "Comments")
	rows := [][]string{headerRow}

	for _, key := range langs.Base.keys {
		row := []string{key, langs.Base.dict[key].Value}

		for _, l := range langs.Translations {
			row = append(row, l.dict[key].Value)
		}

		row = append(row, langs.Base.dict[key].Comment)

		rows = append(rows, row)
	}

	return rows
}

// Format specifies to format of the files and what languages are available
type Format struct {
	FileType           int
	AvailableLanguages []string
}

// LanguageEntries is a slice of `LanguageEntry`, which contains Keys, Values and Comments
type LanguageEntries []LanguageEntry

// LanguageEntry represents a single entry, with `Key`, `Value` and `Comment`
type LanguageEntry struct {
	Key     string
	Value   string
	Comment string
}

// Languages contains a Base Language and its translations
type Languages struct {
	Base         Language
	Translations []Language
}

// Language contains a localisation.Language code, a dictionary of localisation keys and translated values, as well as a plain slice of keys
type Language struct {
	dict         languageDict
	languageCode string
	keys         []string
}

type languageDict map[string]entry

type entry struct {
	Value   string
	Comment string
}
