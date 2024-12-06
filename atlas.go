package main

import (
	"fmt"
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

func parseAtlasFile(filename string, imgRef int) ([]boxpack.BoxTranslation, error) {
	boxes := make([]boxpack.BoxTranslation, 0)
	fp, err := os.Open(filename)
	if err != nil {
		return boxes, fmt.Errorf("error whilst trying to open (%s): %w", filename, err)
	}
	defer fp.Close()
	atlasRegions, err := atlas.ParseAtlasFile(fp)
	if err != nil {
		return boxes, fmt.Errorf("error whilst trying to parse (%s): %w", filename, err)
	}
	return atlasToBoxes(imgRef, atlasRegions), nil
}

// converts atlasRegions type to []boxpack.BoxTranslation
func atlasToBoxes(refImage int, ar atlas.AtlasRegions) []boxpack.BoxTranslation {
	boxes := make([]boxpack.BoxTranslation, 0, len(ar))
	for _, v := range ar {
		boxes = append(boxes, boxpack.BoxFromRect(refImage, v.Rectangle, v.RotateRequired))
	}
	return boxes
}
