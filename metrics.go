package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type producerMetrics struct {
	sentBatches uint64
	errors      uint64
	done        chan bool
	wg          *sync.WaitGroup
}

func CreateMetrics() *producerMetrics {
	var wg sync.WaitGroup
	return &producerMetrics{0, 0, make(chan bool, 1), &wg}
}

func (metrics *producerMetrics) StartReporting() {
	started := metrics.SentBatches()
	ticker := time.NewTicker(time.Second)
	go func() {
		metrics.wg.Add(1)
		defer metrics.wg.Done()
		for range ticker.C {
			current := metrics.SentBatches()
			LogInfo(fmt.Sprintf("%d batches/s", current-started))
			started = current
			select {
			case <-metrics.done:
				return
			default:
				continue
			}
		}
	}()
}

func (metrics *producerMetrics) StopReporting() {
	metrics.done <- true
	metrics.wg.Wait()
	LogInfo(fmt.Sprintf("Sent batches: %d", metrics.SentBatches()))
	LogInfo(fmt.Sprintf("Errors while sending: %d", metrics.Errors()))
}

func (metrics *producerMetrics) AddBatch() {
	atomic.AddUint64(&metrics.sentBatches, 1)
}

func (metrics *producerMetrics) AddError() {
	atomic.AddUint64(&metrics.errors, 1)

}

func (metrics *producerMetrics) Errors() uint64 {
	return atomic.LoadUint64(&metrics.errors)
}

func (metrics *producerMetrics) SentBatches() uint64 {
	return atomic.LoadUint64(&metrics.sentBatches)
}
