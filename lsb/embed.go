package lsb

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"stego/encryption"
)

// message should be encoded as [whole message length]\x00[nonce length]\x00[nonce][salt (fixed 8 bytes)][encrypted message]

func Embed(origImg image.Image, message []byte, encrypted bool, encKey []byte) (image.Image, error) {
	if encrypted {
		key := encKey[8:] // cut 8 bytes of salt
		encMessage, nonce, err := encryption.Encrypt(message, key)
		if err != nil {
			return nil, err
		}
		message = append(byteSplit(len(nonce)), 0x00) // nonce length + \x00
		message = append(message, nonce...)           // nonce
		message = append(message, encKey[:8]...)      // salt
		message = append(message, encMessage...)      // encrypted message
	}

	message = addLen(message)
	newImg := image.NewRGBA(origImg.Bounds())
	imgWidth := origImg.Bounds().Max.X
	imgHeight := origImg.Bounds().Max.Y

	ch := make(chan bool)

	go streamBits(ch, message)

	// Check if the message is too long
	if len(message)*8 > imgWidth*imgHeight*3 {
		return nil, fmt.Errorf("message is too long to embed in the image")
	}

	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			origColor := origImg.At(x, y)
			r, g, b, a := origColor.RGBA()
			newColor := color.RGBA{
				R: writeLSB(uint8(r>>8), <-ch),
				G: writeLSB(uint8(g>>8), <-ch),
				B: writeLSB(uint8(b>>8), <-ch),
				A: uint8(a >> 8),
			}
			newImg.Set(x, y, newColor)
		}
	}

	return newImg, nil
}

func addLen(message []byte) []byte {
	length := append(byteSplit(len(message)), 0x00)
	return append(length, message...)
}

// streamBits sends the bits of the message to the channel
// in the form of true (1) or false (0)
// and then sends random bits
func streamBits(ch chan bool, message []byte) {
	for _, char := range message {
		for i := 0; i < 8; i++ {
			ch <- char&(1<<uint(7-i)) != 0
		}
	}
	// todo: maybe implement return after imgWidth*imgHeight*3 bits
	for {
		ch <- rand.Intn(2) == 1
	}
}

func writeLSB(byte uint8, bit bool) uint8 {
	if bit {
		return byte | 1 // ex: 1010 | 0001 = 1011
	}
	return byte &^ 1 // (AND NOT) ex: 1011 &^ 0001 = 1010
}

func byteSplit(num int) []byte {
	result := make([]byte, 0)
	for {
		if num > 0xff {
			result = append(result, 0xff)
			num -= 0xff
		} else {
			result = append(result, byte(num))
			break
		}
	}
	return result
}
