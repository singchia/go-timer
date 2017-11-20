# go-timer
A implementation of hierarchical timing wheels, since build-in timer in golang has serveral limitations:

* one timer only notify once or one specific duration.
* build-in timer can't keep any states.
* build-in timer can't customize channel, **_n_** timer will create **_n_** channel.

In many senarios, using buind-in timer with goroutines seems very common, but goroutines grow really fast as more timer be started which may reach millions of orders of magnitude. that's the reason **go-timer** be needed. And **go-timer** supports:

* no extra goroutines be started.
* channel customizing.
* timer data depositing.
* timer data modification at runtime.
* delegate function customizing.
* can be paused at runtime.

## How-to-use
**go-timer** supplies several easy-understand and easy-integrate interfaces, Let's see a easy sample.

```golang
package main

import (
    "fmt"
    "time"

    timer "github.com/singchia/go-timer"
)

func main() {
    t := timer.NewTimer()
    t.Start()
    ch := make(chan interface{})
    old := time.Now()
    t.Time(5, 1, ch, nil)
    <-ch
    elapse := time.Now().Sub(old)
    fmt.Printf("time diff: %d", elapse)
}
```

## Installation

If you don't have the Go development environment installed, visit the [Getting Started](https://golang.org/doc/install) document and follow the instructions. Once you're ready, execute the following command:
```
go get -u github.com/singchia/go-timer
```

## API

## Bench

## Details