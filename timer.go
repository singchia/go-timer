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
	//Time must be called after Timer.Start.
	Time(d uint64, data interface{}, C chan interface{}, handler Handler) (Tick, error)

	//Start to start timer.
	Start()

	//Stop to stop timer,
	//all ticks set would be discarded.
	Stop()

	//Pause the timer,
	//all ticks won't continue after Timer.Movenon().
	Pause()

	//Continue the paused timer.
	Moveon()
}

//Tick that set in Timer can be required from Timer.Time()
type Tick interface {
	//To reset the data set at Timer.Time()
	Reset(data interface{}) error

	//To cancel the tick
	Cancel() error

	//Delay the tick if not timeout
	Delay(d uint64) error

	//To get the channel called at Timer.Time(),
	//you will get the same channel if set, if not and handler is nil,
	//then a new created channel will be returned.
	Tunnel() <-chan interface{}
}

//Entry
func NewTimer() Timer {
	return newTimingwheel()
}
