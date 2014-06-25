package main

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
)

var hashes = make(map[string]string)

func calcFileHash(filename string) string {
	bytes, _ := ioutil.ReadFile(filename)

	h := fnv.New32()
	h.Write(bytes)

	return fmt.Sprint(h.Sum32())
}

func (*TmplData) Bust(path string) string {
	filename := path

	if filename[0] == '/' {
		filename = filename[1:]
	}

	if _, found := hashes[filename]; !found {
		hashes[filename] = fmt.Sprintf("%s?v=%s", path, calcFileHash(filename))
	}

	return hashes[filename]
}
