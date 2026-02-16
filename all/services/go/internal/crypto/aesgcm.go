package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

type AESGCM struct {
	aead cipher.AEAD
}

func NewAESGCM(rawKey []byte) (*AESGCM, error) {
	if len(rawKey) != 32 {
		return nil, errors.New("aesgcm: key must be 32 bytes")
	}
	block, err := aes.NewCipher(rawKey)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &AESGCM{aead: aead}, nil
}

func (c *AESGCM) Encrypt(plaintext, aad []byte) (nonce []byte, ciphertext []byte, err error) {
	nonce = make([]byte, c.aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	ciphertext = c.aead.Seal(nil, nonce, plaintext, aad)
	return nonce, ciphertext, nil
}

func (c *AESGCM) Decrypt(nonce, ciphertext, aad []byte) ([]byte, error) {
	return c.aead.Open(nil, nonce, ciphertext, aad)
}
