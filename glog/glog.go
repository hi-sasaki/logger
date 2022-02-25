package glog

import (
	"context"

	"logger"

	"github.com/golang/glog"
)

type glogLogger struct{}

// NewLogger glogインスタンスを生成する
func NewLogger() logger.Interface {
	return new(glogLogger)
}

func (g *glogLogger) Debug(_ context.Context, e logger.Entry) {
	glog.Infof("(%s) [%s]: %s\n%s", e.TraceID, e.Caller, e.Message(), e.Trace)
}

func (g *glogLogger) Info(_ context.Context, e logger.Entry) {
	glog.Infof("(%s) [%s]: %s\n%s", e.TraceID, e.Caller, e.Message(), e.Trace)
}

func (g *glogLogger) Warn(_ context.Context, e logger.Entry) {
	glog.Warningf("(%s) [%s]: %s\n%s", e.TraceID, e.Caller, e.Message(), e.Trace)
}

func (g *glogLogger) Error(_ context.Context, e logger.Entry) {
	glog.Errorf("(%s) [%s]: %s\n%s", e.TraceID, e.Caller, e.Message(), e.Trace)
}

func (g *glogLogger) Critical(_ context.Context, e logger.Entry) {
	glog.Errorf("(%s) [%s]: %s\n%s", e.TraceID, e.Caller, e.Message(), e.Trace)
}

func (g *glogLogger) Finalize() error {
	return nil
}
