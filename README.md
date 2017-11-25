# go-timer
An implementation of hierarchical timing wheels, since build-in timer in golang has serveral limitations:

* one timer only notify once or one specific duration.
* build-in timer can't keep any states.
* build-in timer can't customize channel, **_n_** timer will create **_n_** channel.

In many senarios, using buind-in timer with goroutines seems very common, but goroutines grow really fast as more timer be started which may reach millions of orders of magnitude. that's the reason **go-timer** be needed. And **go-timer** supports:

* no extra goroutines will be started.
* channel customizing.
* timer data depositing.
* timer data modification at runtime.
* delegate function customizing.
* can be paused at runtime.

## How-to-use
**go-timer** supplies several easy-understand and easy-integrate interfaces, Let's see an easy sample.

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
### _func NewTimer() Timer_

### _Timer_

```golang
//SetMaxTicks sets the max number of ticks,
//it should be called before Start if you want to customize it.
//Default 1024*1024*1024 ticks
SetMaxTicks(max uint64)

//SetInterval set the time lapse between two ticks.
//it should be called before Start if you want to customize it.
//Default 1 second.
//And the whole time lapse should be maxTicks*interval.
//Note that the max tick you can preset is (maxTicks*interval-1),
//because current tick cannot be set.
SetInterval(interval time.Duration)

//Time preset a Tick which will be triggered after d ticks,
//you can set channel C, and after d ticks, data would be consumed from C.
//Or you can set func handler, after d ticks, data would be handled
//by handler in go-timer. If neither one be set, go-timer will generate a channel,
//it's attatched with return value Tick, get it by Tick.Tunnel().
//Time must be called after Timer.Start.
Time(d uint64, data interface{}, C chan interface{}, handler Handler) (Tick, error)

//Start to start timer.
Start()

//Stop to stop timer,
//all ticks set would be discarded.
Stop()

//Pause the timer,
//all ticks won't continue after Timer.Movenon().
Pause()

//Continue the paused timer.
Moveon()
```

### _Tick_

```golang
//To reset the data set at Timer.Time()
Reset(data interface{}) error

//To cancel the tick
Cancel() error

//Delay the tick if not timeout
Delay(d uint64) error

//To get the channel called at Timer.Time(),
//you will get the same channel if set, if not and handler is nil,
//then a new created channel will be returned.
Tunnel() <-chan interface{}
```


## Benchmarks
**my bench env:**   

* MacBook Pro (13-inch, Late 2016, Four Thunderbolt 3 Ports)
* 2.9 GHz Intel Core i5
* 16 GB 2133 MHz LPDDR3

**bench code:**   
see [here](bench/main.go).   

**conclusion:**   
if (amount of goroutines) > 30w, **go-timer** got less errors, else **build-in timer** is better.   
if you don't want too much goroutines and  less precise can be acceptable, I recommend **go-timer**.

