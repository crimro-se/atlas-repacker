/*	This script counts the occurance of attachment names in all .atlas
	files under the specified input directory and prints them out.
	Usage:
	go run atlas_counter.go /input/directory
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/exp/constraints"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("specify a folder to search on the commandline!")
		return
	}
	dirPath := os.Args[1]

	lineCounts := make(map[string]int)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".atlas") {
			file, err := os.Open(path)
			if err != nil {
				fmt.Printf("Failed to open file %s: %v\n", path, err)
				return nil
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			lineNumber := 0

			for scanner.Scan() {
				lineNumber++
				line := scanner.Text()

				if lineNumber <= 2 {
					continue // Skip the first two lines
				}

				firstNonWhitespaceIndex := strings.IndexFunc(line, func(r rune) bool { return strings.ContainsRune(" \t", r) })
				if firstNonWhitespaceIndex == -1 {
					lineCounts[line]++
				}
			}

			if err := scanner.Err(); err != nil {
				fmt.Printf("Error reading file %s: %v\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Failed to walk directory: %v\n", err)
	}
	sortedKeys := SortKeysByValues(lineCounts)

	for _, key := range sortedKeys {
		fmt.Printf("%s: %d\n", key, lineCounts[key])
	}
}

func SortKeysByValues[K comparable, V constraints.Integer | constraints.Float](m map[K]V) []K {
	var keys []K
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	return keys
}
