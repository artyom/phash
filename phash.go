// Package phash computes a phash string for a JPEG image and retrieves
// the hamming distance between two phash strings.
package phash

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"math/bits"

	"github.com/disintegration/imaging"
)

func toGray(img image.Image) *image.Gray {
	if g, ok := img.(*image.Gray); ok {
		return g
	}
	dst := image.NewGray(img.Bounds())
	draw.Draw(dst, img.Bounds(), img, image.Point{}, draw.Over)
	return dst
}

// GetHash returns a phash string for a JPEG image
func GetHash(reader io.Reader) (string, error) {
	img, err := imaging.Decode(reader)

	if err != nil {
		return "", err
	}

	img = imaging.Resize(img, mSize, mSize, imaging.Lanczos)

	imageMatrixData := getImageMatrix(toGray(img))
	dctMatrix := getDCTMatrix(imageMatrixData)

	smallDctMatrix := reduceMatrix(dctMatrix)
	dctMeanValue := calculateMeanValue(smallDctMatrix)
	return hashToString(buildHash(smallDctMatrix, dctMeanValue)), nil
}

const mSize = 32

func ImageHash(img image.Image) (string, error) {
	x, err := ImageHashUint(img)
	if err != nil {
		return "", err
	}
	return hashToString(x), nil
}

func ImageHashUint(img image.Image) (uint64, error) {
	return Get(img, func(img image.Image, w, h int) image.Image {
		return imaging.Resize(img, w, h, imaging.Lanczos)
	})
}

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

// Get calculates perceptual hash for an image. It uses scalefunc to scale
// image to 32×32 size if needed.
//
// An example of a scalefunc, using github.com/disintegration/imaging package:
//
//  scalefunc := func(img image.Image, w, h int) image.Image {
//      return imaging.Resize(img, w, h, imaging.Lanczos)
//  }
//
// Note that hash value depends highly on a scaling algorithm, smoother scaling
// algorithms usually work better.
func Get(img image.Image, scalefunc func(img image.Image, width, height int) image.Image) (uint64, error) {
	if p := img.Bounds().Size(); p.X != mSize || p.Y != mSize {
		img = scalefunc(img, mSize, mSize)
	}
	return processedImageHash(toGray(img))
}

// Distance distance between two hashes (number of bits that differ)
func Distance(h1, h2 uint64) int { return bits.OnesCount64(h1 ^ h2) }

// GetDistance returns the hamming distance between two hashes
func GetDistance(hash1, hash2 string) int {
	distance := 0
	for i := 0; i < len(hash1); i++ {
		if hash1[i] != hash2[i] {
			distance++
		}
	}

	return distance
}

func hashToString(x uint64) string { return fmt.Sprintf("%064b", x) }

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
