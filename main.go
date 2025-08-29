package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Repository (swap this with a DB-backed repo in production)
	repo := NewMemoryRepo()

	// Seed sample tasks (only if enabled via env var)
	if getEnv("SEED_DATA", "false") == "true" {
		seed(repo)
		log.Println("✅ Seed data loaded")
	}

	// Service
	svc := NewTaskService(repo)

	// Handler / Router
	h := NewHandler(svc)

	// Start server
	addr := ":" + getEnv("PORT", "8080")
	log.Printf("🚀 Task service starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, h.Routes()))
}

func seed(r Repository) {
	_, _ = r.Create(nil, &Task{
		Title:       "Buy groceries",
		Description: "Milk, eggs, bread",
		Status:      StatusPending,
	})
	_, _ = r.Create(nil, &Task{
		Title:       "Deploy release",
		Description: "Deploy v1.2 to staging",
		Status:      StatusInProgress,
	})
	_, _ = r.Create(nil, &Task{
		Title:       "Retro meeting",
		Description: "Sprint retro",
		Status:      StatusCompleted,
	})
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
