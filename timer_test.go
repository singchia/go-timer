package timer

import (
	"math/rand"
	"testing"
	"time"
)

func TestGoTimer(t *testing.T) {
	tw := NewTimer()
	tw.Start()
	tick, _ := tw.Time(1, time.Now(), nil, nil)
	o := <-tick.Tunnel()
	elapse := time.Since(o.(time.Time))
	t.Log(elapse)
	return
}

func Test_Stop(t *testing.T) {
	tw := NewTimer()
	tw.Start()
	_, err := tw.Time(2, struct{}{}, nil, func(data interface{}) error {
		t.Log(data.(struct{}))
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	tw.Stop()
	time.Sleep(500 * time.Millisecond)
}

	tick := tw.Time(time.Second)
	<-tick.Channel()
	t.Logf("%v %v", tick.InsertTime(), time.Now())
}

func BenchmarkGoTimer(b *testing.B) {
	tw := NewTimer()
	tw.Start()

	for i := 0; i < b.N; i++ {
		tw.Time(time.Second, WithHandler(func(data interface{}) error {
			return nil
		}))
	}
}

func BenchmarkBuildinTimer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.NewTimer(time.Second)
	}
}

func Benchmark_Delay(b *testing.B) {
	tw := NewTimer()
	tw.Start()
	fired := 0
	for i := 0; i < b.N; i++ {
		second := uint64(rand.Intn(10) + 1)
		tick, err := tw.Time(second, time.Now(), nil, func(data interface{}) error {
			elapse := time.Since(data.(time.Time).Add(time.Duration(second) * time.Second))
			b.Log(elapse)
			fired++
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
		tick.Delay(2)
	}
	time.Sleep(13 * time.Second)
	b.Log(b.N, fired)
}
