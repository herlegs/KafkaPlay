package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/herlegs/KafkaPlay/testratelimiter/rate"
)

var globalCount uint32

const (
	qps = 3000

	workerNum = 20
)

// main
func main() {
	limiter := rate.NewLimiter(rate.Limit(qps), 1)

	for i := 0; i < workerNum; i++ {
		countWithLimit(limiter)
	}

	countChecker()

	<-time.After(time.Minute)
}

// countWithLimit
func countWithLimit(limiter *rate.Limiter) {
	go func() {
		for {
			if limiter.Allow() {
				atomic.AddUint32(&globalCount, 1)
			}
		}
	}()
}

// countChecker checks whether counting speed is close to qps
func countChecker() {
	go func() {
		for {
			prev := atomic.LoadUint32(&globalCount)
			<-time.After(time.Second)
			curr := atomic.LoadUint32(&globalCount)
			fmt.Printf("increased per sec: %v\n", curr-prev)
		}
	}()
}
