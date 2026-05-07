package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ctf-platform/config"
	"ctf-platform/db"
	"ctf-platform/handlers"
	"ctf-platform/middleware"
	"ctf-platform/services"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:3000" || origin == "http://localhost:5173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := config.Load()
	database := db.Connect(cfg.DBURL)

	if err := runMigrations(database); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	challengeSvc := services.NewChallengeService(database)
	submissionSvc := services.NewSubmissionService(database)
	scoreboardSvc := services.NewScoreboardService(database)
	adminSvc := services.NewAdminService(database)
	if err := adminSvc.EnsureDefaultAdmin(services.DefaultAdminInput{
		Username: cfg.AdminUsername,
		Email:    cfg.AdminEmail,
		Password: cfg.AdminPassword,
	}); err != nil {
		log.Fatalf("default admin seed failed: %v", err)
	}

	authH := handlers.NewAuthHandler(database, cfg.JWTSecret)
	challengeH := handlers.NewChallengeHandler(challengeSvc)
	submissionH := handlers.NewSubmissionHandler(submissionSvc)
	scoreboardH := handlers.NewScoreboardHandler(scoreboardSvc)
	adminH := handlers.NewAdminHandler(adminSvc)

	r := chi.NewRouter()
	r.Use(corsMiddleware)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			writeHealth(w)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.JWTSecret))

			r.Get("/challenges", challengeH.List)
			r.Get("/challenges/{id}", challengeH.GetByID)
			r.Post("/submissions", submissionH.Submit)
			r.Get("/scoreboard", scoreboardH.Leaderboard)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.JWTSecret))
			r.Use(middleware.AdminOnly)

			r.Get("/admin/challenges", adminH.ListChallenges)
			r.Post("/admin/challenges", adminH.CreateChallenge)
			r.Put("/admin/challenges/{id}", adminH.UpdateChallenge)
			r.Delete("/admin/challenges/{id}", adminH.DeleteChallenge)
			r.Patch("/admin/challenges/{id}/visibility", adminH.ToggleVisibility)
			r.Get("/admin/submissions", adminH.ListSubmissions)
			r.Patch("/admin/users/{id}/disable", adminH.DisableUser)
			r.Get("/admin/users", adminH.ListUsers)
			r.Post("/admin/seed", adminH.Seed)
		})
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func writeHealth(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func runMigrations(database *sqlx.DB) error {
	sql, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("reading migration file: %w", err)
	}
	_, err = database.Exec(string(sql))
	return err
}
