package filestore

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type DownloadService struct{}

func NewDownloadService() DownloadService {
	return DownloadService{}
}

// DownloadCoverFile downloads a book cover file from a remote URL,
// and stores it to the provided folder. Returns a stored filepath.
func (ds DownloadService) DownloadCoverFile(coverURL, coverOutputFolder, coverFileName string) (string, error) {
	bookCoverOutputPath := filepath.Join(coverOutputFolder, coverFileName)
	cover, err := os.Create(bookCoverOutputPath)
	if err != nil {
		return "", err
	}
	defer closeResource(cover)

	response, err := http.Get(coverURL)
	if err != nil {
		return "", err
	}
	defer closeResponseBody(response.Body)

	writtenBytes, err := io.Copy(cover, response.Body)
	if err != nil {
		return "", err
	}

	if err = cover.Sync(); err != nil {
		return "", err
	}
	log.Printf("[INFO] - Written %d bytes of book cover", writtenBytes)

	return bookCoverOutputPath, err
}

func closeResponseBody(f io.Closer) {
	err := f.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}
