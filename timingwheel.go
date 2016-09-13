package gotimer

import (
	"time"
)

const (
	//10 milliseconds is min interval gotimer can accept
	MinTickInterval time.Duration = 10 * time.Millisecond
	BigTickInterval time.Duration = time.Second
)

//timingWheel can hold timer for one year
//so number of wheels should be calculated by min interval
type timingWheel struct {
	wheels [5]*wheel
	numWheels int
	interval time.Duration

	numSlotsOfWheel0 int
	numSlotsOfWheelN int
}

func NewTimingWheel(interval time.Duration) *timingWheel{
	if interval < MinTickInterval {
		return nil
	}
	tw := &timingWheel{}
	if interval > BigTickInterval {
		tw.numSlotsOfWheel0 = 64
		tw.numSlotsOfWheelN = 16
	} else {
		tw.numSlotsOfWheel0 = 128
		tw.numSlotsOfWheelN = 64
	}
}

//circular linked list
type wheel struct {
	//pointer to head of slots
	headSlot *slot
	//pointer to tail of slots, it should be same with headSlot
	tailSlot *slot
	//pointer to current slot
	curSlot *slot
	numSlots int
}

type slot struct {
	//pointer to head of list
	headTimer *tWtimer
	//pointer to next slot in a wheel
	nextSlot *slot
}

//doubly linked list
type tWtimer struct {
	baseTimer
	prev *tWtimer
	next *tWtimer
}
