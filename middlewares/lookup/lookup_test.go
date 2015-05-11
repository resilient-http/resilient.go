package lookup

import (
	"gopkg.in/resilient-http/resilient.go.v0"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicLookup(t *testing.T) {
	client := resilient.New()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello World"))
	}))
	defer ts.Close()

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`["` + ts.URL + `"]`))
	}))
	defer ts2.Close()

	client.Use(New(Options{
		Servers: []string{ts2.URL},
	}))

	res, err := client.Get("/foo")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatal("Invalid status code")
	}

}
