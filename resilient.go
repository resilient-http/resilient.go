package resilient

import (
	"gopkg.in/resilient-http/resilient.go.v0/client"
	"gopkg.in/resilient-http/resilient.go.v0/middlewares"
	"io"
	"net/http"
	"sync"
)

type Strategy func() StrategyHandler

type StrategyHandler func(*http.Request, *http.Response, error) bool

type Resilient struct {
	sync.Mutex
	Servers     []string
	strategies  []Strategy
	middlewares *middlewares.Middlewares
}

func New() *Resilient {
	return &Resilient{
		middlewares: middlewares.New(),
	}
}

func (r *Resilient) UseServers(servers []string) {
	r.Lock()
	defer r.Unlock()
	r.Servers = servers
}

func (r *Resilient) Use(m middlewares.Handler) error {
	r.Lock()
	defer r.Unlock()

	handler, err := m(r)
	if err == nil {
		r.middlewares.Add(handler)
	}

	return err
}

func (r *Resilient) UseStrategy(s Strategy) {
	r.Lock()
	defer r.Unlock()
	r.strategies = append(r.strategies, s)
}

func (r *Resilient) Get(url string) (*http.Response, error) {
	o := &client.Options{
		Method: "GET",
		URL:    url,
	}
	return r.Send(o)
}

func (r *Resilient) Post(url string, body io.Reader) (*http.Response, error) {
	o := &client.Options{
		Method: "POST",
		URL:    url,
		Body:   body,
	}
	return r.Send(o)
}

func (r *Resilient) Put(url string, body io.Reader) (*http.Response, error) {
	o := &client.Options{
		Method: "POST",
		URL:    url,
		Body:   body,
	}
	return r.Send(o)
}

func (r *Resilient) Delete(url string, body io.Reader) (*http.Response, error) {
	o := &client.Options{
		Method: "DELETE",
		URL:    url,
		Body:   body,
	}
	return r.Send(o)
}

func (r *Resilient) Patch(url string, body io.Reader) (*http.Response, error) {
	o := &client.Options{
		Method: "PATCH",
		URL:    url,
		Body:   body,
	}
	return r.Send(o)
}

func (r *Resilient) Options(url string, body io.Reader) (*http.Response, error) {
	o := &client.Options{
		Method: "OPTIONS",
		URL:    url,
		Body:   body,
	}
	return r.Send(o)
}

func (r *Resilient) Custom() *client.Options {
	return &client.Options{Method: "GET"}
}

func (r *Resilient) Send(o *client.Options) (*http.Response, error) {
	return NewRequest(r, o)
}
