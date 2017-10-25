package timer

import (
	"errors"
	"time"

	"github.com/singchia/go-hammer/doublinker"
	scheduler "github.com/singchia/go-scheduler"
)

const (
	//10 milliseconds is min interval gotimer can accept
	MinTickInterval     time.Duration = 10 * time.Millisecond
	DefaultTickInterval time.Duration = time.Second
	DefaultMaxTicks     uint64        = 1024 * 1024 * 1024
)

type timingwheel struct {
	wheels   []*wheel
	interval time.Duration
	signal   chan struct{}
	max      uint64
	sch      scheduler.Scheduler
}

func newTimingwheel() *timingwheel {
	return &timingwheel{interval: DefaultTickInterval, sch: scheduler.NewScheduler()}
}

func (t *timingwheel) setMaxTicks(max uint64) {
	t.max = max
	nums := calcuQuotients(max)
	t.wheels = make([]*wheel, 0, len(nums))
	for position, num := range nums {
		t.wheels = append(t.wheels, newWheel(num, uint(position)))
	}
	return
}

func (t *timingwheel) setInterval(interval time.Duration) {
	t.interval = interval
}

func (t *timingwheel) time(d uint64, data interface{}, C chan interface{}, handler Handler) (*Tick, error) {
	if d == 0 || d > t.max-1 {
		return nil, errors.New("invalid duration")
	}
	ipw := t.indexesPerWheel(d)
	tick := &Tick{data: data, C: C, handler: handler, ipw: ipw, duration: d}
	t.wheels[len(ipw)-1].add(ipw[len(ipw-1)], tick)
	return tick, nil
}

func (t *timingwheel) start() {
	if t.wheels == nil {
		t.setMaxTicks(DefaultMaxTicks)
	}
	go func() {
		driver := time.NewTicker(t.interval)
		select {
		case <-driver.C:
			for _, wheel := range t.wheels {
				//the linker is the list of Tickers
				linker := wheel.incN(1)
				t.sch.PublishRequest(&scheduler.Request{Data: linker, Handler: t.handler})
			}
		case t.signal:
			//TODO
		}
	}()
}

func (t *timingwheel) handler(data interface{}) error {
	linker, ok := data.(*doublinker.Doublinker)
	//TODO
}

func (t *timingwheel) indexesPerWheel(d uint64) []uint {
	var ipw []uint
	var reminder uint64
	var quotient = d
	for _, wheel := range t.wheels {
		if quotient == 0 {
			break
		}
		reminder = quotient % uint64(wheel.numSlots)
		quotient = quotient / uint64(wheel.numSlots)
		ipw = append(ipw, uint(reminder))
	}
	return ipw
}

func calcuQuotients(num uint64) []uint {
	var quos []uint
	for {
		quo := num / 16
		rem := num % 16
		if quo != 0 {
			quos = append(quos, 16)
		} else {
			quos = append(quos, uint(rem))
			break
		}
		if rem != 0 {
			quo = quo + 1
		}
		num = quo
	}
	return quos
}
