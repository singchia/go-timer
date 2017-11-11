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
