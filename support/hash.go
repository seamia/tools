// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package support

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/deatil/go-encoding/base62"
)

func Hash(raw []byte) string {
	h := sha1.New()
	h.Write(raw)
	return hex.EncodeToString(h.Sum(nil))
}

func NameToId(name string, length int) string {
	return Hash([]byte(name))[:length]
}

func Hash62(raw []byte) string {
	encodedString := base62.StdEncoding.EncodeToString(data)
	return encodedString
}