package gui

import (
	"strings"

	"github.com/Cryptacular/localisation-export/localisation"
	"github.com/tadvi/winc"
)

const (
	langResx  = iota
	langXliff = iota
	langXML   = iota
)

var (
	availableLanguages = []string{}
	langCheckboxes     = []langCheckbox{}
)

// Props defines the properties and callback methods `Create` needs
type Props interface {
	OnExecute(int, string, []string) error
	DetectFormat(string) (localisation.Format, error)
}

// Create shows the Windows GUI to export localised files
func Create(props Props) {
	mainWindow := winc.NewForm(nil)
	mainWindow.SetText("Localisation Export")
	resizeWindow(mainWindow, availableLanguages)

	label := winc.NewLabel(mainWindow)
	label.SetPos(10, 20)
	label.SetText("Folder Path:")

	pathTextBox := winc.NewEdit(mainWindow)
	pathTextBox.SetPos(10, 40)

	fileTypeDropdown := winc.NewComboBox(mainWindow)
	fileTypeDropdown.SetPos(220, 40)
	fileTypeDropdown.SetSize(80, 20)
	fileTypeDropdown.InsertItem(langResx, "ISW")
	fileTypeDropdown.InsertItem(langXliff, "iOS")
	fileTypeDropdown.InsertItem(langXML, "Android")

	createLangCheckboxes(mainWindow, availableLanguages)

	createButton := winc.NewPushButton(mainWindow)
	createButton.SetText("Create Excel")
	createButton.SetSize(290, 40)
	positionButton(createButton, availableLanguages)

	responseLabel := winc.NewLabel(mainWindow)
	responseLabel.SetSize(400, 30)
	responseLabel.Hide()
	positionResponseLabel(responseLabel, availableLanguages)

	pathTextBox.OnChange().Bind(func(e *winc.Event) {
		path := pathTextBox.Text()
		f, err := props.DetectFormat(path)

		if err != nil {
			return
		}

		fileTypeDropdown.SetSelectedItem(f.FileType)
		availableLanguages = f.AvailableLanguages

		createLangCheckboxes(mainWindow, availableLanguages)
		resizeWindow(mainWindow, availableLanguages)
		positionButton(createButton, availableLanguages)
		positionResponseLabel(responseLabel, availableLanguages)
	})

	createButton.OnClick().Bind(func(e *winc.Event) {
		responseLabel.SetText("")
		responseLabel.Hide()
		languagesToInclude := []string{}
		for _, c := range langCheckboxes {
			if c.checkbox.Checked() {
				languagesToInclude = append(languagesToInclude, c.languageCode)
			}
		}
		err := props.OnExecute(fileTypeDropdown.SelectedItem(), pathTextBox.Text(), languagesToInclude)
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

func createLangCheckboxes(parent winc.Controller, availableLanguages []string) {
	for _, c := range langCheckboxes {
		c.checkbox.Close()
	}

	checkboxes := []langCheckbox{}
	for i, l := range availableLanguages {
		c := winc.NewCheckBox(parent)
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

	langCheckboxes = checkboxes
}

func resizeWindow(window *winc.Form, langs []string) {
	window.SetSize(330, 200+(len(langs)/2+1)*30)
}

func positionButton(button *winc.PushButton, langs []string) {
	button.SetPos(10, 70+(len(langs)/2+1)*30)
}

func positionResponseLabel(label *winc.Label, langs []string) {
	label.SetPos(10, 120+(len(langs)/2+1)*30)
}

type langCheckbox struct {
	languageCode string
	checkbox     *winc.CheckBox
}

func wndOnClose(arg *winc.Event) {
	winc.Exit()
}
