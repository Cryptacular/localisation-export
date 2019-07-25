package main

import (
	"errors"
	"io/ioutil"
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
	props := guiProps{
		readers: map[int]localisation.Reader{},
	}
	gui.Create(props)
}

type guiProps struct {
	readers map[int]localisation.Reader
}

func (g guiProps) DetectFormat(path string) (localisation.Format, error) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return localisation.Format{}, err
	}

	for _, f := range files {
		name := f.Name()
		isDirectory := f.IsDir()

		if !isDirectory && strings.HasSuffix(name, ".resx") {
			return g.detectLanguages(langResx, path)
		}

		if isDirectory && strings.Contains(name, "xcloc") {
			return g.detectLanguages(langXliff, path)
		}

		if isDirectory && strings.Contains(name, "values-") {
			return g.detectLanguages(langXML, path)
		}
	}

	return localisation.Format{}, errors.New("Could not detect file type")
}

func (g guiProps) GetReader(fileType int) (localisation.Reader, error) {
	lr := g.readers[fileType]

	if lr != nil {
		return lr, nil
	}

	lr2, err := localisationReaderFactory(fileType)

	if err != nil {
		return nil, err
	}

	g.readers[fileType] = lr2

	return lr2, err
}

func (g guiProps) OnExecute(fileType int, path string, languagesToInclude []string) error {
	localisationReader, err := g.GetReader(fileType)

	if err != nil {
		return errors.New("No valid reader found")
	}

	langs, err := localisationReader.Read(path, languagesToInclude)

	if err != nil {
		return err
	}

	out := localisation.ConvertToSpreadsheet(langs)

	filename := buildFilename(path, languagesToInclude)
	excel.Write(out, filename)

	return nil
}

func (g guiProps) detectLanguages(fileType int, path string) (localisation.Format, error) {
	lr, err := g.GetReader(fileType)
	if err != nil {
		return localisation.Format{}, err
	}

	langs := lr.DetectLanguages(path)

	return localisation.Format{
		FileType:           fileType,
		AvailableLanguages: langs,
	}, nil
}

func localisationReaderFactory(fileType int) (localisation.Reader, error) {
	if fileType == langResx {
		return resx.Reader{}, nil
	} else if fileType == langXliff {
		return xliff.Reader{}, nil
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
