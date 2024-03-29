package timer

import "time"

type TimerOption func(*timerOption)

func WithTimeInterval(interval time.Duration) TimerOption {
	return func(to *timerOption) {
		to.interval = interval
	}
}

func WithOperationBufferSize(n int) TimerOption {
	return func(to *timerOption) {
		to.operationBufferSize = n
	}
}

type TickOption func(*tickOption)

func WithData(data interface{}) TickOption {
	return func(to *tickOption) {
		to.data = data
	}
}

func WithCyclically() TickOption {
	return func(to *tickOption) {
		to.cyclically = true
	}
}

func WithChan(C chan *Event) TickOption {
	return func(to *tickOption) {
		to.ch = C
		to.chOutside = true
	}
}

func WithHandler(handler func(*Event)) TickOption {
	return func(to *tickOption) {
		to.handler = handler
	}
}

type Timer interface {
	Add(d time.Duration, opts ...TickOption) Tick

	// Close to close timer,
	// all ticks set would be discarded.
	Close()

	// Pause the timer,
	// all ticks won't continue after Timer.Movenon().
	Pause()

	// Continue the paused timer.
	Moveon()
}

// Tick that set in Timer can be required from Timer.Add()
type Tick interface {
	//To reset the data set at Timer.Time()
	Reset(data interface{}) error

	//To cancel the tick
	Cancel() error

	//Delay the tick
	Delay(d time.Duration) error

	//To get the channel called at Timer.Time(),
	//you will get the same channel if set, if not and handler is nil,
	//then a new created channel will be returned.
	C() <-chan *Event

	// Insert time
	InsertTime() time.Time

	// The tick duration original set
	Duration() time.Duration

	// Fired count
	Fired() int64
}

type Event struct {
	Duration   time.Duration
	InsertTIme time.Time
	Data       interface{}
	Error      error
}

// Entry
func NewTimer(opts ...TimerOption) Timer {
	return newTimingwheel(opts...)
}
