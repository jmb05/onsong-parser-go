package html

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/jmb05/onsong-parser-go/onsong"
	"github.com/jmb05/styling"
)

func CreateSongHtml(song onsong.Song, templatePath string) string {
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

type SongLibrary struct {
	Title string
	Songs []SongEntry
}

type SongEntry struct {
	Name     string
	Location string
}

func CreateListHtml(songLibrary SongLibrary, templatePath string) string {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		parseError(err)
		return ""
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, &songLibrary)
	if err != nil {
		parseError(err)
		return ""
	}
	return buffer.String()
}

func parseError(err error) {
	fmt.Println(styling.Style("red", "regular", "Error applying template"))
	fmt.Println(styling.Style("red", "regular", err.Error()))
}
