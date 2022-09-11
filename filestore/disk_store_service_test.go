package filestore

import (
	"github.com/golang/mock/gomock"
	"github.com/sdreger/lib-file-processor-go/domain/book"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const (
	testCoverURL     = "https://cover.com/12345.png"
	testPdfBookName  = "1.pdf"
	testEpubBookName = "1.epub"
	testCoverName    = "1.png"
	testArchiveName  = "1.zip"
	tempInputDir     = "/tmp/book-input"
)

func TestDiskFileStore_PrepareBookFiles(t *testing.T) {

	t.Log("Given the need to test book files preparing.")
	mockCtrl := gomock.NewController(t)
	mockCoverDownloader := NewMockCoverDownloader(mockCtrl)
	mockBookCompressor := NewMockBookCompressor(mockCtrl)
	diskStore := NewDiskStoreService(mockBookCompressor, mockCoverDownloader, log.Default())
	defer mockCtrl.Finish()

	tempOutputDir, err := os.MkdirTemp("", "output-dir-*")
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create output folder: %v", failed, err)
	}
	defer os.RemoveAll(tempOutputDir)

	testArchivePath := filepath.Join(tempOutputDir, testArchiveName)
	bookArchive, err := os.Create(testArchivePath)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create a book archive file: %v", failed, err)
	}
	if _, err = bookArchive.WriteString("File content"); err != nil {
		t.Fatalf("\t\t%s\tShould be able to write data to the book archive: %v", failed, err)
	}
	bookArchive.Close()

	parsedData := book.ParsedData{
		BookFileName:  testArchiveName,
		CoverFileName: testCoverName,
		CoverURL:      testCoverURL,
	}

	testCoverPath := filepath.Join(tempOutputDir, testCoverName)
	mockCoverDownloader.EXPECT().DownloadCoverFile(testCoverURL, tempOutputDir, testCoverName).
		Return(testCoverPath, nil).Times(1)

	namesInArchive := []string{testPdfBookName, testEpubBookName}
	mockBookCompressor.EXPECT().CompressBookFiles(tempInputDir, tempOutputDir, parsedData.BookFileName).
		Return(testArchivePath, namesInArchive, nil).Times(1)

	tempFilesData, err := diskStore.PrepareBookFiles(parsedData, tempInputDir, tempOutputDir)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to prepare book files: %v", failed, err)
	}
	if tempFilesData.CoverFilePath != testCoverPath {
		t.Fatalf("\t\t%s\tShould get a %q cover path: %q", failed, testCoverPath, tempFilesData.CoverFilePath)
	}
	if tempFilesData.BookArchivePath != testArchivePath {
		t.Fatalf("\t\t%s\tShould get a %q book archive path: %q", failed,
			testArchivePath, tempFilesData.BookArchivePath)
	}
	if len(tempFilesData.BookFormats) != len(namesInArchive) {
		t.Fatalf("\t\t%s\tShould get %d file formats: %d", failed, len(namesInArchive), len(tempFilesData.BookFormats))
	}
	if tempFilesData.BookSize == 0 {
		t.Fatalf("\t\t%s\tBook archive size should not be 0", failed)
	}

	t.Logf("\t\t%s\tShould successfully prepare book files", succeed)
}

func TestDiskStore_StoreBookArchive(t *testing.T) {
	t.Log("Given the need to test book archive storing.")
	t.Run("The output folder does not exist", testStoreBookArchiveOutputFolderDoesNotExist)
	t.Run("The output folder exists", testStoreBookArchiveOutputFolderExists)
}

func testStoreBookArchiveOutputFolderDoesNotExist(t *testing.T) {
	outputSubDir := "sub"
	diskStore := NewDiskStoreService(nil, nil, log.Default())

	// Create book input folder with book files inside
	bookInputDir := createBookInputFolder(t)
	defer os.RemoveAll(bookInputDir)

	// Create temp folder with a book archive inside
	bookTempDir := createBookTempFolder(t)
	defer os.RemoveAll(bookTempDir)

	// Create book output folder to place book archive file into
	bookOutputDir := createBookOutputDir(t)
	defer os.RemoveAll(bookOutputDir)

	bookInputPath := filepath.Join(bookTempDir, testArchiveName)
	bookOutputPath := filepath.Join(bookOutputDir, outputSubDir, testArchiveName)
	err := diskStore.StoreBookArchive(bookInputDir, bookInputPath, bookOutputPath)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to store a book archive file: %v", failed, err)
	}
	assertBookStoreFoldersContent(t, bookInputDir, bookTempDir, filepath.Join(bookOutputDir, outputSubDir))

	t.Logf("\t\t%s\tShould successfully store a book archive into non-existing folder", succeed)
}

