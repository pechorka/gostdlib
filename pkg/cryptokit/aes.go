package cryptokit

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/pechorka/stdlib/pkg/errs"
)

type AesCryptor struct {
	aesgcm cipher.AEAD
}

func NewEncrypter(key []byte) (*AesCryptor, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errs.Wrap(err, "failed to create cipher")
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errs.Wrap(err, "failed to create gcm")
	}

	return &AesCryptor{
		aesgcm: aesgcm,
	}, nil
}

func (e *AesCryptor) Encrypt(data []byte) ([]byte, error) {
	return e.encrypt(data)
}

func (e *AesCryptor) EncryptString(data string) (string, error) {
	encrypted, err := e.encrypt([]byte(data))
	if err != nil {
		return "", err
	}

	return string(encrypted), nil
}

func (e *AesCryptor) encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, e.aesgcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, errs.Wrap(err, "failed to generate nonce")
	}

	ciphertext := e.aesgcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (e *AesCryptor) DecryptString(data string) (string, error) {
	decrypted, err := e.decrypt([]byte(data))
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func (e *AesCryptor) Decrypt(data []byte) ([]byte, error) {
	return e.decrypt(data)
}

func (e *AesCryptor) decrypt(data []byte) ([]byte, error) {
	nonceSize := e.aesgcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errs.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := e.aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errs.Wrap(err, "failed to decrypt")
	}

	return plaintext, nil
}
