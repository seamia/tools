// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math/rand"
	"sort"
	"strings"
)

func resolveAny(src interface{}, key string) (bool, string) {

	if src == nil {
		return notFound() // todo: could it be a legit return value
	}
	switch actual := src.(type) {
	case msi:
		return resolveMap(actual, key)
	case slice:
		return resolveSlice(actual, key)
	case string:
		if len(key) == 0 {
			return true, actual
		} else {
			quit("stil have non empty path [%s] for terminal value [%s]", key, actual)
		}
	default:
		quit("unhandled type: %v", actual)
	}
	return notFound()
}

func breakPath(src string) (string, string) {
	bits := strings.Split(src, itemsSeparator)
	first := bits[0]
	remainder := strings.Join(bits[1:], itemsSeparator)
	return first, remainder
}

func breakParam(src string) (string, string) {
	bits := strings.Split(src, ":")
	switch len(bits) {
	case 1:
		return bits[0], ""
	case 2:
		return bits[0], bits[1]
	default:
		quit("too many parts in [%s]", src)
	}
	return "", ""
}

func resolveMap(src msi, key string) (bool, string) {
	first, remainder := breakPath(key)
	if data, found := src[first]; found {
		/*
			if txt, okay := data.(string); okay {
				return true, txt
			}*/
		return resolveAny(data, remainder)
	}

	return notFound()
}

func resolveSlice(src slice, key string) (bool, string) {
	if len(src) == 0 {
		return notFound()
	}

	cmds := strings.Split(key, ";")
	if len(cmds) > 1 {
		exact := make([]string, 0)
		inexact := make([]string, 0)
		for _, cmd := range cmds {
			if strings.Contains(cmd, "=") {
				exact = append(exact, cmd)
			} else if strings.Contains(cmd, ":") {
				inexact = append(inexact, cmd)
			} else {
				quit("illegal value [%s], which a part of [%s]", cmd, key)
			}
		}
		if len(inexact) > 1 {
			quit("more than one selectors: [%s], which is a part of [%s]", strings.Join(inexact, "; "), key)
		}

		for _, one := range exact {
			src = reduceSlice(src, one)
		}

		if len(inexact) == 1 {
			return resolveSlice(src, inexact[0])
		} else {
			quit("cannot resolve the condition [%s]", key)
		}

		return notFound()
	}

	first, remainder := breakPath(key)
	name, options := breakParam(first)

	index := indexInvalid
	switch options {
	case "first":
		index = 0
	case "last":
		index = len(src) - 1
	case "random":
		index = rand.Intn(len(src))
	}

	if index == indexInvalid {
		selector := make(map[string]int)
		for index, item := range src {
			if success, value := resolveAny(item, name); success {
				selector[value] = index
			} else {
				// ?????
			}
		}

		var keys []string
		for k := range selector {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		index = indexInvalid
		switch options {
		case "latest":
			index = selector[keys[len(keys)-1]]
		case "earliest":
			index = selector[keys[0]]
		}
	}

	if index != indexInvalid {
		return resolveAny(src[index], remainder)
	} else {

	}

	comment(true, "********************* %s; %s; %s;", name, options, remainder)

	return notFound()
}

func reduceSlice(src slice, key string) slice {
	field, value, evaluator := findEvaluator(key)
	result := make(slice, 0, len(src))

	for _, entry := range src {
		if hasFieldValue(entry, field, value, evaluator) {
			result = append(result, entry)
		}
	}

	debug("-- reduced slice from %v to %v using [%s] condition", len(src), len(result), key)
	return result
}

func notFound() (bool, string) {
	return false, ""
}
