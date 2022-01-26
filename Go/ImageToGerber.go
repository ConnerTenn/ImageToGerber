package main

import (
	"fmt"
	"image"
	"image/png"
	"math/rand"
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
	processlist, ditherlist := ParseConfig(options.ConfigFileName)
	imageCache := make(map[string]image.Image)

	var maxBounds = image.Point{0, 0}

	for _, process := range processlist {
		rand.Seed(process.Seed)

		if _, exist := imageCache[process.Infile]; !exist {
			//Open Image
			splitpath := strings.Split(process.Infile, "/")
			imageName := splitpath[len(splitpath)-1]
			Print(TERM_GREY+"Opening File \"%s\"\n"+TERM_RESET, imageName)
			file, err := os.Open(process.Infile)
			CheckError(err)
			img, err := png.Decode(file)
			CheckError(err)

			Print("Image \"%s\" Resolution: %dx%d", imageName, img.Bounds().Dx(), img.Bounds().Dy())

			if img.Bounds().Dx() > maxBounds.X {
				maxBounds.X = img.Bounds().Dx()
			}
			if img.Bounds().Dy() > maxBounds.Y {
				maxBounds.Y = img.Bounds().Dy()
			}

			imageCache[process.Infile] = img
			// imageCache[process.Infile] = ImageCache{Img: img, Dithers: make(map[float64]DitherImg)}
		}
	}

	for _, dither := range ditherlist {
		ditherImg := GenerateDither(maxBounds.X, maxBounds.Y, dither.Factor, dither.Scale)
		WriteDitherToFile(ditherImg, dither.OutFile)
	}

	done := make(chan bool)
	for _, process := range processlist {
		go func(process Process) {
			printer := NewPrinter()

			//Get Cache
			img := imageCache[process.Infile]

			for i, rule := range process.Selection {
				if len(rule.FillStr) > 0 {
					fmt.Println(rule.FillStr)
					fillfile, err := os.Open(rule.FillStr)
					CheckError(err)
					process.Selection[i].FillImg, err = png.Decode(fillfile)
					CheckError(err)
					fmt.Println("Ajdbwkjdjwadwd")

					defer fillfile.Close()
				}
			}

			//Select Config
			newimg := SelectColors(img, &process.Selection, printer)

			// //Post Process the fill type
			// if process.Fill != "Solid" {
			// 	tokens := strings.Split(process.Fill, "\"")
			// 	if len(tokens) != 3 {
			// 		CheckError("Invalid Fill specifier")
			// 	}
			// 	fillfile, err := os.Open(tokens[1])
			// 	CheckError(err)
			// 	fillimg, err := png.Decode(fillfile)
			// 	CheckError(err)

			// 	//Process image
			// 	newimg = FillMask(newimg, fillimg, tokens[2])

			// 	fillfile.Close()
			// }

			//Image preview output
			imgName := process.Outfile + "-" + process.Types[0] + ".png"
			printer.Print(TERM_GREY + "Writing Image File:" + imgName + TERM_RESET)
			ofile := CreateFile(imgName)
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

	ClosePrinter()

	//Done!
	fmt.Println("")
	fmt.Println(TERM_GREEN + "== Done ==" + TERM_RESET)
	fmt.Println("")

}
