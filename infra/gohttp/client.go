package gohttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
	CloseIdleConnections()
}

type (
	// Client is an interface that defines the behavior of an HTTP client.
	Client interface {
		// DoReq executes a HTTP request.
		DoReq(ctx context.Context, rq *RequestAttributes) (*ResponseAttributes, error)
		// Close is a cleanup function that makes the client close any idle connections.
		Close()
	}

	// httpClient is used to make HTTP requests. It adds additional functionality
	// like automatic retries.
	httpClient struct {
		client *http.Client // Internal HTTP client.
		b      *clientBuilder

		clientOnce sync.Once
	}
)

// nolint:all
func (c *httpClient) DoReq(ctx context.Context, httpRequest *RequestAttributes) (*ResponseAttributes, error) {
	var resp *http.Response

	fullHeaders := c.RequestHeaders(httpRequest.Headers)

	payload, err := c.RequestPayload(fullHeaders.Get(HeaderContentType), httpRequest.Body)
	if err != nil {
		return nil, err
	}

	if c.Logger() != nil && httpRequest.TraceID != "" {
		c.AddTraceIDToLogger(httpRequest.TraceID)
	}

	req, err := http.NewRequestWithContext(ctx, httpRequest.Method, httpRequest.URL, bytes.NewBuffer(payload))
	req.Header = fullHeaders

	var attempt int
	var shouldRetry bool

	var doErr, checkErr error

	for i := 0; i < 3; i++ {
		attempt++

		resp, doErr = c.Client().Do(req)

		switch {
		case ctx.Err() != nil:
			c.client.CloseIdleConnections()
			return nil, fmt.Errorf("http client.Send: %v", ctx.Err())
		default:
		}

		shouldRetry, checkErr = c.b.CheckRetry(req.Context(), resp, doErr)

		if doErr != nil && c.Logger() != nil {
			switch valT := c.Logger().(type) {
			case *slog.Logger:
				valT.Error("failed request")
			case Logger:
				valT.Printf("[ERR] %s %s request failed: %v", req.Method, req.URL, err)
			}
		}

		if !shouldRetry {
			break
		}

		// by default we expect only 1 request [without retry]
		remain := c.RetryMax() - i - 1
		if remain <= 0 {
			break
		}

		wait := c.b.Backoff(c.b.RetryWaitMin, c.b.RetryWaitMax, i, resp)

		if c.Logger() != nil {
			desc := fmt.Sprintf("%s %s", req.Method, req.URL) //nolint:all
			if resp != nil {
				desc = fmt.Sprintf("%s (status after retry: %d)", desc, resp.StatusCode)
			}
			switch valT := c.Logger().(type) {
			case *slog.Logger:
				valT.Debug("retrying request")
			case Logger:
				valT.Printf("[DEBUG] %s: retrying in %s (%d left)", desc, wait, remain)
			}
		}

		timer := time.NewTimer(wait)
		select {
		case <-req.Context().Done():
			timer.Stop()
			c.client.CloseIdleConnections()
			return nil, fmt.Errorf("http client.Send: %v", ctx.Err())
		case <-timer.C:
		}
	}

	defer c.Client().CloseIdleConnections()

	if checkErr != nil {
		return nil, fmt.Errorf("check retry policy: %v", checkErr)
	}

	if doErr != nil {
		return nil, fmt.Errorf("http client.Do: %v", doErr)
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 299 && c.Logger() != nil {
		desc := fmt.Sprintf("%s %s", req.Method, req.URL)
		if resp != nil {
			desc = fmt.Sprintf("%s (status: %d)", desc, resp.StatusCode)
		}

		switch valT := c.Logger().(type) {
		case *slog.Logger:
			valT.Warn(fmt.Sprintf("http client got the response with status code: %d", resp.StatusCode))
		case Logger:
			valT.Printf("[WARN] %s: response: %s with body: %s", desc, resp.StatusCode, string(respBytes))
		}
	}

	return &ResponseAttributes{
		status:     resp.Status,
		statusCode: resp.StatusCode,
		body:       respBytes,
		headers:    resp.Header,
	}, nil
}

func (c *httpClient) Close() {
	if c.client != nil {
		c.client.CloseIdleConnections()
	}
}

func tlsConfig(clb *clientBuilder) *tls.Config {
	if clb == nil || clb.caCert == "" || clb.crtKey == "" {
		return nil
	}

	cert, err := tls.LoadX509KeyPair(clb.caCert, clb.crtKey)
	if err != nil {
		return nil
	}

	caBundle := x509.NewCertPool()
	caBundle.AppendCertsFromPEM([]byte(clb.caBundle))

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caBundle,
	}
}
