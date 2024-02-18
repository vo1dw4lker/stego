package lsb

import (
	"image"
	"stego/encryption"
)

// message should be encoded as [whole message length]\x00[nonce length]\x00[nonce][salt (fixed 8 bytes)][encrypted message]

func Extract(img image.Image, encrypted bool, password string) (string, error) {
	ch := make(chan byte)
	go streamBytes(ch, img)

	//var message, lengthStr, nonceStr []byte
	var message, lenBytes, nonceBytes []byte
	var length int
	var key, nonce []byte

	// Read the length of the message
	for {
		char := <-ch
		if char == 0 {
			break
		}
		lenBytes = append(lenBytes, char)
	}

	// Read the nonce size
	if encrypted {
		for {
			char := <-ch
			if char == 0 {
				break
			}
			nonceBytes = append(nonceBytes, char)
		}
	}

	// Convert the length to int
	length = byteMerge(lenBytes)

	// Read nonce, make key
	if encrypted {
		nonceLength := byteMerge(nonceBytes)

		nonce = make([]byte, nonceLength)
		for i := 0; i < nonceLength; i++ {
			nonce[i] = <-ch
		}
		salt := make([]byte, 8)
		for i := 0; i < 8; i++ {
			salt[i] = <-ch
		}
		key, _ = encryption.PBKDF2(password, salt)

		length -= nonceLength + 2 + 8
	}

	// Read the message
	for i := 0; i < length; i++ {
		char := <-ch
		message = append(message, char)
	}

	// Decrypt the message
	if encrypted {
		decMessage, err := encryption.Decrypt(message, key, nonce)
		if err != nil {
			return "", err
		}
		return string(decMessage), nil

	}

	return string(message), nil
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

func byteMerge(bytes []byte) int {
	var result int
	for _, b := range bytes {
		result += int(b)
	}
	return result
}
