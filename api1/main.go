package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"net/http"
	"tracing/logging"
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

	r.GET("/ping", func(c *gin.Context) {

		defaultLogger := logging.GetFileLogger(c)
		defaultLogger.Info("this is a test message")

		resp, err := tracelib.HTTPClient(c.Request.Context(), "GET", "http://go-api2:8081/pong", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response tracelib.Response
		unMarshallErr := json.Unmarshal(resp, &response)
		if unMarshallErr != nil {
			fmt.Print("error occured: ", unMarshallErr)
		}

		resp2, err2 := tracelib.HTTPClient(c.Request.Context(), "GET", "http://go-api4:8083/dong", nil)
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response2 tracelib.Response
		json.Unmarshal(resp2, &response2)

		response.Message = response.Message + " : " + response2.Message

		c.JSON(http.StatusOK, response)
	})

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
