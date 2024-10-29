package findislands

import (
	"bytes"
	"encoding/base64"
	"image"
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
	boxes := ImageToIslands(img, false)
	if len(boxes) != 5 {
		t.Fail()
	}
}
