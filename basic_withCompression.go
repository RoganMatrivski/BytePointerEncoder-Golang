package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"io/ioutil"
	"math/rand"
	"os"
	"unsafe"
)

func basic_withCompression() {
	// Generate key file
	key := make([]byte, KEY_LENGTH)
	for true {
		rand.Read(key)
		if checkBytesIfUsable(key) {
			break
		}
	}

	// Hash the key file for verifying the keys
	keyHash := sha256.Sum256(key)

	// Find the int variable byte size. Important for decoding later.
	var intSize int
	variableLength = byte(unsafe.Sizeof(intSize))

	// Load the file to encode to memory
	sourceBytes, err := ioutil.ReadFile("./testEncode.bin")
	check(err)

	// Create an array to hold the result
	// Array length breakdown : SHA256 length + null byte + address length + null byte + encoded bytes
	targetBytes := make([]byte, 32+3+len(sourceBytes)*int(variableLength))

	// Write the hash to array
	for i, hashBytes := range keyHash {
		targetBytes[i] = hashBytes
	}

	// Write the int byte length to array
	targetBytes[31+2] = variableLength

	// Set a pointer for writing to array
	pointer := 31 + 4

	// For each byte of the file
	for _, data := range sourceBytes {
		// Get all the candidates of address that have a value of the byte
		addressCandidates := getAddressesWithValue(key, data)

		// Pick a random address
		randomAddress := addressCandidates[r.Intn(len(addressCandidates))]

		// Convert it to a byte array
		randomAddressByteArray := addressToByteArray(int(randomAddress))

		// Write it to array
		for offset, dataByte := range randomAddressByteArray {
			targetBytes[pointer+offset] = dataByte
		}

		// Increment the pointer
		pointer += len(randomAddressByteArray)
	}

	targetFile, err := os.Create("./encodedFile_compressed.bin")
	check(err)
	gzipWriter := gzip.NewWriter(targetFile)

	// Write all binaries to file
	gzipWriter.Write(targetBytes)
	ioutil.WriteFile("./keyFile.bin", key, 0644)

	gzipWriter.Close()

	// =======================================================

	// Load the encoded file
	encodedFileStream, err := os.Open("./encodedFile_compressed.bin")
	check(err)

	gzipReader, err := gzip.NewReader(encodedFileStream)
	var encodedBuffer bytes.Buffer
	_, err = encodedBuffer.ReadFrom(gzipReader)
	check(err)
	encodedFile := encodedBuffer.Bytes()

	// Load the key file
	keyFile, err := ioutil.ReadFile("./keyFile.bin")
	check(err)

	// Check the hash on encoded file with the hash of the key file
	for i, hashByte := range sha256.Sum256(keyFile) {
		if hashByte != encodedFile[i] {
			panic("Hash mismatch")
		}
	}

	// Set the length of the address bytes
	variableLength = encodedFile[31+2]

	// Create a container for the decoded bytes
	res := make([]byte, (len(encodedFile)-35)/int(variableLength))

	// For each byte slot
	for i := range res {
		// Get the array of bytes containing the address of the byte
		addressByteArray := make([]byte, variableLength)
		for j := range addressByteArray {
			addressByteArray[j] = encodedFile[35+i*len(addressByteArray)+j]
		}

		// Convert the byte array into int
		address := byteArrayToInt(addressByteArray)

		// Convert the address to decoded byte
		res[i] = keyFile[address]
	}

	// Write the decoded file
	ioutil.WriteFile("./decodedFile.bin", res, 0644)
}
