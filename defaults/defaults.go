// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package defaults
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/seamia/tools/debug"
	"github.com/seamia/tools/support"
)

const (
	configFileExtention = ".defaults"
)

func convertToString(what interface{}) (s string, err error) {
	switch value := what.(type) {
	case string:
		// todo: maybe need to decode it?
		s = os.ExpandEnv(value)
	case fmt.Stringer:
		s = value.String()
	case bool:
		s = fmt.Sprintf("%v", what)
	case int64:
		s = fmt.Sprintf("%d", what)
	case float64:
		s = fmt.Sprintf("%f", what)
	default:
		err = errors.New("failed to convert to string")
	}
	return
}

func applySettings(config map[string]interface{}) {

	if config != nil && len(config) > 0 {
		for key, value := range config {
			if data, err := convertToString(value); err == nil {
				if err := flag.Set(key, data); err != nil {
					if key == "unformatted" {
						// unformatted = value
					} else {
						debug.Trace("Found unsupported key (%s) in config file.", key)
					}
				} else {
					debug.Trace("\tkey (%s) = value (%s)", key, value)
				}
			} else {
				debug.Trace("Failed to convert value of the key (%s) to string - ignoring this setting", key)
			}
		}
	}
}

func loadAndApplyDefaults(implicit bool) {

	// here is the plan:
	// 0. environment ???
	// 1. if there is file named the same way as the executable + '.defaults' - we'll try to load the file and apply found settings
	// 2. if there is '-config file' option - we'll try to load the file and apply found settings

	executable := os.Args[0]
	// base := filepath.Base(executable)
	ext := len(filepath.Ext(executable))
	if ext > 0 {
		// remove the extention if present
		// base = base[:len(base)-ext]
		executable = executable[:len(executable)-ext]
	}
	// environment = strings.ToUpper(base + ":defaults")
	executable += configFileExtention

	if settings, err := support.LoadMapFromJsonFile(executable, true); err == nil {
		debug.Trace("Applying settings loaded from file: %s", executable)
		applySettings(settings)
	}

	if name, err := support.FindArgumentByPrefix("/config:"); err == nil {
		if support.Exists(name) {
			if settings, err := support.LoadMapFromJsonFile(name, true); err == nil {
				debug.Trace("Applying settings loaded from file: %s", name)
				applySettings(settings)
			} else {
				debug.Trace("Failed to load settings from file [%s] - %v.", name, err)
			}
		} else {
			debug.Trace("Specified config file [%s] not found.", name)
		}
	}
}

func LoadAndApplyDefaults() {
	loadAndApplyDefaults(false)
}

func init() {
	// this is module init func. it is called before the 'main'
	loadAndApplyDefaults(true)
}
