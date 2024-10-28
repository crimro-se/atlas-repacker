package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"

	"github.com/crimro-se/atlas-repacker/internal/boxpack"
	_ "golang.org/x/image/webp"
)

func init() {
	initLogging()
}

func main() {
	//
	// 1. Flag parsing
	//
	flags, inputFiles := getFlags()
	errs := validateFlags(flags, inputFiles)
	if len(errs) > 0 {
		logErrors(errs)
		flag.Usage()
		os.Exit(1)
	}

	//
	// 2. Box Packing
	//
	images := must1(loadAllImages(inputFiles))
	boxes := boxpack.ImagesToBoxes(images, flags.checkDiagonals)
	if len(boxes) < 1 {
		errHandler(errors.New("no pixel islands detected in the input image"))
	}
	unpacked := boxpack.PackBoxes(boxes, flags.width, flags.height, flags.margin, getOffset(flags))

	// maximum margin finder
	// todo: double margin then backoff in a binary-search fashion
	if flags.maximumMargin && unpacked == 0 {
		for unpacked == 0 {
			flags.margin++
			unpacked = boxpack.PackBoxes(boxes, flags.width, flags.height, flags.margin, getOffset(flags))
		}
		flags.margin--
		fmt.Println("Margin chosen: ", flags.margin)
		unpacked = boxpack.PackBoxes(boxes, flags.width, flags.height, flags.margin, getOffset(flags))
	}

	if unpacked > 0 {
		fmt.Println("Note: ", unpacked, "boxes couldn't be packed")
	}
	outImg := image.NewNRGBA(image.Rect(0, 0, flags.width, flags.height))
	boxpack.RenderNewAtlas(images, boxes, outImg)

	errHandler(saveImage(flags.outputFileName, outImg))
}

func getOffset(flags myFlags) int {
	var offset int
	switch flags.align {
	case 0:
		offset = 0
	case 1:
		offset = flags.margin / 2
	case 2:
		offset = flags.margin
	}
	return offset
}

func saveImage(fileName string, img image.Image) error {
	fp, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer fp.Close()
	err = png.Encode(fp, img)
	return err
}

// loads all input files as image.Image types.
func loadAllImages(files []string) ([]image.Image, error) {
	images := make([]image.Image, 0, len(files))
	for _, inputFile := range files {
		fp, err := os.Open(inputFile)
		if err != nil {
			return images, err
		}
		defer fp.Close()
		img, _, err := image.Decode(fp)
		if err != nil {
			return images, err
		}
		images = append(images, img)
	}
	return images, nil
}
