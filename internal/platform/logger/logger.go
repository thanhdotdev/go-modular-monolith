package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"project-example/internal/platform/config"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(cfg config.LoggingConfig) (*zap.Logger, io.Closer, error) {
	core, closer, err := newCore(cfg)
	if err != nil {
		return nil, nil, err
	}

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Fields(zap.String("service", cfg.ServiceName)),
	)
	zap.ReplaceGlobals(logger)

	return logger, closer, nil
}

func newCore(cfg config.LoggingConfig) (zapcore.Core, io.Closer, error) {
	level := zap.NewAtomicLevelAt(parseLevel(cfg.Level))
	encoder := newEncoder(cfg.Format)
	writeSyncer, closer, err := newWriteSyncer(cfg)
	if err != nil {
		return nil, nil, err
	}

	return zapcore.NewCore(encoder, writeSyncer, level), closer, nil
}

func newEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LevelKey = "L"
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.CallerKey = "C"
	encoderConfig.MessageKey = "M"
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	switch strings.ToLower(format) {
	case "", "console":
		return zapcore.NewConsoleEncoder(encoderConfig)
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

func newWriteSyncer(cfg config.LoggingConfig) (zapcore.WriteSyncer, io.Closer, error) {
	switch strings.ToLower(cfg.Output) {
	case "", "stdout":
		return zapcore.AddSync(os.Stdout), closerFunc(noopClose), nil
	case "file":
		writer, err := newRollingFileWriter(cfg)
		if err != nil {
			return nil, nil, err
		}

		return zapcore.AddSync(writer), writer, nil
	case "both":
		writer, err := newRollingFileWriter(cfg)
		if err != nil {
			return nil, nil, err
		}

		return zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(writer),
		), writer, nil
	default:
		return nil, nil, fmt.Errorf("unsupported LOG_OUTPUT value: %s", cfg.Output)
	}
}

func newRollingFileWriter(cfg config.LoggingConfig) (*lumberjack.Logger, error) {
	if err := ensureLogDir(cfg.FilePath); err != nil {
		return nil, err
	}

	return &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}, nil
}

func ensureLogDir(path string) error {
	if path == "" {
		return fmt.Errorf("LOG_FILE_PATH is empty")
	}

	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return nil
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

type closerFunc func() error

func (c closerFunc) Close() error {
	return c()
}

func noopClose() error {
	return nil
}
