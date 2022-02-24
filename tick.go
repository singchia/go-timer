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
	"github.com/singchia/go-timer/pkg/linker"
)

type tickOption struct {
	data    interface{}
	C       chan interface{}
	handler Handler
}

//the real shit
type tick struct {
	data     interface{}
	C        chan interface{}
	handler  Handler
	id       linker.DoubID
	s        *slot
	ipw      []uint
	duration uint64
}

func (t *tick) Reset(data interface{}) error {
	return t.s.update(t, data)
}

func (t *tick) Cancel() error {
	return t.s.delete(t)
}

func (t *tick) Delay(d uint64) error {
	t.s.delete(t)
	_, err := t.s.w.tw.timeBased(d, t)
	return err
}

func (t *tick) Tunnel() <-chan interface{} {
	return t.C
}
