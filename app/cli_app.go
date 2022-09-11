package app

import (
	"database/sql"
	"github.com/sdreger/lib-file-processor-go/config"
	"github.com/sdreger/lib-file-processor-go/domain/author"
	"github.com/sdreger/lib-file-processor-go/domain/book"
	"github.com/sdreger/lib-file-processor-go/domain/category"
	"github.com/sdreger/lib-file-processor-go/domain/filetype"
	"github.com/sdreger/lib-file-processor-go/domain/lang"
	"github.com/sdreger/lib-file-processor-go/domain/publisher"
	"github.com/sdreger/lib-file-processor-go/domain/tag"
	"github.com/sdreger/lib-file-processor-go/filestore"
	"github.com/sdreger/lib-file-processor-go/scrapper"
	"log"
)

type CliApp struct {
	*core
}

func NewCliApp(config config.AppConfig, db *sql.DB, blobStore filestore.BlobStore, logger *log.Logger) (CliApp, error) {
	compressionService := filestore.NewCompressionService(logger)
	downloadService := filestore.NewDownloadService(logger)
	diskStoreService := filestore.NewDiskStoreService(compressionService, downloadService, logger)

	bookDataScrapper, err := scrapper.NewAmazonScrapper("", logger)
	if err != nil {
		return CliApp{}, err
	}

	// initStores initializes all book-related stores
	authorStore := author.NewPostgresStore(db, logger)
	categoryStore := category.NewPostgresStore(db, logger)
	fileTypeStore := filetype.NewPostgresStore(db, logger)
	tagStore := tag.NewPostgresStore(db, logger)
	publisherStore := publisher.NewPostgresStore(db, logger)
	languageStore := lang.NewPostgresStore(db, logger)
	bookDBStore := book.
		NewPostgresStore(db, publisherStore, languageStore, authorStore, categoryStore, fileTypeStore, tagStore, logger)

	return CliApp{
		core: NewCore(config, bookDBStore, blobStore, diskStoreService, bookDataScrapper, logger),
	}, nil
}

func (a CliApp) Run() {
	for {
		// -------------------- Get book ID --------------------
		bookIDString := a.getBookID()

		// -------------------- Parse book page and prepare book files --------------------
		parsedData, existingData, tempFilesData := a.PrepareBook(bookIDString)
		if tempFilesData == nil {
			a.core.Logger.Printf("[INFO] - The book file name is copied to clipboard: %s",
				parsedData.GetBookFileNameWithoutExtension())
			a.core.Logger.Printf("[WARN] - Attention, there are no book files! Skipping further processing!")
			continue
		}

		// -------------------- Approval checkpoint --------------------
		askForApproval(parsedData, existingData, a.Config.NewLineDelimiter)

		// -------------------- Store book files --------------------
		a.StoreBook(parsedData, existingData, tempFilesData)
	}
}
