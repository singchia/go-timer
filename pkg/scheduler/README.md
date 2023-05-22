# go-scheduler
A simple scheduler for goroutines, **go-scheduler** helps to manage goroutines, only needs to set three optional quotas:  

* the maximum count of goroutines 
* the maximum count of processed requests per interval 
* the maximum value of rate (processed requests per interval/ incoming requests per interval) 

Actually **go-sheduler** only adjust count of goroutines to satisfy those quotas if set, the default strategy works like gradienter, if runtime statistics don't match any quotas, **go-sheduler** starts to work.  
Since scheduler manage goroutines to handle user's **_Request_** which contains **_Data_** and **_Handler_**, the scheduler simple call **_Request.Handler(Request.Data)_**.  
**note:**  three optional quotas are only undercontrolled in **go-scheduler**

## How-to-use
**go-scheduler** supplies several easy-understand and easy-integrate interfaces, Let's see a easy sample.
```	golang
import (
    "fmt"
    "sync/atomic"
    "time"
    scheduler "github.com/singchia/go-scheduler"
)

func main() {
    sch := scheduler.NewScheduler()
	
    sch.SetMaxGoroutines(5000)
    sch.StartSchedule()
	
    var val int64
    for i := 0; i < 10*10000*10000; i++ {
        sch.PublishRequest(&scheduler.Request{Data: val, Handler: SchedulerHandler})
        atomic.AddInt64(&val, 1)
    }
    time.Sleep(time.Second * 5)
    fmt.Printf("maxValue: %d\n", maxValue)
    sch.Close()
}
	
var maxValue int64 = 0
	
func SchedulerHandler(data interface{}) {
    val, ok := data.(int64)
    if ok {
        if val > maxValue {
            maxValue = val
        }
    }
}
```  
It's not a good sample in production environment, but it does illustrate the usage  of **go-scheduler**. After **_SetMaxGoroutines(5000)_**, the max count of scheduler's goroutines shouldn't go beyond the range **_5000_**, use **_StartSchedule_** to start the scheduler, publish the **_Request_** into the scheduler by using **_PublishRequest_**, then scheduler will handle the request undercontrol.

## Installation
If you don't have the Go development environment installed, visit the [Getting Started](https://golang.org/doc/install) document and follow the instructions. Once you're ready, execute the following command:
```
go get -u github.com/singchia/go-scheduler
```

## Interfaces
#### _Scheduler.Interval_
This should be set before call **_StartSchedule_** and bigger than **_500us_**, if not set or less than **_500us_**, default 200ms.

#### _Scheduler.SetMaxGoroutines(int64)_
This limits the max count of goroutines in **go-scheduler**, can be set at any time.

#### _Scheduler.SetMaxProcessedReqs(int64)_
This limits the max processed requests per interval, can be set at any time.

#### _Scheduler.SetMaxRate(float64)_
The rate is the value of processed requests / incoming requests, bigger means you want a faster speed to handle requests, can be set at any time.

#### _Scheduler.SetDefaultHandler(scheduler.Handler)_
If you want set a default handler when **_scheduler.Request.Handler_** not given, can be set at any time.

#### _Scheduler.SetMonitor(scheduler.Monitor)_
You can use this to monitor incoming requests, processed requests, shift（changing of goroutines）, count of goroutines this interval, can be set at any time.

#### _Scheduler.SetStrategy(scheduler.Strategy)_
**_scheduler.Strategy_** is the key deciding how to shift(update) the count of goroutines, you can replace it as your own strategy.

## Strategy
#### _scheduler.Gradienter_
Defaultly **go-scheduler** uses _Gradienter_ as strategy, it behaves like:
```
if incoming requests == 0 then shrink 20%
if any quotas > max quotas then shrink the count of goroutines
	if quotas == 1 then shrink directly to MaxGoroutines
	else shrink 20%  
if all quotas < max quotas then expand randFloat * incomingReqs / (incomingReqs + maxCountGoroutines) * (maxCountGoroutines - currentCountGoroutines)
```

#### other strategies
In scheduler file, a circularLink.go exists, I was trying to look for next goroutines-updating by using history status, but temporarily no idea came up, if you have some idea welcome to contact me.
