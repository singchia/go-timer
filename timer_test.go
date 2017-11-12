package timer

import (
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
	tw.Stop()
	_, err := tw.Time(2, struct{}{}, nil, nil)
	if err == nil {
		t.Error("should be error")
		return
	}
	t.Log(err.Error())
	return
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
