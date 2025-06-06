package api

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
)

func (api *API) BindRoutes() {
	api.Router.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger, api.Sessions.LoadAndSave)

	csrfMiddleware := csrf.Protect(
		[]byte(os.Getenv("GOBID_CSRF_KEY")),
		csrf.Secure(false), // Setado como false para desenvolvimento, deve ser true em produção

	)
	api.Router.Use(csrfMiddleware)

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/csrf-token", api.HandleGetCSRFToken)
			r.Route("/users", func(r chi.Router) {
				r.Post("/signup", api.HandleSignUp)
				r.Post("/login", api.HandleLogin)
				r.With(api.AuthMiddleware).Post("/logout", api.HandleLogout)
			})
		})
	})
}
