/*
Copyright (c) 2019 Dell Inc, or its subsidiaries.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"

	log "github.com/akutz/gournal"
)

func isBinOctetBody(h http.Header) bool {
	return h.Get(headerKeyContentType) == headerValContentTypeBinaryOctetStream
}

func logRequest(ctx context.Context, w io.Writer, req *http.Request, verbose VerboseType) {
	fmt.Fprintln(w, "")
	fmt.Fprint(w, "    -------------------------- ")
	fmt.Fprint(w, "GOISILON HTTP REQUEST")
	fmt.Fprintln(w, " -------------------------")

	switch verbose {
	case Verbose_Low:
		// minimal logging, i.e. print request line only
		fmt.Fprintf(w, "    %s %s %s\r\n", req.Method, req.URL.RequestURI(), req.Proto)
	default:
		// full logging, i.e. print full request message content
		buf, _ := httputil.DumpRequest(req, !isBinOctetBody(req.Header))
		decodedBuf := encryptPassword(buf)
		WriteIndented(w, decodedBuf)
		fmt.Fprintln(w)
	}
}

func logResponse(ctx context.Context, res *http.Response, verbose VerboseType) {
	w := &bytes.Buffer{}

	fmt.Fprintln(w)
	fmt.Fprint(w, "    -------------------------- ")
	fmt.Fprint(w, "GOISILON HTTP RESPONSE")
	fmt.Fprintln(w, " -------------------------")

	var buf []byte

	switch verbose {
	case Verbose_Low:
		// minimal logging, i.e. pirnt status line only
		fmt.Fprintf(w, "    %s %s\r\n", res.Proto, res.Status)
	case Verbose_Medium:
		// print status line + headers
		buf, _ = httputil.DumpResponse(res, false)
	default:
		// print full response message content
		buf, _ = httputil.DumpResponse(res, !isBinOctetBody(res.Header))
	}

	// when DumpResponse gets err, buf will be nil. No message content will be printed
	WriteIndented(w, buf)

	log.Debug(ctx, w.String())
}

// WriteIndentedN indents all lines n spaces.
func WriteIndentedN(w io.Writer, b []byte, n int) error {
	s := bufio.NewScanner(bytes.NewReader(b))
	if !s.Scan() {
		return nil
	}
	l := s.Text()
	for {
		for x := 0; x < n; x++ {
			if _, err := fmt.Fprint(w, " "); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(w, l); err != nil {
			return err
		}
		if !s.Scan() {
			break
		}
		l = s.Text()
		if _, err := fmt.Fprint(w, "\n"); err != nil {
			return err
		}
	}
	return nil
}

// WriteIndented indents all lines four spaces.
func WriteIndented(w io.Writer, b []byte) error {
	return WriteIndentedN(w, b, 4)
}

func encryptPassword(buf []byte) []byte {
	sc := bufio.NewScanner(bytes.NewReader(buf))
	ou := &bytes.Buffer{}
	var l string

	for sc.Scan() {
		match := ""
		separator := ""
		l = sc.Text()

		if strings.Contains(l, "Authorization: Basic") {
			base64str := strings.Split(l, " ")[2]
			decoded, _ := base64.StdEncoding.DecodeString(base64str)
			decodedName := strings.Split(string(decoded), ":")[0]
			l = "Authorization: " + decodedName + ":******"
		} else {
			switch {
			case strings.Contains(l, "password"):
				match = `"password":"`
				separator = `"`
			case strings.Contains(l, "Cookie: isisessid="):
				match = `isisessid=`
				separator = `-`
			case strings.Contains(l, "X-Csrf-Token"):
				match = `X-Csrf-Token:`
				separator = `-`
			}
			if match != "" {
				startIndex, endIndex, matchStrLen := FetchValueIndexForKey(l, match, separator)
				if startIndex >= 0 && endIndex > 0 { // if the separator is present then replace only the characters till separator with the special character
					l = l[:startIndex+matchStrLen] + "****" + l[startIndex+matchStrLen+endIndex:]
				} else if startIndex >= 0 { // if the separator in not present then replace the string to be masked with the special character
					l = l[:startIndex+matchStrLen] + "****"
				}
			}
		}
		fmt.Fprintln(ou, l)
	}

	return ou.Bytes()
}
