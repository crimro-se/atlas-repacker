package main

import (
	"image"

	"github.com/crimro-se/atlas-repacker/internal/boxpack"
)

func rectsToBoxTranslation(rr [][]image.Rectangle) []boxpack.BoxTranslation {
	boxes := make([]boxpack.BoxTranslation, 0, len(rr[0]))
	for i, rects := range rr {
		for _, rect := range rects {
			boxes = append(boxes, boxpack.BoxFromRect(i, rect, false))
		}
	}
	return boxes
}
