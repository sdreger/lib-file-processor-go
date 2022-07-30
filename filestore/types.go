package filestore

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
