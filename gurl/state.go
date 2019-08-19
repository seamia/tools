// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/seamia/libs/resolve"

var (
	baseUrl     = "https://gurl.seamia.net/test"
	curlOptions = "-i"

	printResponseHeaders = printResponseHeadersDefault
	headers              = map[string]string{}
	generateCurlCommands = generateCurlCommandsDefault
	collectTimingInfo    = collectTimingInfoDefault

	resolver = resolve.New()

	currentFile       = ""
	currentLineNumber = 0
	currentCommand    = ""

	responsePrettyPrintBody = responsePrettyPrintBodyDefault

	incrementalCounter int64
)
