package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io/ioutil"
	"math/rand"
	"time"
	"unsafe"
)

const KEY_LENGTH = 1024

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var variableLength byte

func checkBytesIfUsable(bytes []byte) bool {
	for i := byte(0); i < 255; i++ {
		byteFound := false
		for _, byteData := range bytes {
			if byteData == i {
				byteFound = true
				break
			}
		}

		if !byteFound {
			return false
		}
	}

	return true
}

func getAddressesWithValue(bytes []byte, value byte) []int {
	var res []int

	for address, byteData := range bytes {
		if byteData == value {
			res = append(res, address)
		}
	}

	return res
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func byteArrayToInt(bytes []byte) int {
	switch variableLength {
	case 8:
		return int(binary.LittleEndian.Uint64(bytes))
	case 4:
		return int(binary.LittleEndian.Uint32(bytes))
	case 2:
		return int(binary.LittleEndian.Uint16(bytes))
	}

	return -1
}

func addressToByteArray(address int) []byte {
	res := new(bytes.Buffer)

	switch variableLength {
	case 8:
		err := binary.Write(res, binary.LittleEndian, int64(address))
		check(err)
	case 4:
		err := binary.Write(res, binary.LittleEndian, int32(address))
		check(err)
	case 2:
		err := binary.Write(res, binary.LittleEndian, int16(address))
		check(err)
	}

	return res.Bytes()
}

func basic() {
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

	// Write all binaries to file
	ioutil.WriteFile("./encodedFile.bin", targetBytes, 0644)
	ioutil.WriteFile("./keyFile.bin", key, 0644)

	// =======================================================

	// Load the encoded file
	encodedFile, err := ioutil.ReadFile("./encodedFile.bin")
	check(err)

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
