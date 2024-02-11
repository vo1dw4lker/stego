package main

import (
	"fmt"
	"image/png"
	"os"
	"stego/cli"
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
		panic(err)
	}
	defer outfile.Close()

	newImg, err := lsb.Embed(args.Image, args.Text)
	if err != nil {
		panic(err)
	}

	err = png.Encode(outfile, newImg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Text embedded successfully")
}

func extract(args *cli.Args) {
	message, err := lsb.Extract(args.Image)
	if err != nil {
		panic(err)
	}
	fmt.Println(message)
}
