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

	"github.com/disintegration/imaging"
)

const (
// Width      = 60 //30 // Ширина матрицы
// Height     = 60 //30 // Высота матрицы
// pixelWidth = 3 // размеры пикселя матрицы
)

var fallenColor = color.RGBA{0, 146, 247, 255}
var needToFallColor = color.RGBA{0, 0, 255, 255} // color.RGBA{0, 116, 217, 255}
var backColor = color.RGBA{200, 200, 200, 255}
var maxVolatility = 20 // летучесть, чем выше тем ниже скорость падения
// var maxObstaclesPircent float32 = 0.05 // 0.15
var imageRotateAngle float64 = 1.0

type Particle struct {
	X, Y              int
	Fallen            bool
	needToFall        bool
	volatility        int
	initialVolatility int
	rotating          int
}

// var isBegin = true
type matrixParticlesPtrType *[][]Particle

var matrixParticlesWithOutRotating matrixParticlesPtrType
var matrixParticlesRotated matrixParticlesPtrType

func writeParticlesPngToResponse(resp_writer http.ResponseWriter, rotating bool) {

	var img *image.RGBA

	if rotating {
		matrixParticlesRotated, _ = drawParticlesToMatrixMakeOneStep(matrixParticlesRotated, rotating)
		img = convertMatixToImg(matrixParticlesRotated)
	} else {
		matrixParticlesWithOutRotating, _ = drawParticlesToMatrixMakeOneStep(matrixParticlesWithOutRotating, rotating)
		img = convertMatixToImg(matrixParticlesWithOutRotating)
	}

	savematrixParticlesImage(img, resp_writer)
}

func writeParticlesGifToResponse(resp_writer http.ResponseWriter, rotating bool) {

	animGif := gif.GIF{}

	const (
		delayGif = 1 // Задержка между кадрами (единица - 10мс)
	)

	if rotating {
		for isMatrixStartedFromBegin := false; isMatrixStartedFromBegin == false; matrixParticlesRotated, isMatrixStartedFromBegin = drawParticlesToMatrixMakeOneStep(matrixParticlesRotated, rotating) {

			imgRGDA := convertMatixToImg(matrixParticlesRotated)

			bounds := imgRGDA.Bounds()
			palettedImage := image.NewPaletted(bounds, palette.Plan9)
			draw.Draw(palettedImage, palettedImage.Rect, imgRGDA, bounds.Min, draw.Over)

			animGif.Image = append(animGif.Image, palettedImage)
			animGif.Delay = append(animGif.Delay, delayGif)
		}
	} else {
		for isMatrixStartedFromBegin := false; isMatrixStartedFromBegin == false; matrixParticlesWithOutRotating, isMatrixStartedFromBegin = drawParticlesToMatrixMakeOneStep(matrixParticlesWithOutRotating, rotating) {

			imgRGDA := convertMatixToImg(matrixParticlesWithOutRotating)

			bounds := imgRGDA.Bounds()
			palettedImage := image.NewPaletted(bounds, palette.Plan9)
			draw.Draw(palettedImage, palettedImage.Rect, imgRGDA, bounds.Min, draw.Over)

			animGif.Image = append(animGif.Image, palettedImage)
			animGif.Delay = append(animGif.Delay, delayGif)
		}
	}

	gif.EncodeAll(resp_writer, &animGif)
}

func drawParticlesToMatrixMakeOneStep(matrixParticles matrixParticlesPtrType, rotating bool) (matrixParticlesPtrType, bool) {

	if ismatrixParticlesFree(matrixParticles) || ismatrixParticlesFull(matrixParticles) {
		rand.Seed(time.Now().UnixNano())
		matrixParticles = generateMatrixParticles(matrixParticles)
		generateParticlesobstacles(matrixParticles)
		generateBorder(matrixParticles)
		generateParticleToMove(matrixParticles, rotating)
	}

	moveParticles(matrixParticles, rotating)
	generateParticleToMove(matrixParticles, rotating)

	return matrixParticles, ismatrixParticlesFull(matrixParticles)
}

func generateMatrixParticles(matrixParticles matrixParticlesPtrType) matrixParticlesPtrType {

	Height := 60
	Width := 60

	newMatrixParticles := make([][]Particle, Height)
	for i := range newMatrixParticles {
		newMatrixParticles[i] = make([]Particle, Width)
	}

	matrixParticles = &newMatrixParticles

	return matrixParticles
}

