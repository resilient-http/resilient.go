package resilient

import (
	"gopkg.in/resilient-http/resilient.go.v0/client"
	"net/http"
	"sync"
)

type StrategyHandler func(*http.Request, *http.Response, error) bool
type Strategy func() StrategyHandler

type Middleware interface {
	Out()
	In()
}

type Resilient struct {
	sync.Mutex
	Servers          []string
	DiscoveryServers []string
	strategies       []Strategy
	middlewares      []Middleware
}

func New() *Resilient {
	return &Resilient{}
}

func (r *Resilient) UseServers(servers []string) *Resilient {
	r.Lock()
	defer r.Unlock()
	r.Servers = servers
	return r
}

func (r *Resilient) UseDiscoveryServers(servers []string) *Resilient {
	r.Lock()
	defer r.Unlock()
	r.DiscoveryServers = servers
	return r
}

func (r *Resilient) Use(m Middleware) {
	r.Lock()
	defer r.Unlock()
	r.middlewares = append(r.middlewares, m)
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

func (r *Resilient) Send(o *client.Options) (*http.Response, error) {
	return NewRequest(r, o)
}
