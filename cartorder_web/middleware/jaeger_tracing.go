package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"mic-trainning-lesson-part4/internal"
)

func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		jaegerAddr := fmt.Sprintf("%s:%d", internal.AppConf.JaegerConfig.AgentHost,
			internal.AppConf.JaegerConfig.AgentPort)
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: jaegerAddr,
			},
			ServiceName: "xzqMall",
		}
		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		defer closer.Close()
		if err != nil {
			panic(err)
		}
		span := tracer.StartSpan(c.Request.URL.Path)
		c.Set("tracer", tracer)
		c.Set("parentSpan", span)
		c.Next()
	}
}
