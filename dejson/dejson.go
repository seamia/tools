// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	compressedStreamPrefix = 31
)

type (
	msi   = map[string]interface{}
	slice = []interface{}
)

func loadFromJson(name string) interface{} {
	raw, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println("*** Failed to locate/read file ["+name+"], error: ", err)
		// os.Exit(1)
		return nil
	}

	raw, err = decompress(raw)
	if err != nil {
		fmt.Println("*** Failed to decompress file ["+name+"], error: ", err)
		return nil
	}

	var c interface{}
	if err = json.Unmarshal(raw, &c); err != nil {
		fmt.Println("*** Failed to process file ["+name+"], error: ", err)
		c = nil
	}
	return c
}

func saveToJson(what interface{}, name string) {

	b, err := json.MarshalIndent(what, "", "\t")
	if err != nil {
		log.Println(err)
	} else {

		if len(name) == 0 || name == "stdout" || name == "console" || name == "-" || name == "screen" {
			_, err = os.Stdout.Write(b)
		} else {
			err = ioutil.WriteFile(name, b, 0644)
		}

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

	if len(input) == 0 || len(operation) == 0 {
		if len(os.Args) >= 2 {
			input = os.ExpandEnv(os.Args[1])
			if len(os.Args) >= 3 {
				output = os.ExpandEnv(os.Args[2])
			}
		}

		if len(input) == 0 || len(operation) == 0 {
			flag.Usage()
			return
		}
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

func decompress(what []byte) ([]byte, error) {

	if len(what) == 0 {
		// the source is empty - there is nothing here to decompress
		return what, nil
	}

	if what[0] != compressedStreamPrefix {
		// it doesn't seem to be compressed - return the source
		if what[0] != byte('{') && what[0] != byte('[') {
			fmt.Println("hmmmm.... unextected prefix of persisted stream ....")
		}
		return what, nil
	}

	gz, err := gzip.NewReader(bytes.NewBuffer(what))
	if err != nil {
		return nil, fmt.Errorf("Read: %v", err)
	}

	var decompressed bytes.Buffer
	_, err = io.Copy(&decompressed, gz)
	errClose := gz.Close()
	if err != nil {
		return nil, err
	}
	if errClose != nil {
		return nil, errClose
	}

	return decompressed.Bytes(), nil
}
