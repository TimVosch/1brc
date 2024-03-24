package main

import (
	"bytes"
	"testing"
)

func TestParseTemperature(t *testing.T) {
	testCases := []struct {
		input  string
		output int16
	}{
		{
			input:  "-12.5",
			output: -125,
		},
		{
			input:  "12.5",
			output: 125,
		},
		{
			input:  "-2.5",
			output: -25,
		},
		{
			input:  "2.5",
			output: 25,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
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
