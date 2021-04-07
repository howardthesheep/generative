package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	svg "github.com/ajstarks/svgo"
)

// Range is a helper struct for passing
// Min/Max value pairs
//
//		Note: Min and Max should only be number types (int,float)
type Range struct {
	Min interface{}
	Max interface{}
}

const (
	width  = 1920 // Canvas Width
	height = 1080 // Canvas Height
)

func main() {
	file, err := os.OpenFile("generate.svg", os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}

	canvas := svg.New(file)
	canvas.Start(width, height)
	drawBackground(canvas, "white")
	//drawRandomSquares(Range{50,500},canvas)
	start := time.Now()
	drawRandomRecursiveSquares(Range{200, 300}, canvas)
	end := time.Now()

	// Timing
	elapsed := end.Sub(start)
	str := fmt.Sprintf("Time to generate: %s", elapsed)
	println(str)

	canvas.End()
}

// drawBackground sets the SVG background fill to whatever color
// you specify
func drawBackground(canvas *svg.SVG, color string) {
	canvas.Rect(0, 0, width, height, "fill:"+color+";")
}

// drawRandomSquares draws a random amount of squares within the pre-defined Range
func drawRandomSquares(squareCount Range, canvas *svg.SVG) {
	// Create random number generator w/ seed
	rng := rand.New(rand.NewSource(time.Now().Unix()))

	maxSqrs := squareCount.Max.(int)
	minSqrs := squareCount.Min.(int)
	squareCnt := minSqrs + rng.Intn(maxSqrs-minSqrs)

	// Create our squares
	for i := 0; i < squareCnt; i++ {
		// Create random attribute values
		randX := rng.Intn(width)
		randY := rng.Intn(height)
		length := rng.Intn(width / 4)
		r := rng.Intn(255)
		g := rng.Intn(255)
		b := rng.Intn(255)
		opacity := rng.Float32()

		// Draw the squares w/ attributes
		canvas.Square(randX, randY, length, "fill:rgb("+strconv.Itoa(r)+
			","+strconv.Itoa(g)+
			","+strconv.Itoa(b)+
			");opacity:"+strconv.FormatFloat(float64(opacity), 'g', 2, 64)+";")
	}
}

// drawRandomRecursiveSquares draws a random amount of squares within the given Range
// then recursively draws concentric smaller squares of random color
// 		Inspired by: https://www.generativeart.com/on/cic/GA2010/2010_6.pdf
func drawRandomRecursiveSquares(squareCount Range, canvas *svg.SVG) {
	// Create random number generator w/ seed
	rng := rand.New(rand.NewSource(time.Now().Unix()))

	maxSqrs := squareCount.Max.(int)
	minSqrs := squareCount.Min.(int)
	squareCnt := minSqrs + rng.Intn(maxSqrs-minSqrs)

	// Create our squares
	for i := 0; i < squareCnt; i++ {
		// Create random attribute values
		randX := rng.Intn(width)
		randY := rng.Intn(height)
		length := rng.Intn(width / 4)
		r := rng.Intn(255)
		g := rng.Intn(255)
		b := rng.Intn(255)
		opacity := rng.Float32()

		// Cap our opacity
		if opacity > 0.75 {
			opacity = 0.75
		}
		formattedOpac := strconv.FormatFloat(float64(opacity), 'g', 2, 64)

		// Draw the squares w/ attributes
		canvas.Square(randX, randY, length, "fill:rgb("+strconv.Itoa(r)+
			","+strconv.Itoa(g)+
			","+strconv.Itoa(b)+
			");opacity:"+formattedOpac+";")

		// Draw concentric Squares
		randInset := rand.Intn(50) + 5
		for j := length - randInset; j > 0; j -= randInset {
			randInset = rand.Intn(50) + 5
			r2 := int(float32(r)*0.6 + 10)
			g2 := int(float32(g)*0.6 + 10)
			b2 := int(float32(b)*0.6 + 10)
			canvas.Square(randX+randInset, randY+randInset, j, "fill:rgb("+strconv.Itoa(r2)+
				","+strconv.Itoa(g2)+
				","+strconv.Itoa(b2)+");opacity:"+formattedOpac+";")
		}
	}
}

func tileSquares(canvas *svg.SVG) {
	scale := 1.0 / 8.0
	length := height * scale

	canvas.Square(0, 0, int(math.Round(length)))
}
