package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *API) BindRoutes() {
	api.Router.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger, api.Sessions.LoadAndSave)

	// Uncomment the above lines to enable CSRF protection
	// csrfMiddleware := csrf.Protect(
	// 	[]byte(os.Getenv("GOBID_CSRF_KEY")),
	// 	csrf.Secure(false), // Setado como false para desenvolvimento, deve ser true em produção

	// )
	// api.Router.Use(csrfMiddleware)

	// Todo: Implementar rota para listar produtos
	// Todo: Implementar rota para listar usuários
	// Todo: inplementar rota para buscar usuário/produtos por ID

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// r.Get("/csrf-token", api.HandleGetCSRFToken)
			r.Route("/users", func(r chi.Router) {
				r.Post("/signup", api.HandleSignUp)
				r.Post("/login", api.HandleLogin)

				r.Group(func(r chi.Router) {
					r.Use(api.AuthMiddleware)
					r.Post("/logout", api.HandleLogout)
				})
			})
			r.Route("/products", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					// Todas a rotas deste grupo requerem autenticação
					r.Use(api.AuthMiddleware)
					r.Post("/", api.HandlerCreateProduct)
					r.Get("/ws/subscribe/{product_id}", api.HandleSubscribeToAuction)
			})
		})
	})
}
