package app

import (
	"fmt"

	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newZap(cfg Config) (*zap.SugaredLogger, error) {
	var zapCfg zap.Config

	switch cfg.Env {
	case EnvDev:
		zapCfg = zap.Config{
			Level: zap.NewAtomicLevelAt(zapcore.DebugLevel),
			// Development: true,
			Encoding: "console",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:  "T",
				LevelKey: "S",
				NameKey:  "N",
				// CallerKey:      "C",
				FunctionKey: zapcore.OmitKey,
				MessageKey:  "B",
				// StacktraceKey:  "S",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalColorLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	default:
		zapCfg = zap.Config{
			Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
			Development: false,
			Encoding:    "json",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:  "ts",
				LevelKey: "severity",
				NameKey:  "logger",
				// CallerKey:      "caller",
				FunctionKey: zapcore.OmitKey,
				MessageKey:  "body",
				// StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.EpochTimeEncoder,
				EncodeDuration: zapcore.MillisDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}

	l, err := zapCfg.Build(
		zap.Fields(
			zap.String(string(semconv.ServiceNameKey), cfg.Name),
			zap.String(string(semconv.ServiceNamespaceKey), cfg.Project),
			zap.String(string(semconv.ServiceVersionKey), cfg.Version),
			zap.String(string(semconv.ServiceInstanceIDKey), cfg.ServerInstanceID),
			zap.String(string(semconv.DeploymentEnvironmentKey), string(cfg.Env)),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("err initializing zap: %w", err)
	}

	return l.Sugar(), nil
}
