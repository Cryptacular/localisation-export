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

	dictionaries := languages{}
	baseLang := readResx("./UIStrings.resx", "en")
	dictionaries.dicts = []languageDict{baseLang.dict}
	dictionaries.keys = baseLang.keys

	for _, l := range params {
		d := readResx("./UIStrings."+l+".resx", l)
		dictionaries.dicts = append(dictionaries.dicts, d.dict)
	}

	out := convertToSpreadsheet(dictionaries)

	filename := buildFilename(params)
	writeExcel(out, filename)
}

func buildFilename(languages []string) string {
	filename := "./EN"
	for _, l := range languages {
		filename += "-" + strings.ToUpper(l)
	}
	filename += ".xlsx"
	return filename
}
