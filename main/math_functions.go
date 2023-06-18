package main

func GeometricProgression(number int64, ratio int64) int64 {

	var result int64 = 1

	var c int64
	for c = 1; c <= number; c++ {
		result *= ratio
	}

	return result
}

func ArifmeticProgression(number int64) int64 {

	var result int64 = 0

	var c int64
	for c = 0; c <= number; c++ {
		result += c
	}

	return result
}
