// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "strings"

type Evaluate func(left, right string) bool

func evalEqual(left, right string) bool {
	return left == right
}

func evalNotEqual(left, right string) bool {
	return left != right
}

func evalContains(left, right string) bool {
	return strings.Index(left, right) != -1
}

func evalNotContains(left, right string) bool {
	return !evalContains(left, right)
}

func evalSuffix(left, right string) bool {
	return strings.HasSuffix(left, right)
}

func evalPrefix(left, right string) bool {
	return strings.HasPrefix(left, right)
}

func hasFieldValue(what interface{}, field, looking4 string, compare Evaluate) bool {
	if object, converts := what.(msi); converts {
		if data, exists := object[field]; exists {
			if value, converts := data.(string); converts {
				return compare(value, looking4)
			} else {
				// failed to convert to string ....
			}
		}
	}
	return false
}

var evaluateMap = map[string]Evaluate{

	// the key in this map MUST have "=" in it

	"==": evalEqual,
	"!=": evalNotEqual,
	"<=": evalContains,    // < =
	"=>": evalNotContains, // = >
	"=(": evalPrefix,
	"=)": evalSuffix,
}

func findEvaluator(cmd string) (string, string, Evaluate) {
	for key, evalFunc := range evaluateMap {
		if index := strings.Index(cmd, key); index >= 0 {
			field := strings.TrimSpace(cmd[:index])
			value := strings.TrimSpace(cmd[index+len(key):])
			return field, value, evalFunc
		}
	}
	quit("cannot find an appopriate evaluator for [%s] expression", cmd)
	return "", "", nil
}
