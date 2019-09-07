// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/golang/go/src/pkg/io/ioutil"
	"github.com/seamia/tools/assets"
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
		lines = append(lines, expand(scanner.Text()))
	}
	return lines, scanner.Err()
}

func expand(from string) string {
	return os.ExpandEnv(from)
}

func getHeader(headerName string) (string, error) {
	fileName := expand(headerName)
	if len(fileName) == 0 {
		return "", nil
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	txt := expand(string(data))
	return txt, nil
}

func main() {

	var listOfFiles, root, destination, packageName, headerName string
	flag.StringVar(&listOfFiles, "src", "", "Name of the file where files to be included are listed")
	flag.StringVar(&root, "root", "", "Path to be considered as root (within included files)")
	flag.StringVar(&destination, "output", "staticAssets.go", "Name of the output file")
	flag.StringVar(&packageName, "package", "main", "Name of the package to be used in the generated file")
	flag.StringVar(&headerName, "header", "", "Name of the file to be used as a header in the generated file")
	flag.Parse()

	if len(listOfFiles) == 0 {
		fmt.Println("Static Assets Generator (github.com/seamia/tools/assets/cmd/assets)")
		fmt.Println("Please specify appropriate parameters")
		flag.PrintDefaults()
		return
	}

	filenames, err := readLines(expand(listOfFiles))
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to process the list of files: ", err)
		return
	}

	header, err := getHeader(headerName)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to process the header file: ", err)
		return
	}

	err = assets.Generate(filenames, expand(root), expand(destination), expand(packageName), header)
	if err != nil {
		fmt.Fprint(os.Stderr, "There was an error: ", err)
		return
	}
	fmt.Fprint(os.Stdout, "Generated file: ", destination)
}
