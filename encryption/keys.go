package encryption

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

const publicKeyDir = "keys/"

func GetPublicKeyForUser(userID string) (*rsa.PublicKey, error) {
	publicKeyPath := fmt.Sprintf("%s%s_pub.pem", publicKeyDir, userID)

	pubKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not read public key file for user %s: %v", userID, err)
	}

	block, _ := pem.Decode(pubKeyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key for user %s", userID)
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key for user %s: %v", userID, err)
	}

	rsaPubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key for user %s is not an RSA key", userID)
	}

	return rsaPubKey, nil
}
