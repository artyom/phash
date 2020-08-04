package phash

import (
	"image"
)

func getImageMatrix(img *image.Gray) [mSize][mSize]float64 {
	xSize := img.Bounds().Max.X
	ySize := img.Bounds().Max.Y

	var vals [mSize][mSize]float64

	for x := 0; x < xSize; x++ {
		for y := 0; y < ySize; y++ {
			_, _, b, _ := img.GrayAt(x, y).RGBA()
			vals[x][y] = float64(b)
		}
	}

	return vals
}
