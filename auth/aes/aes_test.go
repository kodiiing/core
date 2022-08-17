package auth_aes_test

import (
	"crypto/rand"
	auth_aes "kodiiing/auth/aes"
	"log"
	"os"
	"testing"
)

var aes *auth_aes.Aes

func TestMain(m *testing.M) {
	randKey := make([]byte, 32)
	_, err := rand.Read(randKey)
	if err != nil {
		log.Fatalf("failed to generate key: %v", err)
	}

	aes = auth_aes.NewAes(randKey)

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestEncrypt(t *testing.T) {
	token := "test"
	encrypted := aes.Encrypt(token)
	if encrypted == "" {
		t.Error("encrypted token is empty")
	}

	t.Logf("encrypted: %s", encrypted)
}

func TestDecrypt(t *testing.T) {
	token := "test"
	encrypted := aes.Encrypt(token)
	if encrypted == "" {
		t.Error("encrypted token is empty")
	}

	decrypted := aes.Decrypt(encrypted)
	if decrypted != token {
		t.Error("decrypted token is not equal to original token")
	}
}
