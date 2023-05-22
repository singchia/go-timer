package timer

import (
	"math/rand"
	"testing"
	"time"
)

func Test_SetInterval(t *testing.T) {
	tw := NewTimer()
	tw.SetInterval(time.Second * 5)
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

func Test_Pause(t *testing.T) {
	tw := NewTimer()
	tw.Start()
	tick, _ := tw.Time(2, time.Now(), nil, nil)
	tw.Pause()
	select {
	case <-tick.Tunnel():
		t.Error("not paused")
	case <-time.NewTimer(time.Second * 5).C:
		t.Log("paused and timeout")
	}
	return
}

func Test_Moveon(t *testing.T) {
	tw := NewTimer()
	tw.Start()
	tick, _ := tw.Time(2, time.Now(), nil, nil)
	tw.Pause()
	time.Sleep(time.Second)
	tw.Moveon()
	o := <-tick.Tunnel()
	elapse := time.Since(o.(time.Time))
	t.Log(elapse)
	return
}

func Test_Reset(t *testing.T) {
	tw := NewTimer()
	tw.Start()
	tick, _ := tw.Time(2, 1, nil, nil)
	tick.Reset(2)
	o := <-tick.Tunnel()
	if o.(int) != 2 {
		t.Error(o)
		return
	}
	return
}

func Test_Delay(t *testing.T) {
	tw := NewTimer()
	tw.Start()
	tick, _ := tw.Time(1, time.Now(), nil, nil)
	tick.Delay(2)
	o := <-tick.Tunnel()
	elapse := time.Since(o.(time.Time))
	t.Log(elapse)
	return
}

func Test_Cancel(t *testing.T) {
	tw := NewTimer()
	tw.Start()
	tick, _ := tw.Time(1, time.Now(), nil, nil)
	err := tick.Cancel()
	if err != nil {
		t.Error(err)
		return
	}
	select {
	case <-tick.Tunnel():
		t.Error("not paused")
	case <-time.NewTimer(time.Second * 5).C:
		t.Log("canceled and timeout")
	}
	return

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
