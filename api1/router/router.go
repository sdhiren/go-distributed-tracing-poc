package router

import (
	"bytes"
	"fmt"
	"log"
	"tracing/api1/controller"
	"tracing/tracelib"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func InitRoutes() *gin.Engine {
	if err := tracelib.InitializeTracing("Service1", "http://jaeger:14268/api/traces"); err != nil {
		log.Fatalf("Failed to initialize tracing: %v", err)
	}

	r := gin.Default()
	r.Use(TraceMiddleware())

	apiController1 := controller.NewApiController1()

	r.GET("/ping", apiController1.CallApi2)
	return r
}

func TraceMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        spanName := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)

		ctx := otel.GetTextMapPropagator().Extract(c, propagation.HeaderCarrier(c.Request.Header))
		ctx, span := tracelib.StartSpanFromContext(ctx, spanName)

		defer span.End()

        c.Set("span", span)
        c.Request = c.Request.Clone(ctx)
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
    	c.Writer = blw
        c.Next()
    }
}

type bodyLogWriter struct {
    gin.ResponseWriter
    body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
    w.body.Write(b)
    return w.ResponseWriter.Write(b)
}
