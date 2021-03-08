package console

import (
	"io"

	"go.kuoruan.net/v8go-polyfills/internal"
)

type Option interface {
	internal.Option

	apply(c *console)
}

type funcOption func(c *console)

func (f funcOption) apply(c *console) {
	f(c)
}

func (f funcOption) Polyfill() internal.Polyfill {
	return internal.PolyfillConsole
}

func WithOutput(output io.Writer) Option {
	return funcOption(func(c *console) {
		c.Output = output
	})
}
