package scrapper

import "time"

type BookMeta struct {
	Title         string
	Subtitle      string
	Description   string
	ISBN10        string
	ISBN13        string
	ASIN          string
	Pages         int
	Language      string
	PublisherURL  string
	Publisher     string
	Edition       uint8
	PubDate       time.Time
	Authors       []string
	Categories    []string
	Tags          []string
	Formats       []string
	BookFileSize  int
	BookFileName  string
	CoverFileName string
}
