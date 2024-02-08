package cli

import (
	"flag"
	"fmt"
	"os"
)

type OperationMode uint8

const (
	ModeEncode OperationMode = iota
	ModeDecode
)

type Args struct {
	Mode      OperationMode
	InputFile []byte
	Text      string
}

func newArgs(mode OperationMode, infile []byte, text string) *Args {
	return &Args{
		Mode:      mode,
		InputFile: infile,
		Text:      text,
	}
}

// todo: rename encode to embed

func ParseCli() (*Args, error) {
	eMode := flag.Bool("e", false, "Encode mode")
	dMode := flag.Bool("d", false, "Decode mode")
	infilePath := flag.String("i", "", "Specifies input file")
	text := flag.String("t", "", "Text to hide")
	flag.Parse()

	if !(*eMode || *dMode) || (*eMode && *dMode) {
		return nil, fmt.Errorf("choose either '-e' (encode) or '-d' (decode)")
	}

	mode := ModeEncode
	if *dMode {
		mode = ModeDecode
	}

	if *infilePath == "" {
		return nil, fmt.Errorf("input file path is required")
	}

	if (mode == ModeEncode) && (*text == "") {
		return nil, fmt.Errorf("text is required for encoding")
	}

	fileContent, err := os.ReadFile(*infilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading input file '%v': %w", *infilePath, err)
	}

	return newArgs(mode, fileContent, *text), nil
}
