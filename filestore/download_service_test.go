package filestore

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	coverFileURL = "/cover/1.png"
)

func TestDownloadCoverFile(t *testing.T) {
	coverFileName := "cover-1.png"

	t.Log("Given the need to test cover file download.")
	server := testMockServer(t)
	defer server.Close()

	tempDir, err := os.MkdirTemp("", "cover-dir-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	filePath, err := NewDownloadService().DownloadCoverFile(server.URL+coverFileURL, tempDir, coverFileName)
	if err != nil || filePath == "" {
		t.Fatalf("\t\t%s\tShould be able to download a cover file: %v", failed, err)
	}

	stat, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to get a cover file stat: %v", failed, err)
	}
	if stat.Name() != coverFileName {
		t.Fatalf("\t\t%s\tShould get a %q cover file name: %q", failed, coverFileName, stat)
	}
	if stat.Size() == 0 {
		t.Fatalf("\t\t%s\tCover file size should not be 0", failed)
	}

	dirEntries, err := os.ReadDir(tempDir)
	if err != nil || len(dirEntries) != 1 {
		t.Fatalf("\t\t%s\tShould be able to get directory file list: %v", failed, err)
	}
	if dirEntries[0].Name() != coverFileName {
		t.Fatalf("\t\t%s\tShould get a %q cover file name: %q", failed, coverFileName, dirEntries[0].Name())
	}

	t.Logf("\t\t%s\tShould successfully download a book cover", succeed)
}

func testMockServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc(coverFileURL, func(rw http.ResponseWriter, req *http.Request) {
		file, err := os.Open("testdata/1.png")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		rw.Header().Set("Content-Type", "image/png")
		rw.WriteHeader(http.StatusOK)
		_, err = io.Copy(rw, file)
		if err != nil {
			t.Fatal(err)
		}
	})

	return server
}
