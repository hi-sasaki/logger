package logger

import (
	"context"
	"fmt"
	"runtime"
	"strings"
)

// Entry ログ出力情報
type Entry struct {
	Format  string
	Args    []interface{}
	Caller  string
	Trace   string
	TraceID string
}

func createEntry(traceID, format string, args ...interface{}) Entry {
	pc, file, line, _ := runtime.Caller(2)
	var (
		f        = runtime.FuncForPC(pc)
		parts    = strings.Split(f.Name(), "/")
		funcName = parts[len(parts)-1]
	)
	stack := make([]byte, 2048)
	length := runtime.Stack(stack, false)
	return Entry{
		Format:  format,
		Args:    args,
		Caller:  fmt.Sprintf("%s#%s:L%v", file, funcName, line),
		Trace:   string(stack[:length]),
		TraceID: traceID,
	}
}

func createNoTraceEntry(traceID, format string, args ...interface{}) Entry {
	pc, file, line, _ := runtime.Caller(2)
	var (
		f        = runtime.FuncForPC(pc)
		parts    = strings.Split(f.Name(), "/")
		funcName = parts[len(parts)-1]
	)
	return Entry{
		Format:  format,
		Args:    args,
		Caller:  fmt.Sprintf("%s#%s:L%v", file, funcName, line),
		TraceID: traceID,
	}
}

// Message ログメッセージ取得
func (a Entry) Message() string {
	return fmt.Sprintf(a.Format, a.Args...)
}

// TraceIDGen ログ追跡用のトレースID生成関数定義
type TraceIDGen func(context.Context) string

// Logger ログ出力のファサード
// ログ情報として出力元情報とトレースIDを出力可能
type Logger struct {
	logger         Interface
	extractTraceID TraceIDGen
}

// NewLogger Loggerを生成する
func NewLogger(logger Interface, f TraceIDGen) *Logger {
	return &Logger{
		logger:         logger,
		extractTraceID: f,
	}
}

// Debugf デバッグログを出力をする
func (l Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Debug(ctx, createEntry(l.traceID(ctx), format, args...))
}

func (l Logger) DebugfNoTrace(ctx context.Context, format string, args ...interface{}) {
	l.logger.Debug(ctx, createNoTraceEntry(l.traceID(ctx), format, args...))
}

// Infof 情報ログを出力をする
func (l Logger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.logger.Info(ctx, createEntry(l.traceID(ctx), format, args...))
}

// Warnf 警告ログを出力をする
func (l Logger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Warn(ctx, createEntry(l.traceID(ctx), format, args...))
}

// Errorf エラーログを出力をする
func (l Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.logger.Error(ctx, createEntry(l.traceID(ctx), format, args...))
}

// Criticalf エラーログを出力をする
func (l Logger) Critical(ctx context.Context, summary string, err error) {
	l.logger.Critical(ctx, createEntry(l.traceID(ctx), "%s\n%+v", summary, err))
}

// Finalize Logger終了処理
func (l Logger) Finalize() error {
	return l.logger.Finalize()
}

func (l Logger) traceID(ctx context.Context) string {
	var traceID string
	if l.extractTraceID != nil {
		traceID = l.extractTraceID(ctx)
	}
	return traceID
}

// Interface ログの出力を行うインターフェースの定義
type Interface interface {
	Debug(context.Context, Entry)
	Info(context.Context, Entry)
	Warn(context.Context, Entry)
	Error(context.Context, Entry)
	Critical(context.Context, Entry)
	Finalize() error
}
