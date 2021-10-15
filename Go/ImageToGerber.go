package main

import (
	"fmt"
	"image/png"
	"os"
	"strings"
)

type CmdOptions struct {
	ConfigFileName string
	SelectMethod   string
}

func ShowHelp() {
	fmt.Print(`Usage:
./ImageToGerber.py [-h] [-c <config file>] [-m <method>]
    -h --help       Show the help menu.
                    This argument is optional and will cause the program to
                    exit immediately.
    -c --config     Specify the config file to be used for processing the image.
                    If none is specified, a default config will be used.
                    This argument is optional.
    -m --method     <Dist/Blur/Blocky>
                    Select which pixel selection and processing method should be used.
                    Defaults to Dist

`)
	os.Exit(-1)
}

func main() {
	options := CmdOptions{ConfigFileName: "Default.cfg", SelectMethod: "Blocky"}

	//Options parsing
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-h" || arg == "--help" {
			ShowHelp()
		} else if arg == "-c" || arg == "--config" {
			options.ConfigFileName = os.Args[i+1]
			i++
		} else if arg == "-m" || arg == "--method" {
			options.SelectMethod = os.Args[i+1]
			i++
		} else {
			CheckError("Unexpected Argument " + arg)
		}
	}

	//Parse Config
	processlist := ParseConfig(options.ConfigFileName)

	for _, process := range processlist {
		//Open Image
		file, err := os.Open(process.Infile)
		CheckError(err)
		img, err := png.Decode(file)
		CheckError(err)

		fmt.Printf("Image Resolution: %dx%d\n", img.Bounds().Dx(), img.Bounds().Dy())

		//Select Config
		newimg := SelectColors(img, &process.Selection)

		//Debug output
		directory := process.Outfile[:strings.LastIndex(process.Outfile, "/")]
		if len(directory) > 0 {
			os.MkdirAll(directory, 0755)
		}
		ofile, err := os.Create(process.Outfile + ".png")
		CheckError(err)
		png.Encode(ofile, newimg)
		ofile.Close()
	}
}
