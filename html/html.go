package html

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/jmb05/Onsong-Parser-go/onsong"
	"github.com/jmb05/styling"
)

func CreateHtml(song onsong.Song, templatePath string) string {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		parseError(err)
		return ""
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, &song)
	if err != nil {
		parseError(err)
		return ""
	}
	return buffer.String()
}

func parseError(err error) {
	fmt.Println(styling.Style("red", "regular", "Error applying template"))
	fmt.Println(err)
}
