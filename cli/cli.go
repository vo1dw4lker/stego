package cli

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
)

type OperationMode uint8

const (
	ModeEmbed OperationMode = iota
	ModeExtract
)

type Args struct {
	Mode    OperationMode
	Text    string
	Image   image.Image
	OutFile string
}

func newArgs(mode OperationMode, img image.Image, text, outfile string) *Args {
	return &Args{
		Mode:    mode,
		Text:    text,
		Image:   img,
		OutFile: outfile,
	}
}

func ParseCli() (*Args, error) {
	eMode := flag.Bool("e", false, "Embed mode")
	dMode := flag.Bool("d", false, "Extract mode")
	infilePath := flag.String("i", "", "Specifies input file")
	outfilePath := flag.String("o", "", "Specifies output file")
	text := flag.String("t", "", "Text to hide")
	flag.Parse()

	// Set the mode
	mode := ModeEmbed
	if *dMode {
		mode = ModeExtract
	}

	err := checkFlags(eMode, dMode, infilePath, outfilePath, mode, text)
	if err != nil {
		return nil, err
	}

	// Open the image file
	imgFile, err := os.ReadFile(*infilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	// Check png signature (first 8 bytes)
	err = checkPNG(imgFile)
	if err != nil {
		return nil, err
	}

	img, err := decodeImage(imgFile)
	if err != nil {
		return nil, err
	}

	return newArgs(mode, img, *text, *outfilePath), nil
}

// checkFlags checks the flags for required values
func checkFlags(eMode *bool, dMode *bool, infilePath *string, outfilePath *string, mode OperationMode, text *string) error {
	// Check for required flags
	if !(*eMode || *dMode) || (*eMode && *dMode) {
		return fmt.Errorf("choose either '-e' (encode) or '-d' (decode)")
	}

	// Check for input file
	if *infilePath == "" {
		return fmt.Errorf("input file path is required")
	}

	// Check for output file
	if mode == ModeEmbed && *outfilePath == "" {
		return fmt.Errorf("output file path is required for encoding")

	}

	// Check for required text
	if (mode == ModeEmbed) && (*text == "") {
		return fmt.Errorf("text is required for encoding")
	}
	return nil
}

// checkPNG checks the first 8 bytes of the file to see if it's a PNG
func checkPNG(imgFile []byte) error {
	if len(imgFile) < 8 {
		return fmt.Errorf("file is too small to be a PNG")
	}
	if string(imgFile[:8]) != "\x89PNG\r\n\x1a\n" {
		return fmt.Errorf("file is not a PNG")
	}
	return nil
}

// decodeImage decodes the image from the file
func decodeImage(imgFile []byte) (image.Image, error) {
	img, err := png.Decode(bytes.NewReader(imgFile))
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %v", err)
	}
	return img, nil
}
