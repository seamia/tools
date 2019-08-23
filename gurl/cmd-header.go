// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "strings"

func processHeader(params, options string) {
	comment(echoHeaderCommand, "HEADER command: %s", params)

	// do not expand the header's value - do it right before the call
	key, value := split(params)
	key = strings.TrimRight(key, ":")
	if len(key) == 0 {
		quit("Header name cannot be empty/absent")
	}

	if len(value) == 0 {
		delete(headers, key)
	} else {
		headers[key] = value
	}
}
