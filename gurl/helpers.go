// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
)

func help(txt string) bool {
	if txt == "/help" || strings.HasPrefix(txt, "-") {
		return true
	}
	return false
}

func usage() {
	color.Set(colorUsage)
	fmt.Println("Usage: gurl script.gurl")
	fmt.Println(versionInfo)
	color.Unset()

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

func comment(allow bool, format string, a ...interface{}) {
	if allow && !isSilent() {
		colorPrint(colorComment, format, a...)
	}
}

func section(allow bool, format string, a ...interface{}) {
	const highlight = "================================================= "
	if allow && !isSilent() {
		colorPrint(colorSection, highlight+format, a...)
	}
}

func report(format string, a ...interface{}) {
	colorPrint(colorComment, format, a...)
}

func debug(format string, a ...interface{}) {
	if echoDebug {
		colorPrint(colorDebug, format, a...)
	}
}

func response(format string, a ...interface{}) {
	colorPrint(colorResponse, format, a...)
}

func responseSuccess(format string, a ...interface{}) {
	if !isSilent() {
		colorPrint(colorResponseSuccess, format, a...)
	}
}
func responseFailure(format string, a ...interface{}) {
	colorPrint(colorResponseFailure, format, a...)
}
func responseAttention(format string, a ...interface{}) {
	colorPrint(colorResponseAttention, format, a...)
}

func colorPrint(clr interface{}, format string, a ...interface{}) {
	switch actual := clr.(type) {
	case color.Attribute:
		color.Set(actual)
	case []color.Attribute:
		color.Set(actual...)
	default:
		debug("unhandled type %v", actual)
	}

	_, _ = fmt.Fprintf(os.Stdout, format, a...)
	color.Unset()
	_, _ = fmt.Fprintf(os.Stdout, "\n")
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
	return splitBy(src, wordSeparator)
}

func splitBy(src, separator string) (string, string) {
	cmd, payload := src, ""
	if index := strings.IndexAny(src, separator); index > 0 {
		cmd = src[:index]
		payload = strings.TrimSpace(src[index+1:])
	}
	return strings.TrimSpace(cmd), payload
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

func loadExternalFile(src string) string {
	external, filename := dataPointsToExternalFile(src)
	if !external {
		return src
	}
	data, err := ioutil.ReadFile(filename)
	quitOnError(err, "Opening file [%s]", filename)

	txt := string(data)
	if resolveExternalFiles {
		txt = expand(txt)
	}
	return txt
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

func generate(format string, a ...interface{}) {
	if offline() {
		_, _ = fmt.Fprintf(os.Stdout, ""+format+"\n", a...)
	}
}
