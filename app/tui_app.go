package app

import (
	"database/sql"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
	"strconv"
	"strings"
	"time"
)

const (
	dateLayout = "_2 Jan 2006"
)

type TuiApp struct {
	*core
	bookIDChan <-chan string

	tuiApp        *tview.Application
	grid          *tview.Grid
	bookIDInput   *tview.InputField
	parsedForm    *tview.Form
	equalityTable *tview.Table
	existingTable *tview.Table
	footer        *tview.TextView

	bookIDString       string
	parsedData         *book.ParsedData
	existingData       *book.StoredData
	tempFilesData      *filestore.TempFilesData
	ignoreExistingData bool

	editErrorMap map[string]error
}

func NewTuiApp(config config.AppConfig, db *sql.DB, blobStore filestore.BlobStore, logger *log.Logger,
	bookIDChan <-chan string) (*TuiApp, error) {
	compressionService := filestore.NewCompressionService(logger)
	downloadService := filestore.NewDownloadService(logger)
	diskStoreService := filestore.NewDiskStoreService(compressionService, downloadService, logger)

	bookDataScrapper, err := scrapper.NewAmazonScrapper("", logger)
	if err != nil {
		return nil, err
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

	return &TuiApp{
		core:          NewCore(config, bookDBStore, blobStore, diskStoreService, bookDataScrapper, logger),
		bookIDChan:    bookIDChan,
		tuiApp:        tview.NewApplication(),
		grid:          tview.NewGrid(),
		bookIDInput:   tview.NewInputField(),
		parsedForm:    tview.NewForm().SetItemPadding(0).SetFieldBackgroundColor(tcell.ColorBlack),
		equalityTable: tview.NewTable().SetBorders(false),
		existingTable: tview.NewTable().SetBorders(false),
		footer:        tview.NewTextView().SetScrollable(true),
		editErrorMap:  make(map[string]error),
	}, nil
}

func (t *TuiApp) Run() error {
	t.initBookIDInput(t.bookIDInput)
	t.initGrid(t.grid)
	if err := t.tuiApp.SetRoot(t.grid, true).SetFocus(t.bookIDInput).Run(); err != nil {
		return err
	}

	return nil
}

func (t *TuiApp) initBookIDInput(input *tview.InputField) {
	input.SetLabel("Enter ISBN10/ASIN: ").
		SetFieldWidth(13).
		SetChangedFunc(func(text string) {
			t.bookIDString = text
		}).
		SetDoneFunc(t.bookIDInputHandler)

	go func() {
		for bookID := range t.bookIDChan {
			t.bookIDInput.SetText(bookID)
			t.appendFooterText(fmt.Sprintf("A new file extracted with ID: %s", bookID))
			t.tuiApp.Draw()
		}
	}()
}

func (t *TuiApp) bookIDInputHandler(key tcell.Key) {
	if key != tcell.KeyEnter {
		return
	}
	if len(t.bookIDString) != 10 {
		t.appendFooterText(fmt.Sprintf("The book ID must be of size 10: %q", t.bookIDString))
		return
	}
	parsedData, existingData, tempFilesData := t.PrepareBook(t.bookIDString)
	t.parsedData = parsedData
	t.existingData = existingData
	t.tempFilesData = tempFilesData

	t.clearForms()
	t.fillParsedForm(t.parsedForm, parsedData)
	if existingData != nil {
		t.fillCheckboxTable(t.equalityTable, parsedData, existingData)
		t.fillExisingTable(t.existingTable, parsedData, existingData)
	}
	if tempFilesData != nil {
		t.tuiApp.SetFocus(t.parsedForm)
	} else {
		t.footer.SetText(fmt.Sprintf("The book file name is copied to clipboard!\r\n%s",
			parsedData.GetBookFileNameWithoutExtension())).SetTextColor(tcell.ColorOrange)
		t.restartFlow(false)
	}
}

func (t *TuiApp) initGrid(grid *tview.Grid) {
	var dbAvailable *tview.InputField
	if t.Config.DBAvailable {
		dbAvailable = tview.NewInputField().SetLabel("DB: ").SetFieldBackgroundColor(tcell.ColorBlack).
			SetText("Available").SetFieldTextColor(tcell.ColorGreen)
	} else {
		dbAvailable = tview.NewInputField().SetLabel("DB: ").SetFieldBackgroundColor(tcell.ColorBlack).
			SetText("Unavailable").SetFieldTextColor(tcell.ColorRed)
	}

	var blobSoreAvailable *tview.InputField
	if t.Config.BlobStoreAvailable {
		blobSoreAvailable = tview.NewInputField().SetLabel("BLOB Store: ").SetFieldBackgroundColor(tcell.ColorBlack).
			SetText("Available").SetFieldTextColor(tcell.ColorGreen)
	} else {
		blobSoreAvailable = tview.NewInputField().SetLabel("BLOB Store: ").SetFieldBackgroundColor(tcell.ColorBlack).
			SetText("Unavailable").SetFieldTextColor(tcell.ColorRed)
	}
	parsedFormFrame := tview.NewFrame(t.parsedForm).SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Parsed Data", true, tview.AlignCenter, tcell.ColorYellow)
	equalityFrame := tview.NewFrame(t.equalityTable).SetBorders(0, 0, 0, 0, 1, 1).
		AddText("\u225f", true, tview.AlignCenter, tcell.ColorYellow)
	existingDataFrame := tview.NewFrame(t.existingTable).SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Existing data", true, tview.AlignCenter, tcell.ColorYellow)

	grid.
		SetRows(1, 1, 1, 0, 3).
		SetColumns(0, 3, 0).
		SetBorders(true).
		AddItem(dbAvailable, 0, 0, 1, 3, 0, 0, false).
		AddItem(blobSoreAvailable, 1, 0, 1, 3, 0, 0, false).
		AddItem(t.bookIDInput, 2, 0, 1, 3, 0, 0, false)
	grid.AddItem(t.footer, 4, 0, 1, 3, 0, 0, false)
	grid.AddItem(parsedFormFrame, 3, 0, 1, 1, 0, 0, false).
		AddItem(equalityFrame, 3, 1, 1, 1, 0, 0, false).
		AddItem(existingDataFrame, 3, 2, 1, 1, 0, 0, false)
}

func (t *TuiApp) fillParsedForm(form *tview.Form, parsedData *book.ParsedData) {
	var bookFileNameInputField *tview.InputField

	form.AddInputField("Title:", parsedData.Title, 0, nil, func(text string) {
		parsedData.Title = text
		bookFileNameInputField.SetText(t.parsedData.GetBookFileName())
	})
	form.AddInputField("Subtitle:", parsedData.Subtitle, 0, nil, func(text string) {
		parsedData.Subtitle = text
	})
	form.AddInputField("Description:", parsedData.Description, 0, nil, func(text string) {
		parsedData.Description = text
	})
	form.AddInputField("ISBN10:", parsedData.ISBN10, 0, nil, func(text string) {
		parsedData.ISBN10 = text
		bookFileNameInputField.SetText(t.parsedData.GetBookFileName())
	})
	form.AddInputField("ISBN13:", strconv.FormatInt(parsedData.ISBN13, 10), 0, nil, func(text string) {
		if len(text) != 13 {
			t.editErrorMap["ISBN13"] = fmt.Errorf("ISBN13 length is invalid: %d", len(text))
			return
		}
		isbn13, convErr := strconv.Atoi(text)
		if convErr != nil {
			t.editErrorMap["ISBN13"] = convErr
			return
		}
		delete(t.editErrorMap, "ISBN13")
		parsedData.ISBN13 = int64(isbn13)
	})
	form.AddInputField("ASIN:", parsedData.ASIN, 0, nil, func(text string) {
		parsedData.ASIN = text
		bookFileNameInputField.SetText(t.parsedData.GetBookFileName())
	})
	form.AddInputField("Pages:", strconv.FormatUint(uint64(parsedData.Pages), 10), 0, nil, func(text string) {
		pages, convErr := strconv.Atoi(text)
		if convErr != nil {
			t.editErrorMap["Pages"] = convErr
			return
		}
		delete(t.editErrorMap, "Pages")
		parsedData.Pages = uint16(pages)
	})
	form.AddInputField("Language:", parsedData.Language, 0, nil, func(text string) {
		parsedData.Language = text
	})
	form.AddInputField("Publisher:", parsedData.Publisher, 0, nil, func(text string) {
		parsedData.Publisher = text
		bookFileNameInputField.SetText(t.parsedData.GetBookFileName())
	})
	form.AddInputField("PublisherURL:", parsedData.PublisherURL, 0, nil, func(text string) {
		parsedData.PublisherURL = text
	})
	form.AddInputField("Edition:", strconv.FormatUint(uint64(parsedData.Edition), 10), 0, nil, func(text string) {
		edition, convErr := strconv.Atoi(text)
		if convErr != nil {
			t.editErrorMap["Edition"] = convErr
			return
		}
		delete(t.editErrorMap, "Edition")
		parsedData.Edition = uint8(edition)
		bookFileNameInputField.SetText(t.parsedData.GetBookFileName())
	})
	form.AddInputField("PubDate:", parsedData.PubDate.Format(dateLayout), 0, nil, func(text string) {
		parsedDate, dateErr := time.Parse(dateLayout, text)
		if dateErr != nil {
			t.editErrorMap["PubDate"] = dateErr
			return
		}
		delete(t.editErrorMap, "PubDate")
		parsedData.PubDate = parsedDate
		bookFileNameInputField.SetText(t.parsedData.GetBookFileName())
	})
	form.AddInputField("Authors:", strings.Join(parsedData.Authors, ";"), 0, nil, func(text string) {
		newValues := getNewSliceData(text)
		if len(newValues) == 0 {
			t.editErrorMap["Authors"] = fmt.Errorf("authors slice length is 0")
			return
		}
		delete(t.editErrorMap, "Authors")
		parsedData.Authors = newValues
	})
	form.AddInputField("Categories:", strings.Join(parsedData.Categories, ";"), 0, nil, func(text string) {
		newValues := getNewSliceData(text)
		if len(newValues) == 0 {
			t.editErrorMap["Categories"] = fmt.Errorf("categories slice length is 0")
			return
		}
		delete(t.editErrorMap, "Categories")
		parsedData.Categories = newValues
	})
	form.AddInputField("Tags:", strings.Join(parsedData.Tags, ";"), 0, nil, func(text string) {
		newValues := getNewSliceData(text)
		if len(newValues) == 0 {
			t.editErrorMap["Tags"] = fmt.Errorf("tags slice length is 0")
			return
		}
		delete(t.editErrorMap, "Tags")
		parsedData.Tags = newValues
	})
	form.AddInputField("Formats:", strings.Join(parsedData.Formats, ";"), 0, nil, func(text string) {
		newValues := getNewSliceData(text)
		if len(newValues) == 0 {
			t.editErrorMap["Formats"] = fmt.Errorf("formats slice length is 0")
			return
		}
		delete(t.editErrorMap, "Formats")
		parsedData.Formats = newValues
	})
	form.AddInputField("BookFileName:", parsedData.BookFileName, 0, nil, func(text string) {
		parsedData.BookFileName = text
	})
	form.AddInputField("BookFileSize:", strconv.Itoa(int(parsedData.BookFileSize)), 0, nil, func(text string) {
		size, convErr := strconv.Atoi(text)
		if convErr != nil {
			t.editErrorMap["BookFileSize"] = convErr
			return
		}
		delete(t.editErrorMap, "BookFileSize")
		parsedData.BookFileSize = int64(size)
	})
	form.AddInputField("CoverFileName:", parsedData.CoverFileName, 0, nil, func(text string) {
		parsedData.CoverFileName = text
	})
	form.AddInputField("CoverURL:", parsedData.CoverURL, 0, nil, func(text string) {
		parsedData.CoverURL = text
	})

	form.AddButton("Add", func() {
		t.validateAuthors(parsedData)
		t.validateCategories(parsedData)
		t.ignoreExistingData = true // ignore existing, force to create a new record
		t.saveBook()
	})
	if t.existingData != nil {
		form.AddButton("Update", func() {
			t.validateAuthors(parsedData)
			t.validateCategories(parsedData)
			t.saveBook()
		})
		form.SetFocus(21) // Update Button
	} else {
		form.SetFocus(20) // Add Button
	}
	form.AddButton("Quit", func() {
		t.tuiApp.Stop()
	})
	form.SetButtonsAlign(tview.AlignCenter)

	bookFileNameInputField = form.GetFormItemByLabel("BookFileName:").(*tview.InputField)
}

func (t *TuiApp) validateAuthors(parsedData *book.ParsedData) {
	if strings.Contains(strings.Join(parsedData.Authors, ";"), "author") {
		t.editErrorMap["AuthorName"] = fmt.Errorf("the 'author' word should not be present")
	} else if isContainOneWordItem(parsedData.Authors) {
		t.editErrorMap["AuthorName"] = fmt.Errorf("an author name should contain at least 2 words")
	} else {
		delete(t.editErrorMap, "AuthorName")
	}
}

func (t *TuiApp) validateCategories(parsedData *book.ParsedData) {
	if len(parsedData.Categories) == 0 {
		t.editErrorMap["Categories"] = fmt.Errorf("there should be at least one category")
	} else {
		delete(t.editErrorMap, "Categories")
	}
}

func (t *TuiApp) saveBook() {
	if len(t.editErrorMap) != 0 {
		t.footer.SetText(getErrorText(t.editErrorMap))
		return
	}

	if t.ignoreExistingData {
		t.StoreBook(t.parsedData, nil, t.tempFilesData)
	} else {
		t.StoreBook(t.parsedData, t.existingData, t.tempFilesData)
	}

	if t.existingData == nil || t.ignoreExistingData {
		t.footer.SetText("The book is added successfully").SetTextColor(tcell.ColorGreen)
	} else {
		t.footer.SetText("The book is updated successfully").SetTextColor(tcell.ColorYellow)
	}

	t.restartFlow(true)
}

func (t *TuiApp) fillCheckboxTable(table *tview.Table, parsedData *book.ParsedData, existingData *book.StoredData) {
	table.SetCell(1, 0, equalCell(parsedData.Title, existingData.Title))
	table.SetCell(2, 0, equalCell(parsedData.Subtitle, existingData.Subtitle))
	table.SetCell(3, 0, equalCell(parsedData.Description, existingData.Description))
	table.SetCell(4, 0, equalCell(parsedData.ISBN10, existingData.ISBN10))
	table.SetCell(5, 0, equalCell(parsedData.ISBN13, existingData.ISBN13))
	table.SetCell(6, 0, equalCell(parsedData.ASIN, existingData.ASIN))
	table.SetCell(7, 0, equalCell(parsedData.Pages, existingData.Pages))
	table.SetCell(8, 0, equalCell(parsedData.Language, existingData.Language))
	table.SetCell(9, 0, equalCell(parsedData.Publisher, existingData.Publisher))
	table.SetCell(10, 0, equalCell(parsedData.PublisherURL, existingData.PublisherURL))
	table.SetCell(11, 0, equalCell(parsedData.Edition, existingData.Edition))
	table.SetCell(12, 0, equalCell(parsedData.PubDate.Format(dateLayout), existingData.PubDate.Format(dateLayout)))
	table.SetCell(13, 0, equalCell(strings.Join(parsedData.Authors, ";"), strings.Join(existingData.Authors, ";")))
	table.SetCell(14, 0,
		equalCell(strings.Join(parsedData.Categories, ";"), strings.Join(existingData.Categories, ";")))
	table.SetCell(15, 0, equalCell(strings.Join(parsedData.Tags, ";"), strings.Join(existingData.Tags, ";")))
	table.SetCell(16, 0, equalCell(strings.Join(parsedData.Formats, ";"), strings.Join(existingData.Formats, ";")))
	table.SetCell(17, 0, equalCell(parsedData.BookFileName, existingData.BookFileName))
	table.SetCell(18, 0, equalCell(parsedData.BookFileSize, existingData.BookFileSize))
	table.SetCell(19, 0, equalCell(parsedData.CoverFileName, existingData.CoverFileName))
}

func (t *TuiApp) fillExisingTable(table *tview.Table, parsedData *book.ParsedData, existingData *book.StoredData) {
	table.SetCell(1, 0, tview.NewTableCell(existingData.Title).
		SetTextColor(equalColor(parsedData.Title, existingData.Title)).
		SetAlign(tview.AlignLeft))
	table.SetCell(2, 0, tview.NewTableCell(existingData.Subtitle).
		SetTextColor(equalColor(parsedData.Subtitle, existingData.Subtitle)).
		SetAlign(tview.AlignLeft))
	table.SetCell(3, 0, tview.NewTableCell(existingData.Description).
		SetTextColor(equalColor(parsedData.Description, existingData.Description)).
		SetAlign(tview.AlignLeft))
	table.SetCell(4, 0, tview.NewTableCell(existingData.ISBN10).
		SetTextColor(equalColor(parsedData.ISBN10, existingData.ISBN10)).
		SetAlign(tview.AlignLeft))
	table.SetCell(5, 0, tview.NewTableCell(strconv.FormatInt(existingData.ISBN13, 10)).
		SetTextColor(equalColor(parsedData.ISBN13, existingData.ISBN13)).
		SetAlign(tview.AlignLeft))
	table.SetCell(6, 0, tview.NewTableCell(existingData.ASIN).
		SetTextColor(equalColor(parsedData.ASIN, existingData.ASIN)).
		SetAlign(tview.AlignLeft))
	table.SetCell(7, 0, tview.NewTableCell(strconv.FormatUint(uint64(existingData.Pages), 10)).
		SetTextColor(equalColor(parsedData.Pages, existingData.Pages)).
		SetAlign(tview.AlignLeft))
	table.SetCell(8, 0, tview.NewTableCell(existingData.Language).
		SetTextColor(equalColor(parsedData.Language, existingData.Language)).
		SetAlign(tview.AlignLeft))
	table.SetCell(9, 0, tview.NewTableCell(existingData.Publisher).
		SetTextColor(equalColor(parsedData.Publisher, existingData.Publisher)).
		SetAlign(tview.AlignLeft))
	table.SetCell(10, 0, tview.NewTableCell(existingData.PublisherURL).
		SetTextColor(equalColor(parsedData.PublisherURL, existingData.PublisherURL)).
		SetAlign(tview.AlignLeft))
	table.SetCell(11, 0, tview.NewTableCell(strconv.FormatUint(uint64(existingData.Edition), 10)).
		SetTextColor(equalColor(parsedData.Edition, existingData.Edition)).
		SetAlign(tview.AlignLeft))
	table.SetCell(12, 0, tview.NewTableCell(existingData.PubDate.Format(dateLayout)).
		SetTextColor(equalColor(parsedData.PubDate.Format(dateLayout), existingData.PubDate.Format(dateLayout))).
		SetAlign(tview.AlignLeft))
	existingAuthors := strings.Join(existingData.Authors, ";")
	parsedAuthors := strings.Join(parsedData.Authors, ";")
	table.SetCell(13, 0, tview.NewTableCell(existingAuthors).
		SetTextColor(equalColor(parsedAuthors, existingAuthors)).
		SetAlign(tview.AlignLeft))
	existingCategories := strings.Join(existingData.Categories, ";")
	parsedCategories := strings.Join(parsedData.Categories, ";")
	table.SetCell(14, 0, tview.NewTableCell(existingCategories).
		SetTextColor(equalColor(parsedCategories, existingCategories)).
		SetAlign(tview.AlignLeft))
	existingTags := strings.Join(existingData.Tags, ";")
	parsedTags := strings.Join(parsedData.Tags, ";")
	table.SetCell(15, 0, tview.NewTableCell(existingTags).
		SetTextColor(equalColor(parsedTags, existingTags)).
		SetAlign(tview.AlignLeft))
	existingFormats := strings.Join(existingData.Formats, ";")
	parsedFormats := strings.Join(parsedData.Formats, ";")
	table.SetCell(16, 0, tview.NewTableCell(existingFormats).
		SetTextColor(equalColor(parsedFormats, existingFormats)).
		SetAlign(tview.AlignLeft))
	table.SetCell(17, 0, tview.NewTableCell(existingData.BookFileName).
		SetTextColor(equalColor(parsedData.BookFileName, existingData.BookFileName)).
		SetAlign(tview.AlignLeft))
	table.SetCell(18, 0, tview.NewTableCell(strconv.FormatInt(existingData.BookFileSize, 10)).
		SetTextColor(equalColor(parsedData.BookFileSize, existingData.BookFileSize)).
		SetAlign(tview.AlignLeft))
	table.SetCell(19, 0, tview.NewTableCell(existingData.CoverFileName).
		SetTextColor(equalColor(parsedData.CoverFileName, existingData.CoverFileName)).
		SetAlign(tview.AlignLeft))
}

func (t *TuiApp) appendFooterText(text string) {
	footerText := t.footer.GetText(false)
	if len(strings.TrimSpace(footerText)) == 0 {
		t.footer.SetText(text)
	} else {
		t.footer.SetText(footerText + text)
	}
}

func (t *TuiApp) restartFlow(clearForms bool) {
	t.parsedData = nil
	t.existingData = nil
	t.tempFilesData = nil
	t.ignoreExistingData = false

	if clearForms {
		t.clearForms()
	}

	t.bookIDInput.SetText("")
	t.tuiApp.SetFocus(t.bookIDInput)
}

func (t *TuiApp) clearForms() {
	t.parsedForm.Clear(true)
	t.existingTable.Clear()
	t.equalityTable.Clear()
}

func getErrorText(errorMap map[string]error) string {
	builder := strings.Builder{}
	for key, val := range errorMap {
		builder.WriteString(fmt.Sprintf("%s: %v\n", key, val))
	}

	return builder.String()
}

func getNewSliceData(text string) []string {
	stringSplit := strings.Split(text, ";")
	newValues := make([]string, 0)
	for _, value := range stringSplit {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			newValues = append(newValues, trimmed)
		}
	}

	return newValues
}

func equalColor[T comparable](a, b T) tcell.Color {
	if a == b {
		return tcell.ColorGreen
	}

	return tcell.ColorRed
}

func equalCell[T comparable](a, b T) *tview.TableCell {
	if a == b {
		return tview.NewTableCell("=").SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignCenter)
	}

	return tview.NewTableCell("\u2260").SetTextColor(tcell.ColorRed).SetAlign(tview.AlignCenter)
}

func isContainOneWordItem(items []string) bool {
	for _, item := range items {
		splitItem := strings.Split(item, " ")
		if len(splitItem) < 2 {
			return true
		}
	}

	return false
}
