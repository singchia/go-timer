package timer

import (
	"sync"

	"github.com/singchia/go-hammer/linker"
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
	s.slotMutex.RLock()
	defer s.slotMutex.RUnlock()

	doubID := s.dlinker.Add(tick)
	tick.id = doubID
	tick.s = s
	return tick
}

func (s *slot) delete(tick *tick) error {
	s.slotMutex.RLock()
	defer s.slotMutex.RUnlock()
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
	return s.dlinker.Foreach(handler)
}

type Handler func(data interface{}) error

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
