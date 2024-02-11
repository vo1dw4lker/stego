package lsb

import (
	"image"
	"strconv"
)

func Extract(img image.Image) (string, error) {
	ch := make(chan byte)
	go streamBytes(ch, img)

	var message string
	var lengthStr string
	for {
		char := <-ch
		if char == 0 {
			break
		}
		lengthStr += string(char)
	}

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		char := <-ch
		message += string(char)
	}

	return message, nil
}

// Message is encoded as: [length as str]/x00/[message]
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
