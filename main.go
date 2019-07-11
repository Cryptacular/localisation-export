package main

import (
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		panic("Please enter a language code")
	}

	params := os.Args[1:]

	baseLang := readResx("./UIStrings.resx", "en")
	otherLangs := []language{}

	for _, p := range params {
		l := readResx("./UIStrings."+p+".resx", p)
		otherLangs = append(otherLangs, l)
	}

	langs := buildLanguages(baseLang, otherLangs)

	out := convertToSpreadsheet(langs)

	filename := buildFilename(params)
	writeExcel(out, filename)
}

func buildLanguages(baseLang language, others []language) languages {
	return languages{
		base:         baseLang,
		translations: others,
	}
}

func buildFilename(languages []string) string {
	filename := "./EN"
	for _, l := range languages {
		filename += "-" + strings.ToUpper(l)
	}
	filename += ".xlsx"
	return filename
}
