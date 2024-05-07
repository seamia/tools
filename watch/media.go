// Copyright 2020 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"
)

type Media interface {
	Printf(format string, args ...interface{})
	Error(err error, comment string)
}

type MediaImpl struct {
	strings.Builder
}

func (m *MediaImpl) Printf(format string, args ...interface{}) {
	m.WriteString(fmt.Sprintf(format, args...))
}

func (m *MediaImpl) Error(err error, comment string) {
	m.Printf("Error: %v; %s\n", err, comment)
}

func (m *MediaImpl) Save(filename string) error {
	return os.WriteFile(filename, []byte(m.String()), 0644)
}

func (m *MediaImpl) Get() string {
	return m.String()
}
