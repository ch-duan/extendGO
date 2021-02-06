package rand

import "math/rand"

func RandInt(max, min int) int {
	return min + rand.Intn(max+1-min)
}

func RandFloat(max, min float64) float64 {
	return min + rand.Float64()*(max+1-min)
}
