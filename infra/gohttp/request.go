package gohttp

import (
	"net/http"
)

// RequestAttributes represents the attributes of an HTTP request that can be used by a client to send a request.
type RequestAttributes struct {
	Method  string
	URL     string
	Headers http.Header
	Body    interface{}
	Retry   bool
	TraceID string
}

// HttpClientConfig represents the configuration options for an HTTP client that uses TLS (Transport Layer Security) to secure its connections.
type HttpClientConfig struct {
	CACert   string
	CrtKey   string
	CaBundle string
	Timeout  int
	Audience string
}
