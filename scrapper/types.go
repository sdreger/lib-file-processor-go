package scrapper

import "github.com/sdreger/lib-file-processor-go/domain/book"

//go:generate mockgen -destination=./book_data_scrapper_mock.go -package=scrapper github.com/sdreger/lib-file-processor-go/scrapper BookDataScrapper
type BookDataScrapper interface {
	GetBookData(bookID string) (book.ParsedData, error)
	Close() error
}
