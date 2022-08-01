package ratelimit

import (
	"errors"
	"time"
)

type RateLimiter interface {
	Acquire()                                        // 获取许可，会阻塞直到获得许可
	TryAcquire() bool                                // 尝试获取许可，如果不成功会立即返回false，而不是一直阻塞
	AcquireWithTimeout(duration time.Duration) error // 获取许可，会阻塞直到获得许可或者超时，超时时会返回一个超时异常，成功时返回nil
}

type DefaultRateLimiter struct {
	Num         int64
	rateLimiter chan time.Time
}

func NewRateLimiter(numPerSecond int64) RateLimiter {
	r := new(DefaultRateLimiter)
	r.Num = numPerSecond
	r.rateLimiter = make(chan time.Time)
	go func() {
		d := time.Duration(numPerSecond)
		ticker := time.NewTicker(time.Second / d)
		for t := range ticker.C {
			r.rateLimiter <- t
		}
	}()

	return r
}

func (r *DefaultRateLimiter) Acquire() {
	<-r.rateLimiter
}

func (r *DefaultRateLimiter) TryAcquire() bool {
	select {
	case <-r.rateLimiter:
		return true
	default:
		return false
	}
}

func (r *DefaultRateLimiter) AcquireWithTimeout(timeout time.Duration) error {
	ticker := time.NewTicker(timeout)
	select {
	case <-r.rateLimiter:
		return nil
	case <-ticker.C:
		return errors.New("acquire timeout")

	}
}

type RateLimitWrapper struct {
	global       RateLimiter
	methodLimits map[string]RateLimiter //Service.Method为key
}
