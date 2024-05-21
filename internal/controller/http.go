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

func (hc *httpController) indexHandler(w http.ResponseWriter, req *http.Request) {

	todos, err := hc.repository.ListTODOS(req.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, item := range todos {
		fmt.Fprintf(w, "%s\n", item)
	}

}

func (hc *httpController) postHandler(w http.ResponseWriter, req *http.Request) {

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

func (hc *httpController) putHandler(w http.ResponseWriter, req *http.Request) {

	oldItem := req.URL.Query().Get("old_item")
	if oldItem == "" {
		http.Error(w, "provide the old item to get replaced", http.StatusBadRequest)
		return
	}

	newItem := req.URL.Query().Get("new_item")
	if newItem == "" {
		http.Error(w, "provide the new item to be replaced", http.StatusBadRequest)
		return
	}

	err := hc.repository.ReplaceTODO(req.Context(), newItem, oldItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (hc *httpController) deleteHandler(w http.ResponseWriter, req *http.Request) {

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

	mux.HandleFunc("GET /", hc.indexHandler)
	mux.HandleFunc("POST /", hc.postHandler)
	mux.HandleFunc("PUT /update", hc.putHandler)
	mux.HandleFunc("DELETE /delete", hc.deleteHandler)

	return &hc
}
