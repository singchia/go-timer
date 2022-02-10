package timer

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkGoTimer(b *testing.B) {
	tw := NewTimer()
	tw.Start()
	wait := new(sync.WaitGroup)
	for i := 0; i < b.N; i++ {
		wait.Add(1)
		tw.Time(1, wait, nil, func(data interface{}) error {
			wait.Done()
			return nil
		})
	}
	wait.Wait()
}

func BenchmarkBuildTimer(b *testing.B) {
	wait := new(sync.WaitGroup)
	for i := 0; i < b.N; i++ {
		wait.Add(1)
		timer := time.NewTimer(time.Second)
		go func(t *time.Timer) {
			<-t.C
			wait.Done()
		}(timer)
	}
	wait.Wait()
}

func TestCompare(t *testing.T) {
	count := 10000000
	s := time.Now()
	tw := NewTimer()
	tw.Start()
	wait := new(sync.WaitGroup)
	for i := 0; i < count; i++ {
		wait.Add(1)
		tw.Time(1, wait, nil, func(data interface{}) error {
			wait.Done()
			return nil
		})
	}
	wait.Wait()
	diff := time.Now().Sub(s).Milliseconds()
	t.Log("go-timer:", diff)

	s = time.Now()
	wait = new(sync.WaitGroup)
	for i := 0; i < count; i++ {
		wait.Add(1)
		timer := time.NewTimer(time.Second)
		go func(t *time.Timer) {
			<-t.C
			wait.Done()
		}(timer)
	}
	wait.Wait()
	diff = time.Now().Sub(s).Milliseconds()
	t.Log("buildin timer:", diff)
}

func TestSetInterval(t *testing.T) {
	tw := NewTimer()
	tw.SetInterval(time.Second * 5)
	tw.Start()
	tick, _ := tw.Time(1, time.Now(), nil, nil)
	o := <-tick.Tunnel()
	elapse := time.Since(o.(time.Time))
	t.Log(elapse)
	return
}

func TestStop(t *testing.T) {
	tw := NewTimer()
	tw.Start()
	tw.Stop()
	_, err := tw.Time(2, struct{}{}, nil, nil)
	if err == nil {
		t.Error("should be error")
		return
	}
	t.Log(err.Error())
	return
}

func TestPause(t *testing.T) {
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

func TestMoveon(t *testing.T) {
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

func TestReset(t *testing.T) {
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

func TestDelay(t *testing.T) {
	tw := NewTimer()
	tw.Start()
	tick, _ := tw.Time(1, time.Now(), nil, nil)
	tick.Delay(2)
	o := <-tick.Tunnel()
	elapse := time.Since(o.(time.Time))
	t.Log(elapse)
	return
}

func TestCancel(t *testing.T) {
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
