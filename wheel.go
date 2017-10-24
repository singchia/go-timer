package timer

type wheel struct {
	slots    []*slot
	cur      uint
	numSlots uint
}

func newWheel(numSlots uint) *wheel {
	slots := make([]*slot, 0, numSlots)
	for i := 0; i < numSlots; i++ {
		slot := newSlot()
		slots = append(slots, slot)
	}
	return &wheel{slots: slots, cur: 0, numSlots: numSlots}
}

//forword nth slot to return
func (w *wheel) forwordN(n uint) *slot {
	index := n + w.cur
	if index >= w.numSlots {
		return w.slots[index-w.numSlots]
	}
	return w.slots[n+w.cur]
}

//increace n on cur
func (w *wheel) incN(n uint) {
	index := n + w.cur
	if index >= w.numSlots {
		w.cur = index - w.numSlots
	}
}

func (w *wheel) add() {}
