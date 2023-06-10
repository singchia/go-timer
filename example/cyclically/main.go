package main

import (
	"log"
	"time"

	timer "github.com/singchia/go-timer/v2"
)

func main() {
	t1 := time.Now()

	t := timer.NewTimer()
	tick := t.Add(time.Second, timer.WithCyclically())
	for {
		<-tick.C()
		t2 := time.Now()
		log.Printf("time elapsed: %fs\n", t2.Sub(t1).Seconds())
		t1 = t2
	}
}
