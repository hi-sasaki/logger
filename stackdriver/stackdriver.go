package stackdriver

import (
	"context"
	"fmt"
	"log"

	"logger"

	"cloud.google.com/go/logging"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

// Option ロガーの情報
type Option struct {
	// Required プロジェクトID
	// AppEngine環境では `os.GetEnv("GOOGLE_CLOUD_PROJECT")` で取得可能
	Project string

	// Required AppEngineサービス名
	// AppEngine環境では `os.GetEnv("GAE_SERVICE")` で取得可能
	Service string

	// Required AppEngineバージョン
	// AppEngine環境では `os.GetEnv("GAE_VERSION")` で取得可能
	Version string

	// Required ログタイプ
	// AppEngine用ログ: gae_app
	Type string

	// Required ログID StackDriverでのログ絞り込みで利用する
	// e.g. app_logs
	Name string
}

type stackdriverLogger struct {
	opt    Option
	logger *logging.Logger
	client *logging.Client
}

// NewLogger Stackdriver loggingにログ出力を行うロガーを生成する
func NewLogger(opt Option) (logger.Interface, error) {
	client, err := logging.NewClient(context.Background(), fmt.Sprintf("projects/%s", opt.Project))
	if err != nil {
		return nil, err
	}

	return &stackdriverLogger{
		opt:    opt,
		logger: client.Logger(opt.Name),
		client: client,
	}, nil
}

func (l *stackdriverLogger) Debug(ctx context.Context, e logger.Entry) {
	l.log(ctx, logging.Debug, e)
}

func (l *stackdriverLogger) Info(ctx context.Context, e logger.Entry) {
	l.log(ctx, logging.Info, e)
}

func (l *stackdriverLogger) Warn(ctx context.Context, e logger.Entry) {
	l.log(ctx, logging.Warning, e)
}

func (l *stackdriverLogger) Error(ctx context.Context, e logger.Entry) {
	l.log(ctx, logging.Error, e)
}

func (l *stackdriverLogger) Critical(ctx context.Context, e logger.Entry) {
	l.log(ctx, logging.Critical, e)
}

func (l stackdriverLogger) Finalize() error {
	go func() {
		if err := l.logger.Flush(); err != nil {
			log.Printf("failed to flush log: %#v", err)
		}
		if err := l.client.Close(); err != nil {
			log.Printf("failed to close client: %#v", err)
		}
	}()
	return nil
}

func (l stackdriverLogger) log(ctx context.Context, severity logging.Severity, e logger.Entry) {
	l.logger.Log(logging.Entry{
		Severity: severity,
		Trace:    fmt.Sprintf("projects/%s/traces/%s", l.opt.Project, e.TraceID), // request_logsと同じフォーマットにして紐付けができるようにする
		Payload: map[string]interface{}{
			"message": e.Message(),
			"caller":  e.Caller,
			"trace":   e.Trace,
		},
		Resource: &monitoredres.MonitoredResource{
			Labels: map[string]string{
				"module_id":  l.opt.Service,
				"project_id": l.opt.Project,
				"version_id": l.opt.Version,
			},
			Type: l.opt.Type,
		},
	})
}
