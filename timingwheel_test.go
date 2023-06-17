package timer

import (
	"testing"
	"time"
)

func Test_indexesPerWheel(t *testing.T) {
	tw := &timingwheel{
		timerOption: &timerOption{
			interval: time.Millisecond,
		},
	}
	ipw := tw.indexesPerWheel(1023)
	t.Log(ipw)
	return
}

func Test_calcuWheels(t *testing.T) {
	max, length := calcuWheels(defaultTicks)
	t.Log(max, length)
}
