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
	"sync"
	"time"

	"github.com/singchia/go-timer/pkg/scheduler"
)

const (
	defaultTickInterval time.Duration = time.Millisecond
	defaultTicks        uint64        = 1024 * 1024 * 1024 * 1024
	defaultSlots        uint          = 256
)

var (
	ErrDurationOutOfRange = errors.New("duration out of range")
	ErrTimerNotStarted    = errors.New("timer not started")
	ErrTimerForceClosed   = errors.New("timer force closed")
)

const (
	twStatusInit = iota
	twStatusStarted
	twStatusPaused
	twStatusStoped
)

type opertype int

const (
	operdel opertype = iota
	operadd
	operdelay
)

type operation struct {
	tick     *tick
	opertype opertype
}

type timerOption struct {
	interval time.Duration
}

type timingwheel struct {
	*timerOption

	mtx      sync.RWMutex
	twStatus int

	wheels     []*wheel
	max        uint64
	sch        *scheduler.Scheduler
	operations chan *operation

	pause, quit chan struct{}
}

func newTimingwheel(opts ...TimerOption) *timingwheel {
	max, length := calcuWheels(defaultTicks)
	tw := &timingwheel{
		twStatus:   twStatusInit,
		sch:        scheduler.NewScheduler(),
		operations: make(chan *operation, 1024),
		timerOption: &timerOption{
			interval: defaultTickInterval,
		},
		pause: make(chan struct{}),
		quit:  make(chan struct{}),
	}
	for _, opt := range opts {
		opt(tw.timerOption)
	}
	tw.setWheels(max, length)
	tw.sch.StartSchedule()
	return tw
}

func (tw *timingwheel) setWheels(max uint64, length int) {
	tw.wheels = make([]*wheel, 0, length)
	tw.max = max
	for i := uint(0); i < uint(length); i++ {
		tw.wheels = append(tw.wheels, newWheel(tw, defaultSlots, i))
	}
}

func (tw *timingwheel) Add(d time.Duration, opts ...TickOption) Tick {
	tick := &tick{
		duration:   d,
		insertTime: time.Now(),
		tickOption: &tickOption{},
	}
	for _, opt := range opts {
		opt(tick.tickOption)
	}
	if tick.C == nil && tick.handler == nil {
		tick.C = make(chan interface{}, 1)
	}
	tw.operations <- &operation{tick, operadd}
	return tick
}

func (t *timingwheel) Start() {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if t.twStatus != twStatusInit && t.twStatus != twStatusPaused {
		return
	}
	t.twStatus = twStatusStarted
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

func (t *timingwheel) Close() {
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
		case operation := <-t.operations:
			switch operation.opertype {
			case operadd:
				ipw := t.indexesPerWheel(operation.tick.duration)
				operation.tick.ipw = ipw
				t.wheels[len(ipw)-1].add(ipw[len(ipw)-1], operation.tick)

			case operdel:
				operation.tick.s.delete(operation.tick)

			case operdelay:
				operation.tick.s.delete(operation.tick)
				ipw := t.indexesPerWheelBased(operation.tick.delay, operation.tick.ipw)
				operation.tick.ipw = ipw
				operation.tick.duration += operation.tick.delay
				t.wheels[len(ipw)-1].add(ipw[len(ipw)-1], operation.tick)
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
		}
	}
}

func (t *timingwheel) iterate(data interface{}) error {
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

// TODO
func (t *timingwheel) forceClose(data interface{}) error {
	tick, _ := data.(*tick)
	t.sch.PublishRequest(&scheduler.Request{Data: tick, Handler: t.handle})
	return nil
}

func (t *timingwheel) handle(data interface{}) {
	tick, _ := data.(*tick)
	if tick.C == nil {
		tick.handler(tick.data, nil)
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
		if count < num {
			count *= uint64(defaultSlots)
			length += 1
			continue
		}
		break
	}
	return count, length
}