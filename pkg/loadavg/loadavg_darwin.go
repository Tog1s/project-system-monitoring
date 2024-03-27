//go:build darwin
// +build darwin

package loadavg

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func Get() (*LoadAverage, error) {
	type loadavg struct {
		load  [3]uint32
		scale int
	}

	b, err := unix.SysctlRaw("vm.loadavg")
	if err != nil {
		return nil, err
	}

	load := *(*loadavg)(unsafe.Pointer(&b[0]))
	scale := float64(load.scale)

	metric := &LoadAverage{
		LoadAvg1:  float64(load.load[0]) / scale,
		LoadAvg5:  float64(load.load[1]) / scale,
		LoadAvg15: float64(load.load[2]) / scale,
	}

	return metric, nil
}
