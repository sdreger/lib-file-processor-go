package app

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/sdreger/lib-file-processor-go/config"
	"github.com/sdreger/lib-file-processor-go/domain/book"
	"github.com/sdreger/lib-file-processor-go/filestore"
	"github.com/sdreger/lib-file-processor-go/scrapper"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCore_PrepareBook(t *testing.T) {
	t.Log("Given the need to test book files preparing.")
	t.Run("Book exists in the DB, and book files present", testWithBookAndWithFiles)
	t.Run("Book does not exist in the DB, and book files present", testWithoutBookAndWithFiles)
	t.Run("Book exists in the DB, and book files absent", testWithBookAndWithoutFiles)
	t.Run("Book does not exist in the DB, and book files absent", testWithoutBookAndWithoutFiles)
	t.Logf("\t%s\tShould successfully prepare book files", succeed)
}
func testWithBookAndWithFiles(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appConfig := config.GetAppConfig()
	appConfig.DBAvailable = true
	appConfig.BlobStoreAvailable = true

	mockBlobStore := filestore.NewMockBlobStore(ctrl)

	mockBookDataScrapper := scrapper.NewMockBookDataScrapper(ctrl)
	testParsedData := getTestParsedData()
	mockBookDataScrapper.EXPECT().GetBookData(testBookID).Return(testParsedData, nil).Times(1)

	mockBookDBStore := book.NewMockStore(ctrl)
	searchRequest := book.SearchRequest{
		Title:   testParsedData.Title,
		Edition: testParsedData.Edition,
		ISBN10:  testParsedData.ISBN10,
		ISBN13:  testParsedData.ISBN13,
		ASIN:    testParsedData.ASIN,
	}
	testStoredData := getTestStoredData()
	// book found
	mockBookDBStore.EXPECT().Find(gomock.Any(), gomock.Eq(searchRequest)).Return(&testStoredData, nil).Times(1)

	mockDiskStore := filestore.NewMockDiskStore(ctrl)
	// book files present
	mockDiskStore.EXPECT().IsFolderEmpty(appConfig.BookInputFolder).Return(false, nil).Times(1)
	testTempFilesData := getTestTempFilesData()
	mockDiskStore.EXPECT().PrepareBookFiles(testParsedData, appConfig.BookInputFolder, appConfig.TempInputFolder).
		Return(testTempFilesData, nil).Times(1)

	coreApp := NewCore(appConfig, mockBookDBStore, mockBlobStore, mockDiskStore, mockBookDataScrapper)

	updatedParsedData, storedData, tempFilesData := coreApp.PrepareBook(testBookID)

	if updatedParsedData.BookFileSize != testBookFileSize {
		t.Fatalf("\t\t%s\tShould get %d book file size: %d", failed, testBookFileSize, updatedParsedData.BookFileSize)
	}
	if !reflect.DeepEqual(updatedParsedData.Formats, testBookFormats) {
		t.Fatalf("\t\t%s\tShould get %v book formats: %v", failed, testBookFormats, updatedParsedData.Formats)
	}

	if storedData == nil || !reflect.DeepEqual(*storedData, testStoredData) {
		t.Fatalf("\t\t%s\tStored data invalid", failed)
	}

	if tempFilesData == nil || !reflect.DeepEqual(*tempFilesData, testTempFilesData) {
		t.Fatalf("\t\t%s\tTemp files data invalid", failed)
	}
}

func testWithoutBookAndWithFiles(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appConfig := config.GetAppConfig()
	appConfig.DBAvailable = true
	appConfig.BlobStoreAvailable = true

	mockBlobStore := filestore.NewMockBlobStore(ctrl)

	mockBookDataScrapper := scrapper.NewMockBookDataScrapper(ctrl)
	testParsedData := getTestParsedData()
	mockBookDataScrapper.EXPECT().GetBookData(testBookID).Return(testParsedData, nil).Times(1)

	mockBookDBStore := book.NewMockStore(ctrl)
	searchRequest := book.SearchRequest{
		Title:   testParsedData.Title,
		Edition: testParsedData.Edition,
		ISBN10:  testParsedData.ISBN10,
		ISBN13:  testParsedData.ISBN13,
		ASIN:    testParsedData.ASIN,
	}

	// book not found
	mockBookDBStore.EXPECT().Find(gomock.Any(), gomock.Eq(searchRequest)).Return(nil, nil).Times(1)

	mockDiskStore := filestore.NewMockDiskStore(ctrl)
	// book files present
	mockDiskStore.EXPECT().IsFolderEmpty(appConfig.BookInputFolder).Return(false, nil).Times(1)
	testTempFilesData := getTestTempFilesData()
	mockDiskStore.EXPECT().PrepareBookFiles(testParsedData, appConfig.BookInputFolder, appConfig.TempInputFolder).
		Return(testTempFilesData, nil).Times(1)

	coreApp := NewCore(appConfig, mockBookDBStore, mockBlobStore, mockDiskStore, mockBookDataScrapper)

	updatedParsedData, storedData, tempFilesData := coreApp.PrepareBook(testBookID)

	if updatedParsedData.BookFileSize != testBookFileSize {
		t.Fatalf("\t\t%s\tShould get %d book file size: %d", failed, testBookFileSize, updatedParsedData.BookFileSize)
	}
	if !reflect.DeepEqual(updatedParsedData.Formats, testBookFormats) {
		t.Fatalf("\t\t%s\tShould get %v book formats: %v", failed, testBookFormats, updatedParsedData.Formats)
	}

	if storedData != nil {
		t.Fatalf("\t\t%s\tStored data invalid", failed)
	}

	if tempFilesData == nil || !reflect.DeepEqual(*tempFilesData, testTempFilesData) {
		t.Fatalf("\t\t%s\tTemp files data invalid", failed)
	}
}

