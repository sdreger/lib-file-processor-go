package parser

import (
	"fmt"
	"github.com/mantidtech/wordnumber"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// 1st Edition / 5th Edition
	editionCardinalRegex = regexp.MustCompile(`(?i),? ?\(?(\d+)(st|nd|rd|th) (Edition)\)?`)
	// Second Edition / Fifth Edition
	editionOrdinalRegex = regexp.MustCompile(`(?i),? ?\(?([a-zA-Z]{3,}) (Edition)\)?`)
	// Apress; 1st ed. edition (November 1, 2020) / Wiley; 1st edition (October 16, 2017)
	pubRegexp01 = regexp.MustCompile(`(^[A-Z][^;]+); ((\d+)(st|nd|rd|th))?([^(]+)\((\w+) (\d+), (\d+)\)`)
	// No Starch Press (November 5, 2020)
	pubRegexp02 = regexp.MustCompile(`(^[A-Z][^(;]+) \((\w+) (\d+), (\d+)\)`)
	// Packt Publishing; 3rd edition (17 May 2021)
	pubRegexp03 = regexp.MustCompile(`(^[A-Z][^;]+); ((\d+)(st|nd|rd|th))?([^(]+)\((\d+) (\w+)\.? (\d+)\)`)
	// No Starch Press (1 Oct. 2020)
	pubRegexp04 = regexp.MustCompile(`(^[A-Z][^(;]+) \((\d+) (\w+)\.? (\d+)\)`)
	// 522 pages
	lengthRegex = regexp.MustCompile(`(^\d+) pages`)
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
		// If there is no colon in the title string, try to extract the part in parentheses into subtitle string
		parenthesisStartIndex := strings.Index(titleString, "(")
		parenthesisEndIndex := strings.LastIndex(titleString, ")")
		if parenthesisStartIndex != -1 && parenthesisEndIndex == len(titleString)-1 {
			return titleString[:parenthesisStartIndex-1], titleString[parenthesisStartIndex:]
		}
		return titleString, ""
	}

	return titleString[:lastColonIndex], titleString[lastColonIndex+2:]
}

// ParsePublisherString parses a book publisher string and returns publisher-related meta information.
func ParsePublisherString(publisherString string) (BookPublishMeta, error) {
	var publisher string
	var edition = 1
	var day int
	var month string
	var year int

	// US format 1
	subMatch := pubRegexp01.FindStringSubmatch(publisherString)
	if len(subMatch) == 9 {
		publisher = subMatch[1]
		if subMatch[3] != "" {
			edition, _ = strconv.Atoi(subMatch[3])
		}
		month = subMatch[6]
		day, _ = strconv.Atoi(subMatch[7])
		year, _ = strconv.Atoi(subMatch[8])
	}

	// US format 2
	subMatch = pubRegexp02.FindStringSubmatch(publisherString)
	if len(subMatch) == 5 {
		publisher = subMatch[1]
		month = subMatch[2]
		day, _ = strconv.Atoi(subMatch[3])
		year, _ = strconv.Atoi(subMatch[4])
	}

	// EU format 1
	subMatch = pubRegexp03.FindStringSubmatch(publisherString)
	if len(subMatch) == 9 {
		publisher = subMatch[1]
		if subMatch[3] != "" {
			edition, _ = strconv.Atoi(subMatch[3])
		}
		day, _ = strconv.Atoi(subMatch[6])
		month = subMatch[7]
		year, _ = strconv.Atoi(subMatch[8])
	}

	// EU format 2
	subMatch = pubRegexp04.FindStringSubmatch(publisherString)
	if len(subMatch) == 5 {
		publisher = subMatch[1]
		day, _ = strconv.Atoi(subMatch[2])
		month = subMatch[3]
		year, _ = strconv.Atoi(subMatch[4])
	}

	if publisher == "" {
		return BookPublishMeta{}, fmt.Errorf("the publisher string '%s' can not be parsed", publisherString)
	}

	monthVal, err := time.Parse("Jan", month[:3])
	if err != nil {
		return BookPublishMeta{}, err
	}

	return BookPublishMeta{
		Publisher: strings.TrimSpace(publisher),
		Edition:   uint8(edition),
		PubDate:   time.Date(year, monthVal.Month(), day, 0, 0, 0, 0, time.UTC),
	}, nil
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

// ParseLengthString parses a book length string and returns its numeric value.
func ParseLengthString(lengthString string) uint16 {
	subMatch := lengthRegex.FindStringSubmatch(lengthString)
	if subMatch != nil {
		bookLength, _ := strconv.Atoi(subMatch[1])
		return uint16(bookLength)
	}

	return 0
}
