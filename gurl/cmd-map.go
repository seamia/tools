// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func processMap(params string) {
	comment(echoMapCommand, "MAP command: %s", params)
	key, value := split(expand(params))
	resolver.Add(key, value)
}
