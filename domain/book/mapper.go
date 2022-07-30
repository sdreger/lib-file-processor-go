package book

// mapToStoredData map an SQL dot product row to a StoredData struct. Collection fields are appended.
func mapToStoredData(rowData dotProductRow, storedData *StoredData) {
	if storedData.ID == 0 {
		storedData.ID = rowData.ID
	}
	if storedData.Title == "" {
		storedData.Title = rowData.Title
	}
	if storedData.Subtitle == "" && rowData.Subtitle.Valid {
		storedData.Subtitle = rowData.Subtitle.String
	}
	if storedData.Description == "" {
		storedData.Description = rowData.Description
	}
	if storedData.ISBN10 == "" && rowData.ISBN10.Valid {
		storedData.ISBN10 = rowData.ISBN10.String
	}
	if storedData.ISBN13 == 0 && rowData.ISBN13.Valid {
		storedData.ISBN13 = rowData.ISBN13.Int64
	}
	if storedData.ASIN == "" && rowData.ASIN.Valid {
		storedData.ASIN = rowData.ASIN.String
	}
	if storedData.Pages == 0 {
		storedData.Pages = rowData.Pages
	}
	if storedData.Language == "" {
		storedData.Language = rowData.Language
	}
	if storedData.Publisher == "" {
		storedData.Publisher = rowData.Publisher
	}
	if storedData.PublisherURL == "" {
		storedData.PublisherURL = rowData.PublisherURL
	}
	if storedData.Edition == 0 {
		storedData.Edition = rowData.Edition
	}
	if storedData.PubDate.IsZero() {
		storedData.PubDate = rowData.PubDate
	}
	if storedData.BookFileName == "" {
		storedData.BookFileName = rowData.BookFileName
	}
	if storedData.BookFileSize == 0 {
		storedData.BookFileSize = rowData.BookFileSize
	}
	if storedData.CoverFileName == "" {
		storedData.CoverFileName = rowData.CoverFileName
	}
	if storedData.CreatedAt.IsZero() {
		storedData.CreatedAt = rowData.CreatedAt
	}
	if storedData.UpdatedAt.IsZero() {
		storedData.UpdatedAt = rowData.UpdatedAt
	}
	if rowData.AuthorName.Valid {
		storedData.Authors = append(storedData.Authors, rowData.AuthorName.String)
	}
	if rowData.CategoryName.Valid {
		storedData.Categories = append(storedData.Categories, rowData.CategoryName.String)
	}
	if rowData.FileTypeName.Valid {
		storedData.Formats = append(storedData.Formats, rowData.FileTypeName.String)
	}
	if rowData.TagName.Valid {
		storedData.Tags = append(storedData.Tags, rowData.TagName.String)
	}
}

// mapToParsedDate map a StoredData struct to ParsedData struct. Updates empty non-reference fields only.
func mapToParsedDate(existingData *StoredData, parsedData *ParsedData) {
	if parsedData.Title == "" {
		parsedData.Title = existingData.Title
	}
	if parsedData.Subtitle == "" {
		parsedData.Subtitle = existingData.Subtitle
	}
	if parsedData.Description == "" {
		parsedData.Description = existingData.Description
	}
	if parsedData.ISBN10 == "" {
		parsedData.ISBN10 = existingData.ISBN10
	}
	if parsedData.ISBN13 == 0 {
		parsedData.ISBN13 = existingData.ISBN13
	}
	if parsedData.ASIN == "" {
		parsedData.ASIN = existingData.ASIN
	}
	if parsedData.Pages == 0 {
		parsedData.Pages = existingData.Pages
	}
	if parsedData.PublisherURL == "" {
		parsedData.PublisherURL = existingData.PublisherURL
	}
	if parsedData.Edition == 0 {
		parsedData.Edition = existingData.Edition
	}
	if parsedData.PubDate.IsZero() {
		parsedData.PubDate = existingData.PubDate
	}
	if parsedData.BookFileName == "" {
		parsedData.BookFileName = existingData.BookFileName
	}
	if parsedData.BookFileSize == 0 {
		parsedData.BookFileSize = existingData.BookFileSize
	}
	if parsedData.CoverFileName == "" {
		parsedData.CoverFileName = existingData.CoverFileName
	}
}

// deduplicateMappedData removes duplicates from a string slice.
func deduplicateMappedData(slice []string) (result []string) {
	uniqueMap := make(map[string]bool)
	for _, entry := range slice {
		if _, ok := uniqueMap[entry]; !ok {
			uniqueMap[entry] = true
			result = append(result, entry)
		}
	}

	return
}
