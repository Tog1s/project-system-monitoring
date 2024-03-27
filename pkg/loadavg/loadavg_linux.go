//go:build linux
// +build linux

package loadavg

import (
	"os"
	"strconv"
	"strings"
)

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

	loadAvg1, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return nil, err
	}

	loadAvg5, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return nil, err
	}

	loadAvg15, err := strconv.ParseFloat(values[2], 64)
	if err != nil {
		return nil, err
	}

	return &LoadAverage{
		LoadAvg1:  loadAvg1,
		LoadAvg5:  loadAvg5,
		LoadAvg15: loadAvg15,
	}, nil
}

func readFromFile() ([]string, error) {
	f, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	return strings.Split(string(f), " "), nil
}
