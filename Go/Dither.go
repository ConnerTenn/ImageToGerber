package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
)

type DitherImg [][]float64

func CreateKernel() [][]float64 {
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

	return kernel
}

func GenerateDither(width int, height int, factor float64) DitherImg {
	Print("Generating Dither [%dx%d] factor:%f\n", width, height, factor)

	dither := make(DitherImg, height)
	for y := range dither {
		dither[y] = make([]float64, width)
		for x := range dither[y] {
			//Add a bit of random adjustment as a 'seed' for the dither
			//Without it, the pattern would be too uniform
			dither[y][x] = (0.1 * (rand.Float64() - 0.5)) + factor
		}
	}

	//https://en.wikipedia.org/wiki/Floyd%E2%80%93Steinberg_dithering
	//Floyd-Steinberg Dithering
	kernel := CreateKernel()

	//Perform a convolution-like operation over the image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			old := dither[y][x]

			//Clamp to 0 or 1
			new := 0.0
			if old >= 0.5 {
				new = 1.0
			}

			//Set the pixel
			dither[y][x] = new

			//calculate the error
			quantError := old - float64(new)

			//Convolve with the kernel
			for yoff, kernelRow := range kernel {
				if y+yoff < height { //Bounds check

					for xoff, kfactor := range kernelRow {
						xoff = xoff - len(kernelRow)/2     //Kernel is centered horizontally about the current pixel
						if x+xoff >= 0 && x+xoff < width { //Bounds check

							//Carry over error to neghboring pixels, weighted by the kernel factor
							dither[y+yoff][x+xoff] = dither[y+yoff][x+xoff] + quantError*kfactor
						}
					}
				}
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
