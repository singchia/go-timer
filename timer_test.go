package timer

import (
	"testing"
	"time"
)

func TestGoTimer(t *testing.T) {
	tw := NewTimer()
	tw.Start()

	tick := tw.Add(time.Second)
	<-tick.C()
	t.Logf("%v %v", tick.InsertTime(), time.Now())
}

func BenchmarkGoTimer(b *testing.B) {
	tw := NewTimer()
	tw.Start()

	for i := 0; i < b.N; i++ {
		tw.Add(time.Second, WithHandler(func(event *Event) {
		}))
	}
}

func BenchmarkBuildinTimer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.NewTimer(time.Second)
	}
}
