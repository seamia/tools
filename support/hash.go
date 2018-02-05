// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package support

import (
	"crypto/sha1"
	"encoding/hex"
)

func Hash(raw []byte) string {
	h := sha1.New()
	h.Write(raw)
	return hex.EncodeToString(h.Sum(nil))
}
