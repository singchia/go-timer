package scheduler

import (
	"sync"
	"sync/atomic"
	"time"
)

type Scheduler struct {
	//the interval of goroutines number changing, min 200ms, default 1s
	//should be set before Schedule called
	Interval time.Duration

	incomingChan chan *Request

	//gorouties close chan
	closeChan chan struct{}

	//for shutting down whole scheduler
	allCloseChan chan struct{}

	//count of whole incoming requests
	countIncomingReqs int64

	//count of whole incoming requests until (now - interval)
	countIncomingReqsL int64

	//count of whole processed requests
	countProcessedReqs int64

	//count of whole processed requests until (now - interval)
	countProcessedReqsL int64

	//number of goroutines
	numActives int64

	defaultHandler Handler
	monitor        Monitor
	strategy       Strategy

	//runtimeLock locks strategy, defaultHandler, any changes refer to them locked
	runtimeLock sync.RWMutex
}

type Monitor func(incomingReqsLastInterval, processedReqsLastInterval, shift, numActives int64)
type Handler func(data interface{})

type Request struct {
	Data    interface{}
	Handler Handler
}

func NewScheduler() *Scheduler {
	var initialNum int64 = 1
	scheduler := &Scheduler{
		Interval:            time.Millisecond * 200,
		strategy:            NewGradienter(),
		countIncomingReqs:   0,
		countIncomingReqsL:  0,
		countProcessedReqs:  0,
		countProcessedReqsL: 0,
		numActives:          0,
		incomingChan:        make(chan *Request, 1024),
		closeChan:           make(chan struct{}, 1024),
		allCloseChan:        make(chan struct{})}
	scheduler.expandGoRoutines(initialNum)
	return scheduler
}

func (s *Scheduler) SetDefaultHandler(handler Handler) {
	s.runtimeLock.Lock()
	defer s.runtimeLock.Unlock()
	s.defaultHandler = handler
}

func (s *Scheduler) SetMonitor(monitor Monitor) {
	s.runtimeLock.Lock()
	defer s.runtimeLock.Unlock()
	s.monitor = monitor
}

func (s *Scheduler) SetStrategy(strategy Strategy) {
	s.runtimeLock.Lock()
	defer s.runtimeLock.Unlock()
	s.strategy = strategy
}

// max number of active goroutines, -1 means not limited
// default -1, can be set at any runtime
func (s *Scheduler) SetMaxGoroutines(maxCountGoroutines int64) {
	s.runtimeLock.Lock()
	defer s.runtimeLock.Unlock()
	s.strategy.SetMaxActives(maxCountGoroutines)
}

// max number of request processed per second, -1 means not limited
// default -1, can be set at any runtime
func (s *Scheduler) SetMaxProcessedReqs(maxProcessedReqs int64) {
	s.runtimeLock.Lock()
	defer s.runtimeLock.Unlock()
	s.strategy.SetMaxProcessedReqs(maxProcessedReqs)
}

// max rate of (incoming requests)/(processed requests), should between 0 and 1
// default 1
func (s *Scheduler) SetMaxRate(rate float64) {
	s.runtimeLock.Lock()
	defer s.runtimeLock.Unlock()
	s.strategy.SetMaxRate(rate)
}

func (s *Scheduler) StartSchedule() {
	go s.control()
}

func (s *Scheduler) control() {
	if s.Interval < time.Microsecond*500 {
		s.Interval = time.Millisecond * 200
	}

	timer := time.NewTicker(s.Interval)
	for {
		select {
		case _, ok := <-timer.C:
			if !ok {
				return
			}
			numActives := atomic.LoadInt64(&s.numActives)
			incomingReqsDiff := s.countIncomingReqs - s.countIncomingReqsL
			processedReqsDiff := s.countProcessedReqs - s.countProcessedReqsL
			s.countIncomingReqsL = s.countIncomingReqs
			s.countProcessedReqsL = s.countProcessedReqs

			s.runtimeLock.RLock()
			shift := s.strategy.ExpandOrShrink(incomingReqsDiff, processedReqsDiff, numActives)
			s.runtimeLock.RUnlock()

			if shift < 0 && numActives == 1 {
				// at least reserve 1 to handle request
				if s.monitor != nil {
					s.monitor(incomingReqsDiff, processedReqsDiff, shift, numActives)
				}
				continue
			} else {
				if s.monitor != nil {
					s.monitor(incomingReqsDiff, processedReqsDiff, shift, numActives+shift)
				}
			}
			if shift > 0 {
				s.expandGoRoutines(shift)
			} else {
				s.shrinkGoRoutines(-shift)
			}
		case <-s.allCloseChan:
			return
		}
	}
}

func (s *Scheduler) shrinkGoRoutines(num int64) {
	var i int64
	for i = 0; i < num; i++ {
		s.closeChan <- struct{}{}
	}
}

func (s *Scheduler) expandGoRoutines(num int64) {
	var i int64
	for i = 0; i < num; i++ {
		go func() {
			atomic.AddInt64(&s.numActives, 1)
			for {
				select {
				case r, ok := <-s.incomingChan:
					if !ok {
						atomic.AddInt64(&s.numActives, -1)
						return
					}
					if r.Handler == nil && s.defaultHandler != nil {
						s.runtimeLock.RLock()
						s.defaultHandler(r.Data)
						s.runtimeLock.RUnlock()

					} else if r.Handler != nil {
						r.Handler(r.Data)
					}
					atomic.AddInt64(&s.countProcessedReqs, 1)
				case <-s.closeChan:
					atomic.AddInt64(&s.numActives, -1)
					return
				}
			}
		}()
	}
}

func (s *Scheduler) PublishRequest(req *Request) {
	s.incomingChan <- req
	atomic.AddInt64(&s.countIncomingReqs, 1)
}

func (s *Scheduler) Close() {
	close(s.incomingChan)
	close(s.allCloseChan)
}
