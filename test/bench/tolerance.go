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

var count int
var grain int64
var iterateTolerance bool

func main() {
	flag.IntVar(&count, "c", 1000000, "count of tick")
	flag.Int64Var(&grain, "g", 1000, "grain of tick in Microsecond")
	flag.BoolVar(&iterateTolerance, "t", false, "print iterate tolerance")
	flag.Parse()

	tolerance := buildinTime()
	fmt.Println(time.Duration(tolerance))
	tolerance = goTime()
	fmt.Println(time.Duration(tolerance))
}

func buildinTime() (tolerance int64) {
	var wait sync.WaitGroup
	for i := 0; i < count; i++ {
		wait.Add(1)
		t := time.Now()
		timer := time.NewTimer(time.Duration(grain * 1000))
		go func(t time.Time, timer *time.Timer) {
			defer wait.Done()
			<-timer.C
			thisToler := int64(math.Abs(float64(time.Since(t) - time.Duration(grain*1000))))
			if iterateTolerance {
				fmt.Printf("buildin time %d toler: %d\n", i, thisToler)
			}
			atomic.AddInt64(&tolerance, thisToler)
		}(t, timer)
	}
	wait.Wait()
	return tolerance / int64(count)
}

func goTime() (tolerance int64) {
	var wait sync.WaitGroup
	tw := gotime.NewTimer()
	tw.Start()
	for i := 0; i < count; i++ {
		wait.Add(1)
		t := time.Now()
		tw.Time(1, gotime.WithData(t), gotime.WithHandler(func(data interface{}) error {
			defer wait.Done()
			t = data.(time.Time)
			thisToler := int64(math.Abs(float64(time.Since(t) - time.Duration(grain*1000))))
			if iterateTolerance {
				fmt.Printf("buildin time %d toler: %d\n", i, thisToler)
			}
			atomic.AddInt64(&tolerance, thisToler)
			return nil
		}))
	}
	wait.Wait()
	return tolerance / int64(count)
}

// we cannot get tolerance from no-context timer
func buildinReflectTime() (tolerance int64) {
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
				wait.Done()
			}
		}
	}()
	wait.Wait()
	return tolerance / int64(count)
}
