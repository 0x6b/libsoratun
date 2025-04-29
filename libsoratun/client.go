package libsoratun

import (
	"C"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"golang.zx2c4.com/wireguard/device"
)
import (
	"context"
	"net"
)

var Revision = "dev"

const (
	UnifiedEndpointHostname = "100.127.69.42"
	UnifiedEndpointPort     = 80
)

// UnifiedEndpointHTTPClient is an HTTP client that can be used to communicate with SORACOM Unified Endpoint.
type UnifiedEndpointHTTPClient struct {
	httpClient *http.Client
	endpoint   *url.URL
	headers    []string
	logger     *device.Logger
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

// UnifiedEndpointUDPClient is an UDP client that can be used to communicate with SORACOM Unified Endpoint.
type UnifiedEndpointUDPClient struct {
	dialcontext func(ctx context.Context, network string, addr string) (net.Conn, error)
	logger      *device.Logger
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
	logger := device.NewLogger(config.LogLevel, "(libsoratun/client) ") // lazily use device.Logger for logging
	t, err := newTunnel(&config)
	if err != nil {
		return nil, err
	}

	endpoint, err := url.Parse(fmt.Sprintf("http://%s:%d", UnifiedEndpointHostname, UnifiedEndpointPort))
	if err != nil {
		return nil, err
	}

	ua := fmt.Sprintf("User-Agent: libsoratun/%s", Revision)
	logger.Verbosef("Soracom Unified Endpoint URL: %s", endpoint)
	logger.Verbosef("%s", ua)

	return &UnifiedEndpointHTTPClient{
		httpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: t.DialContext,
			},
		},
		endpoint: endpoint,
		headers:  []string{ua},
		logger:   logger,
	}, nil
}

// NewUnifiedEndpointUDPClient creates a new UDP client specific for the SORACOM Unified Endpoint.
//
// Args:
// - `config`: A Config object containing the configuration for the SORACOM Arc connection.
//
// Returns:
// - `*UnifiedEndpointUDPClient`: A pointer to the created UnifiedEndpointUDPClient and an error object. If there's any error occurred during setting up the tunnel connection or parsing the endpoint URL, the error will be returned.
//
// The created client uses the udp.
func NewUnifiedEndpointUDPClient(config Config) (*UnifiedEndpointUDPClient, error) {
	logger := device.NewLogger(config.LogLevel, "(libsoratun/client) ") // lazily use device.Logger for logging
	t, err := newTunnel(&config)
	if err != nil {
		return nil, err
	}

	endpoint, err := url.Parse(fmt.Sprintf("http://%s:%d", UnifiedEndpointHostname, UnifiedEndpointPort))
	if err != nil {
		return nil, err
	}

	ua := fmt.Sprintf("User-Agent: libsoratun/%s", Revision)
	logger.Verbosef("Soracom Unified Endpoint URL: %s", endpoint)
	logger.Verbosef("%s", ua)

	return &UnifiedEndpointUDPClient{
		dialcontext: t.DialContext,
		logger:      logger,
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
	method := strings.ToUpper(strings.TrimSpace(params.Method))
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

	r, err := httputil.DumpRequest(req, true)
	if err == nil {
		c.logger.Verbosef("Sent HTTP request:\n%s", r)
	} else {
		c.logger.Errorf("Failed to dump HTTP request", err)
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

	r, err := httputil.DumpResponse(res, true)
	if err == nil {
		c.logger.Verbosef("Received HTTP response:\n%s", r)
	} else {
		c.logger.Errorf("Failed to dump HTTP response", err)
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

// DoUDPRequest sends an UDP request and returns the response.
//
// Args:
// - `Body`: A request body
// - `port`: A string representing the address to send the request to.
//
// Returns:
// - `string`: A response string, and an error object.
func (c *UnifiedEndpointUDPClient) DoUDPRequest(body []byte, port int16) (string, error) {
	addr := fmt.Sprintf("%s:%d", UnifiedEndpointHostname, port)
	conn, err := c.dialcontext(context.Background(), "udp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	len, err := conn.Write([]byte(body))
	if err != nil {
		c.logger.Errorf("Failed to Write packet", err)
		return "", err
	}
	c.logger.Verbosef("UDP sent %d bytes", len)

	res := make([]byte, 1024)

	len, err = conn.Read(res)
	if err != nil {
		return "", err
	}
	c.logger.Verbosef("UDP received %d bytes", len)
	return string(res), nil
}
