package timer

import "time"

type Timer interface {
	SetMaxTicks(max uint64)
	SetInterval(interval time.Duration)
	Time(d uint64, data interface{}, C chan interface{}, handler Handler) (Tick, error)
	Start()
	Stop()
	Pause()
	Moveon()
}

type Tick interface {
	Reset(data interface{}) error
	Cancel() error
	Delay(d uint64) error
}

func NewTimer() Timer {
	return newTimingwheel()
}
