package timer

import (
	"sync"

	"github.com/singchia/go-hammer/doublinker"
)

type slot struct {
	dlinker   *doublinker.Doublinker
	w         *wheel
	slotMutex sync.RWMutex
}

func newSlot(w *wheel) *slot {
	return &slot{w: w, dlinker: doublinker.NewDoublinker()}
}

func (s *slot) add(tick *Tick) *Tick {
	s.slotMutex.RLock()
	defer s.slotMutex.RUnlock()

	doubID := s.dlinker.Add(tick)
	tick.id = doubID
	tick.s = s
	return tick
}

func (s *slot) delete(tick *Tick) error {
	s.slotMutex.RLock()
	defer s.slotMutex.RUnlock()
	return s.dlinker.Delete(tick.id)
}

func (s *slot) update(tick *Tick, data interface{}) error {
	tick.data = data
	return nil
}

func (s *slot) remove() *doublinker.Doublinker {
	s.slotMutex.Lock()
	defer s.slotMutex.Unlock()
	temp := s.dlinker
	s.dlinker = doublinker.NewDoublinker()
	return temp
}

type Handler func(interface{}) error

//the real shit
type Tick struct {
	data     interface{}
	C        chan interface{}
	handler  Handler
	id       doublinker.DoubID
	s        *slot
	ipw      []uint
	duration uint64
}

func (t *Tick) Reset(data interface{}) error {
	return t.s.update(t, data)
}

func (t *Tick) Cancel() error {
	return t.s.delete(t)
}

func (t *Tick) Delay() error {
	return nil
}
