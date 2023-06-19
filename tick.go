package timer

import (
	"sync/atomic"
	"time"

	"github.com/singchia/go-timer/v2/pkg/linker"
)

type tickOption struct {
	data       interface{}
	ch         chan *Event
	chOutside  bool
	handler    func(*Event)
	cyclically bool
}

type status int

const (
	statusAdd = iota
	statusWait
	statusFire
	statusCanceled
)

// the real shit
type tick struct {
	// user interface
	*tickOption

	// location
	id  linker.DoubID
	s   *slot
	tw  *timingwheel
	ipw []uint

	// meta
	duration   time.Duration
	delay      time.Duration
	insertTime time.Time

	// status
	fired  int64
	status status
}

// TODO revision
func (t *tick) Reset(data interface{}) error {
	ch := make(chan *operationRet)
	t.tw.operations <- &operation{
		tick:     t,
		operType: operReset,
		retCh:    ch,
		data:     data,
	}
	ret, ok := <-ch
	if !ok {
		return ErrOperationForceClosed
	}
	return ret.err
}

// Don't cancel the tick in timeout handler, it may cause the block when buffer is fulfilled.
// The cancel operation may not work when the tick is about to time up.
func (t *tick) Cancel() error {
	t.tw.mtx.RLock()
	defer t.tw.mtx.RUnlock()

	ch := make(chan *operationRet)
	t.tw.operations <- &operation{
		tick:     t,
		operType: operCancel,
		retCh:    ch,
	}
	ret, ok := <-ch
	if !ok {
		return ErrOperationForceClosed
	}
	return ret.err
}

// TODO revision
func (t *tick) Delay(d time.Duration) error {
	if t.cyclically {
		return ErrDelayOnCyclically
	}
	ch := make(chan *operationRet)
	t.tw.operations <- &operation{
		tick:     t,
		operType: operDelay,
		delay:    d,
		retCh:    ch,
	}
	ret, ok := <-ch
	if !ok {
		return ErrOperationForceClosed
	}
	return ret.err
}

func (t *tick) Fired() int64 {
	return atomic.LoadInt64(&t.fired)
}

func (t *tick) C() <-chan *Event {
	return t.ch
}

func (t *tick) InsertTime() time.Time {
	return t.insertTime
}

func (t *tick) Duration() time.Duration {
	return t.duration
}
