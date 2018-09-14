// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package support

import (
	"encoding/json"
	// "github.com/seamia/tools/debug"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		return os.IsExist(err)
	}
	return true
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// annotate("Creating directory: "+dir)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			//debug.Trace("*** Failed to create directory [" + dir + "], error: " + err.Error())
			panic(err)
		}
	}
}

// Load json data from given file
// if 'strict' is false - attempt to load file + '.json' extention
func LoadMapFromJsonFile(name string, strict bool) (map[string]interface{}, error) {

	if !strict && !Exists(name) && !strings.HasSuffix(name, ".json") {
		plus := name + ".json"
		if Exists(plus) {
			name = plus
		}
	}

	raw, err := ioutil.ReadFile(name)
	if err != nil {
		//debug.Trace("*** Failed to locate file [" + name + "], error: " + err.Error())
		// os.Exit(1)
		return nil, err
	}
	// annotate("Read file: " + name)

	var c map[string]interface{}
	if err = json.Unmarshal(raw, &c); err != nil {
		//debug.Trace("*** Failed to process file [" + name + "], error: " + err.Error())
		c = nil
	}
	return c, err
}

func FindArgumentByPrefix(prefix string) (string, error) {

	l := len(prefix)
	if len(os.Args) > 1 && l > 1 {
		for _, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, prefix) {
				return arg[l:], nil
			}
		}
	}
	return "", errors.New("Argument with prefix [" + prefix + "] not found.")
}
