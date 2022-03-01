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
	"sync"

	"github.com/singchia/go-timer/pkg/linker"
)

type slot struct {
	dlinker   *linker.Doublinker
	w         *wheel
	slotMutex sync.RWMutex
}

func newSlot(w *wheel) *slot {
	return &slot{w: w, dlinker: linker.NewDoublinker()}
}

func (s *slot) add(tick *tick) *tick {
	doubID := s.dlinker.Add(tick)
	tick.id = doubID
	tick.s = s
	return tick
}

func (s *slot) delete(tick *tick) {
	s.dlinker.Delete(tick.id)
}

func (s *slot) update(tick *tick, data interface{}) {
	tick.data = data
}

func (s *slot) remove() *linker.Doublinker {
	temp := s.dlinker
	s.dlinker = linker.NewDoublinker()
	return temp
}

func (s *slot) foreach(handler linker.ForeachFunc) error {
	return s.dlinker.Foreach(handler)
}
