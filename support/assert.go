// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package support

import (
	"fmt"
)

func assert(msg string) {
	// todo: remove the output
	fmt.Println("Assert: " + msg)
}
