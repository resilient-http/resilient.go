# resilient.go [![Build Status](https://travis-ci.org/resilient-http/resilient.go.png)](https://travis-ci.org/resilient-http/resilient.go) [![GitHub release](http://img.shields.io/github/tag/resilient-http/resilient.go.svg?style=flat-square)](https://github.com/resilient-http/resilient.go/releases) [![GoDoc](https://godoc.org/github.com/resilient-http/resilient.go?status.svg)](https://godoc.org/github.com/resilient-http/resilient.go) [![Coverage Status](https://coveralls.io/repos/resilient-http/resilient.go/badge.svg?branch=master)](https://coveralls.io/r/resilient-http/resilient.go?branch=master)

Middleware-oriented and plugable HTTP client with superpowers designed 
for distributed and [reactive](http://www.reactivemanifesto.org/) systems

This can be considered as a free-style port in Go based on resilient.js, 
but focused on a better design, extensibility and simplicity.

**Work in progress**

## Superpowers

- Full featured HTTP client build on top of `net/http`
- No third-party dependencies (only native packages are used)
- Smart fail over (via multiple strategies)
- Dynamic server discovery (via middleware)
- Retries cycles 
- Client side based balancing and load distribution (via)

## Installation

```
go get gopkg.in/resilient-http/resilient.go.v0
```

## Addons

### Strategies

- Failover

### Middlewares

- Server discovery

## Rationale

Distributed and reactive system architectures are growing during last years, but they are mostly hard to support since most of the complexibility is delegated in the consumer side.

Resilient aims to simplify this to deal elegantly with distributed systems

## License 

MIT - Tomas Aparicio
