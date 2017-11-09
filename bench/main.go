package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	errors := buildinTime()
	fmt.Print(errors)
}

func buildinTime() (errors int64) {
	var wait sync.WaitGroup
	for i := 0; i < 1024; i++ {
		go func() {
			wait.Add(1)
			defer wait.Done()
			t := time.Now()
			timer := time.NewTimer(time.Millisecond * 200)
			select {
			case <-timer.C:
				elapse := int64(time.Since(t) - time.Millisecond*200)
				atomic.AddInt64(&errors, elapse)
			}
		}()
	}
	wait.Wait()
	return
}

func goTime() (errors int64) {

}
