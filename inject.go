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

package polyfills

import (
	"go.kuoruan.net/v8go-polyfills/base64"
	"go.kuoruan.net/v8go-polyfills/console"
	"go.kuoruan.net/v8go-polyfills/fetch"
	"go.kuoruan.net/v8go-polyfills/internal"
	"go.kuoruan.net/v8go-polyfills/timers"
	"go.kuoruan.net/v8go-polyfills/url"

	"rogchap.com/v8go"
)

func InjectToGlobalObject(iso *v8go.Isolate, global *v8go.ObjectTemplate, opt ...interface{}) error {
	var fetchOpts []fetch.Option

	for _, o := range opt {
		switch t := o.(type) {
		case fetch.Option:
			fetchOpts = append(fetchOpts, t)
		}
	}

	if err := fetch.InjectTo(iso, global, fetchOpts...); err != nil {
		return err
	}

	if err := base64.InjectTo(iso, global); err != nil {
		return err
	}

	if err := timers.InjectTo(iso, global); err != nil {
		return err
	}

	return nil
}

func InjectToContext(ctx *v8go.Context, opt ...interface{}) error {
	var consoleOpts []console.Option

	for _, o := range opt {
		switch t := o.(type) {
		case console.Option:
			consoleOpts = append(consoleOpts, t)
		}
	}

	for _, p := range []func(*v8go.Context) error{
		url.InjectTo,
	} {
		if err := p(ctx); err != nil {
			return err
		}
	}

	if err := console.InjectTo(ctx, consoleOpts...); err != nil {
		return err
	}

	return nil
}

func Version() string {
	return internal.Version
}
