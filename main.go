package main

import (
	"fmt"
	"image/png"
	"os"
	"stego/cli"
	"stego/encryption"
	"stego/lsb"
)

func main() {
	args, err := cli.ParseCli()
	if err != nil {
		panic(err)
	}

	switch args.Mode {
	case cli.ModeEmbed:
		embed(args)
	case cli.ModeExtract:
		extract(args)
	}
}

func embed(args *cli.Args) {
	outfile, err := os.Create(args.OutFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outfile.Close()

	var encKey []byte = nil
	if args.Encrypted {
		key, salt := encryption.PBKDF2(args.EncPasswd, nil)
		encKey = append(salt, key...)
	}

	newImg, err := lsb.Embed(args.Image, []byte(args.Text), args.Encrypted, encKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = png.Encode(outfile, newImg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Text embedded successfully")
}

func extract(args *cli.Args) {
	message, err := lsb.Extract(args.Image, args.Encrypted, args.EncPasswd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(message)
}
