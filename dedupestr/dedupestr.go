// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: dedupestr from to\n")
		return
	}

	data, err := ioutil.ReadFile(os.Args[1])
	failure("input file", err)

	lines := strings.Split(string(data), lineSeparator)
	for i, str := range lines {
		lines[i] = strings.Trim(str, trimCutSet)
	}
	sort.Strings(lines)

	unique := []string{lines[0]}
	for i := 1; i < len(lines); i++ {
		if lines[i] != lines[i-1] {
			unique = append(unique, lines[i])
		}
	}

	err = ioutil.WriteFile(os.Args[2], []byte(strings.Join(unique, lineSeparator)), filePremissions)
	failure("output file", err)
}

func failure(what string, err error) {
	if err != nil {
		fmt.Printf("Failure: %s (%v)\n", what, err)
		os.Exit(7)
	}
}

const (
	lineSeparator   = "\n"
	trimCutSet      = " \t\r"
	filePremissions = 0644
)
