/*
 * Copyright (c) 2021 Austin Zhai <singchia@163.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation; either version 2 of
 * the License, or (at your option) any later version.
 */
package timer

import (
	"errors"
	"fmt"
	"sync"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/singchia/go-timer/pkg/scheduler"
)

const (
	defaultTickInterval time.Duration = time.Millisecond
	defaultTicks        uint64        = 1024 * 1024 * 1024 * 1024
	defaultSlots        uint          = 256
)

var (
	ErrInvalidDuration  = errors.New("invalid duration")
	ErrTimerUnavailable = errors.New("timer unavailable")
)

const (
	twStatusInit = iota
	twStatusStarted
	twStatusPaused
	twStatusStoped
)

type timingwheel struct {
	mtx      sync.RWMutex
	twStatus int

	wheels   []*wheel
	interval time.Duration
	max      uint64
	sch      *scheduler.Scheduler

	pause, quit chan struct{}
}

func newTimingwheel() *timingwheel {
	tw := &timingwheel{
		interval: DefaultTickInterval,
		twStatus: twStatusInit,
		sch:      scheduler.NewScheduler(),
		pause:    make(chan struct{}),
		quit:     make(chan struct{})}
	tw.sch.SetMaxGoroutines(1000)
	tw.sch.StartSchedule()
	return tw
}

func (t *timingwheel) SetMaxTicks(max uint64) {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	if t.twStatus != twStatusInit {
		return
	}

	t.setMaxTicks(max)
}

func (t *timingwheel) setMaxTicks(max uint64) {
	t.max = max
	nums := calcuQuotients(max)
	t.wheels = make([]*wheel, 0, len(nums))
	for position, num := range nums {
		if position == len(nums)-1 {
			num = num + 1
		}
		t.wheels = append(t.wheels, newWheel(t, num, uint(position)))
	}
}

func (t *timingwheel) SetInterval(interval time.Duration) {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	if t.twStatus != twStatusInit {
		return
	}
	t.interval = interval
}

func (t *timingwheel) Time(d uint64, data interface{}, C chan interface{}, handler Handler) (Tick, error) {
	if d == 0 || d > t.max-1 {
		return nil, ErrInvalidDuration
	}
	for _, opt := range opts {
		opt(tw.timerOption)
	}
	t.mtx.RLock()
	if t.twStatus != twStatusStarted && t.twStatus != twStatusPaused {
		t.mtx.RUnlock()
		return nil, ErrTimerUnavailable
	}
	t.mtx.RUnlock()

	ipw := t.indexesPerWheel(d)
	tick := &tick{data: data, C: C, handler: handler, ipw: ipw, duration: d}
	t.wheels[len(ipw)-1].add(ipw[len(ipw)-1], tick)
	return tick, nil
}

func (t *timingwheel) timeBased(d uint64, tick *tick) (*tick, error) {
	ipw := t.indexesPerWheelBased(d, tick.ipw)
	tick.ipw = ipw

	tick.duration += d
	t.wheels[len(ipw)-1].add(ipw[len(ipw)-1], tick)
	return tick, nil
}

func (t *timingwheel) Start() {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if t.twStatus != twStatusInit && t.twStatus != twStatusPaused {
		return
	}
	t.twStatus = twStatusStarted
	if t.wheels == nil {
		t.setMaxTicks(DefaultMaxTicks)
	}
	go t.drive()
	return
}

func (t *timingwheel) Pause() {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if t.twStatus != twStatusStarted {
		return
	}
	t.twStatus = twStatusPaused
	t.pause <- struct{}{}
	return
}

func (t *timingwheel) Moveon() {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if t.twStatus != twStatusPaused {
		return
	}
	t.twStatus = twStatusStarted
	go t.drive()
}

func (t *timingwheel) Stop() {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if t.twStatus != twStatusStarted {
		return
	}
	t.quit <- struct{}{}
	t.twStatus = twStatusStoped
}

func (t *timingwheel) drive() {
	driver := time.NewTicker(t.interval)
	for {
		select {
		case <-driver.C:
			for _, wheel := range t.wheels {
				linker := wheel.incN(1)
				if linker.Length() > 0 {
					linker.Foreach(t.iterate)
				}
				if wheel.cur != 0 {
					break
				}
			}
		case <-t.pause:
			return
		case <-t.quit:
			for _, wheel := range t.wheels {
				for _, slot := range wheel.slots {
					slot.foreach(t.forceClose)
				}
			}
			t.sch.Close()
			return
		}
	}
}

func (t *timingwheel) iterate(data interface{}) error {
	t.sch.PublishRequest(&scheduler.Request{Data: data, Handler: t.handle})
	return nil
}

// TODO
func (t *timingwheel) forceClose(data interface{}) error {
	t.sch.PublishRequest(&scheduler.Request{Data: data, Handler: func(data interface{}) {
		tick, _ := data.(*tick)
		if tick.C == nil {
			tick.handler(tick.data)
		} else {
			tick.C <- tick.data
		}
	}})
	return nil
}

func (t *timingwheel) handle(data interface{}) {
	tick, _ := data.(*tick)
	position := tick.s.w.position
	for position > 0 {
		position--
		if tick.ipw[position] > 0 {
			t.wheels[position].add(tick.ipw[position], tick)
			return nil
		}
	}
	t.sch.PublishRequest(&scheduler.Request{Data: tick, Handler: t.handle})
	return nil
}

func (t *timingwheel) handle(data interface{}) {
	tick, _ := data.(*tick)
	if tick.C == nil {
		tick.handler(tick.data)
	} else {
		tick.C <- tick.data
	}
}

func (t *timingwheel) indexesPerWheel(d time.Duration) []uint {
	var ipw []uint
	var reminder uint64
	var quotient = uint64((d + t.interval - 1) / t.interval)
	for i, wheel := range t.wheels {
		if quotient == 0 {
			break
		}
		quotient += uint64(t.wheels[i].cur)
		reminder = quotient % uint64(wheel.numSlots)
		quotient = quotient / uint64(wheel.numSlots)
		ipw = append(ipw, uint(reminder))
	}
	return ipw
}

func (t *timingwheel) indexesPerWheelBased(d time.Duration, base []uint) []uint {
	var ipw []uint
	var reminder uint64
	var quotient = uint64((d + t.interval - 1) / t.interval)
	for i, wheel := range t.wheels {
		if quotient == 0 {
			break
		}
		if len(base) <= i {
			ipw = append(ipw, uint(quotient))
			break
		}
		quotient += uint64(base[i])
		reminder = quotient % uint64(wheel.numSlots)
		quotient = quotient / uint64(wheel.numSlots)
		ipw = append(ipw, uint(reminder))
	}
	return ipw
}

func calcuWheels(num uint64) (uint64, int) {
	count := uint64(defaultSlots)
	length := 1
	for {
		quo := num / 16
		rem := num % 16
		if quo != 0 {
			quos = append(quos, 16)
		} else {
			quos = append(quos, uint(rem))
			break
		}
		if rem != 0 {
			quo = quo + 1
		}
		num = quo
	}
	return quos
}

// for debug
func (t *timingwheel) Topology() ([]byte, error) {
	j := simplejson.New()
	var ws []*simplejson.Json
	for i, wheel := range t.wheels {
		var ss []*simplejson.Json
		for j, slot := range wheel.slots {
			var ts []interface{}
			slot.foreach(func(data interface{}) error {
				t := data.(*tick)
				ts = append(ts, t.data)
				return nil
			})
			s := simplejson.New()
			s.Set(fmt.Sprintf("slot%d", j), ts)
			ss = append(ss, s)
		}
		break
	}
	return count, length
}
