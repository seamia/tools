// Copyright 2020 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type EmailConfig struct {
	Username string `json:"user"`
	Password string `json:"pass"`
	SMTP     string `json:"smtp"`
	Port     string `json:"port"`
	Message  string `json:"message"`
	From     string `json:"from"`
	Subject  string `json:"subject"`
}

func (ec EmailConfig) Address() string {
	return ec.SMTP + ":" + ec.Port
}

var (
	emailConfig EmailConfig
)

func sendEmail(to []string, txt string) {
	if len(emailConfig.Username) == 0 {
		if data, err := os.ReadFile(emailConfigFileName); err == nil {
			if err := json.Unmarshal(data, &emailConfig); err != nil {
				// no point to go on, since we have no config
				fmt.Println("failed to load config")
				return
			}
		}
	}

	message := emailConfig.Message
	if len(txt) > 0 {
		message = txt
	}

	msg := []byte("To:" + strings.Join(to, ";") +
		"\r\nFrom: " + emailConfig.From +
		"\r\nSubject: " + emailConfig.Subject +
		"\r\nContent-Type: text/plain\r\n\r\n" +
		message)

	// Authentication.
	auth := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.SMTP)

	// Sending email.
	err := smtp.SendMail(emailConfig.Address(), auth, emailConfig.Username, to, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}

func sendTestEmail(to []string) {
	msg := "if you are getting this ... the mail delivery is working"
	sendEmail(to, msg)
	fmt.Fprintf(os.Stdout, "sending test emails...\n")
}
