package loadavg

import (
	"os"
	"strconv"
	"strings"
)

type LoadAverage struct {
	OneMinutes     float64
	FiveMinutes    float64
	FifteenMinutes float64
	// ProcRuning     int
	// ProcAll        int
	// LastProcessId  int
}

func Get() (*LoadAverage, error) {
	stats, err := loadAvgFromFile()
	if err != nil {
		return nil, err
	}
	return stats, err
}

func loadAvgFromFile() (*LoadAverage, error) {
	values, err := readFromFile()
	if err != nil {
		return nil, err
	}

	oneMinutes, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return nil, err
	}

	fiveMinutes, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return nil, err
	}

	fifteenMinutes, err := strconv.ParseFloat(values[2], 64)
	if err != nil {
		return nil, err
	}

	return &LoadAverage{
		OneMinutes:     oneMinutes,
		FiveMinutes:    fiveMinutes,
		FifteenMinutes: fifteenMinutes,
	}, nil
}

func readFromFile() ([]string, error) {
	f, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	return strings.Split(string(f), " "), nil
}
