package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	logWriter "github.com/typedef-tokyo/lessonlink-backend/internal/infrastructure/logger"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/constants"
	"github.com/typedef-tokyo/lessonlink-backend/internal/usecase/repository"
)

func AuthRequiredMiddleware(env configs.EnvConfig, sessionRepository repository.SessionRepository, logger *logWriter.LogWriter) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			cookie, err := c.Cookie(env.SessionName)
			if err != nil {
				// Cookieが無い場合・取得できなかった場合
				return c.JSON(http.StatusUnauthorized, map[string]string{"msg": "ログインしてください"})
			}

			sessionID := cookie.Value
			if sessionID == "" {
				// クッキーはあったが中身が空の場合
				return c.JSON(http.StatusUnauthorized, map[string]string{"msg": "ログインしてください"})
			}

			// セッションIDからセッション情報を取得
			sessionEntity, err := sessionRepository.Find(c.Request().Context(), sessionID)
			if err != nil {
				status, msg := logger.WriteErrLog(c, err)
				return c.JSON(status, map[string]interface{}{
					"msg": msg,
				})
			}

			if sessionEntity == nil || sessionEntity.ExpiresAt.Before(time.Now()) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"msg": "ログインしてください"})
			}

			// 有効期間を1時間伸ばす
			const LIMITED = 3600 * time.Second
			sessionEntity.ExpiresAt = time.Now().Add(LIMITED)
			err = sessionRepository.Update(c.Request().Context(), nil, *sessionEntity)
			if err != nil {
				status, msg := logger.WriteErrLog(c, err)
				return c.JSON(status, map[string]interface{}{
					"msg": msg,
				})
			}

			// クッキーの有効時間も伸ばす
			newCookie := &http.Cookie{
				Name:     env.SessionName,
				Value:    sessionID,
				Path:     "/",
				Domain:   "",
				MaxAge:   3600,
				Expires:  time.Now().Add(LIMITED),
				Secure:   os.Getenv("LOCAL_TEST") != "true",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}

			c.SetCookie(newCookie)

			c.Set(constants.USER_IDENTIFIER, sessionEntity.UserID)
			c.Set(constants.ROLE_IDENTIFIER, sessionEntity.RoleKey)

			return next(c)
		}
	}
}
