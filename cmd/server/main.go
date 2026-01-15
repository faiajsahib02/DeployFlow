package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sahib002/deployflow/internal/adapter/handler"
	"github.com/sahib002/deployflow/internal/adapter/proxy" // Import Proxy
	"github.com/sahib002/deployflow/internal/adapter/runtime"
	"github.com/sahib002/deployflow/internal/adapter/storage/postgres"
	"github.com/sahib002/deployflow/internal/core/services"
)

func main() {
	// 1. Database
	dsn := "postgres://postgres:postgres@localhost:5435/deployflow?sslmode=disable"
	repo, err := postgres.NewRepository(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL")

	// 2. Docker
	docker, err := runtime.NewDockerClient()
	if err != nil {
		log.Fatalf("Failed to connect to Docker: %v", err)
	}
	log.Println("🐳 Connected to Docker Engine")

	// 3. Setup the Management API (Port 8080)
	svc := services.NewDeploymentService(repo, docker)
	h := handler.NewHandler(repo, svc)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Mount("/api/v1", h.Routes())

	// Start API in a background Goroutine so it doesn't block
	go func() {
		log.Println("🚀 Management API listening on :8080")
		if err := http.ListenAndServe(":8080", r); err != nil {
			log.Fatal(err)
		}
	}()

	// 4. Setup the Proxy Server (Port 8000)
	// This will handle requests like "cat-dog.localhost:8000"
	p := proxy.NewProxyServer(repo)

	log.Println("🌍 Proxy Server listening on :8000 (Use project-name.localhost:8000)")
	if err := http.ListenAndServe(":8000", p); err != nil {
		log.Fatal(err)
	}
}
