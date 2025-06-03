package api

import (
	"net/http"

	"github.com/bruguedes/gobid/internal/jsonutils"
)

func (api *API) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), "AuthenticatedUserID") {
			jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]any{
				"error": "must be logged in to access this resource",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
