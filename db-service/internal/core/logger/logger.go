package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(logLevel string, path string) (*zap.Logger, func() error, error) {
	// set depth lever for logs
	lvl := zap.NewAtomicLevel()
	if err := lvl.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, nil, fmt.Errorf("unmarshal log: %w", err)
	}

	// create dir with permissions to read and execute folders (for users)
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	// make filepath
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(path, fmt.Sprintf("%s.log", timestamp))

	// create and open file (with permissions)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file: %w", err)
	}

	// create zap config
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15-04-05.000000")

	// human encoding format
	encoder := zapcore.NewConsoleEncoder(cfg)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), lvl),
		zapcore.NewCore(encoder, zapcore.AddSync(logFile), lvl),
	)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger, func() error {
		logger.Sync()
		return logFile.Close()
	}, err
}
