package logging

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
)

type customLoggerTestSuite struct{
	suite.Suite
	context context.Context

}

func TestGetDefaultLogger_ValidContext(t *testing.T) {
	// Mock Gin context
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = httptest.NewRequest("GET", "/", nil)

	span := trace.SpanFromContext(context.Request.Context())
	spanId := span.SpanContext().SpanID().String()
	traceId := span.SpanContext().TraceID().String()

	os.Setenv("LOG_LEVEL", "INFO")

	handlerOpts := &slog.HandlerOptions{		
		Level:     slog.LevelInfo,
	}
	logger2 := slog.New(slog.NewJSONHandler(os.Stdout, handlerOpts))
	slog.SetDefault(logger2)
	
	slogLogger := logger2.With(
								slog.String("trace_id", traceId),
								slog.String("span_id", spanId),
								slog.String("end_point", "/"),
							  )
	expectedLogger := NewCustomLogger(slogLogger, context)						  
							  
	logger := GetDefaultLogger(context)
	assert.Equal(t, expectedLogger, logger)
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.logger)

}

func TestGetDefaultLogger_NilContext(t *testing.T) {
	logger := GetDefaultLogger(nil)
	assert.Nil(t, logger)
}

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name          string
		logLevel      string
		expectedLevel slog.Level
	}{
		{"DebugLevel", "DEBUG", slog.LevelDebug},
		{"InfoLevel", "INFO", slog.LevelInfo},
		{"WarnLevel", "WARN", slog.LevelWarn},
		{"ErrorLevel", "ERROR", slog.LevelError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv("LOG_LEVEL", test.logLevel)
			level := getLogLevel()
			assert.Equal(t, test.expectedLevel, level)
		})
	}
}
