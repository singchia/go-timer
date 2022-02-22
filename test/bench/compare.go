package main

import (
	"flag"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	gotime "github.com/singchia/go-timer"
)

var points int
var grain int64
var ratio int
var start int

func main() {
	flag.IntVar(&start, "s", 1, "start")
	flag.IntVar(&points, "p", 10, "points")
	flag.IntVar(&ratio, "r", 5, "ratio of tick")      // 系数
	flag.Int64Var(&grain, "g", 1000, "grain of tick") // default 1s
	flag.Parse()

	index := 0
	for index < points {
		count := start * int(math.Pow(float64(ratio), float64(index)))
		tolerance := goTime(count)
		fmt.Println(count, tolerance)
		index++
	}

	index = 0
	for index < points {
		count := start * int(math.Pow(float64(ratio), float64(index)))
		tolerance := buildinTime(count)
		fmt.Println(count, tolerance)
		index++
	}
}

func buildinTime(count int) (tolerance int64) {
	var wait sync.WaitGroup
	for i := 0; i < count; i++ {
		wait.Add(1)
		t := time.Now()
		timer := time.NewTimer(time.Duration(grain * 1000))
		go func(t time.Time, timer *time.Timer) {
			defer wait.Done()
			<-timer.C
			thisToler := int64(math.Abs(float64(time.Since(t) - time.Duration(grain*1000))))
			atomic.AddInt64(&tolerance, thisToler)
		}(t, timer)
	}
	wait.Wait()
	return tolerance / int64(count)
}

func goTime(count int) (tolerance int64) {
	var wait sync.WaitGroup
	tw := gotime.NewTimer()
	tw.SetInterval(time.Duration(grain * 1000))
	tw.Start()
	for i := 0; i < count; i++ {
		wait.Add(1)
		t := time.Now()
		tw.Time(1, t, nil, func(data interface{}) error {
			defer wait.Done()
			t = data.(time.Time)
			thisToler := int64(math.Abs(float64(time.Since(t) - time.Duration(grain*1000))))
			atomic.AddInt64(&tolerance, thisToler)
			return nil
		})
	}
	wait.Wait()
	tw.Stop()
	return tolerance / int64(count)
}
