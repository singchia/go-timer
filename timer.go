package gotimer

import (
	"time"
)

const (
	TimingWheel string = "timingwheel"
	MinHeap     string = "minheap"
)

type TimerID int64

type TimerHandler func(data interface{})

type baseTimer struct {
	timerCenter    TimerCenter
	d              time.Duration
	data           interface{}
	timerId        TimerID
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

//TimerCenter is a interface for managering basic timers,
//for now, only timingWheel implements the interface,
//users can create a timingWheel instance by calling
//NewTimerCenter("gotimer.TimingWheel"), it's kind of
//factory method, users can add timer to TimerCenter with
//no actions, it means when timer is up, TimerCenter will
//delete it with no actions. or users can add timer with
//one or more TimerHandler, when timer is up, timer's all
//TimerHandler will be called. TimerCenter won't start
//until StartTimer be called, it won't stop until
//StopTimer be called.
type TimerCenter interface {
	//
	AddTimer(d time.Duration) Timer

	AddTimerWithHandler(d time.Duration, data interface{}, handler TimerHandler) Timer

	AddTimerWithHandlers(d time.Duration, data interface{}, handler []TimerHandler) Timer

	DelTimer(timer Timer)

	//TimerCenter won't start until StartTimer be called
	StartTimer() error

	//TimerCenter won't stop until StopTimer be called
	StopTimer() error

	//set max amount of goroutines that handle TimerHandler,
	//a new groutine will be started when exist routines are
	//in busy, num < 0 means unlimited
	SetMaxGoRoutines(num int)
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
