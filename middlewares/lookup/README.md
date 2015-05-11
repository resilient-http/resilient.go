# Lookup middleware

Dynamic servers lookup middleware for Resilient.go

## Features

- Dynamic servers lookup
- Automatic servers refresh (via goroutine)

## Adaptors

- Consul

## Installation

```
go get gopkg.in/resilient-http/resilient.go.v0/middlewares/lookup
```

## Usage

```go
package main

import (
  "time"
  "fmt"
  "gopkg.in/resilient-http/resilient.go.v0"
  "gopkg.in/resilient-http/resilient.go.v0/middlewares/lookup"
)

func main() {
  r := resilient.New()
  
  r.Use(lookup.New(lookup.Options{
    Servers: []string{"http://foo"},
    Refresh: 100 * time.Millisecond,
  }))

  res, err := r.Get("/foo")
  if err != nil {
    fmt.Printf("Error: %#v", err)
  }
  fmt.Printf("Status: %d", res.StatusCode)
}
```
