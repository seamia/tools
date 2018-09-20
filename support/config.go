// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package support

import (
	"encoding/json"
	"errors"
	"github.com/seamia/tools/assets"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

const configRuntimeBranch = "runtime"

func bytes2interface(raw []byte) (c interface{}, err error) {
	if err = json.Unmarshal(raw, &c); err != nil {
		assert("*** Failed to process file json data, error: " + err.Error())
		c = nil
	}
	return c, err
}

func loadFromJson(name string) (interface{}, error) {
	raw, err := ioutil.ReadFile(name)
	if err != nil {
		assert("*** Failed to locate/read file [" + name + "], error: " + err.Error())
		return nil, err
	}

	return bytes2interface(raw)
}

func loadConfig(fullname string) (config map[string]interface{}, err error) {
	if back, err := loadFromJson(fullname); err != nil {
		return nil, err
	} else {
		config = back.(map[string]interface{})
	}

	if config != nil {
		runtime := make(map[string]string)

		runtime["config.name"] = fullname

		if value, err := filepath.Abs(filepath.Dir(fullname)); err == nil {
			runtime["config.dir"] = value
			os.Setenv("config.dir", value)
		}
		if value, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
			runtime["app.dir"] = value
			os.Setenv("app.dir", value)
		}
		if value, err := os.Getwd(); err == nil {
			runtime["start.dir"] = value
		}

		if _, found := config[configRuntimeBranch]; !found {
			config[configRuntimeBranch] = runtime
		}
	}
	return config, nil
}

func LoadConfig(name string, fallbackToAssetsOnFailure bool) (map[string]interface{}, error) {

	if assets.IsAssetFile(name) {
		var err error
		if reader, err := assets.Open(name); err == nil {
			if data, err := ioutil.ReadAll(reader); err == nil {
				if raw, err := bytes2interface(data); err == nil {
					return raw.(map[string]interface{}), nil
				} else {
					assert("failed to parse asset:" + name)
				}
			} else {
				assert("failed to read asset:" + name)
			}
		} else {
			assert("failed to open asset:" + name)
		}
		return nil, err
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		fullname := filepath.Join(dir, name)
		if Exists(fullname) { // let's look into locations related to the location of the app
			return loadConfig(fullname)
		} else if Exists(name) { // let's test the path as we're given (in case it is an absolute one, or ...)
			return loadConfig(name)
		}

		if fallbackToAssetsOnFailure {
			return LoadConfig(assets.AssetUriPrefix + name, false)
		}

		err = errors.New("failed to find file [" + fullname + "] or [" + name + "].")
	}
	return nil, err
}

func GetLocation(config map[string]interface{}, what string) (string, error) {
	if locs, present := config["locations"]; present {
		locations := locs.(map[string]interface{})
		if locations != nil {
			if value, found := locations[what]; found {
				evalue := os.ExpandEnv(value.(string))
				return evalue, nil
			} else {
				return "", errors.New("Warning: Location setting (" + what + ") was not found in the provided config file.")
			}
		} else {
		}
	} else {
	}
	return "", errors.New("Warning: Location settings were not found in the provided config file.")
}

func init() {
	if runtime.GOOS == "windows" {
		if os.Getenv("HOME") == "" {
			os.Setenv("HOME", os.Getenv("USERPROFILE")) // making 'HOME' env var available everywhere
		}
	}
}
