package filestore

import (
	"context"
	"github.com/sdreger/lib-file-processor-go/domain/book"
)

//go:generate mockgen -destination=./book_compressor_mock.go -package=filestore github.com/sdreger/lib-file-processor-go/filestore BookCompressor
type BookCompressor interface {
	CompressBookFiles(filesFolder, archiveOutputFolder, archiveFileName string) (string, []string, error)
}

//go:generate mockgen -destination=./book_extractor_mock.go -package=filestore github.com/sdreger/lib-file-processor-go/filestore BookExtractor
type BookExtractor interface {
	ExtractZipFile(zipFilePath string, pathToExtract string) error
}

//go:generate mockgen -destination=./download_service_mock.go -package=filestore github.com/sdreger/lib-file-processor-go/filestore CoverDownloader
type CoverDownloader interface {
	DownloadCoverFile(coverURL, coverOutputFolder, coverFileName string) (string, error)
}

//go:generate mockgen -destination=./disk_store_mock.go -package=filestore github.com/sdreger/lib-file-processor-go/filestore DiskStore
type DiskStore interface {
	PrepareBookFiles(bookMeta book.ParsedData, bookInputFolder, outputFolder string) (TempFilesData, error)
	StoreBookArchive(bookInputFolder, bookTempFilePath, bookOutputPath string) error
	StoreCoverFile(tempFilePath, coverOutputPath string) error
	IsFolderEmpty(path string) (bool, error)
}

//go:generate mockgen -destination=./blob_store_mock.go -package=filestore github.com/sdreger/lib-file-processor-go/filestore BlobStore
type BlobStore interface {
	CreateBucket(ctx context.Context, bucketName string) error
	StoreObject(ctx context.Context, bucketName string, fileName, filePath string) (string, error)
}
