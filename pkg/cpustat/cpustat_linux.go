package cpustat

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type CPUStat struct {
	User   float64
	System float64
	Idle   float64
}

var idle string

func (c CPUStat) String() string {
	return fmt.Sprintf("User: %f System: %f Idle: %f", c.User, c.System, c.Idle)
}

func Get() (*CPUStat, error) {
	cmd := exec.Command("top", "-b", "-n1")
	r, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	fields := strings.FieldsFunc(string(r), func(r rune) bool {
		return r == '\n'
	})

	cpuStat := strings.Fields(fields[2])

	cpuUser, err := strconv.ParseFloat(strings.ReplaceAll(cpuStat[1], ",", "."), 64)
	if err != nil {
		return nil, err
	}

	cpuSystem, err := strconv.ParseFloat(strings.ReplaceAll(cpuStat[3], ",", "."), 64)
	if err != nil {
		return nil, err
	}

	if cpuStat[7] == "id," {
		idle = cpuStat[6]
	}

	cpuIdle, err := strconv.ParseFloat(strings.ReplaceAll(idle, ",", "."), 64)
	if err != nil {
		return nil, err
	}

	return &CPUStat{
		User:   cpuUser,
		System: cpuSystem,
		Idle:   cpuIdle,
	}, nil
}
