// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rs/xid"
)

func loadDefaults() {

	// init RNG
	rand.Seed(time.Now().UnixNano())

	// add runtime-based resolutions
	resolver.Add(mapSessionKeyName, xid.New().String())

	if fullpath, err := filepath.Abs(filename()); err == nil {
		resolver.Add(mapScripFileName, filepath.Base(fullpath))
		resolver.Add(mapScripFullFileName, fullpath)
	}
	setResolverFilters()
	processCmdLine()

	// attempt to locate and load the defaults config file
	location := os.Getenv(envDefaultsLocation)
	location = expand(location)
	if len(location) == 0 {
		return
	}

	data, err := ioutil.ReadFile(location)
	if err != nil {
		reportError(err, "Loading file [%s]", location)
		return
	}

	var settings msi
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
			debug("ignoring[%s]", key)
			// generateCurlCommands = getBoolean(txt, generateCurlCommandsDefault)
		case "collect.timing.info":
			collectTimingInfo = getBoolean(txt, collectTimingInfoDefault)
		case "color":
			fmt.Println("setting color", getBoolean(txt, true))
			color.NoColor = !getBoolean(txt, true)
		default:
			if strings.HasPrefix(lower(key), configurationHeaderPrefix) {
				headerKey := key[len(configurationHeaderPrefix):]
				headers[headerKey] = txt
			}
		}
	}
	report("loaded default settings from %s", location)
}

func processCmdLine() {
	// generateCurlCommands

	if len(os.Args) > 2 {
		for i := 2; i < len(os.Args); i++ {
			param := os.Args[i]
			switch lower(param) {
			case "-silent":
				goSilent()

			case "-debug":
				// enable debug features here
				echoDebug = true

			case "-curl":
				generateCurlCommands = true

				// turn off extra reporting:
				echoHeaderCommand = false
				echoMapCommand = false
				echoGetCommand = false
				echoPostCommand = false
				echoPatchCommand = false
				echoDeleteCommand = false

				debug("enabling curl commands generations")
			default:
				debug("don't know how to handle param [%s]", param)
			}
		}
	}
}

func goSilent() {
	debug("switching to silent mode")
	echoSilent = true
}

func isSilent() bool {
	return echoSilent
}
