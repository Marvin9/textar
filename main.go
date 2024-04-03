package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Marvin9/textar/pkg"
)

func main() {
	type Item struct {
		ID   string `json:"id"`
		Data []struct {
			Type   string `json:"type"`
			String string `json:"string"`
		} `json:"data"`
	}

	// Open the file
	file, err := os.Open("../textar-data/github.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Parse the JSON content
	var items []Item
	if err := json.NewDecoder(file).Decode(&items); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	var dictionaries []pkg.Dictionary

	for _, dict := range items {
		index := make([]pkg.Index, 0)

		for _, data := range dict.Data {
			index = append(index, pkg.Index{
				Type:   data.Type,
				String: data.String,
			})
		}

		dictionaries = append(dictionaries, *pkg.NewDictionary(dict.ID, index))
	}

	var index *pkg.DictionaryIndex
	indexTime := pkg.MeasureTime(func() {
		index = pkg.NewDictionaryIndex(dictionaries)
	})

	// Continuously read for search input
	fmt.Printf("Index constructed in %s\n", indexTime)
	fmt.Println("Enter search query (type 'exit' to quit):")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Search: ")
		query, _ := reader.ReadString('\n')
		query = strings.TrimSpace(query)

		if query == "exit" {
			break
		}

		var searchRes pkg.SearchResult
		latency := pkg.MeasureTime(func() {
			searchRes = index.SearchWithOpts(query, pkg.SearchOpts{
				IncludeParentRelationShip: true,
			})
		})

		fmt.Printf("%d search results found in %s\n", searchRes.Occurrences, latency)
		searchRes.Display(pkg.SearchResultOpts{
			SuffixLength: 100,
			ShowParent:   true,
		})
	}
}
