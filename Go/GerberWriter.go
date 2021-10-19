package main

import (
	"fmt"
	"image"
	"math"
	"os"
	"strings"
	"time"
)

func FinishFile(file *os.File, buff []byte) {
	buff = append(buff, []byte("M02*")...)
	file.Write(buff)
	file.Close()
}

func NumRepr(num float64) string {
	integer := int(num)
	fraction := int(math.Mod(num, 1.0))
	return fmt.Sprintf("%04d%06d", integer, fraction)
	// return strings.ReplaceAll(fmt.Sprintf("%011.6f", num), ".", "")
}

func WriteHeader(buff []byte, gerberType string) []byte {
	buff = append(buff, []byte("%%TF.GenerationSoftware,ImageToGerber,Design,1.1.0*%\n")...)
	buff = append(buff, []byte("%%TF.SameCoordinates,Original*%\n")...)
	if gerberType == "F_Cu" {
		// Front Copper
		buff = append(buff, []byte("%TF.FileFunction,Copper,L1,Top*%\n")...)
		buff = append(buff, []byte("%TF.FilePolarity,Positive*%\n")...)
	} else if gerberType == "B_Cu" {
		// Back Copper
		buff = append(buff, []byte("%TF.FileFunction,Copper,L2,Bot*%\n")...)
		buff = append(buff, []byte("%TF.FilePolarity,Positive*%\n")...)
	} else if gerberType == "F_Mask" {
		// Front Solder Mask
		buff = append(buff, []byte("%TF.FileFunction,Soldermask,Top*%\n")...)
		buff = append(buff, []byte("%TF.FilePolarity,Negative*%\n")...)
	} else if gerberType == "B_Mask" {
		// Back Solder Mask
		buff = append(buff, []byte("%TF.FileFunction,Soldermask,Bot*%\n")...)
		buff = append(buff, []byte("%TF.FilePolarity,Negative*%\n")...)
	} else if gerberType == "F_SilkS" {
		// Front Silk screen
		buff = append(buff, []byte("%TF.FileFunction,Legend,Top*%\n")...)
		buff = append(buff, []byte("%TF.FilePolarity,Positive*%\n")...)
	} else if gerberType == "B_SilkS" {
		// Back Silk screen
		buff = append(buff, []byte("%TF.FileFunction,Legend,Bot*%\n")...)
		buff = append(buff, []byte("%TF.FilePolarity,Positive*%\n")...)
	} else if gerberType == "Edge_Cuts" {
		// Edge Cuts
		buff = append(buff, []byte("%%TF.FileFunction,Profile,NP*%\n")...)
	}

	buff = append(buff, []byte("%FSLAX46Y46*%\n")...) //Number Format
	buff = append(buff, []byte("%MOMM*%\n")...)       //Millimeters
	buff = append(buff, []byte("%LPD*%\n")...)        //Dark polarity
	buff = append(buff, []byte("G04 == End Header ==*\n")...)

	return buff
}

func GerberCreateLine(buff []byte, scaleWidth float64, scaleHeight float64, img *image.Image, x int, y int, rectsize int) []byte {

	buff = append(buff, []byte(fmt.Sprintf("D%d*\n", rectsize+10))...) //Use Rectangle Object
	buff = append(buff, []byte(
		fmt.Sprintf("X%sY%sD03*\n",
			NumRepr(scaleWidth*(float64(x)-float64(rectsize)/2)),
			NumRepr(scaleHeight*(float64((*img).Bounds().Dy()-y)+1.0/2)),
		),
	)...) //Place rectangle

	return buff
}

func GenerateGerberFillLines(img image.Image, gerberWidth float64, gerberHeight float64, filename string, gerberType string, printer Printer) {
	// file = CreateFile(filename, gerberType)
	file := CreateFile(filename + "-" + gerberType + ".gbr")

	scaleWidth := gerberWidth / float64(img.Bounds().Dx())
	scaleHeight := gerberHeight / float64(img.Bounds().Dy())

	printer.Print("Writing Gerber")
	buff := make([]byte, 0)
	buff = WriteHeader(buff, gerberType)

	for i := 1; i <= img.Bounds().Dx(); i++ {
		buff = append(buff, []byte(fmt.Sprintf("%%ADD%dR,%.6fX%.6f*%%\n", i+10, scaleWidth*float64(i), scaleHeight))...) //Rectangle Object
	}

	tStart := time.Now()
	tLast := time.Time{}
	rectsize := 0
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			//Get pixel value at this point
			r, _, _, _ := img.At(x, y).RGBA()
			if r == 65535 {
				rectsize += 1
				// Handle line
			} else if rectsize > 0 {
				buff = GerberCreateLine(buff, scaleWidth, scaleHeight, &img, x, y, rectsize)
				// DrawLine(x, y, rectsize)
				rectsize = 0
			}
		}
		// Handle last run
		if rectsize > 0 {
			buff = GerberCreateLine(buff, scaleWidth, scaleHeight, &img, img.Bounds().Dx()-1, y, rectsize)
			// DrawLine(img.Bounds().Dx()-1, y, rectsize)
			rectsize = 0
		}
		tNow := time.Now()

		if tNow.Sub(tLast) > 100*time.Millisecond {
			progress := 10 * y / img.Bounds().Dy()
			bar := strings.Repeat("=", progress) + strings.Repeat(" ", 10-progress)
			printer.Print(fmt.Sprintf("Writing Gerber %s [%s] Time:%v", , bar, tNow.Sub(tStart)))
			tLast = tNow
		}
	}
	// fmt.Println("Time:", t2.Sub(t1))

	FinishFile(file, buff)
}

func GenerateGerberTrace(img image.Image, gerberWidth float64, gerberHeight float64, filename string, gerberType string) {
	segments := LineDetection(img)
	_ = segments

	scaleWidth := gerberWidth / float64(img.Bounds().Dx())
	scaleHeight := gerberHeight / float64(img.Bounds().Dy())

	file := CreateFile(filename + "-" + gerberType + ".gbr")

	fmt.Println(TERM_BLUE + "== Writing Gerber ==" + TERM_RESET)
	buff := make([]byte, 0)
	buff = WriteHeader(buff, gerberType)

	buff = append(buff, []byte("G01*\n")...) //Linear interpolation mode
	// # file.write("%TA.AperFunction,Profile*%")
	buff = append(buff, []byte("%ADD10C,0.200000*%")...) //Circle aperture
	// # file.write("%TD*%")
	buff = append(buff, []byte("D10*")...) //Use aperture

	for _, segment := range segments {
		buff = append(buff, []byte(
			fmt.Sprintf("X%sY%sD02*\n",
				NumRepr(scaleWidth*segment.P1.X),
				NumRepr(scaleHeight*(float64(img.Bounds().Dy())-segment.P1.Y)),
			),
		)...)
		buff = append(buff, []byte(
			fmt.Sprintf("X%sY%sD02*\n",
				NumRepr(scaleWidth*segment.P2.X),
				NumRepr(scaleHeight*(float64(img.Bounds().Dy())-segment.P2.Y)),
			),
		)...)
	}

	FinishFile(file, buff)
	fmt.Println(TERM_GREY + "== Done Writing Gerber ==" + TERM_RESET)
}
