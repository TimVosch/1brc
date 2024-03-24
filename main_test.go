package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTemperature(t *testing.T) {
	testCases := []struct {
		input  []byte
		output int16
	}{
		{
			input:  []byte("-12.5"),
			output: -125,
		},
		{
			input:  []byte("12.5"),
			output: 125,
		},
		{
			input:  []byte("-2.5"),
			output: -25,
		},
		{
			input:  []byte("2.5"),
			output: 25,
		},
	}
	for _, tC := range testCases {
		t.Run(string(tC.input), func(t *testing.T) {
			got := ParseTemperature(tC.input)
			if got != tC.output {
				t.Logf("Expected %s to equal %d but got %d\n", tC.input, tC.output, got)
				t.Fail()
			}
		})
	}
}

func TestPrintTemperature(t *testing.T) {
	testCases := []struct {
		input  int16
		output string
	}{
		{
			output: "-12.5",
			input:  -125,
		},
		{
			output: "12.5",
			input:  125,
		},
		{
			output: "-2.5",
			input:  -25,
		},
		{
			output: "2.5",
			input:  25,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.output, func(t *testing.T) {
			var buf bytes.Buffer
			PrintTemperature(&buf, tC.input)
			if buf.String() != tC.output {
				t.Logf("Expected %d to equal %s but got %s\n", tC.input, tC.output, buf.String())
				t.Fail()
			}
		})
	}
}

func TestCut(t *testing.T) {
	input := []byte("Eindhoven;-15.3")
	city, temperature := Cut(input)
	assert.EqualValues(t, []byte("Eindhoven"), city)
	assert.EqualValues(t, []byte("-15.3"), temperature)
}
