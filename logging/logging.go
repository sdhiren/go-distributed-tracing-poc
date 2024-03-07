package logging

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/exp/slog"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type CustomLogger struct {
	logger *slog.Logger
	context *gin.Context
}

func (c *CustomLogger) Info(msg string, args ...any) {
	span := trace.SpanFromContext(c.context.Request.Context())
	spanId := span.SpanContext().SpanID().String()
	traceId := span.SpanContext().TraceID().String()
	var attributes []attribute.KeyValue

	attributes = append(attributes, attribute.String("traceId", string(traceId)), attribute.String("spanId", spanId), attribute.String("time", time.Now().String()), attribute.String("level", slog.LevelInfo.String()))
	span.AddEvent(msg, trace.WithAttributes(attributes...), trace.WithTimestamp(time.Now()))
	// c.logger.InfoContext(c.context.Request.Context(), msg, args...)
	c.logger.Info(msg, args...)
}

func (c *CustomLogger) Error(msg string, args ...any) {
	span := trace.SpanFromContext(c.context.Request.Context())
	spanId := span.SpanContext().SpanID().String()
	traceId := span.SpanContext().TraceID().String()
	var attributes []attribute.KeyValue

	attributes = append(attributes, attribute.String("traceId", string(traceId)), attribute.String("spanId", spanId), attribute.String("time", time.Now().String()), attribute.String("level", slog.LevelInfo.String()))
	span.AddEvent(msg, trace.WithAttributes(attributes...), trace.WithTimestamp(time.Now()))
	// c.logger.InfoContext(c.context.Request.Context(), msg, args...)
	c.logger.Error(msg, args...)
}


func GetLogger(cotext *gin.Context) *slog.Logger {
	LOG_FILE := os.Getenv("LOG_FILE")
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error :", err)
	}
	defer file.Close()

	handlerOpts := &slog.HandlerOptions{		
		Level:     slog.LevelDebug,		
	}
	logger := slog.New(slog.NewJSONHandler(file, handlerOpts))
	slog.SetDefault(logger)

	return logger
}

func GetDefaultLogger(context *gin.Context) *CustomLogger {
	span := trace.SpanFromContext(context.Request.Context())
	spanId := span.SpanContext().SpanID().String()
	traceId := span.SpanContext().TraceID().String()

	fmt.Println("traceId: ", traceId)

	handlerOpts := &slog.HandlerOptions{		
		Level:     slog.LevelDebug,		
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, handlerOpts))
	slog.SetDefault(logger)

	
	defaultLogger := logger.With(slog.String("trace_id", string(traceId)),
	slog.String("span_id", spanId),
	slog.String("method_name", "methodName"),
	slog.String("class_name", "className"),)

	customLogger := &CustomLogger{}
	customLogger.context = context
	customLogger.logger = defaultLogger

	return customLogger
}

func GetFileLogger(context *gin.Context) *CustomLogger  {
	span := trace.SpanFromContext(context.Request.Context())
	spanId := span.SpanContext().SpanID().String()
	traceId := span.SpanContext().TraceID().String()

	LOG_FILE := os.Getenv("LOG_FILE")
	fmt.Println("Log file: ",LOG_FILE)
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error :", err)
	}
	// defer file.Close()

	handlerOpts := &slog.HandlerOptions{		
		Level:     slog.LevelDebug,		
	}
	logger := slog.New(slog.NewJSONHandler(file, handlerOpts))
	slog.SetDefault(logger)

	
	defaultLogger := logger.With(slog.String("trace_id", string(traceId)),
	slog.String("span_id", spanId),
	slog.String("method_name", "methodName"),
	slog.String("class_name", "className"),)

	customLogger := &CustomLogger{}
	customLogger.context = context
	customLogger.logger = defaultLogger

	return customLogger
}