func testWithBookAndWithoutFiles(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appConfig := config.GetAppConfig()
	appConfig.DBAvailable = true
	appConfig.BlobStoreAvailable = true

	mockBlobStore := filestore.NewMockBlobStore(ctrl)

	mockBookDataScrapper := scrapper.NewMockBookDataScrapper(ctrl)
	testParsedData := getTestParsedData()
	mockBookDataScrapper.EXPECT().GetBookData(testBookID).Return(testParsedData, nil).Times(1)

	mockBookDBStore := book.NewMockStore(ctrl)
	searchRequest := book.SearchRequest{
		Title:   testParsedData.Title,
		Edition: testParsedData.Edition,
		ISBN10:  testParsedData.ISBN10,
		ISBN13:  testParsedData.ISBN13,
		ASIN:    testParsedData.ASIN,
	}
	testStoredData := getTestStoredData()
	// book found
	mockBookDBStore.EXPECT().Find(gomock.Any(), gomock.Eq(searchRequest)).Return(&testStoredData, nil).Times(1)

	mockDiskStore := filestore.NewMockDiskStore(ctrl)
	// book files not present
	mockDiskStore.EXPECT().IsFolderEmpty(appConfig.BookInputFolder).Return(true, nil).Times(1)

	coreApp := NewCore(appConfig, mockBookDBStore, mockBlobStore, mockDiskStore, mockBookDataScrapper)

	updatedParsedData, storedData, tempFilesData := coreApp.PrepareBook(testBookID)

	if updatedParsedData.BookFileSize != 0 {
		t.Fatalf("\t\t%s\tShould get 0 book file size: %d", failed, updatedParsedData.BookFileSize)
	}
	if len(updatedParsedData.Formats) > 0 {
		t.Fatalf("\t\t%s\tShould get no book formats: %v", failed, updatedParsedData.Formats)
	}

	if storedData == nil || !reflect.DeepEqual(*storedData, testStoredData) {
		t.Fatalf("\t\t%s\tStored data invalid", failed)
	}

	if tempFilesData != nil {
		t.Fatalf("\t\t%s\tTemp files data invalid", failed)
	}
}

func testWithoutBookAndWithoutFiles(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appConfig := config.GetAppConfig()
	appConfig.DBAvailable = true
	appConfig.BlobStoreAvailable = true

	mockBlobStore := filestore.NewMockBlobStore(ctrl)

	mockBookDataScrapper := scrapper.NewMockBookDataScrapper(ctrl)
	testParsedData := getTestParsedData()
	mockBookDataScrapper.EXPECT().GetBookData(testBookID).Return(testParsedData, nil).Times(1)

	mockBookDBStore := book.NewMockStore(ctrl)
	searchRequest := book.SearchRequest{
		Title:   testParsedData.Title,
		Edition: testParsedData.Edition,
		ISBN10:  testParsedData.ISBN10,
		ISBN13:  testParsedData.ISBN13,
		ASIN:    testParsedData.ASIN,
	}
	// book not found
	mockBookDBStore.EXPECT().Find(gomock.Any(), gomock.Eq(searchRequest)).Return(nil, nil).Times(1)

	mockDiskStore := filestore.NewMockDiskStore(ctrl)
	// book files not present
	mockDiskStore.EXPECT().IsFolderEmpty(appConfig.BookInputFolder).Return(true, nil).Times(1)

	coreApp := NewCore(appConfig, mockBookDBStore, mockBlobStore, mockDiskStore, mockBookDataScrapper)

	updatedParsedData, storedData, tempFilesData := coreApp.PrepareBook(testBookID)

	if updatedParsedData.BookFileSize != 0 {
		t.Fatalf("\t\t%s\tShould get 0 book file size: %d", failed, updatedParsedData.BookFileSize)
	}
	if len(updatedParsedData.Formats) > 0 {
		t.Fatalf("\t\t%s\tShould get no book formats: %v", failed, updatedParsedData.Formats)
	}

	if storedData != nil {
		t.Fatalf("\t\t%s\tStored data invalid", failed)
	}

	if tempFilesData != nil {
		t.Fatalf("\t\t%s\tTemp files data invalid", failed)
	}
}

