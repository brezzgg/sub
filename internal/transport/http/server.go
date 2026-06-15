package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/brezzgg/go-packages/lg"
	"github.com/brezzgg/sub/internal/pkg/errors"
	"github.com/brezzgg/sub/internal/usecase"
)

type Handler struct {
	usec    *usecase.Usecase
	host    string
	srv     *http.Server
	pattern string
}

func New(usec *usecase.Usecase, host string, handlePattern string) *Handler {
	return &Handler{usec: usec, host: host, pattern: handlePattern}
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
	pl, err := h.usec.GetPayload(id)
	if err != nil {
		if errors.CodeIs(err, errors.CodeNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write headers
	for k, v := range pl.Headers {
		w.Header().Set(k, v)
	}

	// write body
	_, err = w.Write([]byte(pl.Body))
	if err != nil {
		lg.Error("get", lg.Ef("failed to write response: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
