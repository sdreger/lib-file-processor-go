package scrapper

import (
	"fmt"
	"time"
)

type BookMeta struct {
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

func (bm BookMeta) GetPrimaryId() string {
	if bm.ISBN10 != "" {
		return bm.ISBN10
	}
	if bm.ASIN != "" {
		return bm.ASIN
	}

	return fmt.Sprintf("%d", bm.ISBN13)
}
