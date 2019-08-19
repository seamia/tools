// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/fatih/color"
)

func processSet(params string) {
	comment("SET command: %s", params)
	key, value := split(expand(params))

	switch lower(key) {
	case "baseurl":
		baseUrl = value
	case "producecurl":
		generateCurlCommands = getBoolean(value, true)
	case "prettyprintbody":
		responsePrettyPrintBody = getBoolean(value, responsePrettyPrintBodyDefault)
	case "nocolor":
		color.NoColor = true // disables colorized output
	case "color":
		color.NoColor = false
	default:
		quit("Unknown SET: [%s]", key)
	}
}
