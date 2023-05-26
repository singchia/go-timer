package main

import (
	"fmt"
	"time"

	timer "github.com/singchia/go-timer/v2"
)

func main() {
	t := timer.NewTimer()
	t.Start()
	ch := make(chan interface{})
	old := time.Now()

	tick := t.Add(5, timer.WithData(1), timer.WithChan(ch))
	<-tick.Chan()
	elapse := time.Now().Sub(old)
	fmt.Printf("time diff: %d", elapse)
}
