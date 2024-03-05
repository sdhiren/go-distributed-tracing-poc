package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tracing/tracelib"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	if err := tracelib.InitializeTracing("Service2", "http://jaeger:14268/api/traces"); err != nil {
		log.Fatalf("Failed to initialize tracing: %v", err)
	}

	r := gin.Default()
	r.Use(TraceMiddleware())
	r.GET("/pong", func(c *gin.Context) {

		// spanName := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)

		// ctx := otel.GetTextMapPropagator().Extract(c, propagation.HeaderCarrier(c.Request.Header))
		// _, span := tracelib.StartSpanFromContext(ctx, spanName)

		// defer span.End()

		resp, err := tracelib.HTTPClient(c.Request.Context(), "GET", "http://go-api3:8082/ding", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response tracelib.Response
		json.Unmarshal(resp, &response)
		fmt.Println("Response from service3: ", response)
		c.JSON(http.StatusOK, response)
	})

	log.Fatal(r.Run(":8081"))
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
        // c.Request = c.Request.WithContext(ctx)
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