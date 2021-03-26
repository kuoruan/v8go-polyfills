package console

import (
	"fmt"
	"io"
	"os"

	"rogchap.com/v8go"
)

type console struct {
	Output io.Writer
}

/**
Inject basic console.log support.
*/
func Inject(ctx *v8go.Context, opt ...Option) error {
	c := console{Output: os.Stdout}

	for _, o := range opt {
		o.apply(&c)
	}

	iso, _ := ctx.Isolate()
	global := ctx.Global()

	console, _ := v8go.NewObjectTemplate(iso)
	logFn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) > 0 {
			inputs := make([]interface{}, len(args))
			for i, input := range args {
				inputs[i] = input
			}

			fmt.Fprintln(c.Output, inputs...)
		}

		return nil
	})

	if err := console.Set("log", logFn); err != nil {
		return fmt.Errorf("set console.log: %w", err)
	}

	if err := global.Set("console", console); err != nil {
		return fmt.Errorf("set console: %w", err)
	}

	return nil
}
