package main

import (
	"fmt"
	"image"
	"image/color"
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

type ImageCache struct {
	Img     image.Image
	Dithers map[float64]DitherImg
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
	imageMap := make(map[string]ImageCache)

	for _, process := range processlist {
		rand.Seed(process.Seed)

		if _, exist := imageMap[process.Infile]; !exist {
			//Open Image
			splitpath := strings.Split(process.Infile, "/")
			imageName := splitpath[len(splitpath)-1]
			fmt.Printf(TERM_GREY+"Opening File \"%s\"\n"+TERM_RESET, imageName)
			file, err := os.Open(process.Infile)
			CheckError(err)
			img, err := png.Decode(file)
			CheckError(err)

			Print("Image \"%s\" Resolution: %dx%d", imageName, img.Bounds().Dx(), img.Bounds().Dy())

			imageMap[process.Infile] = ImageCache{Img: img, Dithers: make(map[float64]DitherImg)}
		}

		//Generate Dithers
		if process.HasDither {
			ditherExist := false
			//Check if a dither for this image and config has already been generated
			for dither := range imageMap[process.Infile].Dithers {
				if dither == process.Dither {
					ditherExist = true
				}
			}

			if !ditherExist {
				fmt.Println("Generate Dither")

				cache := imageMap[process.Infile]

				cache.Dithers[process.Dither] = GenerateDither(
					cache.Img.Bounds().Dx(),
					cache.Img.Bounds().Dy(),
					process.Dither)
			}
		}
	}

	done := make(chan bool)
	for _, process := range processlist {
		go func(process Process) {
			printer := NewPrinter()

			//Get Cache
			cache := imageMap[process.Infile]
			img := cache.Img

			//Select Config
			newimg := SelectColors(img, &process.Selection, printer)

			//Post Process the fill type
			if process.Fill == "Dither+" || process.Fill == "Dither-" {
				//Process the dither selection
				if process.HasDither {
					dither := cache.Dithers[process.Dither]

					//Process image
					for y := 0; y < img.Bounds().Dy(); y++ {
						for x := 0; x < img.Bounds().Dx(); x++ {
							//Blank pixel according to dither pattern
							if process.Fill == "Dither+" && !dither[y][x] ||
								process.Fill == "Dither-" && dither[y][x] {
								newimg.Set(x, y, color.RGBA{0, 0, 0, 255})
							}
						}
					}

				} else {
					CheckError("Cannot fill Dither if no dither has been specified")
				}
			}

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
