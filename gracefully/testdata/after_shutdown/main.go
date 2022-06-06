package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/ojizero/hades/gracefully"
)

func main() {
	var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.Write([]byte("hello there"))
	}
	gracefully.ServeHandler(h, ":4000", gracefully.AfterShutdown(func(ctx context.Context) error {
		os.Create("./testdata/after_shutdown/after_shutdown.txt")
		return nil
	}))
}
