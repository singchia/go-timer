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
	"sync/atomic"
	"time"
)

const (
	defaultTickInterval time.Duration = time.Millisecond
	defaultPoolSize     int           = 10
	defaultTicks        uint64        = 1024 * 1024 * 1024 * 1024
	defaultSlots        uint          = 2
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
	poolSize int
}

type timingwheel struct {
	*timerOption
	wheels     []*wheel
	max        uint64
	quit       chan struct{}
	operations chan *operation
	doneTickCh chan interface{}
	ticking    int32 // atomic bool
}

func newTimingwheel(opts ...TimerOption) *timingwheel {
	max, length := calcuWheels(defaultTicks)

	tw := &timingwheel{
		wheels:     make([]*wheel, 0, length),
		max:        max,
		quit:       make(chan struct{}),
		operations: make(chan *operation, 1024),
		timerOption: &timerOption{
			interval: defaultTickInterval,
			poolSize: defaultPoolSize,
		},
		doneTickCh: make(chan interface{}, 1024),
		ticking:    0,
	}
	for _, opt := range opts {
		opt(tw.timerOption)
	}
	for i := uint(0); i < uint(length); i++ {
		tw.wheels = append(tw.wheels, newWheel(tw, defaultSlots, i))
	}
	return tw
}

func (t *timingwheel) Time(d time.Duration, opts ...TickOption) Tick {
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
	t.operations <- &operation{tick, operadd}
	return tick
}

func (t *timingwheel) Start() {
	atomic.StoreInt32(&t.ticking, 1)
	go t.drive()
	go t.handle()
}

func (t *timingwheel) Pause() {
	paused := atomic.CompareAndSwapInt32(&t.ticking, 1, 0)
	if paused {
		close(t.quit)
	}
}

func (t *timingwheel) Moveon() {
	moveon := atomic.CompareAndSwapInt32(&t.ticking, 0, 1)
	if moveon {
		t.quit = make(chan struct{})
		go t.drive()
	}
}

func (t *timingwheel) Stop() {
	close(t.quit)
	t.wheels = nil
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

		case <-t.quit:
			return
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
	t.doneTickCh <- tick
	return nil
}

func (t *timingwheel) handle() {
	for i := 0; i < t.poolSize; i++ {
		go func() {
			for {
				select {
				case data := <-t.doneTickCh:
					tick, _ := data.(*tick)
					if tick.C == nil {
						tick.handler(tick.data)
					} else {
						tick.C <- tick.data
					}
				case <-t.quit:
					return
				}
			}
		}()
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
