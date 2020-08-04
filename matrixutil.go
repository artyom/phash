package phash

const sSize = 8

func reduceMatrix(dctMatrix [mSize][mSize]float64) [sSize][sSize]float64 {
	var newMatrix [sSize][sSize]float64
	for x := 0; x < sSize; x++ {
		for y := 0; y < sSize; y++ {
			newMatrix[x][y] = dctMatrix[x][y]
		}
	}
	return newMatrix
}

func calculateMeanValue(dctMatrix [sSize][sSize]float64) float64 {
	var total float64
	for x := 0; x < sSize; x++ {
		for y := 0; y < sSize; y++ {
			total += dctMatrix[x][y]
		}
	}
	total -= dctMatrix[0][0]
	avg := total / float64((sSize*sSize)-1)
	return avg
}
