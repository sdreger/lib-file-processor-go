package scrapper

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

const (
	succeed = "\u2713"
	failed  = "\u2717"

	testBookID01           = "1234567890"
	testBookTitle          = "Test Title"
	testBookSubtitle       = "Test Subtitle"
	testBookDescription    = `<p><span>Test description</span></p>`
	testBookISBN10         = "1234567890"
	testBookISBN13         = 9781234567890
	testBookASIN           = "B08HG2JYS2"
	testBookPages          = 355
	testBookLanguage       = "English"
	testBookPublisher      = "Test Publisher"
	testBookEdition        = 4
	testBookFileName       = "Test.Publisher.Test.Title.4th.Edition.1234567890.Apr.2022.zip"
	testCoverFileName      = "1234567890.png"
	testCoverURL           = "https://cover.com/1.png"
	testBookAuthorName01   = "First Author"
	testBookAuthorName02   = "Second Author"
	testBookAuthorName03   = "Third Author"
	testBookCategoryName01 = "Computers & Technology"
	testBookCategoryName02 = "Programming"

	testBookResponse01 = "testdata/book_full.html"
)

var (
	handlersMap = map[string]string{
		testBookID01: testBookResponse01,
	}
	testBookPubDate    = time.Date(2022, 4, 6, 0, 0, 0, 0, time.UTC)
	testBookAuthors    = []string{testBookAuthorName01, testBookAuthorName02, testBookAuthorName03}
	testBookCategories = []string{testBookCategoryName01, testBookCategoryName02}
	testBookTags       []string
)

func TestGetBookDataISBN10(t *testing.T) {
	t.Log("Given the need to test book page scrapping.")
	server := testMockServer(t)
	defer server.Close()

	testBookPublisherURL := server.URL + "/" + testBookID01
	amazonScrapper, err := NewAmazonScrapper(server.URL+"/", log.Default())
	if err != nil {
		log.Fatalln(err)
	}
	bookMeta, err := amazonScrapper.GetBookData(testBookID01)
	if err != nil {
		log.Fatalln(err)
	}

	if bookMeta.Title != testBookTitle {
		t.Fatalf("\t\t%s\tShould get a %q book title: %q", failed, testBookTitle, bookMeta.Title)
	}
	if bookMeta.Subtitle != testBookSubtitle {
		t.Fatalf("\t\t%s\tShould get a %q book subtitle: %q", failed, testBookSubtitle, bookMeta.Subtitle)
	}
	if bookMeta.Description != testBookDescription {
		t.Fatalf("\t\t%s\tShould get a %q book description: %q", failed, testBookDescription, bookMeta.Description)
	}
	if bookMeta.ISBN10 != testBookISBN10 {
		t.Fatalf("\t\t%s\tShould get a %q book ISBN10: %q", failed, testBookISBN10, bookMeta.ISBN10)
	}
	if bookMeta.ISBN13 != testBookISBN13 {
		t.Fatalf("\t\t%s\tShould get a %d book ISBN13: %d", failed, testBookISBN13, bookMeta.ISBN13)
	}
	if bookMeta.ASIN != testBookASIN {
		t.Fatalf("\t\t%s\tShould get a %q book ASIN: %q", failed, testBookASIN, bookMeta.ASIN)
	}
	if bookMeta.Pages != testBookPages {
		t.Fatalf("\t\t%s\tShould get a %d book pages: %d", failed, testBookPages, bookMeta.Pages)
	}
	if bookMeta.Language != testBookLanguage {
		t.Fatalf("\t\t%s\tShould get a %q book language: %q", failed, testBookLanguage, bookMeta.Language)
	}
	if bookMeta.Publisher != testBookPublisher {
		t.Fatalf("\t\t%s\tShould get a %q book publisher: %q", failed, testBookPublisher, bookMeta.Publisher)
	}
	if bookMeta.PublisherURL != testBookPublisherURL {
		t.Fatalf("\t\t%s\tShould get a %q book publisher URL: %q", failed, testBookPublisherURL, bookMeta.PublisherURL)
	}
	if bookMeta.Edition != testBookEdition {
		t.Fatalf("\t\t%s\tShould get a %d book edition: %d", failed, testBookEdition, bookMeta.Edition)
	}
	if !bookMeta.PubDate.Equal(testBookPubDate) {
		t.Fatalf("\t\t%s\tShould get a %v book publish date: %v", failed, testBookPubDate, bookMeta.PubDate)
	}
	if !reflect.DeepEqual(bookMeta.Authors, testBookAuthors) {
		t.Fatalf("\t\t%s\tShould get %v book authors: %v", failed, testBookAuthors, bookMeta.Authors)
	}
	if !reflect.DeepEqual(bookMeta.Categories, testBookCategories) {
		t.Fatalf("\t\t%s\tShould get %v book categories: %v", failed, testBookCategories, bookMeta.Categories)
	}
	if !reflect.DeepEqual(bookMeta.Tags, testBookTags) {
		t.Fatalf("\t\t%s\tShould get %v book tags: %v", failed, testBookTags, bookMeta.Tags)
	}
	if bookMeta.BookFileName != testBookFileName {
		t.Fatalf("\t\t%s\tShould get a %q book file name: %q", failed, testBookFileName, bookMeta.BookFileName)
	}
	if bookMeta.CoverFileName != testCoverFileName {
		t.Fatalf("\t\t%s\tShould get a %q book cover file name: %q", failed, testCoverFileName, bookMeta.CoverFileName)
	}
	if bookMeta.CoverURL != testCoverURL {
		t.Fatalf("\t\t%s\tShould get a %q book cover URL: %q", failed, testCoverURL, bookMeta.CoverURL)
	}

	t.Logf("\t\t%s\tShould be able to scrape book data.", succeed)
}

func testMockServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	for path, responseFile := range handlersMap {
		mux.HandleFunc("/"+path, func(rw http.ResponseWriter, req *http.Request) {
			file, err := ioutil.ReadFile(responseFile)
			if err != nil {
				t.Fatal(err)
			}
			rw.Header().Set("Content-Type", "text/html")
			rw.WriteHeader(http.StatusOK)
			_, err = rw.Write(file)
			if err != nil {
				t.Fatal(err)
			}
		})
	}

	return server
}
