package gotimer

import (
	"time"
)

const (
	//10 milliseconds is min interval gotimer can accept
	MinTickInterval time.Duration = 10 * time.Millisecond
	BigTickInterval time.Duration = time.Second
)

const (
	TimerDeleteSucceed = iota
	TimerDeleteFailed
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
	var ticksPerWheel []uint
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
	length := len(ticksPerWheel)
	return tw.wheel[length-1].backN(tickPerWheel[length-1]).addTimer(d, tw.wheel[length-1], tickPerWheel)
}

func (tw *timingWheel) StartTimerWithTicks(ticks int) *Timer{

}

//circular linked list
type linkSlotswheel struct {
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

func newlinkSlotsWheel(numSlotsOfWheel int, tW *timingWheel) *linkSlotswheel{
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

	return &linkSlotswheel{headSlot:headSlot, curSlot:curSlot, curIndex: 0, numSlots:numSlotsOfWheel, tW: tW}
}

type wheel struct{
	tW *timingWheel
	slots []*slot
	curIndex int
	numSlots int
}

func newWheel(numSlotsOfWheel int, tW *timingWheel) *wheel {
	slots := make([]*slot, numSlotsOfWheel)
	return &wheel{tW: tW, slots: slots, curIndex: 0, numSlots: numSlotsOfWheel}
}

//back nth slot to return
func (w *wheel) backN(n uint) *slot {
	index := n + w.curIndex
	if index > w.numSlots {
		return w.slots[index - w.numSlots]
	}
	return w.slots[n+w.curIndex]
}

//increace n to curIndex
func (w *wheel) incN(n uint) {
	index := n + w.curIndex
	if index > w.numSlots {
		w.curIndex = index - w.numSlots
	}
}

type slot struct {
	//pointer to head of list
	headTimer *tWtimer
	//point to tail of list
	tailTimer *tWtimer
	//pointer to next slot in a wheel
	nextSlot *slot
	//numTimers
	numTimers int
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
	s.numTimers++
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

func (s *slot) delTimer(timer *Timer) int {

	var delflag bool = false

	for temp := s.headTimer; temp != s.tailTimer; temp = temp.next {
		if temp == timer && temp.next == nil && temp.prev == nil {
			//only one timer in list
			s.headTimer = nil
			s.tailTimer = nil
			delflag = true
			break
		} else if temp == timer && temp.next == nil {
			//timer is tail timer in list
			s.tailTimer = temp.prev
			s.tailTimer.next = nil
			temp.prev = nil
			delflag = true
			break
		} else if temp == timer && temp.next != nil && temp.prev == nil {
			//timer is head timer in list
			s.headTimer = temp.next
			s.headTimer.prev = nil
			temp.prev = nil
			delflag = true
			break
		} else if temp == timer {
			//timer is in middle
			temp.prev.next = temp.next
			temp.next.prev = temp.prev
			temp.next, temp.prev = nil, nil
			delflag = true
			break
		}
	}

	if !delflag {
		//no such timer in list
		return TimerDeleteFailed
	}
	s.numTimers--
	return TimerDeleteSucceed
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
