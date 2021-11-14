package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

type Reader struct {}

func New() *Reader {
	return &Reader{}
}

func (r *Reader) FromFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func (r *Reader) FromStdio() string {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func (r *Reader) Read() string {

	filePath := flag.String("f", "", "Path to json file")
	flag.Parse()

	if *filePath != "" {
		return r.FromFile(*filePath)
	} else {
		return r.FromStdio()
	}
}

