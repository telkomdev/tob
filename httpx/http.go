package httpx

import (
	"errors"
	"io"
	"net/http"
	"time"
)

var (
	// ErrorStatusNot200 error type
	ErrorStatusNot200 = errors.New("error: HTTP Response code not 200")
)

// HTTPRequest will executes HTTP request and returns the response body.
// Any errors or non-200 status code result in an error.
func HTTPRequest(req *http.Request, timeout int) (*http.Response, error) {
	// set default http client timeout to 5 seconds
	if timeout <= 0 {
		timeout = 5
	}

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// HTTPPost will execute HTTP POST
func HTTPPost(url string, body io.Reader, headers map[string]string, timeout int) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	return HTTPRequest(req, timeout)
}

// HTTPGet will execute HTTP GET
func HTTPGet(url string, headers map[string]string, timeout int) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	return HTTPRequest(req, timeout)
}
