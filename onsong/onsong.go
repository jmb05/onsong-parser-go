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

func Parse(content string) (Song, bool) {
	if strings.TrimSpace(content) == "" {
		return Song{}, false
	}
	paragraphs := SplitParagraphs(content)
	metadata, copyright := ParseMetadata(paragraphs)
	song := Song{
		Title:     ParseTitle(paragraphs),
		Artist:    ParseArtist(paragraphs),
		Meta:      metadata,
		Sections:  ParseSections(paragraphs),
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

func isBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

func ParseTitle(paragraphs [][]string) string {
	return strings.Replace(paragraphs[0][0], string([]byte{255, 254}), "", 1)
}

func ParseArtist(paragraphs [][]string) string {
	return paragraphs[0][1]
}

func ParseSections(paragraphs [][]string) []Section {
	sections := []Section{}
	if len(paragraphs) > 1 {
		for i := 1; i < len(paragraphs); i++ {
			if len(paragraphs[i]) > 1 {
				lines := []Line{}

				for _, lineStr := range paragraphs[i][1:] {
					line := Line{
						Parts: parseLineParts(lineStr),
					}
					lines = append(lines, line)
				}

				section := Section{
					Title: paragraphs[i][0],
					Lines: lines,
				}
				sections = append(sections, section)
			}
		}
	}
	return sections
}

func parseLineParts(lineStr string) []LinePart {
	lineParts := []LinePart{}
	splitLineU := splitLine(lineStr)
	splitLine := removeBlankParts(splitLineU)
	for i, part := range splitLine {
		if isChord(part) {
			padding := 0
			if i > 0 {
				partBefore := splitLine[i-1]
				if isChord(partBefore) {
					padding = len(partBefore) * 15
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

func ParseMetadata(paragraphs [][]string) ([]string, string) {
	lines := []string{}
	var copyright string
	for i := 2; i < len(paragraphs[0]); i++ {
		if strings.HasPrefix(paragraphs[0][i], "Copyright") {
			copyright = strings.TrimSpace(strings.Split(paragraphs[0][i], ":")[1])
		} else if !strings.HasPrefix(paragraphs[0][i], "Keywords") && !strings.HasPrefix(paragraphs[0][i], "CCLI") {
			lines = append(lines, paragraphs[0][i])
		}
	}
	return lines, copyright
}
