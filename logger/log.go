package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// custom logger encoder
func newSimpleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:    "T",
		LevelKey:   "L",
		MessageKey: "M",
		CallerKey:  "C",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeLevel: func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			switch level {
			case zapcore.InfoLevel:
				enc.AppendString("[INFO]")
			case zapcore.ErrorLevel:
				enc.AppendString("[ERROR]")
			case zapcore.WarnLevel:
				enc.AppendString("[WARN]")
			default:
				enc.AppendString("[" + level.String() + "]")
			}
		},
		EncodeCaller: zapcore.ShortCallerEncoder, // e.g. logger/logger.go:52
	})
}

// New logger instance
func NewLogger(args ...interface{}) (*Logger, error) {
	var s string
	var config LoggerConfig
	var hasS, hasConfig bool

	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			s = v
			hasS = true
		case LoggerConfig:
			config = v
			hasConfig = true
		default:
			fmt.Printf("未知の型: %T\n", v)
		}
	}

	switch {
	case hasS:
		return newLogger(s)
	case !hasS && hasConfig:
		return newLoggerWithLumberjack(config)
	case hasS && hasConfig:
		// configが指定されている場合は、configを使用してログファイルを設定
		return newLoggerWithLumberjack(config)
	}

	return newLogger()
}

// LoggerConfig defines the configuration for the logger with lumberjack rotation.
func newLoggerWithLumberjack(config LoggerConfig) (*Logger, error) {
	encoder := newSimpleEncoder()

	if config.Filename == "" {
		return nil, fmt.Errorf("filename must be specified in LoggerConfig")
	}

	if config.MaxSize <= 0 {
		config.MaxSize = 1 // default value: 1 MB
	}

	if config.MaxBackups <= 0 {
		config.MaxBackups = 1 // default value: 1 backup
	}

	// lumberjack ローテーション設定
	rotator := &lumberjack.Logger{
		Filename:   config.Filename,   // e.g. "logs/app.log"
		MaxSize:    config.MaxSize,    // Max size in MB
		MaxBackups: config.MaxBackups, // Max number of old log files to keep
		MaxAge:     config.MaxAge,     // Max age in days
		Compress:   config.Compress,   // compress old log files
	}

	// console output
	consoleSync := zapcore.Lock(os.Stdout)
	consoleCore := zapcore.NewCore(encoder, consoleSync, zapcore.InfoLevel)

	// file output with rotation
	fileSync := zapcore.AddSync(rotator)
	fileCore := zapcore.NewCore(encoder, fileSync, zapcore.InfoLevel)

	// combine cores with Tee
	core := zapcore.NewTee(consoleCore, fileCore)
	// initialize zap logger with caller info (skip 2 levels)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

	return &Logger{logger: zapLogger}, nil
}

// newLogger creates a Logger that logs to console and optionally to a specified file.
// If paths are provided, the first path is used as the log file path.
func newLogger(paths ...string) (*Logger, error) {

	encoder := newSimpleEncoder()

	consoleSync := zapcore.Lock(os.Stdout)

	var cores []zapcore.Core

	// If a log file path is provided, set up file logging
	if len(paths) > 0 && paths[0] != "" {
		path := paths[0]
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		fileSync := zapcore.AddSync(file)
		cores = append(cores, zapcore.NewCore(encoder, fileSync, zapcore.InfoLevel))
	}

	cores = append(cores, zapcore.NewCore(encoder, consoleSync, zapcore.InfoLevel))

	core := zapcore.NewTee(cores...)

	// initialize zap logger with caller info (skip 2 levels)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

	return &Logger{logger: zapLogger}, nil
}

// Info logs an info message.
func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

// Error logs an error message.
func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string) {
	l.logger.Debug(msg)
}
