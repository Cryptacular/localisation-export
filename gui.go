package main

import (
	"strings"

	"github.com/tadvi/winc"
)

func createGui(onExecute func(string, []string) error) {
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
		err := onExecute(pathTextBox.Text(), languagesToInclude)
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
