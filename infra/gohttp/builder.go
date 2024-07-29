package gohttp

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	// Default retry configuration.
	defaultRetryWaitMin = 1 * time.Second
	defaultRetryWaitMax = 30 * time.Second
	defaultRetryMax     = 1

	// defaultLogger is the logger provided with defaultClient.
	defaultLogger = log.New(os.Stderr, "", log.LstdFlags) //nolint:all
)

// The ClientBuilder interface defines a set of methods
// that can be used to configure an HTTP client for making HTTP requests.
// These methods allow the user to set various properties of the client
// such as the HTTP headers, the user agent string,
// the connection timeout, the response header timeout,
// and the TLS configuration.
type ClientBuilder interface {
	// DisableTimeouts can be used to disable or enable timeouts for the HTTP client.
	// If disable is set to true, the client will not time out on any requests.
	DisableTimeouts(disable bool) *clientBuilder

	// SetHTTPClient sets the underlying HTTP client that will be used to make HTTP requests.
	// The user can pass in their own custom HTTP client with specific configurations.
	SetHTTPClient(client *http.Client) *clientBuilder

	// SetHeaders sets the HTTP headers that will be sent with every HTTP request made by the client.
	// The headers are provided as an http.Header object.
	SetHeaders(headers http.Header) *clientBuilder

	// SetUserAgent sets the User-Agent header that will be sent with every HTTP request made by the client.
	// The user can provide their own custom user agent string.
	SetUserAgent(userAgent string) *clientBuilder

	// SetConnectionTimeOut sets the connection timeout for the HTTP client.
	// If the client is unable to establish a connection within the specified duration, it will return an error.
	SetConnectionTimeOut(timeout time.Duration) *clientBuilder

	// This method sets the response header timeout for the HTTP client.
	// If the server takes longer than the specified duration to send the response headers,
	// the client will return an error.
	SetResponseHeaderTimeout(timeout time.Duration) *clientBuilder

	// SetTLSVariables sets the TLS configuration for the HTTP client.
	// The user can provide their own custom TLS configuration by passing in an HttpClientConfig object.
	SetTLSVariables(conf HttpClientConfig) *clientBuilder

	// SetLogger sets the logger if it can be important.
	SetLogger(l Logger) *clientBuilder
}

// Backoff specifies a policy for how long to wait between retries.
// It is called after a failing request to determine the amount of time
// that should pass before trying again.
type Backoff func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration

// CheckRetry specifies a policy for handling retries. It is called
// following each request with the response and error values returned by
// the http.Client. If CheckRetry returns false, the Client stops retrying
// and returns the response to the caller. If CheckRetry returns an error,
// that error value is returned in lieu of the error from the request. The
// Client will close any response body when retrying, but if the retry is
// aborted it is up to the CheckRetry callback to properly close any
// response body before returning.
type CheckRetry func(ctx context.Context, resp *http.Response, err error) (bool, error)

type clientBuilder struct {
	client  *http.Client
	baseUrl string //nolint:all

	headers http.Header

	connectionTimeout     time.Duration
	responseHeaderTimeOut time.Duration
	disableTimeouts       bool

	caCert   string
	crtKey   string
	caBundle string

	audience string

	RetryWaitMin time.Duration // Minimum time to wait
	RetryWaitMax time.Duration // Maximum time to wait
	RetryMax     int           // Maximum number of retries

	Logger Logger // Customer logger instance. Can be either Logger or LeveledLogger

	// Backoff specifies the policy for how long to wait between retries
	Backoff Backoff

	// CheckRetry specifies the policy for handling retries, and is called
	// after each request. The default policy is DefaultRetryPolicy.
	CheckRetry CheckRetry
}

// NewClientBuilder returns a new instance of the clientBuilder struct, which can be used to configure and create HTTP clients.
// The clientBuilder struct provides several methods for configuring the client, including setting the base URL, adding default headers,
// and setting the HTTP transport. Once the client has been configured, call the Build method to create an HTTP client instance.
func NewClientBuilder() *clientBuilder {
	return &clientBuilder{
		RetryWaitMin: defaultRetryWaitMin,
		RetryWaitMax: defaultRetryWaitMax,
		RetryMax:     defaultRetryMax,
		Backoff:      DefaultBackoff,
		CheckRetry:   DefaultRetryPolicy,
	}
}

// Build creates a new httpClient instance with the current configuration settings and returns it.
func (b *clientBuilder) Build() *httpClient {
	return &httpClient{
		b: b,
	}
}

// SetHTTPClient sets the HTTP client to use for making requests.
func (b *clientBuilder) SetHTTPClient(client *http.Client) *clientBuilder {
	b.client = client

	return b
}

// SetRetryWaitMin sets the RetryWaitMin to the builder.
func (b *clientBuilder) SetRetryWaitMin(arg time.Duration) *clientBuilder {
	b.RetryWaitMin = arg

	return b
}

// SetRetryWaitMax sets the RetryWaitMax to the builder.
func (b *clientBuilder) SetRetryWaitMax(arg time.Duration) *clientBuilder {
	b.RetryWaitMax = arg

	return b
}

// SetRetryMax sets the RetryMax to the builder.
func (b *clientBuilder) SetRetryMax(arg int) *clientBuilder {
	b.RetryMax = arg

	return b
}

// SetAudience sets the audience value to use for generating ID tokens.
func (b *clientBuilder) SetAudience(audience string) *clientBuilder {
	b.audience = audience

	return b
}

// SetHeaders sets the default HTTP headers to include with every request.
func (b *clientBuilder) SetHeaders(headers http.Header) *clientBuilder {
	if b.headers == nil {
		b.headers = make(http.Header, 0)
	}

	for k, v := range headers {
		for _, dv := range v {
			b.headers.Add(k, dv)
		}
	}

	return b
}

// SetConnectionTimeOut sets the timeout for establishing a connection to the server.
func (b *clientBuilder) SetConnectionTimeOut(timeout time.Duration) *clientBuilder {
	b.connectionTimeout = timeout

	return b
}

// SetConnectionTimeOut sets the timeout for establishing a connection to the server.
func (b *clientBuilder) SetLogger(l Logger) *clientBuilder {
	b.Logger = l

	return b
}

// SetResponseHeaderTimeout sets the timeout for receiving response headers from the server.
func (b *clientBuilder) SetResponseHeaderTimeout(timeout time.Duration) *clientBuilder {
	b.responseHeaderTimeOut = timeout

	return b
}

// DisableTimeouts sets whether to disable timeouts for requests.
func (b *clientBuilder) DisableTimeouts(disable bool) *clientBuilder {
	b.disableTimeouts = disable

	return b
}

// SetTLSVariables sets the TLS variables for the HTTP client.
func (b *clientBuilder) SetTLSVariables(conf *HttpClientConfig) *clientBuilder {
	b.caCert = conf.CACert
	b.crtKey = conf.CrtKey
	b.caBundle = conf.CaBundle

	return b
}

// SetUserAgent sets the User-Agent header to include with every request.
func (b *clientBuilder) SetUserAgent(userAgent string) *clientBuilder {
	if val := b.headers.Get(HeaderUserAgent); val != "" {
		return b
	}

	b.headers.Set(HeaderUserAgent, userAgent)
	return b
}
