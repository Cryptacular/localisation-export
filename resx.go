package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func readResx(filename, lang string) language {
	file := readFile(filename)
	r := convertToResx(file)
	return convertToLanguageDict(r, lang)
}

func readFile(filename string) []byte {
	f, err := ioutil.ReadFile(filename)

	if err != nil {
		panic("File not found: " + filename)
	}

	return f
}

func convertToResx(content []byte) resx {
	r := resx{}
	xml.Unmarshal([]byte(content), &r)
	return r
}

func convertToLanguageDict(r resx, lang string) language {
	ld := languageDict{
		language: lang,
		entries:  map[string]entry{},
	}

	keys := []string{}

	for i := range r.Entries {
		e := r.Entries[i]
		ld.entries[e.Key] = entry{e.Value, e.Comment}
		keys = append(keys, e.Key)
	}

	return language{
		dict: ld,
		keys: keys,
	}
}

func convertToSpreadsheet(langs languages) [][]string {
	headerRow := []string{"Key"}

	for _, l := range langs.dicts {
		headerRow = append(headerRow, strings.ToUpper(l.language))
	}

	headerRow = append(headerRow, "Comments")
	rows := [][]string{headerRow}

	for _, key := range langs.keys {
		row := []string{key}

		for _, l := range langs.dicts {
			row = append(row, l.entries[key].Value)
		}

		row = append(row, langs.dicts[0].entries[key].Comment)

		rows = append(rows, row)
	}

	return rows
}

func writeExcel(content [][]string, path string) {
	f := excelize.NewFile()

	for i, row := range content {
		cell := "A" + strconv.Itoa(i+1)
		f.SetSheetRow("Sheet1", cell, &row)
	}

	err := f.SaveAs(path)
	if err != nil {
		fmt.Println(err)
	}
}

type resx struct {
	Entries []resxEntry `xml:"data"`
}

type resxEntry struct {
	Key     string `xml:"name,attr"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

type language struct {
	dict languageDict
	keys []string
}

type languages struct {
	dicts []languageDict
	keys  []string
}

type languageDict struct {
	language string
	entries  map[string]entry
}

type entry struct {
	Value   string
	Comment string
}
