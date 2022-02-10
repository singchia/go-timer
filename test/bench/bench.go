package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"sync"
	"time"

	gotime "github.com/singchia/go-timer"
)

var count int

func main() {
	flag.IntVar(&count, "count", 100000, "count of tick")
	flag.Parse()
	goTime()
}

func goTime() {
	tw := gotime.NewTimer()
	tw.Start()
	wait := new(sync.WaitGroup)
	for i := 0; i < count; i++ {
		time.Sleep(time.Microsecond)
		wait.Add(1)
		tw.Time(1, wait, nil, func(data interface{}) error {
			wait.Done()
			return nil
		})
	}
	wait.Wait()
	fmt.Println("the end")
}
