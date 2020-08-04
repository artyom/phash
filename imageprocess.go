package phash

import (
	"image"
)

func getImageMatrix(img *image.Gray) [mSize][mSize]float64 {
	var vals [mSize][mSize]float64
	// TODO: make sure img.Rect starts at zero point
	for x := 0; x < mSize; x++ {
		for y := 0; y < mSize; y++ {
			_, _, b, _ := img.GrayAt(x, y).RGBA()
			vals[x][y] = float64(b)
		}
	}
	return vals
}
