package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"math"
	"math/rand"
	"net/http"
	"time"
)

const (
	Width      = 60 //30 // Ширина матрицы
	Height     = 60 //30 // Высота матрицы
	pixelWidth = 3  // размеры пикселя матрицы
)

var fallenColor = color.RGBA{0, 146, 247, 255}
var needToFallColor = color.RGBA{0, 0, 255, 255} // color.RGBA{0, 116, 217, 255}
var backColor = color.RGBA{200, 200, 200, 255}
var maxVolatility = 20 // летучесть, чем выше тем ниже скорость падения
var maxObstaclesPircent float32 = 0.15

type Particle struct {
	X, Y              int
	Fallen            bool
	needToFall        bool
	volatility        int
	initialVolatility int
	rotating          int
}

var isBegin = true
var matrixParticles [][]Particle

func writeParticlesPngToResponse(resp_writer http.ResponseWriter) {

	drawParticlesToMatrixMakeOneStep()

	img := convertMatixToImg(&matrixParticles)

	savematrixParticlesImage(img, resp_writer)
}

func writeParticlesGifToResponse(resp_writer http.ResponseWriter) {

	animGif := gif.GIF{}

	const (
		delayGif = 1 // Задержка между кадрами (единица - 10мс)
	)

	for isMatrixStartedFromBegin := false; isMatrixStartedFromBegin == false; isMatrixStartedFromBegin = drawParticlesToMatrixMakeOneStep() {

		imgRGDA := convertMatixToImg(&matrixParticles)

		bounds := imgRGDA.Bounds()
		palettedImage := image.NewPaletted(bounds, palette.Plan9)
		draw.Draw(palettedImage, palettedImage.Rect, imgRGDA, bounds.Min, draw.Over)

		animGif.Image = append(animGif.Image, palettedImage)
		animGif.Delay = append(animGif.Delay, delayGif)
	}

	gif.EncodeAll(resp_writer, &animGif)
}

func drawParticlesToMatrixMakeOneStep() (isMatrixStartedFromBegin bool) {

	if isBegin {
		isBegin = false
		rand.Seed(time.Now().UnixNano())
		matrixParticles = generatematrixParticles()
		generateParticlesobstacles(matrixParticles)
		generateParticleToMove(matrixParticles)
	}

	isFirstRowsFilled := false
	for idxRow := 0; idxRow < 10; idxRow++ {
		if rowIsFilled(matrixParticles, idxRow) {
			isFirstRowsFilled = true
			break
		}
	}

	if ismatrixParticlesFull(matrixParticles) || isFirstRowsFilled {
		// firstRowFilled(matrixParticles) ||
		// secondRowFilled(matrixParticles) ||
		// thirdRowFilled(matrixParticles) {
		isBegin = true
	} else {
		moveParticles(matrixParticles)

		// if len(getPaticlesToMove(matrixParticles)) == 0 {
		// 	generateParticleToMove(matrixParticles)
		// }
		if !(ismatrixParticlesFull(matrixParticles) || secondRowFilled(matrixParticles)) {
			generateParticleToMove(matrixParticles)
		}
	}

	return isBegin
}

func generatematrixParticles() [][]Particle {
	matrixParticles := make([][]Particle, Height)
	for i := range matrixParticles {
		matrixParticles[i] = make([]Particle, Width)
	}
	return matrixParticles
}

func generateParticlesobstacles(matrixParticles [][]Particle) {

	countObstacles := rand.Intn(int(Width * Height * maxObstaclesPircent))

	for c := 0; c < countObstacles; c++ {
		x := rand.Intn(Width)
		y := rand.Intn(Height)
		matrixParticles[y][x].Fallen = true
		matrixParticles[y][x].X = x
		matrixParticles[y][x].Y = y
		matrixParticles[y][x].needToFall = false
	}
}

func generateParticleToMove(matrixParticles [][]Particle) {
	x := rand.Intn(Width)
	matrixParticles[0][x].Fallen = false
	matrixParticles[0][x].X = x
	matrixParticles[0][x].Y = 0
	matrixParticles[0][x].needToFall = true
	matrixParticles[0][x].volatility = rand.Intn(maxVolatility)
	matrixParticles[0][x].initialVolatility = matrixParticles[0][x].volatility
}

func ismatrixParticlesFull(matrixParticles [][]Particle) bool {
	for _, row := range matrixParticles {
		for _, particle := range row {
			if !particle.Fallen {
				return false
			}
		}
	}
	return true
}

func firstRowFilled(matrixParticles [][]Particle) bool {
	return rowIsFilled(matrixParticles, 0)
}

func secondRowFilled(matrixParticles [][]Particle) bool {
	return rowIsFilled(matrixParticles, 1)
}

func thirdRowFilled(matrixParticles [][]Particle) bool {
	return rowIsFilled(matrixParticles, 2)
}

