// Copyright 2020 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var (
	alreadyDownloaded = make(map[string]string)
)

func DownloadHTML(what string) (string, error) {

	if value, found := alreadyDownloaded[what]; found {
		report("\treusing previously downloaded: %s", what)
		return value, nil
	}

	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var text string
	if err := chromedp.Run(ctx, getHtml(what, &text)); err != nil {
		return "", err
	}

	alreadyDownloaded[what] = text
	return text, nil
}

func getHtml(urlstr string, res *string) chromedp.Tasks {

	headers := msi{}
	headers["User-Agent"] = userAgent

	headers["Sec-Ch-Ua"] = `" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`
	headers["Accept-Language"] = `en-US,en;q=0.9`
	headers["Sec-Ch-Ua-Mobile"] = `?0`

	return chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(headers)),

		chromedp.Navigate(urlstr),
		chromedp.OuterHTML("html", res),
	}
}
