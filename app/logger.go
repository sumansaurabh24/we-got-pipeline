package app

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitializeLogger() (*zap.SugaredLogger, error) {
	level, err := zap.ParseAtomicLevel("debug")
	if err != nil {
		return nil, err
	}
	cfg := &zap.Config{
		Encoding:         "json",
		ErrorOutputPaths: []string{"stderr"},
		OutputPaths:      []string{"stdout", "./logs/app.log"},
		Level:            level,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "timestamp",
			NameKey:        "name",
			CallerKey:      "caller",
			FunctionKey:    "function",
			StacktraceKey:  "stack_trace",
			SkipLineEnding: false,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger := zap.Must(cfg.Build()).Sugar()
	defer logger.Sync()

	logger.Infow("Logger setup done.", "level", level.String())
	return logger, nil
}
