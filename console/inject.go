/*
 * Copyright (c) 2021 Xingwang Liao
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package console

import (
	"errors"
	"fmt"

	"rogchap.com/v8go"
)

/**
Inject basic console.log support.
*/
func InjectTo(ctx *v8go.Context, opt ...Option) error {
	if ctx == nil {
		return errors.New("v8go-polyfills/console: ctx is required")
	}

	iso, err := ctx.Isolate()
	if err != nil {
		return fmt.Errorf("v8go-polyfills/console: %w", err)
	}

	c := NewConsole(opt...)

	con, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		return fmt.Errorf("v8go-polyfills/console: %w", err)
	}

	logFn, err := v8go.NewFunctionTemplate(iso, c.GetLogFunctionCallback())
	if err != nil {
		return fmt.Errorf("v8go-polyfills/console: %w", err)
	}

	if err := con.Set("log", logFn, v8go.ReadOnly); err != nil {
		return fmt.Errorf("v8go-polyfills/console: %w", err)
	}

	conObj, err := con.NewInstance(ctx)
	if err != nil {
		return fmt.Errorf("v8go-polyfills/console: %w", err)
	}

	global := ctx.Global()

	if err := global.Set("console", conObj); err != nil {
		return fmt.Errorf("v8go-polyfills/console: %w", err)
	}

	return nil
}
