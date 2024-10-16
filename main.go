package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"

	"github.com/crimro.se/atlas-repacker/internal/boxpack"
	_ "golang.org/x/image/webp"
)

type myFlags struct {
	outputFileName               string
	checkDiagonals               bool
	width, height, margin, align int
}

func main() {
	//
	// 1. Flag parsing
	//
	var flags myFlags
	flag.StringVar(&flags.outputFileName, "o", "output.png", "filename of output")
	flag.BoolVar(&flags.checkDiagonals, "diagonal", false,
		"when set, diagonally adjacent pixels are considered connected during island detection.")
	flag.IntVar(&flags.width, "w", 512, "width of output image")
	flag.IntVar(&flags.height, "h", 512, "height of output image")
	flag.IntVar(&flags.margin, "margin", 1, "margin to use for each box")
	flag.IntVar(&flags.align, "align", 1, "how to align a box within its margin?\n0 = top left, 1 = center, 2 = bottom right")

	flag.Usage = Usage
	flag.Parse()
	inputFiles := flag.Args()

	isValid := validate(flags, inputFiles)
	if !isValid {
		os.Exit(1)
		return
	}

	var offset int
	switch flags.align {
	case 0:
		offset = 0
	case 1:
		offset = flags.margin / 2
	case 2:
		offset = flags.margin
	}

	//
	// 2. Box Packing
	//
	images, err := loadAllImages(inputFiles)
	errHandler(err)
	boxes := boxpack.ImagesToBoxes(images, flags.checkDiagonals)
	if len(boxes) < 1 {
		fmt.Println("!! No pixel islands were located in input images. Aborting.")
		return
	}
	unpacked := boxpack.PackBoxes(boxes, flags.width, flags.height, flags.margin, offset)
	if unpacked > 0 {
		fmt.Println("Note: ", unpacked, "boxes couldn't be packed")
	}
	outImg := image.NewNRGBA(image.Rect(0, 0, flags.width, flags.height))
	boxpack.RenderNewAtlas(images, boxes, outImg)

	err = saveImage(flags.outputFileName, outImg)
	errHandler(err)
}

func Usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n\n", os.Args[0])
	fmt.Fprintln(flag.CommandLine.Output(), os.Args[0], "[flags]", "[input.png] [input2.png ...]")
	fmt.Fprintln(flag.CommandLine.Output(), "Flags:")
	flag.PrintDefaults()
}

// attempt to validate flags. Tells the user about any issues
func validate(flags myFlags, inputs []string) bool {
	valid := true
	if len(inputs) < 1 {
		fmt.Println("!! No input files specified!")
		valid = false
	}

	if flags.align < 0 || flags.align > 2 {
		fmt.Println("!! invalid alignment. Should be 0, 1 or 2")
		valid = false
	}

	if flags.margin < 0 || flags.width < 1 || flags.height < 1 {
		fmt.Println("!! An input parameter specified is too small or negative")
		valid = false
	}
	if !valid {
		flag.Usage()
	}
	return valid
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

func loadAllImages(files []string) ([]image.Image, error) {
	images := make([]image.Image, 0, len(files))
	for _, inputFile := range files {
		fp, err := os.Open(inputFile)
		if err != nil {
			return images, err
		}
		img, _, err := image.Decode(fp)
		if err != nil {
			return images, err
		}
		images = append(images, img)
		err = fp.Close()
		if err != nil {
			return images, err
		}
	}
	return images, nil
}

// if there's an error, show it to the user and stop execution.
func errHandler(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
