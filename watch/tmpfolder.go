package main

import (
	"fmt"
	"os"
	"strings"
)

func createTmpFolder(format string) string {
	name := format
	if strings.Contains(format, "%") {
		name = fmt.Sprintf(format, os.Getpid())
	}

	if err := os.MkdirAll(name, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create folder (%s) due to: %v\n", name, err)
		return ""
	}

	os.Setenv("TMPDIR", name)
	report("using (%s) as tmp folder", name)
	return name
}

func removeTmpFolder(name string) {
	if len(name) == 0 {
		return
	}

	if err := os.RemoveAll(name); err != nil {
		fmt.Fprintf(os.Stderr, "failed to remove folder (%s) due to: %v\n", name, err)
	}
}
