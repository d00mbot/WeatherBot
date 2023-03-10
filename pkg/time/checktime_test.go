package time

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CheckTime(t *testing.T) {

	tc := []struct {
		input    string
		expected string
	}{
		{input: "00", expected: "00"},
		{input: "01", expected: "01"},
		{input: "02", expected: "02"},
		{input: "03", expected: "03"},
		{input: "04", expected: "04"},
		{input: "05", expected: "05"},
		{input: "06", expected: "06"},
		{input: "07", expected: "07"},
		{input: "08", expected: "08"},
		{input: "09", expected: "09"},
		{input: "10", expected: "10"},
		{input: "11", expected: "11"},
		{input: "12", expected: "12"},
		{input: "13", expected: "13"},
		{input: "14", expected: "14"},
		{input: "15", expected: "15"},
		{input: "16", expected: "16"},
		{input: "17", expected: "17"},
		{input: "18", expected: "18"},
		{input: "19", expected: "19"},
		{input: "20", expected: "20"},
		{input: "21", expected: "21"},
		{input: "22", expected: "22"},
		{input: "23", expected: "23"},
		{input: "", expected: ""},
		{input: "aaa", expected: ""},
		{input: "/", expected: "1"},
		{input: "1", expected: "24"},
		{input: "001", expected: "0000"},
		{input: "!!@#$%", expected: "::::"},
	}
	for _, tc := range tc {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			result := CheckTime(tc.input)

			t.Logf("Calling checkTime(%v), result: %s", tc.input, result)

			assert.Equal(t, tc.expected, result,
				fmt.Sprintf("Incorrect result. Expected: %s", result))
		})
	}
}
