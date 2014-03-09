package main

import (
	"fmt"
	"math"
	"os"
	"time"
	
	"github.com/brunetto/goutils/debug"
)

var Debug = false

type body struct {
	x, y, z, vx, vy, vz, ax, ay, az, a0x, a0y, a0z, m float64
}

var (
	bodies = make([]*body, 0)
	rij = new(struct{x, y, z float64})				// distance between two particles
)

func acceleration () () {
	if Debug{defer debug.TimeMe(time.Now())}
	// Reset acceleration
	for i:=0; i<len(bodies); i++ {
		bodies[i].ax = 0
		bodies[i].ay = 0
		bodies[i].az = 0
	}

	
	for i:=0; i<len(bodies); i++ {
		for j:=i+1; j<len(bodies); j++ {
			rij.x = bodies[i].x - bodies[j].x
			rij.y = bodies[i].y - bodies[j].y
			rij.z = bodies[i].z - bodies[j].z
			
			RdotR := (rij.x * rij.x) + (rij.y * rij.y) + (rij.z * rij.z)
			apre  := 1.0 / math.Sqrt(RdotR * RdotR * RdotR)
		
			//Update acceleration
			bodies[i].ax -= bodies[j].m * apre * rij.x
			bodies[i].ay -= bodies[j].m * apre * rij.y
			bodies[i].az -= bodies[j].m * apre * rij.z
		}
	}
}

// Update positions
func updatePositions (dt float64) () {
	if Debug{defer debug.TimeMe(time.Now())}
	
	for i:=0; i<len(bodies); i++ {
		// Update the positions, based on the calculated accelerations and velocities
		bodies[i].a0x = bodies[i].ax
		bodies[i].a0y = bodies[i].ay
		bodies[i].a0z = bodies[i].az
		// For each axis (x/y/z)
		bodies[i].x += dt * bodies[i].vx + 0.5 * dt * dt * bodies[i].a0x
		bodies[i].y += dt * bodies[i].vy + 0.5 * dt * dt * bodies[i].a0y
		bodies[i].z += dt * bodies[i].vz + 0.5 * dt * dt * bodies[i].a0z
	}
	
}

// Update velocities based on previous and new accelerations
func updateVelocities (dt float64) () {
	if Debug{defer debug.TimeMe(time.Now())}
	
	//Update the velocities based on the previous and old accelerations
	for i:=0; i<len(bodies); i++ {
		bodies[i].vx += 0.5 * dt * (bodies[i].a0x + bodies[i].ax)
		bodies[i].vy += 0.5 * dt * (bodies[i].a0y + bodies[i].ay)
		bodies[i].vz += 0.5 * dt * (bodies[i].a0z + bodies[i].az)
		
		// Update accelerations 
		bodies[i].a0x = bodies[i].ax
		bodies[i].a0y = bodies[i].ay
		bodies[i].a0z = bodies[i].az
	}  
} 


// Compute the energy of the system, 
// contains an expensive O(N^2) part which can be moved to the acceleration part
// where this is already calculated
func energies () (EKin, EPot float64) {
	if Debug{defer debug.TimeMe(time.Now())}
	
	EKin = 0
	
	//Kinetic energy
	for i:=0; i<len(bodies); i++ {
		EKin += 0.5 * bodies[i].m * bodies[i].vx * bodies[i].vx + bodies[i].vy * bodies[i].vy + bodies[i].vz * bodies[i].vz
	}
	
	//Potential energy
	for i:=0; i<len(bodies); i++ {
		for j:=i+1; j<len(bodies); j++ {
			//Distance between the two stars
			rij.x = bodies[i].x - bodies[j].x
			rij.y = bodies[i].y - bodies[j].y
			rij.z = bodies[i].z - bodies[j].z

			EPot -= (bodies[i].m * bodies[j].m) / math.Sqrt((rij.x * rij.x) + (rij.y * rij.y) + (rij.z * rij.z))
		}
	}
	return EKin, EPot
}


func main () {
	defer debug.TimeMe(time.Now())
	
	var (
		t float64    = 0.0
		tend float64 = 1.0
		dt float64   = 1e-3
		k int    = 0
		kinEnergy, potEnergy, totEnergy, totEnergy0, dE float64
		inFileName string
		inFile *os.File
		outFile *os.File
		err error
		minusOne int
	)
	
	inFileName = os.Args[1]
	
	if inFile, err = os.Open(inFileName); err != nil {panic(err)}
	defer inFile.Close()
	
	for {
		bd := new(body)
		if _, err := fmt.Fscanf(inFile, "%d %f %f %f %f %f %f %f\n", 
			&minusOne, &(bd.m), &(bd.x), &(bd.y), &(bd.z), &(bd.vx), &(bd.vy), &(bd.vx)); err != nil {
			break
		}
		bodies = append(bodies, bd)
	}
	fmt.Println("Read ", len(bodies), "lines")
	// Compute initial energy of the system
	kinEnergy, potEnergy = energies()
	totEnergy0 = kinEnergy+potEnergy
	
	fmt.Printf("Starting: Etot0=%f Ek0=%f Ep0=%f\n", totEnergy0, kinEnergy, potEnergy)
	
	//Initialize the accelerations
	acceleration()
	//Start the main loop
	for {
		if t > tend {break}
		
		// Update positions based on velocities and accelerations
		updatePositions(dt)
		
		// Get new accelerations
		acceleration()
		
		// Update velocities
		updateVelocities(dt)
		
		// Update time
		t += dt
		k += 1
		
		if k%10 == 0 {
			kinEnergy, potEnergy = energies()
			totEnergy = kinEnergy+potEnergy
			dE = (totEnergy-totEnergy0) / totEnergy0
			
			fmt.Printf("\rt= %f Etot=%f Etot0=%f Ek=%f Ep=%f dE=%f", t, totEnergy, totEnergy0, kinEnergy, potEnergy, dE)
		}	
		if Debug{os.Exit(1)}
	}	
	fmt.Println()
	
	// Write results
	if outFile, err = os.Create("babelDumpBase.dat"); err != nil {panic(err)}
	defer outFile.Close()
	
	for i:=0; i<len(bodies); i++ {
		fmt.Fprintf(outFile, "%d %f %f %f %f %f %f %f\n", 
			minusOne, bodies[i].m, bodies[i].x, bodies[i].y, bodies[i].z, bodies[i].vx, bodies[i].vy, bodies[i].vx)
	}
}











