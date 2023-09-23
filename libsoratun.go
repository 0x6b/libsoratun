package main

import (
	"C"
	"fmt"
	"io"

	"github.com/0x6b/libsoratun/libsoratun"
)

func main() {}

// Send is a function that uses SORACOM Arc configuration provided as a JSON string,
// an HTTP method, a path, and a payload to create and execute an HTTP request to the SORACOM Unified Endpoint.
// The function returns the response as a C language string, or nil if an error occurs.
// Input parameters are expected to be C language strings.
//
// Parameters:
//   - `configJson`: C string representation of SORACOM Arc configuration
//   - `method`: HTTP method to be used in the request. Only GET and POST are supported.
//   - `path`: path to be used in the HTTP request
//   - `body`: body of the HTTP request
//
// Returns:
//   - C string representation of the HTTP response body, or nil if an error occurred
//
// Usage:
//
//	response := Send(configJson, method, path, body)
//
//export Send
func Send(configJson, method, path, body *C.char) *C.char {
	config, err := libsoratun.NewConfig([]byte(C.GoString(configJson)))
	if err != nil {
		return nil
	}

	c, err := libsoratun.NewUnifiedEndpointHTTPClient(*config)
	if err != nil {
		return nil
	}

	req, err := c.MakeRequest(&libsoratun.Params{
		Method: C.GoString(method),
		Path:   C.GoString(path),
		Body:   C.GoString(body),
	})
	if err != nil {
		return nil
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil
	}

	response, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("failed to read response from Unified Endpoint: %v\n", err)
		panic(err)
	}

	return C.CString(string(response))
}
