package log

import (
	"context"
	"errors"
	"fmt"
	"github.com/caoyingjunz/pixiu/pkg/db"
	"sync"
	"time"

	klog "github.com/sirupsen/logrus"
)

var once sync.Once

type LogFormat string

const (
	LogFormatJson LogFormat = "json"
	LogFormatText LogFormat = "text"
)

var ErrInvalidLogFormat = errors.New("invalid log format")

type LogLevel = klog.Level

// Providing 3 log levels now.
const (
	ErrorLevel LogLevel = klog.ErrorLevel
	InfoLevel  LogLevel = klog.InfoLevel
	DebugLevel LogLevel = klog.DebugLevel
)

type LogOptions struct {
	LogFormat `yaml:"log_format"`
	LogSQL    bool `yaml:"log_sql"`
	LogLevel  `yaml:"log_level"`
}

// DefaultLogOptions returns the default configs.
func DefaultLogOptions() *LogOptions {
	return &LogOptions{
		LogFormat: LogFormatJson,
		LogSQL:    false,
		LogLevel:  InfoLevel,
	}
}

func (o *LogOptions) Valid() error {
	switch o.LogFormat {
	case LogFormatJson, LogFormatText:
		return nil
	default:
		return ErrInvalidLogFormat
	}
}

// Init sets the log format only once.
func (o *LogOptions) Init() {
	once.Do(func() {
		klog.SetLevel(o.LogLevel)
		switch o.LogFormat {
		case LogFormatJson:
			klog.SetFormatter(&klog.JSONFormatter{
				TimestampFormat: time.RFC3339Nano,
			})
		default:
			klog.SetFormatter(&klog.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: time.RFC3339Nano,
			})
		}
	})
}

const (
	SuccessMsg = "SUCCESS"
	ErrorMsg   = "ERROR"
	FailMsg    = "FAIL"
)

type Logger struct {
	startTime time.Time
	logSQL    bool
	logEntry  *klog.Entry
}

func NewLogger(cfg *LogOptions) *Logger {
	return &Logger{
		startTime: time.Now(),
		logSQL:    cfg.LogSQL,
		logEntry:  klog.NewEntry(klog.StandardLogger()),
	}
}

func (l *Logger) WithLogField(key string, value interface{}) {
	l.logEntry = l.logEntry.WithField(key, value)
}

func (l *Logger) WithLogFields(fields map[string]interface{}) {
	l.logEntry = l.logEntry.WithFields(fields)
}

func (l *Logger) Log(ctx context.Context, level LogLevel, err error) {
	fields := make(map[string]interface{})
	if l.logSQL {
		if sqls := db.GetSQLs(ctx); len(sqls) > 0 {
			fields["sqls"] = sqls
		}
	}
	fields["latency"] = fmt.Sprintf("%dµs", time.Since(l.startTime).Microseconds())

	if err != nil {
		fields["error"] = err
		l.logEntry.WithFields(fields).Error(FailMsg)
		return
	}

	switch level {
	case DebugLevel:
		l.logEntry.WithFields(fields).Debug(SuccessMsg)
	case InfoLevel:
		l.logEntry.WithFields(fields).Info(SuccessMsg)
	}
}
