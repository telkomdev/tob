package httpx

import (
	"errors"
	"io"
	"net/http"
)

var (
	// ErrorStatusNot200 error type
	ErrorStatusNot200 = errors.New("error: HTTP Response code not 200")
)

// HTTPRequest will executes HTTP request and returns the response body.
// Any errors or non-200 status code result in an error.
func HTTPRequest(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// HTTPPost will execute HTTP POST with json payload
func HTTPPost(url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	return HTTPRequest(req)
}

// HTTPGet will execute HTTP GET
func HTTPGet(url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	return HTTPRequest(req)
}
