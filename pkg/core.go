package pkg

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/achim-k/go-vebt"
)

type DictionaryLevelIndex struct {
	// h1, h2, p....
	Type string
	End  int
}

type LevelIndex = map[string]map[int64]DictionaryLevelIndex

type DictionaryIndex struct {
	dictIndexes        []Dictionary
	mappedSuffixArrays []SuffixArray
	levelIndex         LevelIndex
	// map[dictionary-id]
	vebIndex map[string]*vebt.VEB
}

func NewDictionaryIndex(dictionary []Dictionary) *DictionaryIndex {
	newDictionary := make([]Dictionary, len(dictionary))

	copy(newDictionary, dictionary)

	AttachOffsets(newDictionary)

	di := &DictionaryIndex{
		levelIndex: make(map[string]map[int64]DictionaryLevelIndex),
		vebIndex:   make(map[string]*vebt.VEB),
	}

	di.dictIndexes = make([]Dictionary, len(dictionary))
	di.mappedSuffixArrays = make([]SuffixArray, len(dictionary))

	for idx, dict := range newDictionary {
		dictRaw := dict.Raw()
		dict.CachedRaw = dictRaw
		di.dictIndexes[idx] = dict
		di.vebIndex[dict.Id] = vebt.CreateTree(dictRaw.Len())
		di.mappedSuffixArrays[idx] = *NewSuffixArray([]rune(dictRaw.String()), dict.Id)

		for _, index := range dict.Indexes {
			_, exist := di.levelIndex[dict.Id]

			if !exist {
				di.levelIndex[dict.Id] = make(map[int64]DictionaryLevelIndex)
			}

			di.levelIndex[dict.Id][index.Offset] = DictionaryLevelIndex{
				Type: index.Type,
				End:  len(index.String),
			}
			di.vebIndex[dict.Id].Insert(int(index.Offset))
		}
	}

	return di
}

type matchedResult struct {
	dictionaryId       string
	index              int64
	originalString     strings.Builder
	largestPredecessor struct {
		// [startIndex, endIndex]
		indexRange []int
		_type      string
	}
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
	// show the largest matched predecessor
	ShowParent bool
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

		orig := matched.originalString.String()

		pred := ""
		if opts.ShowParent {
			pred = fmt.Sprintf(
				"predecessor - %s - %s\n\n",
				matched.largestPredecessor._type,
				orig[matched.largestPredecessor.indexRange[0]:matched.largestPredecessor.indexRange[1]],
			)
		}

		if opts.PrefixLength > 0 {
			ini := int(math.Max(0, float64(int(matched.index)-opts.PrefixLength)))
			str.WriteString(string(orig[ini:matched.index]))
		}

		str.WriteString(string(sr.inputString))

		if opts.SuffixLength > 0 {
			end := int(math.Min(float64(len(orig)), float64(matched.index+int64(len(sr.inputString))+int64(opts.SuffixLength)+1)))
			str.WriteString(string(orig[int(matched.index)+len(sr.inputString) : end]))
		}

		fmt.Println(fmt.Sprintf(displayTemplate, matched.dictionaryId, str.String())+"\n", pred)
	}
}

func (di *DictionaryIndex) RawLevelIndex() string {
	str := strings.Builder{}
	for dictId, dict := range di.levelIndex {
		str.WriteString(fmt.Sprintf("dictionary id - %s\n", dictId))
		for idx, levels := range dict {
			str.WriteString(fmt.Sprintf("type - %s\n", levels.Type))
			for _, dict := range di.dictIndexes {
				if dictId == dict.Id {
					_dictRaw := dict.Raw()
					dictRaw := _dictRaw.String()
					end := int(math.Min(float64(len(dictRaw)), float64(int(idx)+levels.End)))
					str.WriteString(fmt.Sprintf("value - %s\n", string(dictRaw[idx:end])))
				}
			}
		}
	}
	return str.String()
}

func (di *DictionaryIndex) Search(str string) SearchResult {
	return di.search(str, SearchOpts{IncludeParentRelationShip: false})
}

func (di *DictionaryIndex) SearchWithOpts(str string, opts SearchOpts) SearchResult {
	return di.search(str, opts)
}

type SearchOpts struct {
	IncludeParentRelationShip bool
}

func (di *DictionaryIndex) search(str string, opts SearchOpts) SearchResult {
	searchRes := SearchResult{
		matched:     make([]matchedResult, 0),
		inputString: str,
		Occurrences: 0,
	}

	for idx, dictionary := range di.dictIndexes {
		matched := di.mappedSuffixArrays[idx].Search(str)

		_origStr := dictionary.Raw()
		veb := di.vebIndex[dictionary.Id]
		for _, match := range matched {
			largestPredecessor := []int{-1, -1}
			largestPredecessorType := ""
			if opts.IncludeParentRelationShip {
				// find first predecessor if type is null because this is in the middle of block
				startBlockOfMatch := match
				if _, ok := di.levelIndex[dictionary.Id][startBlockOfMatch]; !ok {
					startBlockOfMatch = int64(veb.Predecessor(int(match)))
				}

				predecessor := veb.Predecessor(int(startBlockOfMatch))

				for predecessor != -1 &&
					di.levelIndex[dictionary.Id][int64(predecessor)].Type == di.levelIndex[dictionary.Id][startBlockOfMatch].Type {
					predecessor = veb.Predecessor(predecessor)
				}

				if predecessor != -1 {
					largestPredecessor = []int{
						predecessor,
						predecessor + di.levelIndex[dictionary.Id][int64(predecessor)].End,
					}
					largestPredecessorType = di.levelIndex[dictionary.Id][int64(predecessor)].Type
				}
			}

			searchRes.matched = append(searchRes.matched, matchedResult{
				dictionaryId:   dictionary.Id,
				index:          match,
				originalString: _origStr,
				largestPredecessor: struct {
					indexRange []int
					_type      string
				}{
					indexRange: largestPredecessor,
					_type:      largestPredecessorType,
				},
			})
		}
	}

	searchRes.Occurrences = int64(len(searchRes.matched))
	return searchRes
}
