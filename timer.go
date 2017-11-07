package timer

import "time"

type Timer interface {
	//SetMaxTicks sets the max number of ticks,
	//it should be called before Start if you want to customize it.
	//Default 1024*1024*1024 ticks
	SetMaxTicks(max uint64)

	//SetInterval set the time lapse between two ticks.
	//it should be called before Start if you want to customize it.
	//Default 1 second.
	//And the whole time lapse should be maxTicks*interval.
	//Note that the max tick you can preset is (maxTicks*interval-1),
	//because current tick cannot be set.
	SetInterval(interval time.Duration)

	//Time preset a Tick which will be triggered after d ticks,
	//you can set channel C, and after d ticks, data would be consumed from C.
	//Or you can set func handler, after d ticks, data would be handled
	//by handler in go-timer. If neither one be set, go-timer will generate a channel,
	//it's attatched with return value Tick, get it by Tick.Tunnel().
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
