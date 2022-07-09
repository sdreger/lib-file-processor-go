package publisher

import "testing"

func TestMapPublisherName(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{input: "ACM Books", output: "ACM"},
		{input: "Academic Press", output: "AP"},
		{input: "Apress", output: "Apress"},
		{input: "Addison-Wesley", output: "AW"},
		{input: "Addison-Wesley Professional", output: "AW"},
		{input: "BCS", output: "BCS"},
		{input: "BPB", output: "BPB"},
		{input: "Cisco press", output: "Cisco"},
		{input: "Cengage Learning", output: "CL"},
		{input: "Course Technology", output: "CL"},
		{input: "South-Western College Publishing", output: "CL"},
		{input: "Apple academic press", output: "CRC"},
		{input: "A K Peters/CRC Press", output: "CRC"},
		{input: "Auerbach Publications", output: "CRC"},
		{input: "Chapman and Hall/CRC", output: "CRC"},
		{input: "CRC Press", output: "CRC"},
		{input: "Taylor & Francis", output: "CRC"},
		{input: "Taylor and Francis", output: "CRC"},
		{input: "Cambridge University Press", output: "CUP"},
		{input: "De Gruyter", output: "DG"},
		{input: "De Gruyter Oldenbourg", output: "DG"},
		{input: "DK", output: "DK"},
		{input: "Dorling Kindersley", output: "DK"},
		{input: "DK Children", output: "DK"},
		{input: "DK Publishing (Dorling Kindersley)", output: "DK"},
		{input: "Esri Press", output: "Esri"},
		{input: "For Dummies", output: "FD"},
		{input: "IET Standards", output: "IET"},
		{input: "Ivy Press", output: "Ivy"},
		{input: "The Institution of Engineering and Technology", output: "IET"},
		{input: "Institution of Engineering and Technology", output: "IET"},
		{input: "Inst of Engineering & Technology", output: "IET"},
		{input: "Jones & Bartlett Learning", output: "JBL"},
		{input: "Jones & Bartlett Publishers", output: "JBL"},
		{input: "Jones and Bartlett Publishers", output: "JBL"},
		{input: "Manning", output: "Manning"},
		{input: "Make Community, LLC", output: "Make"},
		{input: "Maker Media, Inc", output: "Make"},
		{input: "Morgan & Claypool Publishers", output: "MaC"},
		{input: "Morgan & Claypool", output: "MaC"},
		{input: "Morgan and Claypool", output: "MaC"},
		{input: "MIT Press", output: "MIT"},
		{input: "The MIT Press", output: "MIT"},
		{input: "Microsoft Press", output: "Microsoft"},
		{input: "McGraw Hill", output: "MGH"},
		{input: "McGraw-Hill Education", output: "MGH"},
		{input: "Mercury Learning & Information ", output: "ML"},
		{input: "Mercury Learning and Information", output: "ML"},
		{input: "Morgan Kaufmann", output: "MK"},
		{input: "Morgan Kaufmann Publishers", output: "MK"},
		{input: "Newnes", output: "Newnes"},
		{input: "Nova Science Pub Inc", output: "Nova"},
		{input: "Nova Science Publishers, Inc", output: "Nova"},
		{input: "No Starch Press", output: "NSP"},
		{input: "OReilly", output: "OReilly"},
		{input: "O′Reilly", output: "OReilly"},
		{input: "O'Reilly Media", output: "OReilly"},
		{input: "O'Reilly Media, Inc, USA", output: "OReilly"},
		{input: "O'Reilly UK Limited", output: "OReilly"},
		{input: "Oxford University Press", output: "OUP"},
		{input: "Oxford University Press Inc", output: "OUP"},
		{input: "Oxford University Press, Usa", output: "OUP"},
		{input: "OUP Oxford", output: "OUP"},
		{input: "Packt Publishing", output: "Packt"},
		{input: "Pearson", output: "Pearson"},
		{input: "Pearson College Div", output: "Pearson"},
		{input: "Pearson Education", output: "Pearson"},
		{input: "Pearson Education ESL", output: "Pearson"},
		{input: "Pragmatic Bookshelf", output: "Pragmatic"},
		{input: "Razeware LLC", output: "Razeware"},
		{input: "Sams", output: "Sams"},
		{input: "Sams Publishing", output: "Sams"},
		{input: "Springer", output: "Springer"},
		{input: "Wiley", output: "Wiley"},
		{input: "Wiley-Blackwell", output: "Wiley"},
		{input: "Unknown Weird Name", output: "Unknown Weird Name"},
	}

	t.Log("Given the need to test publisher name mapping.")
	for i, tt := range tests {
		t.Logf("\tTest: %d\tWhen checking %q for mapped value %q\n", i, tt.input, tt.output)
		mappedPublisher := MapPublisherName(tt.input)
		if mappedPublisher != tt.output {
			t.Errorf("\t\t%s\tShould get a %q mapped value: %q", failed, tt.output, mappedPublisher)
		} else {
			t.Logf("\t\t%s\tShould be able to map publisher name.", succeed)
		}
	}
}
