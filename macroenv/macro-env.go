// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	separator              = "="
	allowUnresolvedEntries = true
)

func main() {

	quit(len(os.Args) != 4, "USAGE: macroenv source substitutes result\n(see github.com/seamia/tools/macroenv)")

	srcFile := os.Args[1]
	subFile := os.Args[2]
	resFile := os.Args[3]

	src := readFile(srcFile)
	subst := loadSubstitutes(subFile)
	result := replace(src, subst)
	saveFile(resFile, result)
}

func replace(src string, subs map[string]string) string {
	return os.Expand(src, func(key string) string {
		if value, found := subs[key]; found {
			return value
		} else if allowUnresolvedEntries {
			return fmt.Sprintf("${%s}", key)
		} else {
			quit(true, "failed to find a substitute for key [%s]", key)
			return "disaster" // to make compile 'happy'
		}
	})
}

func loadSubstitutes(names string) map[string]string {
	subst := make(map[string]string)
	for _, name := range strings.Split(names, ";") {
		loadSubstitute(name, &subst)
	}
	return subst
}

func loadSubstitute(from string, to *map[string]string) {
	substitutes := readFile(from)

	for n, line := range strings.Split(substitutes, "\n") {
		if len(trim(line)) == 0 {
			// empty string --> nothing to do
			continue
		}
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			// this is a comment --> ignore it
			continue
		}

		equal := strings.Index(line, separator)
		quit(equal <= 0, "failed to find separator; file: %s; line %d", from, n)

		key := trim(line[:equal])
		value := trim(line[equal+len(separator):])
		quit(len(key) <= 0, "the key is empty; file: %s; line %d", from, n)

		// allow insert of environment variables
		value = os.ExpandEnv(value)

		if evalue, exists := (*to)[key]; exists {
			quit(true, "found a dup key [%s] inside file [%s], set to [%s]", key, from, evalue)
		}
		(*to)[key] = value
	}
}

func quit(condition bool, format string, a ...interface{}) {
	if condition {
		fmt.Fprintf(os.Stderr, "Fatal error: "+format+"\n", a...)
		os.Exit(11)
	}
}

func quitOnError(err error, format string, a ...interface{}) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "Fatal error: "+format+"\n", a...)
	fmt.Fprintf(os.Stderr, "(details: %v)\n", err)

	os.Exit(7)
}

func readFile(name string) string {

	data, err := ioutil.ReadFile(name)
	quitOnError(err, "failed to read file: %s", name)
	return string(data)
}

func trim(src string) string {
	return strings.Trim(src, " \t\r")
}

func saveFile(name string, data string) {
	err := ioutil.WriteFile(name, []byte(data), 0644)
	quitOnError(err, "failed to save file %s", name)
}
