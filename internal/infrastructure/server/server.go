package server

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
)

type Server struct {
	Engine *echo.Echo
}

func NewServer(env configs.EnvConfig) *Server {

	e := echo.New()
	e.Use(session.Middleware(createStore(env)))

	origins := []string{
		"https://lessonlink-frontend-382133459414.asia-northeast1.run.app", // GCP cloud run 環境
	}

	if env.Environment == "local" {
		origins = append(origins, "http://localhost:3031")
	}

	corsConfig := middleware.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAuthorization,
		},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}
	e.Use(middleware.CORSWithConfig(corsConfig))

	e.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	return &Server{
		Engine: e,
	}
}
