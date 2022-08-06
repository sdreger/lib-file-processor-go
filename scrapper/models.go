package scrapper

type scrappedRawData struct {
	authors         []string
	categories      []string
	tags            []string
	detailsBlock    map[string]string
	detailsCarousel map[string]string
	titleString     string
	subtitleString  string
	description     string
	ISBN10String    string
	ISBN13String    string
	coverURL        string
}

func newScrappedRawData() scrappedRawData {
	rawData := scrappedRawData{}
	rawData.init()
	return rawData
}

func (srd *scrappedRawData) init() {
	srd.authors = make([]string, 0)
	srd.categories = make([]string, 0)
	srd.tags = make([]string, 0)
	srd.detailsBlock = make(map[string]string)
	srd.detailsCarousel = make(map[string]string)
	srd.titleString = ""
	srd.subtitleString = ""
	srd.description = ""
	srd.ISBN10String = ""
	srd.ISBN13String = ""
	srd.coverURL = ""
}
