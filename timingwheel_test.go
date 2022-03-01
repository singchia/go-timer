package timer

import (
	"testing"
)

func Test_indexesPerWheel(t *testing.T) {
	tw := &timingwheel{}
	ipw := tw.indexesPerWheel(1023)
	t.Log(ipw)
	return
}
