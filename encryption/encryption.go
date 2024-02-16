package encryption

import (
	"crypto/aes"
	"crypto/cipher"
)

func Encrypt(message string, encKey []byte) (ciphertext, nonce []byte, err error) {
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce = RandomBytes(aead.NonceSize())
	cipherbytes := aead.Seal(nil, nonce, []byte(message), nil)

	return cipherbytes, nonce, nil
}

func Decrypt(ciphertext string, encKey []byte, nonce []byte) (message []byte, err error) {
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aead.Open(nil, nonce, []byte(ciphertext), nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
