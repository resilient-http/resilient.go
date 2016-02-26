package resilient

import (
	"gopkg.in/eapache/go-resiliency.v1/retrier"
	"gopkg.in/h2non/gentleman-retry.v0"
	"gopkg.in/h2non/gentleman.v0"
	"net/http"
)

func New() *gentleman.Client {
	return gentleman.New()
}

func NewRetryClient() *gentleman.Client {
	cli := New()
	cli.Use(retry.New(nil))
	return cli
}

func NewExponentialRetryClient() *gentleman.Client {
	cli := New()
	cli.Use(retry.New(retry.New(retrier.New(retrier.ExponentialBackoff(3, 100*time.Millisecond), nil))))
	return cli
}

// TODO: add consul client
