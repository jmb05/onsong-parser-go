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
	Lines []string
}

func Parse(content string) Song {
	paragraphs := SplitParagraphs(content)
	song := Song{
		Title:     ParseTitle(paragraphs),
		Artist:    ParseArtist(paragraphs),
		Meta:      ParseMetadata(paragraphs),
		Sections:  ParseSections(paragraphs),
		Copyright: "Test",
	}
	return song
}

func SplitParagraphs(content string) [][]string {
	lines := strings.Split(content, "\n")
	out := [][]string{}
	currentParagraph := []string{}
	for i := range lines {
		if len(strings.TrimSpace(lines[i])) == 0 {
			out = append(out, currentParagraph)
			currentParagraph = []string{}
		} else {
			currentParagraph = append(currentParagraph, lines[i])
		}
	}
	out = append(out, currentParagraph)
	return out
}

func ParseTitle(paragraphs [][]string) string {
	return paragraphs[0][0]
}

func ParseArtist(paragraphs [][]string) string {
	return paragraphs[0][1]
}

func ParseSections(paragraphs [][]string) []Section {
	sections := []Section{}
	if len(paragraphs) > 1 {
		for i := 1; i < len(paragraphs); i++ {
			if len(paragraphs[i]) > 1 {
				section := Section{
					Title: paragraphs[i][0],
					Lines: paragraphs[i][1:],
				}
				sections = append(sections, section)
			}
		}
	}
	return sections
}

func ParseMetadata(paragraphs [][]string) []string {
	lines := []string{}
	for i := 2; i < len(paragraphs[0]); i++ {
		lines = append(lines, paragraphs[0][i])
	}
	return lines
}

func splitFunc(r rune) bool {
	return r == '[' || r == ']'
}