func testStoreBookArchiveOutputFolderExists(t *testing.T) {

	diskStore := NewDiskStoreService(nil, nil, log.Default())

	// Create book input folder with book files inside
	bookInputDir := createBookInputFolder(t)
	defer os.RemoveAll(bookInputDir)

	// Create temp folder with a book archive inside
	bookTempDir := createBookTempFolder(t)
	defer os.RemoveAll(bookTempDir)

	// Create book output folder to place book archive file into
	bookOutputDir := createBookOutputDir(t)
	defer os.RemoveAll(bookOutputDir)

	bookInputPath := filepath.Join(bookTempDir, testArchiveName)
	bookOutputPath := filepath.Join(bookOutputDir, testArchiveName)
	err := diskStore.StoreBookArchive(bookInputDir, bookInputPath, bookOutputPath)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to store a book archive file: %v", failed, err)
	}
	assertBookStoreFoldersContent(t, bookInputDir, bookTempDir, bookOutputDir)

	t.Logf("\t\t%s\tShould successfully store a book archive into existing folder", succeed)
}

func createBookInputFolder(t *testing.T) string {
	bookInputDir, err := os.MkdirTemp("", "book-input-dir-*")
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create a book input folder: %v", failed, err)
	}

	bookFilePdf, err := os.Create(filepath.Join(bookInputDir, testPdfBookName))
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create a PDF book file: %v", failed, err)
	}
	bookFilePdf.Close()
	bookFileEpub, err := os.Create(filepath.Join(bookInputDir, testEpubBookName))
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create a EPUB book file: %v", failed, err)
	}
	bookFileEpub.Close()

	return bookInputDir
}

func createBookTempFolder(t *testing.T) string {
	bookTempDir, err := os.MkdirTemp("", "book-temp-dir-*")
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create a temp folder: %v", failed, err)
	}
	bookArchive, err := os.Create(filepath.Join(bookTempDir, testArchiveName))
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create a book archive file: %v", failed, err)
	}
	bookArchive.Close()

	return bookTempDir
}

func createBookOutputDir(t *testing.T) string {
	bookOutputDir, err := os.MkdirTemp("", "book-output-dir-*")
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create a book output folder: %v", failed, err)
	}

	return bookOutputDir
}

func assertBookStoreFoldersContent(t *testing.T, bookInputDir, bookTempDir, bookOutputDir string) {
	// Check if book archive is present in the output folder
	outputDirEntries, err := os.ReadDir(bookOutputDir)
	if err != nil || len(outputDirEntries) == 0 {
		t.Fatalf("\t\t%s\tShould be able to get output folder contents: %v", failed, err)
	}
	if outputDirEntries[0].Name() != testArchiveName {
		t.Fatalf("\t\t%s\tThe output folder should contain a book archive: %v", failed, err)
	}

	// Check if book files are removed from the input folder
	inputDirEntries, err := os.ReadDir(bookInputDir)
	if err != nil || len(inputDirEntries) != 0 {
		t.Fatalf("\t\t%s\tThe input folder should be empty: %v", failed, err)
	}

	// Check if book archive is removed from the temp folder
	tempDirEntries, err := os.ReadDir(bookTempDir)
	if err != nil || len(tempDirEntries) != 0 {
		t.Fatalf("\t\t%s\tThe temp folder should be empty: %v", failed, err)
	}
}

func TestDiskStore_StoreCoverFile(t *testing.T) {
	t.Log("Given the need to test book cover storing.")
	t.Run("The output folder does not exist", testStoreBookCoverOutputFolderDoesNotExist)
	t.Run("The output folder exists", testStoreBookCoverOutputFolderExists)
	//diskStore := NewDiskStoreService(nil, nil)
}

