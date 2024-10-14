// package to make stb_rect_pack.h usable in Go.
package boxes

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
)

type Box struct {
	imgSrc     int             // which input image is this box from?
	sourceRect image.Rectangle // pixel locations on original input image
	destRect   image.Rectangle // destination rect.
	wasPacked  bool            // true if this box has been successfully packed
}

// Packs boxes, with multiple output sheets. Input slice isn't modified.
// boxMargin - additinal padding to provide each box in total, pixels.
// offset - ammount to offset each box. useful values are half of margin, =margin, or zero.
// returns result and count of any remaining unpacked (size with margin > W or H)
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
