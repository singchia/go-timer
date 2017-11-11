package timer

import (
	"errors"
	"fmt"
	"time"

	simplejson "github.com/bitly/go-simplejson"
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
	sch      *scheduler.Scheduler
}

func newTimingwheel() *timingwheel {
	return &timingwheel{interval: DefaultTickInterval, sch: scheduler.NewScheduler(), signal: make(chan struct{})}
}

func (t *timingwheel) SetMaxTicks(max uint64) {
	t.max = max
	nums := calcuQuotients(max)
	t.wheels = make([]*wheel, 0, len(nums))
	for position, num := range nums {
		t.wheels = append(t.wheels, newWheel(num, uint(position)))
	}
	return
}

func (t *timingwheel) SetInterval(interval time.Duration) {
	t.interval = interval
}

func (t *timingwheel) Time(d uint64, data interface{}, C chan interface{}, handler Handler) (Tick, error) {
	if d == 0 || d > t.max-1 {
		return nil, errors.New("invalid duration")
	}
	if data == nil {
		return nil, errors.New("invalid data")
	}
	if C == nil && handler == nil {
		C = make(chan interface{}, 1)
	}
	if t.wheels == nil {
		return nil, errors.New("timer not started")
	}
	ipw := t.indexesPerWheel(d)
	tick := &tick{data: data, C: C, handler: handler, ipw: ipw, duration: d}
	t.wheels[len(ipw)-1].add(ipw[len(ipw)-1], tick)
	return tick, nil
}

func (t *timingwheel) timeBased(d uint64, tick *tick) (*tick, error) {
	ipw := t.indexesPerWheelBased(d, tick.ipw)
	tick.ipw = ipw
	tick.duration += d
	t.wheels[len(ipw)-1].add(ipw[len(ipw)-1], tick)
	return tick, nil
}

func (t *timingwheel) Start() {
	if t.wheels == nil {
		t.SetMaxTicks(DefaultMaxTicks)
	}
	go t.drive()
	return
}

func (t *timingwheel) Pause() {
	t.signal <- struct{}{}
	return
}

func (t *timingwheel) Moveon() {
	go t.drive()
}

func (t *timingwheel) Stop() {
	t.signal <- struct{}{}
	t.wheels = nil
	t.sch.Close()
}

func (t *timingwheel) drive() {
	driver := time.NewTicker(t.interval)
	for {
		select {
		case <-driver.C:
			for _, wheel := range t.wheels {
				linker := wheel.incN(1)
				linker.Foreach(t.iterate)
				if wheel.cur != 0 {
					break
				}
			}
		case <-t.signal:
			return
		}
	}
	return
}

func (t *timingwheel) iterate(data interface{}) error {
	t.sch.PublishRequest(&scheduler.Request{Data: data, Handler: t.handle})
	return nil
}

func (t *timingwheel) handle(data interface{}) {
	tick, _ := data.(*tick)
	position := tick.s.w.position
	for position > 0 {
		position--
		if tick.ipw[position] > 0 {
			t.wheels[position].add(tick.ipw[position], tick)
			return
		}
	}
	if tick.C == nil {
		tick.handler(tick.data)
	} else {
		tick.C <- tick.data
	}
}

func (t *timingwheel) indexesPerWheel(d uint64) []uint {
	var ipw []uint
	var reminder uint64
	var quotient = d
	for i, wheel := range t.wheels {
		if quotient == 0 {
			break
		}
		quotient += uint64(t.wheels[i].cur)
		reminder = quotient % uint64(wheel.numSlots)
		quotient = quotient / uint64(wheel.numSlots)
		ipw = append(ipw, uint(reminder))
	}
	return ipw
}

func (t *timingwheel) indexesPerWheelBased(d uint64, base []uint) []uint {
	var ipw []uint
	var reminder uint64
	var quotient = d
	for i, wheel := range t.wheels {
		if quotient == 0 {
			break
		}
		quotient += uint64(base[i])
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

//for debug
func (t *timingwheel) Topology() ([]byte, error) {
	j := simplejson.New()
	var ws []*simplejson.Json
	for i, wheel := range t.wheels {
		var ss []*simplejson.Json
		for j, slot := range wheel.slots {
			var ts []interface{}
			slot.foreach(func(data interface{}) error {
				t := data.(*tick)
				ts = append(ts, t.data)
				return nil
			})
			s := simplejson.New()
			s.Set(fmt.Sprintf("slot%d", j), ts)
			ss = append(ss, s)
		}
		w := simplejson.New()
		w.Set(fmt.Sprintf("wheel%d", i), ss)
		ws = append(ws, w)
	}
	j.Set("wheels", ws)
	return j.MarshalJSON()
}
