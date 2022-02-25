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
	"fmt"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	scheduler "github.com/singchia/go-scheduler"
)

const (
	defaultTickInterval time.Duration = time.Millisecond
	defaultTicks        uint64        = 1024 * 1024 * 1024 * 1024
	defaultSlots        uint64        = 256
)

type timerOption struct {
	interval time.Duration
}

type timingwheel struct {
	*timerOption
	wheels []*wheel
	max    uint64
	sch    *scheduler.Scheduler
	quit   chan struct{}
}

func newTimingwheel(opts ...TimerOption) *timingwheel {
	max, length := calcuWheels(defaultTicks)

	tw := &timingwheel{
		wheels: make([]*wheel, length),
		max:    max,
		sch:    scheduler.NewScheduler(),
		quit:   make(chan struct{}),
		timerOption: &timerOption{
			interval: defaultTickInterval,
		},
	}
	for _, opt := range opts {
		opt(tw.timerOption)
	}
	for i := 0; i < length; i++ {
		tw.wheels = append(tw.wheels, newWheel(tw, defaultSlots, i))
	}
	return tw
}

func (t *timingwheel) Time(d time.Duration, opts ...Option) Tick {
	tick := &tick{}
	for _, opt := range opts {
		opt(tick)
	}
	if tick.C == nil && tick.handler == nil {
		tick.C = make(chan interface{}, 1)
	}
	ipw := t.indexesPerWheel(d)

	tick := &tick{
		data:     data,
		C:        C,
		handler:  handler,
		ipw:      ipw,
		duration: d}
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
	if t.wheels == nil {
		t.SetMaxTicks(DefaultMaxTicks)
	}
	go t.drive()
	return
}

func (t *timingwheel) Pause() {
	t.quit <- struct{}{}
	return
}

func (t *timingwheel) Moveon() {
	go t.drive()
}

func (t *timingwheel) Stop() {
	t.signal <- struct{}{}
	t.wheels = nil
	t.sch.Close()
}

func (t *timingwheel) drive() {
	driver := time.NewTicker(t.interval)
	for {
		select {
		case <-driver.C:
			for _, wheel := range t.wheels {
				linker := wheel.incN(1)
				linker.Foreach(t.iterate)
				if wheel.cur != 0 {
					break
				}
			}
		case <-t.signal:
			return
		}
	}
	return
}

func (t *timingwheel) iterate(data interface{}) error {
	t.sch.PublishRequest(&scheduler.Request{Data: data, Handler: t.handle})
	return nil
}

func (t *timingwheel) handle(data interface{}) {
	tick, _ := data.(*tick)
	position := tick.s.w.position
	for position > 0 {
		position--
		if tick.ipw[position] > 0 {
			t.wheels[position].add(tick.ipw[position], tick)
			return
		}
	}
	if tick.C == nil {
		tick.handler(tick.data)
	} else {
		tick.C <- tick.data
	}
}

func (t *timingwheel) indexesPerWheel(d uint64) []uint {
	var ipw []uint
	var reminder uint64
	var quotient = d
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

func (t *timingwheel) indexesPerWheelBased(d uint64, base []uint) []uint {
	var ipw []uint
	var reminder uint64
	var quotient = d
	for i, wheel := range t.wheels {
		if quotient == 0 {
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
	count := defaultSlotsPerWheel
	length := 1
	for {
		if count < num {
			count *= defaultSlotsPerWheel
			length += 1
			continue
		}
		break
	}
	return count, length
}

//for debug
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
		w := simplejson.New()
		w.Set(fmt.Sprintf("wheel%d", i), ss)
		ws = append(ws, w)
	}
	j.Set("wheels", ws)
	return j.MarshalJSON()
}
