package parser

import (
	"encoding/xml"
	"regexp"
)

var refRegex = regexp.MustCompile(`[0-9]+ [0-9]+ R`)

type PDFBytes []byte

type Reference struct {
	ObjectID int
	GenID    int
}

type Info struct {
	XMLName  xml.Name `xml:"xmpmeta"`
	Author   string   `xml:"RDF>Description>creator"`
	Creator  string   `xml:"RDF>Description>CreatorTool"`
	Producer string   `xml:"RDF>Description>Producer"`
}

type XRefObject struct {
	ObjectID int
	Offset   int64
}

type XRef struct {
	StartID   int
	Count     int
	ObjectRef []XRefObject
}
