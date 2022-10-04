package logging

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var rootLogger *zap.Logger
var logger *zap.SugaredLogger
var isConfigured bool
var defaultOutputPaths []string
var logLevel zapcore.Level

func init() {
	defaultOutputPaths = []string{"stdout"}
	configure("debug", "root", defaultOutputPaths, defaultOutputPaths, "console")
}

func configure(level, name string, outputPaths, errOutputPaths []string, encoder string) {
	l := zap.NewAtomicLevel()
	err := l.UnmarshalText([]byte(level))
	if err != nil {
		panic(fmt.Sprintf("Incorrect logging level: `%s`.", level))
	}
	logLevel = l.Level()
	// If no path config is passed, we will default to `stdout`.  Otherwise, the logs
	// get swallowed...
	if len(outputPaths) == 0 {
		outputPaths = defaultOutputPaths
	}
	if len(errOutputPaths) == 0 {
		errOutputPaths = defaultOutputPaths
	}
	conf := zap.Config{
		Level:            l,
		Development:      (logLevel < zapcore.InfoLevel),
		Encoding:         encoder,
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      outputPaths,
		ErrorOutputPaths: errOutputPaths,
	}
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	newLogger, _ := conf.Build()
	rootLogger = newLogger.Named(name)
	logger = rootLogger.Sugar()
}

func Configure(level, name string, outputPaths, errOutputPaths []string, encoder ...string) {
	enc := "console"
	if encoder != nil {
		enc = encoder[0]
	}
	configure(level, name, outputPaths, errOutputPaths, enc)
	isConfigured = true
}

func GetRootLogger() *zap.Logger {
	if !isConfigured {
		panic("Can not access root logger before it is configured.")
	}
	return GetRootLoggerUnsafe()
}

func GetLogger() *zap.SugaredLogger {
	if !isConfigured {
		panic("Can not access logger before it is configured.")
	}
	return GetLoggerUnsafe()
}

func GetRootLoggerUnsafe() *zap.Logger {
	return rootLogger
}

func GetLoggerUnsafe() *zap.SugaredLogger {
	return logger
}

func GetContextLogger(c context.Context) (logger *zap.SugaredLogger) {
	if val := c.Value("logger"); val != nil {
		logger, _ = val.(*zap.SugaredLogger)
	}
	return
}

func LevelValue() zapcore.Level {
	return logLevel
}
