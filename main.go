package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/TomOnTime/utfutil"
	"github.com/jmb05/onsong-parser-go/html"
	"github.com/jmb05/onsong-parser-go/onsong"
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
	fmt.Println(styling.Style("white", "bold", "Onsong to HTML Parser"))
	fmt.Println("Copyright (C) 2022, Josiah Bennett, Jared M. Bennett")
	var skip int
	templatePath := DEFAULT_TEMPLATE_PATH
	onsongFiles := []string{}
	recursive := false
	padding := DEFAULT_PADDING
	paddingSens := 0
	fixUmlauts := true
	for i, arg := range os.Args[1:] {
		if skip > 0 {
			skip--
			continue
		}
		switch arg {
		case "--dont-fix-umlauts":
			fixUmlauts = false
		case "-t", "--template":
			templatePath = os.Args[i+2]
			skip = 1
			if !exists(templatePath) {
				fmt.Printf(styling.Style("yellow", "", "Warning: Template \""))
				fmt.Printf(styling.Style("yellow", "italic", templatePath))
				fmt.Printf(styling.Style("yellow", "", "\" does not exist! Using default...\n"))
				templatePath = DEFAULT_TEMPLATE_PATH
			}
		case "-r", "--recursive":
			recursive = true
		case "-m", "--metadata-tags":
			metadataKeys = strings.Split(os.Args[i+2], " ")
			skip = 1
		case "-p", "--padding-size":
			paddingCp, err := strconv.Atoi(os.Args[i+2])
			padding = paddingCp
			if err != nil {
				panic(err)
			}
			skip = 1
		case "--padding-sensitivity":
			paddingSensCp, err := strconv.Atoi(os.Args[i+2])
			paddingSens = paddingSensCp
			if err != nil {
				panic(err)
			}
			skip = 1
		case "-h", "--help":
			printHelp()
			os.Exit(0)
		default:
			onsongFiles = append(onsongFiles, arg)
		}
	}
	if len(onsongFiles) > 0 {
		fmt.Printf("Using Template: " + styling.Style("cyan", "", "\""))
		fmt.Printf(styling.Style("cyan", "italic", templatePath))
		fmt.Printf(styling.Style("cyan", "", "\"\n"))
		var filesCreated int
		for _, path := range onsongFiles {
			fileInfo, err := os.Stat(path)
			if err != nil {
				fileError(err, path)
				continue
			}

			if fileInfo.IsDir() {
				filesCreated += parseFolder(path, templatePath, metadataKeys, recursive, padding, paddingSens, fixUmlauts)
			} else {
				if parseOnsongFile(path, templatePath, metadataKeys, padding, paddingSens, fixUmlauts) {
					filesCreated++
				}
			}
		}
		if filesCreated == 1 {
			fmt.Println(styling.Style("green", "bold", "Created "+strconv.Itoa(filesCreated)+" File"))
		} else {
			fmt.Println(styling.Style("green", "bold", "Created "+strconv.Itoa(filesCreated)+" Files"))
		}
	} else {
		fmt.Println(styling.Style("red", "bold", "No Files selected"))
	}
}

func parseFolder(path string, templatePath string, metadataKeys []string, recursive bool, padding int, paddingSens int, fixUmlauts bool) int {
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
				parseFolder(filePath, templatePath, metadataKeys, recursive, padding, paddingSens, fixUmlauts)
			} else {
				fmt.Printf("Skipping directory: \"")
				fmt.Printf(styling.Style("white", "italic", filePath))
				fmt.Printf("\" Add parameter \"-r\" to parse recursively\n")
				continue
			}
		}

		if !strings.HasSuffix(filePath, ".onsong") {
			continue
		}

		if parseOnsongFile(filePath, templatePath, metadataKeys, padding, paddingSens, fixUmlauts) {
			filesCreated++
		}
	}
	return filesCreated
}

func parseOnsongFile(path string, templatePath string, metadataKeys []string, padding int, paddingSens int, fixUmlauts bool) bool {
	if !strings.HasSuffix(path, ".onsong") {
		fmt.Printf(styling.Style("yellow", "", "Warning: File \""))
		fmt.Printf(styling.Style("yellow", "italic", path))
		fmt.Printf(styling.Style("yellow", "", "\" doesn't have \".onsong\" ending!\n"))
	}
	song, success := onsong.Parse(readFile(path), metadataKeys, padding, paddingSens)
	if !success {
		fmt.Println(styling.Style("yellow", "", "Skipping..."))
		return false
	}
	html := html.CreateHtml(song, templatePath)
	if fixUmlauts {
		path = replaceUmlauts(path)
	}
	htmlPath := strings.Replace(path, ".onsong", ".html", 1)
	os.WriteFile(htmlPath, []byte(html), 0666)
	fmt.Printf("Created File: ")
	fmt.Printf(styling.Style("green", "", "\""))
	fmt.Printf(styling.Style("green", "italic", htmlPath))
	fmt.Printf(styling.Style("green", "italic", "\"\n"))
	return true
}

//replaces the umlauts from OnSongs (imo broken) filenames
//can be diabled with "--dont-fix-umlauts"
//(and yes "umlauts" is the correct english plural, check if you don't trust me)
func replaceUmlauts(s string) string {
	s = strings.Replace(s, "u"+string([]byte{226, 149, 160, 208, 152}), "ü", -1)
	s = strings.Replace(s, "o"+string([]byte{226, 149, 160, 208, 152}), "ö", -1)
	s = strings.Replace(s, "a"+string([]byte{226, 149, 160, 208, 152}), "ä", -1)
	s = strings.Replace(s, string([]byte{226, 148, 156, 208, 175}), "ß", -1)
	return s
}

func fileError(err error, path string) {
	fmt.Printf(styling.Style("red", "", "Error reading file: \""))
	fmt.Printf(styling.Style("red", "italic", path))
	fmt.Printf(styling.Style("red", "", "\"\n"))
	log.Println(err)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func printHelp() {
	fmt.Println("\nUsage: Onsong-Parser-go [OPTION]... [FILE/FOLDER]... ")
	fmt.Println("Parses *.onsong files to *.html files\n")
	fmt.Println("Options:")
	fmt.Println("    --dont-fix-umlauts    don't fix the broken OnSong filenames")
	fmt.Println("-h, --help                show this info")
	fmt.Println("-m, --metadata-tags       which metadata tags should be shown ")
	fmt.Println("                          (e.g.: \"Key Duration Keywords\")")
	fmt.Println("-r, --recursive           search recursive (in subfolders)")
	fmt.Println("-p, --padding-size        size of the padding between chords (per character)")
	fmt.Println("    --padding-sensitivity change the padding sensitivity")
	fmt.Println("-t, --template            choose a custom template")
}
