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

type CompressionService struct {
	logger *log.Logger
}

func NewCompressionService(logger *log.Logger) CompressionService {
	return CompressionService{
		logger: logger,
	}
}

// CompressBookFiles creates a compressed archive with files from the 'filesFolder'
// and returns a file path of the created archive, file names in archive.
func (cs CompressionService) CompressBookFiles(filesFolder, archiveOutputFolder, archiveFileName string) (string, []string, error) {

	// TODO: recursively handle folders
	fileNames, err := cs.getFilesForCompression(filesFolder)
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
	defer cs.closeResource(zipArchive)

	err = cs.compressZip(zipArchive, filesFolder, fileNames)
	if err != nil {
		return "", nil, err
	}

	return bookArchiveOutputPath, fileNames, nil
}

// ExtractZipFile extracts all book files from a compressed 'zip' archive located in the 'zipFilePath'
// to the 'pathToExtract'. Removes the source 'zip' file after successful extraction. Returns an error if any.
func (cs CompressionService) ExtractZipFile(zipFilePath string, pathToExtract string) error {
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}

	for _, f := range zipReader.File {
		outputPath := filepath.Join(pathToExtract, f.Name)
		cs.logger.Printf("unzipping file: %q", outputPath)

		if f.FileInfo().IsDir() {
			cs.logger.Printf("creating directory: %q", pathToExtract)
			err := os.MkdirAll(pathToExtract, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		srcFile, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}

		err = os.Chtimes(outputPath, f.Modified, f.Modified)
		if err != nil {
			return err
		}

		err = dstFile.Close()
		if err != nil {
			return err
		}

		err = srcFile.Close()
		if err != nil {
			return err
		}
	}

	err = zipReader.Close()
	if err != nil {
		return err
	}

	err = os.Remove(zipFilePath)
	if err != nil {
		return err
	}

	return nil
}

// getFilesForCompression returns a list of files to be compressed from a particular directory.
func (cs CompressionService) getFilesForCompression(fileDir string) ([]string, error) {
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
func (cs CompressionService) compressZip(archive *os.File, filesFolder string, fileNames []string) error {
	zipWriter := zip.NewWriter(archive)
	// Register a custom Deflate compressor.
	zipWriter.RegisterCompressor(compressionMethod, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, compressionLevel)
	})
	defer cs.closeResource(zipWriter)

	for _, fileName := range fileNames {
		err := cs.addFileToZip(zipWriter, filepath.Join(filesFolder, fileName), fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

// addFileToZip adds a single file to a zip archive
func (cs CompressionService) addFileToZip(zipWriter *zip.Writer, path, fileName string) error {
	fileToZip, err := os.Open(path)
	if err != nil {
		return err
	}
	defer cs.closeResource(fileToZip)

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

func (cs CompressionService) closeResource(f io.Closer) {
	err := f.Close()
	if err != nil {
		cs.logger.Fatal(err.Error())
	}
}
