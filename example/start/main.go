package main

import (
	"log"
	"time"

	timer "github.com/singchia/go-timer/v2"
)

func main() {
	t1 := time.Now()

	t := timer.NewTimer()
	tick := t.Add(time.Second)
	<-tick.C()

	log.Printf("time elapsed: %fs\n", time.Now().Sub(t1).Seconds())
}
