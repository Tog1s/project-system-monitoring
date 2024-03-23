package loadavg

import (
	"os"
	"strconv"
	"strings"
)

type LoadAverage struct {
	OneMinutes     float32
	FiveMinutes    float32
	FifteenMinutes float32
	ProcRuning     int
	ProcAll        int
	LastProcessId  int
}

func Read() LoadAverage {
	f, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		panic(err)
	}
	var loadAvg LoadAverage
	data := strings.Split(string(f), " ")

	oneMinutes, _ := strconv.ParseFloat(data[0], 32)
	loadAvg.OneMinutes = float32(oneMinutes)

	return loadAvg
}
