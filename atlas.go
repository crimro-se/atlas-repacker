package main

import (
	"fmt"
	"os"

	"github.com/crimro-se/atlas-repacker/internal/atlas"
	"github.com/crimro-se/atlas-repacker/internal/boxpack"
)

func parseAtlasFile(filename string, imgRef int) ([]NamedBox, error) {
	boxes := make([]NamedBox, 0)
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
func atlasToBoxes(refImage int, ar atlas.AtlasRegions) []NamedBox {
	boxes := make([]NamedBox, 0, len(ar))
	for name, v := range ar {
		boxes = append(boxes, NamedBoxFromBoxpack(boxpack.BoxFromRect(refImage, v.Rectangle, v.RotateRequired), name))
	}
	return boxes
}