func testStoreBookCoverOutputFolderDoesNotExist(t *testing.T) {

	diskStore := NewDiskStoreService(nil, nil, log.Default())

	// Create temp folder with an image file inside
	coverTempDir := createCoverTempFolder(t)
	defer os.RemoveAll(coverTempDir)

	// Create book output folder to place cover file into
	coverOutputDir := createCoverOutputDir(t)
	defer os.RemoveAll(coverOutputDir)

	coverInputPath := filepath.Join(coverTempDir, testCoverName)
	coverOutputPath := filepath.Join(coverOutputDir, testCoverName)
	err := diskStore.StoreCoverFile(coverInputPath, coverOutputPath)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to store a cover file: %v", failed, err)
	}
	assertCoverStoreFoldersContent(t, coverTempDir, coverOutputDir)

	t.Logf("\t\t%s\tShould successfully store a book cover into existing folder", succeed)
}

func testStoreBookCoverOutputFolderExists(t *testing.T) {
	outputSubDir := "sub"
	diskStore := NewDiskStoreService(nil, nil, log.Default())

	// Create temp folder with an image file inside
	coverTempDir := createCoverTempFolder(t)
	defer os.RemoveAll(coverTempDir)

	// Create book output folder to place cover file into
	coverOutputDir := createCoverOutputDir(t)
	defer os.RemoveAll(coverOutputDir)

	coverInputPath := filepath.Join(coverTempDir, testCoverName)
	coverOutputPath := filepath.Join(coverOutputDir, outputSubDir, testCoverName)
	err := diskStore.StoreCoverFile(coverInputPath, coverOutputPath)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to store a cover file: %v", failed, err)
	}
	assertCoverStoreFoldersContent(t, coverTempDir, filepath.Join(coverOutputDir, outputSubDir))

	t.Logf("\t\t%s\tShould successfully store a book cover into non-existing folder", succeed)
}

func createCoverTempFolder(t *testing.T) string {
	coverTempDir, err := os.MkdirTemp("", "book-temp-dir-*")
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create a temp folder: %v", failed, err)
	}
	coverFile, err := os.Create(filepath.Join(coverTempDir, testCoverName))
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create a book cover file: %v", failed, err)
	}
	coverFile.Close()

	return coverTempDir
}

func createCoverOutputDir(t *testing.T) string {
	coverOutputDir, err := os.MkdirTemp("", "cover-output-dir-*")
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create a cover output folder: %v", failed, err)
	}

	return coverOutputDir
}

func assertCoverStoreFoldersContent(t *testing.T, coverTempDir, coverOutputDir string) {
	// Check if book cover is present in the output folder
	outputDirEntries, err := os.ReadDir(coverOutputDir)
	if err != nil || len(outputDirEntries) == 0 {
		t.Fatalf("\t\t%s\tShould be able to get output folder contents: %v", failed, err)
	}
	if outputDirEntries[0].Name() != testCoverName {
		t.Fatalf("\t\t%s\tThe output folder should contain a book cover: %v", failed, err)
	}

	// Check if book cover is removed from the temp folder
	tempDirEntries, err := os.ReadDir(coverTempDir)
	if err != nil || len(tempDirEntries) != 0 {
		t.Fatalf("\t\t%s\tThe temp folder should be empty: %v", failed, err)
	}
}

func TestDiskStore_IsFolderEmpty(t *testing.T) {

	t.Log("Given the need to test book files preparing.")
	diskStore := NewDiskStoreService(nil, nil, log.Default())

	tempDir, err := os.MkdirTemp("", "temp-dir-to-check-content-*")
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create a temp folder: %v", failed, err)
	}
	defer os.RemoveAll(tempDir)

	isEmpty, err := diskStore.IsFolderEmpty(tempDir)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to check if folder is empty: %v", failed, err)
	}
	if !isEmpty {
		t.Fatalf("\t\t%s\tThe folder should be empty", failed)
	}
	t.Logf("\t\t%s\tShould successfully check empty folder", succeed)

	bookArchive, err := os.Create(filepath.Join(tempDir, testArchiveName))
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create a book archive file: %v", failed, err)
	}
	bookArchive.Close()

	isEmpty, err = diskStore.IsFolderEmpty(tempDir)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to check if folder is empty: %v", failed, err)
	}
	if isEmpty {
		t.Fatalf("\t\t%s\tThe folder should not be empty", failed)
	}

	t.Logf("\t\t%s\tShould successfully check non-empty folder", succeed)
}
