package main

import (
// 	"bufio"
	"fmt"
	"log"
	nbt "./nBodyTools"
// 	"os"
// 	"path/filepath"
// 	"reflect"
// 	"regexp"
// 	"sort"
// 	"strconv"
// 	"strings"
	"time"
)

func main() {
	tGlob0 := time.Now()

	var inPath string
	var inFile string
	var cluster nbt.Cluster
	
	inPath = "inputFilesSmall"
	inFile = "input1k"

	cluster = *new(nbt.Cluster)
	
	
	cluster.Populate(inPath, inFile)
	fmt.Println(cluster.Particles)
	cluster.Print()
	fmt.Println(cluster.KineticEnergy())

	tGlob1 := time.Now()
	fmt.Println()
	log.Println("Wall time for all ", tGlob1.Sub(tGlob0))
} //End main




