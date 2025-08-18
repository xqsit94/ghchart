package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xqsit94/ghchart/internal/chart"
	"github.com/xqsit94/ghchart/internal/handlers"
	"github.com/xqsit94/ghchart/internal/services"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	githubService := services.NewGitHubService()
	chartGenerator := chart.NewGenerator(githubService)

	indexHandler, err := handlers.NewIndexHandler()
	if err != nil {
		log.Fatal("Failed to initialize index handler:", err)
	}

	chartHandler := handlers.NewChartHandler(chartGenerator)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Cache-Control", "public, max-age=86400"))

	r.Get("/", indexHandler.Index)
	r.Get("/{username}", chartHandler.DefaultChart)
	r.Get("/{themeColor}/{username}", func(w http.ResponseWriter, r *http.Request) {
		themeColor := chi.URLParam(r, "themeColor")
		if strings.Contains(themeColor, ":") {
			chartHandler.ThemeColorChart(w, r)
		} else {
			chartHandler.CustomColorChart(w, r)
		}
	})

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
