package timer

import (
	"github.com/singchia/go-timer/v2/pkg/linker"
)

type wheel struct {
	//keep the reference to for Tick(Tick-->slot-->wheel-->timingwheel-->wheels)
	tw       *timingwheel
	slots    []*slot
	cur      uint
	numSlots uint
	//position in whole timer
	position uint
}

func newWheel(t *timingwheel, numSlots uint, position uint) *wheel {
	slots := make([]*slot, numSlots, numSlots)
	wil := &wheel{tw: t, cur: 0, numSlots: numSlots, position: position}
	wil.slots = slots
	return wil
}

func (w *wheel) add(n uint, tick *tick) *tick {
	if w.slots[n] == nil {
		w.slots[n] = newSlot(w)
	}
	return w.slots[n].add(tick)
}

// increace n on cur
func (w *wheel) incN(n uint) *linker.Doublinker {
	w.cur += n
	if w.cur >= w.numSlots {
		w.cur -= w.numSlots
	}
	if w.slots[w.cur] == nil || w.slots[w.cur].length() == 0 {
		return nil
	}
	return w.slots[w.cur].remove()
}
