# steGo

steGo is a steganography tool that allows you to hide files in images.
It is written in Go and uses the `image` package to manipulate images.

## Installation

To build from source, you need to have Go installed. Then, clone the repo, cd to the directory and run `go build .`.
If you want to reduce the size of the binary, compile with the `-ldflags "-s -w"` flags.

## Usage

```
./stego {flags}
Available flags:
  -h:   Show help message
  -d:   Extract mode
  -e:   Embed mode
  
  -i string
        Specifies input file
  -o string
        Specifies output file (embed mode only)
  -t string
        Text to hide (embed mode only)
```

To pass multi-word text to the `-t` flag, use quotes.

## Examples

Embedding:
```
./stego -e -i image.png -o output.png -t "Hello, world!"
```

Extracting:
```
./stego -d -i output.png
```

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
