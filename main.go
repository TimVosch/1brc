package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	cpuprofile       = flag.Bool("profile", false, "Whether to profile")
	measurementsPath = flag.String("measurements", "ramdisk/measurements.txt", "Where to read measurements from")
)

func main() {
	flag.Parse()
	if *cpuprofile {
		f, err := os.Create("profiles/" + time.Now().Format("200601021504") + ".pprof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		log.Println("Profiling")
		defer pprof.StopCPUProfile()
	}
	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
	}
}

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
	Min   int16
	Max   int16
	Total int16
	Count int64
}

func Run() error {
	stations := make(map[string]Station, 10000)

	file, err := os.Open(*measurementsPath)
	if err != nil {
		return err
	}

	// Process
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()
		stationName, temperatureString, _ := strings.Cut(line, ";")
		temperature := ParseTemperature(temperatureString)
		station, ok := stations[stationName]
		if !ok {
			station = Station{
				Min:   temperature,
				Max:   temperature,
				Total: temperature,
				Count: 1,
			}
		} else {
			station.Count++
			station.Total += temperature

			station.Min = Min(station.Min, temperature)
			station.Max = Max(station.Max, temperature)
		}
		stations[stationName] = station
	}
	file.Close()

	// Sort stationNames
	stationNames := make([]string, len(stations))
	ix := 0
	for stationName := range stations {
		stationNames[ix] = stationName
		ix++
	}
	slices.Sort(stationNames)

	// Output
	file, _ = os.OpenFile("output.txt", os.O_CREATE|os.O_WRONLY, 644)
	file.Write([]byte("{"))
	for _, stationName := range stationNames {
		file.Write([]byte(stationName + "="))
		PrintStation(file, stations[stationName])
	}
	file.Seek(-2, 1)
	file.Write([]byte("}"))
	file.Close()

	return nil
}

func PrintTemperature(f io.Writer, value int16) {
	str := []byte(strconv.Itoa(int(value)))
	str = append(str, str[len(str)-1])
	str[len(str)-2] = '.'
	f.Write(str)
}

func PrintStation(f io.Writer, station Station) {
	PrintTemperature(f, station.Min)
	f.Write([]byte{'/', ' '})
	PrintTemperature(f, int16(int64(station.Total)/station.Count))
	PrintTemperature(f, station.Max)
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

// __0 1 2 3 4
//   - 1 2 . 3
//     1 2 . 3
//   - 2 . 3
//     2 . 3
//func ParseTemperature2(str string) int16 {
//	var out int16
//	chars := []byte(str)
//	end := len(chars) - 1
//	// These three characters are always the same digit dot digit
//	out += int16((chars[end] - 0x30))
//	// out += int16((chars[end-1] - 0x30) * 10) // ignore '.'
//	out += int16((chars[end-2] - 0x30) * 10)
//	if end-3 < 0 {
//		return out
//	}
//	// The following digits are optional: ( [negation + digit] | [negation | digit] )
//	if chars[end-3] == '-' {
//		return out * -1
//	}
//	out += int16((chars[end-3] - 0x30) * 100)
//	if end-4 >= 0 && chars[end-4] == '-' {
//		return out * -1
//	}
//
//	return out
//}

func ParseTemperature(str string) int16 {
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
