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

var count int
var groupsCount int
var pauseElapse int
var duration uint64

func main() {
	flag.IntVar(&count, "count", 1024, "count of timer")
	flag.IntVar(&groupsCount, "groupsCount", 10, "groups count of timer")
	flag.IntVar(&pauseElapse, "pauseElapse", 10, "pause elapse of two groups")
	flag.Uint64Var(&duration, "duration", 200, "duration you book")
	flag.Parse()

	errors := goTime()
	fmt.Println(errors)
	errors = buildinTime()
	fmt.Println(errors)

}

func buildinTime() (errors int64) {
	var wait sync.WaitGroup
	for i := 0; i < count; i++ {
		go func() {
			wait.Add(1)
			defer wait.Done()
			t := time.Now()
			timer := time.NewTimer(time.Millisecond * time.Duration(duration))
			select {
			case <-timer.C:
				elapse := int64(time.Since(t) - time.Millisecond*time.Duration(duration))
				atomic.AddInt64(&errors, elapse)
			}
		}()
		if i%(count/groupsCount) == 0 {
			time.Sleep(time.Millisecond * time.Duration(pauseElapse))
		}
	}
	wait.Wait()
	return
}

var errors int64

func goTime() int64 {
	gotimer := gotime.NewTimer()
	gotimer.SetInterval(time.Microsecond * 100)
	gotimer.Start()
	c := make(chan interface{}, 1024*1024)
	cls := make(chan struct{})
	for i := 0; i < 400; i++ {
		go func() {
			for {
				select {
				case data := <-c:
					t := data.(time.Time)
					elapse := int64(math.Abs(float64(time.Since(t)) - float64(time.Millisecond*time.Duration(duration))))
					atomic.AddInt64(&errors, elapse)
				case <-cls:
					return
				}
			}
		}()
	}

	for i := 0; i < count; i++ {
		gotimer.Time(duration, time.Now(), c, nil)
		if i%(count/groupsCount) == 0 {
			time.Sleep(time.Millisecond * time.Duration(pauseElapse))
		}
	}

	time.Sleep(time.Second * 10)
	close(cls)
	return errors
}

func handler(data interface{}) error {
	t := data.(time.Time)
	elapse := int64(math.Abs(float64(time.Since(t)) - float64(time.Millisecond*time.Duration(duration))))
	atomic.AddInt64(&errors, elapse)
	return nil
}
