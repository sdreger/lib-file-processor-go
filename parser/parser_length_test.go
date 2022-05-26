package parser

import (
	"testing"
)

func TestParseLengthString(t *testing.T) {
	tests := []struct {
		input  string
		length uint16
	}{
		{
			input:  "544 pages",
			length: 544,
		},
		{
			input:  "",
			length: 0,
		},
		{
			input:  "ten pages",
			length: 0,
		},
	}

	t.Log("Given the need to test book length string parsing.")
	for i, tt := range tests {
		t.Logf("\tTest: %d\tWhen checking %q for book length %d\n", i, tt.input, tt.length)
		length := ParseLengthString(tt.input)
		if length != tt.length {
			t.Errorf("\t\t%s\tShould get a %d book length: %d", failed, tt.length, length)
		} else {
			t.Logf("\t\t%s\tShould be able to get correct book length value.", succeed)
		}
	}
}
