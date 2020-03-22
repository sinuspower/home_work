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

func regexpReplace(s, rs string, re *regexp.Regexp) string {
	return re.ReplaceAllString(s, rs)
}

// Top10 returns 10 the most frequent words in the input string.
func Top10(in string) []string {
	if in == "" {
		return nil
	}
	words := strings.Fields(in)
	dict := make(map[string]int)
	re := regexp.MustCompile(`(^[!?,.:\-"']+|[!?,.:\-"']+$)`)
	for _, word := range words {
		word = regexpReplace(word, "", re)
		if word != "" {
			dict[strings.ToLower(word)]++
		}
	}
	sorted := make([]frequency, 0, len(dict))
	for word, freq := range dict {
		sorted = append(sorted, frequency{word, freq})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].freq > sorted[j].freq
	})
	limit := 10
	if len(sorted) >= limit {
		sorted = sorted[:limit]
	}
	limit = len(sorted)
	res := make([]string, 0, limit)
	for i := range sorted {
		res = append(res, sorted[i].word)
	}
	return res
}
