package client

import (
	"sync/atomic"
	"time"
)

type CircuitBreaker interface {
	AllowRequest() bool
	Success()
	Fail(err error)
}

// DefaultCircuitBreaker 暂时实现了简单的基于时间窗口的熔断器
type DefaultCircuitBreaker struct {
	lastFail  time.Time
	fails     uint64
	threshold uint64
	window    time.Duration
}

func (cb *DefaultCircuitBreaker) AllowRequest() bool {
	if time.Since(cb.lastFail) > cb.window {
		cb.reset()
		return true
	}
	failures := atomic.LoadUint64(&cb.fails)
	return failures < cb.threshold
}

func (cb *DefaultCircuitBreaker) Success() {
	cb.reset()
}

func (cb *DefaultCircuitBreaker) Fail() {
	atomic.AddUint64(&cb.fails, 1)
	cb.lastFail = time.Now()
}

func (cb *DefaultCircuitBreaker) reset() {
	atomic.StoreUint64(&cb.fails, 0)
	cb.lastFail = time.Now()
}

func NewDefaultCircuitBreaker(threshold uint64, window time.Duration) *DefaultCircuitBreaker {
	return &DefaultCircuitBreaker{
		threshold: threshold,
		window:    window,
	}
}
