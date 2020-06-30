package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	var user User
	var matched bool
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		if err := user.UnmarshalJSON(scanner.Bytes()); err != nil {
			return nil, err
		}
		matched = strings.HasSuffix(user.Email, "."+domain)
		if matched && strings.Contains(user.Email, "@") {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
