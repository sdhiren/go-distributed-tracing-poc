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
	if err := tracelib.InitializeTracing("Service1", "http://localhost:14268/api/traces"); err != nil {
		log.Fatalf("Failed to initialize tracing: %v", err)
	}

	r := gin.Default()
	r.Use(TraceMiddleware())
	r.GET("/ping", func(c *gin.Context) {

		// start a new span for the incoming request
		// spanName := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		// ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		// ctx, span := tracelib.StartSpanFromContext(ctx, spanName)
		// defer span.End()


		// create a separate request/response span
		// requestResponseSpanName := fmt.Sprintf("%s %s %s", c.Request.Method, c.Request.URL.Path, "request | response" )
		// _, span2 := tracelib.StartSpanFromContext(ctx, requestResponseSpanName)
		// defer span2.End()

		// adding request in the request/response span
		// tracelib.AddRequestToSpan(span2, c.Request, nil)

		resp, err := tracelib.HTTPClient(c.Request.Context(), "GET", "http://localhost:8081/pong", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response tracelib.Response
		json.Unmarshal(resp, &response)
		fmt.Println("Response from service2: ", response)

		// adding response in the request/response span
		// tracelib.AddResponseToSpan(span2, response, nil)

		resp2, err2 := tracelib.HTTPClient(c.Request.Context(), "GET", "http://localhost:8083/dong", nil)
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response2 tracelib.Response
		json.Unmarshal(resp2, &response2)
		fmt.Println("Response from service4: ", response2)

		response.Message = response.Message + " : " + response2.Message

		c.JSON(http.StatusOK, response)
	})

	log.Fatal(r.Run(":8080"))
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
