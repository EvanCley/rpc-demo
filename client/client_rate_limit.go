package client

import (
	"context"
	"errors"
	"github.com/megaredfan/rpc-demo/share/ratelimit"
)

type RateLimitInterceptor struct {
	defaultClientInterceptor
	Limit ratelimit.RateLimiter
}

var ErrRateLimited = errors.New("request limited")

func (r *RateLimitInterceptor) WrapCall(option *SGOption, callFunc CallFunc) CallFunc {
	return func(ctx context.Context, ServiceMethod string, arg interface{}, reply interface{}) error {
		if r.Limit != nil {
			if r.Limit.TryAcquire() { // 进行尝试获取，获取失败时直接返回限流异常
				return callFunc(ctx, ServiceMethod, arg, reply)
			} else {
				return ErrRateLimited
			}
		} else { // 若限流器为 nil 则不进行限流
			return callFunc(ctx, ServiceMethod, arg, reply)
		}
	}
}

func (r *RateLimitInterceptor) WrapGo(option *SGOption, goFunc GoFunc) GoFunc {
	return func(ctx context.Context, ServiceMethod string, arg interface{}, reply interface{}, done chan *Call) *Call {
		if r.Limit != nil {
			if r.Limit.TryAcquire() { // 进行尝试获取，获取失败时直接返回限流异常
				return goFunc(ctx, ServiceMethod, arg, reply, done)
			} else {
				call := &Call{
					ServiceMethod: ServiceMethod,
					Args:          arg,
					Reply:         nil,
					Error:         ErrRateLimited,
					Done:          done,
				}
				done <- call
				return call
			}
		} else { // 若限流器为nil则不进行限流
			return goFunc(ctx, ServiceMethod, arg, reply, done)
		}
	}
}
