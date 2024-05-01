# Textar

Document search using Suffix Array and Van Embde Boas Trees.

This was final project of the subject Advanced Data Structure (CSCI 6057) in supervision of professor [Dr. Meng He](https://www.dal.ca/faculty/computerscience/faculty-staff/meng-he.html)

## Use Case

If you have levelled documents in JSON format, then you can use this SDK to index all documents and search through it efficiently. Most suitable use-case (and this project is tested on) is website search. On the build time of the static website or documentation site, you can omit the index data in JSON which can then be used by the Golang server using this SDK.

```go
import "github.com/Marvin9/textar/pkg"

func main() {
    index := pkg.NewDictionaryIndex(
		[]pkg.Dictionary{
			{
				Id: "home-page",
				Indexes: []pkg.Index{
					{
						Type:   "h1",
						String: "This is heading",
					},
					{
						Type:   "p",
						String: "This is paragraph",
					},
				},
			},
		},
	)

	s := index.SearchWithOpts("paragraph", pkg.SearchOpts{
		IncludeParentRelationShip: true,
	})

	s.Display(pkg.SearchResultOpts{
		PrefixLength: 10,
		SuffixLength: 10,
		ShowParent:   true,
	})
}
```

## Mechanism

We have a string "S" of length "n". This string is divided into blocks, each denoted as "Bi". These blocks represent different parts of the text, such as paragraphs, headings, or subheadings. Each block "Bi" contains a substring of "S" starting from index "j", denoted as "Bi[j:]", along with an offset indicating the index of its first character in the original string "S", and a type "T".

Let "Di" represent the document that contains the string "S" and its blocks.

This simplification allows us to focus on the fundamental structure of the text and how it's organised into blocks within documents.

We have "n" Documents with String "S" and "m" Blocks:

"i" = 0...n
"D0, D1, D2, ..., Dn"
For each document "Di":

"j" = 0...m
"B0, B1, B2, ..., Bm"
For each block "Bj":

"Bj.T" represents the type of the block (e.g., paragraph, heading, subheading)
"Bj.Offset" is calculated as:
"Bj.Offset = j > 0 ? Bj-1.Offset + length(Bj-1.String) : 0"
Within each document "Di", we have substrings:

"k" = j...o
"DiSj, DiSj+1, ..., DiSo"

We'll create a suffix array index for each dictionary "Di" and a Van Emde Boas tree for each dictionary "Di". Then, for each search query "Q" of string "Str", we'll search for all occurrences of "Str" in all dictionaries "Di=0...n" using the suffix array index of dictionary "Di". For each matching substring "Substr" in "Di" and block "Bj", we'll find the first block "Bk" where "k < j" and "Bk.T != Bj.T", utilising Van Emde Boas trees.
