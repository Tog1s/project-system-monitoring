package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/tog1s/project-system-monitoring/pkg/loadavg"
)

func reader(ch <-chan loadavg.LoadAverage, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		val, ok := <-ch
		if !ok {
			fmt.Println("Channel Closed")
			return
		}
		fmt.Println(val)
	}
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan loadavg.LoadAverage)

	wg.Add(1)
	go reader(ch, &wg)

	for {
		loadAvg, err := loadavg.Get()
		if err != nil {
			fmt.Println(err)
		}
		ch <- *loadAvg
		time.Sleep(1 * time.Second)
	}
}
