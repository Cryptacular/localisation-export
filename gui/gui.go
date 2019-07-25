package gui

import (
	"strings"

	"github.com/tadvi/winc"
)

const (
	langResx  = iota
	langXliff = iota
	langXML   = iota
)

var availableLanguages = []string{"de", "el", "es", "fr", "it", "ja", "ko", "nl", "pt-br", "ro", "sv", "th", "zh"}

// Create shows the Windows GUI to export localised files
func Create(onExecute func(int, string, []string) error) {
	mainWindow := winc.NewForm(nil)
	mainWindow.SetSize(300, 200+(len(availableLanguages)/2+1)*30)
	mainWindow.SetText("RESX Exporter")

	label := winc.NewLabel(mainWindow)
	label.SetPos(10, 20)
	label.SetText("Folder Path:")

	pathTextBox := winc.NewEdit(mainWindow)
	pathTextBox.SetPos(10, 40)

	fileType := langResx
	fileTypeDropdown := winc.NewComboBox(mainWindow)
	fileTypeDropdown.SetPos(220, 40)
	fileTypeDropdown.SetSize(50, 20)
	fileTypeDropdown.InsertItem(langResx, "ISW")
	fileTypeDropdown.InsertItem(langXliff, "iOS")
	fileTypeDropdown.InsertItem(langXML, "Android")
	fileTypeDropdown.SetSelectedItem(langResx)
	fileTypeDropdown.OnClose().Bind(func(e *winc.Event) {
		fileType = fileTypeDropdown.SelectedItem()
	})

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

	responseLabel := winc.NewLabel(mainWindow)
	responseLabel.SetSize(400, 30)
	responseLabel.SetPos(10, 120+(len(availableLanguages)/2+1)*30)
	responseLabel.Hide()

	createButton.OnClick().Bind(func(e *winc.Event) {
		responseLabel.SetText("")
		responseLabel.Hide()
		languagesToInclude := []string{}
		for _, c := range checkboxes {
			if c.checkbox.Checked() {
				languagesToInclude = append(languagesToInclude, c.languageCode)
			}
		}
		err := onExecute(fileType, pathTextBox.Text(), languagesToInclude)
		if err != nil {
			responseLabel.SetText("Oops! Please enter a valid folder path :)")
			responseLabel.Show()
		} else {
			responseLabel.SetText("Success!")
			responseLabel.Show()
		}
	})

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(wndOnClose)

	winc.RunMainLoop()
}

type langCheckbox struct {
	languageCode string
	checkbox     *winc.CheckBox
}

func wndOnClose(arg *winc.Event) {
	winc.Exit()
}
