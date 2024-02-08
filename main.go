package main

import (
	"fmt"
	"stego/cli"
)

func main() {
	args, err := cli.ParseCli()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", args)
}
