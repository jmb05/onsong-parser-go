package onsong

import (
	"strings"
)

type Song struct {
	Title     string
	Artist    string
	Meta      []string
	Sections  []Section
	Copyright string
}

type Section struct {
	Title string
	Lines []Line
}

type Line struct {
	Parts []LinePart
}

type LinePart struct {
	IsChord bool
	Chord   Chord
	Text    string
}

type Chord struct {
	Value   string
	Padding int
}

func Parse(content string, metadataKeys []string, defaultPadding int, paddingSensitivity int) (Song, bool) {
	if isBlank(content) {
		return Song{}, false
	}
	paragraphs := SplitParagraphs(content)
	metadata, copyright := ParseMetadata(paragraphs[0], metadataKeys)
	song := Song{
		Title:     paragraphs[0][0],
		Artist:    paragraphs[0][1],
		Meta:      metadata,
		Sections:  ParseSections(paragraphs, defaultPadding, paddingSensitivity),
		Copyright: copyright,
	}
	return song, true
}

func SplitParagraphs(content string) [][]string {
	lines := strings.Split(content, "\n")
	out := [][]string{}
	currentParagraph := []string{}
	for _, line := range lines {
		if isBlank(line) {
			out = append(out, currentParagraph)
			currentParagraph = []string{}
		} else {
			currentParagraph = append(currentParagraph, line)
		}
	}
	out = append(out, currentParagraph)
	return out
}

func ParseMetadata(paragraph []string, includedKeys []string) ([]string, string) {
	metadata := []string{}
	var copyright string
	for _, line := range paragraph[2:] {
		if strings.HasPrefix(line, "Copyright") {
			copyright = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if doesKeyExist(includedKeys, strings.Split(line, ":")[0]) {
			metadata = append(metadata, line)
		}
	}
	return metadata, copyright
}

func doesKeyExist(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func ParseSections(paragraphs [][]string, defaultPadding int, paddingSensitivity int) []Section {
	sections := []Section{}
	if len(paragraphs) > 1 {
		for i := 1; i < len(paragraphs); i++ {
			if len(paragraphs[i]) > 1 {
				lines := []Line{}

				var title string
				var startFrom int
				if isSectionTitle(paragraphs[i][0]) {
					startFrom = 1
					title = paragraphs[i][0]
				} else {
					startFrom = 0
				}

				for _, lineStr := range paragraphs[i][startFrom:] {
					line := Line{
						Parts: parseLineParts(lineStr, defaultPadding, paddingSensitivity),
					}
					lines = append(lines, line)
				}

				section := Section{
					Title: title,
					Lines: lines,
				}
				sections = append(sections, section)
			}
		}
	}
	return sections
}

func isSectionTitle(line string) bool {
	return !(strings.Contains(line, "[") && strings.Contains(line, "]"))
}

func parseLineParts(lineStr string, defaultPadding int, paddingSensitivity int) []LinePart {
	lineParts := []LinePart{}
	splitLineU := splitLine(lineStr)
	splitLine := removeBlankParts(splitLineU)
	for i, part := range splitLine {
		if isChord(part) {
			padding := 0
			if i > 0 {
				partBefore := splitLine[i-1]
				if isChord(partBefore) {
					padding = len(strings.TrimSpace(partBefore)) * defaultPadding
				} else if i > 1 {
					chordBefore := splitLine[i-2]
					padding = (len(chordBefore) - len(partBefore) + (paddingSensitivity - 1)) * defaultPadding
					if padding < 0 {
						padding = 0
					}
				}
			}
			lineParts = append(lineParts, LinePart{
				IsChord: true,
				Chord:   Chord{strings.Replace(strings.Replace(part, "[", "", 1), "]", "", 1), padding},
			})
		} else {
			lineParts = append(lineParts, LinePart{
				Text: part,
			})
		}
	}
	return lineParts
}

func isChord(part string) bool {
	return strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]")
}

func removeBlankParts(partsIn []string) []string {
	outParts := []string{}
	for _, part := range partsIn {
		if len(part) != 0 {
			outParts = append(outParts, part)
		}
	}
	return outParts
}

func splitLine(line string) []string {
	parts1 := splitBefore(line, "[")
	parts := []string{}
	for _, part := range parts1 {
		parts2 := strings.SplitAfter(part, "]")
		for _, p := range parts2 {
			parts = append(parts, p)
		}
	}
	return parts
}

func splitBefore(in string, delimeter string) []string {
	split := strings.Split(in, delimeter)
	out := []string{split[0]}
	for _, part := range split[1:] {
		out = append(out, delimeter+part)
	}
	return out
}

func isBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}
