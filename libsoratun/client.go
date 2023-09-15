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
	Body io.Reader
}

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

func (c *UnifiedEndpointHTTPClient) MakeRequest(params *Params) (*http.Request, error) {
	if !(params.Method == http.MethodGet || params.Method == http.MethodPost) {
		return nil, fmt.Errorf("only GET or POST is supported")
	}

	req, err := http.NewRequest(
		params.Method,
		fmt.Sprintf("%s://%s:%s/%s",
			c.endpoint.Scheme,
			c.endpoint.Hostname(),
			c.endpoint.Port(),
			params.Path,
		),
		params.Body,
	)
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
