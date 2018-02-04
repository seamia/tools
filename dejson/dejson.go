// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"os"
	"log"
)

func loadFromJson(name string) interface{} {
	raw, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println("*** Failed to locate/read file ["+name+"], error: "+err.Error())
		// os.Exit(1)
		return nil
	}

	var c interface{}
	if err = json.Unmarshal(raw, &c); err != nil {
		fmt.Println("*** Failed to process file ["+name+"], error: "+err.Error())
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
	operationFlag = flag.String("operation", "pretty", "What operation to perform (e.g. pretty, compact, etc.)")
)

func main() {

	flag.Parse()

	input := os.ExpandEnv(*inputFlag)
	output := os.ExpandEnv(*outputFlag)
	operation := os.ExpandEnv(*operationFlag)

	if len(input) == 0 || len(output) == 0 || len(operation) == 0{
		flag.Usage()
		return
	}

	data := loadFromJson(input)
	saveToJson(data, output)
}