func generateParticlesobstacles(matrixParticles matrixParticlesPtrType) {

	Height := len(*matrixParticles)
	Width := len((*matrixParticles)[Height-1])
	maxObstaclesPircent := 0.05

	countObstacles := rand.Intn(int(float64(Width) * float64(Height) * maxObstaclesPircent))

	for c := 0; c < countObstacles; c++ {
		x := rand.Intn(Width)
		y := rand.Intn(Height)
		(*matrixParticles)[y][x].Fallen = true
		(*matrixParticles)[y][x].X = x
		(*matrixParticles)[y][x].Y = y
		(*matrixParticles)[y][x].needToFall = false
	}
}

func generateBorder(matrixParticles matrixParticlesPtrType) {

	Height := len(*matrixParticles)
	Width := len((*matrixParticles)[Height-1])

	for row := 0; row < Height; row++ {
		for col := 0; col < Width; col++ {

			if row != 0 && row != (Height-1) {
				if col != 0 && col != (Width-1) {
					continue
				}
			}

			// var cell *Particle = &matrixParticles[col][row]
			// cell.Fallen = true
			// cell.X = col
			// cell.Y = row
			// cell.needToFall = false

			(*matrixParticles)[row][col].Fallen = true
			(*matrixParticles)[row][col].X = col
			(*matrixParticles)[row][col].Y = row
			(*matrixParticles)[row][col].needToFall = false
		}
	}
}

func generateParticleToMove(matrixParticles matrixParticlesPtrType, generateInCenter bool) {

	Height := len((*matrixParticles))
	Width := len((*matrixParticles)[Height-1])

	var x int
	var y int

	if generateInCenter {
		x = rand.Intn(int(float64(Width)*0.4)) + int(float64(Width)/2.0) - int(float64(Width)*0.2)
		y = int(Height / 2.0)
	} else {
		x = rand.Intn(Width-10) + 5
		y = 1
	}

	(*matrixParticles)[y][x].Fallen = false
	(*matrixParticles)[y][x].X = x
	(*matrixParticles)[y][x].Y = y
	(*matrixParticles)[y][x].needToFall = true
	(*matrixParticles)[y][x].volatility = rand.Intn(maxVolatility)
	(*matrixParticles)[y][x].initialVolatility = (*matrixParticles)[y][x].volatility

}

func ismatrixParticlesFree(matrixParticles matrixParticlesPtrType) bool {

	if matrixParticles == nil {
		return true
	}

	for y := 1; y < len(*matrixParticles)-1; y++ {
		for x := 1; x < len((*matrixParticles)[y])-1; x++ {
			if (*matrixParticles)[y][x].Fallen {
				return false
			}
		}
	}

	return true
}

func ismatrixParticlesFull(matrixParticles matrixParticlesPtrType) bool {

	if matrixParticles == nil {
		return false
	}

	// for _, row := range matrixParticles {
	// 	for _, particle := range row {
	// 		if !particle.Fallen {
	// 			return false
	// 		}
	// 	}
	// }
	// return true

	Height := len(*matrixParticles)
	Width := len((*matrixParticles)[Height-1])

	totalCountOfParticles := Height * Width
	countParticlesToBeFree := int(float64(totalCountOfParticles) * 0.8)
	countFallenParticles := 0

	for _, row := range *matrixParticles {
		for _, particle := range row {
			if particle.Fallen {
				countFallenParticles = countFallenParticles + 1
			}
		}
	}

	isMatixFull := countFallenParticles >= countParticlesToBeFree

	return isMatixFull
}

func firstRowFilled(matrixParticles matrixParticlesPtrType) bool {
	return rowIsFilled(matrixParticles, 0)
}

func secondRowFilled(matrixParticles matrixParticlesPtrType) bool {
	return rowIsFilled(matrixParticles, 1)
}

func thirdRowFilled(matrixParticles matrixParticlesPtrType) bool {
	return rowIsFilled(matrixParticles, 2)
}

func rowIsFilled(matrixParticles matrixParticlesPtrType, rowIdx int) bool {

	for _, particle := range (*matrixParticles)[rowIdx] {
		if !particle.Fallen {
			return false
		}
	}
	return true
}

