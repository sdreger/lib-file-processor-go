package filestore

import (
	"archive/zip"
	"compress/flate"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	compressionMethod = zip.Deflate
	compressionLevel  = flate.BestCompression
)

type CompressionService struct{}

func NewCompressionService() CompressionService {
	return CompressionService{}
}

// CompressBookFiles creates a compressed archive with files from the 'filesFolder'
// and returns a file path of the created archive, file names in archive.
func (cs CompressionService) CompressBookFiles(filesFolder, archiveOutputFolder, archiveFileName string) (string, []string, error) {

	// TODO: recursively handle folders
	fileNames, err := getFilesForCompression(filesFolder)
	if err != nil {
		return "", nil, err
	}
	if len(fileNames) == 0 {
		return "", nil, fmt.Errorf("there are no files to compress")
	}

	bookArchiveOutputPath := filepath.Join(archiveOutputFolder, archiveFileName)
	zipArchive, err := os.Create(bookArchiveOutputPath)
	if err != nil {
		return "", nil, err
	}
	defer closeResource(zipArchive)

	err = compressZip(zipArchive, filesFolder, fileNames)
	if err != nil {
		return "", nil, err
	}

	return bookArchiveOutputPath, fileNames, nil
}

// getFilesForCompression returns a list of files to be compressed from a particular directory.
func getFilesForCompression(fileDir string) ([]string, error) {
	dirEntries, err := os.ReadDir(fileDir)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(dirEntries))
	for i, dirEntry := range dirEntries {
		result[i] = dirEntry.Name()
	}
	return result, nil
}

// compressZip compresses file list into a zip archive.
func compressZip(archive *os.File, filesFolder string, fileNames []string) error {
	zipWriter := zip.NewWriter(archive)
	// Register a custom Deflate compressor.
	zipWriter.RegisterCompressor(compressionMethod, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, compressionLevel)
	})
	defer closeResource(zipWriter)

	for _, fileName := range fileNames {
		err := addFileToZip(zipWriter, filepath.Join(filesFolder, fileName), fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

// addFileToZip adds a single file to a zip archive
func addFileToZip(zipWriter *zip.Writer, path, fileName string) error {
	fileToZip, err := os.Open(path)
	if err != nil {
		return err
	}
	defer closeResource(fileToZip)

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = fileName
	header.Method = compressionMethod

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		return err
	}

	return nil
}

func closeResource(f io.Closer) {
	err := f.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}
