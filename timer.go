package gotimer

import(
	"time"
)

const (
	TimingWheel string = "timingwheel"
	MinHeap string = "minheap"
)

type TimerID int64

type TimerHander func(data interface{})

type baseTimer struct {
	timerMngr *TimerCenter
	d time.Duration
	data interface{}
	timerId TimerID
	expiryHandlers []TimerHandler
}

func (t *BaseTimer) ResetData(data interface{}) {
	t.data = data
}

func (t *BaseTimer) AddHander(handler TimerHandler) {
	t.expiryHandlers = append(t.expiryHanders, handler)
}

type Timer interface {
	ResetData(data interface{})
	AddHandler()
	Stop()
	Delay(d time.Duration)
}

//for managering timers
type TimerCenter interface {
	StartTimer(d time.Duration) *Timer

	StartTimerWithHander(d time.Duration, data interface{}, handler TimerHander) *Timer

	StartTimerWithHanders(d time.Duration, data interface{}, hander []TimerHander) *Timer

	StartTimerWithActions(d time.Duration, data interface{}, actions Actions) *Timer

	StopTimer(timer *Timer)
}

func NewTimerCenter(timerType string) TimerCenter {
	switch timerType {
	case TimingWheel:
	case MinHeap:
	default:
		return nil
	}
	return nil
}

//user should implements this
type Actions interface {
	Expiry(data interface{})
}

