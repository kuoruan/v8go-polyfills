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
