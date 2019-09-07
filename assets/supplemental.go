// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package assets

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type AssetInfo struct {
	// name  string
	Data  string
	Mime  string
	Mtime int64
	Size  int64
	Hash  string
}

type AssetRoot map[string]*AssetInfo

const AssetUriPrefix = "assets://"

var g_staticAssets AssetRoot = nil

func Assign(assets AssetRoot) {
	if g_staticAssets == nil || len(g_staticAssets) == 0 {
		g_staticAssets = assets
	} else {
		for k, v := range assets {
			g_staticAssets[k] = v
		}
	}
}

func Open(name string) (io.ReadCloser, error) {

	// normalize the name (remove Uri prefix)
	if strings.HasPrefix(name, AssetUriPrefix) {
		name = name[len(AssetUriPrefix):]
	}

	f, ok := g_staticAssets[name]
	if !ok {
		return nil, fmt.Errorf("Asset %s not found", name)
	}

	if f.Size == 0 {
		return ioutil.NopCloser(strings.NewReader(f.Data)), nil
	}
	return gzip.NewReader(strings.NewReader(f.Data))
}

func OpenSeeker(name string, expansion interface{}) (io.ReadSeeker, error) {

	// normalize the name (remove Uri prefix)
	if strings.HasPrefix(name, AssetUriPrefix) {
		name = name[len(AssetUriPrefix):]
	}

	f, ok := g_staticAssets[name]
	if !ok {
		return nil, fmt.Errorf("Asset %s not found", name)
	}

	var data []byte = []byte(f.Data)
	if f.Size > 0 {
		reader, err := gzip.NewReader(strings.NewReader(f.Data))
		if err != nil {
			return nil, err
		}
		data = make([]byte, f.Size)
		read, err := reader.Read(data)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if int64(read) != f.Size {
			return nil, fmt.Errorf("failed to decompress expected data size. got %v instead of %v", read, f.Size)
		}
	}

	var err error
	if data, err = substitute(data, name, expansion); err != nil {
		return nil, err
	}

	return CreateReader(data)
}

func substitute(data []byte, name string, expansion interface{}) ([]byte, error) {
	if expansion == nil {
		// nothing to expand with
		return data, nil
	}

	// Create a new template and parse the letter into it.
	templ, err := template.New(name).Parse(string(data))
	if err != nil {
		return nil, err
	}

	var rendered bytes.Buffer
	if err := templ.Execute(&rendered, expansion); err != nil {
		return nil, err
	}

	return rendered.Bytes(), nil
}

func exists(filePath string) (exists bool) {
	exists = true

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		exists = false
	}

	return
}

func ExtractAssets(destination string, overwriteExisting bool) error {

	// 0. If not 'destination' was specified - let's assume the app's location
	if len(destination) == 0 {
		if appDir, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
			destination = appDir
		} else {
			return err
		}
	}

	// 1. check if any of the potential destinations are already present
	if !overwriteExisting {
		alreadyExist := make([]string, 0)
		for name, _ := range g_staticAssets {
			fullname := path.Join(destination, name)
			if root, err := filepath.Abs(fullname); err == nil {
				if exists(root) {
					alreadyExist = append(alreadyExist, fullname)
				}
			}
		}

		if len(alreadyExist) > 0 {
			return errors.New(fmt.Sprint("At least some of the files will be overwritten. The following file(s) are already present: ", alreadyExist))
		}
	}

	createdFiles := make([]string, 0, len(g_staticAssets))
	doNotRollbackChanges := false

	defer func() { // 'rollback' in case of a failure
		if !overwriteExisting && !doNotRollbackChanges && createdFiles != nil && len(createdFiles) > 0 {
			for _, one := range createdFiles {
				if err := os.Remove(one); err != nil {
					// looks like there was an error during the attemp to delete a file. not much we can do about it at this point
					_ = err
				}
			}
		}
	}()

	for name, _ := range g_staticAssets {
		fullname := path.Join(destination, name)
		if root, err := filepath.Abs(filepath.Dir(fullname)); err == nil {
			if !exists(root) {
				if err := os.MkdirAll(root, 0775); err != nil {
					return err
				}
			}

			out, err := os.Create(fullname)
			if err == nil {
				createdFiles = append(createdFiles, fullname)
				defer out.Close()
			} else {
				return err
			}

			if source, err := Open(name); err == nil {
				defer source.Close()
				if _, err := io.Copy(out, source); err != nil {
					return err
				}
			} else {
				return err
			}

		} else {
			return err
		}
	}

	doNotRollbackChanges = true
	return nil
}

func IsAssetFile(name string) bool {

	if g_staticAssets != nil && len(g_staticAssets) > 0 {
		if strings.HasPrefix(name, AssetUriPrefix) {
			assetName := name[len(AssetUriPrefix):]
			if _, found := g_staticAssets[assetName]; found {
				return true
			}
		}
	}
	return false
}
