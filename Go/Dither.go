package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
)

//x: E [0,1]
//y: E [0,1]
func InsideRegion(tl bool, tr bool, bl bool, br bool, x float64, y float64) bool {
	x = x - 0.5
	y = y - 0.5

	// # #
	// - -
	if tl && tr && !bl && !br {
		return y <= 0
	}
	// - -
	// # #
	if !tl && !tr && bl && br {
		return y >= 0
	}
	// # -
	// # -
	if tl && !tr && bl && !br {
		return x <= 0
	}
	// - #
	// - #
	if !tl && tr && !bl && br {
		return x >= 0
	}

	// # -
	// - -
	if tl && !tr && !bl && !br {
		return -x-y >= 0.5
	}
	// - #
	// - -
	if !tl && tr && !bl && !br {
		return +x-y >= 0.5
	}
	// - -
	// # -
	if !tl && !tr && bl && !br {
		return -x+y >= 0.5
	}
	// - -
	// - #
	if !tl && !tr && !bl && br {
		return +x+y >= 0.5
	}

	// - #
	// # #
	if !tl && tr && bl && br {
		return +x+y >= -0.5
	}
	// # -
	// # #
	if tl && !tr && bl && br {
		return -x+y >= -0.5
	}
	// # #
	// - #
	if tl && tr && !bl && br {
		return +x-y >= -0.5
	}
	// # #
	// # -
	if tl && tr && bl && !br {
		return -x-y >= -0.5
	}

	// # -
	// - #
	if tl && !tr && !bl && br {
		return -x-y >= 0.5 || +x+y >= 0.5
	}
	// - #
	// # -
	if !tl && tr && bl && !br {
		return +x-y >= 0.5 || -x+y >= 0.5
	}

	// # #
	// # #
	if tl && tr && bl && br {
		return true
	}

	return false
}

type DitherImg [][]float64

func CreateKernel() ([][]float64, int, int) {
	// kernel := [][]float64{
	// 	{0.0, 0.0, 7.0},
	// 	{3.0, 5.0, 1.0},
	// }
	// kernel := [][]float64{
	// 	{0.0, 0.0, 0.0, 7.0, 5.0},
	// 	{3.0, 5.0, 7.0, 5.0, 3.0},
	// 	{1.0, 3.0, 5.0, 3.0, 1.0},
	// }
	kernel := [][]float64{
		{0.0, 0.0, 0.0, 0.0, 9.0, 7.0, 5.0},
		{3.0, 5.0, 7.0, 9.0, 7.0, 5.0, 3.0},
		{1.0, 3.0, 5.0, 7.0, 5.0, 3.0, 1.0},
		{0.0, 1.0, 3.0, 5.0, 3.0, 1.0, 0.0},
	}
	// kernel := [][]float64{
	// 	{0.0, 0.0, 0.0, 0.0, 13.0, 11.0, 7.0},
	// 	{5.0, 7.0, 11.0, 13.0, 11.0, 7.0, 5.0},
	// 	{3.0, 5.0, 7.0, 11.0, 7.0, 5.0, 3.0},
	// 	{1.0, 3.0, 5.0, 7.0, 5.0, 3.0, 1.0},
	// }

	//All elements of the kernel must add up to be 1

	sum := 0.0
	//Sum all elements of the kernel
	for y, row := range kernel {
		for x := range row {
			sum += kernel[y][x]
		}
	}
	//Divide each element by the total sum
	for y, row := range kernel {
		for x, val := range row {
			kernel[y][x] = val / sum
		}
	}

	return kernel, len(kernel[0]), len(kernel) + 1
}

