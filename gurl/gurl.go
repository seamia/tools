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
	debug("args: %v, %v", len(os.Args), os.Args)

	// todo: del this
	os.Setenv("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjb20ubWFnaWNsZWFwLndlYi5kZXZlbG9wZXIiLCJleHAiOjE1NjY1MTU2MzIsImp0aSI6IjRka3kyamszZG03cjFxM3MiLCJpYXQiOjE1NjY1MTIwMzIsImlzcyI6Imh0dHBzOi8vYXV0aC5kZXYubWFnaWNsZWFwLmJsdWUiLCJzdWIiOiJ1cy1lYXN0LTE6MDBiMzJkNGQtNmU3MC00NmViLTljMTgtNDI4Y2FiNzFhZWI3IiwidG9rZW5fdXNlIjoiYWNjZXNzIiwiY2xpZW50X2lkIjoiY29tLm1hZ2ljbGVhcC53ZWIuZGV2ZWxvcGVyIiwic2NvcGUiOlsic3NvOnplbmRlc2siLCJjb2duaXRvIiwicGhvbmUiLCJlbWFpbCIsInByb2ZpbGUiXSwicnRpZCI6IjFjZmUwNGExLTJjN2UtNDg3Yy1hODg2LTYxMDI5MGM0YWM3ZSIsInN0eXAiOiJ1c2VyIn0.ZHucmolVM0ltS484xPMDzL8nsjpcUDVLe_-YLZv-HwM")

	loadDefaults()

	printer.Set(debug) // todo: revisit this

	data, err := ioutil.ReadFile(filename())
	quitOnError(err, "Opening file %s", filename())

	comment(echoProgress, "Processing file %s", filename())
	generate("# generating curls commands from %s", filename())
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
	insideCommentBlock := false

	for lineNumber, line := range lines {
		currentLineNumber = lineNumber + offset

		// ignore whitespace
		line = strings.TrimLeft(line, leadingWhiteSpace)
		line = strings.TrimRight(line, trainingWhiteSpace)

		if insideCommentBlock {
			if strings.HasSuffix(line, "*/") {
				insideCommentBlock = false
			}
			continue
		}

		if strings.HasPrefix(line, "/*") {
			if !strings.HasSuffix(line, "*/") {
				insideCommentBlock = true
			}
			continue
		}

		// ignore comments
		if strings.HasPrefix(line, commentPrefix) {
			continue
		}
		if pound := strings.Index(line, commentPrefix); pound > 0 {
			line = strings.TrimSpace(line[:pound])
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
	fullcmd, payload := split(command)
	cmd, options := splitBy(fullcmd, ":")

	if handler, found := handlers[lower(cmd)]; found {
		handler(payload, options)
	} else {
		quit("Unknown command [%s]", fullcmd)
	}

	// fmt.Println("========", command)
}

type cmdHandler func(params, options string)

var handlers = map[string]cmdHandler{
	"set":    processSet,
	"map":    processMap,
	"header": processHeader,

	"get":    processGet,
	"patch":  processPatch,
	"post":   processPost,
	"delete": processDelete,

	"echo":    processEcho,
	"require": processRequire,
	"load":    processLoad,
	"section": processSection,
}
