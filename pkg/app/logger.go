package app

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func InitializeLogger() (*zap.SugaredLogger, error) {
	level, err := zap.ParseAtomicLevel("debug")
	if err != nil {
		return nil, err
	}
	cfg := &zap.Config{
		Encoding:         "json",
		ErrorOutputPaths: []string{"stderr"},
		OutputPaths:      []string{"stdout"},
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

	// initialize the rotator
	logFile := "logs/app-%Y-%m-%d-%H.log"
	rotator, err := rotatelogs.New(
		logFile,
		rotatelogs.WithMaxAge(60*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
	}

	writer := zapcore.AddSync(rotator)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		writer,
		level,
	)

	logger := zap.New(core).Sugar()
	logger.Infow("Logger setup done.", "level", level.String())
	return logger, nil
}
