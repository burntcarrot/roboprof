package utils

import (
	"sync"
	"time"
)

var timerPool = sync.Pool{}

func Sleep(d time.Duration, cancel <-chan struct{}) {
	timer, _ := timerPool.Get().(*time.Timer)
	if timer == nil {
		timer = time.NewTimer(d)
	} else {
		timer.Reset(d)
	}

	select {
	case <-timer.C:
	case <-cancel:
		if !timer.Stop() {
			<-timer.C
		}
	}

	timerPool.Put(timer)
}
