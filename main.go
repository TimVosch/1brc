package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
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
	Min   float64
	Max   float64
	Total float64
	Count int64
}

func Run() error {
	stations := make(map[string]Station, 10000)

	file, err := os.Open("measurements.txt")
	if err != nil {
		return err
	}

	// Process
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()
		stationName, temperatureString, _ := strings.Cut(line, ";")
		temperature, _ := strconv.ParseFloat(temperatureString, 64)
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
			station.Min = math.Min(station.Min, temperature)
			station.Max = math.Max(station.Max, temperature)
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
	file, _ = os.OpenFile("output.txt", os.O_CREATE|os.O_WRONLY, 0)
	file.Write([]byte("{"))
	for _, stationName := range stationNames {
		station := stations[stationName]
		fmt.Fprintf(file, "%s=%.1f/%.1f/%.1f, ", stationName, station.Min, station.Total/float64(station.Count), station.Max)
	}
	file.Seek(-2, 1)
	file.Write([]byte("}"))
	file.Close()

	return nil
}
