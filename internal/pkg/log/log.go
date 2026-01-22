package log

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strings"
)

const MUDULE_NAME = "github.com/typedef-tokyo/lessonlink-backend/internal"

const (
	INTERNAL_ERROR = http.StatusInternalServerError
	BAD_REQUEST    = http.StatusBadRequest
	NOT_FOUND      = http.StatusNotFound
	UNAUTHORIZED   = http.StatusUnauthorized
	CONFLICT       = http.StatusConflict
	FORBIDDEN      = http.StatusForbidden
)

var (
	LogLevelMap = map[int]slog.Level{
		INTERNAL_ERROR: slog.LevelError,
		BAD_REQUEST:    slog.LevelWarn,
		NOT_FOUND:      slog.LevelWarn,
		UNAUTHORIZED:   slog.LevelWarn,
		CONFLICT:       slog.LevelWarn,
		FORBIDDEN:      slog.LevelWarn,
	}

	LogSeverityMap = map[int]string{
		INTERNAL_ERROR: "ERROR",
		BAD_REQUEST:    "WARNING",
		NOT_FOUND:      "WARNING",
		UNAUTHORIZED:   "WARNING",
		CONFLICT:       "WARNING",
		FORBIDDEN:      "WARNING",
	}
)

type Error struct {
	StatusCode int
	Message    string
	StackTrace string
}

func (e *Error) Error() string {
	return e.Message
}

func Errorf(format string, a ...any) error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func wrapError(err error, statusCode ...int) error {
	if err == nil {
		return nil
	}

	var stackErr *Error
	if errors.As(err, &stackErr) {
		// すでに *Error 型の場合
		if stackErr.StackTrace == "" {
			stackErr.StackTrace = buildStackTrace()
		}
		// ステータスコードが指定された場合だけ上書き
		if len(statusCode) > 0 {
			stackErr.StatusCode = statusCode[0]
		}
		return stackErr
	}

	// 通常の error の場合
	newErr := &Error{
		Message: err.Error(),
	}
	if len(statusCode) > 0 {
		newErr.StatusCode = statusCode[0]
	}
	newErr.StackTrace = buildStackTrace()
	return newErr
}

func buildStackTrace() string {
	const depth = 32
	pc := make([]uintptr, depth)
	n := runtime.Callers(4, pc[:])
	frames := runtime.CallersFrames(pc[:n])

	var sb strings.Builder
	for {
		frame, more := frames.Next()
		if strings.Contains(frame.Function, MUDULE_NAME) {
			fmt.Fprintf(&sb, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		}
		if !more {
			break
		}
	}
	return sb.String()
}

func WrapErrorWithStackTrace(err error) error { return wrapError(err) }
func WrapErrorWithStackTraceInternalServerError(err error) error {
	return wrapError(err, INTERNAL_ERROR)
}
func WrapErrorWithStackTraceBadRequest(err error) error   { return wrapError(err, BAD_REQUEST) }
func WrapErrorWithStackTraceNotFound(err error) error     { return wrapError(err, NOT_FOUND) }
func WrapErrorWithStackTraceUnauthorized(err error) error { return wrapError(err, UNAUTHORIZED) }
func WrapErrorWithStackTraceConflict(err error) error     { return wrapError(err, CONFLICT) }
func WrapErrorWithStackTraceForbidden(err error) error    { return wrapError(err, FORBIDDEN) }
