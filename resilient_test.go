package resilient

import (
	//"io/ioutil"
	"gopkg.in/resilient-http/resilient.go.v0/client"
	"gopkg.in/resilient-http/resilient.go.v0/middlewares"
	"gopkg.in/resilient-http/resilient.go.v0/strategies"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStrategy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("error"))
	}))
	defer ts.Close()

	resilient := New()
	resilient.UseServers([]string{ts.URL, ts.URL, ts.URL})

	retries := 3
	resilient.UseStrategy(func() strategies.Handler {
		return func(req *http.Request, res *http.Response, err error) bool {
			if retries == 0 {
				return false
			}
			retries -= 1
			return true
		}
	})

	res, err := resilient.Get("/foo")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 500 {
		t.Fatalf("Invalid status code", res.StatusCode)
	}

	if retries != 0 {
		t.Fatal("Invalid number of retries attempts")
	}
}

type middleware struct {
	resilient middlewares.Resilient
	servers   []string
}

func (m *middleware) Out(o *client.Options) error {
	m.resilient.UseServers(m.servers)
	return nil
}

func (m *middleware) In(*http.Request, *http.Response, error) error {
	return nil
}

func TestMiddleware(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("error"))
	}))
	defer ts.Close()

	resilient := New()

	resilient.Use(func(r middlewares.Resilient) (middlewares.Middleware, error) {
		return &middleware{servers: []string{ts.URL}, resilient: r}, nil
	})

	res, err := resilient.Get("/foo")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 500 {
		t.Fatalf("Invalid status code", res.StatusCode)
	}
}
