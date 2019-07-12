package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func readResx(filename, lang string) (language, error) {
	file, err := readFile(filename)

	if err != nil {
		return language{}, err
	}

	r := convertToResx(file)
	return convertToLanguage(r, lang), nil
}

func readFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func convertToResx(content []byte) resx {
	r := resx{}
	xml.Unmarshal([]byte(content), &r)
	return r
}

func convertToLanguage(r resx, lang string) language {
	ld := map[string]entry{}

	keys := []string{}

	for i := range r.Entries {
		e := r.Entries[i]
		ld[e.Key] = entry{e.Value, e.Comment}
		keys = append(keys, e.Key)
	}

	return language{
		dict:         ld,
		keys:         keys,
		languageCode: strings.ToUpper(lang),
	}
}

func convertToSpreadsheet(langs languages) [][]string {
	headerRow := []string{"Key", langs.base.languageCode}

	for _, l := range langs.translations {
		headerRow = append(headerRow, l.languageCode)
	}

	headerRow = append(headerRow, "Comments")
	rows := [][]string{headerRow}

	for _, key := range langs.base.keys {
		row := []string{key, langs.base.dict[key].Value}

		for _, l := range langs.translations {
			row = append(row, l.dict[key].Value)
		}

		row = append(row, langs.base.dict[key].Comment)

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

type languages struct {
	base         language
	translations []language
}

type language struct {
	dict         languageDict
	languageCode string
	keys         []string
}

type languageDict map[string]entry

type entry struct {
	Value   string
	Comment string
}
