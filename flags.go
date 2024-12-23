package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type myFlags struct {
	outputFileName                                      string
	checkDiagonals, maximumMarginMode, loadAtlas, debug bool
	width, height, margin, align, minimumSquareMode     int

	atlasFilter string
}

func initFlags() {
	flag.Usage = usage
}

func getFlags() (myFlags, []string) {
	var flags myFlags
	flag.StringVar(&flags.outputFileName, "o", "output.png",
		"Filename of output.")
	flag.StringVar(&flags.atlasFilter, "filter", "",
		"Comma separated string of attachment names in the atlas file to allow. Case insensitive.")
	flag.BoolVar(&flags.loadAtlas, "atlas", false,
		"When set, loads pixel region information from .atlas files with same name.")
	flag.BoolVar(&flags.debug, "debug", false,
		"When set, writes a debug.png image demonstrating all detected/loaded islands.")
	flag.BoolVar(&flags.checkDiagonals, "diagonal", false,
		"When set, diagonally adjacent pixels are considered connected during island detection.")
	flag.BoolVar(&flags.maximumMarginMode, "findmaxmargin", false,
		"When set, will find the largest margin value for which all islands still fit in the output.")
	flag.IntVar(&flags.minimumSquareMode, "findminsquare", 0,
		"If set > 0, finds the smallest output image size for which w and h is a multiple of this value.")
	flag.IntVar(&flags.width, "w", 512,
		"Width of output image.")
	flag.IntVar(&flags.height, "h", 512,
		"Height of output image.")
	flag.IntVar(&flags.margin, "margin", 1,
		"Margin to use for each box.")
	flag.IntVar(&flags.align, "align", 1,
		"How to align a box within its margin?\n0 = top left, 1 = center, 2 = bottom right.")

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