func moveParticles(matrixParticles matrixParticlesPtrType, rotating bool) {

	decreaseVolatility(matrixParticles)

	var particleToMove []Particle = getPaticlesToMove(matrixParticles)
	if len(particleToMove) == 0 {
		return
	}

	for _, particle := range particleToMove {
		moveParticle(&particle, matrixParticles, rotating)
	}
}

func decreaseVolatility(matrixParticles matrixParticlesPtrType) {

	var particle *Particle

	Height := len(*matrixParticles)
	Width := len((*matrixParticles)[Height-1])

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			particle = &(*matrixParticles)[y][x]
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

func getPaticlesToMove(matrixParticles matrixParticlesPtrType) []Particle {

	var particle *Particle
	var particleToMove []Particle = nil

	Height := len((*matrixParticles))
	Width := len((*matrixParticles)[Height-1])

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			particle = &(*matrixParticles)[y][x]
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

func realRotateAngle() float64 {
	return imageRotateAngle + 90
}

func moveParticle(particle *Particle, matrixParticles matrixParticlesPtrType, rotating bool) {

	if rotating {
		moveParticleRotate(particle, matrixParticles)
	} else {
		moveParticleWithOutRotate(particle, matrixParticles)
	}

}

func moveParticleRotate(particle *Particle, matrixParticles matrixParticlesPtrType) {

	Height := len((*matrixParticles))
	Width := len((*matrixParticles)[Height-1])

	y := particle.Y
	x := particle.X

	var currCellY int = int(math.Round(float64(y)))
	// var currCellX int = int(math.Round(float64(x)))

	var downCellY int = int(math.Round(float64(y) + 1.0*math.Sin(realRotateAngle()*math.Pi/180)))
	var downCellX int = int(math.Round(float64(x) + 1.0*math.Cos(realRotateAngle()*math.Pi/180)))

	var leftCellX int = int(math.Round(float64(x) - 1.0*math.Cos(realRotateAngle()*math.Pi/180)))
	var rightCellX int = int(math.Round(float64(x) + 1.0*math.Cos(realRotateAngle()*math.Pi/180)))

	// currCellY = int(math.Min(float64(currCellY), Height-1))
	// downCellY = int(math.Min(float64(downCellY), Height-1))

	// currCellX = int(math.Min(float64(currCellX), Width-1))
	// leftCellX = int(math.Min(float64(leftCellX), Width-1))
	// rightCellX = int(math.Min(float64(rightCellX), Width-1))

	// currCellY = int(math.Max(float64(currCellY), 0))
	// downCellY = int(math.Max(float64(downCellY), 0))

	// currCellX = int(math.Max(float64(currCellX), 0))
	// leftCellX = int(math.Max(float64(leftCellX), 0))
	// rightCellX = int(math.Max(float64(rightCellX), 0))

	// var downCellY int = y + 1
	// var leftCellX int = x - 1
	// var rightCellX int = x + 1

	movingParticle := &(*matrixParticles)[y][x]

	// bottomCellFill := y == Height-1 || matrixParticles[downCellY][x].Fallen
	// bottomRightCellFill := y == Height-1 || matrixParticles[downCellY][int(math.Min(float64(rightCellX), float64(Width-1)))].Fallen || matrixParticles[downCellY][int(math.Min(float64(rightCellX), float64(Width-1)))].needToFall
	// bottomLeftCellFill := y == Height-1 || matrixParticles[downCellY][int(math.Max(float64(leftCellX), float64(0)))].Fallen || matrixParticles[downCellY][int(math.Max(float64(leftCellX), float64(0)))].needToFall
	// rightCellFill := matrixParticles[y][int(math.Min(float64(rightCellX), float64(Width-1)))].Fallen || matrixParticles[y][int(math.Min(float64(rightCellX), float64(Width-1)))].needToFall
	// leftCellFill := matrixParticles[y][int(math.Max(float64(leftCellX), float64(0)))].Fallen || matrixParticles[y][int(math.Max(float64(leftCellX), float64(0)))].needToFall

	bottomCellFill := (*matrixParticles)[downCellY][downCellX].Fallen
	bottomRightCellFill := (*matrixParticles)[downCellY][rightCellX].Fallen || (*matrixParticles)[downCellY][rightCellX].needToFall
	bottomLeftCellFill := (*matrixParticles)[downCellY][leftCellX].Fallen || (*matrixParticles)[downCellY][leftCellX].needToFall
	rightCellFill := (*matrixParticles)[currCellY][rightCellX].Fallen || (*matrixParticles)[currCellY][rightCellX].needToFall
	leftCellFill := (*matrixParticles)[currCellY][leftCellX].Fallen || (*matrixParticles)[currCellY][leftCellX].needToFall

	// if y == Height-1 ||
	// 	bottomCellFill && bottomRightCellFill && bottomLeftCellFill ||
	// 	bottomCellFill && rightCellFill && leftCellFill {

	// 	movingParticle.Fallen = true
	// 	movingParticle.needToFall = false
	// } else {

	if bottomCellFill && bottomRightCellFill && bottomLeftCellFill ||
		bottomCellFill && rightCellFill && leftCellFill {

		movingParticle.Fallen = true
		movingParticle.needToFall = false

	} else {

		movingParticle.Fallen = false
		movingParticle.needToFall = false

		// // if bottomCellFill && !bottomLeftCellFill && !leftCellFill {
		// // 	x = int(math.Round(float64(x) - 1.0*math.Cos(realRotateAngle()*math.Pi/180)))
		// // 	y = int(math.Round(float64(y) + 1.0*math.Sin(realRotateAngle()*math.Pi/180)))
		// // 	// x = x - 1
		// // 	// y = y + 1
		// // } else if bottomCellFill && !bottomRightCellFill && !rightCellFill {
		// // 	x = int(math.Round(float64(x) + 1.0*math.Cos(realRotateAngle()*math.Pi/180)))
		// // 	y = int(math.Round(float64(y) + 1.0*math.Sin(realRotateAngle()*math.Pi/180)))
		// // 	// x = x + 1
		// // 	// y = y + 1
		// // } else if !bottomCellFill {
		// x = int(math.Round(float64(x) + 1.0*math.Cos(realRotateAngle()*math.Pi/180)))
		// y = int(math.Round(float64(y) + 1.0*math.Sin(realRotateAngle()*math.Pi/180)))
		// // y = y + 1
		// // } else if bottomCellFill {
		// // 	movingParticle.Fallen = true
		// // 	movingParticle.needToFall = false
		// // }

		// if bottomCellFill && !bottomLeftCellFill && !leftCellFill {
		// 	x = int(math.Round(float64(x) - 1.0*math.Cos(realRotateAngle()*math.Pi/180)))
		// 	y = int(math.Round(float64(y) + 1.0*math.Sin(realRotateAngle()*math.Pi/180)))
		// 	// x = x - 1
		// 	// y = y + 1
		// } else if bottomCellFill && !bottomRightCellFill && !rightCellFill {
		// 	x = int(math.Round(float64(x) + 1.0*math.Cos(realRotateAngle()*math.Pi/180)))
		// 	y = int(math.Round(float64(y) + 1.0*math.Sin(realRotateAngle()*math.Pi/180)))
		// 	// x = x + 1
		// 	// y = y + 1
		// } else

		if !bottomCellFill {
			x = int(math.Round(float64(x) + 1.0*math.Cos(realRotateAngle()*math.Pi/180)))
			y = int(math.Round(float64(y) + 1.0*math.Sin(realRotateAngle()*math.Pi/180)))
			// y = y + 1
		} else if bottomCellFill {
			movingParticle.Fallen = true
			movingParticle.needToFall = false
		}

		// if x < 0 || y < 0 || x >= Width || y >= Height {
		// 	return
		// }

		y = int(math.Min(float64(y), float64(Height-1)))
		x = int(math.Min(float64(x), float64(Width-1)))

		y = int(math.Max(float64(y), 0))
		x = int(math.Max(float64(x), 0))

		if !movingParticle.Fallen {
			initialVolatility := movingParticle.initialVolatility
			rotating := movingParticle.rotating
			movingParticle = &(*matrixParticles)[y][x]
			movingParticle.X = x
			movingParticle.Y = y
			movingParticle.Fallen = false
			movingParticle.needToFall = true
			movingParticle.initialVolatility = initialVolatility
			movingParticle.rotating = rotating
		}

	}
}

func moveParticleWithOutRotate(particle *Particle, matrixParticles matrixParticlesPtrType) {

	Height := len((*matrixParticles))
	Width := len((*matrixParticles)[Height-1])

	y := particle.Y
	x := particle.X

	var currCellY int = int(math.Round(float64(y)))
	// var currCellX int = int(math.Round(float64(x)))

	var downCellY int = int(math.Round(float64(y) + 1.0*math.Sin(realRotateAngle()*math.Pi/180)))
	var downCellX int = int(math.Round(float64(x) + 1.0*math.Cos(realRotateAngle()*math.Pi/180)))

	var leftCellX int = int(math.Round(float64(x) - 1.0*math.Cos(realRotateAngle()*math.Pi/180)))
	var rightCellX int = int(math.Round(float64(x) + 1.0*math.Cos(realRotateAngle()*math.Pi/180)))

	// currCellY = int(math.Min(float64(currCellY), Height-1))
	// downCellY = int(math.Min(float64(downCellY), Height-1))

	// currCellX = int(math.Min(float64(currCellX), Width-1))
	// leftCellX = int(math.Min(float64(leftCellX), Width-1))
	// rightCellX = int(math.Min(float64(rightCellX), Width-1))

	// currCellY = int(math.Max(float64(currCellY), 0))
	// downCellY = int(math.Max(float64(downCellY), 0))

	// currCellX = int(math.Max(float64(currCellX), 0))
	// leftCellX = int(math.Max(float64(leftCellX), 0))
	// rightCellX = int(math.Max(float64(rightCellX), 0))

	// var downCellY int = y + 1
	// var leftCellX int = x - 1
	// var rightCellX int = x + 1

	movingParticle := &(*matrixParticles)[y][x]

	// bottomCellFill := y == Height-1 || matrixParticles[downCellY][x].Fallen
	// bottomRightCellFill := y == Height-1 || matrixParticles[downCellY][int(math.Min(float64(rightCellX), float64(Width-1)))].Fallen || matrixParticles[downCellY][int(math.Min(float64(rightCellX), float64(Width-1)))].needToFall
	// bottomLeftCellFill := y == Height-1 || matrixParticles[downCellY][int(math.Max(float64(leftCellX), float64(0)))].Fallen || matrixParticles[downCellY][int(math.Max(float64(leftCellX), float64(0)))].needToFall
	// rightCellFill := matrixParticles[y][int(math.Min(float64(rightCellX), float64(Width-1)))].Fallen || matrixParticles[y][int(math.Min(float64(rightCellX), float64(Width-1)))].needToFall
	// leftCellFill := matrixParticles[y][int(math.Max(float64(leftCellX), float64(0)))].Fallen || matrixParticles[y][int(math.Max(float64(leftCellX), float64(0)))].needToFall

	bottomCellFill := (*matrixParticles)[downCellY][downCellX].Fallen
	bottomRightCellFill := (*matrixParticles)[downCellY][rightCellX].Fallen || (*matrixParticles)[downCellY][rightCellX].needToFall
	bottomLeftCellFill := (*matrixParticles)[downCellY][leftCellX].Fallen || (*matrixParticles)[downCellY][leftCellX].needToFall
	rightCellFill := (*matrixParticles)[currCellY][rightCellX].Fallen || (*matrixParticles)[currCellY][rightCellX].needToFall
	leftCellFill := (*matrixParticles)[currCellY][leftCellX].Fallen || (*matrixParticles)[currCellY][leftCellX].needToFall

	// if y == Height-1 ||
	// 	bottomCellFill && bottomRightCellFill && bottomLeftCellFill ||
	// 	bottomCellFill && rightCellFill && leftCellFill {

	// 	movingParticle.Fallen = true
	// 	movingParticle.needToFall = false
	// } else {

	if bottomCellFill && bottomRightCellFill && bottomLeftCellFill ||
		bottomCellFill && rightCellFill && leftCellFill {

		movingParticle.Fallen = true
		movingParticle.needToFall = false

	} else {

		movingParticle.Fallen = false
		movingParticle.needToFall = false

		// // if bottomCellFill && !bottomLeftCellFill && !leftCellFill {
		// // 	x = int(math.Round(float64(x) - 1.0*math.Cos(realRotateAngle()*math.Pi/180)))
		// // 	y = int(math.Round(float64(y) + 1.0*math.Sin(realRotateAngle()*math.Pi/180)))
		// // 	// x = x - 1
		// // 	// y = y + 1
		// // } else if bottomCellFill && !bottomRightCellFill && !rightCellFill {
		// // 	x = int(math.Round(float64(x) + 1.0*math.Cos(realRotateAngle()*math.Pi/180)))
		// // 	y = int(math.Round(float64(y) + 1.0*math.Sin(realRotateAngle()*math.Pi/180)))
		// // 	// x = x + 1
		// // 	// y = y + 1
		// // } else if !bottomCellFill {
		// x = int(math.Round(float64(x) + 1.0*math.Cos(realRotateAngle()*math.Pi/180)))
		// y = int(math.Round(float64(y) + 1.0*math.Sin(realRotateAngle()*math.Pi/180)))
		// // y = y + 1
		// // } else if bottomCellFill {
		// // 	movingParticle.Fallen = true
		// // 	movingParticle.needToFall = false
		// // }

		// if bottomCellFill && !bottomLeftCellFill && !leftCellFill {
		// 	x = int(math.Round(float64(x) - 1.0*math.Cos(realRotateAngle()*math.Pi/180)))
		// 	y = int(math.Round(float64(y) + 1.0*math.Sin(realRotateAngle()*math.Pi/180)))
		// 	// x = x - 1
		// 	// y = y + 1
		// } else if bottomCellFill && !bottomRightCellFill && !rightCellFill {
		// 	x = int(math.Round(float64(x) + 1.0*math.Cos(realRotateAngle()*math.Pi/180)))
		// 	y = int(math.Round(float64(y) + 1.0*math.Sin(realRotateAngle()*math.Pi/180)))
		// 	// x = x + 1
		// 	// y = y + 1
		// } else

		if !bottomCellFill {
			x = int(math.Round(float64(x) + 1.0*math.Cos(realRotateAngle()*math.Pi/180)))
			y = int(math.Round(float64(y) + 1.0*math.Sin(realRotateAngle()*math.Pi/180)))
			// y = y + 1
		} else if bottomCellFill {
			movingParticle.Fallen = true
			movingParticle.needToFall = false
		}

		// if x < 0 || y < 0 || x >= Width || y >= Height {
		// 	return
		// }

		y = int(math.Min(float64(y), float64(Height-1)))
		x = int(math.Min(float64(x), float64(Width-1)))

		y = int(math.Max(float64(y), 0))
		x = int(math.Max(float64(x), 0))

		if !movingParticle.Fallen {
			initialVolatility := movingParticle.initialVolatility
			rotating := movingParticle.rotating
			movingParticle = &(*matrixParticles)[y][x]
			movingParticle.X = x
			movingParticle.Y = y
			movingParticle.Fallen = false
			movingParticle.needToFall = true
			movingParticle.initialVolatility = initialVolatility
			movingParticle.rotating = rotating
		}

	}
}

func convertMatixToImg(matrixParticles matrixParticlesPtrType) *image.RGBA {

	Height := len(*matrixParticles)
	Width := len((*matrixParticles)[Height-1])
	pixelWidth := 3

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

	img = rotateImage(img, imageRotateAngle)
	imageRotateAngle = imageRotateAngle + 1

	return img
}

func rotateImage(src *image.RGBA, angle float64) *image.RGBA {
	// Создаем новое изображение с теми же размерами
	dst := image.NewRGBA(src.Bounds())

	// Поворачиваем изображение
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)
	rotated := imaging.Rotate(dst, angle, color.Transparent)

	//return rotated.(*image.RGBA)
	return convertNRGBAToRGBA(rotated)

}

func convertNRGBAToRGBA(src *image.NRGBA) *image.RGBA {
	bounds := src.Bounds()
	// width := bounds.Dx()
	// height := bounds.Dy()

	// Создаем новое изображение с теми же размерами
	dst := image.NewRGBA(bounds)

	// Копируем каждый пиксель из исходного изображения в новое изображение
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)

	return dst
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
