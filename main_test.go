package main

import "testing"

func TestBuildFilenameForZeroLanguages(t *testing.T) {
	languages := []string{}

	actual := buildFilename(".", languages)

	if actual != "./EN.xlsx" {
		t.Errorf(`Actual: "%s"; Expected: "./EN.xlsx"`, actual)
	}
}

func TestBuildFilenameForOneLanguage(t *testing.T) {
	languages := []string{"nl"}

	actual := buildFilename(".", languages)

	if actual != "./EN-NL.xlsx" {
		t.Errorf(`Actual: "%s"; Expected: "./EN-NL.xlsx"`, actual)
	}
}

func TestBuildFilenameForMultipleLanguages(t *testing.T) {
	languages := []string{"nl", "de", "pt-br"}

	actual := buildFilename(".", languages)

	if actual != "./EN-NL-DE-PT-BR.xlsx" {
		t.Errorf(`Actual: "%s"; Expected: "./EN-NL-DE-PT-BR.xlsx"`, actual)
	}
}
