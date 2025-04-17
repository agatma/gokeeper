package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
)

type Encrypter struct{}

func NewEncrypter() *Encrypter {
	return &Encrypter{}
}

func (e *Encrypter) EncryptMessage(msg []byte, secrets ...string) ([]byte, error) {
	var key [32]byte
	for _, secret := range secrets {
		byteSecret := []byte(secret)
		byteSecret = append(byteSecret, key[:]...)
		key = sha256.Sum256(byteSecret)
	}

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	dst := aesgcm.Seal(nil, nonce, msg, nil)
	return dst, nil
}

func (e *Encrypter) DecryptMessage(msg []byte, secrets ...string) ([]byte, error) {
	var key [32]byte
	for _, secret := range secrets {
		byteSecret := []byte(secret)
		byteSecret = append(byteSecret, key[:]...)
		key = sha256.Sum256(byteSecret)
	}

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	decrypted, err := aesgcm.Open(nil, nonce, msg, nil)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
