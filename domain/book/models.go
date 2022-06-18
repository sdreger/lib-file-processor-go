package book

import (
	"fmt"
	"strings"
	"time"
)

type ParsedData struct {
	Title         string
	Subtitle      string
	Description   string
	ISBN10        string
	ISBN13        int64
	ASIN          string
	Pages         uint16
	Language      string
	Publisher     string
	PublisherURL  string
	Edition       uint8
	PubDate       time.Time
	Authors       []string
	Categories    []string
	Tags          []string
	Formats       []string
	BookFileName  string
	BookFileSize  int64
	CoverFileName string
	CoverURL      string
}

func (bm ParsedData) GetPrimaryId() string {
	if bm.ISBN10 != "" {
		return bm.ISBN10
	}
	if bm.ASIN != "" {
		return bm.ASIN
	}

	return fmt.Sprintf("%d", bm.ISBN13)
}

func (bm ParsedData) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintln("book.ParsedData: {"))
	b.WriteString(fmt.Sprintf("\tTitle: %q\n", bm.Title))
	b.WriteString(fmt.Sprintf("\tSubtitle: %q\n", bm.Subtitle))
	b.WriteString(fmt.Sprintf("\tDescription: \"%s...%s\"\n", bm.Description[:40], bm.Description[len(bm.Description)-40:]))
	b.WriteString(fmt.Sprintf("\tISBN10: %q\n", bm.ISBN10))
	b.WriteString(fmt.Sprintf("\tISBN13: %d\n", bm.ISBN13))
	b.WriteString(fmt.Sprintf("\tASIN: %q\n", bm.ASIN))
	b.WriteString(fmt.Sprintf("\tPages: %d\n", bm.Pages))
	b.WriteString(fmt.Sprintf("\tLanguage: %q\n", bm.Language))
	b.WriteString(fmt.Sprintf("\tPublisher: %q\n", bm.Publisher))
	b.WriteString(fmt.Sprintf("\tPublisherURL: %q\n", bm.PublisherURL))
	b.WriteString(fmt.Sprintf("\tEdition: %d\n", bm.Edition))
	b.WriteString(fmt.Sprintf("\tPubDate: %q\n", bm.PubDate.Format("_2 Jan 2006")))
	b.WriteString(fmt.Sprintf("\tAuthors: %q\n", strings.Join(bm.Authors, ",")))
	b.WriteString(fmt.Sprintf("\tCategories: %q\n", strings.Join(bm.Categories, ",")))
	b.WriteString(fmt.Sprintf("\tTags: %q\n", strings.Join(bm.Tags, ",")))
	b.WriteString(fmt.Sprintf("\tFormats: %q\n", strings.Join(bm.Formats, ",")))
	b.WriteString(fmt.Sprintf("\tBookFileName: %q\n", bm.BookFileName))
	b.WriteString(fmt.Sprintf("\tBookFileSize: %d\n", bm.BookFileSize))
	b.WriteString(fmt.Sprintf("\tCoverFileName: %q\n", bm.CoverFileName))
	b.WriteString(fmt.Sprintf("\tCoverURL: %q\n", bm.CoverURL))
	b.WriteString(fmt.Sprintln("}"))

	return b.String()
}
