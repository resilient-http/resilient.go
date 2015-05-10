package resilient

import (
	//"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStrategy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(500)
		w.Write([]byte("error"))
	}))
	defer ts.Close()

	resilient := New()
	resilient.UseServers([]string{ts.URL, ts.URL, ts.URL})

	retries := 3
	resilient.UseStrategy(func() StrategyHandler {
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
