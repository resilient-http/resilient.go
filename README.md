# resilient.go [![Build Status](https://travis-ci.org/resilient-http/resilient.go.png)](https://travis-ci.org/resilient-http/resilient.go) [![GoDoc](https://godoc.org/github.com/resilient-http/resilient.go?status.svg)](https://godoc.org/github.com/resilient-http/resilient.go) [![Go Report Card](https://goreportcard.com/badge/github.com/resilient-http/resilient.go)](https://goreportcard.com/report/github.com/resilient-http/resilient.go)

resilient HTTP client for [Go](http://golang.org) built on top of [gentleman](https://github.com/h2non/gentleman) HTTP toolkit.

Designed for distributed systems with fault tolerance capabilities. 

Still beta.

## Features

- Simple, expressive, fuent API.
- Plugin driven architecture.
- Middleware oriented.
- Built on top of `net/http` package.
- Retry strategies supporting exponential or constant back off.
- Service discovery capabitilies via Consul (in progress).

## Installation

```bash
go get -u gopkg.in/resilient-http/resilient.go.v0
```

## Example

```go
package main

import (
  "fmt"
  "gopkg.in/resilient-http/resilient.go.v0"
  "gopkg.in/h2non/gentleman-consul.v0"
)

func main() {
  // Create a new constant retry client
  cli := resilient.NewRetryClient()
  
  // Use Consul for service discovery
  cli.Use(consul.New(consul.NewConfig("demo.consul.io", "web")))

  // Define base URL
  cli.URL("http://httpbin.org")

  // Create a new request based on the current client
  req := cli.Request()

  // Define the URL path at request level
  req.Path("/status/503")

  // Set a new header field
  req.SetHeader("Client", "gentleman")

  // Perform the request
  res, err := req.Send()
  if err != nil {
    fmt.Printf("Request error: %s\n", err)
    return
  }
  if !res.Ok {
    fmt.Printf("Invalid server response: %d\n", res.StatusCode)
    return
  }

  // Print response status and body as string
  fmt.Printf("Status: %d\n", res.StatusCode)
  fmt.Printf("Body: %s\n", res.String())
}
```

#### Custom retry strategy

```go
package main

import (
  "fmt"
  "gopkg.in/resilient-http/resilient.go.v0"
  "gopkg.in/eapache/go-resiliency.v1/retrier"
)

func main() {
  // Create a simple client
  cli := resilient.New()

  // Define base URL
  cli.URL("http://httpbin.org")

  // Register the retry plugin, using a custom exponential retry strategy
  cli.Use(retry.New(retrier.New(retrier.ExponentialBackoff(3, 100*time.Millisecond), nil)))

  // Create a new request based on the current client
  req := cli.Request()

  // Define the URL path at request level
  req.Path("/status/503")

  // Set a new header field
  req.SetHeader("Client", "gentleman")

  // Perform the request
  res, err := req.Send()
  if err != nil {
    fmt.Printf("Request error: %s\n", err)
    return
  }
  if !res.Ok {
    fmt.Printf("Invalid server response: %d\n", res.StatusCode)
    return
  }

  // Print response status and body as string
  fmt.Printf("Status: %d\n", res.StatusCode)
  fmt.Printf("Body: %s\n", res.String())
}
```

## License 

MIT - Tomas Aparicio
