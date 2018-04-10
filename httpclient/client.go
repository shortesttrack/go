package httpclient

import (
	"net/http"
	"bytes"
	"net/url"
)

var client = &http.Client{}

func GETJSON(url string, params map[string]string, headers map[string]string) (*Response, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	return GET(url, params, headers)
}

func GET(urlString string, params map[string]string, headers map[string]string) (*Response, error) {
	payload := url.Values{}
	for k, v := range params {
		payload.Set(k, v)
	}
	if params != nil && len(params) > 0 {
		urlString = urlString + "?" + payload.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, urlString, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return &Response{res}, err
}

func POSTJSON(url string, headers map[string]string, body []byte) (*Response, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	return POST(url, headers, body)
}

func POST(url string, headers map[string]string, body []byte) (*Response, error) {
	var buf bytes.Buffer
	buf.Write(body)

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)
	return &Response{res}, err
}
