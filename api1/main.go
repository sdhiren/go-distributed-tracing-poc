package main

import (
	"bytes"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"tracing/api1/controller"
	"tracing/tracelib"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {

	if err := tracelib.InitializeTracing("Service1", "http://jaeger:14268/api/traces"); err != nil {
		log.Fatalf("Failed to initialize tracing: %v", err)
	}

	r := gin.Default()
	r.Use(TraceMiddleware())

	apiController1 := controller.NewApiController1()

	r.GET("/ping", apiController1.CallApi2)

	log.Fatal(r.Run(":8080"))
}

func testHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     nil,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.TimeValue(time.Time{})
			}

			return a
		},
	})
}

func TraceMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        spanName := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)

		ctx := otel.GetTextMapPropagator().Extract(c, propagation.HeaderCarrier(c.Request.Header))
		ctx, span := tracelib.StartSpanFromContext(ctx, spanName)

		defer span.End()

		requestResponseSpanName := fmt.Sprintf("%s %s %s", c.Request.Method, c.Request.URL.Path, "request | response" )
		_, span2 := tracelib.StartSpanFromContext(ctx, requestResponseSpanName)
		defer span2.End()

		tracelib.AddRequestToSpan(span2, c.Request, nil)

        c.Set("span", span)
        c.Request = c.Request.Clone(ctx)
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
    	c.Writer = blw
        c.Next()

		tracelib.AddResponseToSpan(span2, blw.body.String(), nil)

        spanCtx := span.SpanContext()
        c.Writer.Header().Set("traceparent", spanCtx.TraceID().String())

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
