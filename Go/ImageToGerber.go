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

	InitPrinter()

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

	done := make(chan bool)
	for _, process := range processlist {
		go func(process Process) {
			printer := NewPrinter()

			//Open Image
			splitpath := strings.Split(process.Infile, "/")
			printer.Print("Opening File \"%s\"", splitpath[len(splitpath)-1])
			file, err := os.Open(process.Infile)
			CheckError(err)
			img, err := png.Decode(file)
			CheckError(err)

			Print("Image Resolution: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())

			//Select Config
			newimg := SelectColors(img, &process.Selection, printer)
			_ = newimg

			//Image preview output
			printer.Print("Writing Image File")
			ofile := CreateFile(process.Outfile + ".png")
			CheckError(err)
			png.Encode(ofile, newimg)
			ofile.Close()

			//Generate a gerber for every gerber type requested
			for _, gerbertype := range process.Types {
				boardWidth := process.BoardWidth
				boardHeight := boardWidth / float64(img.Bounds().Dx()) * float64(img.Bounds().Dy())
				if gerbertype == "Edge_Cuts" {
					GenerateGerberTrace(newimg, boardWidth, boardHeight, process.Outfile, gerbertype, printer)
				} else {
					GenerateGerberFillLines(newimg, boardWidth, boardHeight, process.Outfile, gerbertype, printer)
				}
			}

			//Signal done process
			printer.Close()
			done <- true
		}(process)
	}

	//Wait for all processes to finish
	for i := 0; i < len(processlist); i++ {
		<-done
	}
	close(done)

	//Done!
	Print("")
	Print("== Done ==")
	Print("")

	ClosePrinter()
}
