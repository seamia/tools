// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"
)

// Require ${response:status} HEALTHY

func processLoad(params string) {
	comment(echoLoadCommand, "LOAD: %s", params)

	parts := strings.Split(params, " ")
	if len(parts) >= 3 {
		entry := parts[0]
		filename := parts[1]
		format := lower(parts[2])
		key := parts[3]

		fullfilename, err := expandPath(filename)
		quitOnError(err, "Failed to process file [%s]", filename)

		if format != "json" {
			quit("LOAD command has wrong format [%s]", format)
		}

		data, err := ioutil.ReadFile(fullfilename)
		quitOnError(err, "reading file [%s]", filename)

		var receiver interface{}
		err = json.Unmarshal(data, &receiver)
		quitOnError(err, "parsing content of file [%s]", filename)

		if success, value := resolveAny(receiver, key); success {

			value = expand(value)
			resolver.Add(entry, value)

		} else {
			quit("Cannot resolve key [%s] inside of the content of file [%s]", key, filename)
		}

	} else {
		quit("LOAD command has wrong arguments [%s]", params)
	}
}

// Dir returns the home directory for the executing user.
// An error is returned if a home directory cannot be detected.
func dir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	if currentUser.HomeDir == "" {
		return "", errors.New("cannot find user-specific home dir")
	}

	return currentUser.HomeDir, nil
}

// Expand expands the path to include the home directory if the path
// is prefixed with `~`. If it isn't prefixed with `~`, the path is
// returned as-is.
func expandPath(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand user-specific home dir")
	}

	dir, err := dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, path[1:]), nil
}
