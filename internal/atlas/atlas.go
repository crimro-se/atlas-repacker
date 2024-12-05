package atlas

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"io"
	"path/filepath"
	"strconv"
	"strings"
)

type AtlasRegions map[string]RotatableRect

// a rect that may require rotation
type RotatableRect struct {
	image.Rectangle
	RotateRequired bool // True if this region should be rotated 90 degrees clockwise
}

// Edits filenames, replacing extensions with .atlas
func FilepathsToDotAtlas(filenames []string) []string {
	modifiedFilenames := make([]string, len(filenames))
	for i, filename := range filenames {
		// Extract the file name without the extension
		filename = filename[:len(filename)-len(filepath.Ext(filename))]
		modifiedFilenames[i] = fmt.Sprintf("%s.atlas", filename)

	}
	return modifiedFilenames
}

// parse atlas file data
func ParseAtlasFile(data io.Reader) (AtlasRegions, error) {
	filedata, err := parseAtlasFileToMap(data)
	if err != nil {
		return nil, err
	}

	regions := make(AtlasRegions)
	for currentRegion, attrs := range filedata {
		rotate := attrs["rotate"] == "true" || attrs["rotate"] == "90"
		if bounds, ok := attrs["bounds"]; ok {
			x, y, w, h, err := parse4Ints(bounds)
			if err != nil {
				return nil, err
			}
			regions[currentRegion] = buildRect(x, y, w, h, rotate)
		} else if xy, ok := attrs["xy"]; ok {
			size, ok := attrs["size"]
			if !ok {
				return nil, fmt.Errorf("error in atlas file, xy attribute presented but size is missing")
			}
			x, y, err := parse2Ints(xy)
			if err != nil {
				return nil, err
			}
			w, h, err := parse2Ints(size)
			if err != nil {
				return nil, err
			}
			regions[currentRegion] = buildRect(x, y, w, h, rotate)
		} else {
			return nil, fmt.Errorf("error in atlas file, boundary completely unknown for '%s'", currentRegion)
		}
	}
	return regions, nil
}

// digests the file into a map of string (region name) to (map of (parameter name) to parameter value.)
func parseAtlasFileToMap(data io.Reader) (map[string]map[string]string, error) {
	mappedData := make(map[string]map[string]string, 0)
	scanner := bufio.NewScanner(data)
	currentRegion := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 || strings.HasSuffix(line, ".png") || strings.HasSuffix(line, ".webp") || strings.HasSuffix(line, ".gif") {
			continue
		}
		if !strings.Contains(line, ":") {
			currentRegion = line
			if _, exists := mappedData[currentRegion]; !exists {
				mappedData[currentRegion] = make(map[string]string)
			}
			continue
		}

		// If we haven't encountered a region name yet, skip this line
		if currentRegion == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			return nil, errors.New("problem parsing atlas file")
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		mappedData[currentRegion][key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return mappedData, nil
}

// builds a rect that might require a deferred rotation
func buildRect(x, y, w, h int, rotate bool) RotatableRect {
	var r RotatableRect
	r.RotateRequired = rotate
	r.Rectangle = image.Rect(x, y, x+w, y+h)
	return r
}

func parse2Ints(str string) (int, int, error) {
	coords := strings.Split(str, ",")
	if len(coords) != 2 {
		return 0, 0, errors.New("expected 2 ints whilst parsing atlas file")
	}
	x, err := strconv.Atoi(strings.TrimSpace(coords[0]))
	if err != nil {
		return 0, 0, err
	}
	y, err := strconv.Atoi(strings.TrimSpace(coords[1]))
	if err != nil {
		return 0, 0, err
	}
	return x, y, nil
}

func parse4Ints(str string) (int, int, int, int, error) {
	coords := strings.Split(str, ",")
	if len(coords) != 4 {
		return 0, 0, 0, 0, errors.New("expected 4 ints whilst parsing atlas file")
	}
	var fourInts [4]int
	var err error
	for i, v := range coords {
		fourInts[i], err = strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return 0, 0, 0, 0, err
		}
	}
	return fourInts[0], fourInts[1], fourInts[2], fourInts[3], nil
}
