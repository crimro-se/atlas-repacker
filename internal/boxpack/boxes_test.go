package boxpack

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"testing"
)

var testimg_b64 string = `
iVBORw0KGgoAAAANSUhEUgAAAQAAAACAAgMAAACZ21+ZAAABhGlDQ1BJQ0MgcHJvZmlsZQAAKJF9
kT1Iw0AcxV9TpVIqDnbQ4pChOtlFRQSXWoUiVAi1QqsOJpd+QZOGpMXFUXAtOPixWHVwcdbVwVUQ
BD9AnB2cFF2kxP8lhRYxHhz34929x907QGhWmGb1xAFNr5npZELM5lbFwCuCiEDALIZlZhlzkpSC
5/i6h4+vdzGe5X3uz9Gv5i0G+ETiODPMGvEG8fRmzeC8TxxmJVklPiceN+mCxI9cV1x+41x0WOCZ
YTOTnicOE4vFLla6mJVMjXiKOKpqOuULWZdVzluctUqdte/JXxjK6yvLXKc5giQWsQQJIhTUUUYF
NcRo1UmxkKb9hIc/4vglcinkKoORYwFVaJAdP/gf/O7WKkxOuEmhBND7Ytsfo0BgF2g1bPv72LZb
J4D/GbjSO/5qE5j5JL3R0aJHwMA2cHHd0ZQ94HIHGHoyZFN2JD9NoVAA3s/om3LA4C0QXHN7a+/j
9AHIUFepG+DgEBgrUva6x7v7unv790y7vx/Op3LLMNEO0gAAAAlwSFlzAAAuIwAALiMBeKU/dgAA
AAd0SU1FB+gKGgMNCY3j5KUAAAAZdEVYdENvbW1lbnQAQ3JlYXRlZCB3aXRoIEdJTVBXgQ4XAAAA
CVBMVEUAAAAAAAD///+D3c/SAAAAAXRSTlMAQObYZgAAAAFiS0dEAmYLfGQAAABfSURBVGiB7c3B
CcAgEEVBm7AvD27/rQQ8hRBiDhLUzDsun9mUpEHlaFUAAPAIRPQeALYAOv0NiEsVAFgJyG1QABMD
LwIAAB8C52uZD0i3MwAAAAAAAAAAoB9gACBp+w67yiWwcnXERQAAAABJRU5ErkJggg==`

func loadTestImg() (image.Image, error) {
	pngBytes, err := base64.StdEncoding.DecodeString(testimg_b64)
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(bytes.NewReader(pngBytes))
	return img, err
}

func TestIslands(t *testing.T) {
	img, err := loadTestImg()
	if err != nil {
		t.Error(err)
	}
	boxes := ImageToBoxes(img, false)
	if len(boxes) != 5 {
		t.Fail()
	}
}

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

func drawRects(boxes []BoxTranslation, W, H int) image.Image {
	img := image.NewRGBA64(image.Rect(0, 0, W, H))
	src := image.NewUniform(color.RGBA{0, 100, 255, 255})
	for _, b := range boxes {
		draw.Draw(img, b.destRect, src, image.ZP, draw.Src)
	}
	return img
}
