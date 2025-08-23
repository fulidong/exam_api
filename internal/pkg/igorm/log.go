package igorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"strings"
	"time"
)

func NewLogger(l log.Logger, level gormlogger.LogLevel) gormlogger.Interface {
	return &logger{
		log:   log.NewHelper(l),
		level: level,
	}
}

type logger struct {
	log           *log.Helper
	level         gormlogger.LogLevel
	SlowThreshold time.Duration
	Name          string
}

func (l *logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	l.level = level
	return l
}

func (l *logger) Info(ctx context.Context, s string, i ...interface{}) {
	if l.level < gormlogger.Info {
		return
	}

	ll := log.Context(ctx)
	ll.Infof(s, i...)
}

func (l *logger) Warn(ctx context.Context, s string, i ...interface{}) {
	if l.level < gormlogger.Warn {
		return
	}

	ll := log.Context(ctx)
	ll.Warnf(s, i...)
}

func (l *logger) Error(ctx context.Context, s string, i ...interface{}) {
	if l.level < gormlogger.Error {
		return
	}

	ll := log.Context(ctx)
	ll.Errorf(s, i...)
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level <= 0 {
		return
	}

	var (
		ll        = log.Context(ctx)
		elapsed   = time.Since(begin)
		sql, rows = fc()
	)

	switch {
	case err != nil && l.level >= gormlogger.Error:
		if errors.Is(err, gorm.ErrRecordNotFound) && strings.HasPrefix(sql, "SELECT") {
			ll.Infow("duration", fmt.Sprintf("%v", elapsed),
				"rows", rows,
				"sql", sql)
			return
		}
		ll.Errorw("duration", fmt.Sprintf("%v", elapsed),
			"rows", rows,
			"sql", sql)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.level >= gormlogger.Warn:
		ll.Warnw("duration", fmt.Sprintf("%v", elapsed),
			"rows", rows,
			"sql", sql)
	case l.level >= gormlogger.Info:
		ll.Infow("duration", fmt.Sprintf("%v", elapsed),
			"rows", rows,
			"sql", sql)
	}
}
