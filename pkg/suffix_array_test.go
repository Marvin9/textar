package pkg

import (
	"slices"
	"testing"
)

func TestSuffixArray(t *testing.T) {
	sa := NewSuffixArray([]rune("foobar"))

	if string(sa.OriginalString) != "foobar" {
		t.Errorf("expected original string %s, got %s", "foobar", string(sa.OriginalString))
	}

	expectedSa := SuffixArray{
		Suffixes: []Suffix{
			{
				Char:          rune('a'),
				OriginalIndex: 4,
			},
			{
				Char:          rune('b'),
				OriginalIndex: 3,
			},
			{
				Char:          rune('f'),
				OriginalIndex: 0,
			},
			{
				Char:          rune('o'),
				OriginalIndex: 1,
			},
			{
				Char:          rune('o'),
				OriginalIndex: 2,
			},
			{
				Char:          rune('r'),
				OriginalIndex: 5,
			},
		},
	}

	for idx, expSa := range expectedSa.Suffixes {
		if sa.Suffixes[idx].Char != expSa.Char {
			t.Errorf("got char %s, expected %s at suffix array %d", string(sa.Suffixes[idx].Char), string(expSa.Char), idx)
		}

		if sa.Suffixes[idx].OriginalIndex != expSa.OriginalIndex {
			t.Errorf("got index %d, expected %d at suffix array %d", sa.Suffixes[idx].OriginalIndex, expSa.OriginalIndex, idx)
		}
	}

	wrongEvaluation := func() {
		t.Errorf("wrong evaluation for substring bar")
	}

	if !sa.match("ar", sa.Suffixes[0]) {
		wrongEvaluation()
	}

	if !sa.match("bar", sa.Suffixes[1]) {
		wrongEvaluation()
	}

	if sa.match("z", sa.Suffixes[3]) {
		wrongEvaluation()
	}

	sa = NewSuffixArray([]rune("aabcaadaccdabcdac"))

	matched := sa.Search("aa")

	expectedMatchedAt := []int64{0, 4}

	matcher := func() {

		for _, exp := range expectedMatchedAt {
			if !slices.Contains(matched, exp) {
				t.Errorf("expected %d to be matched", exp)
			}
		}
	}

	matcher()

	matched = sa.Search("dac")

	expectedMatchedAt = []int64{6, 14}

	matcher()

	matched = sa.Search("ac")

	expectedMatchedAt = []int64{7, 15}

	matcher()
}
