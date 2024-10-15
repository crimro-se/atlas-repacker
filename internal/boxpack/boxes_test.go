package boxpack

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"
)

func TestPack(t *testing.T) {
	var boxes []Box
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 20, 10)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 10, 20)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 30, 30)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 20, 10)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 10, 10)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 10, 20)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 30, 30)})
	boxes = append(boxes, Box{sourceRect: image.Rect(0, 0, 50, 50)})
	//fmt.Println(PackBoxes(boxes, 100, 40, 1, 0))
	ab, _ := PackAllBoxes(boxes, 100, 40, 1, 0)

	img := DrawRects(ab[0], 100, 40)

	// test ability to find box in img
	//visited := make(map[image.Point]bool)
	visited := NewVisitedArray(img.Bounds())
	fmt.Println(findConnectedPixels(img, 0, 0, false, visited))
	for i := 0; i < 100; i++ {
		//fmt.Println(visited.At(i, 1))
	}

	file, err := os.Create("test.png")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	png.Encode(file, img)
}

func DrawRects(boxes []Box, W, H int) image.Image {
	img := image.NewRGBA64(image.Rect(0, 0, W, H))
	src := image.NewUniform(color.RGBA{0, 100, 255, 255})
	for _, b := range boxes {
		draw.Draw(img, b.destRect, src, image.ZP, draw.Src)
	}
	return img
}
