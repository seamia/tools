// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package assets

import (
	"net/http"
	"time"
)

func ServeAsset(w http.ResponseWriter, req *http.Request, asset string) error {
	return ServeAssetTemplate(w, req, asset, nil)
}

func ServeAssetTemplate(w http.ResponseWriter, req *http.Request, asset string, values interface{}) error {
	reader, err := OpenSeeker(asset, values)
	if err != nil {
		return err
	}
	http.ServeContent(w, req, asset, time.Now(), reader)
	return nil
}
