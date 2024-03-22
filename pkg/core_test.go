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
	for i := 0; i < len(text); i += 1500 {
		end := i + 1500
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

	b.Run("search", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			index.Search("In singleness the parts that thou shouldst bear")
			index.Search("How can I then return in happy plight")
			index.Search("To see his active")
			index.Search("to")
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

	b.Run("search", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			brute.Search(text, "In singleness the parts that thou shouldst bear")
			brute.Search(text, "How can I then return in happy plight")
			brute.Search(text, "To see his active")
			brute.Search(text, "to")
		}
	})
}
