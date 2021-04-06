package main

import (
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
	width  = 500  // Canvas Width
	height = 1000 // Canvas Height
)

func main() {
	file, err := os.OpenFile("generate.svg", os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}

	canvas := svg.New(file)
	canvas.Start(width, height)
	//canvas.Circle(width/2, height/2, 100, "fill:blue;")
	//canvas.Text(width/2, height/2, "Hello, SVG", "text-anchor:middle;font-size:30px;fill:white")
	drawBackground(canvas, "white")
	//drawRandomSquares(Range{50,500},canvas)
	canvas.End()
}

// drawBackground sets the SVG background fill to whatever color
// you specify
func drawBackground(canvas *svg.SVG, color string) {
	canvas.Rect(0,0, width, height, "fill:"+color+";")
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
		length := rng.Intn(width/4)
		r := rng.Intn(255)
		g := rng.Intn(255)
		b := rng.Intn(255)
		opacity := rng.Float32()

		// Draw the squares w/ attributes
		canvas.Square(randX, randY, length, "fill:rgb("+strconv.Itoa(r)+
			"," + strconv.Itoa(g) +
			"," + strconv.Itoa(b) +
			");opacity:"+ strconv.FormatFloat(float64(opacity),'g',2,64)+";")
	}
}

func tileSquares(canvas *svg.SVG) {
	scale := 1.0/8.0
	length := height * scale

	canvas.Square(0,0, int(math.Round(length)))
}