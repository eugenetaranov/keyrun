package main

import (
	"io/ioutil"
	"log"
	"strings"
)

const extension = ".enc"

func findFiles() []string {
	var filteredFiles []string
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), extension) {
			fname := strings.TrimSuffix(file.Name(), extension)
			filteredFiles = append(filteredFiles, fname)
		}
	}
	return filteredFiles
}
