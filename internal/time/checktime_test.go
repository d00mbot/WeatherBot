package time

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckTime(t *testing.T) {
	tc := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "TestCase 1", input: "00", expected: "00"},
		{name: "TestCase 2", input: "01", expected: "01"},
		{name: "TestCase 3", input: "02", expected: "02"},
		{name: "TestCase 4", input: "03", expected: "03"},
		{name: "TestCase 5", input: "04", expected: "04"},
		{name: "TestCase 6", input: "05", expected: "05"},
		{name: "TestCase 7", input: "06", expected: "06"},
		{name: "TestCase 8", input: "07", expected: "07"},
		{name: "TestCase 9", input: "08", expected: "08"},
		{name: "TestCase 10", input: "09", expected: "09"},
		{name: "TestCase 11", input: "10", expected: "10"},
		{name: "TestCase 12", input: "11", expected: "11"},
		{name: "TestCase 13", input: "12", expected: "12"},
		{name: "TestCase 14", input: "13", expected: "13"},
		{name: "TestCase 15", input: "14", expected: "14"},
		{name: "TestCase 16", input: "15", expected: "15"},
		{name: "TestCase 17", input: "16", expected: "16"},
		{name: "TestCase 18", input: "17", expected: "17"},
		{name: "TestCase 19", input: "18", expected: "18"},
		{name: "TestCase 20", input: "19", expected: "19"},
		{name: "TestCase 21", input: "20", expected: "20"},
		{name: "TestCase 22", input: "21", expected: "21"},
		{name: "TestCase 23", input: "22", expected: "22"},
		{name: "TestCase 24", input: "23", expected: "23"},
		{name: "TestCase 25", input: "", expected: ""},
		{name: "TestCase 26", input: "aaa", expected: ""},
		{name: "TestCase 27", input: "0/1", expected: "1"},
		{name: "TestCase 28", input: "1", expected: "1"},
		{name: "TestCase 29", input: "001", expected: "001"},
		{name: "TestCase 30", input: "!!@#$%", expected: "::::"},
		{name: "TestCase 31", input: "01 ", expected: "01"},
		{name: "TestCase 32", input: " 01", expected: "01"},
		{name: "TestCase 33", input: "0 1", expected: ""},
	}

	for _, tc := range tc {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := CheckTime(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}
