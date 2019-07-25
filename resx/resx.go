package resx

import (
	"encoding/xml"
	"io/ioutil"
	"strings"

	"github.com/Cryptacular/resx-exporter/localisation"
)

// Reader implements the `LocalisationReader` interface to read RESX files
type Reader struct{}

func (r Reader) Read(path string, languagesToInclude []string) (localisation.Languages, error) {
	return buildLanguagesFromResx(path, languagesToInclude)
}

// DetectLanguages looks at a folder path and checks what languages are available
func (r Reader) DetectLanguages(path string) []string {
	langs := []string{}
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return langs
	}

	for _, f := range files {
		name := f.Name()
		isDirectory := f.IsDir()

		if !isDirectory && strings.HasSuffix(name, ".resx") {
			parts := strings.Split(name, ".")
			length := len(parts)

			if length <= 2 {
				continue
			}

			language := parts[length-2]
			langs = append(langs, language)
		}
	}

	return langs
}

func buildLanguagesFromResx(path string, languagesToInclude []string) (localisation.Languages, error) {
	baseLang, err := readResx(path+"/UIStrings.resx", "en")

	if err != nil {
		return localisation.Languages{}, err
	}

	otherLangs := []localisation.Language{}

	for _, p := range languagesToInclude {
		l, err := readResx(path+"/UIStrings."+p+".resx", p)
		if err != nil {
			return localisation.Languages{}, err
		}
		otherLangs = append(otherLangs, l)
	}

	return localisation.Languages{
		Base:         baseLang,
		Translations: otherLangs,
	}, nil
}

func readResx(filename, lang string) (localisation.Language, error) {
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		return localisation.Language{}, err
	}

	r := convertToResx(file)
	return localisation.ConvertToLanguage(r, lang), nil
}

func convertToResx(content []byte) localisation.LanguageEntries {
	le := localisation.LanguageEntries{}
	r := resx{}
	xml.Unmarshal([]byte(content), &r)

	for _, e := range r.Entries {
		le = append(le, localisation.LanguageEntry{
			Key:     e.Key,
			Value:   e.Value,
			Comment: e.Comment,
		})
	}

	return le
}

type resx struct {
	Entries []resxEntry `xml:"data"`
}

type resxEntry struct {
	Key     string `xml:"name,attr"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}
