package app

import (
	"context"
)

func LogInfo(ctx context.Context) {

}

//
// func newZap(env Environment) error {
// 	var logger *zap.Logger
// 	var err error
//
// 	switch env {
// 	case EnvDev:
// 		logger, err = zap.NewDevelopment(
// 			zap.AddStacktrace(zapcore.ErrorLevel),
// 			zap.WrapCore(func(core zapcore.Core) zapcore.Core {
// 				zapcore.NewTee()
// 			}),
// 		)
// 	default:
// 		logger, err = zap.NewProduction()
// 	}
//
// 	if err != nil {
// 		return fmt.Errorf("err init zap: %w", err)
// 	}
//
// 	logger.Sugar()
//
// 	return nil
// }
