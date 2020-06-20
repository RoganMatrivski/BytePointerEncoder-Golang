package main

import (
	"crypto/sha256"
	"io/ioutil"
)

func main() {
	basic()
	// basic_withCompression()
	// basic_bufio()

	a, _ := ioutil.ReadFile("./testEncode.bin")
	b, _ := ioutil.ReadFile("./decodedFile.bin")

	aHash := sha256.Sum256(a)
	bHash := sha256.Sum256(b)

	for i, aData := range aHash {
		if aData != bHash[i] {
			panic("Result file hash mismatch")
		}
	}
}
