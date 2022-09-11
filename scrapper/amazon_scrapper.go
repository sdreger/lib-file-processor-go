package scrapper

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	cookiejar "github.com/juju/persistent-cookiejar"
	"github.com/sdreger/lib-file-processor-go/domain/book"
	"github.com/sdreger/lib-file-processor-go/domain/publisher"
	"github.com/sdreger/lib-file-processor-go/parser"
	"log"
	"strconv"
	"strings"
	"unicode"
)

const (
	defaultBasePath = "https://www.amazon.com/dp/"

	//userAgentSafari = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Safari/605.1.15"
	userAgentEdge = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.134 Safari/537.36 Edg/103.0.1264.77"

	publisherKey         = "Publisher"
	editionKey           = "Edition"
	pubDateKey           = "Publication date"
	pagesKey             = "Print pages"
	printLengthKey       = "Print length"
	hardcoverKey         = "Hardcover"
	paperbackKey         = "Paperback"
	printedAccessCodeKey = "Printed Access Code"
	sourceISBNKey        = "Page numbers source ISBN"
	ISBN10Key            = "ISBN-10"
	ISBN13Key            = "ISBN-13"
	ASINKey              = "ASIN"
	languageKey          = "Language"

	bookTitleSelector       = `span[id=productTitle]`
	bookSubtitleSelector    = `span[id=productSubtitle]`
	bookDescriptionSelector = `div[id=bookDescription_feature_div]>div[data-a-expander-name=book_description_expander]>div.a-expander-content`
	bookCategoriesSelector  = `div[id=wayfinding-breadcrumbs_feature_div]>ul>li>span>a`
	bookAuthorsSelector     = `div[id=bylineInfo]>span.author>span.a-declarative>a,div[id=bylineInfo]>span.author>a`
	bookDetailsSelector     = `div[id=detailBullets_feature_div]>ul>li`
	bookCarouserSelector    = `li.rpi-carousel-attribute-card>div.rpi-attribute-content`
	bookISBNBlockSelector   = `div[id=isbn_feature_div]>div.a-section>div.a-row, div[id=printEditionIsbn_feature_div]>div.a-section>div.a-row`
	bookCoverURLSelector    = `img[id=imgBlkFront], img[id=ebooksImgBlkFront]`
)

type AmazonScrapper struct {
	basePath        string
	cookieJar       *cookiejar.Jar
	collector       *colly.Collector
	scrappedRawData *scrappedRawData
	logger          *log.Logger
}

func NewAmazonScrapper(basePath string, logger *log.Logger) (*AmazonScrapper, error) {
	if basePath == "" || !(strings.HasPrefix(basePath, "http") || strings.HasPrefix(basePath, "file")) {
		basePath = defaultBasePath
	}

	cookieJar, err := cookiejar.New(&cookiejar.Options{Filename: "cookie.db"})
	if err != nil {
		return nil, err
	}

	collector := colly.NewCollector()
	collector.AllowURLRevisit = true
	//extensions.RandomUserAgent(collector)
	collector.UserAgent = userAgentEdge
	collector.SetCookieJar(cookieJar)

	scrappedRawData := newScrappedRawData()
	initCallbacks(collector, &scrappedRawData, logger)

	return &AmazonScrapper{
		basePath:        basePath,
		cookieJar:       cookieJar,
		collector:       collector,
		scrappedRawData: &scrappedRawData,
		logger:          logger,
	}, nil
}

