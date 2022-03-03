package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
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
	flag.IntVar(&ratio, "r", 5, "ratio of tick") // 系数
	flag.Int64Var(&grain, "g", 1000, "grain of tick")
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

func goTime(count int) (tolerance int64) {
	tw := gotime.NewTimer(gotime.WithTimeInterval(100 * time.Millisecond))
	tw.Start()
	var wait sync.WaitGroup
	for i := 0; i < count; i++ {
		wait.Add(1)
		t := time.Now()
		d := time.Duration((rand.Int63n(int64(count)) + 1) * grain)
		tw.Time(d, gotime.WithData(t), gotime.WithHandler(func(data interface{}) error {
			defer wait.Done()
			t = data.(time.Time)
			thisToler := int64(math.Abs(float64(time.Since(t) - d)))
			atomic.AddInt64(&tolerance, thisToler)
			return nil
		}))
	}
	wait.Wait()
	tw.Stop()
	return tolerance / int64(count)
}

func buildinTime(count int) (tolerance int64) {
	var wait sync.WaitGroup
	for i := 0; i < count; i++ {
		wait.Add(1)
		t := time.Now()
		d := time.Duration((rand.Int63n(int64(count) + 1)) * grain)
		timer := time.NewTimer(d)
		go func(t time.Time, timer *time.Timer) {
			defer wait.Done()
			<-timer.C
			thisToler := int64(math.Abs(float64(time.Since(t) - d)))
			atomic.AddInt64(&tolerance, thisToler)
		}(t, timer)
	}
	wait.Wait()
	return tolerance / int64(count)
}
