package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/TomOnTime/utfutil"
	"github.com/jmb05/Onsong-Parser-go/html"
	"github.com/jmb05/Onsong-Parser-go/onsong"
	"github.com/jmb05/styling"
)

const defaultTemplatePath = "html/template.gohtml"

func readFile(filePath string) string {
	content, err := utfutil.ReadFile(filePath, utfutil.UTF8)
	if err != nil {
		fmt.Println(styling.Color("red", "Error reading file: \"") + styling.ColorItalic("red", filePath) + styling.Color("red", "\""))
		log.Println(err)
		return ""
	}
	return string(content)
}

func main() {
	fmt.Println(styling.ColorBold("white", "Onsong to HTML Parser"))
	fmt.Println("Copyright (C) 2022, Jared M. Bennett")
	var skip int
	templatePath := defaultTemplatePath
	onsongFiles := []string{}
	for i, arg := range os.Args[1:] {
		if skip > 0 {
			skip--
			continue
		}
		if arg == "-t" {
			templatePath = os.Args[i+1]
			skip = 1
		} else {
			onsongFiles = append(onsongFiles, arg)
		}
	}
	if len(onsongFiles) > 0 {
		fmt.Println("Using Template: " + styling.Color("cyan", "\"") + styling.ColorItalic("cyan", templatePath) + styling.Color("cyan", "\""))
		var filesCreated int
		for _, path := range onsongFiles {
			if !strings.HasSuffix(path, ".onsong") {
				fmt.Println(styling.Color("yellow", "Warning: File \"") + styling.ColorItalic("yellow", path) + styling.Color("yellow", "\" doesn't have \".onsong\" ending"))
			}
			song, success := onsong.Parse(readFile(path))
			if !success {
				fmt.Println(styling.Color("yellow", "Skipping..."))
				continue
			}
			html := html.CreateHtml(song, templatePath)
			htmlPath := strings.Replace(path, ".onsong", ".html", 1)
			os.WriteFile(htmlPath, []byte(html), 0666)
			fmt.Println("Created File: " + styling.Color("green", "\"") + styling.ColorItalic("green", htmlPath) + styling.Color("green", "\""))
			filesCreated++
		}
		if filesCreated > 1 {
			fmt.Println(styling.ColorBold("green", "Created "+strconv.Itoa(filesCreated)+" Files"))
		} else {
			fmt.Println(styling.ColorBold("green", "Created "+strconv.Itoa(filesCreated)+" File"))
		}
	} else {
		fmt.Println(styling.ColorBold("red", "No Files selected"))
	}
}
