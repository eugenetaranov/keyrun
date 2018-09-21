package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, secret string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(secret)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, secret string) []byte {
	key := []byte(createHash(secret))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func encryptFile(spath string, secret string) error {
	dpath := fmt.Sprintf("%s.enc", spath)
	sdata, err := ioutil.ReadFile(spath)
	if err != nil {
		log.Fatalf("Failed to read %s\n", spath)
	}

	f, _ := os.Create(dpath)
	defer f.Close()
	f.Write(encrypt(sdata, secret))
	return nil
}

func decryptFile(spath string, secret string) error {
	dpath := strings.TrimSuffix(spath, ".enc")
	encdata, err := ioutil.ReadFile(spath)
	if err != nil {
		log.Fatalf("Failed to read %s\n", spath)
	}

	f, _ := os.Create(dpath)
	defer f.Close()
	f.Write(decrypt(encdata, secret))
	return nil
}

func decryptFileString(spath string, secret string) error {
	encdata, err := ioutil.ReadFile(spath)
	if err != nil {
		log.Fatalf("Failed to read %s\n", spath)
	}

	fmt.Print(string(decrypt(encdata, secret)))
	return nil
}
