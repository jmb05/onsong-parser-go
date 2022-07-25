package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jmb05/Onsong-Parser-go/html"
	"github.com/jmb05/Onsong-Parser-go/onsong"
	"github.com/jmb05/styling"
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
	fmt.Println(styling.Color("green", ""))
	song := onsong.Parse(readFile("test-song.onsong"))
	html := html.CreateHtml(song)
	os.WriteFile("song.html", []byte(html), 0666)
}
