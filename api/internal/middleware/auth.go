package middleware

import (
	"context"
	"net/http"

	"github.com/anhtr13/synth-socket/api-service/internal/conf"
	"github.com/anhtr13/synth-socket/api-service/internal/util"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		access_token := r.Header.Get("Authorization")
		if access_token == "" {
			cookie, err := r.Cookie("access_token")
			if err == nil {
				access_token = cookie.Value
			}
		}
		if access_token == "" {
			next.ServeHTTP(w, r)
			return
		}
		userClaim, err := util.VerifyJWT(access_token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), conf.USER_CTX_KEY, userClaim)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
