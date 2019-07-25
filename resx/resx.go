package resx

import (
	"encoding/xml"
	"io/ioutil"

	"github.com/Cryptacular/resx-exporter/localisation"
)

// ResxReader implements the `LocalisationReader` interface to read RESX files
type ResxReader struct{}

func (r ResxReader) Read(path string, languagesToInclude []string) (localisation.Languages, error) {
	return buildLanguagesFromResx(path, languagesToInclude)
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
