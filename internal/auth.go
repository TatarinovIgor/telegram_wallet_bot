package internal

import (
	"crypto/x509"
	"encoding/pem"
	jwt "github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

type Payload struct {
	PublicKeyId string `json:"public_key_id"`
	MerchantId  string `json:"merchant_id"`
	ExternalId  string `json:"external_id"`
}

func GenerateAuthJWT(keyPath, merchantID, pubKeyID, externalID string) string {
	// Parsing private key
	privateKey, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("could not read private key: %s, error: %v", keyPath, err)
	}
	block, _ := pem.Decode(privateKey)
	if block == nil {
		panic("failed to parse PEM block containing the private key")
	}
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("could not parse private key env variable: %s, error: %v", keyPath, err)
	}
	tokenPayload := Payload{pubKeyID, merchantID, externalID}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"payload": tokenPayload,
			"sub":     "1234567890",
			"iat":     time.Now().Unix(),
			"exp":     time.Now().Add(time.Hour).Unix(),
		})
	jwtToken, err := t.SignedString(private)
	if err != nil {
		log.Fatalf("could not sign token, error: %v", err)
	}
	return jwtToken
}
