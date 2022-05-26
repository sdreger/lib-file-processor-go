package parser

import (
	"github.com/mantidtech/wordnumber"
	"regexp"
	"strconv"
	"strings"
)

var (
	// 1st Edition / 5th Edition
	editionCardinalRegex = regexp.MustCompile(`(?i),? ?\(?(\d+)(st|nd|rd|th) (Edition)\)?`)
	// Second Edition / Fifth Edition
	editionOrdinalRegex = regexp.MustCompile(`(?i),? \(?([a-zA-Z]+) (Edition)\)?`)
)

// ParseTitleString parses a book title string and returns separate title and subtitle strings.
func ParseTitleString(titleString string) (string, string) {
	// Try to cut out the ordinal 'edition' part, like '3rd Edition'
	cardinalEditionSubMatch := editionCardinalRegex.FindStringSubmatch(titleString)
	if cardinalEditionSubMatch != nil {
		titleString = strings.ReplaceAll(titleString, cardinalEditionSubMatch[0], "")
	}

	// Try to cut out the cardinal 'edition' part, like 'Second Edition'
	ordinalEditionSubMatch := editionOrdinalRegex.FindStringSubmatch(titleString)
	if ordinalEditionSubMatch != nil {
		ordinalValue, _ := wordnumber.OrdinalToInt(ordinalEditionSubMatch[1])
		if ordinalValue > 0 {
			titleString = strings.ReplaceAll(titleString, ordinalEditionSubMatch[0], "")
		}
	}

	// Extract 'titleString' and 'subtitle' values
	lastColonIndex := strings.LastIndex(titleString, ": ")
	if lastColonIndex == -1 {
		return titleString, ""
	}

	return titleString[:lastColonIndex], titleString[lastColonIndex+2:]
}

// ParseEditionString parses a book edition string and returns its numeric value.
func ParseEditionString(editionString string) (uint8, error) {
	// Try to parse cardinal value like '2nd Edition'
	cardinalEditionSubMatch := editionCardinalRegex.FindStringSubmatch(editionString)
	if cardinalEditionSubMatch != nil {
		ed, _ := strconv.Atoi(cardinalEditionSubMatch[1])
		return uint8(ed), nil
	}

	// Try to parse ordinal edition value like 'Second Edition'
	ordinalEditionSubMatch := editionOrdinalRegex.FindStringSubmatch(editionString)
	if ordinalEditionSubMatch != nil {
		ordinalValue, err := wordnumber.OrdinalToInt(ordinalEditionSubMatch[1])
		if err != nil {
			return 0, err
		}
		return uint8(ordinalValue), nil
	}

	return 0, nil
}