func GenerateDither(width int, height int, factor float64, scale float64) DitherImg {
	Print("Generating Dither [%dx%d] factor:%f scale:%f\n", width, height, factor, scale)

	//https://en.wikipedia.org/wiki/Floyd%E2%80%93Steinberg_dithering
	//Floyd-Steinberg Dithering
	kernel, kDeadZoneWidth, kDeadZoneHeight := CreateKernel()

	ditherWidth := int(float64(width)/scale) + kDeadZoneWidth*2
	ditherHeight := int(float64(width)/scale) + kDeadZoneHeight

	rawDither := make(DitherImg, ditherHeight)
	for y := range rawDither {
		rawDither[y] = make([]float64, ditherWidth)
		for x := range rawDither[y] {
			//Add a bit of random adjustment as a 'seed' for the dither
			//Without it, the pattern would be too uniform
			rawDither[y][x] = (0.1 * (rand.Float64() - 0.5)) + factor
		}
	}

	//Perform a convolution-like operation over the image
	for y := 0; y < ditherHeight; y++ {
		for x := 0; x < ditherWidth; x++ {
			old := rawDither[y][x]

			//Clamp to 0 or 1
			new := 0.0
			if old >= 0.5 {
				new = 1.0
			}

			//Set the pixel
			rawDither[y][x] = new

			//calculate the error
			quantError := old - float64(new)

			//Convolve with the kernel
			for yoff, kernelRow := range kernel {
				if y+yoff < ditherHeight { //Bounds check

					for xoff, kfactor := range kernelRow {
						xoff = xoff - len(kernelRow)/2           //Kernel is centered horizontally about the current pixel
						if x+xoff >= 0 && x+xoff < ditherWidth { //Bounds check

							//Carry over error to neghboring pixels, weighted by the kernel factor
							rawDither[y+yoff][x+xoff] = rawDither[y+yoff][x+xoff] + quantError*kfactor
						}
					}
				}
			}
		}
	}

	//Chop off kDeadZoneWidth pixels from the left and right sides of each line
	for y := 0; y < ditherHeight; y++ {
		rawDither[y] = rawDither[y][kDeadZoneWidth : ditherWidth-kDeadZoneWidth]
	}
	//Chop off top kDeadZoneWidth pixels from image
	rawDither = rawDither[kDeadZoneHeight:]
	//Recalculate Dimensions
	ditherWidth = len(rawDither[0])
	ditherHeight = len(rawDither)

	//Make new buffer for the final dither output
	dither := make(DitherImg, height)
	for y := range dither {
		dither[y] = make([]float64, width)
		for x := range dither[y] {
			//Init to zero
			dither[y][x] = 0
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			bufX := float64(ditherWidth) * float64(x) / float64(width)
			bufY := float64(ditherHeight) * float64(y) / float64(height)

			bufX_l := int(bufX)
			bufX_r := int(math.Ceil(bufX))
			bufY_t := int(bufY)
			bufY_b := int(math.Ceil(bufY))

			var tl float64 = 0.0
			var tr float64 = 0.0
			var bl float64 = 0.0
			var br float64 = 0.0

			tl = rawDither[bufY_t][bufX_l]
			if bufX_r < ditherWidth {
				tr = rawDither[bufY_t][bufX_r]
			}

			if bufY_b < ditherHeight {
				bl = rawDither[bufY_b][bufX_l]
				if bufX_r < ditherWidth {
					br = rawDither[bufY_b][bufX_r]
				}
			}

			if InsideRegion(tl >= 0.5, tr >= 0.5, bl >= 0.5, br >= 0.5, math.Mod(bufX, 1.0), math.Mod(bufY, 1.0)) {
				dither[y][x] = 1.0
			}
		}
	}

	return dither
}

func WriteDitherToFile(ditherImg DitherImg, filename string) {
	Print(TERM_GREEN+"Outputting File \"%s\"\n"+TERM_RESET, filename)

	img := image.NewRGBA(image.Rect(0, 0, len(ditherImg[0]), len(ditherImg)))

	for y, row := range ditherImg {
		for x, val := range row {
			pixel := color.RGBA{0, 0, 0, 255}
			if val >= 0.5 {
				pixel = color.RGBA{255, 255, 255, 255}
			}
			img.Set(x, y, pixel)
		}
	}

	ofile := CreateFile(filename)
	png.Encode(ofile, img)
	ofile.Close()
}
