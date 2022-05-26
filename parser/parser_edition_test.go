package parser

import (
	"testing"
)

func TestParseEdition(t *testing.T) {
	tests := []struct {
		input   string
		edition uint8
	}{
		{
			input:   "Hands-on Kubernetes on Azure: Use Azure Kubernetes Service to automate management, scaling, and deployment of containerized applications, 3rd Edition",
			edition: 3,
		},
		{
			input:   "Real-World Python: A Hacker's Guide to Solving Problems with Code",
			edition: 0,
		},
		{
			input:   "1st edition",
			edition: 1,
		},
		{
			input:   " 2nd Edition ",
			edition: 2,
		},
		{
			input:   "Kubernetes in Action, Second Edition",
			edition: 2,
		},
		{
			input:   "The Pragmatic Programmer: Your Journey To Mastery, 20th Anniversary Edition (2nd Edition)",
			edition: 2,
		},
		{
			input:   "Python Crash Course, 2nd Edition: A Hands-On, Project-Based Introduction to Programming",
			edition: 2,
		},
	}

	t.Log("Given the need to test edition string parsing.")
	for i, tt := range tests {
		t.Logf("\tTest: %d\tWhen checking %q for edition %d\n", i, tt.input, tt.edition)
		edition, err := ParseEditionString(tt.input)
		if err != nil {
			t.Fatalf("\t\t%s\tShould be able to get edition value: %v", failed, err)
		}
		if edition != tt.edition {
			t.Errorf("\t\t%s\tShould get a %d edition: %d", failed, tt.edition, edition)
		}
		t.Logf("\t\t%s\tShould be able to get correct edition value.", succeed)
	}
}
