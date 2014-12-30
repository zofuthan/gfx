// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build 386,!gles2 amd64,!gles2

package glc

func NewContext() *Context {
	return &Context{
		GetError: getError,
	}
}
