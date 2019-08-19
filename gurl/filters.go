// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/rs/xid"
)

func setResolverFilters() {
	resolver.SetFilter(preFilter, true)
}

func preFilter(key string) (bool, string) {
	ley := lower(key)
	switch ley {
	case "random":
		return true, xid.New().String()
	case "increment":
		return true, strconv.FormatInt(atomic.AddInt64(&incrementalCounter, 1), 10)
	default:
		if strings.HasPrefix(ley, mappingResponseValues) {
			return responseValue(key[len(mappingResponseValues):])
		}
		return false, key
	}
}

func responseValue(key string) (bool, string) {
	return false, key
}
