// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/seamia/libs/resolve"

var (
	baseUrl     = "https://gurl.seamia.net/test"
	curlOptions = "-i"

	headers              = map[string]string{}
	printResponseHeaders = printResponseHeadersDefault
	generateCurlCommands = generateCurlCommandsDefault
	collectTimingInfo    = collectTimingInfoDefault
	resolveExternalFiles = resolveExternalFilesDefault

	echoProgress = echoDefault
	echoMapCommand = echoDefault
	echoSetCommand = echoDefault
	echoGetCommand = echoDefault
	echoPostCommand = echoDefault
	echoPatchCommand = echoDefault
	echoDeleteCommand = echoDefault
	echoHeaderCommand = echoDefault

	resolver = resolve.New()

	currentFile       = ""
	currentLineNumber = 0
	currentCommand    = ""

	responsePrettyPrintBody = responsePrettyPrintBodyDefault

	incrementalCounter int64
)

var (
	dials = map[string]*bool{
		"print.response.headers": &printResponseHeaders,
		"generate.curl.commands": &generateCurlCommands,
		"collect.timing.info":    &collectTimingInfo,
		"resolve.external.files": &resolveExternalFiles,
		"pretty.print.body": &responsePrettyPrintBody,

		"echo.map": &echoMapCommand,
		"echo.set": &echoSetCommand,
		"echo.get": &echoGetCommand,
		"echo.post": &echoPostCommand,
		"echo.patch": &echoPatchCommand,
		"echo.delete": &echoDeleteCommand,
		"echo.header": &echoHeaderCommand,
		"echo.progress": &echoProgress,
	}
)

var (
	savedResponse []byte
)