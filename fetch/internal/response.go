package internal

import (
	"io/ioutil"
	"net/http"
)

/*
 Response keeps the *http.Response
*/
type Response struct {
	Headers    map[string]string `json:"headers"`
	Status     int32             `json:"status"`
	StatusText string            `json:"statusText"`
	OK         bool              `json:"ok"`
	Redirected bool              `json:"redirected"`
	URL        string            `json:"url"`
	Body       string            `json:"body"`
}

/*
 Handle the *http.Response, return *Response
*/
func HandleHttpResponse(res *http.Response, url string, redirected bool) (*Response, error) {
	// convert the http.Header to map
	resHeaders := make(map[string]string)
	for k, v := range res.Header {
		for _, vv := range v {
			resHeaders[k] = vv
			break
		}
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Headers:    resHeaders,
		Status:     int32(res.StatusCode), // int type is not support by v8go
		StatusText: res.Status,
		OK:         res.StatusCode >= 200 && res.StatusCode < 300,
		Redirected: redirected,
		URL:        url,
		Body:       string(resBody),
	}, nil
}
