package console

import (
	"fmt"
	"io"
	"os"

	"rogchap.com/v8go"
)

type Console interface {
	GetLogFunctionCallback() v8go.FunctionCallback
}

type console struct {
	Output io.Writer
}

func NewConsole(opt ...Option) Console {
	c := &console{
		Output: os.Stdout,
	}

	for _, o := range opt {
		o.apply(c)
	}

	return c
}

func (c *console) GetLogFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		if args := info.Args(); len(args) > 0 {
			inputs := make([]interface{}, len(args))
			for i, input := range args {
				inputs[i] = input
			}

			fmt.Fprintln(c.Output, inputs...)
		}

		return nil
	}
}
