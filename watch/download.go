// Copyright 2020 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

func DownloadFile(filepath string, url string) error {
	fmt.Println("downloading:", url, "--->>", filepath)

	// Get the data
	// resp, err := http.Get(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("user-agent", userAgent)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	// _, err = io.Copy(out, resp.Body)

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{All: uint64(resp.ContentLength)}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))

	fmt.Println("done, err:", err)
	return err
}

func GetContent(from string) (string, error) {

	if got, err := chromeGetContent(from); err == nil {
		return got, err
	}

	j, _ := cookiejar.New(nil)
	getter := &http.Client{Jar: j}

	// Get the data
	resp, err := getter.Get(from)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	return string(data), err
}

func findInBetween(src, left, right string, unique bool) []string {
	results := make([]string, 0)

	for {
		l := strings.Index(src, left)
		if l < 0 {
			break
		}
		src = src[l+len(left):]
		r := strings.Index(src, right)
		if r < 0 {
			break
		}
		found := src[:r]
		src = src[r+len(right):]
		results = append(results, found)
	}

	if unique && len(results) > 1 {
		sort.Strings(results)
		distinct := make([]string, 0, len(results))
		for index, value := range results {
			if index == 0 || value != results[index-1] {
				distinct = append(distinct, value)
			}
		}
		return distinct
	}

	return results
}

var (
	errWrongOS = errors.New("this OS is not supported")
)

// google-chrome --headless --disable-gpu --dump-dom --user-agent="Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.50 Safari/537.36" https://www.amd.com/en/direct-buy/us > amd

func chromeGetContent(from string) (string, error) {
	if runtime.GOOS == "windows" {
		return "", errWrongOS
	}

	args := []string{
		"--headless",
		"--disable-gpu",
		"--dump-dom",
		"--user-agent=\"" + userAgent + "\"",
		from,
	}

	if output, err := exec.Command("google-chrome", args...).Output(); err == nil {
		return string(output), nil
	} else {
		return "", err
	}
}
