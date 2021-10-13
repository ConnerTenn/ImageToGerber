package main

import (
	"fmt"
	"os"
)

func ShowHelp() {
	fmt.Print(`Usage:
./ImageToGerber.py [-h] [-c <config file>] [-m <method>] <image file>
    -h --help       Show the help menu.
                    This argument is optional and will cause the program to
                    exit immediately.
    -c --config     Specify the config file to be used for processing the image.
                    If none is specified, a default config will be used.
                    This argument is optional.
    -m --method     <Dist/Blur/Blocky>
                    Select which pixel selection and processing method should be used.
                    Defaults to Dist
    <image file>    The image file to process. This argument is reqired.

`)
	os.Exit(-1)
}

func main() {
	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" {
			ShowHelp()
		}
	}
}
