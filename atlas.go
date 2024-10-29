package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/crimro-se/atlas-repacker/internal/boxpack"
)

type atlasRegions map[string]image.Rectangle

// Edits filenames, replacing extensions with .atlas
func filesToDotAtlas(filenames []string) []string {
	modifiedFilenames := make([]string, len(filenames))
	for i, filename := range filenames {
		// Extract the file name without the extension
		filename = filename[:len(filename)-len(filepath.Ext(filename))]
		modifiedFilenames[i] = fmt.Sprintf("%s.atlas", filename)

	}
	return modifiedFilenames
}

func loadAllAtlas(files []string) ([]boxpack.BoxTranslation, error) {
	boxes := make([]boxpack.BoxTranslation, 0, max(len(files), 10))
	for i, inputFile := range files {
		fp, err := os.Open(inputFile)
		if err != nil {
			return boxes, err
		}
		defer fp.Close()

		atlas, err := parseAtlasFile(fp)
		if err != nil {
			return boxes, err
		}
		boxes = append(boxes, atlasToBoxes(i, atlas)...)
	}
	return boxes, nil
}

// converts atlasRegions type to []boxpack.BoxTranslation
func atlasToBoxes(refImage int, ar atlasRegions) []boxpack.BoxTranslation {
	boxes := make([]boxpack.BoxTranslation, 0, len(ar))
	for _, v := range ar {
		boxes = append(boxes, boxpack.BoxFromRect(refImage, v))
	}
	return boxes
}

// parseAtlasFile reads an atlas file and returns a map of region names to image.Rectangles
func parseAtlasFile(data io.Reader) (atlasRegions, error) {

	regions := make(atlasRegions)
	currentRegion := ""
	currentRect := image.Rectangle{}
	rotate := false

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and the first line (assuming it's the image file name)
		if len(line) == 0 || strings.HasSuffix(line, ".png") || strings.HasSuffix(line, ".webp") {
			continue
		}

		// Check if the line starts a new region
		if !strings.Contains(line, ":") {
			currentRegion = line
			continue
		}

		// If we haven't encountered a region name yet, skip this line
		if currentRegion == "" {
			continue
		}

		// Parse attributes for the current region
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			log.Printf("Warning: Skipping malformed line - '%s'\n", line)
			continue
		}

		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "rotate":
			rotate = (value == "true")

		case "xy":
			coords := strings.Split(value, ",")
			if len(coords) != 2 {
				return nil, errors.New("invalid xy coordinates")
			}
			x, err := strconv.Atoi(strings.TrimSpace(coords[0]))
			if err != nil {
				return nil, err
			}
			y, err := strconv.Atoi(strings.TrimSpace(coords[1]))
			if err != nil {
				return nil, err
			}
			// Store the xy for later use with size
			currentRect.Min = image.Pt(x, y)
			regions[currentRegion] = currentRect
		case "size":
			dims := strings.Split(value, ",")
			if len(dims) != 2 {
				return nil, errors.New("invalid size dimensions")
			}
			w, err := strconv.Atoi(strings.TrimSpace(dims[0]))
			if err != nil {
				return nil, err
			}
			h, err := strconv.Atoi(strings.TrimSpace(dims[1]))
			if err != nil {
				return nil, err
			}
			// NB: we add Min to this later.
			// this enables xy clause to be before or after size.
			// depends on rotate coming first still, but it seems exporters always do this anyway.
			if rotate {
				currentRect.Max = image.Pt(h, w)
			} else {
				currentRect.Max = image.Pt(w, h)
			}
			regions[currentRegion] = currentRect
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	for name, rect := range regions {
		rect.Max = rect.Max.Add(rect.Min)
		regions[name] = rect
	}

	return regions, nil
}
