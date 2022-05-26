package parser

import (
	"testing"
)

func TestParseTitle(t *testing.T) {
	tests := []struct {
		input    string
		title    string
		subtitle string
	}{
		{
			input:    "Hands-on Kubernetes on Azure: Use Azure Kubernetes Service to automate management, scaling, and deployment of containerized applications, 3rd Edition",
			title:    "Hands-on Kubernetes on Azure",
			subtitle: "Use Azure Kubernetes Service to automate management, scaling, and deployment of containerized applications",
		},
		{
			input:    "Real-World Python: A Hacker's Guide to Solving Problems with Code",
			title:    "Real-World Python",
			subtitle: "A Hacker's Guide to Solving Problems with Code",
		},
		{
			input:    "Kubernetes in Action",
			title:    "Kubernetes in Action",
			subtitle: "",
		},
		{
			input:    "Kubernetes in Action, Second Edition",
			title:    "Kubernetes in Action",
			subtitle: "",
		},
		{
			input:    "The Official BBC micro:bit User Guide",
			title:    "The Official BBC micro:bit User Guide",
			subtitle: "",
		},
		{
			input:    "The Pragmatic Programmer: Your Journey To Mastery, 20th Anniversary Edition (2nd Edition)",
			title:    "The Pragmatic Programmer",
			subtitle: "Your Journey To Mastery, 20th Anniversary Edition",
		},
		{
			input:    "Python Crash Course, 2nd Edition: A Hands-On, Project-Based Introduction to Programming",
			title:    "Python Crash Course",
			subtitle: "A Hands-On, Project-Based Introduction to Programming",
		},
	}

	t.Log("Given the need to test title string parsing.")
	for i, tt := range tests {
		t.Logf("\tTest: %d\tWhen checking %q for title %q and subtitle %q\n", i, tt.input, tt.title, tt.subtitle)
		title, subtitle := ParseTitleString(tt.input)
		if title != tt.title {
			t.Errorf("\t\t%s\tShould get a %q title: %v", failed, tt.title, title)
		}
		t.Logf("\t\t%s\tShould be able to get correct title.", succeed)
		if subtitle != tt.subtitle {
			t.Errorf("\t\t%s\tShould get a %q subtitle: %v", failed, tt.subtitle, subtitle)
		}
		t.Logf("\t\t%s\tShould be able to get correct subtitle.", succeed)
	}
}
