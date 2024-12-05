package main

import (
	"os"

	"github.com/crimro-se/atlas-repacker/internal/atlas"
	"github.com/crimro-se/atlas-repacker/internal/boxpack"
)

func loadAllAtlas(files []string) ([]boxpack.BoxTranslation, error) {
	boxes := make([]boxpack.BoxTranslation, 0, max(len(files), 10))
	for i, inputFile := range files {
		fp, err := os.Open(inputFile)
		if err != nil {
			return boxes, err
		}
		defer fp.Close()

		atlas, err := atlas.ParseAtlasFile(fp)
		if err != nil {
			return boxes, err
		}
		boxes = append(boxes, atlasToBoxes(i, atlas)...)
	}
	return boxes, nil
}

// converts atlasRegions type to []boxpack.BoxTranslation
func atlasToBoxes(refImage int, ar atlas.AtlasRegions) []boxpack.BoxTranslation {
	boxes := make([]boxpack.BoxTranslation, 0, len(ar))
	for _, v := range ar {
		boxes = append(boxes, boxpack.BoxFromRect(refImage, v.Rectangle, v.RotateRequired))
	}
	return boxes
}
