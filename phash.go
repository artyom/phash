// Package phash computes a perceptual hash of an image.
//
// Hash values are highly dependent on the scaling function used.
package phash

import (
	"image"
	"image/draw"
	"math/bits"
)

// Get calculates perceptual hash for an image. It uses scalefunc to scale
// image to 32×32 size if needed.
func Get(img image.Image, scalefunc ScaleFunc) (uint64, error) {
	if p := img.Bounds().Size(); p.X != mSize || p.Y != mSize {
		img = scalefunc(img, mSize, mSize)
	}
	return processedImageHash(toGray(img))
}

// Distance distance between two hashes (number of bits that differ)
func Distance(h1, h2 uint64) int { return bits.OnesCount64(h1 ^ h2) }

// ScaleFunc is a function Get uses to scale image to a 32×32 if needed.
//
// An example of a ScaleFunc, using github.com/disintegration/imaging package:
//
//  scalefunc := func(img image.Image, w, h int) image.Image {
//      return imaging.Resize(img, w, h, imaging.Lanczos)
//  }
//
// Note that hash value depends highly on a scaling algorithm, smoother scaling
// algorithms usually work better.
//
// Function must return a 32×32 image with its top left point having 0, 0
// coordinates.
type ScaleFunc func(img image.Image, width, height int) image.Image

const mSize = 32

// processedImageHash must be called on a 32×32 greyscale image
func processedImageHash(img *image.Gray) (uint64, error) {
	if p := img.Rect.Size(); p.X != mSize || p.Y != mSize {
		panic("image dimensions are not 32×32")
	}
	imageMatrixData := getImageMatrix(img)
	dctMatrix := getDCTMatrix(imageMatrixData)
	smallDctMatrix := reduceMatrix(dctMatrix)
	dctMeanValue := calculateMeanValue(smallDctMatrix)
	return buildHash(smallDctMatrix, dctMeanValue), nil
}

func buildHash(dctMatrix [sSize][sSize]float64, dctMeanValue float64) uint64 {
	var b uint64
	var i int = 63
	for x := 0; x < sSize; x++ {
		for y := 0; y < sSize; y++ {
			if dctMatrix[x][y] > dctMeanValue {
				b ^= 1 << i
			}
			i--
		}
	}
	return b
}

func toGray(img image.Image) *image.Gray {
	if g, ok := img.(*image.Gray); ok {
		return g
	}
	dst := image.NewGray(img.Bounds())
	draw.Draw(dst, img.Bounds(), img, image.Point{}, draw.Over)
	return dst
}
