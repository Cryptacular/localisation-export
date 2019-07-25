package localisation

import (
	"strings"
	"testing"
)

func TestConvertToSpreadsheetHasHeaderRow(t *testing.T) {
	langs := Languages{
		Base: Language{
			languageCode: "EN",
		},
		Translations: []Language{
			Language{
				languageCode: "NL",
			},
			Language{
				languageCode: "DE",
			},
		},
	}

	actual := ConvertToSpreadsheet(langs)

	key := actual[0][0]
	en := actual[0][1]
	nl := actual[0][2]
	de := actual[0][3]

	if key != "Key" || en != "EN" || nl != "NL" || de != "DE" {
		t.Errorf(`Actual: %s"; Expected: "Key EN NL DE"`, strings.Join([]string{key, en, nl, de}, " "))
	}
}

func TestConvertToSpreadsheetReturnsRowsInCorrectOrder(t *testing.T) {
	r := LanguageEntries{
		LanguageEntry{
			Key: "One",
		},
		LanguageEntry{
			Key: "Two",
		},
		LanguageEntry{
			Key: "Three",
		},
	}

	l := ConvertToLanguage(r, "en")

	langs := Languages{
		Base: Language{
			dict: l.dict,
			keys: l.keys,
		},
	}

	actual := ConvertToSpreadsheet(langs)

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
	german := LanguageEntries{
		LanguageEntry{
			Key:   "One",
			Value: "Ein",
		},
	}
	dutch := LanguageEntries{
		LanguageEntry{
			Key:   "One",
			Value: "Een",
		},
	}
	lGerman := ConvertToLanguage(german, "de")
	lDutch := ConvertToLanguage(dutch, "nl")

	langs := Languages{
		Base: lGerman,
		Translations: []Language{
			lDutch,
		},
	}

	actual := ConvertToSpreadsheet(langs)
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
	english := LanguageEntries{
		LanguageEntry{
			Key:     "One",
			Value:   "One",
			Comment: "Hi there",
		},
	}
	dutch := LanguageEntries{
		LanguageEntry{
			Key:     "One",
			Value:   "Een",
			Comment: "Not this one",
		},
	}
	lEnglish := ConvertToLanguage(english, "en")
	lDutch := ConvertToLanguage(dutch, "nl")

	langs := Languages{
		Base: lEnglish,
		Translations: []Language{
			lDutch,
		},
	}

	actual := ConvertToSpreadsheet(langs)
	row := actual[1]

	if len(row) != 4 {
		t.Errorf(`Wrong number of cells in row. Actual: %d; Expected: 4`, len(row))
	}

	if comment := row[3]; comment != "Hi there" {
		t.Errorf(`Comment is wrong. Actual: %s"; Expected: "Hi there"`, comment)
	}
}

func TestConvertToSpreadsheetWorksWhenLanguageIsMissingKeys(t *testing.T) {
	english := LanguageEntries{
		LanguageEntry{
			Key:   "One",
			Value: "One",
		},
		LanguageEntry{
			Key:   "Two",
			Value: "Two",
		},
		LanguageEntry{
			Key:   "Three",
			Value: "Three",
		},
	}
	dutch := LanguageEntries{
		LanguageEntry{
			Key:   "One",
			Value: "Een",
		},
		LanguageEntry{
			Key:   "Three",
			Value: "Drie",
		},
	}
	lEnglish := ConvertToLanguage(english, "en")
	lDutch := ConvertToLanguage(dutch, "nl")

	langs := Languages{
		Base: Language{
			dict: lEnglish.dict,
			keys: lEnglish.keys,
		},
		Translations: []Language{
			lDutch,
		},
	}

	actual := ConvertToSpreadsheet(langs)
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

func TestConvertToLanguageDictSetsCorrectLanguage(t *testing.T) {
	r := LanguageEntries{}

	actual := ConvertToLanguage(r, "en")

	if actual.languageCode != "EN" {
		t.Errorf(`Actual: "%s"; Expected: "EN"`, actual.languageCode)
	}
}

func TestConvertToLanguageDictMatchesCorrectKeyToValue(t *testing.T) {
	r := LanguageEntries{
		LanguageEntry{
			Key:   "A nose on...",
			Value: "your face",
		},
		LanguageEntry{
			Key:   "A toe on...",
			Value: "your foot",
		},
	}

	actual := ConvertToLanguage(r, "en")

	if value := actual.dict["A nose on..."].Value; value != "your face" {
		t.Errorf(`Actual: %s; Expected: "A nose on": "your face"`, `"A nose on...": "`+value+`"`)
	}
}
