package pkg

import (
	"math"
	"sort"
	"strings"
)

type Suffix struct {
	Char          rune
	OriginalIndex int64
}

type SuffixArray struct {
	Suffixes       []Suffix
	OriginalString []rune
}

func NewSuffixArray(rawString []rune) *SuffixArray {
	builder := make([]rune, len(rawString))

	copy(builder, rawString)

	suffixArr := &SuffixArray{
		Suffixes:       make([]Suffix, len(rawString)),
		OriginalString: builder,
	}

	indexer := map[rune][]int64{}

	for i, rawRune := range rawString {
		_, exist := indexer[rawRune]

		if !exist {
			indexer[rawRune] = []int64{}
		}

		indexer[rawRune] = append(indexer[rawRune], int64(i))
	}

	sort.Slice(rawString, func(i, j int) bool {
		return rawString[i] < rawString[j]
	})

	for idx, rn := range rawString {
		if rn == rune(' ') || rn == rune('\n') {
			continue
		}

		item := indexer[rn]

		itemTopIdx := len(item) - 1
		suffixArr.Suffixes[idx] = Suffix{
			Char:          rn,
			OriginalIndex: item[itemTopIdx],
		}

		item = item[:itemTopIdx]
		indexer[rn] = item
	}

	return suffixArr
}

func (sa *SuffixArray) Search(str string) []int64 {
	low := 0
	high := len(sa.Suffixes) - 1

	indexes := []int64{}

	for low < high {
		mid := (low + high) / 2

		if rune(str[0]) < sa.Suffixes[mid].Char {
			high = mid - 1
		} else if rune(str[0]) > sa.Suffixes[mid].Char {
			low = mid + 1
		} else {
			start := mid
			end := mid + 1

			for start >= 0 && rune(str[0]) == sa.Suffixes[start].Char {
				if sa.match(str, sa.Suffixes[start]) {
					indexes = append(indexes, sa.Suffixes[start].OriginalIndex)
				}
				start--
			}

			for end < len(sa.Suffixes) && rune(str[0]) == sa.Suffixes[end].Char {
				if sa.match(str, sa.Suffixes[end]) {
					indexes = append(indexes, sa.Suffixes[end].OriginalIndex)
				}
				end++
			}

			return indexes
		}
	}

	return indexes
}

func (sa *SuffixArray) match(str string, suffix Suffix) bool {
	origIdx := suffix.OriginalIndex
	end := math.Min(float64(len(sa.OriginalString)), float64(int(origIdx)+len(str)))
	return string(sa.OriginalString[origIdx:int64(end)]) == str
}

func (sa *SuffixArray) Raw() string {
	saRaw := &strings.Builder{}

	for _, suffixes := range sa.Suffixes {
		saRaw.WriteRune(suffixes.Char)
	}

	return saRaw.String()
}
