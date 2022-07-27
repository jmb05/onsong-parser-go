package main

import (
	"fmt"
	"io/ioutil"
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

func readFile(path string) string {
	content, err := utfutil.ReadFile(path, utfutil.UTF8)
	if err != nil {
		fileError(err, path)
		return ""
	}
	return string(content)
}

func main() {
	fmt.Println(styling.ColorBold("white", "Onsong to HTML Parser"))
	fmt.Println("Copyright (C) 2022, Bennett")
	var skip int
	templatePath := defaultTemplatePath
	onsongFiles := []string{}
	recursive := false
	for i, arg := range os.Args[1:] {
		if skip > 0 {
			skip--
			continue
		}
		if arg == "-t" {
			templatePath = os.Args[i+2]
			skip = 1
			if !exists(templatePath) {
				fmt.Println(styling.Color("yellow", "Warning: Template \"") + styling.ColorItalic("yellow", templatePath) + styling.Color("yellow", "\" does not exist! Using default..."))
				templatePath = defaultTemplatePath
			}
		} else if arg == "-r" {
			recursive = true
		} else {
			onsongFiles = append(onsongFiles, arg)
		}
	}
	if len(onsongFiles) > 0 {
		fmt.Println("Using Template: " + styling.Color("cyan", "\"") + styling.ColorItalic("cyan", templatePath) + styling.Color("cyan", "\""))
		var filesCreated int
		for _, path := range onsongFiles {
			fileInfo, err := os.Stat(path)
			if err != nil {
				fileError(err, path)
				continue
			}

			if fileInfo.IsDir() {
				filesCreated += parseFolder(path, templatePath, recursive)
			} else {
				if parseOnsongFile(path, templatePath) {
					filesCreated++
				}
			}
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

func parseFolder(path string, templatePath string, recursive bool) int {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fileError(err, path)
		return 0
	}
	var filesCreated int
	for _, file := range files {
		var filePath string
		if strings.HasSuffix(path, "/") {
			filePath = path + file.Name()
		} else {
			filePath = path + "/" + file.Name()
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fileError(err, filePath)
			continue
		}

		if fileInfo.IsDir() {
			if recursive {
				parseFolder(filePath, templatePath, recursive)
			} else {
				fmt.Println("Skipping directory: \"" + styling.ColorItalic("white", filePath) + "\" Add parameter \"-r\" to parse recursively")
				continue
			}
		}

		if !strings.HasSuffix(filePath, ".onsong") {
			continue
		}

		if parseOnsongFile(filePath, templatePath) {
			filesCreated++
		}
	}
	return filesCreated
}

func parseOnsongFile(path string, templatePath string) bool {
	if !strings.HasSuffix(path, ".onsong") {
		fmt.Println(styling.Color("yellow", "Warning: File \"") + styling.ColorItalic("yellow", path) + styling.Color("yellow", "\" doesn't have \".onsong\" ending"))
	}
	song, success := onsong.Parse(readFile(path))
	if !success {
		fmt.Println(styling.Color("yellow", "Skipping..."))
		return false
	}
	html := html.CreateHtml(song, templatePath)
	htmlPath := strings.Replace(path, ".onsong", ".html", 1)
	os.WriteFile(htmlPath, []byte(html), 0666)
	fmt.Println("Created File: " + styling.Color("green", "\"") + styling.ColorItalic("green", htmlPath) + styling.Color("green", "\""))
	return true
}

func fileError(err error, path string) {
	fmt.Println(styling.Color("red", "Error reading file: \"") + styling.ColorItalic("red", path) + styling.Color("red", "\""))
	log.Println(err)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
