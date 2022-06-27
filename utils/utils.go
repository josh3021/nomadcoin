// Package utils contain functions to be used globally in application
package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToJSON(i interface{}) []byte {
	d, err := json.Marshal(i)
	HandleErr(err)
	return d
}

func ToBytes(i interface{}) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	HandleErr(encoder.Encode(i))
	return buffer.Bytes()
}

// FromBytes returns decoded data in interface i
func FromBytes(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(decoder.Decode(i))
}

// Hash returns sha256 hex-encoded data
func Hash(i interface{}) string {
	s := fmt.Sprintf("%v", i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func EncodeBigInts(bigA, bigB *big.Int) string {
	bytes := append(bigA.Bytes(), bigB.Bytes()...)
	return fmt.Sprintf("%x", bytes)
}

func RestoreBigInts(payload string) (*big.Int, *big.Int, error) {
	payloadBytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	bigA, bigB := big.Int{}, big.Int{}
	halfLength := len(payloadBytes) / 2
	firstHalfBytes := payloadBytes[:halfLength]
	secondHalfBytes := payloadBytes[halfLength:]
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)
	return &bigA, &bigB, nil
}

func Splitter(s, sep string, index int) string {
	r := strings.Split(s, sep)
	if len(r)-1 < index {
		return ""
	}
	return r[index]
}
