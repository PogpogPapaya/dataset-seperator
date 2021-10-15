package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	inputDir := "dataset"
	outputDir := "papaya"
	dataClasses, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalln("Unable to read dir in inputDir", err)
	}

	datasetClasses := []string{"train", "validate", "test"}
	numberDatasetClass := []int{70, 15, 15}
	numberDatasetPerClass := 200

	for _, dataClass := range dataClasses {
		fileLists, err := os.ReadDir(fmt.Sprintf("%s/%s", inputDir, dataClass.Name()))
		if err != nil {
			log.Fatalln("Unable to read files in ", dataClass.Name(), err)
		}
		numberOfFiles := len(fileLists)
		if numberDatasetPerClass <= 0 {
			numberDatasetPerClass = numberOfFiles
		}

		for j, datasetClass := range datasetClasses {
			destDir := fmt.Sprintf("%s/%s/%s", outputDir, datasetClass, dataClass.Name())
			if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
				log.Fatalln("Unable to create class dir", datasetClass, err)
			}
			for i := 0; i<numberDatasetClass[j] * numberDatasetPerClass / 100; i++ {
				fileName := fileLists[randomInt(numberOfFiles)].Name()
				inputFile := fmt.Sprintf("%s/%s/%s", inputDir, dataClass.Name(),fileName)
				outputFile := fmt.Sprintf("%s/%s", destDir, fileName)
				_, err = copy(inputFile, outputFile)
				if err != nil {
					log.Fatalln("Unable to copy file", inputFile, err)
				}
			}
		}
	}
}

func randomInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}