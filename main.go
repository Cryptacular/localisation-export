package main

import (
	"strings"
)

var availableLanguages = []string{"de", "el", "es", "fr", "it", "ja", "ko", "nl", "pt-br", "ro", "sv", "th", "zh"}

func main() {
	createGui(execute)
}

func execute(path string, languagesToInclude []string) error {
	langs, err := buildLanguages(path, languagesToInclude)

	if err != nil {
		return err
	}

	out := convertToSpreadsheet(langs)

	filename := buildFilename(path, languagesToInclude)
	writeExcel(out, filename)

	return nil
}

func buildLanguages(path string, languagesToInclude []string) (languages, error) {
	baseLang, err := readResx(path+"/UIStrings.resx", "en")

	if err != nil {
		return languages{}, err
	}

	otherLangs := []language{}

	for _, p := range languagesToInclude {
		l, err := readResx(path+"/UIStrings."+p+".resx", p)
		if err != nil {
			return languages{}, err
		}
		otherLangs = append(otherLangs, l)
	}

	return languages{
		base:         baseLang,
		translations: otherLangs,
	}, nil
}

func buildFilename(path string, languages []string) string {
	filename := path + "/EN"
	for _, l := range languages {
		filename += "-" + strings.ToUpper(l)
	}
	filename += ".xlsx"
	return filename
}
