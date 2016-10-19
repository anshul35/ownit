package Utilities

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"github.com/anshul35/ownit/Settings/Constants"
)

func DecryptRSA(cipherText []byte) ([]byte, error) {

	//Decode the base64 encoding
	data, err := base64.URLEncoding.DecodeString(string(cipherText))
	if err != nil {
		return nil, err
	}

	key := make([]byte, 32)

	//Get the RSA private key from os
	pk, err := ioutil.ReadFile(Constants.RSAPrivateKeyFile)
	block, _ := pem.Decode(pk)
	if block == nil {
		return nil, errors.New("Decrypt RSA: Block parse error")
	}
	private_key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	//Decrypt using RSA OAEP with sha256 as hash fucntion, rand reader as
	//source of entropy, and label as nil
	if key, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, private_key, data, nil); err != nil {
		return nil, err
	}
	return key, nil
}

func DecryptAES(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
