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

const DEFAULT_TEMPLATE_PATH = "html/template.gohtml"
const DEFAULT_PADDING = 15

func readFile(path string) string {
	content, err := utfutil.ReadFile(path, utfutil.UTF8)
	if err != nil {
		fileError(err, path)
		return ""
	}
	return string(content)
}

func main() {
	metadataKeys := []string{"Key", "Time", "Tempo"}
	fmt.Println(styling.ColorBold("white", "Onsong to HTML Parser"))
	fmt.Println("Copyright (C) 2022, Josiah Bennett, Jared M. Bennett")
	var skip int
	templatePath := DEFAULT_TEMPLATE_PATH
	onsongFiles := []string{}
	recursive := false
	padding := DEFAULT_PADDING
	paddingSens := 0
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
				templatePath = DEFAULT_TEMPLATE_PATH
			}
		} else if arg == "-r" || arg == "--recursive" {
			recursive = true
		} else if arg == "-m" || arg == "--metadata-tags" {
			metadataKeys = strings.Split(os.Args[i+2], " ")
			skip = 1
		} else if arg == "-p" || arg == "--padding-size" {
			paddingN, err := strconv.Atoi(os.Args[i+2])
			padding = paddingN
			if err != nil {
				panic(err)
			}
			skip = 1
		} else if arg == "--padding-sensitivity" {
			paddingSensCp, err := strconv.Atoi(os.Args[i+2])
			paddingSens = paddingSensCp
			if err != nil {
				panic(err)
			}
			skip = 1
		} else if arg == "-h" || arg == "--help" {
			fmt.Println("\nUsage: Onsong-Parser-go [OPTION]... [FILE/FOLDER]... ")
			fmt.Println("Parses *.onsong files to *.html files\n")
			fmt.Println("Options:")
			fmt.Println("-h, --help\t\t\t show this info")
			fmt.Println("-m, --metadata-tags\t\t which metadata tags should be shown ")
			fmt.Println("\t\t\t\t (e.g.: \"Key Duration Keywords\")")
			fmt.Println("-r, --recursive\t\t\t search recursive (in subfolders)")
			fmt.Println("-p, --padding-size\t\t size of the padding between chords (per character)")

			fmt.Println("    --padding-sensitivity\t change the padding sensitivity")
			os.Exit(0)
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
				filesCreated += parseFolder(path, templatePath, metadataKeys, recursive, padding, paddingSens)
			} else {
				if parseOnsongFile(path, templatePath, metadataKeys, padding, paddingSens) {
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

func parseFolder(path string, templatePath string, metadataKeys []string, recursive bool, padding int, paddingSens int) int {
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
				parseFolder(filePath, templatePath, metadataKeys, recursive, padding, paddingSens)
			} else {
				fmt.Println("Skipping directory: \"" + styling.ColorItalic("white", filePath) + "\" Add parameter \"-r\" to parse recursively")
				continue
			}
		}

		if !strings.HasSuffix(filePath, ".onsong") {
			continue
		}

		if parseOnsongFile(filePath, templatePath, metadataKeys, padding, paddingSens) {
			filesCreated++
		}
	}
	return filesCreated
}

func parseOnsongFile(path string, templatePath string, metadataKeys []string, padding int, paddingSens int) bool {
	if !strings.HasSuffix(path, ".onsong") {
		fmt.Println(styling.Color("yellow", "Warning: File \"") + styling.ColorItalic("yellow", path) + styling.Color("yellow", "\" doesn't have \".onsong\" ending"))
	}
	song, success := onsong.Parse(readFile(path), metadataKeys, padding, paddingSens)
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
