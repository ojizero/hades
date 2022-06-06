package gracefully_test

import (
	"net/http"

	"github.com/ojizero/hades/gracefully"
)

func ExampleServeHandler() {
	var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}
	addr := ":4000"
	gracefully.ServeHandler(h, addr)
}

func ExampleServe() {
	var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}
	addr := ":4000"
	gracefully.Serve(&http.Server{
		Addr:    addr,
		Handler: h,
	})
}
