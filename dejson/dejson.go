// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type (
	msi   = map[string]interface{}
	slice = []interface{}
)

func loadFromJson(name string) interface{} {
	raw, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println("*** Failed to locate/read file [" + name + "], error: " + err.Error())
		// os.Exit(1)
		return nil
	}

	var c interface{}
	if err = json.Unmarshal(raw, &c); err != nil {
		fmt.Println("*** Failed to process file [" + name + "], error: " + err.Error())
		c = nil
	}
	return c
}

func saveToJson(what interface{}, name string) {

	b, err := json.MarshalIndent(what, "", "\t")
	if err != nil {
		log.Println(err)
	} else {
		err = ioutil.WriteFile(name, b, 0644)
		if err != nil {
			log.Println(err)
		}
	}
}

func warning(what string) {
	fmt.Printf("Warning: %s\n", what)
}

func trace(txt string) {
	fmt.Printf("Trace: %s\n", txt)
}

var (
	inputFlag     = flag.String("input", "", "Name of the input file")
	outputFlag    = flag.String("output", "", "Name of the output file")
	filterFlag    = flag.String("filter", "", "Name of the entry to exclude (; separated list)")
	operationFlag = flag.String("operation", "pretty", "What operation to perform (e.g. pretty, compact, etc.)")
)

func main() {

	flag.Parse()

	input := os.ExpandEnv(*inputFlag)
	output := os.ExpandEnv(*outputFlag)
	operation := os.ExpandEnv(*operationFlag)
	filter := *filterFlag

	if len(input) == 0 || len(output) == 0 || len(operation) == 0 {
		flag.Usage()
		return
	}

	data := loadFromJson(input)
	data = filterOut(data, filter)
	saveToJson(data, output)
}

func filterOut(data interface{}, filters string) interface{} {
	parts := strings.Split(filters, ";")
	if len(parts) > 0 {
		for _, filter := range parts {
			filter = strings.Trim(filter, " \t\r\n")
			data, _ = filterOutKey(data, filter)
		}
	}
	return data
}

func filterOutKey(data interface{}, filter string) (interface{}, bool) {
	if len(filter) > 0 {
		dirty := false
		if typed, converts := data.(msi); converts {
			if _, present := typed[filter]; present {
				delete(typed, filter)
				dirty = true
			}
			for k, v := range typed {
				if modified, changed := filterOutKey(v, filter); changed {
					typed[k] = modified
					dirty = true
				}
			}
			return typed, dirty

		} else if list, converts := data.(slice); converts {
			replacement := make(slice, 0, len(list))
			for _, one := range list {
				modified, _ := filterOutKey(one, filter)
				replacement = append(replacement, modified)
			}
			return replacement, true
		} else {

		}
	}
	return data, false
}
