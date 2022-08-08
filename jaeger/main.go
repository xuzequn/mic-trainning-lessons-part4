package main

import (
	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"time"
)

func main() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831",
		},
		ServiceName: "xzqMall",
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	defer closer.Close()
	if err != nil {
		panic(err)
	}
	parentSpan := tracer.StartSpan("order_web")

	caetSpan := tracer.StartSpan("cart_srv", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(1 * time.Second)
	caetSpan.Finish()

	productSpan := tracer.StartSpan("product_srv", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(2 * time.Second)
	productSpan.Finish()

	stockSpan := tracer.StartSpan("stock_srv", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(2 * time.Second)
	stockSpan.Finish()

	parentSpan.Finish()
}
