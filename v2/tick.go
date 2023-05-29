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
	"time"

	"github.com/singchia/go-timer/v2/pkg/linker"
)

type tickOption struct {
	data    interface{}
	ch      chan *Event
	handler func(*Event)
}

// the real shit
type tick struct {
	*tickOption
	id         linker.DoubID
	s          *slot
	ipw        []uint
	duration   time.Duration
	delay      time.Duration
	insertTime time.Time
}

func (t *tick) Reset(data interface{}) {
	t.s.update(t, data)
}

func (t *tick) Cancel() {
	t.s.w.tw.operations <- &operation{
		tick:     t,
		opertype: operdel,
	}
}

func (t *tick) Delay(d time.Duration) {
	t.delay = d
	t.s.w.tw.operations <- &operation{
		tick:     t,
		opertype: operdelay,
	}
}

func (t *tick) C() <-chan *Event {
	return t.ch
}

func (t *tick) InsertTime() time.Time {
	return t.insertTime
}

func (t *tick) Duration() time.Duration {
	return t.duration
}
