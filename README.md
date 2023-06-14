# GO-TIMER

[![Go](https://github.com/singchia/go-xtables/actions/workflows/go.yml/badge.svg)](https://github.com/singchia/go-xtables/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/singchia/go-timer/v2)](https://goreportcard.com/report/github.com/singchia/go-timer/v2)

## Overview

A high performance timer with minimal goroutines.

### How it works

![](docs/overview.jpg)

### Features

* One goroutine runs all
* Goroutine safe
* Clean and simple, no third-party deps at all
* High performance with timing-wheels algorithm
* Minimal resources use
* Managed data and handler
* Customizing channel
* Well tested

## Usage

### Quick Start

```golang
package main

import (
	"log"
	"time"

	timer "github.com/singchia/go-timer/v2"
)

func main() {
	t1 := time.Now()
	// new timer
	t := timer.NewTimer()
	// add a tick in 1s
	tick := t.Add(time.Second)
	// wait for it
	<-tick.C()
	// tick fired as time is up, calcurate and print the elapse
	log.Printf("time elapsed: %fs\n", time.Now().Sub(t1).Seconds())
}
```

### Async handler

```golang
package main

import (
	"log"
	"sync"
	"time"

	timer "github.com/singchia/go-timer/v2"
)

func main() {
	// we need a wait group since using async handler
	wg := sync.WaitGroup{}
	wg.Add(1)
	// new timer
	t := timer.NewTimer()
	// add a tick in 1s with current time and a async handler
	t.Add(time.Second, timer.WithData(time.Now()), timer.WithHandler(func(event *timer.Event) {
		defer wg.Done()
		// tick fired as time is up, calcurate and print the elapse
		log.Printf("time elapsed: %fs\n", time.Now().Sub(event.Data.(time.Time)).Seconds())
	}))

	wg.Wait()
}
```

### With cyclically

```golang
package main

import (
	"log"
	"time"

	timer "github.com/singchia/go-timer/v2"
)

func main() {
	t1 := time.Now()
	// new timer
	t := timer.NewTimer()
	// add cyclical tick in 1s
	tick := t.Add(time.Second, timer.WithCyclically())
	for {
		// wait for it cyclically
		<-tick.C()
		t2 := time.Now()
		// calcurate and print the elapse
		log.Printf("time elapsed: %fs\n", t2.Sub(t1).Seconds())
		t1 = t2
	}
}
```

## License

© Austin Zhai, 2015-2025

Released under the [Apache License 2.0](https://github.com/singchia/go-timer/blob/master/LICENSE)