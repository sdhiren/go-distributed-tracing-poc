package logging

import (
	"fmt"
	"os"

	"golang.org/x/exp/slog"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type CustomLogger struct {
	logger *slog.Logger
	context *gin.Context
}

func NewCustomLogger(logger *slog.Logger, context *gin.Context) *CustomLogger {
	return &CustomLogger{
		logger: logger,
		context: context,
	}
}

func (c *CustomLogger) Info(msg string, args ...any) {
	c.logger.Info(msg, args...)
}

func (c *CustomLogger) Error(msg string, args ...any) {
	c.logger.Error(msg, args...)
}

func (c *CustomLogger) With(key string, value string) *CustomLogger {
	logger := c.logger
	logger = logger.With(slog.String(key, value))
	c.logger = logger
	return c
}

func GetDefaultLogger(context *gin.Context) *CustomLogger {
	if context == nil {
		return nil
	}

	span := trace.SpanFromContext(context.Request.Context())
	spanId := span.SpanContext().SpanID().String()
	traceId := span.SpanContext().TraceID().String()

	end_point := context.Request.URL.Path
	logLevel := getLogLevel()

	handlerOpts := &slog.HandlerOptions{		
		Level:     logLevel,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, handlerOpts))
	slog.SetDefault(logger)
	
	defaultLogger := logger.With(
								slog.String("trace_id", traceId),
								slog.String("span_id", spanId),
								slog.String("end_point", end_point),
							  )

	return NewCustomLogger(defaultLogger, context)
}

func getLogLevel() slog.Level {
	log_level := os.Getenv("LOG_LEVEL")
	var logLevel slog.Level
	switch log_level {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	}
	return logLevel
}

func GetFileLogger(context *gin.Context) *CustomLogger  {
	span := trace.SpanFromContext(context.Request.Context())
	spanId := span.SpanContext().SpanID().String()
	traceId := span.SpanContext().TraceID().String()

	LOG_FILE := os.Getenv("LOG_FILE")
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error :", err)
	}
	// defer file.Close()

	end_point := context.Request.URL.Path
	logLevel := getLogLevel()

	handlerOpts := &slog.HandlerOptions{		
		Level:     logLevel,
	}
	logger := slog.New(slog.NewJSONHandler(file, handlerOpts))
	slog.SetDefault(logger)

	
	defaultLogger := logger.With(
								 slog.String("trace_id", string(traceId)),
								 slog.String("span_id", spanId),
								 slog.String("end_point", end_point),
								)

	customLogger := &CustomLogger{}
	customLogger.context = context
	customLogger.logger = defaultLogger

	return customLogger
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

