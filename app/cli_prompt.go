package app

import (
	"bufio"
	"fmt"
	"github.com/sdreger/lib-file-processor-go/domain/book"
	"log"
	"os"
	"strings"
)

func askForApproval(parsedData *book.ParsedData, existingData *book.StoredData, newLineDelimiter byte) bool {
	if existingData != nil {
		fmt.Println(strings.Repeat("~", 30), "Updating existing book", strings.Repeat("~", 30))
		fmt.Printf("%s   Title: %q => %q\n", equals(existingData.Title, parsedData.Title),
			existingData.Title, parsedData.Title)
		fmt.Printf("%s   Subtitle: %q => %q\n", equals(existingData.Subtitle, parsedData.Subtitle),
			existingData.Subtitle, parsedData.Subtitle)
		storedDescription := fmt.Sprintf("%s...%s", existingData.Description[:20],
			existingData.Description[len(existingData.Description)-20:])
		parsedDescription := fmt.Sprintf("%s...%s", parsedData.Description[:20],
			parsedData.Description[len(parsedData.Description)-20:])
		fmt.Printf("%s   Description: %q => %q\n", equals(storedDescription, parsedDescription),
			storedDescription, parsedDescription)
		fmt.Printf("%s   ISBN10: %q => %q\n", equals(existingData.ISBN10, parsedData.ISBN10),
			existingData.ISBN10, parsedData.ISBN10)
		fmt.Printf("%s   ISBN13: %d => %d\n", equals(existingData.ISBN13, parsedData.ISBN13),
			existingData.ISBN13, parsedData.ISBN13)
		fmt.Printf("%s   ASIN: %q => %q\n", equals(existingData.ASIN, parsedData.ASIN),
			existingData.ASIN, parsedData.ASIN)
		fmt.Printf("%s   Pages: %d => %d\n", equals(existingData.Pages, parsedData.Pages),
			existingData.Pages, parsedData.Pages)
		fmt.Printf("%s   Language: %q => %q\n", equals(existingData.Language, parsedData.Language),
			existingData.Language, parsedData.Language)
		fmt.Printf("%s   Publisher: %q => %q\n", equals(existingData.Publisher, parsedData.Publisher),
			existingData.Publisher, parsedData.Publisher)
		fmt.Printf("%s   PublisherURL: %q => %q\n", equals(existingData.PublisherURL, parsedData.PublisherURL),
			existingData.PublisherURL, parsedData.PublisherURL)
		fmt.Printf("%s   Edition: %d => %d\n", equals(existingData.Edition, parsedData.Edition),
			existingData.Edition, parsedData.Edition)
		existingPubDate := existingData.PubDate.Format("_2 Jan 2006")
		parsedPubDate := parsedData.PubDate.Format("_2 Jan 2006")
		fmt.Printf("%s   PubDate: %q => %q\n", equals(existingPubDate, parsedPubDate),
			existingPubDate, parsedPubDate)
		storedAuthors := strings.Join(existingData.Authors, ",")
		parsedAuthors := strings.Join(parsedData.Authors, ",")
		fmt.Printf("%s   Authors: %q => %q\n", equals(storedAuthors, parsedAuthors),
			storedAuthors, parsedAuthors)
		storedCategories := strings.Join(existingData.Categories, ",")
		parsedCategories := strings.Join(parsedData.Categories, ",")
		fmt.Printf("%s   Categories: %q => %q\n", equals(storedCategories, parsedCategories),
			storedCategories, parsedCategories)
		storedFormats := strings.Join(existingData.Formats, ",")
		parsedFormats := strings.Join(parsedData.Formats, ",")
		fmt.Printf("%s   Formats: %q => %q\n", equals(storedFormats, parsedFormats),
			storedFormats, parsedFormats)
		storedTags := strings.Join(existingData.Tags, ",")
		parsedTags := strings.Join(parsedData.Tags, ",")
		fmt.Printf("%s   Tags: %q => %q\n", equals(storedTags, parsedTags),
			storedTags, parsedTags)
		fmt.Printf("%s   BookFileName: %q => %q\n", equals(existingData.BookFileName, parsedData.BookFileName),
			existingData.BookFileName, parsedData.BookFileName)
		fmt.Printf("%s   BookFileSize: %d => %d\n", equals(existingData.BookFileSize, parsedData.BookFileSize),
			existingData.BookFileSize, parsedData.BookFileSize)
		fmt.Printf("%s   CoverFileName: %q => %q\n", equals(existingData.CoverFileName, parsedData.CoverFileName),
			existingData.CoverFileName, parsedData.CoverFileName)
		fmt.Println(strings.Repeat("~", 30), strings.Repeat("~", 30))
	} else {
		fmt.Println(strings.Repeat("+", 30), "Adding a new book", strings.Repeat("+", 30))
		fmt.Println(parsedData)
		fmt.Println(strings.Repeat("+", 50))
	}

	fmt.Print("Continue? (Press Return) | Edit Authors? (Type 'e' and press Return)")
	var twoChars [2]byte // a letter and \n
	_, err := os.Stdin.Read(twoChars[:])
	if err != nil {
		log.Fatal(err)
	}

	// Edit authors
	if string(twoChars[0]) == "e" || isEditAuthorsRequired(parsedData.Authors) {
		editAuthors(parsedData, newLineDelimiter)
	}

	return true
}

func isEditAuthorsRequired(authors []string) bool {
	for _, authorName := range authors {
		// If an author name contains just one word, then edit is required
		if len(strings.Split(authorName, " ")) < 2 {
			return true
		}
	}

	return false
}

func editAuthors(parsedData *book.ParsedData, newLineDelimiter byte) {
	fmt.Print("Enter author list (; separated): ")
	reader := bufio.NewReader(os.Stdin)
	authorsList, err := reader.ReadString(newLineDelimiter)
	if err != nil {
		log.Fatal(err)
	}
	authorsSplitString := strings.Split(authorsList, ";")
	for i, el := range authorsSplitString {
		authorsSplitString[i] = strings.TrimSpace(el)
	}
	if len(authorsSplitString) > 0 {
		parsedData.Authors = authorsSplitString
		fmt.Printf("New author list is: %v\n", strings.Join(parsedData.Authors, ","))
	}
}

func equals[T comparable](a, b T) string {
	if a == b {
		return "\u2713"
	}

	return "\u2717"
}
