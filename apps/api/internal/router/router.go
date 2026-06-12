package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/ticketing/api/internal/auth"
	"github.com/ticketing/api/internal/config"
	"github.com/ticketing/api/internal/handler"
	"github.com/ticketing/api/internal/repository"
)

func New(cfg *config.Config, db *gorm.DB, rdb *redis.Client) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// JWT manager
	jwtMgr := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiresHours)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)

	// Handlers
	authH := handler.NewAuthHandler(userRepo, jwtMgr)
	eventH := handler.NewEventHandler(eventRepo)
	catH := handler.NewTicketCategoryHandler(catRepo)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Public auth
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)

			// Protected
			r.Group(func(r chi.Router) {
				r.Use(jwtMgr.Authenticator)
				r.Get("/me", authH.Me)
				r.Put("/me", authH.UpdateMe)
			})
		})

		// Public events
		r.Get("/events", eventH.ListPublic)
		r.Get("/events/{id}", eventH.DetailPublic)
		r.Get("/events/{eventId}/categories", catH.ListByEvent)

		// Admin events
		r.Route("/admin/events", func(r chi.Router) {
			r.Use(jwtMgr.Authenticator, auth.RequireAdmin)
			r.Get("/", eventH.ListAdmin)
			r.Post("/", eventH.Create)
			r.Get("/{id}", eventH.DetailAdmin)
			r.Put("/{id}", eventH.Update)
			r.Delete("/{id}", eventH.Delete)
			r.Post("/{id}/publish", eventH.Publish)
		})

		// Admin ticket categories
		r.Route("/admin", func(r chi.Router) {
			r.Use(jwtMgr.Authenticator, auth.RequireAdmin)
			r.Post("/events/{eventId}/categories", catH.Create)
			r.Put("/categories/{id}", catH.Update)
		})

		// Mount war, booking, payment, admin (other) here in later sections
	})

	return r
}
