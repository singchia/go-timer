package timer

import (
	"sync"

	"github.com/singchia/go-hammer/doublinker"
)

type wheel struct {
	slots    []*slot
	cur      uint
	numSlots uint
	//position in whole timer
	position int
	//mutex for cur
	wheelMutex sync.RWMutex
}

func newWheel(numSlots uint, position uint) *wheel {
	slots := make([]*slot, 0, numSlots)
	for i := 0; i < int(numSlots); i++ {
		slot := newSlot()
		slots = append(slots, slot)
	}
	return &wheel{slots: slots, cur: 0, numSlots: numSlots, position: position}
}

func (w *wheel) add(n uint, tick *Tick) *Tick {
	w.wheelMutex.RLock()
	defer w.wheelMutex.RUnlock()
	index := n + w.cur
	if index >= w.numSlots {
		index = index - w.numSlots
	}
	return w.slots[index].add(tick)
}

//increace n on cur
func (w *wheel) incN(n uint) *doublinker.Doublinker {
	w.wheelMutex.Lock()
	defer w.wheelMutex.Unlock()

	index := n + w.cur
	if index >= w.numSlots {
		w.cur = index - w.numSlots
	}
	return w.slots[w.cur].remove()
}
