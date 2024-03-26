package pkg

import (
	"slices"
	"testing"
)

func TestSuffixArray(t *testing.T) {
	sa := NewSuffixArray([]rune("aabcaadaccdabcdac"), "0")

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

	matched = sa.Search("c")

	expectedMatchedAt = []int64{3, 8, 9, 13, 16}

	matcher()

	matched = sa.Search("dac")

	expectedMatchedAt = []int64{6, 14}

	matcher()

	matched = sa.Search("ac")

	expectedMatchedAt = []int64{7, 15}

	matcher()

	sa = NewSuffixArray([]rune(`Hi all,

	I need to travel to present a research paper next week, so I won't be able to lecture on March 19 or 21. Normally I would find a guest lecturer, but it happens that those who know the current topics well are on sabbatical leaves and/or also traveling for conferences; there happen to be two relevant conferences in the same week.
	
	Therefore, I have to cancel the lectures on March 19 and 21. I won't be able to hold office hours on those days, either, but feel free to email me when you have questions. The email response times may be higher than usual when I travel, but I will make sure to reply.
	
	The lectures for the week after next will continue as usual. If you finish assignment 4 early next week (it's actually due after I come back) and plan to spend more time on this course, perhaps you could do the following: If you are an undergraduate student, it would be a good idea to start to review the lectures to prepare for the final exam which will cover all lectures. If you are a graduate students, you could try to make more progress on your project; when giving a presentation during the last lecture (April 4), you will be expected to show the work you have completed by then.`), "0")

	matched = sa.Search("Hi all")

	expectedMatchedAt = []int64{0}

	matcher()

	sa = NewSuffixArray([]rune("Office Hours 11am-12pm Mondays"), "0")

	matched = sa.Search("Office")

	expectedMatchedAt = []int64{0}

	matcher()
}
