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
	"sync/atomic"
	"time"

	"github.com/singchia/go-timer/pkg/scheduler"
)

const (
	defaultTickInterval time.Duration = 100 * time.Millisecond
	defaultTicks        uint64        = 1024 * 1024 * 1024 * 1024
	defaultSlots        uint          = 256
)

var (
	ErrDurationOutOfRange   = errors.New("duration out of range")
	ErrTimerNotStarted      = errors.New("timer not started")
	ErrTimerForceClosed     = errors.New("timer force closed")
	ErrOperationForceClosed = errors.New("operation force closed")
	ErrDelayOnCyclically    = errors.New("cannot delay on cyclically tick")
	ErrCancelOnNonWait      = errors.New("cannot cancel on firing or fired tick")
)

const (
	twStatusInit = iota
	twStatusStarted
	twStatusPaused
	twStatusStoped
)

// operation type and return
type operType int

const (
	operAdd operType = iota
	operCancel
	operDelay
	operReset
)

type operation struct {
	tick     *tick
	operType operType
	retCh    chan *operationRet
	// for delay
	delay time.Duration
	// for reset
	data interface{}
}

type operRetType int

const (
	operOK operRetType = iota
	operErr
)

type operationRet struct {
	operRetType operRetType
	err         error
}

// timer options
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
		tw:         tw,
		insertTime: time.Now(),
		tickOption: &tickOption{},
		status:     statusAdd,
	}
	for _, opt := range opts {
		opt(tick.tickOption)
	}
	if tick.ch == nil && tick.handler == nil {
		tick.ch = make(chan *Event, 1)
	}
	tw.operations <- &operation{
		tick:     tick,
		operType: operAdd,
	}
	return tick
}

func (tw *timingwheel) Start() {
	tw.mtx.Lock()
	defer tw.mtx.Unlock()
	if tw.twStatus != twStatusInit && tw.twStatus != twStatusPaused {
		return
	}
	tw.twStatus = twStatusStarted
	go tw.drive()
	return
}

func (tw *timingwheel) Pause() {
	tw.mtx.Lock()
	defer tw.mtx.Unlock()
	if tw.twStatus != twStatusStarted {
		return
	}
	tw.twStatus = twStatusPaused
	tw.pause <- struct{}{}
	return
}

func (tw *timingwheel) Moveon() {
	tw.mtx.Lock()
	defer tw.mtx.Unlock()
	if tw.twStatus != twStatusPaused {
		return
	}
	tw.twStatus = twStatusStarted
	go tw.drive()
}

func (tw *timingwheel) Close() {
	tw.mtx.Lock()
	defer tw.mtx.Unlock()
	if tw.twStatus != twStatusStarted {
		return
	}
	tw.quit <- struct{}{}
	tw.twStatus = twStatusStoped
}

func (tw *timingwheel) drive() {
	driver := time.NewTicker(tw.interval)
	for {
		select {
		case <-driver.C:
			// fire all ready tick
			for _, wheel := range tw.wheels {
				linker := wheel.incN(1)
				if linker.Length() > 0 {
					linker.Foreach(tw.iterate)
				}
				if wheel.cur != 0 {
					break
				}
			}
		case operation := <-tw.operations:
			switch operation.operType {
			case operAdd:
				// the specific tick's add operation must be before the del or delay operation
				// so we don't care the status cause it must be statusAdd
				ipw := tw.indexesPerWheel(operation.tick.duration)
				operation.tick.ipw = ipw
				tw.wheels[len(ipw)-1].add(ipw[len(ipw)-1], operation.tick)
				operation.tick.status = statusWait

			case operCancel:
				// only in wait status can delete the tick
				if operation.tick.status != statusWait {
					operation.retCh <- &operationRet{
						operRetType: operErr,
						err:         ErrCancelOnNonWait,
					}
					continue
				}
				operation.tick.s.delete(operation.tick)
				operation.retCh <- &operationRet{
					operRetType: operOK,
				}

			case operDelay:
				if operation.tick.status != statusWait {
					operation.retCh <- &operationRet{
						operRetType: operErr,
						err:         ErrCancelOnNonWait,
					}
					continue
				}
				operation.tick.s.delete(operation.tick)
				ipw := tw.indexesPerWheelBased(operation.delay, operation.tick.ipw)
				operation.tick.ipw = ipw
				operation.tick.duration += operation.delay
				tw.wheels[len(ipw)-1].add(ipw[len(ipw)-1], operation.tick)
				operation.retCh <- &operationRet{
					operRetType: operOK,
				}

			case operReset:
				if operation.tick.status != statusWait {
					operation.retCh <- &operationRet{
						operRetType: operErr,
						err:         ErrCancelOnNonWait,
					}
					continue
				}
				operation.tick.data = operation.data
				operation.retCh <- &operationRet{
					operRetType: operOK,
				}
			}
		case <-tw.pause:
			return
		case <-tw.quit:
			for _, wheel := range tw.wheels {
				for _, slot := range wheel.slots {
					slot.foreach(tw.forceClose)
				}
			}
			tw.sch.Close()
		}
	}
}

func (tw *timingwheel) iterate(data interface{}) error {
	tick, _ := data.(*tick)
	position := tick.s.w.position
	for position > 0 {
		position--
		if tick.ipw[position] > 0 {
			tw.wheels[position].add(tick.ipw[position], tick)
			return nil
		}
	}
	tick.status = statusFire
	tw.sch.PublishRequest(&scheduler.Request{Data: tick, Handler: tw.handleNormal})
	if tick.cyclically {
		ipw := tw.indexesPerWheelBased(tick.duration, tick.ipw)
		tick.ipw = ipw
		tw.wheels[len(ipw)-1].add(ipw[len(ipw)-1], tick)
	}
	tick.status = statusWait
	return nil
}

// TODO
func (tw *timingwheel) forceClose(data interface{}) error {
	tick, _ := data.(*tick)
	tick.status = statusFire
	tw.sch.PublishRequest(&scheduler.Request{Data: tick, Handler: tw.handleError})
	return nil
}

func (tw *timingwheel) handleError(data interface{}) {
	tick, _ := data.(*tick)
	event := &Event{
		Duration:   tick.duration,
		InsertTIme: tick.insertTime,
		Data:       tick.data,
		Error:      ErrTimerForceClosed,
	}
	if tick.ch == nil {
		tick.handler(event)
	} else {
		tick.ch <- event
	}
}

func (tw *timingwheel) handleNormal(data interface{}) {
	tick, _ := data.(*tick)
	atomic.AddInt64(&tick.fired, 1)
	event := &Event{
		Duration:   tick.duration,
		InsertTIme: tick.insertTime,
		Data:       tick.data,
		Error:      nil,
	}
	if tick.ch == nil {
		tick.handler(event)
	} else {
		tick.ch <- event
	}
}

func (tw *timingwheel) indexesPerWheel(d time.Duration) []uint {
	var ipw []uint
	var reminder uint64
	var quotient = uint64((d + tw.interval - 1) / tw.interval)
	for i, wheel := range tw.wheels {
		if quotient == 0 {
			break
		}
		quotient += uint64(tw.wheels[i].cur)
		reminder = quotient % uint64(wheel.numSlots)
		quotient = quotient / uint64(wheel.numSlots)
		ipw = append(ipw, uint(reminder))
	}
	return ipw
}

func (tw *timingwheel) indexesPerWheelBased(d time.Duration, base []uint) []uint {
	var ipw []uint
	var reminder uint64
	var quotient = uint64((d + tw.interval - 1) / tw.interval)
	for i, wheel := range tw.wheels {
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
