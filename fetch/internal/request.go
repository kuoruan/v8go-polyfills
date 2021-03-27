package internal

import (
	"net/http"
	"net/url"
)

const (
	RequestRedirectError  = "error"
	RequestRedirectFollow = "follow"
	RequestRedirectManual = "manual"
)

type RequestInit struct {
	Body     string            `json:"body"`
	Headers  map[string]string `json:"headers"`
	Method   string            `json:"method"`
	Redirect string            `json:"redirect"`
}

type Request struct {
	Body     string
	Method   string
	Redirect string

	Header     http.Header
	URL        *url.URL
	RemoteAddr string
}
