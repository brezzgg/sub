package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/brezzgg/go-packages/lg"
	"github.com/brezzgg/sub/internal/models"
	"github.com/brezzgg/sub/internal/models/payload/v1"
)

type Handler struct {
	repo    models.Repo
	host    string
	srv     *http.Server
	pattern string
}

func New(host string, repo models.Repo, handlePattern string) *Handler {
	return &Handler{repo: repo, host: host, pattern: handlePattern}
}

func (h *Handler) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(h.pattern, func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.get(w, r)
		lg.Info("http request", lg.C{
			"addr":     r.RemoteAddr,
			"method":   r.Method,
			"pattern":  r.Pattern,
			"duration": time.Since(start).Abs().String(),
		})
	})

	h.srv = &http.Server{
		Addr:    h.host,
		Handler: mux,
	}

	return fmt.Errorf("listen error: %s", h.srv.ListenAndServe())
}

func (h *Handler) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := h.srv.Shutdown(ctx)
	_ = h.srv.Close()
	return err
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// get sub
	sub, err := h.repo.Get(id)
	if err != nil {
		lg.Error("get", err)
		if errors.Is(err, models.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// validate sub
	if !sub.Validate() {
		lg.Error("get", models.ErrExpired)
		_ = h.repo.Remove(id)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(models.ErrExpired.Error()))
		return
	}

	// unmarshal payload
	pl, err := payload.Unmarshal(sub.Payload, false)
	if err != nil {
		lg.Error("failed to unmarshal payload", err)
		_ = h.repo.Remove(id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write headers
	for k, v := range pl.GetHeaders() {
		w.Header().Set(k, v)
	}

	// write body
	_, err = w.Write([]byte(pl.GetBody()))
	if err != nil {
		lg.Error("get", lg.Ef("failed to write response: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
