package main

import (
	"errors"
	"strings"

	"github.com/Cryptacular/resx-exporter/excel"
	"github.com/Cryptacular/resx-exporter/gui"
	"github.com/Cryptacular/resx-exporter/localisation"
	"github.com/Cryptacular/resx-exporter/resx"
	"github.com/Cryptacular/resx-exporter/xliff"
)

const (
	langResx  = iota
	langXliff = iota
	langXML   = iota
)

func main() {
	gui.Create(execute)
}

func execute(fileType int, path string, languagesToInclude []string) error {
	builder, err := languageBuilder(fileType)

	if err != nil {
		return err
	}

	langs, err := builder.Read(path, languagesToInclude)

	if err != nil {
		return err
	}

	out := localisation.ConvertToSpreadsheet(langs)

	filename := buildFilename(path, languagesToInclude)
	excel.Write(out, filename)

	return nil
}

func languageBuilder(fileType int) (localisation.Reader, error) {
	if fileType == langResx {
		return resx.ResxReader{}, nil
	} else if fileType == langXliff {
		return xliff.XliffReader{}, nil
	}

	return nil, errors.New("No valid file types found")
}

func buildFilename(path string, Languages []string) string {
	filename := path + "/EN"
	for _, l := range Languages {
		filename += "-" + strings.ToUpper(l)
	}
	filename += ".xlsx"
	return filename
}
