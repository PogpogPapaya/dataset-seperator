package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

func main() {
	inputDir := "dataset"
	outputDir := "papaya"

	datasetClasses := []string{"train", "validate", "test"}
	numberDatasetClass := []int{70, 15, 15}
	numberDatasetPerClass := 0

	err := os.RemoveAll(outputDir)
	if err != nil {
		log.Fatalln("Unable to cleanup output dir", err)
	}

	dataClasses, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatalln("Unable to read dir in inputDir", err)
	}
	// Balance or not
	numberPerClass, min := findNumberPerClass(dataClasses, inputDir)
	for i, num := range numberPerClass {
		fmt.Printf("%s: %d\n", dataClasses[i].Name(), num)
	}
	fmt.Printf("The dataset will be balance to %d data per class", min)
	numberDatasetPerClass = min

	for _, dataClass := range dataClasses {
		if !dataClass.IsDir() {
			continue
		}
		fileLists, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", inputDir, dataClass.Name()))
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
			for i := 0; i < numberDatasetClass[j]*numberDatasetPerClass/100; i++ {
				fileIndex := randomInt(numberOfFiles)
				inputFile := fmt.Sprintf("%s/%s/%s", inputDir, dataClass.Name(), fileLists[fileIndex].Name())
				outputFile := fmt.Sprintf("%s/%s", destDir, fileLists[fileIndex].Name())
				_, err = copy(inputFile, outputFile)
				if err != nil {
					log.Fatalln("Unable to copy file", inputFile, err)
				}
				fileLists = append(fileLists[:fileIndex], fileLists[fileIndex+1:]...)
				numberOfFiles = len(fileLists)
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

func findNumberPerClass(classes []os.FileInfo, inputDir string) ([]int, int) {
	min := math.MaxInt32
	var count []int
	for _, class := range classes {
		if !class.IsDir() {
			continue
		}
		fileList, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", inputDir, class.Name()))
		if err != nil {
			log.Fatalln("Unable to read dir", class.Name(), err)
		}
		numberOfFiles := len(fileList)
		if numberOfFiles < min {
			min = numberOfFiles
		}
		count = append(count, numberOfFiles)
	}
	return count, min
}
