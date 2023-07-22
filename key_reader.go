package main

import (
	"bytes"
	_ "embed"
)

//go:embed pub.txt
var pub []byte

//go:embed priv.txt
var priv []byte

var (
	Pub  = bytes.NewReader(pub)
	Priv = bytes.NewReader(priv)
)
