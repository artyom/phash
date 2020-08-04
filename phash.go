// Package phash computes a phash string for a JPEG image and retrieves
// the hamming distance between two phash strings.
package phash

import (
	"fmt"
	"image"
	"io"

	"github.com/disintegration/imaging"
)

// GetHash returns a phash string for a JPEG image
func GetHash(reader io.Reader) (string, error) {
	image, err := imaging.Decode(reader)

	if err != nil {
		return "", err
	}

	image = imaging.Resize(image, 32, 32, imaging.Lanczos)
	image = imaging.Grayscale(image)

	imageMatrixData := getImageMatrix(image)
	dctMatrix := getDCTMatrix(imageMatrixData)

	smallDctMatrix := reduceMatrix(dctMatrix)
	dctMeanValue := calculateMeanValue(smallDctMatrix)
	return buildHashString(smallDctMatrix, dctMeanValue), nil
}

const mSize = 32

func ImageHash(img image.Image) (string, error) {
	p := img.Bounds().Size()
	if p.X != mSize || p.Y != mSize {
		img = imaging.Resize(img, mSize, mSize, imaging.Lanczos)
	}
	img = imaging.Grayscale(img)
	return processedImageHash(img)
}

// processedImageHash must be called on a 32×32 greyscale image
func processedImageHash(img image.Image) (string, error) {
	imageMatrixData := getImageMatrix(img)
	dctMatrix := getDCTMatrix(imageMatrixData)
	smallDctMatrix := reduceMatrix(dctMatrix)
	dctMeanValue := calculateMeanValue(smallDctMatrix)
	return buildHashString(smallDctMatrix, dctMeanValue), nil
}

// GetDistance returns the hamming distance between two phashes
func GetDistance(hash1, hash2 string) int {
	distance := 0
	for i := 0; i < len(hash1); i++ {
		if hash1[i] != hash2[i] {
			distance++
		}
	}

	return distance
}

func buildHashString(dctMatrix [sSize][sSize]float64, dctMeanValue float64) string {
	return fmt.Sprintf("%064b", buildHash(dctMatrix, dctMeanValue))
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
