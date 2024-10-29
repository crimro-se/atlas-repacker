package boxpack

import (
	"image"
	"testing"
)

func TestPacking(t *testing.T) {
	var boxes []BoxTranslation
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 20, 10)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 10, 20)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 30, 30)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 20, 10)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 10, 20)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 30, 30)})
	boxes = append(boxes, BoxTranslation{sourceRect: image.Rect(0, 0, 50, 50)})
	unpacked := PackBoxes(boxes, 900, 40, 1, 0)
	if unpacked != 1 {
		t.Fail()
	}
	unpacked = PackBoxes(boxes, 900, 51, 1, 0)
	if unpacked != 0 {
		t.Fail()
	}

	// verify exactly one box occupies 0,0
	c := 0
	for _, box := range boxes {
		if box.destRect.Min.X == 0 && box.destRect.Min.Y == 0 {
			c++
		}
	}
	if c != 1 {
		t.Fail()
	}

	// verify src boxes still all at 0,0
	c = 0
	for _, box := range boxes {
		if box.sourceRect.Min.X == 0 && box.sourceRect.Min.Y == 0 {
			c++
		}
	}
	if c != len(boxes) {
		t.Fail()
	}
}
