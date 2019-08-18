// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/seamia/libs/printer"
)

func main() {
	if len(os.Args) < 2 || help(filename()) {
		usage()
	}

	loadDefaults()

	printer.Set(printer.Stderr) // todo: revisit this

	data, err := ioutil.ReadFile(filename())
	quitOnError(err, "Opening file %s", filename())

	comment("Processing file %s", filename())
	currentFile = filename()
	processScript(string(data))

	os.Exit(exitCodeOnSuccess)
}

func processScript(script string) {
	lines := strings.Split(script, lineSeparator)

	// skip shebang
	offset := 1
	if strings.HasPrefix(lines[0], shebang) {
		lines = lines[1:]
		offset++
	}

	command := []string{}
	for lineNumber, line := range lines {
		currentLineNumber = lineNumber + offset

		// ignore whitespace
		line = strings.TrimLeft(line, leadingWhiteSpace)
		line = strings.TrimRight(line, trainingWhiteSpace)

		// ignore comments
		if strings.HasPrefix(line, commentPrefix) {
			continue
		}

		if len(line) == 0 {
			processCommand(strings.Join(command, " "))
			command = []string{}
		} else {
			if len(command) == 0 {
				cmd, _ := split(line)
				if !multiLineCommand(cmd) {
					processCommand(line)
					continue
				}
			}
			command = append(command, line)
		}
	}

	// deal with the remains ...
	if len(command) > 0 {
		processCommand(strings.Join(command, " "))
		command = []string{}
	}
}

func processCommand(command string) {
	command = strings.TrimSpace(command)
	if len(command) == 0 {
		return
	}

	currentCommand = command
	cmd, payload := split(command)

	switch lower(cmd) {
	case "set":
		processSet(payload)
	case "map":
		processMap(payload)
	case "header":
		processHeader(payload)
	case "get":
		processGet(payload)
	case "patch":
		processPatch(payload)
	case "post":
		processPost(payload)
	default:
		quit("Unknown command [%s]", cmd)
	}

	// fmt.Println("========", command)
}