package parser

import (
	"github.com/mantidtech/wordnumber"
	"regexp"
	"strings"
)

var (
	// 1st Edition / 5th Edition
	editionCardinalRegex = regexp.MustCompile(`(?i),? \(?(\d+)(st|nd|rd|th) (Edition)\)?`)
	// Second Edition / Fifth Edition
	editionOrdinalRegex = regexp.MustCompile(`(?i),? \(?([a-zA-Z]+) (Edition)\)?`)
)

// ParseTitle parses a book title string and returns separate title and subTitle strings.
func ParseTitle(title string) (string, string) {
	// Try to cut out the ordinal 'edition' part, like '3rd Edition'
	cardinalEditionSubMatch := editionCardinalRegex.FindStringSubmatch(title)
	if cardinalEditionSubMatch != nil {
		title = strings.ReplaceAll(title, cardinalEditionSubMatch[0], "")
	}

	// Try to cut out the cardinal 'edition' part, like 'Second Edition'
	ordinalEditionSubMatch := editionOrdinalRegex.FindStringSubmatch(title)
	if ordinalEditionSubMatch != nil {
		ordinalValue, _ := wordnumber.OrdinalToInt(ordinalEditionSubMatch[1])
		if ordinalValue > 0 {
			title = strings.ReplaceAll(title, ordinalEditionSubMatch[0], "")
		}
	}

	// Extract 'title' and 'subTitle' values
	lastColonIndex := strings.LastIndex(title, ": ")
	if lastColonIndex == -1 {
		return title, ""
	}

	return title[:lastColonIndex], title[lastColonIndex+2:]
}
