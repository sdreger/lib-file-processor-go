package parser

import (
	"testing"
	"time"
)

func TestParsePublisherString(t *testing.T) {
	tests := []struct {
		input             string
		meta              BookPublishMeta
		shouldReturnError bool
	}{
		{
			input: "Wiley; 1st edition (October 16, 2017)",
			meta: BookPublishMeta{
				Publisher: "Wiley",
				Edition:   1,
				PubDate:   time.Date(2017, 10, 16, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: "No Starch Press; 2nd edition (May 3, 2019)",
			meta: BookPublishMeta{
				Publisher: "No Starch Press",
				Edition:   2,
				PubDate:   time.Date(2019, 5, 3, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: "No Starch Press (November 5, 2020)",
			meta: BookPublishMeta{
				Publisher: "No Starch Press",
				Edition:   1,
				PubDate:   time.Date(2020, 11, 5, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: "Apress; 1st ed. edition (November 1, 2020)",
			meta: BookPublishMeta{
				Publisher: "Apress",
				Edition:   1,
				PubDate:   time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: "Packt Publishing; 3rd edition (17 May 2021)",
			meta: BookPublishMeta{
				Publisher: "Packt Publishing",
				Edition:   3,
				PubDate:   time.Date(2021, 5, 17, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: "Pearson; 3rd edition (2 Jun. 2022)",
			meta: BookPublishMeta{
				Publisher: "Pearson",
				Edition:   3,
				PubDate:   time.Date(2022, 6, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: "No Starch Press (25 May 2019)",
			meta: BookPublishMeta{
				Publisher: "No Starch Press",
				Edition:   1,
				PubDate:   time.Date(2019, 5, 25, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: "No Starch Press (1 Oct. 2020)",
			meta: BookPublishMeta{
				Publisher: "No Starch Press",
				Edition:   1,
				PubDate:   time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: "DK Children; Workbook edition (March 7, 2017)",
			meta: BookPublishMeta{
				Publisher: "DK Children",
				Edition:   1,
				PubDate:   time.Date(2017, 3, 7, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: "DK Publishing (Dorling Kindersley); Workbook edition (7 Mar. 2017)",
			meta: BookPublishMeta{
				Publisher: "DK Publishing (Dorling Kindersley)",
				Edition:   1,
				PubDate:   time.Date(2017, 3, 7, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: "Esri Press; Fourth edition (December 28, 2021)",
			meta: BookPublishMeta{
				Publisher: "Esri Press",
				Edition:   4,
				PubDate:   time.Date(2021, 12, 28, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: "Esri Press; Fourth edition (10 Feb. 2022)",
			meta: BookPublishMeta{
				Publisher: "Esri Press",
				Edition:   4,
				PubDate:   time.Date(2022, 2, 10, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: "Springer; 2nd ed. 2023 edition",
			meta: BookPublishMeta{
				Publisher: "Springer",
				Edition:   2,
			},
			shouldReturnError: true,
		},
		{
			input: "Packt Publishing",
			meta: BookPublishMeta{
				Publisher: "Packt Publishing",
				Edition:   1,
			},
			shouldReturnError: true,
		},
		{
			input: "Unknown; 2nd edition (Unknown 15, 2021)",
			meta: BookPublishMeta{
				Publisher: "Unknown",
				Edition:   2,
			},
			shouldReturnError: true,
		},
		{
			input:             "2nd edition: Unknown (15 June 2021)",
			meta:              BookPublishMeta{},
			shouldReturnError: true,
		},
	}

	t.Log("Given the need to test publisher string parsing.")
	for i, tt := range tests {
		t.Logf("\tTest: %d\tWhen checking %q for publisher metadata: %v\n", i, tt.input, tt.meta)
		meta, err := ParsePublisherString(tt.input)
		if !tt.shouldReturnError && err != nil {
			t.Fatalf("\t\t%s\tShould be able to parse publisher string: %v", failed, err)
		}
		if meta != tt.meta {
			t.Errorf("\t\t%s\tShould get a %v publisher metadata: %v", failed, tt.meta, meta)
		} else {
			t.Logf("\t\t%s\tShould be able to get correct publisher metadata.", succeed)
		}
	}
}
