package main

import (
	"fmt"
	"time"

	timer "github.com/singchia/go-timer/v2"
)

func main() {
	t := timer.NewTimer(timer.WithTimeInterval(time.Millisecond))
	ch := make(chan *timer.Event)
	old := time.Now()

	tick := t.Add(5, timer.WithData(1), timer.WithChan(ch))
	<-tick.C()
	elapse := time.Now().Sub(old).Seconds()
	fmt.Printf("time diff: %f", elapse)
}
