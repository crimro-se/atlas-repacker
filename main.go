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

	"github.com/crimro-se/atlas-repacker/internal/atlas"
	"github.com/crimro-se/atlas-repacker/internal/boxpack"
	_ "golang.org/x/image/webp"
)

func init() {
	initLogging()
	initFlags()
}

func main() {
	errored := 0
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
	images, err := loadAllImages(inputFiles)
	errHandler(err, "an error occured whilst loading images")

	// find pixel islands via atlas file or look at the pixels.
	var boxes []boxpack.BoxTranslation
	boxes, err = boxPack(images, inputFiles, flags)
	errHandler(err)
	if len(boxes) < 1 {
		errHandler(errors.New("no pixel islands detected in the input image(s)"))
	}

	if flags.debug {
		img := boxpack.DebugViewRects(boxes, images[0].Bounds().Dx(), images[0].Bounds().Dy(), true, 0)
		errHandler(saveImage("debug.png", img))
		msg("debug.png has been written")
		if len(inputFiles) > 1 {
			msg("NOTE: only the first image you loaded has been debugged.")
		}
	}

	var unpacked int
	unpacked = boxpack.PackBoxes(boxes, flags.width, flags.height, flags.margin, getOffset(flags))

	//
	// 2.1 bruteforce w,h if requested
	//
	if flags.minimumSquareMode > 0 {
		wh := (boxpack.EstimateOutputWH(boxes, flags.margin) / flags.minimumSquareMode) * flags.minimumSquareMode
		unpacked = boxpack.PackBoxes(boxes, wh, wh, flags.margin, getOffset(flags))
		for unpacked > 0 {
			wh += flags.minimumSquareMode
			unpacked = boxpack.PackBoxes(boxes, wh, wh, flags.margin, getOffset(flags))
		}
		flags.width = wh
		flags.height = wh
		msg(fmt.Sprintf("Calculated output size (W&H): %d", wh))
	}

	//
	// 2.2 maximum margin finder
	//     todo: double margin then backoff in a binary-search fashion
	if flags.maximumMarginMode && unpacked == 0 {
		// we need a clone of the boxes for our tests
		boxes2 := make([]boxpack.BoxTranslation, len(boxes))
		copy(boxes2, boxes)
		for unpacked == 0 {
			flags.margin++
			unpacked = boxpack.PackBoxes(boxes2, flags.width, flags.height, flags.margin, getOffset(flags))
			if unpacked == 0 {
				boxes = boxes2
			}
		}
		unpacked = 0
		flags.margin--
		msg(fmt.Sprintf("Margin chosen: %d", flags.margin))
	}

	if flags.maximumMarginMode && unpacked > 0 {
		msg("Note: margin detection skipped as we already can't pack everything")
	}

	if unpacked > 0 {
		msg(fmt.Sprintf("Note: %d boxes couldn't be packed", unpacked))
		errored = 1
	}

	//
	// 2.3 save output
	//
	outImg := image.NewNRGBA(image.Rect(0, 0, flags.width, flags.height))
	boxpack.RenderAll(images, boxes, outImg)
	errHandler(saveImage(flags.outputFileName, outImg))

	// exit status
	os.Exit(errored)
}

func boxPack(images []image.Image, filenames []string, cfg myFlags) ([]boxpack.BoxTranslation, error) {
	boxes := make([]boxpack.BoxTranslation, 0, 8)
	atlasFiles := atlas.FilepathsToDotAtlas(filenames)
	for i, img := range images {
		detectRequired := true // disabled if we successfully load from atlas.
		if cfg.loadAtlas {
			b, e := parseAtlasFile(atlasFiles[i], i)
			if e == nil {
				boxes = append(boxes, b...)
				detectRequired = false
			}
		}
		if detectRequired {
			b, e := detectIslands(img, i, cfg.checkDiagonals)
			if e != nil {
				return boxes, e
			}
			boxes = append(boxes, b...)
		}
	}
	return boxes, nil
}

// resolves the exact pixel offset to apply based on cli flags
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
