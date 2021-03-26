package console

import (
	"io"
)

type Option interface {
	apply(c *console)
}

type funcOption func(c *console)

func (f funcOption) apply(c *console) {
	f(c)
}

func (f funcOption) Polyfill() string {
	return "console"
}

func WithOutput(output io.Writer) Option {
	return funcOption(func(c *console) {
		c.Output = output
	})
}
