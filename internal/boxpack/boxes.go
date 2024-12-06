// package to make stb_rect_pack.h usable in Go.
package boxpack

/*
	#define STB_RECT_PACK_IMPLEMENTATION

	#include "stb_rect_pack.h"
	#include <stdlib.h>
	#include <stdio.h>

struct stbrp_rect* allocateRects(int n) {
   struct stbrp_rect* array;

   // Allocate memory for the array of n structs
   array = (struct stbrp_rect*)calloc(n, sizeof(struct stbrp_rect));

   if (array == NULL) {
      fprintf(stderr, "Memory allocation for rects failed\n");
      exit(EXIT_FAILURE);
   }

   return array;
}

struct stbrp_node* allocateNodes(int n) {
   struct stbrp_node* array;

   // Allocate memory for the array of n structs
   array = (struct stbrp_node*)calloc(n, sizeof(struct stbrp_node));

   if (array == NULL) {
      fprintf(stderr, "Memory allocation for nodes failed\n");
      exit(EXIT_FAILURE);
   }

   return array;
}

struct stbrp_context* allocateCTX() {
   struct stbrp_context* array;

   // Allocate memory for the array of n structs
   array = (struct stbrp_context*)calloc(1, sizeof(struct stbrp_context));

   if (array == NULL) {
      fprintf(stderr, "Memory allocation for ctx failed\n");
      exit(EXIT_FAILURE);
   }

   return array;
}

// no bounds checking btw
void assignValue(struct stbrp_rect* array, int index, struct stbrp_rect* value) {
   array[index] = *value;
}

void getValue(struct stbrp_rect* array, int index, struct stbrp_rect* value) {
   *value= array[index];
}

void myFree(void *mem){
	free(mem);
}
*/
import "C"
import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"unsafe"

	"github.com/disintegration/imaging"
)

// A box translation tracks its source image number and rect,
// its new location and if it was successfully repacked.
type BoxTranslation struct {
	imgSrc         int             // which input image is this box from?
	sourceRect     image.Rectangle // pixel locations on original input image
	destRect       image.Rectangle // destination rect.
	wasPacked      bool            // true if this box has been successfully packed
	deferredRotate bool            // rotate 90 clockwise when rendering if true
}

// returns the sum of area required for all sourceRect boxes
func getSourceArea(boxes []BoxTranslation, margin int) int {
	area := 0
	for _, box := range boxes {
		area += (box.sourceRect.Dx() + margin) * (box.sourceRect.Dy() + margin)
	}
	return area
}

func BoxFromRect(imgref int, r image.Rectangle, rotate bool) BoxTranslation {
	return BoxTranslation{imgSrc: imgref, sourceRect: r, wasPacked: false, deferredRotate: rotate}
}

// Estimates an appropriate w & h for output based on the input squares
func EstimateOutputWH(boxes []BoxTranslation, margin int) int {
	maxWH := 0
	for _, box := range boxes {
		maxWH = max(maxWH, box.sourceRect.Dx()+margin, box.sourceRect.Dy()+margin)
	}
	areaSqrt := int(math.Sqrt(float64(getSourceArea(boxes, margin))))
	return max(maxWH, areaSqrt)
}

// Packs boxes, with multiple output sheets. Input slice isn't modified.
// boxMargin - additinal padding to provide each box in total, pixels.
// offset - ammount to offset each box. useful values are half of margin, =margin, or zero.
// returns result and count of any remaining unpacked (size+margin > W or H)
func PackAllBoxes(boxesImmutable []BoxTranslation, W, H, boxMargin, offset int) ([][]BoxTranslation, int) {
	allBoxes := make([][]BoxTranslation, 0)
	boxes := make([]BoxTranslation, len(boxesImmutable)) // working copy
	copy(boxes, boxesImmutable)
	i := 0
	unpacked := len(boxes)
	for unpacked > 0 {
		previous := unpacked
		unpacked = PackBoxes(boxes, W, H, boxMargin, offset)
		if previous == unpacked {
			// no progress since previous iteration
			break
		}
		remainer := make([]BoxTranslation, 0, unpacked)
		allBoxes = append(allBoxes, make([]BoxTranslation, 0, previous-unpacked))
		for _, box := range boxes {
			if box.wasPacked {
				allBoxes[i] = append(allBoxes[i], box)
			} else {
				remainer = append(remainer, box)
			}
		}
		boxes = remainer
		i++
	}
	return allBoxes, unpacked
}

