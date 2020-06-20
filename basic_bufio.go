package main

import (
	"bufio"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"unsafe"

	"golang.org/x/exp/errors/fmt"
)

func basic_bufio() {
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

	// Open a file to read
	file, err := os.Open("./testEncode.bin")
	check(err)

	// Create a reader
	fileReader := bufio.NewReader(file)

	// Create a file to write to
	target, err := os.Create("./encodedFile_bufio.bin")
	check(err)

	// Create a writer
	fileWriter := bufio.NewWriter(target)

	// Write the hash to file
	fileWriter.Write(keyHash[:])

	// Write the variable length with padding
	fileWriter.Write([]byte{0, variableLength, 0})

	// Read per chunks
	buffer := make([]byte, 1024)
	for {
		// Load file to buffer
		bytesReadAmount, err := fileReader.Read(buffer)

		// Break from loop if End of File reached
		if err == io.EOF {
			break
		}

		// For each buffers
		for _, data := range buffer[0:bytesReadAmount] {
			// Get all the candidates of address that have a value of the byte
			addressCandidates := getAddressesWithValue(key, data)

			// Pick a random address
			randomAddress := addressCandidates[r.Intn(len(addressCandidates))]

			// Convert it to a byte array
			randomAddressByteArray := addressToByteArray(int(randomAddress))

			// Write it to file
			fileWriter.Write(randomAddressByteArray)
		}
	}

	// Write all binaries to file
	fileWriter.Flush()
	ioutil.WriteFile("./keyFile.bin", key, 0644)

	// =======================================================

	// Load the encoded file
	encodedFile, err := os.Open("./encodedFile_bufio.bin")
	check(err)

	// Create file reader
	encodedFileReader := bufio.NewReader(encodedFile)

	// Load the key file
	keyFile, err := ioutil.ReadFile("./keyFile.bin")
	check(err)

	header := make([]byte, 35)
	_, err = encodedFileReader.Read(header)
	if err == io.EOF {
		panic("Unexpected EOF")
	}

	// Check the hash on encoded file with the hash of the key file
	for i, hashByte := range sha256.Sum256(keyFile) {
		if hashByte != header[i] {
			panic("Hash mismatch")
		}
	}

	// Set the length of the address bytes
	variableLength = header[31+2]

	// Create a container for the decoded bytes
	decodedFile, err := os.Create("./decodedFile.bin")
	check(err)
	decodedFileWriter := bufio.NewWriter(decodedFile)

	buffer2 := make([]byte, 1024)
	byteArrayContainer := make([]byte, int(variableLength))
	t := make([]byte, int(variableLength))
	byteArrayPointer := 0

	for {
		_, err := encodedFileReader.Read(buffer2)

		for _, bufferByte := range buffer2 {
			byteArrayContainer[byteArrayPointer] = bufferByte
			byteArrayPointer++

			if byteArrayPointer == int(variableLength) {
				// Convert the byte array into int
				address := byteArrayToInt(byteArrayContainer)

				// Convert the address to decoded byte
				decodedFileWriter.WriteByte(keyFile[address])
				fmt.Println(t)

				byteArrayPointer = 0
				t = byteArrayContainer
				byteArrayContainer = make([]byte, int(variableLength))
			}
		}

		if err == io.EOF {
			break
		}
	}

	// Write the decoded file
	decodedFileWriter.Flush()
}
