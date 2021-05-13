/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// cryptoChaincode is allows the following transactions
//    "put", "key", val - returns "OK" on success
//    "get", "key" - returns val stored previously
type cryptoChaincode struct {
	contractapi.Contract
}

const (
	// AESKeyLength is the default AES key length
	AESKeyLength = 32
	// NonceSize is the default NonceSize
	NonceSize = 24
)


//
// genAESKey returns a random AES key of length AESKeyLength
// 3 Functions to support Encryption and Decryption
// GENAESKey() - Generates AES symmetric key
func (t *cryptoChaincode) genAESKey() ([]byte, error) {
	key := make([]byte, AESKeyLength)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (t *cryptoChaincode) encrypt(key []byte, byteArray []byte) []byte {

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Empty array of 16 + byteArray length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(byteArray))

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// encrypt bytes from byteArray to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], byteArray)

	return ciphertext
}

func (t *cryptoChaincode) decrypt(key []byte, ciphertext []byte) []byte {

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		panic("Text is too short")
	}

	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]

	// Remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// decrypt bytes from ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext
}

func (t *cryptoChaincode) encryptAndDecrypt(arg string) []byte {
	AESKey, _ := t.genAESKey()
	AESEnc := t.encrypt(AESKey, []byte(arg))

	value := t.decrypt(AESKey, AESEnc)
	return value
}

func (t *cryptoChaincode) Put(ctx contractapi.TransactionContextInterface, k, v string) error {
	cryptoArg := t.encryptAndDecrypt(v)
	return ctx.GetStub().PutState(k, cryptoArg)
}

func (t *cryptoChaincode) Get(ctx contractapi.TransactionContextInterface, k string) (string,error) {
	val, err := ctx.GetStub().GetState(k)
	if err != nil {
		return "",err
	}
	return string(val),err
}


func main() {
	cc, err := contractapi.NewChaincode(new(cryptoChaincode))
	if err != nil {
		panic(err.Error())
	}
	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting CryptoChaincode chaincode: %s", err)
	}
}