// Packs the boxes parameter in-place, updating destRect and wasPacked accordingly
// boxMargin - additinal padding to provide each box in total, pixels.
// offset - ammount to offset each box. useful values are half of margin, =margin, or zero.
// returns the number of unpacked rects remaining
// nb: although this looks like boxes is passed by-value, a slice type is just accounting ints and a ptr to its own data.
// extra dereferencing wouldn't benefit us as we don't append or remove from the slice.
func PackBoxes(boxes []BoxTranslation, W, H, boxMargin, offset int) int {
	stbr := C.allocateRects(C.int(len(boxes)))
	defer C.myFree(unsafe.Pointer(stbr))
	boxesToSTBR(boxes, stbr, boxMargin)
	ctx := C.allocateCTX()
	defer C.myFree(unsafe.Pointer(ctx))
	nodeCount := max(512, W, len(boxes))
	nodes := C.allocateNodes(C.int(nodeCount))
	defer C.myFree(unsafe.Pointer(nodes))
	C.stbrp_init_target(ctx, C.int(W), C.int(H), nodes, C.int(nodeCount))
	C.stbrp_pack_rects(ctx, stbr, C.int(len(boxes)))

	var box C.stbrp_rect
	var unpacked int
	for i := 0; i < len(boxes); i++ {
		C.getValue(stbr, C.int(i), &box)
		if box.was_packed > 0 {
			id := int(box.id)
			boxes[id].wasPacked = true
			w, h := boxes[id].sourceRect.Dx(), boxes[id].sourceRect.Dy()
			boxes[id].destRect.Min.X = int(box.x) + offset
			boxes[id].destRect.Min.Y = int(box.y) + offset
			boxes[id].destRect.Max.X = boxes[id].destRect.Min.X + w
			boxes[id].destRect.Max.Y = boxes[id].destRect.Min.Y + h
		} else {
			unpacked++
		}
	}
	return unpacked
}

// Creates a new atlas image based on the input images and packed boxes.
// typically used after ImageToBoxes and PackBoxes
func RenderNewAtlas(images []image.Image, boxes []BoxTranslation, outImg draw.Image) {
	dx, dy := getMaxSourceRectSizes(boxes)
	nrgba := image.NewNRGBA(image.Rect(0, 0, dx, dy))
	for _, box := range boxes {
		if box.deferredRotate {
			// this has the bizzare implication that the source W & H need to be swapped first.
			// we left them "wrong" prior to packing in order to produce a correct destination rect
			box.sourceRect.Max = image.Point{
				X: box.sourceRect.Min.X + box.sourceRect.Dy(),
				Y: box.sourceRect.Min.Y + box.sourceRect.Dx(),
			}

			// rotation
			bufferRect := image.Rect(0, 0, box.sourceRect.Dx(), box.sourceRect.Dy())
			draw.Draw(nrgba, bufferRect, images[box.imgSrc], box.sourceRect.Min, draw.Src)
			croppedBuffer := nrgba.SubImage(bufferRect)
			rotatedImage := imaging.Rotate270(croppedBuffer)

			draw.Draw(outImg, box.destRect, rotatedImage, image.Point{0, 0}, draw.Src)
		} else {
			draw.Draw(outImg, box.destRect, images[box.imgSrc], box.sourceRect.Min, draw.Src)
		}
	}
}

// returns the maximum width and heights represented in the set of source rects.
func getMaxSourceRectSizes(boxes []BoxTranslation) (int, int) {
	dx, dy := 0, 0
	for _, b := range boxes {
		dx = max(b.sourceRect.Dx(), dx)
		dy = max(b.sourceRect.Dy(), dy)
	}
	return dx, dy
}

/*
Converts a slice of Box into a C array of stbrp_rect via the dimensions of box.sourceRect
stbr pointer is presumed to point to an array of sufficient size.
*/
func boxesToSTBR(boxes []BoxTranslation, stbr *C.stbrp_rect, margin int) {
	var box C.stbrp_rect
	for i := 0; i < len(boxes); i++ {
		box.id = C.int(i)
		box.w = C.int(boxes[i].sourceRect.Dx() + margin)
		box.h = C.int(boxes[i].sourceRect.Dy() + margin)
		C.assignValue(stbr, C.int(i), &box)
	}
}

// draws all of either the source or destination set of rects in a []BoxTranslation onto a new RGBA image.
func DebugViewRects(boxes []BoxTranslation, W, H int, drawSrcRects bool, imgSrc int) image.Image {
	img := image.NewRGBA64(image.Rect(0, 0, W, H))
	rectCol := image.NewUniform(color.RGBA{255, 255, 255, 255})

	for _, b := range boxes {
		if b.imgSrc != imgSrc {
			continue
		}
		if drawSrcRects {
			draw.Draw(img, b.sourceRect, rectCol, image.ZP, draw.Src)
		} else {
			draw.Draw(img, b.destRect, rectCol, image.ZP, draw.Src)
		}
	}

	//TODO: make this useful.
	/*
		textCol := image.NewUniform(color.RGBA{255, 0, 0, 255})
		d := &font.Drawer{
			Dst:  img,
			Src:  textCol,
			Face: basicfont.Face7x13,
		}
		// second pass for text
		i := 0 //nb: not using for loop's counter as we may skip drawing some boxes.
		for _, b := range boxes {
			if b.imgSrc != imgSrc {
				continue
			}
			if drawSrcRects {
				d.Dot.X = fixed.I(b.sourceRect.Min.X)
				d.Dot.Y = fixed.I(b.sourceRect.Min.Y + 9)
				d.DrawString(strconv.Itoa(i))
			} else {
				d.Dot.X = fixed.I(b.destRect.Min.X)
				d.Dot.Y = fixed.I(b.destRect.Min.Y + 9)
				d.DrawString(strconv.Itoa(i))
			}
			i++
		}
	*/
	return img
}
