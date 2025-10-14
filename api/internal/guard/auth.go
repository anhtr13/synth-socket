package guard

import (
	"net/http"

	"github.com/anhtr13/synth-socket/api/internal/conf"
	"github.com/anhtr13/synth-socket/api/internal/util"
)

func Auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(conf.USER_CTX_KEY).(map[string]any)
		if !ok {
			util.WriteError(w, 401, "unauthorized: access_token failed")
			return
		}
		handler(w, r)
	}
}
