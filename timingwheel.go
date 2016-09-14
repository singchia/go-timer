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
type timingWheel struct {
	wheels [5]*wheel
	interval time.Duration

	powerSlotsOfWheel0 int
	powerSlotsOfWheelN int
}

func newTimingWheel(interval time.Duration) *timingWheel{
	if interval < MinTickInterval {
		return nil
	}
	tw := &timingWheel{}
	if interval > BigTickInterval {
		tw.powerSlotsOfWheel0 = 6
		tw.powerSlotsOfWheelN = 4
	} else {
		tw.powerSlotsOfWheel0 = 7
		tw.powerSlotsOfWheelN = 6
	}

	for i:=0; i<len(tw.wheels); i++ {
		if i == 0 {
			tw.wheels[0] = newWheel(1 << tw.powerSlotsOfWheel0, tw)
			continue
		}
		tw.wheels[i] = newWheel(1 << tw.powerSlotsOfWheelN, tw)
	}
	tw.interval = interval
	return tw
}

func (tw *timingWheel) StartTimer(d time.Duration) *Timer {
	var ticks int
	if d < MinTickInterval {
		ticks = 1
	} else {
		ticks = d/tw.interval
	}

	//ticksPerWheel keep ticks for each wheel
	var ticksPerWheel []int
	ticks = ticks + tw.wheels[0].curIndex
	reminder := ticks & (tw.powerSlotsOfWheel0 -1 )
	quotient := ticks >> tw.powerSlotsOfWheel0
	ticksPerWheel = append(ticksPerWheel, reminder)

	for i:=1; i<5; i++ {
		if quotient == 0 {
			break
		}
		quotient += tw.wheels[i].curIndex
		reminder = quotient & (tw.powerSlotsOfWheelN -1)
		quotient = qutient >> tw.powerSlotsOfWheelN
		ticksPerWheel = append(ticksPerWheel, reminder)
	}

	//it decides which wheel should be added
	ticksPerWheelLength := len(ticksPerWheel)
}

func (tw *timingWheel) StartTimerWithTicks(ticks int) *Timer{

}

//circular linked list
type wheel struct {
	tW *timingWheel
	//pointer to head of slots
	headSlot *slot
	//pointer to tail of slots, it should be same with headSlot
	//tailSlot *slot
	//pointer to current slot
	curSlot *slot
	curIndex int

	numSlots int

}

func newWheel(numSlotsOfWheel int, tW *timingWheel) *wheel{
	headSlot := &slot{headTimer:nil, tailTimer:nil, nextSlot:nil}
	curSlot := headSlot

	for i:=0; i<numSlotsOfWheel-1; i++ {
		tempSlot := &slot{headTimer:nil, tailTimer:nil, nextSlot:nil}
		curSlot.next = tempSlot
		curSlot = curSlot.next
	}
	curSlot.next = headSlot
	//redirect curSlot to headSlot
	curSlot = headSlot

	return &wheel{headSlot:headSlot, curSlot:curSlot, curIndex: 0, numSlots:numSlotsOfWheel, tW: tW}
}

type slot struct {
	//pointer to head of list
	headTimer *tWtimer
	//point to tail of list
	tailTimer *tWtimer
	//pointer to next slot in a wheel
	nextSlot *slot
}

func (s *slot) addTimer(d time.Duration, wheel *wheel, ticksPerWheel []int) *Timer {

	timer := &tWtimer{wheel: wheel, slot: slot, ticksPerWheel: ticksPerWheel}
	timer.d = d
	timer.timerCenter = wheel.timingWheel

	if s.headTimer == nil {
		s.headTimer = timer
		s.tailTimer = timer
		return timer
	}

	//add timer to last
	s.tailTimer.next = timer
	timer.prev = s.tailTimer
	s.tailTimer = timer
	return
}

func (s *slot) addTimerWithHandler(d time.Duration, wheel *wheel, ticksPerWheel []int, data interface{}, handler TimerHandler) *Timer {
	timer := s.addTimer(d, wheel, ticksPerWheel)
	timer.ResetData(data)
	timer.AddHandler(handler)
	timer
}

func (s *slot) addTimerWithHandlers(d time.Duration, wheel *wheel, ticksPerWheel []int, data interface{}, handlers []TimerHandler) *Timer {

	timer := s.addTimer(d, wheel, ticksPerWheel)
	timer.ResetData(data)
	timer.AddHandlers(handlers)
	return timer
}

func (s *slot) delTimer() {

}

//doubly linked list
type tWtimer struct {
	//point to slot
	slot *slot
	//point to wheel
	wheel *wheel
	//point to timingWheel
	//tW *timingWheel

	baseTimer
	prev *tWtimer
	next *tWtimer
	ticksPerWheel []int
}
