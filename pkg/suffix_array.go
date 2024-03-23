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
	DictionaryId   string
}

func NewSuffixArray(rawString []rune, dictionaryId string) *SuffixArray {
	builder := make([]rune, len(rawString))

	copy(builder, rawString)

	suffixArr := &SuffixArray{
		Suffixes:       make([]Suffix, len(rawString)),
		OriginalString: builder,
		DictionaryId:   dictionaryId,
	}

	bucket := map[rune]int64{}

	// count suffixes
	for _, char := range rawString {
		bucket[char]++
	}

	orderedKeys := make([]rune, 0)

	for bucketKey := range bucket {
		orderedKeys = append(orderedKeys, bucketKey)
	}

	sort.Slice(orderedKeys, func(i, j int) bool {
		return orderedKeys[i] < orderedKeys[j]
	})

	b1a := map[rune]int64{}
	b1b := map[rune]int64{}
	b2a := map[rune]int64{}
	b2b := map[rune]int64{}

	for _, orderedKey := range orderedKeys {
		b1a[orderedKey] = 0
		b1b[orderedKey] = 0
	}

	for i := 0; i < len(rawString)-1; i++ {
		ti := rawString[i]
		tiplusone := rawString[i+1]

		if ti > tiplusone {
			b1a[ti]++
		} else {
			b1b[ti]++
		}
	}
	b1b[0] = 1

	c := int64(0)
	for _, orderedKey := range orderedKeys {
		b2a[orderedKey] = c
		c += b1a[orderedKey]
		c += b1b[orderedKey]
		b2b[orderedKey] = c
	}

	for i := 0; i < len(rawString)-1; i++ {
		ti := rawString[i]
		tiplusone := rawString[i+1]

		if ti <= tiplusone {
			b2b[ti]--
			suffixArr.Suffixes[b2b[ti]] = Suffix{
				Char:          ti,
				OriginalIndex: int64(i),
			}
		}
	}

	// for _, orderedKey := range orderedKeys {
	// 	if b1b[orderedKey] > 1 {
	// 		b := b2b[orderedKey]
	// 		e := b + b1b[orderedKey] - 1

	// 		// sort suffix [b, e]
	// 		sort.Slice(suffixArr.Suffixes[b:e+1], func(i, j int) bool {
	// 			return strings.Compare(string(rawString[suffixArr.Suffixes[i].OriginalIndex:]), string(rawString[suffixArr.Suffixes[j].OriginalIndex:])) < 0
	// 		})
	// 	}
	// }

	for i := 0; i < len(rawString)-1; i++ {
		ai := suffixArr.Suffixes[i].OriginalIndex
		aiminusone := ai - 1

		if aiminusone == -1 {
			aiminusone = int64(len(rawString) - 1)
		}

		if rawString[aiminusone] > rawString[ai] {
			a := rawString[aiminusone]
			suffixArr.Suffixes[b2a[a]] = Suffix{
				OriginalIndex: aiminusone,
				Char:          rawString[aiminusone],
			}
			b2a[a]++
		}
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
