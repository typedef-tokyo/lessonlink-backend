package server

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/typedef-tokyo/lessonlink-backend/internal/configs"
)

func createStore(env configs.EnvConfig) *sessions.CookieStore {

	sameSite := http.SameSiteStrictMode
	if os.Getenv("LOCAL_TEST") == "true" {
		sameSite = http.SameSiteNoneMode
	}

	serssionSecretKey := env.SessionSecretKey
	store := sessions.NewCookieStore([]byte(serssionSecretKey))
	store.Options = &sessions.Options{
		Path:        "/",
		Domain:      "",
		MaxAge:      3600,
		Secure:      os.Getenv("LOCAL_TEST") != "true",
		HttpOnly:    true,
		Partitioned: false,
		SameSite:    sameSite,
	}

	return store
}
