package custom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	neturl "net/url"
	"time"
)

type Client struct {
	RequestParams RequestParams
	Client        *http.Client
}

type RequestParams struct {
	Endpoint string
	Params   map[string]string
	Headers  map[string]string
	Body     string
	Method   string
	Timeout  time.Duration
}

func NewClient() *Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second, // Connection timeout
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second, // TLS handshake timeout
	}

	return &Client{
		Client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}
}

func (s *Client) Do() (*http.Response, error) {
	if s.RequestParams.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required")
	}

	// Build URL with query parameters
	url := s.RequestParams.Endpoint
	if len(s.RequestParams.Params) > 0 {
		url += "?"
		first := true
		for key, value := range s.RequestParams.Params {
			if !first {
				url += "&"
			}

			url += key + "=" + neturl.QueryEscape(value)
			first = false
		}
	}

	// Create request body
	var body io.Reader
	if s.RequestParams.Body != "" {
		body = bytes.NewBufferString(s.RequestParams.Body)
	}

	// Create HTTP request
	req, err := http.NewRequest(s.RequestParams.Method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range s.RequestParams.Headers {
		req.Header.Set(key, value)
	}

	// Set timeout if specified
	ctx := req.Context()
	if s.RequestParams.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.RequestParams.Timeout)
		defer cancel()
		req = req.WithContext(ctx)
	}

	// Execute request
	return s.Client.Do(req)
}

// Helper methods for common HTTP operations
func (s *Client) Get(endpoint string, params map[string]string, headers map[string]string) (*http.Response, error) {
	s.RequestParams = RequestParams{
		Method:   "GET",
		Endpoint: endpoint,
		Params:   params,
		Headers:  headers,
	}
	return s.Do()
}

func (s *Client) Post(endpoint string, body string, headers map[string]string) (*http.Response, error) {
	s.RequestParams = RequestParams{
		Method:   "POST",
		Endpoint: endpoint,
		Body:     body,
		Headers:  headers,
	}
	return s.Do()
}

func (s *Client) Put(endpoint string, body string, headers map[string]string) (*http.Response, error) {
	s.RequestParams = RequestParams{
		Method:   "PUT",
		Endpoint: endpoint,
		Body:     body,
		Headers:  headers,
	}
	return s.Do()
}

func (s *Client) Delete(endpoint string, headers map[string]string) (*http.Response, error) {
	s.RequestParams = RequestParams{
		Method:   "DELETE",
		Endpoint: endpoint,
		Headers:  headers,
	}
	return s.Do()
}

// JSON helper methods
func (s *Client) PostJSON(endpoint string, data interface{}, headers map[string]string) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"

	return s.Post(endpoint, string(jsonData), headers)
}

func (s *Client) PutJSON(endpoint string, data interface{}, headers map[string]string) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"

	return s.Put(endpoint, string(jsonData), headers)
}

// Set timeout for the client
func (s *Client) SetTimeout(timeout time.Duration) {
	s.Client.Timeout = timeout
}

// Set custom HTTP client
func (s *Client) SetClient(client *http.Client) {
	s.Client = client
}
