package main

import (
	"fmt"
	"image"
)

//  +----+   +----+
//  |####| 1 |####|
//  |##+-------+##|
//  +--|-+   +-|--+
//   4 |       | 2
//  +--|-+   +-|--+
//  |##+-------+##|
//  |####| 3 |####|
//  +----+   +----+

type Point struct {
	X float64
	Y float64
}

type Segment struct {
	P1 Point
	P2 Point
}

type LineKernel struct {
	Kernel   [4]bool
	Segments [][2]int
}

var LineData = []LineKernel{
	{
		Kernel: [4]bool{true, true,
			false, false},
		Segments: [][2]int{{2, 4}},
	},
	{
		Kernel: [4]bool{false, false,
			true, true},
		Segments: [][2]int{{2, 4}},
	},
	LineKernel{
		Kernel: [4]bool{true, false,
			true, false},
		Segments: [][2]int{{1, 3}},
	},
	{
		Kernel: [4]bool{false, true,
			false, true},
		Segments: [][2]int{{1, 3}},
	},

	{
		Kernel: [4]bool{false, true,
			true, true},
		Segments: [][2]int{{1, 4}},
	},
	{
		Kernel: [4]bool{true, false,
			true, true},
		Segments: [][2]int{{1, 2}},
	},
	{
		Kernel: [4]bool{true, true,
			false, true},
		Segments: [][2]int{{3, 4}},
	},
	{
		Kernel: [4]bool{true, true,
			true, false},
		Segments: [][2]int{{2, 3}},
	},

	{
		Kernel: [4]bool{true, false,
			false, false},
		Segments: [][2]int{{1, 4}},
	},
	{
		Kernel: [4]bool{false, true,
			false, false},
		Segments: [][2]int{{1, 2}},
	},
	{
		Kernel: [4]bool{false, false,
			true, false},
		Segments: [][2]int{{3, 4}},
	},
	{
		Kernel: [4]bool{false, false,
			false, true},
		Segments: [][2]int{{2, 3}},
	},

	{
		Kernel: [4]bool{true, false,
			false, true},
		Segments: [][2]int{{1, 4}, {2, 3}},
	},
	{
		Kernel: [4]bool{false, true,
			true, false},
		Segments: [][2]int{{1, 2}, {3, 4}},
	},
}

// #Linear interpolation between p1 and p2, based on the relative values of a,b, and 0.5
func InterpolatePos(p1 float64, p2 float64, a float64, b float64) float64 {
	mu := (0.5 - a) / (b - a)
	return p1 + mu*(p2-p1)
}

// #Generate coordinates for a point on a line based on the type of line and the interpolation
func ConvertToCoords(lineIdx int, p Point, tl float64, tr float64, bl float64, br float64) Point {
	switch lineIdx {
	case 1:
		return Point{InterpolatePos(p.X+0, p.X+1, tl, tr), InterpolatePos(p.Y+1, p.Y+1, tl, tr)}
	case 2:
		return Point{InterpolatePos(p.X+1, p.X+0, tl, tr), InterpolatePos(p.Y+1, p.Y+1, tl, tr)}
	case 3:
		return Point{InterpolatePos(p.X+0, p.X+0, tl, tr), InterpolatePos(p.Y+1, p.Y+0, tl, tr)}
	case 4:
		return Point{InterpolatePos(p.X+0, p.X+0, tl, tr), InterpolatePos(p.Y+0, p.Y+1, tl, tr)}
	}
	return Point{0, 0}
}

func GenerateSegments(pos Point, tl float64, tr float64, bl float64, br float64) []Segment {
	var segments []Segment

	// Check if any work has to be done for this selection
	if !((tl >= 0.5) == (tr >= 0.5) == (bl >= 0.5) == (br >= 0.5)) {
		// Check each to see what line data matches the pixel pattern
		for _, dat := range LineData {
			// Check to see if the pixel pattern matches
			if dat.Kernel[0] == (tl >= 0.5) && dat.Kernel[1] == (tr >= 0.5) && dat.Kernel[2] == (bl >= 0.5) && dat.Kernel[3] == (br >= 0.5) {

				// Generate a line for each segment in the pixel pattern
				for _, seg := range dat.Segments {
					p1 := ConvertToCoords(seg[0], pos, tl, tr, bl, br)
					p2 := ConvertToCoords(seg[1], pos, tl, tr, bl, br)

					segments = append(segments, Segment{p1, p2})
				}
				return segments
			}
		}
	}
	return []Segment{}
}

func LineDetection(img image.Image) []Segment {
	var segments []Segment

	fmt.Println(TERM_BLUE + "== Tracing outline" + TERM_RESET)
	for y := 0; y < img.Bounds().Dy()-1; y++ {
		for x := 0; x < img.Bounds().Dx()-1; x++ {
			//Get 2x2 region of image
			tl_i, _, _, _ := img.At(x+0, y+0).RGBA()
			tr_i, _, _, _ := img.At(x+1, y+0).RGBA()
			bl_i, _, _, _ := img.At(x+0, y+1).RGBA()
			br_i, _, _, _ := img.At(x+1, y+1).RGBA()

			tl := float64(tl_i) / 65535
			tr := float64(tr_i) / 65535
			bl := float64(bl_i) / 65535
			br := float64(br_i) / 65535

			segments = append(segments, GenerateSegments(Point{float64(x), float64(y)}, tl, tr, bl, br)...)
		}
	}
	fmt.Println(TERM_GREY + "== Done Tracing outline" + TERM_RESET)

	return segments
}
