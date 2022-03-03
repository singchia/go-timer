/*
 * Copyright (c) 2021 Austin Zhai <singchia@163.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation; either version 2 of
 * the License, or (at your option) any later version.
 */
package timer

import "time"

type Handler func(data interface{}) error

type TimerOption func(*timerOption)

func WithTimeInterval(interval time.Duration) TimerOption {
	return func(to *timerOption) {
		to.interval = interval
	}
}

func WithHandlerPool(size int) TimerOption {
	return func(to *timerOption) {
	}
}

type TickOption func(*tickOption)

func WithData(data interface{}) TickOption {
	return func(to *tickOption) {
		to.data = data
	}
}

func WithChan(C chan interface{}) TickOption {
	return func(to *tickOption) {
		to.C = C
	}
}

func WithHandler(handler Handler) TickOption {
	return func(to *tickOption) {
		to.handler = handler
	}
}

type Timer interface {
	//Time preset a Tick which will be triggered after d ticks,
	//you can set channel C, and after d ticks, data would be consumed from C.
	//Or you can set func handler, after d ticks, data would be handled
	//by handler in go-timer. If neither one be set, go-timer will generate a channel,
	//it's attatched with return value Tick, get it by Tick.Channel().
	//Time must be called after Timer.Start.
	Time(d time.Duration, opts ...TickOption) Tick

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
}

//Tick that set in Timer can be required from Timer.Time()
type Tick interface {
	//To reset the data set at Timer.Time()
	Reset(data interface{})

	//To cancel the tick
	Cancel()

	//Delay the tick if not timeout
	Delay(d time.Duration)

	//To get the channel called at Timer.Time(),
	//you will get the same channel if set, if not and handler is nil,
	//then a new created channel will be returned.
	Channel() <-chan interface{}

	// Insert time
	InsertTime() time.Time

	// Duration
	Duration() time.Duration
}

//Entry
func NewTimer(opts ...TimerOption) Timer {
	return newTimingwheel(opts...)
}
