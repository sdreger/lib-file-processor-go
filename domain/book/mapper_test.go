package book

import (
	"reflect"
	"testing"
)

func TestMapToStoredData(t *testing.T) {
	t.Log("Given the need to test store data mapping.")
	input := getTestProductRow()
	output := StoredData{}
	mapToStoredData(input, &output)

	if output.ID != input.ID {
		t.Errorf("\t\t%s\tShould get a %d mapped value: %d", failed, input.ID, output.ID)
	}
	if output.Title != input.Title {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.Title, output.Title)
	}
	if output.Subtitle != input.Subtitle.String {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.Subtitle.String, output.Subtitle)
	}
	if output.Description != input.Description {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.Description, output.Description)
	}
	if output.ISBN10 != input.ISBN10.String {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.ISBN10.String, output.ISBN10)
	}
	if output.ISBN13 != input.ISBN13.Int64 {
		t.Errorf("\t\t%s\tShould get a %d mapped value: %d", failed, input.ISBN13.Int64, output.ISBN13)
	}
	if output.ASIN != input.ASIN.String {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.ASIN.String, output.ASIN)
	}
	if output.Pages != input.Pages {
		t.Errorf("\t\t%s\tShould get a %d mapped value: %d", failed, input.Pages, output.Pages)
	}
	if output.Language != input.Language {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.Language, output.Language)
	}
	if output.Publisher != input.Publisher {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.Publisher, output.Publisher)
	}
	if output.PublisherURL != input.PublisherURL {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.PublisherURL, output.PublisherURL)
	}
	if output.Edition != input.Edition {
		t.Errorf("\t\t%s\tShould get a %d mapped value: %d", failed, input.Edition, output.Edition)
	}
	if output.PubDate != input.PubDate {
		t.Errorf("\t\t%s\tShould get a %v mapped value: %v", failed, input.PubDate, output.PubDate)
	}
	if output.BookFileName != input.BookFileName {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.BookFileName, output.BookFileName)
	}
	if output.BookFileSize != input.BookFileSize {
		t.Errorf("\t\t%s\tShould get a %d mapped value: %d", failed, input.BookFileSize, output.BookFileSize)
	}
	if output.CoverFileName != input.CoverFileName {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.CoverFileName, output.CoverFileName)
	}
	if output.CreatedAt != input.CreatedAt {
		t.Errorf("\t\t%s\tShould get a %v mapped value: %v", failed, input.CreatedAt, output.CreatedAt)
	}
	if output.UpdatedAt != input.UpdatedAt {
		t.Errorf("\t\t%s\tShould get a %v mapped value: %v", failed, input.UpdatedAt, output.UpdatedAt)
	}
	if output.Authors[0] != input.AuthorName.String {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.AuthorName.String, output.Authors[0])
	}
	if output.Categories[0] != input.CategoryName.String {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.CategoryName.String, output.Categories[0])
	}
	if output.Formats[0] != input.FileTypeName.String {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.FileTypeName.String, output.Authors[0])
	}
	if output.Tags[0] != input.TagName.String {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.TagName.String, output.Tags[0])
	}

	t.Logf("\t\t%s\tShould be able to map dot product row to StoredData", succeed)
}

func TestMapToParsedDate(t *testing.T) {
	t.Log("Given the need to test parsed data mapping.")
	input := getTestStoredData()
	output := ParsedData{}
	mapToParsedDate(&input, &output)

	if output.Title != input.Title {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.Title, output.Title)
	}
	if output.Subtitle != input.Subtitle {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.Subtitle, output.Subtitle)
	}
	if output.Description != input.Description {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.Description, output.Description)
	}
	if output.ISBN10 != input.ISBN10 {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.ISBN10, output.ISBN10)
	}
	if output.ISBN13 != input.ISBN13 {
		t.Errorf("\t\t%s\tShould get a %d mapped value: %d", failed, input.ISBN13, output.ISBN13)
	}
	if output.ASIN != input.ASIN {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.ASIN, output.ASIN)
	}
	if output.Pages != input.Pages {
		t.Errorf("\t\t%s\tShould get a %d mapped value: %d", failed, input.Pages, output.Pages)
	}
	if output.PublisherURL != input.PublisherURL {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.PublisherURL, output.PublisherURL)
	}
	if output.Edition != input.Edition {
		t.Errorf("\t\t%s\tShould get a %d mapped value: %d", failed, input.Edition, output.Edition)
	}
	if output.PubDate != input.PubDate {
		t.Errorf("\t\t%s\tShould get a %v mapped value: %v", failed, input.PubDate, output.PubDate)
	}
	if output.BookFileName != input.BookFileName {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.BookFileName, output.BookFileName)
	}
	if output.BookFileSize != input.BookFileSize {
		t.Errorf("\t\t%s\tShould get a %d mapped value: %d", failed, input.BookFileSize, output.BookFileSize)
	}
	if output.CoverFileName != input.CoverFileName {
		t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, input.CoverFileName, output.CoverFileName)
	}

	t.Logf("\t\t%s\tShould be able to map StoredData to ParsedData", succeed)
}

func TestDeduplicateMappedData(t *testing.T) {
	t.Log("Given the need to test mapped data deduplication.")
	item01 := "AAA"
	item02 := "BBB"
	item03 := "CCC"
	input := []string{item01, item02, item01, item02, item03}
	output := deduplicateMappedData(input)
	expected := []string{item01, item02, item03}
	if len(output) != 3 {
		t.Errorf("\t\t%s\tShould get %s deduplicated slice, got: %s", failed, expected, output)
	}
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("\t\t%s\tShould get %s deduplicated slice, got: %s", failed, expected, output)
	}

	t.Logf("\t\t%s\tShould be able to deduplicate mapped data slice", succeed)
}
