package main

import (
	"strings"
	"testing"
)

func TestConvertToResxReadsCorrectNumberOfEntries(t *testing.T) {
	r := `<root>
			<data name="One"></data>
			<data name="Two"></data>
			<data name="Three"></data>
		  </root>`

	actual := convertToResx([]byte(r))

	if length := len(actual.Entries); length != 3 {
		t.Errorf(`Actual: %d; Expected: 3`, length)
	}
}

func TestConvertToResxReadsEntriesInCorrectOrder(t *testing.T) {
	r := `<root>
			<data name="One"></data>
			<data name="Two"></data>
			<data name="Three"></data>
		  </root>`

	actual := convertToResx([]byte(r))
	one := actual.Entries[0]
	two := actual.Entries[1]
	three := actual.Entries[2]

	if one.Key != "One" ||
		two.Key != "Two" ||
		three.Key != "Three" {
		t.Errorf(`Actual: "%s"; Expected: "One Two Three"`, strings.Join([]string{one.Key, two.Key, three.Key}, " "))
	}
}

func TestConvertToResxReadsKeyCorrectly(t *testing.T) {
	r := `<root>
			<data name="ThisIsTheKey">
				<value>This is the value of the string</value>
				<comment>Now this would be the comment</comment>
			</data>
		  </root>`

	actual := convertToResx([]byte(r))

	if length := len(actual.Entries); length != 1 {
		t.Errorf(`Actual: %d; Expected: 1`, length)
	}

	if key := actual.Entries[0].Key; key != "ThisIsTheKey" {
		t.Errorf(`Actual: "%s"; Expected: "ThisIsTheKey"`, key)
	}
}

func TestConvertToResxReadsValueCorrectly(t *testing.T) {
	r := `<root>
			<data name="ThisIsTheKey">
				<value>This is the value of the string</value>
				<comment>Now this would be the comment</comment>
			</data>
		  </root>`

	actual := convertToResx([]byte(r))

	if value := actual.Entries[0].Value; value != "This is the value of the string" {
		t.Errorf(`Actual: "%s"; Expected: "This is the value of the string"`, value)
	}
}

func TestConvertToResxReadsCommentCorrectly(t *testing.T) {
	r := `<root>
			<data name="ThisIsTheKey">
				<value>This is the value of the string</value>
				<comment>Now this would be the comment</comment>
			</data>
		  </root>`

	actual := convertToResx([]byte(r))

	if comment := actual.Entries[0].Comment; comment != "Now this would be the comment" {
		t.Errorf(`Actual: "%s"; Expected: "Now this would be the comment"`, comment)
	}
}

func TestConvertToResxReadsEmptyValueCorrectly(t *testing.T) {
	r := `<root>
			<data name="EmptyOne">
				<value></value>
			</data>
			<data name="EmptyTwo">
				<value/>
			</data>
		  </root>`

	actual := convertToResx([]byte(r))

	valueOne := actual.Entries[0].Value
	valueTwo := actual.Entries[1].Value

	if valueOne != "" {
		t.Errorf(`Actual: "%s"; Expected: ""`, valueOne)
	} else if valueTwo != "" {
		t.Errorf(`Actual: "%s"; Expected: ""`, valueTwo)
	}
}

func TestConvertToLanguageDictSetsCorrectLanguage(t *testing.T) {
	r := resx{}

	actual := convertToLanguage(r, "en")

	if actual.languageCode != "EN" {
		t.Errorf(`Actual: "%s"; Expected: "EN"`, actual.languageCode)
	}
}

func TestConvertToLanguageDictMatchesCorrectKeyToValue(t *testing.T) {
	r := resx{
		Entries: []resxEntry{
			resxEntry{
				Key:   "A nose on...",
				Value: "your face",
			},
			resxEntry{
				Key:   "A toe on...",
				Value: "your foot",
			},
		},
	}

	actual := convertToLanguage(r, "en")

	if value := actual.dict["A nose on..."].Value; value != "your face" {
		t.Errorf(`Actual: %s; Expected: "A nose on": "your face"`, `"A nose on...": "`+value+`"`)
	}
}

func TestConvertToSpreadsheetHasHeaderRow(t *testing.T) {
	langs := languages{
		base: language{
			languageCode: "EN",
		},
		translations: []language{
			language{
				languageCode: "NL",
			},
			language{
				languageCode: "DE",
			},
		},
	}

	actual := convertToSpreadsheet(langs)

	key := actual[0][0]
	en := actual[0][1]
	nl := actual[0][2]
	de := actual[0][3]

	if key != "Key" || en != "EN" || nl != "NL" || de != "DE" {
		t.Errorf(`Actual: %s"; Expected: "Key EN NL DE"`, strings.Join([]string{key, en, nl, de}, " "))
	}
}

