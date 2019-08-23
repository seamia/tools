// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func processEcho(params, options string) {
	if offline() {
		return
	}
	comment(echoEchoCommand, "ECHO: %s", expand(params))
}

func processSection(params, options string) {
	if offline() {
		return
	}
	section(echoSectionCommand, "%s", expand(params))
}
