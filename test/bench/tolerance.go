package main

import (
	"flag"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

var count int
var grain int64

func main() {
	flag.IntVar(&count, "c", 1000000, "count of tick")
	flag.Int64Var(&grain, "g", 1000, "grain of tick in Microsecond")
	flag.Parse()

	tolerance := buildinTime()
	fmt.Println(time.Duration(tolerance).Seconds())
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
			elapse := int64(math.Abs(float64(time.Since(t) - time.Duration(grain*1000))))
			atomic.AddInt64(&tolerance, elapse)
		}(t, timer)
	}
	wait.Wait()
	return
}
