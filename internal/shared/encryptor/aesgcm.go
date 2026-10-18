package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

type AESGCM struct {
	key []byte
}

func NewAESGCM(key []byte) *AESGCM {
	return &AESGCM{
		key: key,
	}
}

func (a *AESGCM) Encrypt(value string) (string, error) {

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	encrypted := gcm.Seal(
		nonce,
		nonce,
		[]byte(value),
		nil,
	)

	return base64.StdEncoding.EncodeToString(
		encrypted,
	), nil
}

func (a *AESGCM) Decrypt(value string) (string, error) {

	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()

	if len(data) < nonceSize {
		return "", io.ErrUnexpectedEOF
	}

	nonce := data[:nonceSize]

	cipherText := data[nonceSize:]

	plain, err := gcm.Open(
		nil,
		nonce,
		cipherText,
		nil,
	)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
