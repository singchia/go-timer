package timer

import (
	"github.com/singchia/go-timer/v2/pkg/linker"
)

type slot struct {
	dlinker *linker.Doublinker
	w       *wheel
}

func newSlot(w *wheel) *slot {
	return &slot{w: w, dlinker: linker.NewDoublinker()}
}

func (s *slot) length() int64 {
	return s.dlinker.Length()
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

func (s *slot) close() *linker.Doublinker {
	if s == nil {
		return nil
	}
	s.w = nil
	temp := s.dlinker
	s.dlinker = nil
	return temp
}

func (s *slot) foreach(handler linker.ForeachFunc) error {
	if s == nil {
		return nil
	}
	return s.dlinker.Foreach(handler)
}