func rowIsFilled(matrixParticles [][]Particle, rowIdx int) bool {

	for _, particle := range matrixParticles[rowIdx] {
		if !particle.Fallen {
			return false
		}
	}
	return true
}

func moveParticles(matrixParticles [][]Particle) {

	decreaseVolatility(matrixParticles)

	var particleToMove []Particle = getPaticlesToMove(matrixParticles)
	if len(particleToMove) == 0 {
		return
	}

	for _, particle := range particleToMove {
		moveParticle(&particle, matrixParticles)
	}
}

func decreaseVolatility(matrixParticles [][]Particle) {

	var particle *Particle

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			particle = &matrixParticles[y][x]
			if particle.Fallen {
				continue
			} else if particle.volatility > 0 {
				particle.volatility--
			} else if particle.volatility == 0 && particle.initialVolatility > 0 {
				particle.volatility = particle.initialVolatility
			}
		}
	}
}

func getPaticlesToMove(matrixParticles [][]Particle) []Particle {

	var particle *Particle
	var particleToMove []Particle = nil

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			particle = &matrixParticles[y][x]
			if particle.Fallen {
				continue
			}
			if particle.volatility > 0 {
				continue
			}
			if particle.needToFall {
				particleToMove = append(particleToMove, *particle)
			}
		}
	}

	return particleToMove
}

func moveParticle(particle *Particle, matrixParticles [][]Particle) {

	y := particle.Y
	x := particle.X

	movingParticle := &matrixParticles[y][x]

	bottomCellFill := y == Height-1 || matrixParticles[y+1][x].Fallen
	bottomRightCellFill := y == Height-1 || matrixParticles[y+1][int(math.Min(float64(x+1), float64(Width-1)))].Fallen || matrixParticles[y+1][int(math.Min(float64(x+1), float64(Width-1)))].needToFall
	bottomLeftCellFill := y == Height-1 || matrixParticles[y+1][int(math.Max(float64(x-1), float64(0)))].Fallen || matrixParticles[y+1][int(math.Max(float64(x-1), float64(0)))].needToFall
	rightCellFill := matrixParticles[y][int(math.Min(float64(x+1), float64(Width-1)))].Fallen || matrixParticles[y][int(math.Min(float64(x+1), float64(Width-1)))].needToFall
	leftCellFill := matrixParticles[y][int(math.Max(float64(x-1), float64(0)))].Fallen || matrixParticles[y][int(math.Max(float64(x-1), float64(0)))].needToFall

	if y == Height-1 ||
		bottomCellFill && bottomRightCellFill && bottomLeftCellFill ||
		bottomCellFill && rightCellFill && leftCellFill {

		movingParticle.Fallen = true
		movingParticle.needToFall = false
	} else {

		movingParticle.Fallen = false
		movingParticle.needToFall = false

		if bottomCellFill && !bottomLeftCellFill && !leftCellFill {
			x--
			y++
		} else if bottomCellFill && !bottomRightCellFill && !rightCellFill {
			x++
			y++
		} else if !bottomCellFill {
			y++
		} else if bottomCellFill {
			movingParticle.Fallen = true
			movingParticle.needToFall = false
		}

		if !movingParticle.Fallen {
			initialVolatility := movingParticle.initialVolatility
			rotating := movingParticle.rotating
			movingParticle = &matrixParticles[y][x]
			movingParticle.X = x
			movingParticle.Y = y
			movingParticle.Fallen = false
			movingParticle.needToFall = true
			movingParticle.initialVolatility = initialVolatility
			movingParticle.rotating = rotating
		}
	}
}

func convertMatixToImg(matrixParticles *[][]Particle) *image.RGBA {

	width := Width * pixelWidth
	height := Height * pixelWidth

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var curParticiple *Particle

	for y, row := range *matrixParticles {
		for x, particle := range row {

			var col color.RGBA
			isFalling := false

			curParticiple = &(*matrixParticles)[y][x]

			if particle.needToFall {
				col = needToFallColor
				isFalling = true
			} else if particle.Fallen {
				col = fallenColor
			} else {
				col = backColor
			}
			drawParticle(img, pixelWidth, x, y, col, isFalling, curParticiple)
		}
	}
	return img
}

