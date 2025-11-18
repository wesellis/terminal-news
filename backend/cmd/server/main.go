package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/wesellis/terminal-news/backend/internal/api/handlers"
	apimiddleware "github.com/wesellis/terminal-news/backend/internal/api/middleware"
	"github.com/wesellis/terminal-news/backend/internal/database"
	"github.com/wesellis/terminal-news/backend/internal/services"
	"github.com/wesellis/terminal-news/backend/internal/workers"
	"github.com/wesellis/terminal-news/backend/pkg/websocket"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize Redis
	rdb := database.InitRedis()
	defer rdb.Close()

	log.Println("✅ Database and Redis connected")

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()
	log.Println("✅ WebSocket hub running")

	// Initialize services
	authService := services.NewAuthService(db)
	articleService := services.NewArticleService(db, rdb)
	voteService := services.NewVoteService(db, rdb)
	commentService := services.NewCommentService(db)
	classifiedService := services.NewClassifiedService(db)
	paymentService := services.NewPaymentService(db)

	// Initialize handlers
	h := handlers.NewHandler(
		authService,
		articleService,
		voteService,
		commentService,
		classifiedService,
		paymentService,
		wsHub,
	)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // TODO: Restrict in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"terminal-news-api"}`))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			// Authentication
			r.Post("/auth/register", h.HandleRegister)
			r.Post("/auth/login", h.HandleLogin)
			r.Post("/auth/refresh", h.HandleRefreshToken)

			// Articles (public read)
			r.Get("/articles", h.HandleGetArticles)
			r.Get("/articles/hot", h.HandleGetHotArticles)
			r.Get("/articles/controversial", h.HandleGetControversialArticles)
			r.Get("/articles/rising", h.HandleGetRisingArticles)
			r.Get("/articles/{id}", h.HandleGetArticle)
			r.Get("/articles/{id}/comments", h.HandleGetComments)

			// Classifieds (public read)
			r.Get("/classifieds", h.HandleGetClassifieds)
			r.Get("/classifieds/{id}", h.HandleGetClassified)

			// Weather
			r.Get("/weather", h.HandleGetWeather)

			// Webhooks (Stripe)
			r.Post("/webhooks/stripe", h.HandleStripeWebhook)
		})

		// Protected routes (require authentication)
		r.Group(func(r chi.Router) {
			r.Use(apimiddleware.AuthRequired(authService))

			// User
			r.Get("/user/profile", h.HandleGetProfile)
			r.Put("/user/profile", h.HandleUpdateProfile)
			r.Get("/user/activity", h.HandleGetActivity)

			// Voting
			r.Post("/articles/{id}/vote", h.HandleVote)
			r.Delete("/articles/{id}/vote", h.HandleDeleteVote)

			// Comments
			r.Post("/articles/{id}/comments", h.HandlePostComment)
			r.Put("/comments/{id}", h.HandleUpdateComment)
			r.Delete("/comments/{id}", h.HandleDeleteComment)
			r.Post("/comments/{id}/vote", h.HandleVoteComment)

			// Classifieds
			r.Post("/classifieds", h.HandlePostClassified)
			r.Put("/classifieds/{id}", h.HandleUpdateClassified)
			r.Delete("/classifieds/{id}", h.HandleDeleteClassified)
			r.Post("/classifieds/{id}/boost", h.HandleBoostClassified)

			// Payments
			r.Post("/payments/create-intent", h.HandleCreatePaymentIntent)
			r.Get("/payments/history", h.HandleGetPaymentHistory)
		})
	})

	// WebSocket
	r.Get("/ws", h.HandleWebSocket)

	// Start background workers
	scheduler := workers.NewScheduler(db, rdb)
	ctx := context.Background()
	go scheduler.Start(ctx)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", port)
		log.Printf("📡 Health check: http://localhost:%s/health", port)
		log.Printf("📚 API base: http://localhost:%s/api/v1", port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Server stopped gracefully")
}
