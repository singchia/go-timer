package gotimer

import(
	"time"
)

const (
	TimingWheel string = "timingwheel"
	MinHeap string = "minheap"
)

type TimerID int64

type TimerHandler func(data interface{})

type baseTimer struct {
	timerCenter TimerCenter
	d time.Duration
	data interface{}
	timerId TimerID
	expiryHandlers []TimerHandler
}

func (t *baseTimer) ResetData(data interface{}) {
	t.data = data
}

func (t *baseTimer) AddHandler(handler TimerHandler) {
	t.expiryHandlers = append(t.expiryHandlers, handler)
}

func (t *baseTimer) AddHandlers(handlers []TimerHandler) {
	t.expiryHandlers = append(t.expiryHandlers, handlers...)
}

type Timer interface {
	ResetData(data interface{})
	AddHandler(handler TimerHandler)
	AddHandlers(handlers []TimerHandler)
	Stop()
	Delay(d time.Duration)
}

//for managering timers
type TimerCenter interface {
	StartTimer(d time.Duration) Timer

	StartTimerWithHandler(d time.Duration, data interface{}, handler TimerHandler) Timer

	StartTimerWithHandlers(d time.Duration, data interface{}, handler []TimerHandler) Timer

	StartTimerWithActions(d time.Duration, data interface{}, actions Actions) Timer

	StopTimer(timer Timer)
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

