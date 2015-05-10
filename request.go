package resilient

import (
	"gopkg.in/resilient-http/resilient.go.v0/client"
	"net/http"
	"net/url"
)

func NewRequest(r *Resilient, o *client.Options) (*http.Response, error) {
	c, err := client.New(o)
	if err != nil {
		return nil, err
	}

	servers := make([]string, len(r.Servers))
	copy(servers, r.Servers)

	return send(servers, r, c)
}

func send(servers []string, r *Resilient, c *client.Client) (*http.Response, error) {
	req := c.Request

	if req.URL.Host == "" && len(servers) > 0 {
		serverURL, err := url.Parse(servers[0])
		if err != nil {
			return nil, err
		}
		req.URL.Host = serverURL.Host
		req.URL.User = serverURL.User
		req.URL.Scheme = serverURL.Scheme
		req.URL.Opaque = serverURL.Opaque
		req.URL.Path = serverURL.Path + req.URL.Path
	}

	res, err := c.Do()

	var retry bool
	for _, strategy := range r.strategies {
		if strategy(req, res, err) {
			retry = true
		}
	}

	if retry {
		if len(servers) > 1 {
			servers = servers[1:]
		} else {
			servers = servers[0:]
		}
		return send(servers, r, c)
	}

	return res, err
}
