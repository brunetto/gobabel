package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	tGlob0 := time.Now()

	var inPath string
	var inFile string
	var fileObj *os.File
	var nReader *bufio.Reader
	var readLine string
	var err error

	// Open the file
	if fileObj, err = os.Open(filepath.Join(inPath, inFile)); err != nil {
		log.Fatal(os.Stderr, "%v, Can't open %s: error: %s\n", os.Args[0], inFile, err)
		panic(err)
	}
	defer fileObj.Close()
	
	// Create a reader to read the file
	nReader = bufio.NewReader(fileObj)
	
	// Map with key = starId and value StarData struct
	Cluster = make([uint64]Particle)
	
	// Read the file
	line := 0
	for {
		if readLine, err = nReader.ReadString('\n'); err != nil {
			log.Println("Done reading ", line, " lines from file with err", err)
			break
		}
		if readLine[0] == '#' {
			log.Println("Header detected, skip...")
			continue
		}
		// Reading part
		lineSlice := strings.Split(readLine, " ")
		
		
		
		line ++
	}


	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for all ", tGlob1.Sub(tGlob0))
} //End main

type Cluster map[uint]Particle

type Particle struct {
	Id string
	OriginalId string
	Position []float64
	Velocity []float64
	Acceleration [][]float64
}