func drawParticle(img *image.RGBA, pixelWidth int, x int, y int, col color.RGBA, isFalling bool, particle *Particle) {

	if isFalling {

		switch particle.rotating {

		case 0:
			img.Set(x*pixelWidth+1, y*pixelWidth+0, col)

			img.Set(x*pixelWidth+0, y*pixelWidth+1, col)
			img.Set(x*pixelWidth+1, y*pixelWidth+1, col)
			img.Set(x*pixelWidth+2, y*pixelWidth+1, col)

			img.Set(x*pixelWidth+1, y*pixelWidth+2, col)
		case 1:
			img.Set(x*pixelWidth+2, y*pixelWidth+0, col)

			img.Set(x*pixelWidth+0, y*pixelWidth+1, col)
			img.Set(x*pixelWidth+1, y*pixelWidth+1, col)
			img.Set(x*pixelWidth+2, y*pixelWidth+1, col)

			img.Set(x*pixelWidth+1, y*pixelWidth+2, col)
		case 2:
			img.Set(x*pixelWidth+2, y*pixelWidth+0, col)

			img.Set(x*pixelWidth+0, y*pixelWidth+1, col)
			img.Set(x*pixelWidth+1, y*pixelWidth+1, col)

			img.Set(x*pixelWidth+1, y*pixelWidth+2, col)
			img.Set(x*pixelWidth+2, y*pixelWidth+2, col)
		case 3:
			img.Set(x*pixelWidth+2, y*pixelWidth+0, col)

			img.Set(x*pixelWidth+0, y*pixelWidth+1, col)
			img.Set(x*pixelWidth+1, y*pixelWidth+1, col)

			img.Set(x*pixelWidth+0, y*pixelWidth+2, col)
			img.Set(x*pixelWidth+2, y*pixelWidth+2, col)
		case 4:
			img.Set(x*pixelWidth+0, y*pixelWidth+0, col)
			img.Set(x*pixelWidth+2, y*pixelWidth+0, col)

			img.Set(x*pixelWidth+1, y*pixelWidth+1, col)

			img.Set(x*pixelWidth+0, y*pixelWidth+2, col)
			img.Set(x*pixelWidth+2, y*pixelWidth+2, col)

		case 5:
			img.Set(x*pixelWidth+0, y*pixelWidth+0, col)

			img.Set(x*pixelWidth+1, y*pixelWidth+1, col)
			img.Set(x*pixelWidth+2, y*pixelWidth+1, col)

			img.Set(x*pixelWidth+0, y*pixelWidth+2, col)
			img.Set(x*pixelWidth+2, y*pixelWidth+2, col)
		case 6:
			img.Set(x*pixelWidth+0, y*pixelWidth+0, col)

			img.Set(x*pixelWidth+1, y*pixelWidth+1, col)
			img.Set(x*pixelWidth+2, y*pixelWidth+1, col)

			img.Set(x*pixelWidth+0, y*pixelWidth+2, col)
			img.Set(x*pixelWidth+1, y*pixelWidth+2, col)
		case 7:
			img.Set(x*pixelWidth+0, y*pixelWidth+0, col)

			img.Set(x*pixelWidth+0, y*pixelWidth+1, col)
			img.Set(x*pixelWidth+1, y*pixelWidth+1, col)
			img.Set(x*pixelWidth+2, y*pixelWidth+1, col)

			img.Set(x*pixelWidth+1, y*pixelWidth+2, col)
		}

		particle.rotating++

		if particle.rotating > 7 {
			particle.rotating = 0
		}

		// case 0:
		// 	img.Set(x*pixelWidth+0, y*pixelWidth+0, col)
		// 	img.Set(x*pixelWidth+1, y*pixelWidth+0, col)
		// 	img.Set(x*pixelWidth+2, y*pixelWidth+0, col)

		// 	img.Set(x*pixelWidth+0, y*pixelWidth+1, col)
		// 	img.Set(x*pixelWidth+1, y*pixelWidth+1, col)
		// 	img.Set(x*pixelWidth+2, y*pixelWidth+1, col)

		// 	img.Set(x*pixelWidth+0, y*pixelWidth+2, col)
		// 	img.Set(x*pixelWidth+1, y*pixelWidth+2, col)
		// 	img.Set(x*pixelWidth+2, y*pixelWidth+2, col)

		// if particle.rotating != 0 {

		// 	particle.rotating = 0

		// 	img.Set(x*pixelWidth+1, y*pixelWidth+0, col)

		// 	img.Set(x*pixelWidth+0, y*pixelWidth+1, col)
		// 	img.Set(x*pixelWidth+1, y*pixelWidth+1, col)
		// 	img.Set(x*pixelWidth+2, y*pixelWidth+1, col)

		// 	img.Set(x*pixelWidth+1, y*pixelWidth+2, col)

		// } else {

		// 	particle.rotating = 1

		// 	img.Set(x*pixelWidth+0, y*pixelWidth+0, col)
		// 	img.Set(x*pixelWidth+2, y*pixelWidth+0, col)

		// 	img.Set(x*pixelWidth+1, y*pixelWidth+1, col)

		// 	img.Set(x*pixelWidth+0, y*pixelWidth+2, col)
		// 	img.Set(x*pixelWidth+2, y*pixelWidth+2, col)
		// }

	} else {
		for i := 0; i < pixelWidth; i++ {
			for j := 0; j < pixelWidth; j++ {
				img.Set(x*pixelWidth+i, y*pixelWidth+j, col)
			}
		}
	}
}

func savematrixParticlesImage(img *image.RGBA, w http.ResponseWriter) {

	w.Header().Set("Content-Type", "image/png")
	if err := png.Encode(w, img); err != nil {
		fmt.Println(err)
		return
	}
}
