// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
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
	// global savedResponse []byte
	if len(savedResponse) == 0 {
		return false, key
	}

	var holder map[string]interface{}
	if err := json.Unmarshal(savedResponse, holder); err != nil {
		reportError(err, "failed to ingest json from response")
		return false, key
	}

	if data, found := holder[key]; found {
		if txt, okay := data.(string); okay {
			return true, txt
		}
	}

	//last resort
	key = lower(key)
	for name, settings := range holder {
		if lower(name) == key {
			if txt, okay := settings.(string); okay {
				return true, txt
			}
		}
	}

	return false, key
}
