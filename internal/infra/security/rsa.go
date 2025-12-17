package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

type RSADecrypter struct {
	privateKey *rsa.PrivateKey
}

func NewRSADecrypter(privateKeyPath string) (*RSADecrypter, error) {
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try parsing as PKCS8 if PKCS1 fails
		pk8, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err2 != nil {
			return nil, fmt.Errorf("failed to parse private key: %v (also failed pkcs8: %v)", err, err2)
		}
		var ok bool
		priv, ok = pk8.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("key is not of type *rsa.PrivateKey")
		}
	}

	return &RSADecrypter{privateKey: priv}, nil
}

func (r *RSADecrypter) Decrypt(encryptedBase64 []byte) (string, error) {

	decryptedBytes, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		r.privateKey,
		encryptedBase64,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(decryptedBytes), nil
}
