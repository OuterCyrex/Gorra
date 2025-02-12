package GorraAPI

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
	"go.uber.org/zap"
	"io"
)

func InitTracer(name string, host string, port int) (opentracing.Tracer, io.Closer) {
	sampler := jaeger.NewConstSampler(true)
	sender := transport.NewHTTPTransport(fmt.Sprintf("http://%s:%d/api/traces", host, port))

	reporter := jaeger.NewRemoteReporter(sender)

	return jaeger.NewTracer(name, sampler, reporter)
}

func RawContextWithSpan(c *gin.Context) context.Context {
	ctx := context.Background()
	span, ok := c.Get("__span")

	if !ok {
		zap.S().Info(c.Request.URL.Path + "no tracer injected")
		return ctx
	}

	ctx = opentracing.ContextWithSpan(ctx, span.(opentracing.Span))
	return ctx
}
