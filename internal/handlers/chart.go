package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xqsit94/ghchart/internal/chart"
	"github.com/xqsit94/ghchart/internal/services"
)

type ChartHandler struct {
	generator *chart.Generator
}

func NewChartHandler(generator *chart.Generator) *ChartHandler {
	return &ChartHandler{
		generator: generator,
	}
}

func (h *ChartHandler) DefaultChart(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	username = services.RemoveExtension(username, ".svg")
	
	chartData, err := h.generator.Generate(username, "")
	if err != nil {
		http.Error(w, "Failed to generate chart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(chartData)
}

func (h *ChartHandler) CustomColorChart(w http.ResponseWriter, r *http.Request) {
	color := chi.URLParam(r, "color")
	username := chi.URLParam(r, "username")
	
	if username == "" || color == "" {
		http.Error(w, "Username and color are required", http.StatusBadRequest)
		return
	}

	username = services.RemoveExtension(username, ".svg")
	
	chartData, err := h.generator.Generate(username, color)
	if err != nil {
		http.Error(w, "Failed to generate chart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(chartData)
}