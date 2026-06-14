package router

import (
	"log/slog"
	"net/http"
	"os"
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
	"github.com/ticketing/api/internal/service"
)

func init() {
	lvl := slog.LevelInfo
	if os.Getenv("APP_ENV") == "development" {
		lvl = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})))
}

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

	// Services
	warSvc := service.NewWarQueue(rdb, cfg.QueueSessionTTLMinutes, cfg.WarRateLimitPerMin)

	// Repositories
	bookingRepo := repository.NewBookingRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	ticketRepo := repository.NewTicketRepository(db)

	// Handlers
	authH := handler.NewAuthHandler(userRepo, jwtMgr)
	eventH := handler.NewEventHandler(eventRepo)
	catH := handler.NewTicketCategoryHandler(catRepo)
	warH := handler.NewWarHandler(warSvc, eventRepo, catRepo)
	bookingH := handler.NewBookingHandler(bookingRepo, catRepo, warSvc, eventRepo, cfg.BookingTTLMinutes)
	paymentH := handler.NewPaymentHandler(paymentRepo, bookingRepo, catRepo, ticketRepo)
	ticketH := handler.NewTicketHandler(ticketRepo, bookingRepo)
	adminH := handler.NewAdminHandler(bookingRepo, userRepo, eventRepo)

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
			// IMPORTANT: declare nested eventId-scoped routes here (not in
			// the outer `/admin` subroute) because chi matches the
			// `/admin/events` prefix first and won't fall through to a
			// sibling `/admin` route on miss.
			r.Post("/{eventId}/categories", catH.Create)
		})

		// Admin routes
		r.Route("/admin", func(r chi.Router) {
			r.Use(jwtMgr.Authenticator, auth.RequireAdmin)
			r.Put("/categories/{id}", catH.Update)
			r.Get("/stats", adminH.Stats)
			r.Get("/bookings", adminH.ListBookings)
			r.Get("/bookings/{id}", adminH.DetailBooking)
			r.Get("/payments", paymentH.ListAll)
			r.Get("/users", adminH.ListUsers)
			r.Post("/users", adminH.CreateUser)
		})

		// War queue (authenticated)
		r.Route("/war", func(r chi.Router) {
			r.Use(jwtMgr.Authenticator)
			r.Post("/join", warH.Join)
			r.Get("/status", warH.Status)
		})

		// Bookings (authenticated)
		r.Route("/bookings", func(r chi.Router) {
			r.Use(jwtMgr.Authenticator)
			r.Post("/reserve", bookingH.Reserve)
			r.Get("/me", bookingH.ListMy)
			r.Get("/{id}", bookingH.Detail)
			r.Post("/{id}/cancel", bookingH.Cancel)
		})

		// Payments (authenticated)
		r.Route("/payments", func(r chi.Router) {
			r.Use(jwtMgr.Authenticator)
			r.Post("/create", paymentH.Create)
			r.Get("/me", paymentH.ListMy)
			r.Post("/{id}/simulate", paymentH.Simulate)
		})

		// Tickets (authenticated)
		r.Route("/tickets", func(r chi.Router) {
			r.Use(jwtMgr.Authenticator)
			r.Get("/", ticketH.ListMy)
			r.Get("/{id}", ticketH.Detail)
			r.Post("/verify/{code}", ticketH.Verify)
		})
	})

	return r
}
