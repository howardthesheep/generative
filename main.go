package main

import (
	"fmt"
	"image"
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
	file, err := os.OpenFile("generate.svg", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}

	start := time.Now()
	canvas := svg.New(file)
	canvas.Start(width, height)
	drawBackground(canvas, "white")
	//drawRandomSquares(Range{50,500},canvas)
	//drawRandomRecursiveSquares(Range{200, 300}, canvas)
	drawSandScript(canvas)

	// Timing
	str := fmt.Sprintf("Time to generate: %s", time.Since(start))
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

// drawSandScript draws an array of randomly generated 'script-like' characters
// using bezier curves
func drawSandScript(canvas *svg.SVG) {
	type cellBoundary [2]Range                         // Range for Valid X values, Range for Valid Y values //TODO: Is it bad to have type def in a func vs global? Will test to see diff
	const gridFrac = 1.0 / 5                           // How many cells per col/row
	var gridRanges []cellBoundary                      // Reference to cell boundaries
	rng := rand.New(rand.NewSource(time.Now().Unix())) // Create random number generator w/ seed

	// Divide the canvas into grid
	widthStep := width * gridFrac   // cell width
	heightStep := height * gridFrac // cell height

	// Store grid cells using 2 Range's
	for x := float64(0); x < width; x += widthStep {
		for y := float64(0); y < height; y += heightStep {
			xRange := Range{
				Min: x,
				Max: x + widthStep,
			}
			yRange := Range{
				Min: y,
				Max: y + heightStep,
			}
			boundary := cellBoundary{xRange, yRange}
			gridRanges = append(gridRanges, boundary)
		}
	}

	// Generate 3 random points within each grid box
	var vertsPerChar = 3 // Vertices per 'character'
	var characterPoints [][]image.Point

	for _, gridRange := range gridRanges {
		var vertices []image.Point // Stores vertices for single character
		xRange := gridRange[0]
		yRange := gridRange[1]

		for i := 0; i < vertsPerChar; i++ {
			randX := int(xRange.Min.(float64)) + rng.Intn(int(xRange.Max.(float64))-int(xRange.Min.(float64)))
			randY := int(yRange.Min.(float64)) + rng.Intn(int(yRange.Max.(float64))-int(yRange.Min.(float64)))

			//fmt.Printf("Cell #%d Vertex Pair #%d (%d,%d)\n",z,i+1,randX,randY)

			vertex := image.Point{
				X: randX,
				Y: randY,
			}
			vertices = append(vertices, vertex)
		}

		characterPoints = append(characterPoints, vertices)
	}

	// Create a Bezier Curve connecting the 3 points
	for z, vertexSet := range characterPoints {
		controlA := image.Point{} // (cx,cy)
		controlB := image.Point{} // (px,py)

		// Generate random points within grid square for controlA & controlB
		xRange := gridRanges[z][0]
		yRange := gridRanges[z][1]
		controlA.X = int(xRange.Min.(float64)) + rng.Intn(int(xRange.Max.(float64))-int(xRange.Min.(float64)))
		controlA.Y = int(yRange.Min.(float64)) + rng.Intn(int(yRange.Max.(float64))-int(yRange.Min.(float64)))
		controlB.X = int(xRange.Min.(float64)) + rng.Intn(int(xRange.Max.(float64))-int(xRange.Min.(float64)))
		controlB.Y = int(yRange.Min.(float64)) + rng.Intn(int(yRange.Max.(float64))-int(yRange.Min.(float64)))

		points := len(vertexSet) // How many points per character
		for i := 0; i < points; i++ {

			j := i + 1
			if j == points { // Handle if were at the end of set
				j = 0
			}

			// Draw Bezier between 2 points
			canvas.Bezier(vertexSet[i].X, vertexSet[i].Y,
				controlA.X, controlA.Y,
				controlB.X, controlB.Y,
				vertexSet[j].X, vertexSet[j].Y)
		}
	}
}

func tileSquares(canvas *svg.SVG) {
	scale := 1.0 / 8.0
	length := height * scale

	canvas.Square(0, 0, int(math.Round(length)))
}
