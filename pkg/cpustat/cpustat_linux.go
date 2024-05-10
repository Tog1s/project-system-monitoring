//go:build linux
// +build linux

package cpustat

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func (c CPUStat) String() string {
	return fmt.Sprintf("User: %f System: %f Idle: %f", c.User, c.System, c.Idle)
}

func Get() (*CPUStat, error) {
	// cmd := exec.Command("top", "-b", "-n1")

	cmd := exec.Command("mpstat")
	cmd.Env = append(cmd.Environ(), "LC_NUMERIC=en_GB.UTF-8")

	r, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	fields := strings.FieldsFunc(string(r), func(r rune) bool {
		return r == '\n'
	})

	cpuStat := strings.Fields(fields[2])

	cpuUser, err := strconv.ParseFloat(cpuStat[2], 64)
	if err != nil {
		return nil, err
	}

	cpuSystem, err := strconv.ParseFloat(cpuStat[4], 64)
	if err != nil {
		return nil, err
	}

	cpuIdle, err := strconv.ParseFloat(cpuStat[11], 64)
	if err != nil {
		return nil, err
	}

	return &CPUStat{
		User:   cpuUser,
		System: cpuSystem,
		Idle:   cpuIdle,
	}, nil
}
