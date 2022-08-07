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

func NewCliApp(config config.AppConfig, db *sql.DB, blobStore filestore.BlobStore) (CliApp, error) {
	compressionService := filestore.NewCompressionService()
	downloadService := filestore.NewDownloadService()
	diskStoreService := filestore.NewDiskStoreService(compressionService, downloadService)

	bookDataScrapper, err := scrapper.NewAmazonScrapper("")
	if err != nil {
		return CliApp{}, err
	}

	// initStores initializes all book-related stores
	authorStore := author.NewPostgresStore(db)
	categoryStore := category.NewPostgresStore(db)
	fileTypeStore := filetype.NewPostgresStore(db)
	tagStore := tag.NewPostgresStore(db)
	publisherStore := publisher.NewPostgresStore(db)
	languageStore := lang.NewPostgresStore(db)
	bookDBStore := book.
		NewPostgresStore(db, publisherStore, languageStore, authorStore, categoryStore, fileTypeStore, tagStore)

	return CliApp{
		core: NewCore(config, bookDBStore, blobStore, diskStoreService, bookDataScrapper),
	}, nil
}

func (a CliApp) Run() {
	for {
		// -------------------- Get book ID --------------------
		bookIDString := a.getBookID()

		// -------------------- Parse book page and prepare book files --------------------
		parsedData, existingData, tempFilesData := a.PrepareBook(bookIDString)
		if tempFilesData == nil {
			log.Printf("[INFO] - The book file name is copied to clipboard: %s",
				parsedData.GetBookFileNameWithoutExtension())
			log.Printf("[WARN] - Attention, there are no book files! Skipping further processing!")
			continue
		}

		// -------------------- Approval checkpoint --------------------
		askForApproval(parsedData, existingData, a.Config.NewLineDelimiter)

		// -------------------- Store book files --------------------
		a.StoreBook(parsedData, existingData, tempFilesData)
	}
}
