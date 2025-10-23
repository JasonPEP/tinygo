package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	stdhttp "net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tinygo/internal/logger"
	"tinygo/internal/shortener"
	"tinygo/internal/storage"
)

type Handlers struct {
	svc *shortener.Service
}

func NewHandlers(svc *shortener.Service) *Handlers {
	return &Handlers{svc: svc}
}

// Register registers routes on the given mux.
func (h *Handlers) Register(mux *stdhttp.ServeMux) {
	mux.HandleFunc("/healthz", h.health)
	mux.HandleFunc("/api/shorten", h.shorten)
	mux.HandleFunc("/api/links/", h.linkDetail)
	
	// Admin routes for Web UI
	mux.HandleFunc("/admin/shorten", h.shorten)
	mux.HandleFunc("/admin/stats", h.stats)
	mux.HandleFunc("/admin/links/", h.linkDetail)
	
	// Static files
	mux.Handle("/static/", stdhttp.StripPrefix("/static/", stdhttp.FileServer(stdhttp.Dir("web/static/"))))
	
	mux.HandleFunc("/web", h.webUI)
	mux.HandleFunc("/", h.redirect)
}

func (h *Handlers) health(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	w.WriteHeader(stdhttp.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (h *Handlers) ready(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	// In a real application, you might check database connectivity, etc.
	// For now, just return OK
	w.WriteHeader(stdhttp.StatusOK)
	_, _ = w.Write([]byte("ready"))
}

func (h *Handlers) stats(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	// Get all links for stats
	links, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, stdhttp.StatusInternalServerError, err.Error())
		return
	}

	// Calculate stats
	totalLinks := len(links)
	var totalHits int64
	for _, link := range links {
		totalHits += link.HitCount
	}

	stats := map[string]interface{}{
		"total_links": totalLinks,
		"total_hits":  totalHits,
		"links":       links,
	}

	writeJSON(w, stdhttp.StatusOK, stats)
}

func (h *Handlers) webUI(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	// Serve the main HTML page
	htmlPath := filepath.Join("web", "templates", "index.html")
	content, err := os.ReadFile(htmlPath)
	if err != nil {
		writeError(w, stdhttp.StatusInternalServerError, "template not found")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(stdhttp.StatusOK)
	w.Write(content)
}

type shortenRequest struct {
	LongURL    string `json:"long_url"`
	CustomCode string `json:"custom_code"`
}

type shortenResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

func (h *Handlers) shorten(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	logger.Log.Info("shorten handler called", "method", r.Method, "path", r.URL.Path)
	if r.Method != stdhttp.MethodPost {
		writeError(w, stdhttp.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req shortenRequest
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&req); err != nil {
		writeError(w, stdhttp.StatusBadRequest, "invalid json")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	l, err := h.svc.Shorten(ctx, req.LongURL, req.CustomCode)
	if err != nil {
		switch {
		case errors.Is(err, shortener.ErrInvalidURL), errors.Is(err, shortener.ErrInvalidCode):
			writeError(w, stdhttp.StatusBadRequest, err.Error())
		default:
			writeError(w, stdhttp.StatusConflict, err.Error())
		}
		return
	}
	resp := shortenResponse{Code: l.Code, ShortURL: h.svc.ShortURL(l.Code), LongURL: l.LongURL}
	writeJSON(w, stdhttp.StatusCreated, resp)
}

func (h *Handlers) linkDetail(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	// path: /api/links/{code}
	if !strings.HasPrefix(r.URL.Path, "/api/links/") {
		writeError(w, stdhttp.StatusNotFound, "not found")
		return
	}
	code := strings.TrimPrefix(r.URL.Path, "/api/links/")
	switch r.Method {
	case stdhttp.MethodGet:
		l, ok, err := h.svc.Resolve(r.Context(), code)
		if err != nil {
			writeError(w, stdhttp.StatusInternalServerError, err.Error())
			return
		}
		if !ok {
			writeError(w, stdhttp.StatusNotFound, "not found")
			return
		}
		writeJSON(w, stdhttp.StatusOK, l)
	case stdhttp.MethodDelete:
		if err := h.svc.Delete(r.Context(), code); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				writeError(w, stdhttp.StatusNotFound, "not found")
				return
			}
			writeError(w, stdhttp.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(stdhttp.StatusNoContent)
	default:
		writeError(w, stdhttp.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handlers) redirect(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	// Show Web UI for root path
	if r.URL.Path == "/" {
		h.webUI(w, r)
		return
	}
	
	// Skip reserved paths
	if strings.HasPrefix(r.URL.Path, "/api/") ||
		strings.HasPrefix(r.URL.Path, "/admin/") ||
		r.URL.Path == "/healthz" ||
		r.URL.Path == "/readyz" ||
		r.URL.Path == "/web" {
		writeError(w, stdhttp.StatusNotFound, "not found")
		return
	}

	// Extract code from path
	code := strings.TrimPrefix(r.URL.Path, "/")

	// Hit the link (increment counter)
	l, err := h.svc.Hit(r.Context(), code)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, stdhttp.StatusNotFound, "not found")
			return
		}
		writeError(w, stdhttp.StatusInternalServerError, err.Error())
		return
	}

	// Redirect to the long URL
	stdhttp.Redirect(w, r, l.LongURL, stdhttp.StatusFound)
}

// --- helpers ---

func writeJSON(w stdhttp.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w stdhttp.ResponseWriter, status int, msg string) {
	type errResp struct {
		Error string `json:"error"`
	}
	writeJSON(w, status, errResp{Error: msg})
}
