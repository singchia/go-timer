package main

import (
	"log"
	"sync"
	"time"

	timer "github.com/singchia/go-timer/v2"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	t := timer.NewTimer()
	t.Add(time.Second, timer.WithData(time.Now()), timer.WithHandler(func(event *timer.Event) {
		defer wg.Done()
		log.Printf("time elapsed: %fs\n", time.Now().Sub(event.Data.(time.Time)).Seconds())
	}))

	wg.Wait()
}
