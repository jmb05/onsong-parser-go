package html

import (
	"bytes"
	"html/template"
	"log"

	"github.com/jmb05/Onsong-Parser-go/onsong"
)

func CreateHtml(song onsong.Song) string {
	tmpl, err := template.ParseFiles("./html/song.gohtml")
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
