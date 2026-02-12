package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/injector"
)

// @title lessonlink-backend
// @version 1.0
// @license.name typedef-tokyo.
// @description lessonlink-backend バックエンドAPI
// @host localhost:3002
// @BasePath /
func main() {

	// swag init --requiredByDefault コマンドでswagger出力

	// ローカル実行の場合のみ
	environment := os.Getenv("ENVIRONMENT")
	if environment == "local" || environment == "local_test" {
		err := godotenv.Overload("./internal/configs/.env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}

		if environment == "local_test" {
			os.Setenv("DB_NAME", os.Getenv("TEST_DB_NAME"))
		}
	}

	// 環境変数のチェック
	configs.LoadConfig()

	if err := injector.RunInjectedServer(); err != nil {
		log.Fatalln(err)
	}
}
