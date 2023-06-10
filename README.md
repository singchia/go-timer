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
* High performance in timing-wheels algorithm
* Minimal resources use
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
	t := timer.NewTimer()
	tick := t.Add(time.Second, timer.WithCyclically())
	for {
		<-tick.C()
		log.Println(time.Now())
	}
}
```

## License

© Austin Zhai, 2015-2025

Released under the [Apache License 2.0](https://github.com/singchia/go-timer/blob/master/LICENSE)