func TestCore_StoreBook(t *testing.T) {
	t.Log("Given the need to test book files store.")
	t.Run("There is an existing book data", testWithExistingData)
	t.Run("There is no existing book data", testWithoutExistingData)
	t.Logf("\t%s\tShould successfully store book files", succeed)
}

func testWithExistingData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testParsedData := getTestParsedData()
	testStoredData := getTestStoredData()
	testTempFilesData := getTestTempFilesData()

	appConfig := config.GetAppConfig()
	appConfig.DBAvailable = true
	appConfig.BlobStoreAvailable = true

	mockDiskStore := filestore.NewMockDiskStore(ctrl)
	lowerPublisher := strings.ToLower(testParsedData.Publisher)
	bookArchiveOutputPath := filepath.Join(appConfig.BookOutputFolder, lowerPublisher, testParsedData.BookFileName)
	coverOutputPath := filepath.Join(appConfig.CoverOutputFolder, lowerPublisher, testParsedData.CoverFileName)
	mockDiskStore.EXPECT().
		StoreBookArchive(appConfig.BookInputFolder, testTempFilesData.BookArchivePath, bookArchiveOutputPath).
		Return(nil).Times(1)
	mockDiskStore.EXPECT().
		StoreCoverFile(testTempFilesData.CoverFilePath, coverOutputPath).
		Return(nil).Times(1)

	mockBookDBStore := book.NewMockStore(ctrl)
	mockBookDBStore.EXPECT().Update(gomock.Any(), gomock.Eq(&testStoredData), gomock.Eq(&testParsedData)).
		Return(nil).Times(1)

	mockBlobStore := filestore.NewMockBlobStore(ctrl)
	mockBlobStore.EXPECT().StoreObject(gomock.Any(), gomock.Eq(bookBucketName),
		gomock.Eq(fmt.Sprintf("%s/%s", lowerPublisher, testParsedData.BookFileName)), gomock.Eq(bookArchiveOutputPath)).
		Return(testBookEtag, nil).Times(1)
	mockBlobStore.EXPECT().StoreObject(gomock.Any(), gomock.Eq(coverBucketName),
		gomock.Eq(fmt.Sprintf("%s/%s", lowerPublisher, testParsedData.CoverFileName)), gomock.Eq(coverOutputPath)).
		Return(testCoverEtag, nil).Times(1)

	coreApp := NewCore(appConfig, mockBookDBStore, mockBlobStore, mockDiskStore, nil)
	coreApp.StoreBook(&testParsedData, &testStoredData, &testTempFilesData)
}

func testWithoutExistingData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testParsedData := getTestParsedData()
	testTempFilesData := getTestTempFilesData()

	appConfig := config.GetAppConfig()
	appConfig.DBAvailable = true
	appConfig.BlobStoreAvailable = true

	mockDiskStore := filestore.NewMockDiskStore(ctrl)
	lowerPublisher := strings.ToLower(testParsedData.Publisher)
	bookArchiveOutputPath := filepath.Join(appConfig.BookOutputFolder, lowerPublisher, testParsedData.BookFileName)
	coverOutputPath := filepath.Join(appConfig.CoverOutputFolder, lowerPublisher, testParsedData.CoverFileName)
	mockDiskStore.EXPECT().
		StoreBookArchive(appConfig.BookInputFolder, testTempFilesData.BookArchivePath, bookArchiveOutputPath).
		Return(nil).Times(1)
	mockDiskStore.EXPECT().
		StoreCoverFile(testTempFilesData.CoverFilePath, coverOutputPath).
		Return(nil).Times(1)

	mockBookDBStore := book.NewMockStore(ctrl)
	mockBookDBStore.EXPECT().Add(gomock.Any(), gomock.Eq(testParsedData)).Return(testBookIDInt, nil).Times(1)

	mockBlobStore := filestore.NewMockBlobStore(ctrl)
	mockBlobStore.EXPECT().StoreObject(gomock.Any(), gomock.Eq(bookBucketName),
		gomock.Eq(fmt.Sprintf("%s/%s", lowerPublisher, testParsedData.BookFileName)), gomock.Eq(bookArchiveOutputPath)).
		Return(testBookEtag, nil).Times(1)
	mockBlobStore.EXPECT().StoreObject(gomock.Any(), gomock.Eq(coverBucketName),
		gomock.Eq(fmt.Sprintf("%s/%s", lowerPublisher, testParsedData.CoverFileName)), gomock.Eq(coverOutputPath)).
		Return(testCoverEtag, nil).Times(1)

	coreApp := NewCore(appConfig, mockBookDBStore, mockBlobStore, mockDiskStore, nil)
	coreApp.StoreBook(&testParsedData, nil, &testTempFilesData)
}
