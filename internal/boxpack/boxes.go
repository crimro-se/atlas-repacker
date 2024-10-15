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
	"unsafe"

	"golang.org/x/exp/constraints"
)

type Box struct {
	imgSrc     int             // which input image is this box from?
	sourceRect image.Rectangle // pixel locations on original input image
	destRect   image.Rectangle // destination rect.
	wasPacked  bool            // true if this box has been successfully packed
}

// identifies pixel islands in an image
/*
func ImagesToBoxes(images []image.Image) []Box {
	boxes := make([]Box, 0)
	var i int
	for _, img := range images {

		i++
	}
	return boxes
}


func ImageToBoxes(img image.Image) []Box {
	images := make([]image.Image, 0, 1)
	images = append(images, img)
	return ImagesToBoxes(images)
}
*/

// given an image and a starting pixel, finds all connected pixels and returns a square encompassing them.
// 'visited' is used to track progress.
// the diagonal flag enables checking diagonally connected pixels.
func findConnectedPixels(img image.Image, x, y int, diagonal bool, visited visitedArray) image.Rectangle {
	bounds := img.Bounds()
	stack := []image.Point{{X: x, Y: y}}

	var minX, minY, maxX, maxY int
	minX = x
	maxX = x
	minY = y
	maxY = y

	for len(stack) > 0 {
		point := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		visited.Set(point.X, point.Y, true)
		minX = min(minX, point.X)
		minY = min(minY, point.Y)
		maxX = max(maxX, point.X)
		maxY = max(maxY, point.Y)

		pointCheck := func(pt image.Point) {
			if !pt.In(bounds) {
				return
			}
			if !visited.Get(pt.X, pt.Y) {
				_, _, _, a := img.At(pt.X, pt.Y).RGBA()
				if a > 0 {
					stack = append(stack, pt)
				}
			}
		}
		for xOff := -1; xOff <= 1; xOff++ {
			for yOff := -1; yOff <= 1; yOff++ {
				if abs(xOff)+abs(yOff) == 2 && !diagonal {
					continue
				}
				pointCheck(image.Point{point.X + xOff, point.Y + yOff})
			}
		}
	}
	return image.Rect(minX, minY, maxX, maxY)
}

/* Originally used an Image.Grey to track visited pixels, however it involves too much abstraction and indirection. */
type visitedArray struct {
	data []bool
	w, h int
}

func NewVisitedArray(bounds image.Rectangle) visitedArray {
	var va visitedArray
	va.w = bounds.Dx()
	va.h = bounds.Dy()
	va.data = make([]bool, va.w*va.h)
	return va
}

func (va *visitedArray) Get(x, y int) bool {
	return va.data[y*va.w+x]
}

func (va *visitedArray) Set(x, y int, v bool) {
	va.data[y*va.w+x] = v
}

func abs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Packs boxes, with multiple output sheets. Input slice isn't modified.
// boxMargin - additinal padding to provide each box in total, pixels.
// offset - ammount to offset each box. useful values are half of margin, =margin, or zero.
// returns result and count of any remaining unpacked (size+margin > W or H)
func PackAllBoxes(boxesImmutable []Box, W, H, boxMargin, offset int) ([][]Box, int) {
	allBoxes := make([][]Box, 0)
	boxes := make([]Box, len(boxesImmutable)) // working copy
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
		remainer := make([]Box, 0, unpacked)
		allBoxes = append(allBoxes, make([]Box, 0, previous-unpacked))
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
func PackBoxes(boxes []Box, W, H, boxMargin, offset int) int {
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

/*
Converts a slice of Box into a C array of stbrp_rect via the dimensions of box.sourceRect
stbr pointer is presumed to point to an array of sufficient size.
*/
func boxesToSTBR(boxes []Box, stbr *C.stbrp_rect, margin int) {
	var box C.stbrp_rect
	for i := 0; i < len(boxes); i++ {
		box.id = C.int(i)
		box.w = C.int(boxes[i].sourceRect.Dx() + margin)
		box.h = C.int(boxes[i].sourceRect.Dy() + margin)
		C.assignValue(stbr, C.int(i), &box)
	}
}
