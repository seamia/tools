// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package support

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/seamia/tools/assets"
)

const (
	configRuntimeBranch   = "runtime"
	configLocationsBranch = "locations"
	separator             = ";"
)

type (
	msi = map[string]interface{}
)

func bytes2interface(raw []byte) (c interface{}, err error) {
	if err = json.Unmarshal(raw, &c); err != nil {
		assert("*** Failed to process file json data, error: " + err.Error())
		c = nil
	}
	return c, err
}

func loadFromJson(name string) (interface{}, error) {
	raw, err := os.ReadFile(name)
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

	if components := strings.Split(name, separator); len(components) > 1 {
		return loadConfigs(components, fallbackToAssetsOnFailure)
	}

	if assets.IsAssetFile(name) {
		var err error
		if reader, err := assets.Open(name); err == nil {
			if data, err := io.ReadAll(reader); err == nil {
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
			return LoadConfig(assets.AssetUriPrefix+name, false)
		}

		err = errors.New("failed to find file [" + fullname + "] or [" + name + "].")
	}
	return nil, err
}

func GetLocation(config map[string]interface{}, what string) (string, error) {
	if locs, present := config[configLocationsBranch]; present {
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

func SetLocation(config map[string]interface{}, key, value string) {
	if config != nil {
		if locs, present := config[configLocationsBranch]; !present {
			config[configLocationsBranch] = map[string]string{key: value}
		} else {
			locations := locs.(map[string]interface{})
			if locations != nil {
				locations[key] = value
			}
		}
	}
}

func loadConfigs(names []string, fallbackToAssetsOnFailure bool) (msi, error) {

	if runtime.GOOS == "windows" {
		for i, name := range names {
			names[i] = strings.ReplaceAll(name, `\`, `/`)
		}
	}

	directory, _ := path.Split(names[0])

	var result msi
	for _, name := range names {

		if len(directory) > 0 {
			dirname, filename := path.Split(name)
			if len(dirname) == 0 {
				name = path.Join(directory, filename)
			}
		}

		config, err := LoadConfig(name, fallbackToAssetsOnFailure)
		if err != nil {
			return nil, err
		}
		result = mergeConfigs(result, config)
	}

	if len(result) > 0 {
		if rtime := result["runtime"]; rtime != nil {
			if values := rtime.(map[string]string); len(values) > 0 {
				if len(values["config.name"]) > 0 {
					values["config.name"] = strings.Join(names, separator)
				}
			}

		}
	}
	return result, nil
}

func mergeConfigs(left, right msi) msi {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	for key, value := range right {
		if lvalue, exist := left[key]; exist {
			switch actual := value.(type) {
			case map[string]interface{}:
				left[key] = mergeConfigs(lvalue.(msi), actual)

			case map[string]string:
				lmap := lvalue.(map[string]string)
				for k, v := range actual {
					lmap[k] = v
				}
				left[key] = lmap

			case string, bool, int:
				left[key] = actual

			case []interface{}:
				left[key] = append(lvalue.([]interface{}), actual...)

			default:
				panic(fmt.Sprintf("unhandled type %s: %T\n", key, value))
			}

		} else {
			left[key] = value
		}
	}
	return left
}

func init() {
	if runtime.GOOS == "windows" {
		if os.Getenv("HOME") == "" {
			os.Setenv("HOME", os.Getenv("USERPROFILE")) // making 'HOME' env var available everywhere
		}
	}
}
