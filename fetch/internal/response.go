/*
 * Copyright (c) 2021 Xingwang Liao
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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
