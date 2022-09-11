package app

import (
	"bufio"
	"context"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/sdreger/lib-file-processor-go/config"
	"github.com/sdreger/lib-file-processor-go/domain/book"
	"github.com/sdreger/lib-file-processor-go/filestore"
	"github.com/sdreger/lib-file-processor-go/scrapper"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	bookBucketName  = "ebooks"
	coverBucketName = "ebook-covers"
)

type core struct {
	Config           config.AppConfig
	BookDBStore      book.Store
	BookDiskStore    filestore.DiskStore
	BookBlobStore    filestore.BlobStore
	BookDataScrapper scrapper.BookDataScrapper
	Logger           *log.Logger
}

func NewCore(config config.AppConfig, bookDBStore book.Store, blobStore filestore.BlobStore,
	diskStoreService filestore.DiskStore, bookDataScrapper scrapper.BookDataScrapper, logger *log.Logger) *core {
	return &core{
		Config:           config,
		BookDBStore:      bookDBStore,
		BookDiskStore:    diskStoreService,
		BookBlobStore:    blobStore,
		BookDataScrapper: bookDataScrapper,
		Logger:           logger,
	}
}

// PrepareBook scrapes and parse a book page, downloads the book cover image.
// If there are book files - compress and put them into temporary folder,
// otherwise - copies the book file name to clipboard, and skips the compression.
func (c *core) PrepareBook(bookIDString string) (*book.ParsedData, *book.StoredData, *filestore.TempFilesData) {
	var existingData *book.StoredData

	// -------------------- Parse book page --------------------
	parsedData, err := c.BookDataScrapper.GetBookData(bookIDString)
	if err != nil {
		c.Logger.Fatalf("Can not scrape a book metadata: %v", err)
	}

	// -------------------- Check if there are book files --------------------
	folderIsEmpty, err := c.BookDiskStore.IsFolderEmpty(c.Config.BookInputFolder)
	if err != nil {
		c.Logger.Fatalf("Can not check if input folder is empty: %v", err)
	}

	if c.Config.DBAvailable {
		// -------------------- Search for existing book --------------------
		ctx := context.Background()
		existingData, err = c.findExistingBook(ctx, parsedData)
		if err != nil {
			c.Logger.Fatalf("Failed to find a book: %v", err)
		}
	}

	// -------------------- Filename copy mode (if there are no book files) --------------------
	// If there are no files in the input folder, just copy a book file name to the system clipboard
	if folderIsEmpty || c.Config.IsStatelessMode() {
		c.copyToClipboard(parsedData.GetBookFileNameWithoutExtension())
		//c.Logger.Printf("[INFO] - The book file name is copied to clipboard: %s",
		//	parsedData.BookFileName
		//c.Logger.Printf("[WARN] - Attention, there are no book files! Skipping further processing!")
		return &parsedData, existingData, nil
	}

	// -------------------- Prepare book files --------------------
	tempFilesData, err := c.BookDiskStore.PrepareBookFiles(parsedData, c.Config.BookInputFolder, c.Config.TempInputFolder)
	if err != nil {
		c.Logger.Fatalf("Can not prepare book files: %v", err)
	}
	parsedData.BookFileSize = tempFilesData.BookSize
	parsedData.Formats = tempFilesData.BookFormats

	return &parsedData, existingData, &tempFilesData
}

