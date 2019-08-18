// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/seamia/libs/printer"
)

func call(relativeUrl, verb, data string) {

	u, err := url.Parse(expand(baseUrl))
	quitOnError(err, "Parsing url [%s]", baseUrl)
	u.Path = path.Join(u.Path, expand(relativeUrl))
	fullUrl := u.String()

	if generateCurlCommands {
		produceCurlCommand(fullUrl, verb, data)
	} else {
		data = loadExternalFile(data)
		var payload io.Reader
		if len(data) > 0 {
			payload = bytes.NewReader([]byte(data))
		}

		client := &http.Client{}
		request, err := http.NewRequest(strings.ToUpper(verb), fullUrl, payload)

		quitOnError(err, "...")

		for key, value := range headers {
			if len(key) > 0 && len(value) > 0 {
				request.Header.Set(key, expand(value))
			}
		}
		request.Header.Set("User-Agent", userAgent)

		start := time.Now()
		resp, err := client.Do(request)
		if collectTimingInfo {
			duration := time.Now().Sub(start)
			response("the request took %s", duration.String())
		}
		quitOnError(err, "......")

		displayResponse(resp)
	}
}

func displayResponse(resp *http.Response) {
	if resp == nil {
		response("got an empty response")
	}

	print := responseFailure
	if resp.StatusCode < http.StatusBadRequest {
		print = responseSuccess
	}

	// colorPrint(colorResponse, format, a...)

	print("Status: %s", resp.Status)
	displayHeaders(resp, print)

	if resp.Body != nil {
		data, err := ioutil.ReadAll(resp.Body)
		quitOnError(err, "Ingesting response body")

		switch getContentType(resp) {
		case contentTypeJson:
			displayJsonBody(data, print)
		default:
			displayPlainBody(data, print)
		}
	}
}

func displayHeaders(resp *http.Response, print printer.Printer) {
	const format = "\tHeader: [%s] = [%s]"
	if printResponseHeaders && len(resp.Header) > 0 {
		for key := range resp.Header {
			value := resp.Header.Get(key)
			if attentionNeeded(key) {
				responseAttention(format, key, value)
			} else {
				print(format, key, value)
			}
		}
	}
}

func getContentType(resp *http.Response) string {
	return lower(resp.Header.Get(headerContentType))
}

func displayPlainBody(data []byte, print printer.Printer) {
	if len(data) == 0 {
		print("Body is empty.")
	} else {
		print("Body: %s", string(data))
	}
}

func displayJsonBody(data []byte, print printer.Printer) {
	if !responsePrettyPrintBody || len(data) == 0 {
		displayPlainBody(data, print)
		return
	}

	var blank interface{}
	blank = []interface{}{}
	if err := json.Unmarshal(data, &blank); err == nil {
		if pretty, err := json.MarshalIndent(blank, "", "    "); err == nil {
			displayPlainBody(pretty, print)
			return
		} else {
			reportError(err, "marshalling")
		}
	} else {
		reportError(err, "unmarshalling")
	}

	blank = map[string]interface{}{}
	if err := json.Unmarshal(data, &blank); err == nil {
		if pretty, err := json.MarshalIndent(blank, "", "    "); err == nil {
			displayPlainBody(pretty, print)
			return
		} else {
			reportError(err, "marshalling")
		}
	} else {
		reportError(err, "unmarshalling")
	}
	displayPlainBody(data, print)
}

func attentionNeeded(key string) bool {
	if strings.HasSuffix(lower(key), headerAttentionSuffix) {
		return true
	}
	return false
}
