# GO-TIMER

[![Go](https://github.com/singchia/go-xtables/actions/workflows/go.yml/badge.svg)](https://github.com/singchia/go-xtables/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/singchia/go-timer/v2)](https://goreportcard.com/report/github.com/singchia/go-timer/v2)

A high performance timer with minimal goroutines.

## Getting Started

```golang
package main

import (
	"time"

	timer "github.com/singchia/go-timer/v2"
)

func main() {
	t := timer.NewTimer(timer.WithTimeInterval(time.Millisecond))
	ch := make(chan *timer.Event)
	old := time.Now()

	tick := t.Add(5, timer.WithData(1), timer.WithChan(ch))
	<-tick.C()
}
```