package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"1brc/logic"
)

var (
	cpuprofile       = flag.Bool("profile", false, "Whether to profile")
	measurementsPath = flag.String("measurements", "ramdisk/measurements.txt", "Where to read measurements from")
)

func main() {
	flag.Parse()
	if *cpuprofile {
		f, err := os.Create("cpu.pprof")
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
	if err := logic.Run(*measurementsPath); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
	}
}
