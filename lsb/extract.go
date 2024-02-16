package lsb

import (
	"image"
	"stego/encryption"
	"strconv"
)

// message should be encoded as [whole message length]\x00[nonce length]\x00[nonce][salt (fixed 8 bytes)][encrypted message]

func Extract(img image.Image, encrypted bool, password string) (string, error) {
	ch := make(chan byte)
	go streamBytes(ch, img)

	var message, lengthStr, nonceStr string
	var length int
	var key, nonce []byte
	for {
		char := <-ch
		if char == 0 {
			break
		}
		lengthStr += string(char)
	}
	if encrypted {
		for {
			char := <-ch
			if char == 0 {
				break
			}
			nonceStr += string(char)
		}
	}

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", err
	}

	if encrypted {
		nonceLength, err := strconv.Atoi(nonceStr)
		if err != nil {
			return "", err
		}

		nonce = make([]byte, nonceLength)
		for i := 0; i < nonceLength; i++ {
			nonce[i] = <-ch
		}
		salt := make([]byte, 8)
		for i := 0; i < 8; i++ {
			salt[i] = <-ch
		}
		key, _ = encryption.PBKDF2(password, salt)

		length -= nonceLength + 1 + 8
	}

	for i := 0; i < length; i++ {
		char := <-ch
		message += string(char)
	}

	if encrypted {
		decMessage, err := encryption.Decrypt(message, key, nonce)
		if err != nil {
			return "", err
		}
		return string(decMessage), nil

	}

	return message, nil
}

// Message is encoded as: [length as str]/x00/[message]
// If encrypted, the message is: [length as str]/x00/[nonce size]/x00[nonce][encrypted message]
func streamBytes(ch chan byte, img image.Image) {
	imgWidth := img.Bounds().Max.X
	imgHeight := img.Bounds().Max.Y

	var char byte = 0
	var colors [3]uint32
	ctr := 0
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			origColor := img.At(x, y)
			colors[0], colors[1], colors[2], _ = origColor.RGBA()

			for i := range colors {
				char = char<<1 | readLSB(colors[i]>>8)
				ctr++
				if ctr == 8 {
					ch <- char
					char = 0
					ctr = 0
				}
			}

		}
	}
	close(ch)
}

func readLSB(byte uint32) byte {
	bit := byte & 1 // ex: 1010 & 0001 = 0000, 1011 & 0001 = 0001
	if bit == 1 {
		return 1
	}
	return 0
}
