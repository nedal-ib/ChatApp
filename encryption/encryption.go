package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

func EncryptMessage(pubKey *rsa.PublicKey, message []byte) (string, error) {
	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		return "", err
	}

	ciphertext, nonce, err := encryptAES(aesKey, message)
	if err != nil {
		return "", err
	}

	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, aesKey, nil)
	if err != nil {
		return "", err
	}

	finalMessage := append(encryptedKey, nonce...)
	finalMessage = append(finalMessage, ciphertext...)
	return base64.StdEncoding.EncodeToString(finalMessage), nil
}

func DecryptMessage(privKey *rsa.PrivateKey, encryptedMessage string) ([]byte, error) {
	decodedMessage, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return nil, err
	}

	keySize := privKey.Size()
	nonceSize := 12

	if len(decodedMessage) < keySize+nonceSize {
		return nil, errors.New("invalid encrypted message format")
	}

	encryptedKey := decodedMessage[:keySize]
	nonce := decodedMessage[keySize : keySize+nonceSize]
	ciphertext := decodedMessage[keySize+nonceSize:]

	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, encryptedKey, nil)
	if err != nil {
		return nil, err
	}

	return decryptAES(aesKey, nonce, ciphertext)
}
