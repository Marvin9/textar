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

	bucket := map[rune]int64{}

	// count suffixes
	for _, char := range rawString {
		bucket[char]++
	}

	orderedKeys := make([]rune, len(bucket))

	for bucketKey := range bucket {
		orderedKeys = append(orderedKeys, bucketKey)
	}

	sort.Slice(orderedKeys, func(i, j int) bool {
		return orderedKeys[i] < orderedKeys[j]
	})

	// allocate buckets
	bucket2 := map[rune]int64{}
	c := int64(0)
	for _, char := range orderedKeys {
		bucket2[char] = c
		c += bucket[char]
	}

	for i, char := range rawString {
		suffixArr.Suffixes[bucket2[char]] = Suffix{
			Char:          char,
			OriginalIndex: int64(i),
		}
		bucket2[char]++
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
