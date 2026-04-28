package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
	"github.com/typedef-tokyo/lessonlink-backend/internal/modules/session/adapter/handler"
	"github.com/typedef-tokyo/lessonlink-backend/internal/pkg/constants"
	logWriter "github.com/typedef-tokyo/lessonlink-backend/internal/platform/logger"
)

func AuthRequiredMiddleware(env configs.EnvConfig, sessionGetHandler handler.ISessionGetHandler, logger *logWriter.LogWriter) echo.MiddlewareFunc {

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

			userID, roleKey, err := sessionGetHandler.Execute(c.Request().Context(), sessionID)
			if err != nil {
				status, msg := logger.WriteErrLog(c, err)
				return c.JSON(status, map[string]interface{}{
					"msg": msg,
				})
			}

			// クッキーの有効時間も伸ばす
			const LIMITED = 3600 * time.Second
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

			c.Set(constants.USER_IDENTIFIER, userID)
			c.Set(constants.ROLE_IDENTIFIER, roleKey)

			return next(c)
		}
	}
}
