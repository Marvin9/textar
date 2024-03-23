package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Marvin9/textar/pkg"
)

func _main() {
	// Read the text file
	content, err := os.ReadFile("./data/shakespeare.txt")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	text := string(content)

	// Split the text into dictionaries based on character count
	var dictionaries []pkg.Dictionary
	for i := 0; i < len(text); i += 1500 {
		end := i + 1500
		if end > len(text) {
			end = len(text)
		}
		dictionaryText := text[i:end]

		lines := strings.Split(dictionaryText, "\n")

		indexes := make([]pkg.Index, 0)
		for _, line := range lines {
			index := pkg.NewIndex("p", line, int64(i))
			indexes = append(indexes, *index)
		}
		dictionary := pkg.NewDictionary(fmt.Sprintf("Dictionary%d", len(dictionaries)+1), indexes)

		dictionaries = append(dictionaries, *dictionary)
	}

	index := pkg.NewDictionaryIndex(dictionaries)

	index.Search("to").Display(pkg.SearchResultOpts{
		PrefixLength: 10,
		SuffixLength: 10,
	})
}

func main() {
	dictionaries := []pkg.Dictionary{
		{
			Id: "Advanced Data Structures",
			Indexes: []pkg.Index{
				{
					Type:   "h1",
					String: "CSCI6057 CSCI4117 - Advanced Data Structures - Sec: 01 - 2023/2024 Winter",
				},
				{
					Type:   "h2",
					String: "Lectures on March 19 and 21",
				},
				{
					Type: "p",
					String: `Hi all,

					I need to travel to present a research paper next week, so I won't be able to lecture on March 19 or 21. Normally I would find a guest lecturer, but it happens that those who know the current topics well are on sabbatical leaves and/or also traveling for conferences; there happen to be two relevant conferences in the same week.
					
					Therefore, I have to cancel the lectures on March 19 and 21. I won't be able to hold office hours on those days, either, but feel free to email me when you have questions. The email response times may be higher than usual when I travel, but I will make sure to reply.
					
					The lectures for the week after next will continue as usual. If you finish assignment 4 early next week (it's actually due after I come back) and plan to spend more time on this course, perhaps you could do the following: If you are an undergraduate student, it would be a good idea to start to review the lectures to prepare for the final exam which will cover all lectures. If you are a graduate students, you could try to make more progress on your project; when giving a presentation during the last lecture (April 4), you will be expected to show the work you have completed by then.`,
				},
				{
					Type:   "h2",
					String: "Submission open for project proposal (CSCI 6057 only)",
				},
				{
					Type: "p",
					String: `I have opened submission for project proposals. You can find it under assessments->assignments->Project Proposal (CSCI 6057 only).

					Note that this is for students enrolled in CSCI 6057 only.
					
					As scheduled at the beginning of this term, the next assignment will be posted on March 12, so that graduate students will have time for the proposals.`,
				},
			},
		},
		{
			Id: "Algorithm Engineering",
			Indexes: []pkg.Index{
				{
					Type:   "h1",
					String: "CSCI4118 CSCI6105 - Algorithm Engineering - Sec: 01 - 2023/2024 Winter",
				},
				{
					Type:   "h2",
					String: "Snow Day - Schedule Changes",
				},
				{
					Type: "p",
					String: `Hi all,

					With Dal closed today, we'll:
					
					extend the proposal deadline to tomorrow Thursday Feb 15 at 11:59 PM
					extend the current quiz deadline to before Friday's class
					continue the lecture schedule with Lecture 10 on Friday
					reminder: reading week is next week
					reminder: Lab 2 is Monday Feb 26 after the break
					Lab 3 will now be March 11
					Lab 4 will now be March 25
					combine the case study and review lectures on April 5
					Presentation Days April 8 and 9 are unchanged
					I've uploaded a new syllabus with these changes in the tentative schedule on page 6 and will modify the due dates on Brightspace. Please note I did not change all the dates listed on the syllabus - please refer to Brightspace or that tentative schedule.`,
				},
				{
					Type:   "h2",
					String: "Office Hours 11am-12pm Mondays",
				},
				{
					Type: "p",
					String: `Hi all,

					My office hours will be 11am-12pm on Mondays. I can answer questions via MS Teams or in person. Please note if you want to meet me in person in my office in Room 4244 of the Mona Campbell building please reply in MS Teams when you arrive so I can let you in the security doors.
					
					I've selected these times so that the last 30 mins overlap with our scheduled lab time so I know that students are available when we don't have labs. If this time doesn't work for you (particularly during lab weeks) please send me a message on Teams and I can try to answer it during office hours or when I find time, or we can try to find an alternative time to meet.`,
				},
			},
		},
	}

	index := pkg.NewDictionaryIndex(dictionaries)

	index.SearchWithOpts("Hi all", pkg.SearchOpts{IncludeParentRelationShip: true}).Display(pkg.SearchResultOpts{
		// PrefixLength: 10,
		SuffixLength: 10,
		ShowParent:   true,
	})
}
