package timer

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestGoTimer(t *testing.T) {
	tw := NewTimer()

	tick := tw.Add(time.Second)
	<-tick.C()
	t.Logf("%v %v", tick.InsertTime(), time.Now())
}

func BenchmarkBuildinTimer(b *testing.B) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	startAllocs := memStats.Mallocs
	b.ResetTimer()
	b.ReportAllocs()
	wg := sync.WaitGroup{}
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		timer := time.NewTimer(time.Second)
		go func() {
			defer wg.Done()
			<-timer.C
		}()
	}
	wg.Wait()
	runtime.ReadMemStats(&memStats)
	endAllocs := memStats.Mallocs
	b.ReportMetric(float64(endAllocs-startAllocs), "allocs")
	b.ReportMetric(float64(memStats.TotalAlloc)/1024, "allocs_kb")
	b.ReportMetric(float64(memStats.Sys)/1024/1024, "sys_mb")
	b.ReportMetric(float64(memStats.HeapAlloc)/1024, "heap_alloc_kb")
	b.ReportMetric(float64(memStats.HeapSys)/1024/1024, "heap_sys_mb")
	b.ReportMetric(float64(memStats.StackSys)/1024/1024, "stack_sys_mb")
	b.ReportMetric(float64(runtime.NumCPU()), "num_cpu")
}

func BenchmarkGoTimer(b *testing.B) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	startAllocs := memStats.Mallocs
	b.ResetTimer()
	b.ReportAllocs()
	wg := sync.WaitGroup{}
	wg.Add(b.N)

	tw := NewTimer()
	for i := 0; i < b.N; i++ {
		tw.Add(time.Second, WithHandler(func(event *Event) {
			wg.Done()
		}))
	}
	wg.Wait()

	runtime.ReadMemStats(&memStats)
	endAllocs := memStats.Mallocs
	b.ReportMetric(float64(endAllocs-startAllocs), "allocs")
	b.ReportMetric(float64(memStats.TotalAlloc)/1024, "allocs_kb")
	b.ReportMetric(float64(memStats.Sys)/1024/1024, "sys_mb")
	b.ReportMetric(float64(memStats.HeapAlloc)/1024, "heap_alloc_kb")
	b.ReportMetric(float64(memStats.HeapSys)/1024/1024, "heap_sys_mb")
	b.ReportMetric(float64(memStats.StackSys)/1024/1024, "stack_sys_mb")
	b.ReportMetric(float64(runtime.NumCPU()), "num_cpu")
}
