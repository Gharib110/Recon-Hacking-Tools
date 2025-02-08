package encoder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"unicode/utf16"
)

func FromUnicode(d []byte) (string, error) {
	if len(d)%2 > 0 {
		return "", errors.New("unicode specified, but odd data length")
	}
	s := make([]uint16, len(d)/2)
	err := binary.Read(bytes.NewReader(d), binary.LittleEndian, &s)
	if err != nil {
		return "", err
	}
	return string(utf16.Decode(s)), nil
}

func ToUnicode(s string) []byte {
	uints := utf16.Encode([]rune(s))
	b := bytes.Buffer{}
	binary.Write(&b, binary.LittleEndian, &uints)
	return b.Bytes()
}
