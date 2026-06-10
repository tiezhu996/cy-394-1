package utils

import "go.uber.org/zap"

func NewLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}
func ZapError(err error) zap.Field          { return zap.Error(err) }
func ZapString(key, value string) zap.Field { return zap.String(key, value) }
