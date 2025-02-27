// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debug

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/seamia/tools/support"
)

type traceMediaInfo int

const (
	disabled traceMediaInfo = iota
	filename
	stdout
	stderr
)

var (
	media       = disabled
	logFileName = ""
)

func Trace(format string, a ...interface{}) {

	if media == disabled {
		return
	}

	var where io.Writer

	if media == filename {
		f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {

			w := bufio.NewWriter(f)
			where = w
			defer func() {
				w.Flush()
				f.Close()
			}()
		}
	} else if media == stdout {
		where = os.Stdout
	} else if media == stderr {
		where = os.Stderr
	}

	if where != nil {
		fmt.Fprintf(where, format+"\n", a...)
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
		} else if src == "PROGNAME" {
			dst = os.Args[0]
		}
	}
	return dst
}

func processName(name string, first bool) {

	if len(name) == 0 {
		media = disabled
	} else if name == ".log" {
		media = filename
		logFileName = os.Args[0] + ".log"
	} else if name == "(none)" {
		media = disabled
	} else if name == "stdout" {
		media = stdout
	} else if name == "stderr" {
		media = stderr
	} else {
		name = os.Expand(name, mapping)
		if first {
			media = disabled
			processName(name, false)
		} else {
			media = filename
			logFileName = name
		}
	}
}

func init() {

	if name, err := support.FindArgumentByPrefix("/log:"); err == nil {
		// found among args
		processName(name, true)
	} else {
		processName(os.Getenv("LogFile"), true)
	}

	if media != disabled {
		Trace("\n\n------------ app: %s:%d ------------ %s",
			os.Args[0],
			os.Getpid(),
			time.Now().Format(time.RFC850))
	}
}
