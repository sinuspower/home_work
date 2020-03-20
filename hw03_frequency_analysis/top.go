package hw03_frequency_analysis // nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
)

type wordFreq struct {
	word string
	freq int
}

func regexpReplace(s, r string) string {
	return regexp.MustCompile(`(^[!,.:\-"']+|[!,.:\-"']+$)`).ReplaceAllString(s, r)
}

func regexpEquals(s1, s2 string) bool {
	if s1 == "" && s2 == "" {
		return true
	}
	s1 = regexp.MustCompile(`(^[!,.:\-"']+|[!,.:\-"']+$)`).ReplaceAllString(s1, "") // remove all punctuation from the left and the right of the first string
	re := regexp.MustCompile(`(?i)^[\-'"]*` + s1 + `[!,.:\-'"]*$`)                  // build regexp based on the first string
	return re.FindString(s2) != ""                                                  // check the second string by regexp based on the first string
}

func add(wf []wordFreq, w string, f int) []wordFreq {
	var found bool
	for i, item := range wf {
		if regexpEquals(w, item.word) {
			wf[i].freq += f
			found = true
			break
		}
	}
	w = regexpReplace(w, "")
	if !found && w != "-" && w != "" {
		wf = append(wf, wordFreq{strings.ToLower(w), f})
	}
	return wf
}

// Top10 returns 10 the most frequent words in the input string.
func Top10(in string) []string {
	res := make([]string, 0) // Consider preallocating `res` (prealloc)
	words := strings.Fields(in)
	dict := make(map[string]int)
	for _, word := range words {
		dict[word]++
	}
	sorted := make([]wordFreq, 0) // Consider preallocating `sorted` (prealloc)
	for word, freq := range dict {
		sorted = add(sorted, word, freq)
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
