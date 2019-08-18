// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func processGet(params string) {
	comment("GET command: %s", params)
	call(params, "GET", "")
}
