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

package timers

import (
	"fmt"

	"rogchap.com/v8go"
)

func InjectTo(iso *v8go.Isolate, global *v8go.ObjectTemplate) error {
	t := NewTimers()

	for _, f := range []struct {
		Name string
		Func func() v8go.FunctionCallback
	}{
		{Name: "setTimeout", Func: t.GetSetTimeoutFunctionCallback},
		{Name: "setInterval", Func: t.GetSetIntervalFunctionCallback},
		{Name: "clearTimeout", Func: t.GetClearTimeoutFunctionCallback},
		{Name: "clearInterval", Func: t.GetClearIntervalFunctionCallback},
	} {
		fn := v8go.NewFunctionTemplate(iso, f.Func())

		if err := global.Set(f.Name, fn, v8go.ReadOnly); err != nil {
			return fmt.Errorf("v8go-polyfills/timers: %w", err)
		}
	}

	return nil
}
