package book

import (
	"database/sql"
	"time"
)

const (
	succeed = "\u2713"
	failed  = "\u2717"

	testBookId            = 100
	testBookTitle         = "Test title"
	testBookSubtitle      = "Test subtitle"
	testBookDescription   = "Test description"
	testBookISBN10        = "1573273281"
	testBookISBN13        = 9781573273281
	testBookASIN          = "B08HG2JYS2"
	testBookPages         = 355
	testBookLanguage      = "Test language"
	testBookPublisher     = "Test publisher"
	testBookPublisherURL  = "https://test.pub/1573273281"
	testBookEdition       = 3
	testBookFileName      = "Test book name"
	testBookFileSize      = 5000
	testBookCoverFileName = "Test book cover name"
	testBookAuthorName    = "Test author name"
	testBookCategoryName  = "Test category name"
	testBookFileTypeName  = "Test filetype name"
	testBookTagName       = "Test tag name"
)

var (
	testPublishDate, _ = time.Parse(time.RFC3339, "2020-02-20T05:15:45Z")
	testCreateDate, _  = time.Parse(time.RFC3339, "2022-07-30T12:20:00Z")
)

func getTestProductRow() dotProductRow {
	return dotProductRow{
		ID:            testBookId,
		Title:         testBookTitle,
		Subtitle:      sql.NullString{String: testBookSubtitle, Valid: true},
		Description:   testBookDescription,
		ISBN10:        sql.NullString{String: testBookISBN10, Valid: true},
		ISBN13:        sql.NullInt64{Int64: testBookISBN13, Valid: true},
		ASIN:          sql.NullString{String: testBookASIN, Valid: true},
		Pages:         testBookPages,
		Language:      testBookLanguage,
		Publisher:     testBookPublisher,
		PublisherURL:  testBookPublisherURL,
		Edition:       testBookEdition,
		PubDate:       testPublishDate,
		BookFileName:  testBookFileName,
		BookFileSize:  testBookFileSize,
		CoverFileName: testBookCoverFileName,
		CreatedAt:     testCreateDate,
		UpdatedAt:     testCreateDate,
		AuthorName:    sql.NullString{String: testBookAuthorName, Valid: true},
		CategoryName:  sql.NullString{String: testBookCategoryName, Valid: true},
		FileTypeName:  sql.NullString{String: testBookFileTypeName, Valid: true},
		TagName:       sql.NullString{String: testBookTagName, Valid: true},
	}
}

func getTestStoredData() StoredData {
	return StoredData{
		ID:            testBookId,
		Title:         testBookTitle,
		Subtitle:      testBookSubtitle,
		Description:   testBookDescription,
		ISBN10:        testBookISBN10,
		ISBN13:        testBookISBN13,
		ASIN:          testBookASIN,
		Pages:         testBookPages,
		Language:      testBookLanguage,
		Publisher:     testBookPublisher,
		PublisherURL:  testBookPublisherURL,
		Edition:       testBookEdition,
		PubDate:       testPublishDate,
		BookFileName:  testBookFileName,
		BookFileSize:  testBookFileSize,
		CoverFileName: testBookCoverFileName,
		CreatedAt:     testCreateDate,
		UpdatedAt:     testCreateDate,
		Authors:       []string{testBookAuthorName},
		Categories:    []string{testBookCategoryName},
		Tags:          []string{testBookTagName},
		Formats:       []string{testBookFileTypeName},
	}
}
