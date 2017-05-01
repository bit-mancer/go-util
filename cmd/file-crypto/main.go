package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bit-mancer/go-util/crypto"
)

var (
	encrypt    bool
	decrypt    bool
	base64Key  string
	inputFile  string
	outputFile string
)

func init() {
	flag.BoolVar(&encrypt, "e", false, "Encrypt.")
	flag.BoolVar(&decrypt, "d", false, "Decrypt.")
	flag.StringVar(&base64Key, "k", "", "Base64-encoded AES-256 key.")
	flag.StringVar(&inputFile, "i", "", "Input file.")
	flag.StringVar(&outputFile, "o", "", "Output file.")
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-e | -d] -k <key> -i <input-file> -o <output-file>\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	switch {
	case base64Key == "", inputFile == "", outputFile == "":
		fallthrough
	case encrypt && decrypt:
		fallthrough
	case !encrypt && !decrypt:
		flag.Usage()
	}

	key, err := crypto.NewAESKeyFromBase64(base64Key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading the base64-encoded AES-256 key:", err)
		os.Exit(1)
	}

	fileInfo, err := os.Stat(inputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting information on the input file:", err)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input file:", err)
		os.Exit(1)
	}

	var renderedBytes []byte

	if encrypt {
		renderedBytes, err = crypto.Encrypt(data, key)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error encrypting:", err)
			os.Exit(1)
		}
	} else if decrypt {
		renderedBytes, err = crypto.Decrypt(data, key)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error decrypting:", err)
			os.Exit(1)
		}
	} else {
		panic("no mode specified")
	}

	if err = ioutil.WriteFile(outputFile, renderedBytes, fileInfo.Mode()); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing output file:", err)
		os.Exit(1)
	}
}
