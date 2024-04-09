package tracelib

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func TestInitializeTracing(t *testing.T) {
    err := InitializeTracing("test-service", "jaeger-endpoint")

    assert.Nil(t, err)
}

func TestStartSpanFromContext(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		_, span := StartSpanFromContext(ctx, "operation")
		assert.NotNil(t, span)
		span.End()
	})

	t.Run("WithAttributes", func(t *testing.T) {
		attrs := []attribute.KeyValue{
			semconv.HTTPMethodKey.String("GET"),
			semconv.HTTPURLKey.String("/example"),
		}
		_, span := StartSpanFromContext(ctx, "operation", attrs...)
		assert.NotNil(t, span)
		span.End()
	})
}