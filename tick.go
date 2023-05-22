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

	"github.com/singchia/go-timer/pkg/linker"
)

type tickOption struct {
	data    interface{}
	C       chan interface{}
	handler Handler
}

func newSlot(w *wheel) *slot {
	return &slot{w: w, dlinker: linker.NewDoublinker()}
}

func (s *slot) add(tick *tick) *tick {
	s.slotMutex.Lock()
	defer s.slotMutex.Unlock()

	doubID := s.dlinker.Add(tick)
	tick.id = doubID
	tick.s = s
	return tick
}

func (s *slot) delete(tick *tick) error {
	s.slotMutex.Lock()
	defer s.slotMutex.Unlock()
	return s.dlinker.Delete(tick.id)
}

func (s *slot) update(tick *tick, data interface{}) error {
	tick.data = data
	return nil
}

func (s *slot) remove() *linker.Doublinker {
	s.slotMutex.Lock()
	defer s.slotMutex.Unlock()
	temp := s.dlinker
	s.dlinker = linker.NewDoublinker()
	return temp
}

func (s *slot) foreach(handler linker.ForeachFunc) error {
	s.slotMutex.RLock()
	defer s.slotMutex.RUnlock()
	return s.dlinker.Foreach(handler)
}

type Handler func(data interface{}) error

// the real shit
type tick struct {
	data     interface{}
	C        chan interface{}
	handler  Handler
	id       linker.DoubID
	s        *slot
	ipw      []uint
	duration uint64

}

func (t *tick) Delay(d time.Duration) {
	t.delay = d
	t.s.w.tw.operations <- &operation{
		tick:     t,
		opertype: operdelay,
	}
	return
}

func (t *tick) Channel() <-chan interface{} {
	return t.C
}

func (t *tick) InsertTime() time.Time {
	return t.insertTime
}

func (t *tick) Duration() time.Duration {
	return t.duration
}
