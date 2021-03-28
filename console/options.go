package console

import (
	"io"
)

type Option interface {
	apply(c *console)
}

type optionFunc func(c *console)

func (f optionFunc) apply(c *console) {
	f(c)
}

func WithOutput(output io.Writer) Option {
	return optionFunc(func(c *console) {
		c.Output = output
	})
}
