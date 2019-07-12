package main

import (
	"strings"

	"github.com/tadvi/winc"
)

var availableLanguages = []string{"de", "el", "es", "fr", "it", "ja", "ko", "nl", "pt-br", "ro", "sv", "th", "zh"}

func main() {
	mainWindow := winc.NewForm(nil)
	mainWindow.SetSize(300, 200+(len(availableLanguages)/2+1)*30)
	mainWindow.SetText("RESX Exporter")

	label := winc.NewLabel(mainWindow)
	label.SetPos(10, 20)
	label.SetText("Folder Path:")

	pathTextBox := winc.NewEdit(mainWindow)
	pathTextBox.SetPos(10, 40)

	checkboxes := []langCheckbox{}
	for i, l := range availableLanguages {
		c := winc.NewCheckBox(mainWindow)
		c.SetText(strings.ToUpper(l))
		if i%2 == 0 {
			c.SetPos(10, 70+i*15)
		} else {
			c.SetPos(120, 55+i*15)
		}
		checkboxes = append(checkboxes, langCheckbox{
			languageCode: l,
			checkbox:     c,
		})
	}

	createButton := winc.NewPushButton(mainWindow)
	createButton.SetText("Create Excel")
	createButton.SetPos(10, 70+(len(availableLanguages)/2+1)*30)
	createButton.SetSize(200, 40)

	errorLabel := winc.NewLabel(mainWindow)
	errorLabel.SetSize(400, 30)
	errorLabel.SetPos(10, 120+(len(availableLanguages)/2+1)*30)
	errorLabel.Hide()

	createButton.OnClick().Bind(func(e *winc.Event) {
		errorLabel.SetText("")
		errorLabel.Hide()
		languagesToInclude := []string{}
		for _, c := range checkboxes {
			if c.checkbox.Checked() {
				languagesToInclude = append(languagesToInclude, c.languageCode)
			}
		}
		err := execute(pathTextBox.Text(), languagesToInclude)
		if err != nil {
			errorLabel.SetText("Oops! Please enter a valid folder path :)")
			errorLabel.Show()
		}
	})

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	winc.RunMainLoop() // Must call to start event loop.

}

type langCheckbox struct {
	languageCode string
	checkbox     *winc.CheckBox
}

func wndOnClose(arg *winc.Event) {
	winc.Exit()
}

func execute(path string, languagesToInclude []string) error {
	baseLang, err := readResx(path+"/UIStrings.resx", "en")

	if err != nil {
		return err
	}

	otherLangs := []language{}

	for _, p := range languagesToInclude {
		l, err := readResx(path+"/UIStrings."+p+".resx", p)
		if err != nil {
			return err
		}
		otherLangs = append(otherLangs, l)
	}

	langs := buildLanguages(baseLang, otherLangs)

	out := convertToSpreadsheet(langs)

	filename := buildFilename(path, languagesToInclude)
	writeExcel(out, filename)

	return nil
}

func buildLanguages(baseLang language, others []language) languages {
	return languages{
		base:         baseLang,
		translations: others,
	}
}

func buildFilename(path string, languages []string) string {
	filename := path + "/EN"
	for _, l := range languages {
		filename += "-" + strings.ToUpper(l)
	}
	filename += ".xlsx"
	return filename
}
