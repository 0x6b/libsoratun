package libsoratun

import (
	"C"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	UnifiedEndpointHostname = "100.127.69.42"
	UnifiedEndpointPort     = 80
)

// UnifiedEndpointHTTPClient is an HTTP client that can be used to communicate with SORACOM Unified Endpoint.
type UnifiedEndpointHTTPClient struct {
	httpClient *http.Client
	endpoint   *url.URL
	headers    []string
}

// Params is a set of parameters for HTTP request.
type Params struct {
	// Method is an HTTP method of the request. Only GET or POST is supported.
	Method string
	// Path is a path of the request.
	Path string
	// Body is a body of the request.
	Body string
}

// NewUnifiedEndpointHTTPClient creates a new HTTP client specific for the SORACOM Unified Endpoint.
//
// Args:
// - `config`: A Config object containing the configuration for the SORACOM Arc connection.
//
// Returns:
// - `*UnifiedEndpointHTTPClient`: A pointer to the created UnifiedEndpointHTTPClient and an error object. If there's any error occurred during setting up the tunnel connection or parsing the endpoint URL, the error will be returned.
//
// The created client uses the http.Transport with a custom DialContext for making HTTP requests. The User-Agent header for all requests made by this client is set to "libsoratun/0.0.1".
func NewUnifiedEndpointHTTPClient(config Config) (*UnifiedEndpointHTTPClient, error) {
	t, err := newTunnel(&config)
	if err != nil {
		return nil, err
	}

	endpoint, err := url.Parse(fmt.Sprintf("http://%s:%d", UnifiedEndpointHostname, UnifiedEndpointPort))
	if err != nil {
		return nil, err
	}

	return &UnifiedEndpointHTTPClient{
		httpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: t.DialContext,
			},
		},
		endpoint: endpoint,
		headers:  []string{"User-Agent: libsoratun/0.0.1"},
	}, nil
}

// MakeRequest creates a new HTTP request based on provided parameters.
//
// Args:
// - `params`: A Params struct containing the parameters for the request.
//   - `params.Method`: The HTTP method of the request. Only GET or POST is supported. Any extra spaces are trimmed.
//   - `params.Path`: The path of the request. It cannot be empty, and any leading slashes are removed.
//   - `params.Body`: The body of the request. It cannot be empty.
//
// Returns:
// - `*http.Request`: A pointer to the created http.Request, and an error object.
func (c *UnifiedEndpointHTTPClient) MakeRequest(params *Params) (*http.Request, error) {
	method := strings.TrimSpace(params.Method)
	if !(method == http.MethodGet || method == http.MethodPost) {
		return nil, fmt.Errorf("only GET or POST is supported")
	}

	if params.Path == "" {
		return nil, fmt.Errorf("path is required")
	}
	path := strings.TrimPrefix(params.Path, "/")

	if params.Body == "" {
		return nil, fmt.Errorf("body is required")
	}
	body := strings.NewReader(params.Body)

	req, err := http.NewRequest(params.Method, fmt.Sprintf("%s/%s", c.endpoint, path), body)
	if err != nil {
		return nil, err
	}

	for _, h := range c.headers {
		header := strings.SplitN(h, ":", 2)
		if len(header) == 2 {
			req.Header.Set(strings.TrimSpace(header[0]), strings.TrimSpace(header[1]))
		}
	}

	return req, nil
}

// DoRequest sends an HTTP request and returns the response.
//
// Args:
// - `req`: A pointer to the http.Request object to be sent.
//
// Returns:
// - `*http.Response`: A pointer to the http.Response, and an error object.
func (c *UnifiedEndpointHTTPClient) DoRequest(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		defer func() {
			err := res.Body.Close()
			if err != nil {
				fmt.Println("failed to close response", err)
			}
		}()
		r, _ := io.ReadAll(res.Body)
		return res, fmt.Errorf("%s: %s %s: %s", res.Status, req.Method, req.URL, r)
	}

	return res, nil
}
