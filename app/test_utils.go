package app

import (
	"github.com/sdreger/lib-file-processor-go/domain/book"
	"github.com/sdreger/lib-file-processor-go/filestore"
	"time"
)

const (
	succeed = "\u2713"
	failed  = "\u2717"

	testBookID            = "1234567890"
	testBookIDInt         = int64(1234567890)
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
	testBookFileName      = "Test_book_name.1573273281.zip"
	testBookFileSize      = 5000
	testBookCoverFileName = "1573273281.png"
	testBookCoverURL      = "https://cover.com/1.png"
	testBookAuthorName    = "Test author name"
	testBookCategoryName  = "Test category name"
	testBookFileTypeName  = "Test filetype name"
	testBookTagName       = "Test tag name"

	testBookEtag  = "book-1234567890"
	testCoverEtag = "cover-1234567890"
)

var (
	testPublishDate, _ = time.Parse(time.RFC3339, "2020-02-20T05:15:45Z")
	testCreateDate, _  = time.Parse(time.RFC3339, "2022-07-30T12:20:00Z")

	testBookArchivePath = "/out/book/test_book.zip"
	testBookFormats     = []string{"pdf", "epub"}
	testCoverFilePath   = "/out_cover/test_book.png"
)

func getTestParsedData() book.ParsedData {
	return book.ParsedData{
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
		Authors:       []string{testBookAuthorName},
		Categories:    []string{testBookCategoryName},
		Tags:          []string{testBookTagName},
		Formats:       nil,
		BookFileName:  testBookFileName,
		BookFileSize:  0,
		CoverFileName: testBookCoverFileName,
		CoverURL:      testBookCoverURL,
	}
}

func getTestStoredData() book.StoredData {
	return book.StoredData{
		ID:            testBookIDInt,
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

func getTestTempFilesData() filestore.TempFilesData {
	return filestore.TempFilesData{
		BookArchivePath: testBookArchivePath,
		BookFormats:     testBookFormats,
		BookSize:        testBookFileSize,
		CoverFilePath:   testCoverFilePath,
	}
}
