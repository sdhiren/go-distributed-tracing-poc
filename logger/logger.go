package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
    logger *zap.Logger
}

func NewLogger() *Logger {
    cfg := zap.NewProductionConfig()
    cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    logger, _ := cfg.Build()
    defer logger.Sync()
    return &Logger{logger: logger}
}

func (l *Logger) WithContext(context context.Context) *zap.Logger {
    if context == nil {
		return nil
	}
	span := trace.SpanFromContext(context)
	spanId := span.SpanContext().SpanID().String()
	traceId := span.SpanContext().TraceID().String()
    return l.logger.With(zap.String("traceID", traceId), zap.String("spanID", spanId))
}


func (l *Logger) Debug(msg string, fields ...zap.Field) {
    l.logger.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
    l.logger.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
    l.logger.Error(msg, fields...)
}

func (l *Logger) Debugf(msg string, fields ...zap.Field) {
    l.logger.Sugar().Debugf(msg, fields)
}

func (l *Logger) Infof(msg string, fields ...zap.Field) {
    l.logger.Sugar().Infof(msg, fields)
}

func (l *Logger) Errorf(msg string, fields ...zap.Field) {
    l.logger.Sugar().Errorf(msg, fields)
}
