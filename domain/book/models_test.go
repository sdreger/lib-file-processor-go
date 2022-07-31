package book

import (
	"testing"
	"time"
)

func TestParsedData_GetBookFileName(t *testing.T) {
	tests := []struct {
		publisher   string
		title       string
		edition     uint8
		ISBN10      string
		ISBN13      int64
		ASIN        string
		publishDate time.Time
		fileName    string
	}{
		{
			publisher:   "NSP",
			title:       "Awesome Book",
			edition:     1,
			ISBN10:      "1234567890",
			ISBN13:      0,
			ASIN:        "",
			publishDate: testPublishDate,
			fileName:    "NSP.Awesome.Book.1234567890.Feb.2020.zip",
		},
		{
			publisher:   "For, Dummies!",
			title:       "What Is What?!",
			edition:     2,
			ISBN10:      "",
			ISBN13:      0,
			ASIN:        "BH128KL653",
			publishDate: testPublishDate,
			fileName:    "For.Dummies.What.Is.What.2nd.Edition.BH128KL653.Feb.2020.zip",
		},
		{
			publisher:   "DK",
			title:       "C#, For Beginners;",
			edition:     3,
			ISBN10:      "",
			ISBN13:      1234567890123,
			ASIN:        "",
			publishDate: testPublishDate,
			fileName:    "DK.C#.For.Beginners.3rd.Edition.1234567890123.Feb.2020.zip",
		},
		{
			publisher:   "Maker Media",
			title:       "C++ Data-Related Patterns. (Global Edition)",
			edition:     4,
			ISBN10:      "0987654321",
			ISBN13:      0,
			ASIN:        "",
			publishDate: testPublishDate,
			fileName:    "Maker.Media.C++.Data-Related.Patterns.Global.Edition.4th.Edition.0987654321.Feb.2020.zip",
		},
		{
			publisher:   "MK",
			title:       "PHP & MySQL",
			edition:     10,
			ISBN10:      "5432112345",
			ISBN13:      1234567890123,
			ASIN:        "BH128KL653",
			publishDate: testPublishDate,
			fileName:    "MK.PHP.and.MySQL.10th.Edition.5432112345.Feb.2020.zip",
		},
	}

	t.Log("Given the need to test book filename getter.")
	for i, tt := range tests {
		inputData := ParsedData{
			Publisher: tt.publisher,
			Title:     tt.title,
			Edition:   tt.edition,
			ISBN10:    tt.ISBN10,
			ISBN13:    tt.ISBN13,
			ASIN:      tt.ASIN,
			PubDate:   tt.publishDate,
		}
		t.Logf("\tTest: %d\tWhen checking %v for filename %s\n", i, inputData.GetPrimaryId(), tt.fileName)
		fileName := inputData.GetBookFileName()

		if fileName != tt.fileName {
			t.Errorf("\t\t%s\tShould get a %q date: %q", failed, tt.fileName, fileName)
		} else {
			t.Logf("\t\t%s\tShould be able to get correct filename value.", succeed)
		}
	}

	t.Logf("\t\t%s\tShould be able to get correct book filename", succeed)
}

func TestParsedData_GetBookFileNameWithoutExtension(t *testing.T) {
	parsedData := ParsedData{BookFileName: "NSP.Awesome.Book.1234567890.Feb.2020.zip"}
	expected := "NSP.Awesome.Book.1234567890.Feb.2020"
	nameWithoutExtension := parsedData.GetBookFileNameWithoutExtension()
	if nameWithoutExtension != expected {
		t.Errorf("\t\t%s\tShould get %s filename, got: %s", failed, expected, nameWithoutExtension)
	}

	t.Logf("\t\t%s\tShould be able to get book filename without extension", succeed)
}
