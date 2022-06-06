// Package gracefully provides the functionality to server `net/http'
// servers while providing graceful shutdown for them when
// encountering OS interrupt and termination signals.
//
package gracefully

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var defaultOpts = opts{
	log:           defaultLog,
	graceDuration: 30 * time.Second,
}

// ServeHandler wraps the given handler and address as an `http.Server'
// and then executes `Serve' to serve the newly created server
// with graceful shutdown.
//
func ServeHandler(h http.Handler, addr string, opts ...Option) {
	Serve(&http.Server{
		Addr:    addr,
		Handler: h,
	}, opts...)
}

// Serve serves the given `http.Server' while protecting it from abrupt
// shutdowns by capturing OS signals of `SIGINT' and `SIGTERM' and
// gracefully stopping the server given the duration configured.
//
func Serve(srv *http.Server, opts ...Option) {
	cfg := defaultOpts
	cfg.afterShutdown = []StopFunc{srv.Shutdown}
	for _, o := range opts {
		cfg = o(cfg)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go mustListenAndServe(srv, cfg)
	<-ctx.Done()
	stop()
	cfg.log("Shutting down gracefully, press Ctrl+C again to force")
	mustStopWithGrace(cfg)
}

func mustListenAndServe(srv *http.Server, o opts) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		o.log("listen: %v", err)
		os.Exit(1)
	}
}

func mustStopWithGrace(o opts) {
	ctx, cancel := context.WithTimeout(context.Background(), o.graceDuration)
	defer cancel()
	for _, s := range o.afterShutdown {
		if err := s(ctx); err != nil {
			o.log("failed to gracefully stop. Got: %v", err)
			os.Exit(1)
		}
	}
}
