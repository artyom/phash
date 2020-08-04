package phash

import (
	"image"
)

func getImageMatrix(img image.Image) [][]float64 {
	xSize := img.Bounds().Max.X
	ySize := img.Bounds().Max.Y

	vals := make([][]float64, xSize)

	var getXYValue func(x, y int) float64
	switch img2 := img.(type) {
	case *image.NRGBA:
		getXYValue = func(x, y int) float64 {
			_, _, b, _ := img2.NRGBAAt(x, y).RGBA()
			return float64(b)
		}
	default:
		getXYValue = func(x, y int) float64 {
			_, _, b, _ := img.At(x, y).RGBA()
			return float64(b)
		}
	}

	for x := 0; x < xSize; x++ {
		vals[x] = make([]float64, ySize)
		for y := 0; y < ySize; y++ {
			vals[x][y] = getXYValue(x, y)
		}
	}

	return vals
}
