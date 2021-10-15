package main

import (
	"fmt"
	"image"
	"os"
	"strings"
)

func FinishFile(file *os.File) {
	file.Write([]byte("M02*"))
	file.Close()
}

func NumRepr(num float64) string {
	// integer = int(num)
	// fraction = num
	// return F"{num:011.6f}".replace(".","")
	return strings.ReplaceAll(fmt.Sprintf("%011.6f", num), ".", "")
}

func WriteHeader(file *os.File, gerberType string) {
	file.Write([]byte("%%TF.GenerationSoftware,ImageToGerber,Design,1.1.0*%\n"))
	file.Write([]byte("%%TF.SameCoordinates,Original*%\n"))
	if gerberType == "F_Cu" {
		// Front Copper
		file.Write([]byte("%TF.FileFunction,Copper,L1,Top*%\n"))
		file.Write([]byte("%TF.FilePolarity,Positive*%\n"))
	} else if gerberType == "B_Cu" {
		// Back Copper
		file.Write([]byte("%TF.FileFunction,Copper,L2,Bot*%\n"))
		file.Write([]byte("%TF.FilePolarity,Positive*%\n"))
	} else if gerberType == "F_Mask" {
		// Front Solder Mask
		file.Write([]byte("%TF.FileFunction,Soldermask,Top*%\n"))
		file.Write([]byte("%TF.FilePolarity,Negative*%\n"))
	} else if gerberType == "B_Mask" {
		// Back Solder Mask
		file.Write([]byte("%TF.FileFunction,Soldermask,Bot*%\n"))
		file.Write([]byte("%TF.FilePolarity,Negative*%\n"))
	} else if gerberType == "F_SilkS" {
		// Front Silk screen
		file.Write([]byte("%TF.FileFunction,Legend,Top*%\n"))
		file.Write([]byte("%TF.FilePolarity,Positive*%\n"))
	} else if gerberType == "B_SilkS" {
		// Back Silk screen
		file.Write([]byte("%TF.FileFunction,Legend,Bot*%\n"))
		file.Write([]byte("%TF.FilePolarity,Positive*%\n"))
	} else if gerberType == "Edge_Cuts" {
		// Edge Cuts
		file.Write([]byte("%%TF.FileFunction,Profile,NP*%\n"))
	}

	file.Write([]byte("%FSLAX46Y46*%\n")) //Number Format
	file.Write([]byte("%MOMM*%\n"))       //Millimeters
	file.Write([]byte("%LPD*%\n"))        //Dark polarity
	file.Write([]byte("G04 == End Header ==*\n"))
}

func GenerateGerberFillLines(img image.Image, gerberWidth float64, gerberHeight float64, filename string, gerberType string) {
	// file = CreateFile(filename, gerberType)
	file := CreateFile(filename + "-" + gerberType + ".gbr")

	scaleWidth := gerberWidth / float64(img.Bounds().Dx())
	scaleHeight := gerberHeight / float64(img.Bounds().Dy())

	fmt.Println(TERM_BLUE + "== Writing Gerber ==" + TERM_RESET)
	WriteHeader(file, gerberType)

	for i := 1; i <= img.Bounds().Dx(); i++ {
		file.Write([]byte(fmt.Sprintf("%%ADD%dR,%.6fX%.6f*%%\n", i+10, scaleWidth*float64(i), scaleHeight))) //Rectangle Object
	}

	DrawLine := func(x int, y int, rectsize int) {
		file.Write([]byte(fmt.Sprintf("D%d*\n", rectsize+10))) //Use Rectangle Object
		file.Write([]byte(
			fmt.Sprintf("X%sY%sD03*\n",
				NumRepr(scaleWidth*(float64(x)-float64(rectsize)/2)),
				NumRepr(scaleHeight*(float64(img.Bounds().Dy()-y)+1.0/2)),
			),
		)) //Place rectangle
	}

	rectsize := 0
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			//Get pixel value at this point
			r, _, _, _ := img.At(x, y).RGBA()
			if r == 65535 {
				rectsize += 1
				// Handle line
			} else if rectsize > 0 {
				DrawLine(x, y, rectsize)
				rectsize = 0
			}
		}
		// Handle last run
		if rectsize > 0 {
			DrawLine(img.Bounds().Dx()-1, y, rectsize)
			rectsize = 0
		}
	}

	FinishFile(file)
	fmt.Println(TERM_GREY + "== Done Writing Gerber ==" + TERM_RESET)
}

func GenerateGerberTrace(img image.Image, gerberWidth float64, gerberHeight float64, filename string, gerberType string) {
	segments := LineDetection(img)
	_ = segments

	scaleWidth := gerberWidth / float64(img.Bounds().Dx())
	scaleHeight := gerberHeight / float64(img.Bounds().Dy())

	file := CreateFile(filename + "-" + gerberType + ".gbr")

	fmt.Println(TERM_BLUE + "== Writing Gerber ==" + TERM_RESET)
	WriteHeader(file, gerberType)

	file.Write([]byte("G01*\n")) //Linear interpolation mode
	// # file.write("%TA.AperFunction,Profile*%")
	file.Write([]byte("%ADD10C,0.200000*%")) //Circle aperture
	// # file.write("%TD*%")
	file.Write([]byte("D10*")) //Use aperture

	for _, segment := range segments {
		file.Write([]byte(
			fmt.Sprintf("X%sY%sD02*\n",
				NumRepr(scaleWidth*segment.P1.X),
				NumRepr(scaleHeight*(float64(img.Bounds().Dy())-segment.P1.Y)),
			),
		))
		file.Write([]byte(
			fmt.Sprintf("X%sY%sD02*\n",
				NumRepr(scaleWidth*segment.P2.X),
				NumRepr(scaleHeight*(float64(img.Bounds().Dy())-segment.P2.Y)),
			),
		))
	}

	FinishFile(file)
	fmt.Println(TERM_GREY + "== Done Writing Gerber ==" + TERM_RESET)
}