func initCallbacks(collector *colly.Collector, scrappedRawData *scrappedRawData, logger *log.Logger) {

	collector.OnRequest(func(request *colly.Request) {
		logger.Printf("[INFO] - Visiting: %q, using 'User-Agent': %q", request.URL, request.Headers.Get("User-Agent"))
		request.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		request.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		request.Headers.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8,uk;q=0.7")
		request.Headers.Set("Cache-Control", "max-age=0")
		request.Headers.Set("Device-Memory", "8")
		request.Headers.Set("Downlink", "10")
		request.Headers.Set("Dpr", "1")
		request.Headers.Set("Ect", "4g")
		request.Headers.Set("Rtt", "50")
		request.Headers.Set("Upgrade-Insecure-Requests", "1")
		request.Headers.Set("Viewport-Width", "1920")
	})

	collector.OnHTML(bookTitleSelector, func(element *colly.HTMLElement) {
		scrappedRawData.titleString = strings.TrimSpace(element.Text)
	})

	collector.OnHTML(bookSubtitleSelector, func(element *colly.HTMLElement) {
		scrappedRawData.subtitleString = strings.TrimSpace(element.Text)
	})

	collector.OnHTML(bookDescriptionSelector, func(element *colly.HTMLElement) {
		//element.DOM.Find("span.a-text-italic").RemoveClass("a-text-italic").AddClass("text-italic")
		//element.DOM.Find("span.a-text-bold").RemoveClass("a-text-bold").AddClass("text-bold")
		html, _ := element.DOM.Html()
		scrappedRawData.description = strings.ReplaceAll(strings.TrimSpace(html), "<p></p>", "")
	})

	collector.OnHTML(bookCategoriesSelector, func(element *colly.HTMLElement) {
		text := strings.TrimSpace(element.Text)
		if text == "Books" || text == "Kindle Store" || text == "Kindle eBooks" || text == "New, Used & Rental Textbooks" {
			return
		}
		scrappedRawData.categories = append(scrappedRawData.categories, text)
	})

	collector.OnHTML(bookAuthorsSelector, func(element *colly.HTMLElement) {
		if element.Attr("role") == "button" || element.Text == "" {
			return
		}
		scrappedRawData.authors = append(scrappedRawData.authors, strings.TrimSpace(element.Text))
	})

	collector.OnHTML(bookDetailsSelector, func(element *colly.HTMLElement) {
		liText := element.Text
		liParts := strings.Split(liText, ":")
		if len(liParts) < 2 {
			return
		}
		key := strings.TrimSpace(strings.Map(removeNonPrintable, liParts[0]))
		value := strings.TrimSpace(strings.Map(removeNonPrintable, liParts[1]))
		scrappedRawData.detailsBlock[key] = value
	})

	collector.OnHTML(bookCarouserSelector, func(element *colly.HTMLElement) {
		key := strings.TrimSpace(strings.Map(removeNonPrintable, element.ChildText("div.rpi-attribute-label")))
		value := strings.TrimSpace(strings.Map(removeNonPrintable, element.ChildText("div.rpi-attribute-value")))
		if key != "" {
			scrappedRawData.detailsCarousel[key] = value
		}
	})

	collector.OnHTML(bookISBNBlockSelector, func(element *colly.HTMLElement) {
		element.ChildText("span")
		var key, value string
		for i, node := range element.DOM.Children().Nodes {
			if i == 0 {
				key = strings.Trim(node.LastChild.Data, ":")
			}
			value = strings.TrimSpace(node.LastChild.Data)
		}
		if key == ISBN10Key {
			scrappedRawData.ISBN10String = value
		}
		if key == ISBN13Key {
			scrappedRawData.ISBN13String = value
		}
	})

	collector.OnHTML(bookCoverURLSelector, func(element *colly.HTMLElement) {
		scrappedRawData.coverURL = element.Attr("src")
	})
}

func (s *AmazonScrapper) Close() error {
	return s.cookieJar.Save()
}

func (s *AmazonScrapper) GetBookData(bookID string) (book.ParsedData, error) {
	s.scrappedRawData.init()
	err := s.collector.Visit(s.basePath + bookID)
	if err != nil {
		return book.ParsedData{}, err
	}

	authors := s.scrappedRawData.authors
	categories := s.scrappedRawData.categories
	detailsBlock := s.scrappedRawData.detailsBlock
	detailsCarousel := s.scrappedRawData.detailsCarousel
	titleString := s.scrappedRawData.titleString
	subtitleString := s.scrappedRawData.subtitleString
	description := s.scrappedRawData.description
	ISBN10String := s.scrappedRawData.ISBN10String
	ISBN13String := s.scrappedRawData.ISBN13String
	coverURL := s.scrappedRawData.coverURL

	// -------------------- Book title / subtitle --------------------
	title, subtitle := parser.ParseTitleString(titleString)

	// -------------------- Book publisher metadata --------------------
	publishMeta, err := parser.ParsePublisherString(detailsBlock[publisherKey])
	if err != nil {
		s.logger.Printf("[WARN] - %v", err)
	}

	// -------------------- Book ISBN10 --------------------
	isbn10 := getISBN10(ISBN10String, detailsBlock[ISBN10Key], detailsBlock[sourceISBNKey])

	// -------------------- Book ISBN13 --------------------
	isbn13, err := getISBN13(ISBN13String, detailsBlock[ISBN13Key])
	if err != nil {
		return book.ParsedData{}, fmt.Errorf("can not get ISBN13 value: %w", err)
	}

	metadata := book.ParsedData{
		Title:         title,
		Subtitle:      subtitle,
		Description:   description,
		ISBN10:        isbn10,
		ISBN13:        int64(isbn13),
		ASIN:          detailsBlock[ASINKey],
		Pages:         getBookLength(detailsBlock),
		Language:      detailsBlock[languageKey],
		PublisherURL:  "",
		Publisher:     publisher.MapPublisherName(publishMeta.Publisher),
		Edition:       getBookEdition(titleString, subtitleString, publishMeta.Edition),
		PubDate:       publishMeta.PubDate,
		Authors:       authors,
		Categories:    categories,
		Tags:          nil,
		Formats:       nil,
		BookFileSize:  0,
		BookFileName:  "",
		CoverFileName: "",
		CoverURL:      coverURL,
	}

	enrichWithOptionalCarouselData(&metadata, detailsCarousel)
	metadata.PublisherURL = getPublisherURL(s.basePath, metadata.ISBN10, metadata.ASIN)
	primaryBookId := metadata.GetPrimaryId()
	metadata.CoverFileName = fmt.Sprint(primaryBookId, getCoverExtension(metadata.CoverURL))
	metadata.BookFileName = metadata.GetBookFileName()

	return metadata, nil
}

