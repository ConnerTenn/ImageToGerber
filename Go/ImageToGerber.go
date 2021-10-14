package main

import (
	"fmt"
	"image/png"
	"os"
)

type Options struct {
	ConfigFileName string
	SelectMethod   string
	ImageFilename  string
}

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
	options := Options{ConfigFileName: "Default.cfg", SelectMethod: "Blocky"}
	for i := 0; i < len(os.Args); i++ {
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
			options.ImageFilename = arg
		}
	}

	file, err := os.Open(options.ImageFilename)
	defer file.Close()
	if err != nil {
		// fmt.Println("Failed to open image \"" + options.ImageFilename + "\"")
		fmt.Println(err)
		os.Exit(-1)
	}
	img, err := png.Decode(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Printf("Image Resolution: %dx%d\n", img.Bounds().Dx(), img.Bounds().Dy())

	newimg := SelectColors(img)
	ofile, _ := os.Create("Out.png")
	png.Encode(ofile, newimg)
	ofile.Close()
}
