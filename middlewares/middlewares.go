package middlewares

import (
	"gopkg.in/resilient-http/resilient.go.v0/client"
	"net/http"
)

type Handler func(Resilient) (Middleware, error)

type Resilient interface {
	Use(Handler) error
	UseServers([]string)
}

type Middleware interface {
	In(*http.Request, *http.Response, error) error
	Out(*client.Options) error
}

type Middlewares struct {
	pool []Middleware
}

func New() *Middlewares {
	return &Middlewares{
		pool: []Middleware{},
	}
}

func (m *Middlewares) Add(mw Middleware) {
	m.pool = append(m.pool, mw)
}

func (m *Middlewares) DispatchIn(req *http.Request, res *http.Response, resErr error) error {
	var err error
	for _, middleware := range m.pool {
		err = middleware.In(req, res, resErr)
	}
	return err
}

func (m *Middlewares) DispatchOut(o *client.Options) error {
	var err error
	for _, middleware := range m.pool {
		err = middleware.Out(o)
	}
	return err
}
