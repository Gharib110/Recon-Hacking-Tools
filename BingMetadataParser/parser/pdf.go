package parser

import (
	"bytes"
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

func NewPDFData(buf []byte, stripNewLines bool) PDFBytes {
	var ret PDFBytes
	b := bytes.Trim(buf, "\x20\x09\x00\x0C")
	if stripNewLines {
		b = bytes.Replace(b, []byte("\x0A"), []byte{}, -1)
		b = bytes.Replace(b, []byte("\x0D"), []byte{}, -1)
	}
	ret = PDFBytes(b)
	return ret
}
