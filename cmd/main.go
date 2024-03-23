package main

import (
	"fmt"

	"github.com/tog1s/project-system-monitoring/pkg/loadavg"
)

func main() {
	fmt.Println(loadavg.Read())
}
