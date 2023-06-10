package main

import (
	"log"
	"time"

	timer "github.com/singchia/go-timer/v2"
)

func main() {
	t := timer.NewTimer()
	tick := t.Add(time.Second, timer.WithCyclically())
	for {
		<-tick.C()
		log.Println(time.Now())
	}
}
