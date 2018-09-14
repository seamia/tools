package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/seamia/tools/assets"
	"os"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {

	var listOfFiles, root, destination, packageName string
	flag.StringVar(&listOfFiles, "src", "", "Name of the file where files to be included are listed")
	flag.StringVar(&root, "root", "", "Path to be considered as root (within included files)")
	flag.StringVar(&destination, "output", "staticAssets.go", "Name of the output file")
	flag.StringVar(&packageName, "package", "main", "Name of the package to be used in the generated file")
	flag.Parse()

	if len(listOfFiles) == 0 {
		fmt.Println("Static Assets Generator (github.com/seamia/tools/assets/cmd/assets)")
		fmt.Println("Please specify appropriate parameters")
		flag.PrintDefaults()
		return
	}

	filenames, err := readLines(listOfFiles)
	if err != nil {
		fmt.Println("Failed to process the list of files")
		return
	}

	err = assets.Generate(filenames, root, destination, packageName)
	if err != nil {
		fmt.Println("There was an error:", err)
	}
}
