package main

import (
	"flag"
	"fmt"
	"os"
)

type myFlags struct {
	outputFileName                string
	checkDiagonals, maximumMargin bool
	width, height, margin, align  int
}

func getFlags() (myFlags, []string) {
	var flags myFlags
	flag.StringVar(&flags.outputFileName, "o", "output.png", "filename of output")
	flag.BoolVar(&flags.checkDiagonals, "diagonal", false,
		"when set, diagonally adjacent pixels are considered connected during island detection.")
	flag.BoolVar(&flags.maximumMargin, "findmaxmargin", false,
		"when set, will find the largest margin value for which all islands still fit in the output.")
	flag.IntVar(&flags.width, "w", 512, "width of output image")
	flag.IntVar(&flags.height, "h", 512, "height of output image")
	flag.IntVar(&flags.margin, "margin", 1, "margin to use for each box")
	flag.IntVar(&flags.align, "align", 1, "how to align a box within its margin?\n0 = top left, 1 = center, 2 = bottom right")

	flag.Usage = usage
	flag.Parse()
	inputFiles := flag.Args()

	isValid := validate(flags, inputFiles)
	if !isValid {
		os.Exit(1)
	}
	return flags, inputFiles
}

func usage() {
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
