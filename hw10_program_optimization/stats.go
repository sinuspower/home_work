package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"io"
	"strings"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	var email string
	var matched bool
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		email = getEmail(scanner.Text())
		if matched = strings.Contains(email, "."+domain); matched {
			result[strings.ToLower(strings.SplitN(email, "@", 2)[1])]++
		}
	}
	return result, nil
}

func getEmail(in string) string {
	if !strings.Contains(in, `"Email":"`) {
		return ""
	}
	out := strings.SplitN(strings.SplitN(in, `"Email":"`, 2)[1], `"`, 2)[0]
	if !strings.Contains(out, "@") {
		return ""
	}
	return out
}
