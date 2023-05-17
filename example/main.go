package main

import (
	"fmt"
	"time"

	timer "github.com/singchia/go-timer"
)

func main() {
	t := timer.NewTimer()
	t.Start()
	ch := make(chan interface{})
	old := time.Now()

	tick := t.Time(5, timer.WithData(1), timer.WithChan(ch))
	<-tick.Channel()
	elapse := time.Now().Sub(old)
	fmt.Printf("time diff: %d", elapse)
}
