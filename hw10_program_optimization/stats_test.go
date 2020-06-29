// +build !bench

package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

var data string = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

func TestGetDomainStat(t *testing.T) {
	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestGetEmail(t *testing.T) {
	testCases := []struct {
		in       string
		expected string
	}{
		{"", ""},
		{`{"Id":3,"Name":"Brian Olson"}`, ""},
		{"email:email@Quinu.edu", ""},
		{`{"email:FrancesEllis@Quinu.edu"}`, ""},
		{`{"email":"FrancesEllis@Quinu.edu"}`, ""},
		{`{"Email":"FrancesEllis@Quinu.edu"}`, "FrancesEllis@Quinu.edu"},
		{`{"Email":"Email@Quinu.edu"}`, "Email@Quinu.edu"},
		{
			`{"Id":3,"Name":"Brian Olson","Username":"non_quia_id","Email":"FrancesEllis@Quinu.edu","Phone":"237-75-34","Password":"cmEPhX8","Address":"Butterfield Junction 74"}`,
			"FrancesEllis@Quinu.edu",
		},
		{
			`{"Id":3,"Name":"Brian Olson","Username":"non_quia_id","Email":"email","Phone":"237-75-34","Password":"cmEPhX8","Address":"Butterfield Junction 74"}`,
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			require.Equal(t, tc.expected, getEmail(tc.in))
		})
	}
}

func BenchmarkGetDomainStat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetDomainStat(bytes.NewBufferString(data), "com")
	}
}
