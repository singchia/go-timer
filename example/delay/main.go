package main

import (
	"log"
	"math"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/singchia/go-timer"
)

func main() {
	server := &http.Server{
		Addr:    ":6060",
		Handler: nil,
	}
	go server.ListenAndServe()
	n := 100000
	delay := 14
	tw := timer.NewTimer()
	tw.Start()
	fired := int32(0)
	for i := 0; i < n; i++ {
		second := uint64(rand.Intn(10) + 1)
		tick, err := tw.Time(second, time.Now(), nil, func(data interface{}) error {
			elapse := time.Since(data.(time.Time).Add(time.Duration(second) * time.Second)).Seconds()
			abs := int(math.Abs(elapse))
			if abs < delay-1 || abs > delay+1 {
				log.Println(elapse)
			}
			atomic.AddInt32(&fired, 1)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		tick.Delay(uint64(delay))
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-sigs:
			goto END
		case <-time.NewTimer(time.Second).C:
			log.Println(n, fired)
		}
	}
END:
	log.Println(n, fired)
}
