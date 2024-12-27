/*
Copyright 2024 The Vuples Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package log

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"kubevulpes/pkg/db"
)

var once sync.Once

const (
	LogFormatJson = "json"
	LogFormatText = "text"
)

var ErrInvalidLogFormat = errors.New("invalid log format")

type LogLevel = log.Level

// Providing 3 log levels now.
const (
	ErrorLevel LogLevel = log.ErrorLevel
	InfoLevel  LogLevel = log.InfoLevel
	DebugLevel LogLevel = log.DebugLevel
)

type LogOptions struct {
	LogFormat string `config:"log_format"`
	LogSQL    bool   `config:"log_sql"`
	LogLevel  `config:"log_level"`
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
		log.SetLevel(o.LogLevel)
		switch o.LogFormat {
		case LogFormatJson:
			log.SetFormatter(&log.JSONFormatter{
				TimestampFormat: time.RFC3339Nano,
			})
		default:
			log.SetFormatter(&log.TextFormatter{
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
	logEntry  *log.Entry
}

func NewLogger(cfg *LogOptions) *Logger {
	return &Logger{
		startTime: time.Now(),
		logSQL:    cfg.LogSQL,
		logEntry:  log.NewEntry(log.StandardLogger()),
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
