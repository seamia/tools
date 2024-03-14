// Copyright 2020 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	emailConfigFileName = "./email.config"
	userAgent           = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36"
	defaultEmail        = ""
	config              msi
)

func main() {
	config = loadDict("./watch.config")
	emailConfigFileName = get(config, "email.config", emailConfigFileName)
	userAgent = get(config, "agent", userAgent)
	defaultEmail = get(config, "email", defaultEmail)
	script := get(config, "script", "./script.json")

	if len(os.Args) > 1 {
		script = os.Args[1]
		report("using %s as script file...", script)
	}

	report("--------------------------- start. script: %v. time: %v", script, time.Now())
	defer func() {
		report("--------------------------- stop. time: %v", time.Now())
	}()

	list := loadList(script)

	// -----------------------------------------------------------------------------------------------------------------
	to := []string{}
	for _, task := range list {
		action := get(task, "action", "")
		if len(action) == 0 {
			if len(defaultEmail) > 0 {
				action = defaultEmail
			} else {
				continue
			}
		}
		to = append(to, action)
	}

	if getFlag(config, "send.test.email", false) {
		sendTestEmail(to)
	}
	// -----------------------------------------------------------------------------------------------------------------

	for _, task := range list {
		process(task)
	}
}

var (
	// DownloadHTML(what string) (string, error)
	// GetContent(from string) (string, error)

	download func(string) (string, error)
)

func process(dict msi) {

	download = DownloadHTML
	name := dict["name"].(string)
	if data, found := dict["active"]; found {
		if active, converts := data.(bool); converts {
			if !active {
				report("skipping inactive: %s", name)
				return
			}
		}
	}

	action := get(dict, "action", defaultEmail)

	report("checking: %s", name)

	from := dict["url"].(string)
	fullContent, err := download(from)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to download from (%s) with error: %v\n", from, err)
		return
	}

	saveContent("always-"+name, fullContent)

	content := trimPrefix(fullContent, dict["after"])
	content = trimSuffix(content, dict["before"])
	saveContent("trimmed-"+name, content)

	if known, found := dict["known"].(string); found && len(known) > 0 {
		if strings.Compare(content, known) != 0 {
			saveContent(name, fullContent)
			report("BINGO: %s. found %s instead of %s", name, content, known)
			txt := fmt.Sprintf("%s. changed from %s to %s", name, known, content)
			alert(txt, action)
		} else {
			report("--- known is still there (%s)", known)
		}
	} else if missing, found := dict["missing"].(string); found && len(missing) > 0 {
		if !strings.Contains(content, missing) {
			saveContent(name, fullContent)
			report("FOUND: %s", name)
			alert(name, action)
		} else {
			report("--- missing is still there (%s)", missing)
		}
	} else if present, found := dict["present"].(string); found && len(present) > 0 {
		if strings.Contains(content, present) {
			report("FOUND: %s", name)
			saveContent(name, fullContent)
			alert(name, action)
		} else {
			report("--- present is still not there (%s)", present)
		}
	} else {
		report("ERROR: no missing not present entry is found...")
	}
}
