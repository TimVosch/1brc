package logic_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"1brc/logic"
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
			got := logic.ParseTemperature(tC.input)
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
			logic.PrintTemperature(&buf, tC.input)
			if buf.String() != tC.output {
				t.Logf("Expected %d to equal %s but got %s\n", tC.input, tC.output, buf.String())
				t.Fail()
			}
		})
	}
}

func TestCut(t *testing.T) {
	input := []byte("Eindhoven;-15.3")
	city, temperature := logic.Cut(input)
	assert.EqualValues(t, []byte("Eindhoven"), city)
	assert.EqualValues(t, []byte("-15.3"), temperature)
}

func TestTrie(t *testing.T) {
	root := logic.NewTrie()
	station1 := root.Get([]byte{'H', 'E', 'Y'})
	assert.NotNil(t, station1)
	assert.NotNil(t, root.Root.Children[1])
	assert.NotNil(t, root.Root.Children[1].Children[1])
	assert.NotNil(t, root.Root.Children[1].Children[1].Children[1])

	station2 := root.Get([]byte{'H', 'E', 'Z'})
	assert.NotEqual(t, station1, station2)
	assert.NotNil(t, root.Root.Children[1])
	// Only with 0 terminated
	// assert.NotNil(t, root.Root.ChildrenDense[0].ChildrenDense[0].ChildrenDense[0].ChildrenDense[0])
}

func TestTestTrieItems(t *testing.T) {
	root := logic.NewTrie()
	root.Get([]byte{'A', 'B', 'A'})
	root.Get([]byte{'A', 'A', 'A'})
	root.Get([]byte{'H', 'E', 'A'})
	itemChan := make(chan *logic.Station, 1)
	items := make([]*logic.Station, 0)
	go func() {
		root.Items(itemChan)
		close(itemChan)
	}()
	for item := range itemChan {
		items = append(items, item)
	}
	assert.Len(t, items, 3)
	assert.Equal(t, "AAA", items[0].Name)
	assert.Equal(t, "ABA", items[1].Name)
	assert.Equal(t, "HEA", items[2].Name)
}

func TestBinarySearch(t *testing.T) {
	testCases := []struct {
		desc       string
		buffer     [0xff]uint8
		bufferSize uint8
		value      uint8
		expected   uint8
	}{
		{
			desc:       "with empty buffer",
			buffer:     [0xff]uint8{},
			bufferSize: 0,
			value:      10,
			expected:   0,
		},
		{
			desc:       "Start of buffer",
			buffer:     [0xff]uint8{1, 4, 7, 9, 15},
			bufferSize: 5,
			value:      1,
			expected:   0,
		},
		{
			desc:       "End of buffer",
			buffer:     [0xff]uint8{1, 4, 7, 9, 15},
			bufferSize: 5,
			value:      15,
			expected:   4,
		},
		{
			desc:       "Middle of buffer",
			buffer:     [0xff]uint8{1, 4, 7, 9, 15},
			bufferSize: 5,
			value:      7,
			expected:   2,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ix := logic.BinarySearch(tC.buffer, tC.bufferSize, tC.value)
			assert.EqualValues(t, tC.expected, ix)
		})
	}
}

func TestBinarySearchInsert(t *testing.T) {
	testCases := []struct {
		desc       string
		buffer     [0xff]uint8
		bufferSize uint8
		value      uint8
		expected   [0xff]uint8
	}{
		{
			desc:       "with empty buffer",
			buffer:     [0xff]uint8{},
			bufferSize: 0,
			value:      10,
			expected:   [0xff]uint8{10},
		},
		{
			desc:       "Start of buffer",
			buffer:     [0xff]uint8{3, 7, 9},
			bufferSize: 3,
			value:      2,
			expected:   [0xff]uint8{2, 3, 7, 9},
		},
		{
			desc:       "End of buffer",
			buffer:     [0xff]uint8{3, 7, 9},
			bufferSize: 3,
			value:      10,
			expected:   [0xff]uint8{3, 7, 9, 10},
		},
		{
			desc:       "Middle of buffer",
			buffer:     [0xff]uint8{3, 7, 9},
			bufferSize: 3,
			value:      6,
			expected:   [0xff]uint8{3, 6, 7, 9},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			logic.BinarySearchInsert(&tC.buffer, tC.bufferSize, tC.value)
			assert.EqualValues(t, tC.expected, tC.buffer)
		})
	}
}

func TestShift(t *testing.T) {
	buf := &[0xff]uint8{5, 3, 2, 6, 7, 4}
	logic.Shift(buf, 6, 0)
	assert.EqualValues(t, &[0xff]uint8{5, 5, 3, 2, 6, 7, 4}, buf)
}
