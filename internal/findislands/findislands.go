package findislands

import (
	"image"

	"golang.org/x/exp/constraints"
)

// identifies pixel islands in an image
func ImageToIslands(img image.Image, diagonal bool) []image.Rectangle {
	rects := make([]image.Rectangle, 0)
	visited := newVisitedArray(img.Bounds())
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			if !visited.get(x, y) && isVisiblePixel(img, x, y) {
				r := findConnectedPixels(img, x, y, diagonal, visited)
				rects = append(rects, r)
			}
		}
	}
	return rects
}

// identifies pixel islands in images
func ImagesToIslands(images []image.Image, diagonal bool) [][]image.Rectangle {
	boxes := make([][]image.Rectangle, 0, len(images)) // pre-allocate capacity for efficiency
	for _, img := range images {
		// For each image, call the functionally complete ImageToIslands
		islands := ImageToIslands(img, diagonal)
		boxes = append(boxes, islands)
	}
	return boxes
}

// Given an image and a starting pixel, finds all connected pixels and returns a square encompassing them.
// 'visited' is used to track progress.
// the diagonal flag enables checking diagonally connected pixels.
func findConnectedPixels(img image.Image, x, y int, diagonal bool, visited visitedArray) image.Rectangle {
	bounds := img.Bounds()
	stack := []image.Point{{X: x, Y: y}}

	var minX, minY, maxX, maxY int
	minX = x
	maxX = x
	minY = y
	maxY = y

	for len(stack) > 0 {
		point := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		visited.set(point.X, point.Y, true)
		minX = min(minX, point.X)
		minY = min(minY, point.Y)
		maxX = max(maxX, point.X)
		maxY = max(maxY, point.Y)

		// adds pt to the list of pts to check IF it's within bounds
		// and hasn't been visited yet and is visible.
		pointCheck := func(pt image.Point) {
			if !pt.In(bounds) {
				return
			}
			if visited.get(pt.X, pt.Y) {
				return
			}
			if isVisiblePixel(img, pt.X, pt.Y) {
				stack = append(stack, pt)
			}
		}

		for xOff := -1; xOff <= 1; xOff++ {
			for yOff := -1; yOff <= 1; yOff++ {
				if abs(xOff)+abs(yOff) == 2 && !diagonal {
					continue
				}
				pointCheck(image.Point{point.X + xOff, point.Y + yOff})
			}
		}
	}
	return image.Rect(minX, minY, maxX, maxY)
}

func isVisiblePixel(img image.Image, x, y int) bool {
	_, _, _, a := img.At(x, y).RGBA()
	return a > 0
}

/* Originally used an Image.Grey to track visited pixels, however the interface involves too much indirection */
type visitedArray struct {
	data []bool
	w, h int
}

func newVisitedArray(bounds image.Rectangle) visitedArray {
	var va visitedArray
	va.w = bounds.Dx()
	va.h = bounds.Dy()
	va.data = make([]bool, va.w*va.h)
	return va
}

// nb: no parameter boundary validation
func (va *visitedArray) get(x, y int) bool {
	return va.data[(y*va.w)+x]
}

// nb: no parameter boundary validation
func (va *visitedArray) set(x, y int, v bool) {
	va.data[(y*va.w)+x] = v
}

func abs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
