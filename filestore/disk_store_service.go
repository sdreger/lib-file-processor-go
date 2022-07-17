package filestore

import (
	"fmt"
	"github.com/sdreger/lib-file-processor-go/domain/book"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type DiskStoreService struct {
	bookCompressor  BookCompressor
	coverDownloader CoverDownloader
}

func NewDiskStoreService(bookCompressor BookCompressor, coverDownloader CoverDownloader) DiskStoreService {
	return DiskStoreService{
		bookCompressor:  bookCompressor,
		coverDownloader: coverDownloader,
	}
}

// PrepareBookFiles downloads a book cover, compress book files,
// and put both of them to the output folder.
func (ds DiskStoreService) PrepareBookFiles(bookMeta book.ParsedData, bookInputFolder, outputFolder string) (TempFilesData, error) {

	coverFilePath, err := ds.coverDownloader.DownloadCoverFile(bookMeta.CoverURL, outputFolder, bookMeta.CoverFileName)
	if err != nil {
		return TempFilesData{}, fmt.Errorf("can not store a book cover: %w", err)
	}

	archiveFilePath, namesInArchive, err :=
		ds.bookCompressor.CompressBookFiles(bookInputFolder, outputFolder, bookMeta.BookFileName)
	if err != nil {
		return TempFilesData{}, fmt.Errorf("can not compress book files: %w", err)
	}

	size, err := getFileSize(filepath.Join(outputFolder, bookMeta.BookFileName))
	if err != nil {
		return TempFilesData{}, err
	}

	return TempFilesData{
		BookArchivePath: archiveFilePath,
		BookFormats:     getFilesTypes(namesInArchive),
		BookSize:        size,
		CoverFilePath:   coverFilePath,
	}, nil
}

// StoreBookArchive moves a book archive file from the temp folder
// to the output folder. Creates the output folder if not exist.
// After that removes all processed book files from the input folder.
func (ds DiskStoreService) StoreBookArchive(bookInputFolder, bookTempFilePath, bookOutputPath string) error {
	bookOutputFolder := filepath.Dir(bookOutputPath)
	err := createFolderIfNotExist(bookOutputFolder)
	if err != nil {
		return fmt.Errorf("can not create a book archive subfolder: %w", err)
	}

	err = os.Rename(bookTempFilePath, bookOutputPath)
	if err != nil {
		return fmt.Errorf("can not move a book archive to the output folder: %w", err)
	}

	err = cleanup(bookInputFolder)
	if err != nil {
		return fmt.Errorf("can not cleanup the book input folder: %w", err)
	}

	return nil
}

// StoreCoverFile moves a book cover file from the temp folder
// to the output folder. Creates the output folder if not exist.
func (ds DiskStoreService) StoreCoverFile(tempFilePath, coverOutputPath string) error {
	coverOutputFolder := filepath.Dir(coverOutputPath)
	err := createFolderIfNotExist(coverOutputFolder)
	if err != nil {
		return fmt.Errorf("can not create a book cover subfolder: %w", err)
	}

	err = os.Rename(tempFilePath, coverOutputPath)
	if err != nil {
		return fmt.Errorf("can not move a book cover to the output folder: %w", err)
	}

	return nil
}

// IsFolderEmpty returns 'true' if the folder is empty, otherwise returns false.
func (ds DiskStoreService) IsFolderEmpty(path string) (bool, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return true, err
	}

	return len(dirEntries) == 0, nil
}

func getFileSize(filePath string) (int64, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return stat.Size(), nil
}

func getFilesTypes(fileNames []string) []string {
	result := make([]string, len(fileNames))
	for i, name := range fileNames {
		result[i] = name[strings.LastIndex(name, ".")+1:]
	}
	return deduplicateSlice(result)
}

func deduplicateSlice(slice []string) (result []string) {
	uniqueMap := make(map[string]bool)
	for _, entry := range slice {
		if _, ok := uniqueMap[entry]; !ok {
			uniqueMap[entry] = true
			result = append(result, entry)
		}
	}

	return
}

func createFolderIfNotExist(folder string) error {
	_, err := os.Stat(folder)
	if err != nil && os.IsNotExist(err) {
		log.Printf("Creating a new folder: %s", folder)
		err := os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func cleanup(filesFolder string) error {
	err := filepath.WalkDir(filesFolder, func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			removeErr := os.Remove(filepath.Join(filesFolder, entry.Name()))
			if removeErr != nil {
				return removeErr
			}
			log.Printf("[INFO] - removed %q file", entry.Name())
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
