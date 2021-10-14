package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
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

func SelectPixel(pixel color.Color) bool {
	repr := GetColourRepr(pixel)
	_ = repr

	var selection []Rule = []Rule{
		{
			Cond: []Condition{
				{Fmt: "R", Arg: 0.5, TolPos: 1.0, TolNeg: 0.01},
			},
			Inv: false,
		},
	}

	passRules := false

	for _, rule := range selection {
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

func SelectColors(img image.Image) {
	var newimg *image.RGBA = image.NewRGBA(img.Bounds())

	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			if SelectPixel(img.At(x, y)) {
				newimg.Set(x, y, color.White)
			} else {
				newimg.Set(x, y, color.Black)
			}
		}
	}

	ofile, _ := os.Create("Out.png")
	png.Encode(ofile, newimg)
}
