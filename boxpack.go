package main

import "github.com/crimro-se/atlas-repacker/internal/boxpack"

// from boxpack's struct to ours, which also adds a name
func NamedBoxFromBoxpack(b boxpack.BoxTranslation, name string) NamedBox {
	var box NamedBox
	box.BoxTranslation = b
	box.Name = name
	return box
}

// promotes a slice of boxpack.BoxTranslation to []NamedBox. names are optional and can be nil
func NamedBoxFromBoxpackSlice(boxes []boxpack.BoxTranslation, names []string) []NamedBox {
	namedBoxes := make([]NamedBox, 0)
	if names != nil && len(names) != len(boxes) {
		names = nil
	}
	var nb NamedBox
	for i, v := range boxes {
		nb.BoxTranslation = v
		if names != nil {
			nb.Name = names[i]
		}
		namedBoxes = append(namedBoxes, nb)
	}
	return namedBoxes
}

// converts a slice of namedbox to a slice of boxpack.boxtranslation
func BoxpackSliceFromNamedBoxes(boxes []NamedBox) []boxpack.BoxTranslation {
	boxTR := make([]boxpack.BoxTranslation, 0)
	var box boxpack.BoxTranslation
	for _, v := range boxes {
		box = v.BoxTranslation
		boxTR = append(boxTR, box)
	}
	return boxTR
}

// invoke boxpack.PackBoxes whilst adapting []NamedBox to []boxpack.BoxTranslation
func PackNamedBoxes(boxes []NamedBox, W, H, boxMargin, offset int) int {
	boxTR := BoxpackSliceFromNamedBoxes(boxes)
	unpacked := boxpack.PackBoxes(boxTR, W, H, boxMargin, offset)
	// apply results
	for i, _ := range boxes {
		boxes[i].BoxTranslation = boxTR[i]
	}
	return unpacked
}

// adaptor for boxpack.EstimateOutputWH
func EstimateOutputWH(boxes []NamedBox, margin int) int {
	boxTR := BoxpackSliceFromNamedBoxes(boxes)
	return boxpack.EstimateOutputWH(boxTR, margin)
}
