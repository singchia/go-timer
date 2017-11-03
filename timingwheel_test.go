package timer

import (
	"testing"
	"time"
)

func Test_calcuQuotients(t *testing.T) {
	quos := calcuQuotients(1000)
	t.Log(quos)

	quos = calcuQuotients(1024)
	t.Log(quos)

	quos = calcuQuotients(10000000)
	t.Log(quos)
	return
}

func Test_indexesPerWheel(t *testing.T) {
	tw := &timingwheel{}
	tw.SetMaxTicks(1000)
	ipw := tw.indexesPerWheel(1023)
	t.Log(ipw)
	return
}

func Test_topology(t *testing.T) {
	tw := newTimingwheel()
	tw.Start()
	ch := make(chan interface{})
	old := time.Now()
	tw.Time(5, 1, ch, nil)
	tw.Time(5, 2, ch, nil)
	tw.Time(25, 3, ch, nil)
	tw.Time(250, 4, ch, nil)
	tw.Time(250, 5, ch, nil)
	tw.Time(250, 6, ch, nil)
	tw.Time(10000, 7, ch, nil)
	tw.Time(10000, 8, ch, nil)
	topo, err := tw.Topology()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Logf("topology: %s", string(topo))

	<-ch
	elapse := time.Now().Sub(old)
	if elapse/10e6 != 500 {
		t.Errorf("time-diff bigger then 10ms: %d", elapse)
		return
	}
	t.Logf("time-diff: %d", elapse)
}
