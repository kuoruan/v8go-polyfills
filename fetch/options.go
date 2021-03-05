package fetch

import (
	"net/http"

	polyfills "go.kuoruan.net/v8go-polyfills"
)

type Option interface {
	polyfills.Option

	apply(ft *fetcher)
}

type funcOption func(ft *fetcher)

func (f funcOption) apply(ft *fetcher) {
	f(ft)
}

func (f funcOption) ForPolyfill() polyfills.Polyfill {
	return polyfills.PolyfillFetch
}

func WithLocalHandler(handler http.Handler) Option {
	return funcOption(func(ft *fetcher) {
		ft.LocalHandler = handler
	})
}
