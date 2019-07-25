# Localisation Export

This is a Windows GUI tool to convert a number of localisation files for any number of languages into an Excel file with keys and values for each of those languages.

## Build

```shell
go build -ldflags="-H windowsgui"
```

## Run

1. Open `localisation-export.exe` after building
2. Copy the path of the root folder containing `UIStrings.<lang>.RESX` files, or `<lang>.xcloc` folders (XLIFF/iOS), or `values-<lang>` folders (XML/Android)
3. Paste it in the Folder Path input field
4. If it's valid, the type (RESX, iOS or Android) will appear in the dropdown and checkboxes appear with language codes
5. Select the languages you want to export (EN is included by default)
6. Click the Create Excel button
7. You can find the Excel file in the root folder you provided in step 3
