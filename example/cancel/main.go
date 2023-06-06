package main

import (
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync/atomic"
	"time"

	"github.com/singchia/go-timer/v2"
)

func main() {
	server := &http.Server{
		Addr:    ":6060",
		Handler: nil,
	}
	go server.ListenAndServe()
	n := 100000
	tw := timer.NewTimer(timer.WithTimeInterval(time.Millisecond))
	fired := int32(0)
	for i := 0; i < n; i++ {
		second := time.Duration(rand.Intn(10)+1) * time.Second
		tick := tw.Add(second, timer.WithData(time.Now()), timer.WithHandler(func(event *timer.Event) {
			atomic.AddInt32(&fired, 1)
		}))
		err := tick.Cancel()
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	sigs := make(chan os.Signal, 1)
	tick := tw.Add(time.Second, timer.WithCyclically())
	for {
		select {
		case <-sigs:
			goto END
		case <-tick.C():
			log.Println(n, fired)
		}
	}
END:
	log.Println(n, fired)
}
