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
	builder := make([]rune, len(rawString)+1)

	copy(builder, rawString)

	builder[len(rawString)] = '$'

	suffixArr := &SuffixArray{
		Suffixes:       make([]Suffix, len(builder)),
		OriginalString: builder,
		DictionaryId:   dictionaryId,
	}

	bucket := map[rune]int64{}

	// count suffixes
	for _, char := range builder {
		bucket[char]++
	}

	orderedKeys := make([]rune, 0)

	for bucketKey := range bucket {
		orderedKeys = append(orderedKeys, bucketKey)
	}

	sort.Slice(orderedKeys, func(i, j int) bool {
		if orderedKeys[i] == '$' {
			return true
		}
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

	for i := 0; i < len(builder)-1; i++ {
		ti := builder[i]
		tiplusone := builder[i+1]

		if ti > tiplusone || tiplusone == '$' {
			b1a[ti]++
		} else {
			b1b[ti]++
		}
	}
	b1b[0] = 1
	b1a[builder[len(builder)-1]]++

	c := int64(0)
	for _, orderedKey := range orderedKeys {
		b2a[orderedKey] = c
		c += b1a[orderedKey]
		c += b1b[orderedKey]
		b2b[orderedKey] = c
	}

	for i := 0; i < len(builder)-1; i++ {
		ti := builder[i]
		tiplusone := builder[i+1]

		if ti <= tiplusone && tiplusone != '$' {
			b2b[ti]--
			suffixArr.Suffixes[b2b[ti]] = Suffix{
				Char:          ti,
				OriginalIndex: int64(i),
			}
		}
	}

	ti := builder[len(builder)-1]
	suffixArr.Suffixes[0] = Suffix{
		Char:          ti,
		OriginalIndex: int64(len(builder) - 1),
	}

	for _, orderedKey := range orderedKeys {
		if b1b[orderedKey] > 1 {
			b := b2b[orderedKey]
			e := b + b1b[orderedKey] - 1

			// sort suffix [b, e]
			sort.SliceStable(suffixArr.Suffixes[b:e+1], func(i, j int) bool {
				return suffixArr.Suffixes[int(b)+i].OriginalIndex > suffixArr.Suffixes[int(b)+j].OriginalIndex
			})
		}
	}

	for i := 0; i < len(builder); i++ {
		ai := suffixArr.Suffixes[i].OriginalIndex
		aiminusone := ai - 1

		if ai == 0 {
			continue
		}

		if builder[ai] == '$' || builder[aiminusone] > builder[ai] {
			a := builder[aiminusone]
			suffixArr.Suffixes[b2a[a]] = Suffix{
				OriginalIndex: aiminusone,
				Char:          builder[aiminusone],
			}
			b2a[a]++
		}
	}

	// fmt.Println(suffixArr.Raw())

	return suffixArr
}

func (sa *SuffixArray) Search(str string) []int64 {
	low := 0
	high := len(sa.Suffixes) - 1

	indexes := []int64{}

	for low <= high {
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
		saRaw.WriteString(string(sa.OriginalString[suffixes.OriginalIndex:]))
		saRaw.WriteString("\n--------------------------------------------------\n\n")
	}

	return saRaw.String()
}
