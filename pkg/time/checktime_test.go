package time

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckTime(t *testing.T) {
	tc := map[string]struct {
		input    string
		expected string
	}{
		"TestCase 1":  {input: "00", expected: "00"},
		"TestCase 2":  {input: "01", expected: "01"},
		"TestCase 3":  {input: "02", expected: "02"},
		"TestCase 4":  {input: "03", expected: "03"},
		"TestCase 5":  {input: "04", expected: "04"},
		"TestCase 6":  {input: "05", expected: "05"},
		"TestCase 7":  {input: "06", expected: "06"},
		"TestCase 8":  {input: "07", expected: "07"},
		"TestCase 9":  {input: "08", expected: "08"},
		"TestCase 10": {input: "09", expected: "09"},
		"TestCase 11": {input: "10", expected: "10"},
		"TestCase 12": {input: "11", expected: "11"},
		"TestCase 13": {input: "12", expected: "12"},
		"TestCase 14": {input: "13", expected: "13"},
		"TestCase 15": {input: "14", expected: "14"},
		"TestCase 16": {input: "15", expected: "15"},
		"TestCase 17": {input: "16", expected: "16"},
		"TestCase 18": {input: "17", expected: "17"},
		"TestCase 19": {input: "18", expected: "18"},
		"TestCase 20": {input: "19", expected: "19"},
		"TestCase 21": {input: "20", expected: "20"},
		"TestCase 22": {input: "21", expected: "21"},
		"TestCase 23": {input: "22", expected: "22"},
		"TestCase 24": {input: "23", expected: "23"},
		"TestCase 25": {input: "", expected: ""},
		"TestCase 26": {input: "aaa", expected: ""},
		"TestCase 27": {input: "0/1", expected: "1"},
		"TestCase 28": {input: "1", expected: "1"},
		"TestCase 29": {input: "001", expected: "001"},
		"TestCase 30": {input: "!!@#$%", expected: "::::"},
		"TestCase 31": {input: "01 ", expected: "01"},
		"TestCase 32": {input: " 01", expected: "01"},
		"TestCase 33": {input: "0 1", expected: ""},
	}

	for name, tc := range tc {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := CheckTime(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}
