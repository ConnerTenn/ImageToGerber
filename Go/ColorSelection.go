package main

import (
	"image"
	"image/color"
	"math"
)

type Condition struct {
	Fmt    string
	Arg    float64
	TolPos float64
	TolNeg float64
}

type Rule struct {
	Cond []Condition
	Inv  bool
}

type ColorRepr struct {
	R float64
	G float64
	B float64
	A float64

	V  float64 //Value
	C  float64 //Chroma
	L  float64 //Lightness
	H  float64 //Hue
	Sv float64 //Saturation (Value)
	Sl float64 //Saturation (Lightness)
}

func GetColourRepr(pixel color.Color) ColorRepr {
	// https://en.wikipedia.org/wiki/HSL_and_HSV

	r, g, b, a := color.RGBA64Model.Convert(pixel).RGBA() //Slice off Alpha channel

	//Convert rgb from 0->255 to 0->1
	repr := ColorRepr{
		R: float64(r) / 65535.0,
		G: float64(g) / 65535.0,
		B: float64(b) / 65535.0,
		A: float64(a) / 65535.0,

		V:  0, //Value
		C:  0, //Chroma
		L:  0, //Lightness
		H:  0, //Hue
		Sv: 0, //Saturation (Value)
		Sl: 0, //Saturation (Lightness)
	}

	var maxRGB float64 = math.Max(math.Max(repr.R, repr.G), repr.B)
	var minRGB float64 = math.Min(math.Min(repr.R, repr.G), repr.B)

	repr.V = maxRGB
	repr.C = repr.V - minRGB
	repr.L = (maxRGB + minRGB) / 2

	if repr.C == 0 {
	} else if repr.R > repr.G && repr.R > repr.B {
		repr.H = (0 + (repr.G-repr.B)/repr.C) / 6
	} else if repr.G > repr.R && repr.G > repr.B {
		repr.H = (2 + (repr.B-repr.R)/repr.C) / 6
	} else if repr.B > repr.R && repr.B > repr.G {
		repr.H = (4 + (repr.R-repr.G)/repr.C) / 6
	}

	// Sv = (0 if V==0 else C/V)
	if repr.V != 0 {
		repr.Sv = repr.C / repr.V
	}

	// Sl = 0 if L==0 or L==1 else (V-L)/min(L,1-L)
	if repr.L != 0 && repr.L != 1 {
		repr.Sl = repr.V - repr.L/math.Min(repr.L, 1-repr.L)
	}

	return repr
}

func TestTol(value float64, target float64, tolPos float64, tolNeg float64) bool {
	return value-target <= tolPos && target-value <= tolNeg
}

func TestTolWrap(value float64, target float64, tolPos float64, tolNeg float64) bool {
	return math.Mod(value-target, 1.0) <= tolPos && math.Mod(target-value, 1.0) <= tolNeg
}

var Selection *[]Rule

func SelectPixel(pixel color.Color) bool {
	repr := GetColourRepr(pixel)

	passRules := false

	for _, rule := range *Selection {
		passConditions := true

		for _, cond := range rule.Cond {
			if cond.Fmt == "R" {
				passConditions = passConditions && TestTol(repr.R, cond.Arg, cond.TolPos, cond.TolNeg)
			} else if cond.Fmt == "G" {
				passConditions = passConditions && TestTol(repr.G, cond.Arg, cond.TolPos, cond.TolNeg)
			} else if cond.Fmt == "B" {
				passConditions = passConditions && TestTol(repr.B, cond.Arg, cond.TolPos, cond.TolNeg)
			} else if cond.Fmt == "A" {
				passConditions = passConditions && TestTol(repr.A, cond.Arg, cond.TolPos, cond.TolNeg)
			} else if cond.Fmt == "V" {
				passConditions = passConditions && TestTol(repr.V, cond.Arg, cond.TolPos, cond.TolNeg)
			} else if cond.Fmt == "C" {
				passConditions = passConditions && TestTol(repr.C, cond.Arg, cond.TolPos, cond.TolNeg)
			} else if cond.Fmt == "L" {
				passConditions = passConditions && TestTol(repr.L, cond.Arg, cond.TolPos, cond.TolNeg)
			} else if cond.Fmt == "H" {
				passConditions = passConditions && TestTolWrap(repr.H, cond.Arg, cond.TolPos, cond.TolNeg)
			} else if cond.Fmt == "Sv" {
				passConditions = passConditions && TestTol(repr.Sv, cond.Arg, cond.TolPos, cond.TolNeg)
			} else if cond.Fmt == "Sl" {
				passConditions = passConditions && TestTol(repr.Sl, cond.Arg, cond.TolPos, cond.TolNeg)
			}
		}

		if rule.Inv {
			//Selections are excluded
			passRules = passRules && !passConditions
		} else {
			//Selections are added together
			passRules = passRules || passConditions
		}
	}

	return passRules
}

func SelectRow(img image.Image, yidx chan int, doneidx chan int, newimg *image.RGBA, done chan bool) {
	for true {
		//Get next job
		y, more := <-yidx

		if more {
			//Process Job
			for x := 0; x < img.Bounds().Max.X; x++ {
				if SelectPixel(img.At(x, y)) {
					newimg.Set(x, y, color.White)
				} else {
					newimg.Set(x, y, color.Black)
				}
			}

			doneidx <- y
		} else {
			//No more jobs
			done <- true
			return
		}
	}
}

func SelectColors(img image.Image, selection *[]Rule, printer Printer) *image.RGBA {
	Selection = selection

	//New image for selection
	var newimg *image.RGBA = image.NewRGBA(img.Bounds())

	yidx := make(chan int, img.Bounds().Dx()*img.Bounds().Dy())
	doneidx := make(chan int, img.Bounds().Dx()*img.Bounds().Dy())
	done := make(chan bool)

	numThreads := 10

	printer.Print(TERM_GREEN + "Selecting Colors" + TERM_RESET)

	//Spawn threads
	for i := 0; i < numThreads; i++ {
		go SelectRow(img, yidx, doneidx, newimg, done)
	}

	//Queue Jobs
	for y := 0; y < img.Bounds().Dy(); y++ {
		yidx <- y
	}
	close(yidx)

	//Progress print
	i := 0
	barDone := make(chan bool)
	barResp := make(chan bool)
	PrintProgressBar("Selecting Colors", TERM_GREEN, &i, 0, img.Bounds().Dy()-1-1, printer, barDone, barResp)
	//Counter for progress bar
	for i < img.Bounds().Dy() {
		<-doneidx
		i++
	}

	//Wait for all jobs to complete
	for i := 0; i < numThreads; i++ {
		<-done
	}
	close(doneidx)
	close(done)

	//Synchronize with stopping the progress bar
	barDone <- true
	<-barResp
	close(barDone)
	close(barResp)

	return newimg
}
