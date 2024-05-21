package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/kuberise/todolist/internal/gateway"
	"golang.org/x/sync/errgroup"
)

type HTTPConfig struct {
	Port            int `yaml:"port"`
	ShutdownTimeout int `yaml:"shutdown_timeout"`
}

type httpController struct {
	config     *HTTPConfig
	server     *http.Server
	logger     *slog.Logger
	repository gateway.Respository
}

func (hc *httpController) Run(ctx context.Context) error {

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {

		hc.logger.Info("starting the http server", "port", hc.config.Port)
		err := hc.server.ListenAndServe()
		if err == http.ErrServerClosed {
			hc.logger.Info("http server stopped listening to new requests")
			return nil
		}

		return err
	})

	g.Go(func() error {

		<-gCtx.Done()

		hc.logger.Info("http server graceful shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(hc.config.ShutdownTimeout)*time.Second)
		defer cancel()

		return hc.server.Shutdown(ctx)
	})

	return g.Wait()
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
}

func (hc *httpController) indexHandler(w http.ResponseWriter, req *http.Request) {

	enableCORS(w)

	todos, err := hc.repository.ListTODOS(req.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(b))

}

func (hc *httpController) postHandler(w http.ResponseWriter, req *http.Request) {

	enableCORS(w)

	var request struct {
		Item string `json:"item"`
	}

	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = hc.repository.NewTODO(req.Context(), request.Item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (hc *httpController) deleteHandler(w http.ResponseWriter, req *http.Request) {

	enableCORS(w)

	item := req.URL.Query().Get("item")
	if item == "" {
		http.Error(w, "provide the item to be deleted", http.StatusBadRequest)
		return
	}

	err := hc.repository.RemoveTODO(req.Context(), item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewHTTPController(l *slog.Logger, c *HTTPConfig, r gateway.Respository) *httpController {

	mux := http.NewServeMux()

	hc := httpController{
		config:     c,
		logger:     l,
		repository: r,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", c.Port),
			Handler: mux,
		},
	}

	mux.HandleFunc("OPTIONS /", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /", hc.indexHandler)
	mux.HandleFunc("POST /", hc.postHandler)
	mux.HandleFunc("DELETE /delete", hc.deleteHandler)

	return &hc
}