// StoreBook inserts a new book record to database (or updates an existing one if any).
// Moves book archive and book cover to output folder. Stores book archive and book cover to BLOB store.
func (c *core) StoreBook(parsedData *book.ParsedData, existingData *book.StoredData, tempData *filestore.TempFilesData) {

	publisherLowerName := strings.ToLower(parsedData.Publisher)
	coverOutputPath := filepath.Join(c.Config.CoverOutputFolder, publisherLowerName, parsedData.CoverFileName)
	bookArchiveOutputPath := filepath.Join(c.Config.BookOutputFolder, publisherLowerName, parsedData.BookFileName)

	// -------------------- Store book files --------------------
	err := c.storeBookFiles(tempData, bookArchiveOutputPath, coverOutputPath)
	if err != nil {
		c.Logger.Fatalf("Can not store book files: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// -------------------- Store / Update DB data --------------------
	if c.Config.DBAvailable {
		bookID, err := c.upsertBook(ctx, parsedData, existingData)
		if err != nil {
			c.Logger.Fatalf("Can not upsert a book: %v", err)
		}
		_ = bookID
		//fmt.Println("Stored / Updated book ID:", bookID)
	}

	// -------------------- Store book objects --------------------
	if c.Config.BlobStoreAvailable {
		err := c.storeBookObjects(ctx, parsedData.BookFileName, parsedData.CoverFileName,
			publisherLowerName, bookArchiveOutputPath, coverOutputPath)
		if err != nil {
			c.Logger.Fatalf("Can not store book objects: %v", err)
		}
	}

	cancel()
}

func (c *core) getBookID() string {
	fmt.Print("Enter ISBN10/ASIN: ")
	reader := bufio.NewReader(os.Stdin)

	bookIDInput, err := reader.ReadString(c.Config.NewLineDelimiter)
	if err != nil {
		c.Logger.Fatalf("Can not get book identifier: %v", err)
	}
	return strings.TrimSpace(bookIDInput[:len(bookIDInput)-1])
}

func (c *core) findExistingBook(ctx context.Context, parsedData book.ParsedData) (*book.StoredData, error) {
	existingData, err := c.BookDBStore.Find(ctx, book.SearchRequest{
		Title:   parsedData.Title,
		Edition: parsedData.Edition,
		ISBN10:  parsedData.ISBN10,
		ISBN13:  parsedData.ISBN13,
		ASIN:    parsedData.ASIN,
	})
	if err != nil {
		return nil, err
	}

	return existingData, nil
}

func (c *core) copyToClipboard(str string) {
	err := clipboard.WriteAll(str)
	if err != nil {
		c.Logger.Fatalf("Can not write to clipboard: %v", err)
	}
}

func (c *core) storeBookFiles(tempData *filestore.TempFilesData, bookArchiveOutputPath, coverOutputPath string) error {
	err := c.BookDiskStore.
		StoreBookArchive(c.Config.BookInputFolder, tempData.BookArchivePath, bookArchiveOutputPath)
	if err != nil {
		return fmt.Errorf("can not store book archive: %w", err)
	}
	err = c.BookDiskStore.StoreCoverFile(tempData.CoverFilePath, coverOutputPath)
	if err != nil {
		return fmt.Errorf("can not store book cover: %w", err)
	}

	return nil
}

func (c *core) upsertBook(ctx context.Context, parsedData *book.ParsedData, existingData *book.StoredData) (int64, error) {
	var bookID int64
	if existingData != nil {
		// -------------------- Update an existing book --------------------
		updateErr := c.BookDBStore.Update(ctx, existingData, parsedData)
		if updateErr != nil {
			return 0, fmt.Errorf("can not update a book: %w", updateErr)
		}
		bookID = existingData.ID
	} else {
		// -------------------- Add a new book --------------------
		storedBookID, storeErr := c.BookDBStore.Add(ctx, *parsedData)
		if storeErr != nil {
			return 0, fmt.Errorf("can not add a book: %w", storeErr)
		}
		bookID = storedBookID
	}

	return bookID, nil
}

func (c *core) storeBookObjects(ctx context.Context, bookFileName, coverFileName,
	publisherLowerName, bookArchiveOutputPath, coverOutputPath string) error {

	// -------------------- Store book BLOB --------------------
	objectKey := fmt.Sprintf("%s/%s", publisherLowerName, bookFileName)
	_, err := c.BookBlobStore.StoreObject(ctx, bookBucketName, objectKey, bookArchiveOutputPath)
	if err != nil {
		return fmt.Errorf("can not store a book BLOB for the object key: %q. %w", objectKey, err)
	}

	// -------------------- Store cover BLOB --------------------
	objectKey = fmt.Sprintf("%s/%s", publisherLowerName, coverFileName)
	_, err = c.BookBlobStore.StoreObject(ctx, coverBucketName, objectKey, coverOutputPath)
	if err != nil {
		return fmt.Errorf("can not store a cover BLOB for the object key: %q. %w", objectKey, err)
	}

	return nil
}
