// Copyright 2020 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	emailConfigFileName = "./email.config"
	tmpFolderDefaultName = "./tmp.%v"
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

	tmpFolder := createTmpFolder(get(config, "tmp.folder", tmpFolderDefaultName))
	report("--------------------------- start. script: %v. time: %v", script, time.Now())
	defer func() {
		report("--------------------------- stop. time: %v", time.Now())
		removeTmpFolder(tmpFolder)
	}()

	list := loadList(script)

	media := new(MediaImpl)
	media.Printf("using script: %s\n", script)
	media.Printf("(%s)\n", time.Now().Format(time.DateTime))

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
		process(task, media)
	}

	if getFlag(config, "send.summary.email", false) {
		sendEmail([]string{defaultEmail}, media.Get())
	}
}

var (
	// DownloadHTML(what string) (string, error)
	// GetContent(from string) (string, error)

	download func(string) (string, error)
)

func process(dict msi, media Media) {

	download = DownloadHTML
	name := dict["name"].(string)

	trace := func(format string, args ...interface{}) {
		media.Printf("entry [%s]: ", name)
		media.Printf(format+"\n", args...)
	}

	if data, found := dict["active"]; found {
		if active, converts := data.(bool); converts {
			if !active {
				trace("marked as not active")
				report("\tskipping inactive: %s", name)
				return
			}
		}
	}

	action := get(dict, "action", defaultEmail)

	report("checking: %s", name)

	from := dict["url"].(string)

	if status64, found := dict["status"].(float64); found && status64 > 0 {
		status := int(status64)
		if req, err := http.NewRequest("GET", from, nil); err != nil {
			report("\tgot error: %v", err)
		} else {
			// Execute the request.
			txt := ""
			if resp, err := http.DefaultClient.Do(req); err != nil {
				report("\tgot error: %v", err)
				txt = fmt.Sprintf("%s. error: %v", name, err)
			} else {

				// Close response body as required.
				defer resp.Body.Close()

				report("\tgot status: %v", resp.StatusCode)
				if resp.StatusCode != status {
					txt = fmt.Sprintf("%s. instead of expected status (%v) got: %v", name, status, resp.StatusCode)
				}
			}

			if len(txt) > 0 {
				alert(txt, action)
			}
		}
		return
	}

	fullContent, err := download(from)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to download from (%s) with error: %v\n", from, err)
		trace("failed to download from (%s) with error: %v\n", from, err)
		return
	}

	saveContent("always-"+name, fullContent)

	content := trimPrefix(fullContent, dict["after"])
	content = trimSuffix(content, dict["before"])
	if len(content) < len(fullContent) {
		saveContent("trimmed-"+name, content)
	}

	if known, found := dict["known"].(string); found && len(known) > 0 {
		if strings.Compare(content, known) != 0 {
			saveContent(name, fullContent)
			report("\tBINGO: %s. found %s instead of %s", name, content, known)
			// txt := fmt.Sprintf("%s. changed from %s to %s", name, known, content)

			values := map[string]string{
				"Entry":    name,
				"Source":   get(dict, "url", ""),
				"Expected": known,
				"Actual":   content,
			}

			// use "content" instead of "known", so that any future changes (of "content") will result in a new signature
			alertAck(action, name, content, "known.html", values)
		} else {
			trace("current value: %s", known)
			report("\t--- known is still there (%s)", known)
		}
	} else if missing, found := dict["missing"].(string); found && len(missing) > 0 {
		if !strings.Contains(content, missing) {
			saveContent(name, fullContent)
			report("\tFOUND: %s", name)
			alert(name, action)
		} else {
			report("\t--- missing is still there (%s)", missing)
		}
	} else if present, found := dict["present"].(string); found && len(present) > 0 {
		if strings.Contains(content, present) {
			report("\tFOUND: %s", name)
			saveContent(name, fullContent)
			alert(name, action)
		} else {
			report("\t--- present is still not there (%s)", present)
		}
	} else {
		report("\tERROR: missing any recognisable operations...")
	}
}
