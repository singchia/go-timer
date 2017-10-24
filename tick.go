package timer

import (
	"github.com/singchia/go-hammer/doublinker"
)

type slot struct {
	dlinker doublinker.Doublinker
}

func newSlot() *slot {
	return &slot{}
}

func (s *slot) add(data interface{}, Done chan interface{}, handler Handler) *Tick {
	tick := &Tick{data: data, Done: Done, handler: handler}
	doubID := s.dlinker.Add(tick)
	tick.id = doubID
	return tick
}

func (s *slot) delete(id TickID) error {
	doubID := (*Tick)(id).id
	return s.dlinker.Delete(doubID)
}

func (s *slot) update(id TickID, data interface{}) error {
	tick := (*Tick)(id)
	tick.data = data
}

type Handler func(interface{}) error

//the real shit
type Tick struct {
	data    interface{}
	Done    chan interface{}
	handler Hanlder
	id      doublinker.DoubID
	s       *slot
}

func (t *Tick) Reset(data interface{}) error {
	return t.s.Update(t, data)
}

func (t *Tick) Cancel() error {
	return t.s.Delete(t)
}

func (t *Tick) Delay() error {}
