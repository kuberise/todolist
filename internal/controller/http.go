package controller

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type HTTPConfig struct {
	Port            int
	ShutdownTimeout time.Duration
}

type httpController struct {
	config *HTTPConfig
	server *http.Server
	logger *slog.Logger
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

		ctx, cancel := context.WithTimeout(context.Background(), hc.config.ShutdownTimeout*time.Second)
		defer cancel()

		return hc.server.Shutdown(ctx)
	})

	return g.Wait()
}

func (hc *httpController) SetItem(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("wtite item"))
}

func (hc *httpController) RemoveItem(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("remove item"))
}

func (hc *httpController) UpdateItem(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("update item"))
}

func (hc *httpController) ListItems(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("list items"))
}

func NewHTTPController(l *slog.Logger, c *HTTPConfig) *httpController {

	mux := http.NewServeMux()

	hc := httpController{
		config: c,
		logger: l,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", c.Port),
			Handler: mux,
		},
	}

	mux.HandleFunc("/set_item", hc.SetItem)
	mux.HandleFunc("/remove_item", hc.RemoveItem)
	mux.HandleFunc("/update_item", hc.UpdateItem)
	mux.HandleFunc("/list_items", hc.ListItems)

	return &hc
}
