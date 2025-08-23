package ilog

import (
	"exam_api/internal/pkg/icontext"
	"fmt"
	"io"
	"net"
	"os"
	"path"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	k8sJsonLogPath = "/var/jsonlog/%s"
)

type logger struct {
	*zap.SugaredLogger
}

type closer struct {
	sugar *zap.SugaredLogger
}

func (c *closer) Close() error {
	return c.sugar.Sync()
}

func NewLogger(id, name string, opts ...Option) (log.Logger, io.Closer) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	var (
		logPath   = os.Getenv("LOG_PATH")
		_, k8sEnv = os.LookupEnv("KUBERNETES_SERVICE_HOST")
		writer    io.Writer
	)

	if logPath == "" && k8sEnv {
		logPath = fmt.Sprintf(k8sJsonLogPath, name)
	}

	if logPath == "" || o.console {
		writer = os.Stdout
	} else {
		writer = &lumberjack.Logger{
			Filename:  path.Join(logPath, fmt.Sprintf("%s-%s.log", name, getLocalIP())),
			MaxSize:   100,
			MaxAge:    7,
			LocalTime: true,
		}
	}

	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(conf)

	var (
		core         = zapcore.NewCore(encoder, zapcore.AddSync(writer), zapcore.DebugLevel)
		globalZapLog = zap.New(
			core,
			zap.AddStacktrace(zap.ErrorLevel),
		)
	)

	var (
		globalSugarLogger = globalZapLog.Sugar()
		kvs               = []interface{}{
			"service_id", id,
			"service_name", name,
			"trace_id", tracing.TraceID(),
			"span_id", tracing.SpanID(),
		}
	)
	kvs = append(kvs, icontext.LoggerValues()...)

	ll := log.With(&logger{
		SugaredLogger: globalSugarLogger,
	}, kvs...)

	return ll, &closer{
		sugar: globalSugarLogger,
	}
}

func (l *logger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		return nil
	}

	var (
		msg string
		kvs []interface{}
	)

	for i := 0; i < len(keyvals); i += 2 {
		var (
			key   = fmt.Sprint(keyvals[i])
			value = fmt.Sprint(keyvals[i+1])
		)

		if value == "" {
			continue
		}

		if key == log.DefaultMessageKey {
			msg = value
			continue
		}

		kvs = append(kvs, zap.Any(key, value))
	}

	switch level {
	case log.LevelDebug:
		l.Debugw(msg, kvs...)
	case log.LevelInfo:
		l.Infow(msg, kvs...)
	case log.LevelWarn:
		l.Warnw(msg, kvs...)
	case log.LevelError:
		l.Errorw(msg, kvs...)
	case log.LevelFatal:
		l.Fatalw(msg, kvs...)
	}
	return nil
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return uuid.New().String()
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return uuid.New().String()
}
