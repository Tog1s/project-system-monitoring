//go:build darwin
// +build darwin

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
	cmd := exec.Command("iostat")
	// cmd.Env = append(cmd.Environ(), "LC_NUMERIC=en_GB.UTF-8")

	r, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	cpuStat := strings.Fields(string(r))

	cpuUser, err := strconv.ParseFloat(cpuStat[16], 64)
	if err != nil {
		return nil, err
	}

	cpuSystem, err := strconv.ParseFloat(cpuStat[17], 64)
	if err != nil {
		return nil, err
	}

	cpuIdle, err := strconv.ParseFloat(cpuStat[18], 64)
	if err != nil {
		return nil, err
	}

	return &CPUStat{
		User:   cpuUser,
		System: cpuSystem,
		Idle:   cpuIdle,
	}, nil
}
