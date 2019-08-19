// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func processPost(params string) {
	comment("POST command: %s", params)
	relativeUrl, payload := split(expand(params))
	call(relativeUrl, "POST", payload)
}
