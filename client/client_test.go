package client

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		w.Write([]byte("Hello World"))
	}))
	defer ts.Close()

	opts := &Options{URL: ts.URL}
	client, err := New(opts)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Do()
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf("Invalid status code", res.StatusCode)
	}

	body, _ := ioutil.ReadAll(res.Body)
	if string(body) != "Hello World" {
		t.Fatal("Invalid body response")
	}
}
