package client

import (
	"context"
	"github.com/megaredfan/rpc-demo/share/metadata"
	"github.com/megaredfan/rpc-demo/share/trace"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracingLog "github.com/opentracing/opentracing-go/log"
	"log"
)

// OpenTracingInterceptor 目前只做了同步调用支持
type OpenTracingInterceptor struct {
	defaultClientInterceptor
}

func (*OpenTracingInterceptor) WrapCall(option *SGOption, callFunc CallFunc) CallFunc {
	return func(ctx context.Context, ServiceMethod string, arg interface{}, reply interface{}) error {
		var clientSpan opentracing.Span
		if ServiceMethod != "" { // 不是心跳的请求才进行追踪
			// 先从当前context获取已存在的追踪信息
			var parentCtx opentracing.SpanContext
			if parent := opentracing.SpanFromContext(ctx); parent != nil {
				parentCtx = parent.Context()
			}
			// 开始埋点
			clientSpan := opentracing.StartSpan(
				ServiceMethod,
				opentracing.ChildOf(parentCtx),
				ext.SpanKindRPCClient)
			defer clientSpan.Finish()
			// 将追踪信息注入到metadata中，通过rpc传递到下游
			meta := metadata.FromContext(ctx)
			writer := &trace.MetaDataCarrier{&meta}
			injectErr := opentracing.GlobalTracer().Inject(clientSpan.Context(), opentracing.TextMap, writer)
			if injectErr != nil {
				log.Printf("inject trace error: %v", injectErr)
			}
			ctx = metadata.WithMeta(ctx, meta)
		}

		err := callFunc(ctx, ServiceMethod, arg, reply)
		if err != nil && clientSpan != nil {
			clientSpan.LogFields(opentracingLog.String("error", err.Error()))
		}
		return err
	}
}
