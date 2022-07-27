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

// InjectTo injects basic console.log support.
//
// Warning: This method overwrites the previous console object.
// To add more than one console.foo method, use InjectMultipleTo.
func InjectTo(ctx *v8go.Context, opt ...Option) error {
	if ctx == nil {
		return errors.New("v8go-polyfills/console: ctx is required")
	}

	consoleMethod := NewConsole(opt...).(*consoleMethod)

	iso := ctx.Isolate()
	con := v8go.NewObjectTemplate(iso)

	logFn := v8go.NewFunctionTemplate(iso, consoleMethod.GetLogFunctionCallback())

	if err := con.Set(consoleMethod.MethodName, logFn, v8go.ReadOnly); err != nil {
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

// InjectMultipleTo injects one or more console methods to a global object of a context.
//
// Implementing the Console interface will not work in this case.
func InjectMultipleTo(ctx *v8go.Context, consoles ...Console) error {
	if ctx == nil {
		return errors.New("v8go-polyfills/console: ctx is required")
	}

	iso := ctx.Isolate()
	con := v8go.NewObjectTemplate(iso)

	for _, console := range consoles {
		consoleMethod := console.(*consoleMethod)

		logFn := v8go.NewFunctionTemplate(iso, consoleMethod.GetLogFunctionCallback())
		if err := con.Set(consoleMethod.MethodName, logFn, v8go.ReadOnly); err != nil {
			return fmt.Errorf("v8go-polyfills/console: %w", err)
		}
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
