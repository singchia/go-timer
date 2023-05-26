package main

import (
	"net/http"
	_ "net/http/pprof"

	gotime "github.com/singchia/go-timer/v2"
)

func main() {

	tw := gotime.NewTimer()
	tw.Start()
	ch := make(chan struct{}, 1024)

	http.HandleFunc("/bench", func(w http.ResponseWriter, req *http.Request) {
		tw.Add(1, gotime.WithData(w), gotime.WithHandler(func(data interface{}, err error) {
			rw := data.(http.ResponseWriter)
			rw.WriteHeader(http.StatusFound)
			ch <- struct{}{}
		}))
		<-ch
	})

	http.ListenAndServe(":6060", nil)
}
