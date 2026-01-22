package log

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/adapter/controller"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/log"
)

type (
	LogWriter struct {
		env configs.EnvConfig
	}
)

func NewLogWriter(
	env configs.EnvConfig,
) controller.ILogWriter {
	return &LogWriter{
		env: env,
	}
}

func NewLogWriterImplementation(
	env configs.EnvConfig,
) *LogWriter {
	return &LogWriter{
		env: env,
	}
}

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: false,
}))

func (l *LogWriter) WriteErrLog(e echo.Context, err error) (int, string) {
	if err == nil {
		return log.INTERNAL_ERROR, "unknown error"
	}

	var stackErr log.Error
	if eErr := new(log.Error); errors.As(err, &eErr) {
		stackErr = *eErr
	} else {
		stackErr.Message = err.Error()
		stackErr.StatusCode = log.INTERNAL_ERROR
	}

	if stackErr.StatusCode == 0 {
		stackErr.StatusCode = log.INTERNAL_ERROR
	}

	level := log.LogLevelMap[stackErr.StatusCode]
	transactionID := e.Response().Header().Get(echo.HeaderXRequestID)

	if l.env.Environment == "local" {
		fmt.Println(stackErr.Message)
		fmt.Println(stackErr.StackTrace)
		fmt.Println("transactionID", transactionID)
	} else {
		severity := log.LogSeverityMap[stackErr.StatusCode]
		logger.Log(
			e.Request().Context(),
			level,
			stackErr.Message,
			slog.String("transactionID", transactionID),
			slog.String("severity", severity),
			slog.String("stacktrace", stackErr.StackTrace),
		)
	}

	return stackErr.StatusCode, stackErr.Error()
}
