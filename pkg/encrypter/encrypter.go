package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

type Encrypter struct{}

func NewEncrypter() *Encrypter {
	return &Encrypter{}
}

func generateKey(secrets []string) [32]byte {
	var key [32]byte
	for _, secret := range secrets {
		byteSecret := append([]byte(secret), key[:]...)
		key = sha256.Sum256(byteSecret)
	}
	return key
}

func (e *Encrypter) EncryptMessage(msg []byte, secrets ...string) ([]byte, error) {
	key := generateKey(secrets)

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, msg, nil)
	return append(nonce, ciphertext...), nil
}

func (e *Encrypter) DecryptMessage(msg []byte, secrets ...string) ([]byte, error) {
	key := generateKey(secrets)

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(msg) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce := msg[:nonceSize]
	ciphertext := msg[nonceSize:]

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}
