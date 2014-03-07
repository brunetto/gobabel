package nBodyTools

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
// 	"time"
)

type Particle struct {
	Id string
	OriginalId string
	Mass float64
	Position []float64
	Velocity []float64
// 	Acceleration [][]float64
}

func (particle Particle) Print() () {
	fmt.Println(particle.Id,
				particle.OriginalId,
				particle.Mass,
				particle.Position, 
				particle.Velocity,
	)
}

func (particle *Particle) KineticEnergy() (float64) {
	squaredSpeed := 0.
	for _, v := range particle.Velocity {
		squaredSpeed = squaredSpeed + (v*v)
	}
	return 0.5 * particle.Mass * squaredSpeed
}

type Cluster struct {
	PotentialEnergy float64
	Particles []*Particle
}

// Read the file
func (cluster *Cluster) Populate (inPath string, inFile string) () {
	
	var fileObj *os.File
	var nReader *bufio.Reader
	var readLine string
	var err error
	
// 	var cluster.Particles []*Particle
	
	// Open the file
	if fileObj, err = os.Open(filepath.Join(inPath, inFile)); err != nil {
		log.Fatal(os.Stderr, "%v, Can't open %s: error: %s\n", os.Args[0], inFile, err)
		panic(err)
	}
	defer fileObj.Close()
	
	// Create a reader to read the file
	nReader = bufio.NewReader(fileObj)
	
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
		lineSlice := strings.Fields(readLine)
		if len(lineSlice) == 0 {
			fmt.Println("Found empty line ", line)
			continue
		}
		// Create particle
		cluster.Particles = append(cluster.Particles, &Particle{Position: make([]float64, 3), 
											Velocity: make([]float64, 3)})
		fmt.Println(cluster.Particles)
		// Fill particle
		cluster.Particles[line].Id = strconv.Itoa(line)
		cluster.Particles[line].OriginalId = lineSlice[0]
		cluster.Particles[line].Mass, _ = strconv.ParseFloat(lineSlice[1], 64)
		for idx:=0; idx<3; idx++ {
			cluster.Particles[line].Position[idx], _ = strconv.ParseFloat(lineSlice[2+idx], 64)
			cluster.Particles[line].Velocity[idx], _ = strconv.ParseFloat(lineSlice[5+idx], 64)
		}
		fmt.Println(line)
		line ++
	}
	cluster.Particles = make([]*Particle, len(cluster.Particles))
	fmt.Println(cluster.Particles)
}

func (cluster *Cluster) Print() () {
	fmt.Println("Id, OriginalID, Mass, x, y, z, vx, vy, vz")
	fmt.Println("=========================================")
	fmt.Println(len(cluster.Particles))
	for idx, element := range(cluster.Particles) {
		fmt.Print(idx, " ")
		element.Print()
	}
}

func (cluster *Cluster) KineticEnergy() (kineticEnergy float64) {
	kineticEnergy = 0
	for _, particle := range cluster.Particles {
		kineticEnergy = kineticEnergy + particle.KineticEnergy()
	}
	return kineticEnergy
}

