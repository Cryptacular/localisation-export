package xliff

import "testing"

func TestConvertToXliffReadsCorrectNumberOfEntries(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?>
	<xliff xmlns="urn:oasis:names:tc:xliff:document:1.2" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" version="1.2" xsi:schemaLocation="urn:oasis:names:tc:xliff:document:1.2 http://docs.oasis-open.org/xliff/v1.2/os/xliff-core-1.2-strict.xsd">
	  <file original="Healthzone/Classes/Application/Base.lproj/Healthzone.storyboard" source-localisation.Language="en" target-localisation.Language="es" datatype="plaintext">
		<header>
		  <tool tool-id="com.apple.dt.xcode" tool-name="Xcode" tool-version="10.2.1" build-num="10E1001"/>
		</header>
		<body>
		  <trans-unit id="4pM-xO-tTa.title">
			<source>GOALS</source>
			<target>OBJETIVOS</target>
			<note>Class = "UITabBarItem"; title = "GOALS"; ObjectID = "4pM-xO-tTa";</note>
		  </trans-unit>
		</body>
	  </file>
	</xliff>`

	actual := convertToXliff([]byte(x))

	if length := len(actual); length != 1 {
		t.Errorf(`Actual: %d; Expected: 1`, length)
	}
}
