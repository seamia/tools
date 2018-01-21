// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debug

import (
	"os"
	"fmt"
	"bufio"
	"time"
	"strconv"
	"path/filepath"
)

var (
	enabled = false
	logFileName = ""
)

func Trace(format string, a ...interface{}) {
	if enabled {
		f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()

			w := bufio.NewWriter(f)
			fmt.Fprintf(w, format, a...)
			w.Flush()
		}
	}
}

// expands environment variables + a few additional ones (e.g. pid, app-name)
func mapping(src string) string {
	dst := os.Getenv(src)
	if len(dst) == 0 {
		if src == "PID" {
			dst = strconv.FormatInt(int64(os.Getpid()), 10)
		} else if src == "APPNAME" {
			dst = filepath.Base(os.Args[0])
		}
	}
	return dst
}

func init() {
	enabled = true
	name := os.Getenv("LogFile")

	if len(name) == 0 {
		name = os.Args[0] + ".log"
	} else if name == "<none>" {
		enabled = false
	} else {
		name = os.Expand(name, mapping)
	}

	if enabled {
		logFileName = name
		Trace("\n\n------------ app: %s:%d ------------ %s\n",
			os.Args[0],
			os.Getpid(),
			time.Now().Format(time.RFC850))
	}
}
