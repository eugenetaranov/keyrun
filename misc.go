package main

import (
	"io/ioutil"
	"log"
	"strings"
)

const extension = ".enc"

func findFiles() map[string]string {
	filteredFiles := make(map[string]string)

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), extension) {
			fname := strings.TrimSuffix(file.Name(), extension)
			filteredFiles[fname] = ""
		}
	}
	return filteredFiles
}
