package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type myFlags struct {
	outputFileName, inputAtlas                      string
	checkDiagonals, maximumMarginMode               bool
	width, height, margin, align, minimumSquareMode int
}

func getFlags() (myFlags, []string) {
	var flags myFlags
	flag.StringVar(&flags.outputFileName, "o", "output.png", "filename of output")
	//todo:
	//flag.StringVar(&flags.inputAtlas, "atlas", "", "filename of atlas file that defines all the boxes in the input image")
	flag.BoolVar(&flags.checkDiagonals, "diagonal", false,
		"when set, diagonally adjacent pixels are considered connected during island detection.")
	flag.BoolVar(&flags.maximumMarginMode, "findmaxmargin", false,
		"when set, will find the largest margin value for which all islands still fit in the output.")
	flag.IntVar(&flags.minimumSquareMode, "findminsquare", 0,
		"If set > 0, finds the smallest output image size for which w and h is a multiple of this value.")
	flag.IntVar(&flags.width, "w", 512, "width of output image")
	flag.IntVar(&flags.height, "h", 512, "height of output image")
	flag.IntVar(&flags.margin, "margin", 1, "margin to use for each box")
	flag.IntVar(&flags.align, "align", 1, "how to align a box within its margin?\n0 = top left, 1 = center, 2 = bottom right")

	flag.Usage = usage
	flag.Parse()
	inputFiles := flag.Args()
	return flags, inputFiles
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n\n", os.Args[0])
	fmt.Fprintln(flag.CommandLine.Output(), os.Args[0], "[flags]", "[input.png] [input2.png ...]")
	fmt.Fprintln(flag.CommandLine.Output(), "Flags:")
	flag.PrintDefaults()
}

// attempt to validate flags. Any issues returned as errors
func validateFlags(flags myFlags, inputs []string) []error {
	var errs []error

	if len(inputs) < 1 {
		errs = append(errs, errors.New("no input files specified"))
	}

	if flags.align < 0 || flags.align > 2 {
		errs = append(errs, errors.New("invalid alignment. Should be 0, 1 or 2"))
	}

	if flags.margin < 0 || flags.width < 1 || flags.height < 1 {
		errs = append(errs, errors.New("an input parameter specified is too small or negative"))
	}
	return errs
}
