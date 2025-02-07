package encoder

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

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

func parseTags(sf reflect.StructField) (*TagMap, error) {
	ret := &TagMap{
		m:   make(map[string]interface{}),
		has: make(map[string]bool),
	}
	tag := sf.Tag.Get("smb")
	smbTags := strings.Split(tag, ",")
	for _, smbTag := range smbTags {
		tokens := strings.Split(smbTag, ":")
		switch tokens[0] {
		case "len", "offset", "count":
			if len(tokens) != 2 {
				return nil, errors.New("invalid tag")
			}
			ret.Set(tokens[0], tokens[1])

		case "fixed":
			if len(tokens) != 2 {
				return nil, errors.New("invalid tag")
			}
			i, err := strconv.Atoi(tokens[1])
			if err != nil {
				return nil, err
			}
			ret.Set(tokens[0], i)
		case "asn1":
			ret.Set(tokens[0], true)
		}
	}
	return ret, nil
}

func getOffsetByFieldName(fieldName string, meta *Metadata) (uint64, error) {
	if meta == nil || meta.Tags == nil ||
		meta.Parent == nil || meta.Lens == nil {
		return 0, errors.New("invalid metadata")
	}
	var ret uint64
	var found bool

	parentvf := reflect.Indirect(reflect.ValueOf(meta.Parent))
	for i := 0; i < parentvf.NumField(); i++ {
		field := parentvf.Type().Field(i)
		if field.Name == fieldName {
			found = true
			break
		}
		if l, ok := meta.Lens[field.Name]; ok {
			ret += l
		} else {
			buf, err := Marshal(parentvf.Field(i).Interface())
			if err != nil {
				return 0, err
			}
			l := uint64(len(buf))
			meta.Lens[field.Name] = l
			ret += l
		}

	}
	if !found {
		return 0, errors.New("field not found")
	}
	return ret, nil
}

func Marshal(v interface{}) ([]byte, error) {
	return marshal(v, nil)
}

func marshal(v interface{}, meta *Metadata) ([]byte, error) {
	return nil, nil
}
