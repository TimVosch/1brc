package logic

import (
	"io"
	"os"
	"strconv"
)

// Saint-Pierre;11.4
// Ho Chi Minh City;29.3
// Kyoto;19.8
// Atlanta;1.9
// Ljubljana;25.0
// Dakar;36.4
// Palmerston North;10.9
// Cape Town;26.9
// Pyongyang;8.3
// Riga;6.3
//
// {Abha=-23.0/18.0/59.2, Abidjan=-16.2/26.0/67.3, Abéché=-10.0/29.4/69.0, Accra=-10.1/26.4/66.4, Addis Ababa=-23.7/16.0/67.0, Adelaide=-27.8/17.3/58.5, ...}
type Station struct {
	Name  string
	Min   int16
	Max   int16
	Total int64
	Count int64
}

type StationName [100]byte

const (
	readBufferSize = 1 << 24
	maxLineLength  = 110 // 100 name + 1 ; + 1 \n + 5 -12.3
)

var stationTrie = NewTrie()

func Run(path string) error {
	readContents(path)

	// Output
	stations := make(chan *Station, 500)
	go func() {
		stationTrie.Items(stations)
		close(stations)
	}()

	file, _ := os.Create("output.txt")
	file.Write([]byte("{"))
	for station := range stations {
		file.Write([]byte(station.Name + "="))
		PrintStation(file, station)
	}
	file.Seek(-2, 1)
	file.Write([]byte("}"))
	file.Close()

	return nil
}

func process(line []byte) {
	stationName, temperatureString := Cut(line)
	temperature := ParseTemperature(temperatureString)
	station := stationTrie.Get(stationName)
	station.Count++
	station.Total += int64(temperature)
	station.Min = Min(station.Min, temperature)
	station.Max = Max(station.Max, temperature)
}

func readContents(path string) {
	file, _ := os.Open(path)
	var lineStart, lineEnd, n int
	var tail int
	buffer := make([]byte, readBufferSize)
	n, _ = file.Read(buffer)
	for n > 0 {
		for tail = 0; buffer[n-tail-1] != '\n'; tail++ {
		}

		for lineStart = 0; lineStart < n-tail; {
			for lineEnd = lineStart; buffer[lineEnd] != '\n'; lineEnd++ {
			}
			process(buffer[lineStart:lineEnd])
			lineStart = lineEnd + 1
		}

		copy(buffer, buffer[n-tail:n])
		n, _ = file.Read(buffer[tail:])
		n += tail
	}
	file.Close()
}

func Cut(line []byte) ([]byte, []byte) {
	var i int
	for i = len(line) - 4; line[i] != ';'; i-- {
	}
	return line[:i], line[i+1:]
}

func PrintTemperature(f io.Writer, value int16) {
	str := []byte(strconv.Itoa(int(value)))
	str = append(str, str[len(str)-1])
	str[len(str)-2] = '.'
	f.Write(str)
}

func PrintStation(f io.Writer, station *Station) {
	PrintTemperature(f, station.Min)
	f.Write([]byte{'/'})
	PrintTemperature(f, int16(int64(station.Total)/station.Count))
	f.Write([]byte{'/'})
	PrintTemperature(f, station.Max)
	f.Write([]byte{',', ' '})
}

func Min(a, b int16) int16 {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int16) int16 {
	if a < b {
		return b
	}
	return a
}

func ParseTemperature(str []byte) int16 {
	length := len(str)
	out := int16(str[length-1]-'0') + int16(str[length-3]-'0')*10

	// Check if there's a tens place for the integer part or if it's negative.
	// Since we know the structure, we can infer this directly based on the length and content.
	if str[0] == '-' {
		// When the string starts with '-', it's negative.
		// Depending on length, we know if there's a tens digit or just a units digit.
		if length == 5 {
			out += int16(str[1]-'0') * 100
		}
		out = -out
	} else if length == 4 {
		// If it's length 4 and not negative, it must have a tens place.
		out += int16(str[0]-'0') * 100
	}

	return out
}

type TrieChildren [0xff]*TrieNode

type TrieNode struct {
	Parent         *TrieNode
	ChildrenDense  TrieChildren
	ChildrenSparse [0xff]uint8
	NextIndex      uint8
	IsValueNode    bool
	Value          Station
}

type Trie struct {
	Root TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		Root: TrieNode{
			Parent:    nil,
			NextIndex: 1,
		},
	}
}

func (trie *Trie) Get(name []byte) *Station {
	nameIndex := 0
	nameLength := len(name)
	var next *TrieNode
	var key uint8
	node := &trie.Root
	for {
		key = name[nameIndex]
		next = node.ChildrenDense[node.ChildrenSparse[key]]

		// Our name ends so the next node must be(come) our value node
		if nameIndex == nameLength-1 || key == 0 {
			if next == nil {
				next = &TrieNode{
					Parent: node,
					Value: Station{
						Name: string(name),
					},
					IsValueNode: true,
					NextIndex:   1,
				}
				node.ChildrenDense[node.NextIndex] = next
				node.ChildrenSparse[key] = node.NextIndex
				node.NextIndex++
			}
			// TODO: It might happen that a name is part of a bigger name, but a different
			// city, therefor we MUST check if the existing next is not a value node yet
			// ie (next != nil && !next.IsEnd)
			return &next.Value
		}

		// Not returned yet, and next exists so continue traversing
		if next != nil {
			node = next
			nameIndex++
			continue
		}

		// Next does not exist, complete the path
		next = &TrieNode{
			Parent:    node,
			NextIndex: 1,
		}
		node.ChildrenDense[node.NextIndex] = next
		node.ChildrenSparse[key] = node.NextIndex
		node.NextIndex++
		node = next
		nameIndex++
	}
}

func (node *TrieNode) Items(out chan<- *Station) {
	for i := uint8(1); i < node.NextIndex; i++ {
		child := node.ChildrenDense[i]
		if child.IsValueNode {
			out <- &child.Value
		} else {
			child.Items(out)
		}
	}
}

func (trie *Trie) Items(out chan<- *Station) {
	trie.Root.Items(out)
}

func BinarySearch(buf [0xff]uint8, size, value uint8) uint8 {
	var low, mid uint8
	high := size

	for low <= high {
		mid = (low + high) / 2
		if buf[mid] > value {
			high = mid - 1
		} else if buf[mid] < value {
			low = mid + 1
		} else {
			return mid
		}
	}
	return 0
}

func Shift(buf *[0xff]uint8, size, start uint8) {
	for ix := size; ix > start; ix-- {
		buf[ix] = buf[ix-1]
	}
}

func BinarySearchInsert(buf *[0xff]uint8, size, value uint8) {
	var low, mid uint8
	high := size - 1
	potentialIndex := size

	if size == 0 {
		buf[0] = value
		return
	}
	if buf[0] > value {
		potentialIndex = 0
	} else {
		for low < high {
			mid = (low + high) / 2
			if buf[mid] > value {
				potentialIndex = mid
				high = mid - 1
			} else if buf[mid] < value {
				low = mid + 1
			} else {
				panic("Duplicate key")
			}
		}
	}

	Shift(buf, size, potentialIndex)
	buf[potentialIndex] = value
}
