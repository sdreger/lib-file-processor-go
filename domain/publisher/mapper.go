package publisher

import "strings"

var publisherNameMapping = map[string]string{
	"acm books":                        "MaC",
	"alpha science":                    "Alpha",
	"amer radio relay league":          "ARRL",
	"academic press":                   "AP",
	"apress":                           "Apress",
	"addison-wesley":                   "AW",
	"arcler press":                     "Arcler",
	"artech house":                     "Artech",
	"bcs":                              "BCS",
	"big nerd ranch":                   "BNR",
	"bpb":                              "BPB",
	"birkhäuser":                       "Springer",
	"butterworth-heinemann":            "BH",
	"cisco":                            "Cisco",
	"cengage":                          "CL",
	"course technology":                "CL",
	"south-western college publishing": "CL",
	"apple academic press":             "CRC",
	"auerbach":                         "CRC",
	"chapman":                          "CRC",
	"crc":                              "CRC",
	"taylor & francis":                 "CRC",
	"taylor and francis":               "CRC",
	"cambridge university press":       "CUP",
	"cognella academic":                "Cognella",
	"de gruyter":                       "DG",
	"de|g":                             "DG",
	"dk":                               "DK",
	"dk children":                      "DK",
	"dorling kindersley":               "DK",
	"esri":                             "Esri",
	"for dummies":                      "FD",
	"focal press":                      "Focal",
	"hodder":                           "Hodder",
	"iet standards":                    "IET",
	"i/o press":                        "Io",
	"iop publishing":                   "Iop",
	"institute of physics":             "Iop",
	"ivy press":                        "Ivy",
	"of engineering and technology":    "IET",
	"of engineering & technology":      "IET",
	"jones & bartlett":                 "JBL",
	"jones and bartlett":               "JBL",
	"j. ross publishing":               "JRP",
	"manning":                          "Manning",
	"make community":                   "Make",
	"maker media":                      "Make",
	"morgan & claypool":                "MaC",
	"morgan and claypool":              "MaC",
	"mit press":                        "MIT",
	"microsoft":                        "Microsoft",
	"mcgraw-hill":                      "MGH",
	"mcgraw hill":                      "MGH",
	"mercury learning":                 "ML",
	"morgan kaufmann":                  "MK",
	"newnes":                           "Newnes",
	"nova":                             "Nova",
	"no starch":                        "NSP",
	"orange education":                 "Orange",
	"oreilly":                          "OReilly",
	"o'reilly":                         "OReilly",
	"o′reilly":                         "OReilly",
	"oracle":                           "Oracle",
	"oxford university press":          "OUP",
	"oup oxford":                       "OUP",
	"packt":                            "Packt",
	"pearson":                          "Pearson",
	"pragmatic":                        "Pragmatic",
	"princeton":                        "Princeton",
	"que":                              "Que",
	"raspberry pi press":               "RPIP",
	"razeware":                         "Razeware",
	"rheinwerk":                        "Rheinwerk",
	"river publishers":                 "River",
	"sams":                             "Sams",
	"springer":                         "Springer",
	"toronto academic press":           "TAP",
	"visual":                           "Wiley",
	"wiley":                            "Wiley",
	"world scientific":                 "WSPC",
}

// MapPublisherName map the full publisher name to its short form.
// If no mapping found - returns the full publisher name from the input.
func MapPublisherName(publisherFullName string) string {
	lowerPublisherName := strings.ToLower(publisherFullName)
	if shortName, ok := publisherNameMapping[lowerPublisherName]; ok {
		return shortName
	}

	for key, value := range publisherNameMapping {
		if strings.Contains(lowerPublisherName, key) {
			return value
		}
	}

	return publisherFullName
}
