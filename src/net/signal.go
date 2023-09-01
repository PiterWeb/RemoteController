// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

// Package signal contains helpers to exchange the SDP session
// description between examples.
package net

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// Allows compressing offer/answer to bypass terminal input limits.
const compress = false

// signalMustReadStdin blocks until input is received from stdin
func signalMustReadStdin() string {
	r := bufio.NewReader(os.Stdin)

	var in string
	for {
		var err error
		in, err = r.ReadString('\n')
		if err != io.EOF {
			if err != nil {
				panic(err)
			}
		}
		in = strings.TrimSpace(in)
		if len(in) > 0 {
			break
		}
	}

	fmt.Println("")

	return in
}

// signalEncode encodes the input in base64
// It can optionally zip the input before encoding
func signalEncode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	if compress {
		b = signalZip(b)
	}

	return base64.StdEncoding.EncodeToString(b)
}

// signalDecode decodes the input from base64
// It can optionally unzip the input after decoding
func signalDecode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	if compress {
		b = signalUnzip(b)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}

func signalZip(in []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)
	if err != nil {
		panic(err)
	}
	err = gz.Flush()
	if err != nil {
		panic(err)
	}
	err = gz.Close()
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func signalUnzip(in []byte) []byte {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		panic(err)
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		panic(err)
	}
	res, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return res
}