func TestConvertToSpreadsheetReturnsRowsInCorrectOrder(t *testing.T) {
	r := resx{
		Entries: []resxEntry{
			resxEntry{
				Key: "One",
			},
			resxEntry{
				Key: "Two",
			},
			resxEntry{
				Key: "Three",
			},
		},
	}

	l := convertToLanguage(r, "en")

	langs := languages{
		base: language{
			dict: l.dict,
			keys: l.keys,
		},
	}

	actual := convertToSpreadsheet(langs)

	firstKey := actual[1][0]
	secondKey := actual[2][0]
	thirdKey := actual[3][0]

	if firstKey != "One" ||
		secondKey != "Two" ||
		thirdKey != "Three" {
		t.Errorf(`Actual: %s"; Expected: "One Two Three"`, strings.Join([]string{firstKey, secondKey, thirdKey}, " "))
	}
}

func TestConvertToSpreadsheetMatchesKeysToLanguagesCorrectly(t *testing.T) {
	german := resx{
		Entries: []resxEntry{
			resxEntry{
				Key:   "One",
				Value: "Ein",
			},
		},
	}
	dutch := resx{
		Entries: []resxEntry{
			resxEntry{
				Key:   "One",
				Value: "Een",
			},
		},
	}
	lGerman := convertToLanguage(german, "de")
	lDutch := convertToLanguage(dutch, "nl")

	langs := languages{
		base: lGerman,
		translations: []language{
			lDutch,
		},
	}

	actual := convertToSpreadsheet(langs)
	row := actual[1]

	if key := row[0]; key != "One" {
		t.Errorf(`Key is wrong. Actual: %s"; Expected: "One"`, key)
	}

	if de := row[1]; de != "Ein" {
		t.Errorf(`German is wrong. Actual: %s"; Expected: "Ein"`, de)
	}

	if nl := row[2]; nl != "Een" {
		t.Errorf(`Dutch is wrong. Actual: %s"; Expected: "Een"`, nl)
	}
}

func TestConvertToSpreadsheetIncludesCommentFromFirstLanguage(t *testing.T) {
	english := resx{
		Entries: []resxEntry{
			resxEntry{
				Key:     "One",
				Value:   "One",
				Comment: "Hi there",
			},
		},
	}
	dutch := resx{
		Entries: []resxEntry{
			resxEntry{
				Key:     "One",
				Value:   "Een",
				Comment: "Not this one",
			},
		},
	}
	lEnglish := convertToLanguage(english, "en")
	lDutch := convertToLanguage(dutch, "nl")

	langs := languages{
		base: lEnglish,
		translations: []language{
			lDutch,
		},
	}

	actual := convertToSpreadsheet(langs)
	row := actual[1]

	if len(row) != 4 {
		t.Errorf(`Wrong number of cells in row. Actual: %d; Expected: 4`, len(row))
	}

	if comment := row[3]; comment != "Hi there" {
		t.Errorf(`Comment is wrong. Actual: %s"; Expected: "Hi there"`, comment)
	}
}

func TestConvertToSpreadsheetWorksWhenLanguageIsMissingKeys(t *testing.T) {
	english := resx{
		Entries: []resxEntry{
			resxEntry{
				Key:   "One",
				Value: "One",
			},
			resxEntry{
				Key:   "Two",
				Value: "Two",
			},
			resxEntry{
				Key:   "Three",
				Value: "Three",
			},
		},
	}
	dutch := resx{
		Entries: []resxEntry{
			resxEntry{
				Key:   "One",
				Value: "Een",
			},
			resxEntry{
				Key:   "Three",
				Value: "Drie",
			},
		},
	}
	lEnglish := convertToLanguage(english, "en")
	lDutch := convertToLanguage(dutch, "nl")

	langs := languages{
		base: language{
			dict: lEnglish.dict,
			keys: lEnglish.keys,
		},
		translations: []language{
			lDutch,
		},
	}

	actual := convertToSpreadsheet(langs)
	firstRow := actual[1]
	secondRow := actual[2]
	thirdRow := actual[3]

	if en := firstRow[1]; en != "One" {
		t.Errorf(`English is wrong. Actual: %s"; Expected: "One"`, en)
	}

	if nl := firstRow[2]; nl != "Een" {
		t.Errorf(`Dutch is wrong. Actual: %s"; Expected: "Een"`, nl)
	}

	if en := secondRow[1]; en != "Two" {
		t.Errorf(`English is wrong. Actual: %s"; Expected: "Two"`, en)
	}

	if nl := secondRow[2]; nl != "" {
		t.Errorf(`Dutch is wrong. Actual: %s"; Expected: ""`, nl)
	}

	if en := thirdRow[1]; en != "Three" {
		t.Errorf(`English is wrong. Actual: %s"; Expected: "Three"`, en)
	}

	if nl := thirdRow[2]; nl != "Drie" {
		t.Errorf(`Dutch is wrong. Actual: %s"; Expected: "Drie"`, nl)
	}
}