func enrichWithOptionalCarouselData(parsedData *book.ParsedData, detailsCarousel map[string]string) {
	if len(detailsCarousel) == 0 {
		return
	}
	if parsedData.Publisher == "" {
		parsedData.Publisher = publisher.MapPublisherName(detailsCarousel[publisherKey])
	}
	if parsedData.Edition == 0 {
		edition, err := parser.ParseEditionString(detailsCarousel[editionKey] + " Edition")
		if err != nil {
			parsedData.Edition = edition
		} else {
			parsedData.Edition = 1
		}
	}
	if parsedData.PubDate.IsZero() {
		pubDate, err := parser.ParseDateString(detailsCarousel[pubDateKey])
		if err == nil {
			parsedData.PubDate = pubDate
		}
	}
}

func getBookEdition(titleString, subtitleString string, publisherEdition uint8) uint8 {
	titleEdition, err := parser.ParseEditionString(titleString)
	if err == nil && titleEdition > 0 {
		return titleEdition
	}

	subtitleEdition, err := parser.ParseEditionString(subtitleString)
	if err == nil && subtitleEdition > 0 {
		return subtitleEdition
	}

	return publisherEdition
}

func getISBN10(ISBN10, metaISBN, metaASIN string) string {
	var effectiveISBN10 string
	if ISBN10 != "" {
		effectiveISBN10 = ISBN10
	}
	if metaISBN != "" {
		effectiveISBN10 = metaISBN
	}
	if metaASIN != "" {
		effectiveISBN10 = metaASIN
	}

	effectiveISBN10 = strings.ReplaceAll(effectiveISBN10, "-", "")
	if len(effectiveISBN10) == 10 {
		return effectiveISBN10
	}

	return ""
}

func getISBN13(ISBN13, metaISBN13 string) (int, error) {
	var effectiveISBN13 string
	if ISBN13 != "" {
		effectiveISBN13 = ISBN13
	}
	if metaISBN13 != "" {
		effectiveISBN13 = metaISBN13
	}

	effectiveISBN13 = strings.ReplaceAll(effectiveISBN13, "-", "")
	if len(effectiveISBN13) == 13 {
		return strconv.Atoi(effectiveISBN13)
	}

	return 0, nil
}

func getPublisherURL(basePath, ISBN10, ASIN string) string {
	if ISBN10 != "" {
		return basePath + ISBN10
	}

	return basePath + ASIN
}

func getBookLength(detailsBlock map[string]string) uint16 {
	pagesString, paperbackString, hardcoverString, printLength, printedAccessCode :=
		detailsBlock[pagesKey], detailsBlock[paperbackKey],
		detailsBlock[hardcoverKey], detailsBlock[printLengthKey],
		detailsBlock[printedAccessCodeKey]
	if pagesString != "" {
		return parser.ParseLengthString(pagesString)
	}
	if paperbackString != "" {
		return parser.ParseLengthString(paperbackString)
	}
	if hardcoverString != "" {
		return parser.ParseLengthString(hardcoverString)
	}
	if printLength != "" {
		return parser.ParseLengthString(printLength)
	}
	if printedAccessCode != "" {
		return parser.ParseLengthString(printedAccessCode)
	}

	return 0
}

func getCoverExtension(coverURL string) string {
	lastDotIndex := strings.LastIndex(coverURL, ".")
	if lastDotIndex == -1 {
		return ""
	}

	return coverURL[lastDotIndex:]
}

func removeNonPrintable(r rune) rune {
	if unicode.IsPrint(r) {
		return r
	}
	return -1
}
