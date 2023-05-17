package main

import (
	"flag"
	"fmt"
	"math"
	"reflect"
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

func buildinTime(count int) (tolerance int64) {
	var wait sync.WaitGroup
	for i := 0; i < count; i++ {
		wait.Add(1)
		t := time.Now()
		timer := time.NewTimer(time.Duration(grain))
		go func(t time.Time, timer *time.Timer) {
			defer wait.Done()
			<-timer.C
			thisToler := int64(math.Abs(float64(time.Since(t) - time.Duration(grain))))
			atomic.AddInt64(&tolerance, thisToler)
		}(t, timer)
	}
	wait.Wait()
	return tolerance / int64(count)
}

func goTime(count int) (tolerance int64) {
	var wait sync.WaitGroup
	tw := gotime.NewTimer(gotime.WithTimeInterval(time.Microsecond))
	tw.Start()
	for i := 0; i < count; i++ {
		wait.Add(1)
		t := time.Now()
		tw.Time(time.Duration(grain), gotime.WithData(t), gotime.WithHandler(func(data interface{}) error {
			defer wait.Done()
			t = data.(time.Time)
			thisToler := int64(math.Abs(float64(time.Since(t) - time.Duration(grain))))
			atomic.AddInt64(&tolerance, thisToler)
			return nil
		}))
	}
	wait.Wait()
	tw.Stop()
	return tolerance / int64(count)
}

func buildinReflectTime(count int) (tolerance int64) {
	var wait sync.WaitGroup
	cases := []reflect.SelectCase{}
	mu := new(sync.RWMutex)
	for i := 0; i < count; i++ {
		mu.Lock()
		wait.Add(1)
		timer := time.NewTimer(time.Duration(grain * 1000))
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(timer.C)})
		mu.Unlock()
	}

	go func() {
		for {
			mu.RLock()
			i, value, ok := reflect.Select(cases)
			mu.RUnlock()
			if !ok {
				mu.Lock()
				cases = append(cases[:i], cases[i+1:]...)
				mu.Unlock()
			} else {
				fmt.Printf("%v\n", value)
				wait.Done()
			}
		}
	}()
	wait.Wait()
	return tolerance / int64(count)
}
