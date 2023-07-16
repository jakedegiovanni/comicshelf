package main

import (
	_ "embed"
	"io"
)

//go:embed pub.txt
var pub []byte

//go:embed priv.txt
var priv []byte

type pubReader struct {
	readIndex int64
}

func (r *pubReader) Read(p []byte) (n int, err error) {
	if r.readIndex >= int64(len(pub)) {
		err = io.EOF
		return
	}

	n = copy(p, pub[r.readIndex:])
	r.readIndex += int64(n)
	return
}

type privReader struct {
	readIndex int64
}

func (r *privReader) Read(p []byte) (n int, err error) {
	if r.readIndex >= int64(len(priv)) {
		err = io.EOF
		return
	}

	n = copy(p, priv[r.readIndex:])
	r.readIndex += int64(n)
	return
}
