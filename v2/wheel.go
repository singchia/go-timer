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
	slots := make([]*slot, 0, numSlots)
	wil := &wheel{tw: t, cur: 0, numSlots: numSlots, position: position}
	for i := 0; i < int(numSlots); i++ {
		slot := newSlot(wil)
		slots = append(slots, slot)
	}
	wil.slots = slots
	return wil
}

func (w *wheel) add(n uint, tick *tick) *tick {
	return w.slots[n].add(tick)
}

// increace n on cur
func (w *wheel) incN(n uint) *linker.Doublinker {
	w.cur += n
	if w.cur >= w.numSlots {
		w.cur -= w.numSlots
	}
	return w.slots[w.cur].remove()
}
