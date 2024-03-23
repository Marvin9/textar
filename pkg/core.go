package pkg

import (
	"fmt"
	"math"
	"slices"
	"strings"
)

type DictionaryIndex struct {
	dictIndexes        []Dictionary
	mappedSuffixArrays []SuffixArray
}

func NewDictionaryIndex(dictionary []Dictionary) *DictionaryIndex {
	newDictionary := make([]Dictionary, len(dictionary))

	copy(newDictionary, dictionary)

	AttachOffsets(newDictionary)

	di := &DictionaryIndex{}

	di.dictIndexes = make([]Dictionary, len(dictionary))
	di.mappedSuffixArrays = make([]SuffixArray, len(dictionary))

	for idx, dict := range newDictionary {
		di.dictIndexes[idx] = dict
		di.mappedSuffixArrays[idx] = *NewSuffixArray(dict.Raw())
	}

	return di
}

type matchedResult struct {
	dictionaryId   string
	index          int64
	originalString string
}

type SearchResult struct {
	matched     []matchedResult
	inputString string
	Occurrences int64
}

type SearchResultOpts struct {
	// prefix you want to see for matched result
	PrefixLength int
	// suffix you want to see for matched result
	SuffixLength int
	// filter result by dictionary ids
	DictionaryIds []string
}

const displayTemplate = `dictionary-id: %s
match: %s
`

func (sr SearchResult) Display(opts SearchResultOpts) {
	for _, matched := range sr.matched {
		str := strings.Builder{}

		if len(opts.DictionaryIds) > 0 && !slices.Contains(opts.DictionaryIds, matched.dictionaryId) {
			break
		}

		if opts.PrefixLength > 0 {
			ini := int(math.Max(0, float64(int(matched.index)-opts.PrefixLength)))
			str.WriteString(string(matched.originalString[ini:matched.index]))
		}

		str.WriteString(string(sr.inputString))

		if opts.SuffixLength > 0 {
			end := int(math.Min(float64(len(matched.originalString)), float64(matched.index+int64(len(sr.inputString))+int64(opts.SuffixLength)+1)))
			str.WriteString(string(matched.originalString[int(matched.index)+len(sr.inputString) : end]))
		}

		fmt.Println(fmt.Sprintf(displayTemplate, matched.dictionaryId, str.String()))
	}
}

func (di *DictionaryIndex) Search(str string) SearchResult {
	searchRes := SearchResult{
		matched:     make([]matchedResult, 0),
		inputString: str,
		Occurrences: 0,
	}

	for idx, dictionary := range di.dictIndexes {
		matched := di.mappedSuffixArrays[idx].Search(str)

		for _, match := range matched {
			searchRes.matched = append(searchRes.matched, matchedResult{
				dictionaryId: dictionary.Id,
				index:        match,
			})
		}
	}

	searchRes.Occurrences = int64(len(searchRes.matched))
	return searchRes
}
