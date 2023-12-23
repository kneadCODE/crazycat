package internal

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZap(debugMode bool, res *resource.Resource) (*zap.Logger, error) {
	var l *zap.Logger
	var err error

	if debugMode {
		l, err = zap.Config{
			Level: zap.NewAtomicLevelAt(zapcore.DebugLevel),
			// Development: true,
			Encoding: "console",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "T",
				LevelKey:       "S",
				NameKey:        zapcore.OmitKey,
				CallerKey:      zapcore.OmitKey,
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "B",
				StacktraceKey:  zapcore.OmitKey,
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalColorLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}.Build()
	} else {
		l, err = zap.Config{
			Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
			Development: false,
			Encoding:    "json",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "Timestamp",     // OTEL compliant
				LevelKey:       "Severity",      // OTEL compliant
				NameKey:        zapcore.OmitKey, // OTEL compliant
				CallerKey:      zapcore.OmitKey, // OTEL compliant
				FunctionKey:    zapcore.OmitKey, // OTEL compliant
				MessageKey:     "Body",          // OTEL compliant
				StacktraceKey:  zapcore.OmitKey, // OTEL compliant
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalLevelEncoder, // OTEL compliant
				EncodeTime:     zapcore.EpochTimeEncoder,    // OTEL compliant
				EncodeDuration: zapcore.MillisDurationEncoder,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}.Build(zap.Fields(zap.Object("Resource", resourceZapWrapper{res})))
	}
	if err != nil {
		return nil, fmt.Errorf("golib:app:NewZap err initializing zap: %w", err)
	}

	return l, nil
}

type attributesZapWrapper []attribute.KeyValue

func (attr attributesZapWrapper) MarshalLogObject(z zapcore.ObjectEncoder) error {
	for _, a := range attr {
		switch a.Value.Type() {
		case attribute.INT64:
			z.AddInt64(string(a.Key), a.Value.AsInt64())
		case attribute.FLOAT64:
			z.AddFloat64(string(a.Key), a.Value.AsFloat64())
		case attribute.STRING:
			z.AddString(string(a.Key), a.Value.AsString())
		case attribute.BOOL:
			z.AddString(string(a.Key), a.Value.AsString())
		case attribute.STRINGSLICE:
			z.AddString(string(a.Key), strings.Join(a.Value.AsStringSlice(), ","))
		default:
			// Intentionally skipping it as we are not expecting any other type
		}
	}
	return nil
}

type resourceZapWrapper struct {
	*resource.Resource
}

func (r resourceZapWrapper) MarshalLogObject(z zapcore.ObjectEncoder) error {
	return attributesZapWrapper(r.Attributes()).MarshalLogObject(z)
}

func ZapLogEnriched(z *zap.Logger, level zapcore.Level, msg string, span trace.Span, attrs ...attribute.KeyValue) {
	// Ref: OTEL logging fields https://opentelemetry.io/docs/specs/otel/logs/data-model/#field-resource
	z.Log(
		level,
		msg,
		zap.String("TraceId", span.SpanContext().TraceID().String()),
		zap.String("SpanId", span.SpanContext().SpanID().String()),
		zap.String("TraceFlags", span.SpanContext().TraceFlags().String()),
		zap.Object("Attributes", attributesZapWrapper(attrs)),
		// We are unable to get span's attrs and merge it with the given attrs because the span attrs are not exposed.
		// Hence, wherever the span is created or span.SetAttribute is called, we need to set it in zap itself.
	)
}
