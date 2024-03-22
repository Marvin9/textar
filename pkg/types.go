package pkg

import (
	"strings"
)

type Index struct {
	// h1, h2, p....
	Type string
	// actual string
	String string
	// initial character index from start
	Offset int64
}

func NewIndex(indexType string, actualString string, offset int64) *Index {
	return &Index{
		Type:   indexType,
		String: actualString,
		Offset: offset,
	}
}

type Dictionary struct {
	Id      string
	Indexes []Index
}

func NewDictionary(id string, indexes []Index) *Dictionary {
	return &Dictionary{
		Id:      id,
		Indexes: indexes,
	}
}

func (d *Dictionary) Raw() []rune {
	str := strings.Builder{}

	for _, index := range d.Indexes {
		str.WriteString(index.String)
	}

	return []rune(str.String())
}
