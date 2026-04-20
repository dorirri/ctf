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

	authH := handlers.NewAuthHandler(database, cfg.JWTSecret)
	challengeH := handlers.NewChallengeHandler(challengeSvc)
	submissionH := handlers.NewSubmissionHandler(submissionSvc)
	scoreboardH := handlers.NewScoreboardHandler(scoreboardSvc)
	adminH := handlers.NewAdminHandler(adminSvc)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)

	r.Route("/api", func(r chi.Router) {
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

func runMigrations(database *sqlx.DB) error {
	sql, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("reading migration file: %w", err)
	}
	_, err = database.Exec(string(sql))
	return err
}
