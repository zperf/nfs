package logging

import (
	"github.com/willscott/go-nfs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level.SetLevel(zap.InfoLevel)
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
	return logger
}

type MyLogger struct {
	Z *zap.Logger
}

func (m *MyLogger) SetLevel(level nfs.LogLevel) {
	m.Z.Sugar().Panicf("Not supported")
}

func (m *MyLogger) GetLevel() nfs.LogLevel {
	switch m.Z.Level() {
	case zap.DebugLevel:
		return nfs.DebugLevel
	case zap.InfoLevel:
		return nfs.InfoLevel
	case zap.WarnLevel:
		return nfs.WarnLevel
	case zap.ErrorLevel:
		return nfs.ErrorLevel
	case zap.FatalLevel:
		return nfs.FatalLevel
	case zap.DPanicLevel:
		return nfs.PanicLevel
	case zap.PanicLevel:
		return nfs.PanicLevel
	default:
		return nfs.InfoLevel
	}
}

func (m *MyLogger) ParseLevel(level string) (nfs.LogLevel, error) {
	return (*nfs.DefaultLogger)(nil).ParseLevel(level)
}

func (m *MyLogger) Panic(args ...interface{}) {
	m.Z.Sugar().Panic(args)
}

func (m *MyLogger) Fatal(args ...interface{}) {
	m.Z.Sugar().Fatal(args)
}

func (m *MyLogger) Error(args ...interface{}) {
	m.Z.Sugar().Error(args)
}

func (m *MyLogger) Warn(args ...interface{}) {
	m.Z.Sugar().Warn(args)
}

func (m *MyLogger) Info(args ...interface{}) {
	m.Z.Sugar().Info(args)
}

func (m *MyLogger) Debug(args ...interface{}) {
	m.Z.Sugar().Debug(args)
}

func (m *MyLogger) Trace(args ...interface{}) {
	m.Z.Sugar().Debug(args)
}

func (m *MyLogger) Print(args ...interface{}) {
	m.Z.Sugar().Info(args)
}

func (m *MyLogger) Panicf(format string, args ...interface{}) {
	m.Z.Sugar().Panicf(format, args)
}

func (m *MyLogger) Fatalf(format string, args ...interface{}) {
	m.Z.Sugar().Panicf(format, args)
}

func (m *MyLogger) Errorf(format string, args ...interface{}) {
	m.Z.Sugar().Errorf(format, args)
}

func (m *MyLogger) Warnf(format string, args ...interface{}) {
	m.Z.Sugar().Warnf(format, args)
}

func (m *MyLogger) Infof(format string, args ...interface{}) {
	m.Z.Sugar().Infof(format, args)
}

func (m *MyLogger) Debugf(format string, args ...interface{}) {
	m.Z.Sugar().Debugf(format, args)
}

func (m *MyLogger) Tracef(format string, args ...interface{}) {
	m.Z.Sugar().Debugf(format, args)
}

func (m *MyLogger) Printf(format string, args ...interface{}) {
	m.Z.Sugar().Infof(format, args)
}
