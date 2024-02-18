package encryption

import (
	"crypto/aes"
	"crypto/cipher"
)

func Encrypt(message []byte, encKey []byte) (ciphertext, nonce []byte, err error) {
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce = RandomBytes(aead.NonceSize())
	cipherbytes := aead.Seal(nil, nonce, message, nil)

	return cipherbytes, nonce, nil
}

func Decrypt(ciphertext []byte, encKey []byte, nonce []byte) (message []byte, err error) {
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
