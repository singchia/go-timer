package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
	"time"

	gotime "github.com/singchia/go-timer"
)

func main() {

	tw := gotime.NewTimer(gotime.WithTimeInterval(time.Millisecond))
	tw.Start()
	ch := make(chan struct{}, 1024)

	http.HandleFunc("/timer", func(w http.ResponseWriter, req *http.Request) {
		tw.Time(500*time.Millisecond, gotime.WithData(w), gotime.WithHandler(func(data interface{}) error {
			ch <- struct{}{}
			return nil
		}))
		<-ch
	})

	http.HandleFunc("/bench", func(w http.ResponseWriter, req *http.Request) {
		value := req.URL.Query().Get("count")
		count, err := strconv.Atoi(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		go func() {
			f, _ := os.OpenFile("cpu.profile", os.O_CREATE|os.O_RDWR, 0644)
			defer f.Close()
			defer fmt.Println("finished")

			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
			wg := new(sync.WaitGroup)
			wg.Add(count)
			for i := 0; i < count; i++ {
				tw.Time(500*time.Millisecond, gotime.WithData(i), gotime.WithHandler(func(data interface{}) error {
					wg.Done()
					reminder := (data.(int)) % 10000
					if reminder == 0 {
						fmt.Println("10000 done")
					}
					return nil
				}))
			}
			wg.Wait()
		}()
	})

	http.ListenAndServe(":6060", nil)
}
