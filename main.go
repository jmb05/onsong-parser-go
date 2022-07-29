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
	var output string
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
		case "-o", "--output":
			output = os.Args[i+2]
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

			outputExists := true
			var outputIsDirectory bool

			outputInfo, err := os.Stat(output)
			if err != nil {
				outputExists = false
				outputIsDirectory = !exists(getParent(output))
			} else {
				outputIsDirectory = outputInfo.IsDir()
			}

			if fileInfo.IsDir() {
				if outputExists {
					if outputIsDirectory {
						filesCreated += parseFolder(path, output, templatePath, metadataKeys, recursive, padding, paddingSens, fixUmlauts)
					} else {
						fmt.Println(styling.Style("red", "", "Output has to be a directory (when using a directory as input)"))
					}
				} else {
					fmt.Printf(styling.Style("red", "", "Specified output directory \""))
					fmt.Printf(styling.Style("red", "italic", output))
					fmt.Printf(styling.Style("red", "", "\" does not exist!\n"))
					break
				}
			} else {
				if !outputExists {
					if parseOnsongFile(path, output, templatePath, metadataKeys, padding, paddingSens, fixUmlauts) {
						filesCreated++
					}
				} else if outputIsDirectory {
					appendedOutPath := appendFilePaths(output, strings.Replace(getFileName(path), ".onsong", ".html", -1))
					if parseOnsongFile(path, appendedOutPath, templatePath, metadataKeys, padding, paddingSens, fixUmlauts) {
						filesCreated++
					}
				} else {
					fmt.Printf(styling.Style("yellow", "", "Specified output file \""))
					fmt.Printf(styling.Style("yellow", "italic", output))
					fmt.Printf(styling.Style("yellow", "", "\" already exists!\n"))
					fmt.Printf(styling.Style("yellow", "", "Do you want to override? "))
					if readYesNo(false) {
						if parseOnsongFile(path, output, templatePath, metadataKeys, padding, paddingSens, fixUmlauts) {
							filesCreated++
						}
					}
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

func parseFolder(path string, outputPath string, templatePath string, metadataKeys []string, recursive bool, padding int, paddingSens int, fixUmlauts bool) int {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fileError(err, path)
		return 0
	}
	if !exists(outputPath) {
		outputPath = path
	}
	var filesCreated int
	for _, file := range files {
		filePath := appendFilePaths(path, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fileError(err, filePath)
			continue
		}

		if fileInfo.IsDir() {
			if recursive {
				appendedOutPath := appendFilePaths(outputPath, fileInfo.Name())
				if !exists(appendedOutPath) {
					os.Mkdir(appendedOutPath, 0666)
					fmt.Printf("Created missing subdirectory \"")
					fmt.Printf(styling.Style("", "italic", appendedOutPath))
					fmt.Printf("\"\n")
				}
				parseFolder(filePath, appendedOutPath, templatePath, metadataKeys, recursive, padding, paddingSens, fixUmlauts)
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

		appendedOutPath := appendFilePaths(outputPath, strings.Replace(fileInfo.Name(), ".onsong", ".html", -1))
		if parseOnsongFile(filePath, appendedOutPath, templatePath, metadataKeys, padding, paddingSens, fixUmlauts) {
			filesCreated++
		}
	}
	return filesCreated
}

func parseOnsongFile(path string, outputPath string, templatePath string, metadataKeys []string, padding int, paddingSens int, fixUmlauts bool) bool {
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
	os.WriteFile(outputPath, []byte(html), 0666)
	fmt.Printf("Created File: ")
	fmt.Printf(styling.Style("green", "", "\""))
	fmt.Printf(styling.Style("green", "italic", outputPath))
	fmt.Printf(styling.Style("green", "italic", "\"\n"))
	return true
}

func appendFilePaths(p1 string, p2 string) string {
	if strings.HasSuffix(p1, "/") {
		return p1 + p2
	} else {
		return p1 + "/" + p2
	}
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

func getParent(path string) string {
	builder := strings.Builder{}
	parts := strings.Split(path, "/")
	for _, p := range parts[:len(parts)-1] {
		builder.WriteString(p + "/")
	}
	return builder.String()
}

func getFileName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func readYesNo(defaultYes bool) bool {
	var buf string
	if defaultYes {
		fmt.Print("[Y/n] ")
		fmt.Scanln(&buf)
		if buf == "Y" || buf == "y" || buf == "" {
			return true
		}
	} else {
		fmt.Print("[y/N] ")
		fmt.Scanln(&buf)
		if buf == "Y" || buf == "y" {
			return true
		}
	}
	return false
}

func printHelp() {
	fmt.Println("\nUsage: Onsong-Parser-go [OPTION]... [FILE/FOLDER]... ")
	fmt.Println("Parses *.onsong files to *.html files\n")
	fmt.Println("Options:")
	fmt.Println("    --dont-fix-umlauts    don't fix the broken OnSong filenames")
	fmt.Println("-h, --help                show this info")
	fmt.Println("-m, --metadata-tags       which metadata tags should be shown ")
	fmt.Println("                          (e.g.: \"Key Duration Keywords\")")
	fmt.Println("-o, --output              specify output file/directory")
	fmt.Println("-p, --padding-size        size of the padding between chords (per character)")
	fmt.Println("-r, --recursive           search recursive (in subfolders)")
	fmt.Println("    --padding-sensitivity change the padding sensitivity")
	fmt.Println("-t, --template            choose a custom template")
}
