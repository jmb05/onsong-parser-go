package html

import (
	"bytes"
	"log"
	"text/template"

	"github.com/jmb05/Onsong-Parser-go/onsong"
)

func CreateHtml(song onsong.Song, templatePath string) string {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Panic(err)
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, &song)
	if err != nil {
		log.Panic(err)
	}
	return buffer.String()
}
