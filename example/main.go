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
	t.Time(5, 1, ch, nil)
	<-ch
	elapse := time.Now().Sub(old)
	fmt.Printf("time diff: %d", elapse)
}
