package pkg

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/Marvin9/textar/brute"
)

func BenchmarkTextarShakespear(b *testing.B) {
	// Read the text file
	content, err := os.ReadFile("../data/shakespeare.txt")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	text := string(content)

	// Split the text into dictionaries based on character count
	var dictionaries []Dictionary
	for i := 0; i < len(text); i += len(text) {
		end := i + len(text)
		if end > len(text) {
			end = len(text)
		}
		dictionaryText := text[i:end]

		lines := strings.Split(dictionaryText, "\n")

		indexes := make([]Index, 0)
		for _, line := range lines {
			index := NewIndex("p", line, int64(i))
			indexes = append(indexes, *index)
		}
		dictionary := NewDictionary(fmt.Sprintf("Dictionary%d", len(dictionaries)+1), indexes)

		dictionaries = append(dictionaries, *dictionary)
	}

	index := NewDictionaryIndex(dictionaries)

	searches := []string{
		"In singleness the parts that thou shouldst bear",
		"How can I then return in happy plight",
		"To see his active",
		"to",
		"But yet be blamed, if thou thy self deceivest",
		"Thus can my love excuse the slow offence",
		"End of this Etext",
		"whitely",
	}

	b.Run("search", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, search := range searches {
				index.Search(search)
			}
		}
	})
}

func BenchmarkBruteForceShakespear(b *testing.B) {
	// Read the text file
	content, err := os.ReadFile("../data/shakespeare.txt")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	text := string(content)

	searches := []string{
		"In singleness the parts that thou shouldst bear",
		"How can I then return in happy plight",
		"To see his active",
		"to",
		"But yet be blamed, if thou thy self deceivest",
		"Thus can my love excuse the slow offence",
		"End of this Etext",
		"whitely",
	}

	b.Run("search", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, search := range searches {
				brute.Search(text, search)
			}
		}
	})
}
