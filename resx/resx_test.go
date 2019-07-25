package resx

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

	if length := len(actual); length != 3 {
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
	one := actual[0]
	two := actual[1]
	three := actual[2]

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

	if length := len(actual); length != 1 {
		t.Errorf(`Actual: %d; Expected: 1`, length)
	}

	if key := actual[0].Key; key != "ThisIsTheKey" {
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

	if value := actual[0].Value; value != "This is the value of the string" {
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

	if comment := actual[0].Comment; comment != "Now this would be the comment" {
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

	valueOne := actual[0].Value
	valueTwo := actual[1].Value

	if valueOne != "" {
		t.Errorf(`Actual: "%s"; Expected: ""`, valueOne)
	} else if valueTwo != "" {
		t.Errorf(`Actual: "%s"; Expected: ""`, valueTwo)
	}
}
