package auth_aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
)

type Aes struct {
	block cipher.Block
}

func NewAes(key []byte) *Aes {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic("failed to create aes cipher: " + err.Error())
	}

	return &Aes{
		block: block,
	}
}

func (a *Aes) Encrypt(token string) string {
	ciphertext := make([]byte, aes.BlockSize+len(token))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		panic("failed to generate iv: " + err.Error())
	}

	stream := cipher.NewCFBEncrypter(a.block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(token))

	return hex.EncodeToString(ciphertext)
}

func (a *Aes) Decrypt(encrypted string) string {
	ciphertext, _ := hex.DecodeString(encrypted)
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(a.block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext)
}
