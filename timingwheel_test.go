package timer

import "testing"

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
	tw.setMaxTicks(1000)
	ipw := tw.indexesPerWheel(1023)
	t.Log(ipw)
	return
}
