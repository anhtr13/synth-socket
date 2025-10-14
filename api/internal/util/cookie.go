package util

import (
	"net/http"

	"github.com/anhtr13/synth-socket/api/internal/conf"
)

func SetCookie(w http.ResponseWriter, key, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    value,
		MaxAge:   conf.AGE_COOKIE,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Path:     "/",
	})
}

func DeleteCookie(w http.ResponseWriter, key string) {
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Path:     "/",
	})
}
