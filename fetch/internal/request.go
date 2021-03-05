package internal

import "net/http"

type RequestRedirect string

const (
	RequestRedirectError  RequestRedirect = "error"
	RequestRedirectFollow                 = "follow"
	RequestRedirectManual                 = "manual"
)

type RequestInit struct {
	Body     string            `json:"body"`
	Headers  map[string]string `json:"headers"`
	Method   string            `json:"method"`
	Redirect RequestRedirect   `json:"redirect"`
}

type Request struct {
	RequestInit

	Headers http.Header `json:"headers"`
	URL     string      `json:"url"`
	IsLocal bool        `json:"-"`
}
