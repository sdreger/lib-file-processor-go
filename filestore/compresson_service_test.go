package filestore

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestCompressBookFiles(t *testing.T) {
	t.Log("Given the need to test book files compression.")
	archiveFileName := "test-archive.zip"
	inputFileName01 := "input1.pdf"
	inputFileName02 := "input2.epub"

	t.Log("Given the need to test book files compression.")
	tempInputDir, err := os.MkdirTemp("", "input-dir-*")
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create input folder: %v", failed, err)
	}
	defer os.RemoveAll(tempInputDir)

	tempOutputDir, err := os.MkdirTemp("", "output-dir-*")
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create output folder: %v", failed, err)
	}
	defer os.RemoveAll(tempOutputDir)

	tempFile01, err := os.CreateTemp(tempInputDir, "*-"+inputFileName01)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create an input file: %v", failed, err)
	}
	tempFile02, err := os.CreateTemp(tempInputDir, "*-"+inputFileName02)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create an input file: %v", failed, err)
	}

	archiveFilePath, filesInArchive, err :=
		NewCompressionService().CompressBookFiles(tempInputDir, tempOutputDir, archiveFileName)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to comress book files: %v", failed, err)
	}

	if archiveFilePath != filepath.Join(tempOutputDir, archiveFileName) {
		t.Fatalf("\t\t%s\tWrong archive path", failed)
	}
	stat, err := os.Stat(archiveFilePath)
	if err != nil {
		t.Fatalf("\t\t%s\tShold be able to get archive file stats: %v", failed, err)
	}
	if stat.Size() == 0 {
		t.Fatalf("\t\t%s\tThe archive file size should be grater than 0", failed)
	}

	if len(filesInArchive) != 2 {
		t.Fatalf("\t\t%s\tThe archive should contain 2 files", failed)
	}
	tempFileName01 := filepath.Base(tempFile01.Name())
	if filesInArchive[0] != tempFileName01 && filesInArchive[1] != tempFileName01 {
		t.Fatalf("\t\t%s\tThe filesInArchive slice should contain file: %s", failed, tempFileName01)
	}
	tempFileName02 := filepath.Base(tempFile02.Name())
	if filesInArchive[0] != tempFileName02 && filesInArchive[1] != tempFileName02 {
		t.Fatalf("\t\t%s\tThe filesInArchive slice should contain file: %s", failed, tempFileName02)
	}

	t.Logf("\t\t%s\tShould successfully compress book files", succeed)
}

func TestExtractBookFiles(t *testing.T) {
	t.Log("Given the need to test book files extraction.")
	// temporary folder to use as 'zip' file input folder
	tempInputDir, err := os.MkdirTemp("", "extract-input-dir-*")
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create output folder: %v", failed, err)
	}
	defer os.RemoveAll(tempInputDir)

	// test file will be removed, so we need to copy it to a temporary folder
	testZipFile, err := os.Open(filepath.Join("testdata", "test.zip"))
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to open test file: %v", failed, err)
	}
	tempZipFile, err := os.Create(filepath.Join(tempInputDir, "test.zip"))
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create temporary test file: %v", failed, err)
	}
	_, err = io.Copy(tempZipFile, testZipFile)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to copy test file: %v", failed, err)
	}

	// temporary folder for 'zip' file extraction
	tempOutputDir, err := os.MkdirTemp("", "extract-output-dir-*")
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to create output folder: %v", failed, err)
	}
	defer os.RemoveAll(tempOutputDir)

	err = NewCompressionService().ExtractZipFile(tempZipFile.Name(), tempOutputDir)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to extract zip file: %v", failed, err)
	}

	// checks if there is an extracted file in the output folder
	outputDirEntries, err := os.ReadDir(tempOutputDir)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to read output folder: %v", failed, err)
	}
	if len(outputDirEntries) != 1 {
		t.Fatalf("\t\t%s\tOutput folder should contain exactly 1 file", failed)
	}
	if outputDirEntries[0].Name() != "test.txt" {
		t.Fatalf("\t\t%s\tOutput folder contains wrong extracted file: %q", failed, outputDirEntries[0].Name())
	}
	if fileInfo, err := outputDirEntries[0].Info(); err != nil || fileInfo.Size() == 0 {
		t.Fatalf("\t\t%s\tExtracted file is corrupted: %v", failed, err)
	}

	// check if the source 'zip' file was removed
	inputDirEntries, err := os.ReadDir(tempInputDir)
	if err != nil {
		t.Fatalf("\t\t%s\tShould be able to read input folder: %v", failed, err)
	}
	if len(inputDirEntries) != 0 {
		t.Fatalf("\t\t%s\tInput folder should be empty", failed)
	}

	t.Logf("\t\t%s\tShould successfully extract a book file", succeed)
}
