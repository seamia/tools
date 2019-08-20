// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func processPatch(params string) {
	comment(echoPatchCommand, "PATCH command: %s", params)
	relativeUrl, payload := split(expand(params))
	call(relativeUrl, "PATCH", payload)
}
