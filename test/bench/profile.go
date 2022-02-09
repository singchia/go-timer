package main

import (
	"net/http"
	_ "net/http/pprof"

	gotime "github.com/singchia/go-timer"
)

func main() {

	tw := gotime.NewTimer()
	tw.Start()
	ch := make(chan struct{}, 1024)

	http.HandleFunc("/bench", func(w http.ResponseWriter, req *http.Request) {
		tw.Time(1, w, nil, func(data interface{}) error {
			rw := data.(http.ResponseWriter)
			rw.WriteHeader(http.StatusFound)
			ch <- struct{}{}
			return nil
		})
		<-ch
	})

	http.ListenAndServe(":6060", nil)
}
