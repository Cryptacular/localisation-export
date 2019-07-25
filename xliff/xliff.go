package xliff

import (
	"encoding/xml"
	"io/ioutil"
	"strings"

	"github.com/Cryptacular/resx-exporter/localisation"
)

// Reader implements the `LocalisationReader` interface to read XLIFF files
type Reader struct{}

func (r Reader) Read(path string, languagesToInclude []string) (localisation.Languages, error) {
	return buildLanguagesFromXliff(path, languagesToInclude)
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

		if isDirectory && strings.HasSuffix(name, ".xcloc") {
			parts := strings.Split(name, ".")

			if len(parts) > 2 {
				continue
			}

			language := parts[0]
			if language == "en" {
				continue
			}
			langs = append(langs, language)
		}
	}

	return langs
}

func buildLanguagesFromXliff(path string, languagesToInclude []string) (localisation.Languages, error) {
	baseLang, err := readXliff(path+"/en.xcloc/Localized Contents/en.xliff", "en")

	if err != nil {
		return localisation.Languages{}, err
	}

	otherLangs := []localisation.Language{}

	for _, p := range languagesToInclude {
		l, err := readXliff(path+"/"+p+".xcloc/Localized Contents/"+p+".xliff", p)
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

func readXliff(filename, lang string) (localisation.Language, error) {
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		return localisation.Language{}, err
	}

	x := convertToXliff(file)
	return localisation.ConvertToLanguage(x, lang), nil
}

func convertToXliff(content []byte) localisation.LanguageEntries {
	x := localisation.LanguageEntries{}
	xd := xliffDoc{}
	xml.Unmarshal([]byte(content), &xd)

	for _, f := range xd.File {
		for _, tu := range f.Body.TransUnit {
			if tu.Target == nil {
				continue
			}
			x = append(x, localisation.LanguageEntry{
				Key:     tu.ID,
				Value:   tu.Target.Inner,
				Comment: tu.Note,
			})
		}
	}

	return x
}

type xliffSource struct {
	XMLName xml.Name
	Inner   string `xml:",chardata"`
	Lang    string `xml:"lang,attr"`
	Space   string `xml:"space,attr,omitempty"`
	State   string `xml:"state,attr,omitempty"`
}

type xliffTarget xliffSource

type xliffTransUnit struct {
	ID     string       `xml:"id,attr"`
	Source xliffSource  `xml:"source"`
	Target *xliffTarget `xml:"target,omitempty"`
	Note   string       `xml:"note,omitempty"`
}

type xliffBody struct {
	XMLName   xml.Name         `xml:"body"`
	TransUnit []xliffTransUnit `xml:"trans-unit"`
}

type xliffFile struct {
	Original   string    `xml:"original,attr"`
	SourceLang string    `xml:"source-localisation.Language,attr,omitempty"`
	TargetLang string    `xml:"target-localisation.Language,attr,omitempty"`
	DataType   string    `xml:"datatype,attr,omitempty"`
	Body       xliffBody `xml:"body"`
}

type xliffDoc struct {
	XMLName xml.Name    `xml:"xliff"`
	Version string      `xml:"version,attr"`
	Xmlns   string      `xml:"xmlns,attr"`
	File    []xliffFile `xml:"file"`
}
