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
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	RequestRedirectError  = "error"
	RequestRedirectFollow = "follow"
	RequestRedirectManual = "manual"
)

/*
 RequestInit is the fetch API defined object.
 Only supports raw request now.
*/
type RequestInit struct {
	Body     string            `json:"body"`
	Headers  map[string]string `json:"headers"`
	Method   string            `json:"method"`
	Redirect string            `json:"redirect"`
}

/*
 Request is the request object used by fetch
*/
type Request struct {
	Body     string
	Method   string
	Redirect string

	Header     http.Header
	URL        *url.URL
	RemoteAddr string
}

/*
 parse and check the request URL, return *url.URL
*/
func ParseRequestURL(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("url '%s' is not valid, %w", rawURL, err)
	}

	/**
	 * Check the scheme, we only support http and https at this time
	 */
	switch u.Scheme {
	case "http", "https":
	case "": // then scheme is empty, it's a local request
		if !strings.HasPrefix(u.Path, "/") {
			return nil, fmt.Errorf("unsupported relatve path %s", u.Path)
		}
	default:
		return nil, fmt.Errorf("unsupported scheme %s", u.Scheme)
	}

	return u, nil
}
