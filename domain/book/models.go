package book

import (
	"database/sql"
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

func (pd ParsedData) GetPrimaryId() string {
	if pd.ISBN10 != "" {
		return pd.ISBN10
	}
	if pd.ASIN != "" {
		return pd.ASIN
	}

	return fmt.Sprintf("%d", pd.ISBN13)
}

func (pd ParsedData) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintln("book.ParsedData: {"))
	b.WriteString(fmt.Sprintf("\tTitle: %q\n", pd.Title))
	b.WriteString(fmt.Sprintf("\tSubtitle: %q\n", pd.Subtitle))
	b.WriteString(fmt.Sprintf("\tDescription: \"%s...%s\"\n", pd.Description[:40], pd.Description[len(pd.Description)-40:]))
	b.WriteString(fmt.Sprintf("\tISBN10: %q\n", pd.ISBN10))
	b.WriteString(fmt.Sprintf("\tISBN13: %d\n", pd.ISBN13))
	b.WriteString(fmt.Sprintf("\tASIN: %q\n", pd.ASIN))
	b.WriteString(fmt.Sprintf("\tPages: %d\n", pd.Pages))
	b.WriteString(fmt.Sprintf("\tLanguage: %q\n", pd.Language))
	b.WriteString(fmt.Sprintf("\tPublisher: %q\n", pd.Publisher))
	b.WriteString(fmt.Sprintf("\tPublisherURL: %q\n", pd.PublisherURL))
	b.WriteString(fmt.Sprintf("\tEdition: %d\n", pd.Edition))
	b.WriteString(fmt.Sprintf("\tPubDate: %q\n", pd.PubDate.Format("_2 Jan 2006")))
	b.WriteString(fmt.Sprintf("\tAuthors: %q\n", strings.Join(pd.Authors, ",")))
	b.WriteString(fmt.Sprintf("\tCategories: %q\n", strings.Join(pd.Categories, ",")))
	b.WriteString(fmt.Sprintf("\tTags: %q\n", strings.Join(pd.Tags, ",")))
	b.WriteString(fmt.Sprintf("\tFormats: %q\n", strings.Join(pd.Formats, ",")))
	b.WriteString(fmt.Sprintf("\tBookFileName: %q\n", pd.BookFileName))
	b.WriteString(fmt.Sprintf("\tBookFileSize: %d\n", pd.BookFileSize))
	b.WriteString(fmt.Sprintf("\tCoverFileName: %q\n", pd.CoverFileName))
	b.WriteString(fmt.Sprintf("\tCoverURL: %q\n", pd.CoverURL))
	b.WriteString(fmt.Sprintln("}"))

	return b.String()
}

func (pd ParsedData) GetBookFileNameWithoutExtension() string {
	return pd.BookFileName[:strings.LastIndex(pd.BookFileName, ".")]
}

type relationKeys struct {
	publisherID int64
	languageID  int64
	authorIDs   []int64
	categoryIDs []int64
	fileTypeIDs []int64
	tagIDs      []int64
}

type SearchRequest struct {
	Title   string
	Edition uint8
	ISBN10  string
	ISBN13  int64
	ASIN    string
}

type StoredData struct {
	ID            int64
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
	BookFileName  string
	BookFileSize  int64
	CoverFileName string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Authors       []string
	Categories    []string
	Tags          []string
	Formats       []string
}

func (sd StoredData) IsEmpty() bool {
	if sd.ID == 0 {
		return true
	}

	return false
}

type dotProductRow struct {
	ID            int64
	Title         string
	Subtitle      sql.NullString
	Description   string
	ISBN10        sql.NullString
	ISBN13        sql.NullInt64
	ASIN          sql.NullString
	Pages         uint16
	Language      string
	Publisher     string
	PublisherURL  string
	Edition       uint8
	PubDate       time.Time
	BookFileName  string
	BookFileSize  int64
	CoverFileName string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	AuthorName    sql.NullString
	CategoryName  sql.NullString
	FileTypeName  sql.NullString
	TagName       sql.NullString
}
