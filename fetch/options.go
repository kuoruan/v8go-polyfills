package fetch

import (
	"net/http"

	"go.kuoruan.net/v8go-polyfills/internal"
)

type Option interface {
	internal.Option

	apply(ft *fetcher)
}

type funcOption func(ft *fetcher)

func (f funcOption) apply(ft *fetcher) {
	f(ft)
}

func (f funcOption) Polyfill() internal.Polyfill {
	return internal.PolyfillFetch
}

func WithLocalHandler(handler http.Handler) Option {
	return funcOption(func(ft *fetcher) {
		ft.LocalHandler = handler
	})
}
