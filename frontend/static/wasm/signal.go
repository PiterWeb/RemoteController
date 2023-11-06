package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("signalEncode", js.FuncOf(signalEncode))
	js.Global().Set("signalDecode", js.FuncOf(signalDecode))

	<-c
}

// signalEncode encodes the input in base64
// It can optionally zip the input before encoding
func signalEncode(this js.Value, objStr []js.Value) interface{} {

	obj := new(interface{})

	json.Unmarshal([]byte(objStr[0].String()), obj)

	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	b = signalZip(b)

	return js.ValueOf(base64.StdEncoding.EncodeToString(b))
}

// signalDecode decodes the input from base64
// It can optionally unzip the input after decoding
func signalDecode(this js.Value, in []js.Value) interface{} {
	b, err := base64.StdEncoding.DecodeString(in[0].String())
	if err != nil {
		panic(err)
	}

	b = signalUnzip(b)

	obj := new(interface{})

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}

	b, err = json.Marshal(obj)

	if err != nil {
		panic(err)
	}

	return js.ValueOf(string(b))
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
