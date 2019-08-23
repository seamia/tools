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

	echoSilent         = false
	echoDebug          = false
	echoProgress       = echoDefault
	echoMapCommand     = echoDefault
	echoSetCommand     = echoDefault
	echoGetCommand     = echoDefault
	echoPostCommand    = echoDefault
	echoPatchCommand   = echoDefault
	echoDeleteCommand  = echoDefault
	echoHeaderCommand  = echoDefault
	echoEchoCommand    = true
	echoRequireCommand = echoDefault
	echoLoadCommand    = echoDefault
	echoSectionCommand = echoDefault

	resolver = resolve.New()

	currentFile       = ""
	currentLineNumber = 0
	currentCommand    = ""

	responsePrettyPrintBody = responsePrettyPrintBodyDefault

	incrementalCounter int64
)

const (
	echoPrefix = "echo."
)

var (
	dials = map[string]*bool{
		"print.response.headers": &printResponseHeaders,
		//	"generate.curl.commands": &generateCurlCommands,
		"collect.timing.info":    &collectTimingInfo,
		"resolve.external.files": &resolveExternalFiles,
		"pretty.print.body":      &responsePrettyPrintBody,

		echoPrefix + "map":      &echoMapCommand,
		echoPrefix + "set":      &echoSetCommand,
		echoPrefix + "get":      &echoGetCommand,
		echoPrefix + "post":     &echoPostCommand,
		echoPrefix + "patch":    &echoPatchCommand,
		echoPrefix + "delete":   &echoDeleteCommand,
		echoPrefix + "header":   &echoHeaderCommand,
		echoPrefix + "progress": &echoProgress,
		echoPrefix + "echo":     &echoEchoCommand,
		echoPrefix + "require":  &echoRequireCommand,
		echoPrefix + "load":     &echoLoadCommand,
	}
)

var (
	savedResponse []byte
)

func offline() bool {
	return generateCurlCommands
}
