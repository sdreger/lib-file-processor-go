package filestore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompressBookFiles(t *testing.T) {
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
