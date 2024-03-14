// Copyright 2020 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type (
	msi = map[string]interface{}
)

func alert(name, address string) {
	sendEmail([]string{address}, name)
}

func trimPrefix(original string, key interface{}) string {
	if key == nil {
		return original
	}
	var steps []string
	switch actual := key.(type) {
	case string:
		steps = []string{actual}
	case []interface{}:
		for _, one := range actual {
			if txt, converts := one.(string); converts && len(txt) > 0 {
				steps = append(steps, txt)
			}
		}
	default:
		report("unhandled type: %v", actual)
	}

	for _, step := range steps {
		if found := strings.Index(original, step); len(step) > 0 && found > 0 {
			original = original[found+len(step):]
		}
	}

	original = strings.Trim(original, "\r\n\t ")
	return original
}

func trimSuffix(original string, key interface{}) string {
	if key == nil {
		return original
	}
	var steps []string
	switch actual := key.(type) {
	case string:
		steps = []string{actual}
	case []interface{}:
		for _, one := range actual {
			if txt, converts := one.(string); converts && len(txt) > 0 {
				steps = append(steps, txt)
			}
		}
	default:
		report("unhandled type: %v", actual)
	}

	for _, step := range steps {
		if found := strings.Index(original, step); len(step) > 0 && found > 0 {
			original = original[:found]
		}
	}

	original = strings.Trim(original, "\r\n\t ")
	return original
}

func report(format string, arg ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, format+"\n", arg...)
}

func quitOnError(err error, format string, a ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, format+"\n", a...)
		fmt.Fprintf(os.Stderr, "there was an error: %v\n", err)
		panic(err)
	}
}

func saveContent(name, content string) {
	name = strings.ReplaceAll(name, "/", "_")
	filename := ".save/" + name
	if err := os.WriteFile(filename, []byte(content), 0666); err != nil {
		fmt.Fprintf(os.Stderr, "failed to save file (%s) with error: %v\n", filename, err)
	}
}

func loadList(script string) []msi {
	data, err := os.ReadFile(script)
	quitOnError(err, "failed to read file: %s", script)

	var list []interface{}
	err = json.Unmarshal(data, &list)
	quitOnError(err, "failed to unmarshal loaded script - bad format? file: %s", script)

	result := []msi{}
	for _, task := range list {
		dict, converts := task.(msi)
		if converts {
			result = append(result, dict)
		} else {
			// todo: ???
		}
	}
	return result
}

func loadDict(script string) msi {
	data, err := os.ReadFile(script)
	quitOnError(err, "failed to read file: %s", script)

	var what msi
	err = json.Unmarshal(data, &what)
	quitOnError(err, "failed to unmarshal loaded script - bad format? file: %s", script)

	return what
}

func get(from msi, key string, fallback string) string {
	if len(from) > 0 {
		if entry, found := from[key]; found {
			if txt, converts := entry.(string); converts {
				return txt
			}
		}
	}
	return fallback
}

var (
	txt2bool = map[string]bool{
		"yes": true,
		"no":  false,
	}
)

func getFlag(from msi, key string, fallback bool) bool {
	if len(from) > 0 {
		if entry, found := from[key]; found {
			if flag, converts := entry.(bool); converts {
				return flag
			}
		} else if entry, found := from[key]; found {
			if txt, converts := entry.(string); converts {
				if flag, found := txt2bool[strings.ToLower(txt)]; found {
					return flag
				}
			}
		}
	}
	return fallback
}
