package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"net/http"
	"time"
)

const (
	Width      = 30 // Ширина матрицы
	Height     = 30 // Высота матрицы
	pixelWidth = 5
)

var fallenColor = color.RGBA{0, 0, 0, 255}
var needToFallColor = color.RGBA{0, 0, 0, 255}
var backColor = color.RGBA{255, 255, 255, 255}

type Particle struct {
	X, Y       int
	Fallen     bool
	needToFall bool
}

var isBegin = true
var matrix [][]Particle

func drawParticles(resp_writer http.ResponseWriter) {

	if isBegin {
		isBegin = false
		rand.Seed(time.Now().UnixNano())
		matrix = generateMatrix()
		generateParticle(matrix)
	}

	if isMatrixFull(matrix) {
		isBegin = true
	} else {
		updateParticles(matrix)
	}

	saveMatrixImage(matrix, resp_writer)
}

func generateMatrix() [][]Particle {
	matrix := make([][]Particle, Height)
	for i := range matrix {
		matrix[i] = make([]Particle, Width)
	}
	return matrix
}

func isMatrixFull(matrix [][]Particle) bool {
	for _, row := range matrix {
		for _, particle := range row {
			if !particle.Fallen {
				return false
			}
		}
	}
	return true
}

func saveMatrixImage(matrix [][]Particle, w http.ResponseWriter) {
	width := Width * pixelWidth
	height := Height * pixelWidth

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y, row := range matrix {
		for x, particle := range row {
			var col color.RGBA
			if particle.needToFall {
				col = needToFallColor
			} else if particle.Fallen {
				col = fallenColor
			} else {
				col = backColor
			}
			for i := 0; i < pixelWidth; i++ {
				for j := 0; j < pixelWidth; j++ {
					img.Set(x*pixelWidth+i, y*pixelWidth+j, col)
				}
			}
		}
	}

	w.Header().Set("Content-Type", "image/png")
	if err := png.Encode(w, img); err != nil {
		fmt.Println(err)
		return
	}
}

func generateParticle(matrix [][]Particle) {
	x := rand.Intn(Width)
	matrix[0][x].Fallen = false
	matrix[0][x].X = x
	matrix[0][x].Y = 0
	matrix[0][x].needToFall = true
}

func updateParticles(matrix [][]Particle) {

	var particle *Particle

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			particle = &matrix[y][x]
			if particle.Fallen {
				continue
			}
			if particle.needToFall {
				break
			}
		}
		if particle.needToFall {
			break
		}
	}

	if !particle.needToFall {
		return
	}

	y := particle.Y
	x := particle.X

	bottomCellFill := y == Height-1 || matrix[y+1][x].Fallen
	bottomRightCellFill := y == Height-1 || matrix[y+1][int(math.Min(float64(x+1), float64(Width-1)))].Fallen
	bottomLeftCellFill := y == Height-1 || matrix[y+1][int(math.Max(float64(x-1), float64(0)))].Fallen

	if y == Height-1 ||
		bottomCellFill && bottomRightCellFill && bottomLeftCellFill {

		particle.Fallen = true
		particle.needToFall = false
		if !isMatrixFull(matrix) {
			generateParticle(matrix)
		}

	} else {

		particle.Fallen = false
		particle.needToFall = false

		// if x > 0 && matrix[y][x-1].Fallen && !matrix[y+1][x-1].Fallen {
		// 	x--
		// 	y++
		// } else if x < Width-1 && matrix[y][x+1].Fallen && !matrix[y+1][x+1].Fallen {
		// 	x++
		// 	y++
		// } else if x > 0 && x < Width-1 && matrix[y+1][x].Fallen && !matrix[y+1][x-1].Fallen {
		// 	x--
		// 	y++
		// } else if x > 0 && x < Width-1 && matrix[y+1][x].Fallen && !matrix[y+1][x+1].Fallen {
		// 	x++
		// 	y++

		if bottomCellFill && !bottomLeftCellFill {
			x--
			y++
		} else if bottomCellFill && !bottomRightCellFill {
			x++
			y++
			// } else if !bottomLeftCellFill {
			// 	x--
			// 	y++
			// } else if !bottomRightCellFill {
			// 	x++
			// 	y++
		} else if !bottomCellFill {
			y++
		}

		particle = &matrix[y][x]
		particle.X = x
		particle.Y = y
		particle.Fallen = false
		particle.needToFall = true
	}
}
