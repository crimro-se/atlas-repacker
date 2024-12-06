package main

import (
	"image"

	"github.com/crimro-se/atlas-repacker/internal/boxpack"
	"github.com/crimro-se/atlas-repacker/internal/findislands"
)

func detectIslands(img image.Image, imgRef int, diagonalDetection bool) ([]boxpack.BoxTranslation, error) {
	rects := findislands.ImageToIslands(img, diagonalDetection)
	boxes := make([]boxpack.BoxTranslation, 0, len(rects))
	for _, rect := range rects {
		boxes = append(boxes, boxpack.BoxFromRect(imgRef, rect, false))
	}
	return boxes, nil
}

func rectsToBoxTranslation(rr [][]image.Rectangle) []boxpack.BoxTranslation {
	boxes := make([]boxpack.BoxTranslation, 0, len(rr[0]))
	for i, rects := range rr {
		for _, rect := range rects {
			boxes = append(boxes, boxpack.BoxFromRect(i, rect, false))
		}
	}
	return boxes
}
