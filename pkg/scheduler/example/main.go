package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	scheduler "github.com/singchia/go-scheduler"
)

func main() {

	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	sch := scheduler.NewScheduler()
	sch.Interval = time.Millisecond * 100

	sch.SetMonitor(SchedulerMonitor)
	sch.SetMaxRate(0.95)
	sch.SetMaxGoroutines(5000)
	sch.StartSchedule()

	var val int64
	for i := 0; i < 100*10000*10; i++ {
		sch.PublishRequest(&scheduler.Request{Data: val, Handler: SchedulerHandler})
		atomic.AddInt64(&val, 1)
	}
	time.Sleep(time.Second * 10)
	fmt.Printf("maxValue: %d\n", maxValue)

	sch.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)
	go func() {
		<-sigs
		done <- true
	}()
	<-done

	fmt.Println(runtime.NumGoroutine())
}

var maxValue int64 = 0

func SchedulerHandler(data interface{}) {
	val, ok := data.(int64)
	if ok {
		if val > maxValue {
			maxValue = val
		}
	}
}

func SchedulerMonitor(incomingReqsDiff, processedReqsDiff, diff, currentGotoutines int64) {
	fmt.Printf("%d, %d, %d, %d\n", incomingReqsDiff, processedReqsDiff, diff, currentGotoutines)
}
