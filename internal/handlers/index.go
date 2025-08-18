package handlers

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type IndexHandler struct {
	template *template.Template
}

type IndexData struct {
	DomainPrefix string
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
	domainPrefix := os.Getenv("DOMAIN_PREFIX")
	if domainPrefix == "" {
		domainPrefix = "https://ghchart.xqsit94.in"
	}
	
	data := IndexData{
		DomainPrefix: domainPrefix,
	}
	
	w.Header().Set("Content-Type", "text/html")
	if err := h.template.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}