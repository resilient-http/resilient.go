package strategies

import (
	"net/http"
	"time"
)

type Handler func(*http.Request, *http.Response, error) bool

type Strategy func() Handler

// TODO
type Resolver struct {
	Stop  bool
	Retry bool
	RetryWait time.Duration
}

func (r *Resolver) ShouldRetry() bool {
	return r.Retry
}

func (r *Resolver) ShouldStop() bool {
	return r.Stop
}
