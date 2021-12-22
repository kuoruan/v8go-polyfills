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

package base64

import (
	stdBase64 "encoding/base64"

	"rogchap.com/v8go"
)

type Base64 interface {
	GetAtobFunctionCallback() v8go.FunctionCallback
	GetBtoaFunctionCallback() v8go.FunctionCallback
}

type base64 struct {
}

func NewBase64() Base64 {
	return &base64{}
}

/*
 https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/atob
*/
func (b *base64) GetAtobFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		ctx := info.Context()

		if len(args) <= 0 {
			// TODO: v8go can't throw a error now, so we return an empty string
			return newStringValue(ctx, "")
		}

		encoded := args[0].String()

		byts, err := stdBase64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return newStringValue(ctx, "")
		}

		return newStringValue(ctx, string(byts))
	}
}

/*
 https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/btoa
*/
func (b *base64) GetBtoaFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		ctx := info.Context()

		if len(args) <= 0 {
			return newStringValue(ctx, "")
		}

		str := args[0].String()

		encoded := stdBase64.StdEncoding.EncodeToString([]byte(str))
		return newStringValue(ctx, encoded)
	}
}

func newStringValue(ctx *v8go.Context, str string) *v8go.Value {
	iso := ctx.Isolate()
	val, _ := v8go.NewValue(iso, str)
	return val
}
