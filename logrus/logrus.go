package logrus

import (
	"context"

	"logger"

	"github.com/sirupsen/logrus"
)

// Options Loggerオプション構造体
type Options struct {
	formatter logrus.Formatter
	level     logrus.Level
}

// Option オプション設定関数
type Option func(*Options)

// WithFormatter ログフォーマットオプション設定
func WithFormatter(formatter logrus.Formatter) Option {
	return func(ops *Options) {
		ops.formatter = formatter
	}
}

// WithLevel ログレベルオプション設定
func WithLevel(level logrus.Level) Option {
	return func(ops *Options) {
		ops.level = level
	}
}

type logrusLogger struct {
	log *logrus.Logger
}

// NewLogger インスタンス生成
func NewLogger(options ...Option) logger.Interface {
	log := logrus.New()
	opt := Options{
		formatter: new(logrus.TextFormatter),
		level:     logrus.DebugLevel,
	}
	for _, o := range options {
		o(&opt)
	}
	log.Formatter = opt.formatter
	log.Level = opt.level
	return &logrusLogger{log}
}

// Debug Debugレベルのログを出力する関数
func (l logrusLogger) Debug(ctx context.Context, e logger.Entry) {
	l.log.Debugf("[%s]: %s\n%s", e.Caller, e.Message(), e.Trace)
}

// Info Infoレベルのログを出力する関数
func (l logrusLogger) Info(ctx context.Context, e logger.Entry) {
	l.log.Infof("[%s]: %s\n%s", e.Caller, e.Message(), e.Trace)
}

// Warn Warnレベルのログを出力する関数
func (l logrusLogger) Warn(ctx context.Context, e logger.Entry) {
	l.log.Warnf("[%s]: %s\n%s", e.Caller, e.Message(), e.Trace)
}

// Error Errorレベルのログを出力する関数
func (l logrusLogger) Error(ctx context.Context, e logger.Entry) {
	l.log.Errorf("[%s]: %s\n%s", e.Caller, e.Message(), e.Trace)
}

// Critical Criticalレベルのログを出力する関数
func (l logrusLogger) Critical(ctx context.Context, e logger.Entry) {
	l.log.Errorf("[%s]: %s\n%s", e.Caller, e.Message(), e.Trace)
}

// Finalize implements logger.Interface
func (logrusLogger) Finalize() error {
	return nil
}
