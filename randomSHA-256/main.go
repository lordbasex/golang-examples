package main

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"strconv"
	"time"
)

func SetInt64ToString(i int64) string {
	r := strconv.FormatInt(i, 16)
	return r
}

func main() {

	timeNow := time.Now()
	log.Printf("time: %v", timeNow.Unix())
	timestamp := SetInt64ToString(timeNow.Unix())

	sha256Input := sha256.Sum256([]byte(timestamp))
	validateInput := base64.StdEncoding.EncodeToString(sha256Input[:])
	length := len(validateInput)
	log.Printf("Random SHA-256: %s characters: %d", validateInput, length)

}
