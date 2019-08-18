// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rs/xid"
)

func help(txt string) bool {
	if txt == "/help" {
		return true
	}
	return false
}

func usage() {
	color.Set(colorUsage)
	defer color.Unset()

	fmt.Println("Usage: ...")
	os.Exit(exitCodeOnUsage)
}

func quitOnError(err error, format string, a ...interface{}) {
	if err != nil {
		reportError(err, format, a...)
		os.Exit(exitCodeOnError)
	}
}

func reportError(err error, format string, a ...interface{}) {
	color.Set(colorError)
	defer color.Unset()

	_, _ = fmt.Fprintf(os.Stderr, "Got an error: %v, while ", err)
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", a...)

	if len(currentFile) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "(script: [%s], line: %v)\n", currentFile, currentLineNumber)
	}

	if len(currentCommand) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "(command: [%s])\n", currentCommand)
	}
}

func quit(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "Got an error while "+format+"\n", a...)
	os.Exit(exitCodeOnError)
}

func comment(format string, a ...interface{}) {
	colorPrint(colorComment, format, a...)
}

func report(format string, a ...interface{}) {
	colorPrint(colorComment, format, a...)
}

func response(format string, a ...interface{}) {
	colorPrint(colorResponse, format, a...)
}

func responseSuccess(format string, a ...interface{}) {
	colorPrint(colorResponseSuccess, format, a...)
}
func responseFailure(format string, a ...interface{}) {
	colorPrint(colorResponseFailure, format, a...)
}

func colorPrint(clr color.Attribute, format string, a ...interface{}) {
	color.Set(clr)
	defer color.Unset()

	_, _ = fmt.Fprintf(os.Stdout, format+"\n", a...)
}

func filename() string {
	return os.Args[1]
}

func niy() {
	quit("not implemented yet")
}

func expand(from string) string {
	return resolver.Text(from)
}

func split(src string) (string, string) {
	cmd, payload := src, ""
	if index := strings.IndexAny(src, wordSeparator); index > 0 {
		cmd = src[:index]
		payload = strings.TrimSpace(src[index:])
	}
	return cmd, payload
}

func multiLineCommand(cmd string) bool {
	switch lower(cmd) {
	case "post", "get", "put", "patch", "delete":
		return true
	}
	return false
}

func getBoolean(src string, fallback bool) bool {
	switch lower(src) {
	case "true", "yes", "ok", "okay", "please", "do it", "go ahead", "sure", "affirmative", "yeap", "yeah":
		return true
	case "false", "no", "nope", "nada", "no way":
		return false
	}
	return fallback
}

func loadDefaults() {

	resolver.Add(mapSessionKeyName, xid.New().String())

	location := os.Getenv(envDefaultsLocation)
	if len(location) == 0 {
		return
	}

	data, err := ioutil.ReadFile(location)
	if err != nil {
		reportError(err, "Loading file [%s]", location)
		return
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		reportError(err, "Parsing content of [%s]", location)
	}

	for key, value := range settings {
		txt, _ := value.(string)
		switch lower(key) {
		case "base.url":
			baseUrl = txt
		case "curl.options":
			curlOptions = txt
		case "print.response.headers":
			printResponseHeaders = getBoolean(txt, printResponseHeadersDefault)
		case "generate.curl.commands":
			generateCurlCommands = getBoolean(txt, generateCurlCommandsDefault)
		case "collect.timing.info":
			collectTimingInfo = getBoolean(txt, collectTimingInfoDefault)
		default:
			if strings.HasPrefix(lower(key), configurationHeaderPrefix) {
				headerKey := key[len(configurationHeaderPrefix):]
				headers[headerKey] = txt
			}
		}
	}
	report("loaded default settings from %s", location)
}

func loadExternalFile(src string) string {
	external, filename := dataPointsToExternalFile(src)
	if !external {
		return src
	}
	data, err := ioutil.ReadFile(filename)
	quitOnError(err, "Opening file [%s]", filename)
	return string(data)
}

func dataPointsToExternalFile(src string) (bool, string) {
	if strings.HasPrefix(src, externalFilePrefix) {
		return true, src[len(externalFilePrefix):]
	}
	return false, src
}

func lower(src string) string {
	return strings.ToLower(src)
}
