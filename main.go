package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/jmb05/Onsong-Parser-go/html"
	"github.com/jmb05/Onsong-Parser-go/onsong"
)

func readFile(filePath string) string {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(content)
}

func main() {
	song := onsong.Parse(readFile("test-song.onsong"))
	html := html.CreateHtml(song)
	os.WriteFile("song.html", []byte(html), 0666)
}
