package main

import (
	"image"
	"image/color"
	"math"
)

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

	r, g, b, a := pixel.RGBA() //Slice off Alpha channel

	//Convert rgb from 0->255 to 0->1
	repr := ColorRepr{
		R: float64(r) / 255.0,
		G: float64(g) / 255.0,
		B: float64(b) / 255.0,
		A: float64(a) / 255.0,

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

func SelectPixel(pixel color.Color) {
}

func SelectColors(img image.Image) {
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			SelectPixel(img.At(x, y))
		}
	}
}
