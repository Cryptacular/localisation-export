package excel

import (
	"fmt"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// Write takes a 2-dimensional slice of strings (rows and columns) to create an Excel file
func Write(content [][]string, path string) {
	f := excelize.NewFile()

	for i, row := range content {
		cell := "A" + strconv.Itoa(i+1)
		f.SetSheetRow("Sheet1", cell, &row)
	}

	err := f.SaveAs(path)
	if err != nil {
		fmt.Println(err)
	}
}
