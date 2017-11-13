package timer

import (
	"sync"

	"github.com/singchia/go-hammer/doublinker"
)

type wheel struct {
	//keep the reference to for Tick(Tick-->slot-->wheel-->timingwheel-->wheels)
	tw       *timingwheel
	slots    []*slot
	cur      uint
	numSlots uint
	//position in whole timer
	position uint
	//mutex for cur
	wheelMutex sync.RWMutex
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
	w.wheelMutex.RLock()
	defer w.wheelMutex.RUnlock()
	return w.slots[n].add(tick)
}

//increace n on cur
func (w *wheel) incN(n uint) *doublinker.Doublinker {
	w.wheelMutex.Lock()
	defer w.wheelMutex.Unlock()

	w.cur += n
	if w.cur >= w.numSlots {
		w.cur -= w.numSlots
	}
	return w.slots[w.cur].remove()
}
