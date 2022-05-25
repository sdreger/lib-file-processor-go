package parser

import (
	"testing"
)

const (
	succeed = "\u2713"
	failed  = "\u2717"
)

func TestParseTitle(t *testing.T) {
	tests := []struct {
		input    string
		title    string
		subTitle string
	}{
		{
			input:    "Hands-on Kubernetes on Azure: Use Azure Kubernetes Service to automate management, scaling, and deployment of containerized applications, 3rd Edition",
			title:    "Hands-on Kubernetes on Azure",
			subTitle: "Use Azure Kubernetes Service to automate management, scaling, and deployment of containerized applications",
		},
		{
			input:    "Real-World Python: A Hacker's Guide to Solving Problems with Code",
			title:    "Real-World Python",
			subTitle: "A Hacker's Guide to Solving Problems with Code",
		},
		{
			input:    "Kubernetes in Action",
			title:    "Kubernetes in Action",
			subTitle: "",
		},
		{
			input:    "Kubernetes in Action, Second Edition",
			title:    "Kubernetes in Action",
			subTitle: "",
		},
		{
			input:    "The Official BBC micro:bit User Guide",
			title:    "The Official BBC micro:bit User Guide",
			subTitle: "",
		},
		{
			input:    "The Pragmatic Programmer: Your Journey To Mastery, 20th Anniversary Edition (2nd Edition)",
			title:    "The Pragmatic Programmer",
			subTitle: "Your Journey To Mastery, 20th Anniversary Edition",
		},
		{
			input:    "Python Crash Course, 2nd Edition: A Hands-On, Project-Based Introduction to Programming",
			title:    "Python Crash Course",
			subTitle: "A Hands-On, Project-Based Introduction to Programming",
		},
	}

	t.Log("Given the need to test title string parsing.")
	for i, tt := range tests {
		t.Logf("\tTest: %d\tWhen checking %q for title %q and subTitle %q\n", i, tt.input, tt.title, tt.subTitle)
		title, subTitle := ParseTitle(tt.input)
		if title != tt.title {
			t.Errorf("\t\t%s\tShould get a %q title: %v", failed, tt.title, title)
		}
		t.Logf("\t\t%s\tShould be able to get correct title.", succeed)
		if subTitle != tt.subTitle {
			t.Errorf("\t\t%s\tShould get a %q subTitle: %v", failed, tt.subTitle, subTitle)
		}
		t.Logf("\t\t%s\tShould be able to get correct subTitle.", succeed)
	}
}
