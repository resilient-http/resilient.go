package resilient

import (
	"gopkg.in/resilient-http/resilient.go.v0/client"
	"net/http"
)

type Strategy func(*http.Request, *http.Response, error) bool
type Middleware func(params interface{}) bool

type Resilient struct {
	Servers          []string
	DiscoveryServers []string
	strategies       []Strategy
	middlewares      []Middleware
}

func New() *Resilient {
	return &Resilient{}
}

func (r *Resilient) UseServers(servers []string) *Resilient {
	r.Servers = servers
	return r
}

func (r *Resilient) UseDiscoveryServers(servers []string) *Resilient {
	r.DiscoveryServers = servers
	return r
}

func (r *Resilient) Use(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Resilient) UseStrategy(m Strategy) {
	r.strategies = append(r.strategies, m)
}

func (r *Resilient) Get(url string) (*http.Response, error) {
	o := &client.Options{
		Method: "GET",
		URL:    url,
	}

	return r.Send(o)
}

func (r *Resilient) Send(o *client.Options) (*http.Response, error) {
	return NewRequest(r, o)
}
