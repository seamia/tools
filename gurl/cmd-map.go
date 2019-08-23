// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "net/url"

func processMap(params, options string) {
	comment(echoMapCommand, "MAP command: %s", params)
	key, value := split(expand(params))

	if len(options) > 0 {
		if lower(options) == "encode" {
			value = url.QueryEscape(value)
		} else {
			quit("unknown options: %s", options)
		}
	}

	resolver.Add(key, value)

	if offline() {
		generate("%s=%s", key, value)
	}

}
