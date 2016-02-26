package resilient

import (
	"gopkg.in/eapache/go-resiliency.v1/retrier"
	"gopkg.in/h2non/gentleman-retry.v0"
	"gopkg.in/h2non/gentleman.v0"
	"time"
)

// New creates a new gentleman based HTTP client.
func New() *gentleman.Client {
	return gentleman.New()
}

// NewRetryClient creates a new HTTP client with constant
// retry back off strategy in case of failure.
func NewRetryClient() *gentleman.Client {
	cli := New()
	cli.Use(retry.New(nil))
	return cli
}

// NewExponentialRetryClient creates a new HTTP client with
// exponential retry strategy in case of failure.
func NewExponentialRetryClient() *gentleman.Client {
	cli := New()
	cli.Use(retry.New(retrier.New(retrier.ExponentialBackoff(3, 100*time.Millisecond), nil)))
	return cli
}
