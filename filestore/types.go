package filestore

//go:generate mockgen -destination=./compression_service_mock.go -package=filestore github.com/sdreger/lib-file-processor-go/filestore BookCompressor
type BookCompressor interface {
	CompressBookFiles(filesFolder, archiveOutputFolder, archiveFileName string) (string, []string, error)
}

//go:generate mockgen -destination=./download_service_mock.go -package=filestore github.com/sdreger/lib-file-processor-go/filestore CoverDownloader
type CoverDownloader interface {
	DownloadCoverFile(coverURL, coverOutputFolder, coverFileName string) (string, error)
}
