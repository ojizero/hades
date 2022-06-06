package gracefully

import (
	"context"
	"fmt"
	"os"
	"time"
)

type opts struct {
	afterShutdown []StopFunc
	graceDuration time.Duration
	log           func(string, ...interface{})
}

// Option provides the ability to configure the behaviour of `gracefully'.
type Option func(opts) opts

// StopFunc is a function that stops some service and can be used to
// attach additional services to shutdown after shutting down the
// `http.Server' being served.
//
type StopFunc func(context.Context) error

// WithGraceDuration configures the graceful shutdown duration.
func WithGraceDuration(d time.Duration) Option {
	return func(o opts) opts {
		o.graceDuration = d
		return o
	}
}

// AfterShutdown adds additional stop functions to call within the same
// same grace duration after the server stop.
//
// The grace duration is not renewed after the main `http.Server' stops
// and instead is carried from whatever point it reached, meaning
// that the `http.Server' and all `AfterShutdown' hooks must
// be able to finish execution within one grace duration.
//
func AfterShutdown(fs ...StopFunc) Option {
	return func(o opts) opts {
		o.afterShutdown = append(o.afterShutdown, fs...)
		return o
	}
}

// WithLogger configures how to log the messages logged by `gracefully',
// defaults to logging using the `fmt' package to the standard errors.
//
func WithLogger(l func(string, ...interface{})) Option {
	return func(o opts) opts {
		o.log = l
		return o
	}
}

func defaultLog(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, f, args...)
}
