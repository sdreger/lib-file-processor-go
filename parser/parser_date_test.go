package parser

import (
	"testing"
	"time"
)

func TestParseDateString(t *testing.T) {
	tests := []struct {
		input             string
		date              time.Time
		shouldReturnError bool
	}{
		{
			input: "25 Oct. 2022",
			date:  time.Date(2022, 10, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			input: "October 25, 2022",
			date:  time.Date(2022, 10, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			input: "May 17, 2021",
			date:  time.Date(2021, 5, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			input: "17 May 2021",
			date:  time.Date(2021, 5, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			input:             "Sun, 28 Oct 2015",
			date:              time.Time{},
			shouldReturnError: true,
		},
	}

	t.Log("Given the need to test date string parsing.")
	for i, tt := range tests {
		t.Logf("\tTest: %d\tWhen checking %q for date %s\n", i, tt.input, tt.date)
		date, err := ParseDateString(tt.input)
		if !tt.shouldReturnError && err != nil {
			t.Fatalf("\t\t%s\tShould be able to get date value: %v", failed, err)
		}
		if date != tt.date {
			t.Errorf("\t\t%s\tShould get a %s date: %s", failed, tt.date, date)
		} else {
			t.Logf("\t\t%s\tShould be able to get correct date value.", succeed)
		}
	}
}
