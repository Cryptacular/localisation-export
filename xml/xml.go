package xml

import (
	"encoding/xml"
	"io/ioutil"
	"strings"

	"github.com/Cryptacular/localisation-export/localisation"
)

// Reader implements the `LocalisationReader` interface to read RESX files
type Reader struct{}

func (r Reader) Read(path string, languagesToInclude []string) (localisation.Languages, error) {
	return buildLanguagesFromXML(path, languagesToInclude)
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

		if isDirectory && strings.HasPrefix(name, "values") {
			parts := strings.Split(name, "-")

			if len(parts) < 2 {
				continue
			}

			language := strings.Join(parts[1:], "-")
			langs = append(langs, language)
		}
	}

	return langs
}

func buildLanguagesFromXML(path string, languagesToInclude []string) (localisation.Languages, error) {
	baseLang, err := readXML(path+"/values/strings.xml", "en")

	if err != nil {
		return localisation.Languages{}, err
	}

	otherLangs := []localisation.Language{}

	for _, p := range languagesToInclude {
		l, err := readXML(path+"/values-"+p+"/strings.xml", p)
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

func readXML(filename, lang string) (localisation.Language, error) {
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		return localisation.Language{}, err
	}

	r := convertToValues(file)
	return localisation.ConvertToLanguage(r, lang), nil
}

func convertToValues(content []byte) localisation.LanguageEntries {
	le := localisation.LanguageEntries{}
	r := values{}
	xml.Unmarshal([]byte(content), &r)

	for _, e := range r.Entries {
		le = append(le, localisation.LanguageEntry{
			Key:   e.Key,
			Value: e.Value,
		})
	}

	return le
}

type values struct {
	Entries []value `xml:"string"`
}

type value struct {
	Key   string `xml:"name,attr"`
	Value string `xml:",chardata"`
}
