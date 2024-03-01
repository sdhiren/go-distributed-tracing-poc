package tracelib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

func InitializeTracing(serviceName, jaegerEndpoint string) error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("v0.1.0"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return nil
}

func StartSpanFromContext(ctx context.Context, operationName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	tr := otel.Tracer("tracelib")
	ctx, span := tr.Start(ctx, operationName, trace.WithAttributes(attrs...))

	return ctx, span
}

func HTTPClient(ctx context.Context, method, url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	body, error := io.ReadAll(resp.Body)
	if error != nil {
	   fmt.Println(error)
	}
	resp.Body.Close()
	
	return body, nil
}

func AddAttributes(span trace.Span, attributes map[string]interface{}) {
    for key, value := range attributes {
		span.SetAttributes(attribute.String(key, value.(string)))
    }
}

func AddResponseToSpan(span trace.Span, message string, attrs map[string]any) {
	// Convert map[string]any to []attribute.KeyValue
	var attributes []attribute.KeyValue
	span.AddEvent("Response : " + string(message), trace.WithAttributes(attributes...), trace.WithTimestamp(time.Now()))
}

func AddRequestToSpan(span trace.Span, request any, attrs map[string]any) {
	// Convert map[string]any to []attribute.KeyValue
	var attributes []attribute.KeyValue

	m, err := json.Marshal(request)
	if err != nil {
		fmt.Println("error while reading the response")
	}
	
	span.AddEvent("Request : " + string(m), trace.WithAttributes(attributes...), trace.WithTimestamp(time.Now()))
}

type Response struct {
	Message string `json:"message"`
}
