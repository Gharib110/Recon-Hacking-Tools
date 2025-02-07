package encoder

import "errors"

type TagMap struct {
	m   map[string]interface{}
	has map[string]bool
}

type Metadata struct {
	Tags          *TagMap
	Lens          map[string]uint64
	Offsets       map[string]uint64
	Parent        interface{}
	ParentBuf     []byte
	CurrentOffset uint64
	CurrentField  string
}

type BinaryMarshall interface {
	MarshalBinary(*Metadata) ([]byte, error)
	UnmarshalBinary([]byte, *Metadata) error
}

func (t TagMap) Has(key string) bool {
	return t.has[key]
}

func (t TagMap) Set(key string, value interface{}) {
	t.m[key] = value
	t.has[key] = true
}

func (t TagMap) Get(key string) interface{} {
	return t.m[key]
}

func (t TagMap) GetString(key string) (string, error) {
	if !t.Has(key) {
		return "", errors.New("key does not exit in tag")
	}

	return t.Get(key).(string), nil
}

func (t TagMap) GetBool(key string) (int, error) {
	if !t.Has(key) {
		return 0, errors.New("key does not exit in tag")
	}

	return t.Get(key).(int), nil
}
