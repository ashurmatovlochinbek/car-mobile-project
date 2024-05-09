package http

import (
	"car-mobile-project/config"
	"car-mobile-project/internal/middleware"
	"car-mobile-project/internal/user"
	"github.com/go-chi/chi/v5"
)

func MapUserRoutes(userRouter *chi.Mux, h user.Handlers, cfg *config.Config) {
	userRouter.Route("/open/user", func(r chi.Router) {
		r.Post("/register", h.Register())
		r.Post("/login", h.Login())
		r.Get("/verifyOTP", h.VerifyOTP())
	})

	userRouter.Route("/secured/user", func(r chi.Router) {
		r.Use(middleware.AuthJwtMiddleware(cfg.Server.JwtSecretKey))
		r.Get("/refresh", h.Refresh())
		r.Get("/getMessage", h.GetSecuredResource())
	})
}
