package hw03_frequency_analysis // nolint:golint,stylecheck

import (
	"sort"
	"strings"
)

type dictionary struct {
	word string
	freq int
}

// Top10 returns 10 the most frequent words in the input string.
func Top10(in string) []string {
	res := make([]string, 0) // Consider preallocating `res` (prealloc)
	words := strings.Fields(in)
	dict := make(map[string]int)
	for _, word := range words {
		dict[word]++
	}

	sorted := make([]dictionary, 0) // Consider preallocating `sorted` (prealloc)
	for word, freq := range dict {
		sorted = append(sorted, dictionary{word, freq})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].freq > sorted[j].freq
	})

	if len(sorted) >= 10 { // nolint:gomnd
		sorted = sorted[:10]
	}
	for i := range sorted {
		res = append(res, sorted[i].word)
	}
	return res
}
