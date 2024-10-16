package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"

	"github.com/crimro.se/atlas-repacker/internal/boxpack"
	_ "golang.org/x/image/webp"
)

func main() {
	//
	// 1. Flag parsing
	//
	outputNamePtr := flag.String("o", "output.png", "filename of output")
	diagonalPtr := flag.Bool("diagonal", false,
		"when set, diagonally adjacent pixels are considered connected during island detection.")
	widthPtr := flag.Int("w", 512, "width of output image")
	heightPtr := flag.Int("h", 512, "height of output image")
	marginPtr := flag.Int("margin", 1, "margin to use for each box")
	offsetPtr := flag.Int("offset", 0, "ammount to offset each box. Useful values are 0, margin/2, =margin")
	// todo: chroma mask alpha -> real alpha pre-process flag
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), os.Args[0], "[flags]", "[input.png] [input2.png ...]")
		fmt.Fprintln(flag.CommandLine.Output(), "Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Flag Validation
	inputFiles := flag.Args()
	if len(inputFiles) < 1 {
		fmt.Println("!! No input files specified!")
		flag.Usage()
		return
	}

	if *offsetPtr > *marginPtr {
		fmt.Println("!! offset value specified is too large (> margin)")
		return
	}

	if *offsetPtr < 0 || *marginPtr < 0 || *widthPtr < 1 || *heightPtr < 1 {
		fmt.Println("!! An input parameter specified is too small or negative")
		return
	}

	//
	// 2. Box Packing
	//
	images := make([]image.Image, 0, len(inputFiles))
	for _, inputFile := range inputFiles {
		fp, err := os.Open(inputFile)
		errHandler(err)
		img, _, err := image.Decode(fp)
		errHandler(err)
		images = append(images, img)
		errHandler(fp.Close())
	}

	// todo: loading boxes from .atlas file instead of detection
	boxes := boxpack.ImagesToBoxes(images, *diagonalPtr)

	unpacked := boxpack.PackBoxes(boxes, *widthPtr, *heightPtr, *marginPtr, *offsetPtr)
	if unpacked > 0 {
		fmt.Println("Note: ", unpacked, "boxes couldn't be packed")
	}
	outImg := image.NewNRGBA(image.Rect(0, 0, *widthPtr, *heightPtr))
	boxpack.RenderNewAtlas(images, boxes, outImg)

	//save
	fp, err := os.Create(*outputNamePtr)
	errHandler(err)
	defer fp.Close()
	err = png.Encode(fp, outImg)
	errHandler(err)

}

func errHandler(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
