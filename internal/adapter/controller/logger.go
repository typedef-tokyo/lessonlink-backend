package controller

import "github.com/labstack/echo/v4"

type (
	ILogWriter interface {
		WriteErrLog(e echo.Context, err error) (int, string)
	}
)
