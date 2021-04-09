package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"math/rand"
	"os"
	"sync"
)

const (
	_ViewWidth  = 1920
	_ViewHeight = 1080
	_MaxEscape  = 64
)

func main() {
	center := 1.5 + 0i + 0.1 + 0.1i*complex(float64(9), 2*float64(9))
	b := generate(_ViewWidth, _ViewHeight, center, 0.5)
	f, err := os.Create("mandelbrot.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = png.Encode(f, b); err != nil {
		fmt.Println(err)
	}
	if err = f.Close(); err != nil {
		fmt.Println(err)
	}
}

func generate(imgWidth int, imgHeight int, viewCenter complex128, radius float64) image.Image {
	m := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	zoomWidth := radius * 2
	pixelWidth := zoomWidth / float64(imgWidth)
	pixelHeight := pixelWidth
	viewHeight := (float64(imgHeight) / float64(imgWidth)) * zoomWidth
	left := (real(viewCenter) - (zoomWidth / 2)) + pixelWidth/2
	top := (imag(viewCenter) - (viewHeight / 2)) + pixelHeight/2
	escapeColor := color.RGBA{R: 0, G: 0, B: 0, A: 0}

	// Create Color Palette
	palette := make([]color.RGBA, _MaxEscape)
	for i := 0; i < _MaxEscape-1; i++ {
		palette[i] = color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255}
	}

	var wgx sync.WaitGroup
	wgx.Add(imgWidth)
	for x := 0; x < imgWidth; x++ {
		go func(xx int) {
			defer wgx.Done()
			for y := 0; y < imgHeight; y++ {
				coord := complex(left+float64(xx)*pixelWidth, top+float64(y)*pixelHeight)
				f := escape(coord)
				if f == _MaxEscape-1 {
					m.Set(xx, y, escapeColor)
				}
				m.Set(xx, y, palette[f])
			}
		}(x)
	}
	wgx.Wait()
	return m
}

func escape(c complex128) int {
	z := c
	for i := 0; i < _MaxEscape-1; i++ {
		if cmplx.Abs(z) > 2 {
			return i
		}
		z = cmplx.Pow(z, 2) + c
	}
	return _MaxEscape - 1
}
