package logger

import (
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "ts",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

var zapServerConfig = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
	Development:      false,
	Encoding:         "json",
	EncoderConfig:    defaultEncoderConfig,
	OutputPaths:      []string{"stdout"},
	ErrorOutputPaths: []string{"stderr"},
}

// NewRootLogger instantiates zap.Logger with given configuration
func NewRootLogger(debug, developement bool) (*zap.Logger, error) {
	var err error
	var logger *zap.Logger

	if debug {
		// enable debug level
		zapServerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		if !developement {
			return zapServerConfig.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
		}
		zapServerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapServerConfig.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
	}

	if developement {
		zapServerConfig.Encoding = "console"
		zapServerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapServerConfig.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
	}

	log.SetFormatter(&log.JSONFormatter{})
	logger, err = zapServerConfig.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())

	if err != nil {
		panic(err)
	}

	return logger, err
}
