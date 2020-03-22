package hw03_frequency_analysis // nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
)

type frequency struct {
	word string
	freq int
}

func regexpReplace(s, re, rs string) string {
	return regexp.MustCompile(re).ReplaceAllString(s, rs)
}

// Top10 returns 10 the most frequent words in the input string.
func Top10(in string) []string {
	if in == "" {
		return nil
	}
	words := strings.Fields(in)
	dict := make(map[string]int)
	for _, word := range words {
		word = regexpReplace(word, `(^[!?,.:\-"']+|[!?,.:\-"']+$)`, "")
		if word != "" {
			dict[strings.ToLower(word)]++
		}
	}
	sorted := make([]frequency, 0) // Consider preallocating `sorted` (prealloc)
	for word, freq := range dict {
		sorted = append(sorted, frequency{word, freq})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].freq > sorted[j].freq
	})
	if len(sorted) >= 10 { // nolint:gomnd // magic numbers
		sorted = sorted[:10]
	}
	res := make([]string, 0) // Consider preallocating `res` (prealloc)
	for i := range sorted {
		res = append(res, sorted[i].word)
	}
	return res
}
