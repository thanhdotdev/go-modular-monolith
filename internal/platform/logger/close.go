package logger

import (
	"io"

	"go.uber.org/zap"
)

func Close(log *zap.Logger, closer io.Closer) {
	if log != nil {
		_ = log.Sync()
	}

	if closer == nil {
		return
	}

	if err := closer.Close(); err != nil && log != nil {
		log.Error("failed to close log output", zap.Error(err))
	}
}
