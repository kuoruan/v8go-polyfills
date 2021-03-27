package console

import (
	"errors"
	"fmt"

	"rogchap.com/v8go"
)

/**
Inject basic console.log support.
*/
func InjectTo(iso *v8go.Isolate, global *v8go.ObjectTemplate, opt ...Option) error {
	if iso == nil {
		return errors.New("v8go-polyfills/console: isolate is required")
	}
	if global == nil {
		return errors.New("v8go-polyfills/console: global is required")
	}

	c := NewConsole(opt...)

	con, _ := v8go.NewObjectTemplate(iso)
	logFn, _ := v8go.NewFunctionTemplate(iso, c.GetLogFunctionCallback())

	if err := con.Set("log", logFn); err != nil {
		return fmt.Errorf("v8go-polyfills/console: %w", err)
	}

	if err := global.Set("console", con); err != nil {
		return fmt.Errorf("v8go-polyfills/console: %w", err)
	}

	return nil
}
