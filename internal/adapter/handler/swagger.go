//go:build !release

package handler

import (
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/typedef-tokyo/lessonlink-backend/docs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/server"
)

/*
swagger出力コマンド
swag init --requiredByDefault
*/
func initSwagger(env configs.EnvConfig, sever *server.Server) {
	if env.Environment == "local" {
		sever.Engine.GET("/swagger/*", echoSwagger.WrapHandler)
	}
}
