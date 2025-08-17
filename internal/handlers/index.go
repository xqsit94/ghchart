package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type IndexHandler struct {
	template *template.Template
}

func NewIndexHandler() (*IndexHandler, error) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", "index.html"))
	if err != nil {
		return nil, err
	}
	
	return &IndexHandler{
		template: tmpl,
	}, nil
}

func (h *IndexHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := h.template.Execute(w, nil); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}