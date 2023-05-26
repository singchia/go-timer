package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"

	timer "github.com/singchia/go-timer"
)

func main() {
	server := &http.Server{
		Addr:    ":6060",
		Handler: nil,
	}
	go server.ListenAndServe()

	t := timer.NewTimer()
	t.Start()
	ch := make(chan interface{})
	old := time.Now()
	t.Time(10, 1, ch, nil)
	<-ch
	elapse := time.Now().Sub(old)
	t.Stop()

	fmt.Printf("time diff: %d", elapse)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)
	go func() {
		<-sigs
		done <- true
	}()
	<-done

	server.Close()
	signal.Reset()

	time.Sleep(3 * time.Second)
	fmt.Println(runtime.NumGoroutine())
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}
