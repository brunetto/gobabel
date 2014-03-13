// Provided by Michael Jones on the Golang-nuts mailing list
// https://groups.google.com/d/msg/golang-nuts/lhAD3LnfO88/jg6e5OZfSYYJ

package main

import (
	"fmt"
	"math"
	"os"
	
	"github.com/davecheney/profile"
)

type Vec3 struct {
	x, y, z float64
}

type Point struct {
	m  float64
	r  Vec3
	v  Vec3
	a  Vec3
	a0 Vec3
}

func acceleration(point []Point) {
	// Reset acceleration
	for i := range point {
		pi := &point[i]
		pi.a.x = 0
		pi.a.y = 0
		pi.a.z = 0
	}

	for i := range point {
		pi := &point[i]
		for j := i + 1; j < len(point); j++ {
			pj := &point[j]
			rijx := pi.r.x - pj.r.x
			rijy := pi.r.y - pj.r.y
			rijz := pi.r.z - pj.r.z

			RdotR := rijx*rijx + rijy*rijy + rijz*rijz
			apre := 1.0 / math.Sqrt(RdotR*RdotR*RdotR)

			//Update acceleration
			mm := pj.m * apre
			pi.a.x -= mm * rijx
			pi.a.y -= mm * rijy
			pi.a.z -= mm * rijz
		}
	}
}

// Update positions
func updatePositions(point []Point, dt float64) {
	for i := range point {
		// Update the positions, based on the calculated accelerations and velocities
		p := &point[i]
		p.a0.x = p.a.x
		p.a0.y = p.a.y
		p.a0.z = p.a.z

		// For each axis (x/y/z)
		ss := 0.5 * dt * dt
		p.r.x += dt*p.v.x + ss*p.a0.x
		p.r.y += dt*p.v.y + ss*p.a0.y
		p.r.z += dt*p.v.z + ss*p.a0.z
	}

}

// Update velocities based on previous and new accelerations
func updateVelocities(point []Point, dt float64) {
	//Update the velocities based on the previous and old accelerations
	for i := range point {
		p := &point[i]

		ss := 0.5 * dt
		p.v.x += ss * (p.a0.x + p.a.x)
		p.v.y += ss * (p.a0.y + p.a.y)
		p.v.z += ss * (p.a0.z + p.a.z)

		p.a0.x = p.a.x
		p.a0.y = p.a.y
		p.a0.z = p.a.z
	}
}

// Compute the energy of the system,
// contains an expensive O(N^2) part which can be moved to the acceleration part
// where this is already calculated
func energies(point []Point) (float64, float64) {
	//Kinetic energy
	ke := 0.0
	for i := range point {
		p := &point[i]
		ke += 0.5 * p.m * (p.v.x*p.v.x + p.v.y*p.v.y + p.v.z*p.v.z)
	}

	//Potential energy
	pe := 0.0
	for i := range point {
		pi := &point[i]
		for j := i + 1; j < len(point); j++ {
			pj := &point[j]

			//Distance between the two stars
			rijx := pi.r.x - pj.r.x
			rijy := pi.r.y - pj.r.y
			rijz := pi.r.z - pj.r.z

			pe -= pi.m * pi.m / math.Sqrt(rijx*rijx+rijy*rijy+rijz*rijz)
		}
	}
	return ke, pe
}

func main() {
	defer profile.Start(profile.CPUProfile).Stop()
	
	if len(os.Args) < 2 {
		fmt.Printf("usage %s filename\n", os.Args[0])
		os.Exit(1)
	}
	inFileName := os.Args[1]
	inFile, err := os.Open(inFileName)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()
	var point []Point
	for {
		var minusOne int
		var m, x, y, z, vx, vy, vz float64
		if _, err := fmt.Fscanf(inFile, "%d %f %f %f %f %f %f %f\n",
			&minusOne, &m, &x, &y, &z, &vx, &vy, &vz); err != nil {
			break
		}
		point = append(point, Point{m, Vec3{x, y, z}, Vec3{vx, vy, vz}, Vec3{0, 0, 0}, Vec3{0, 0, 0}})
	}
	fmt.Println("Read ", len(point), "lines")

	// Compute initial energy of the system
	kinEnergy, potEnergy := energies(point)
	totEnergy0 := kinEnergy + potEnergy
	fmt.Printf("Starting: Etot0=%f Ek0=%f Ep0=%f\n", totEnergy0, kinEnergy, potEnergy)

	//Initialize the accelerations
	acceleration(point)

	//Start the main loop
	k := 0
	t := 0.0
	tend := 1.0
	dt := 1.0e-3
	for {
		if t > tend {
			break
		}

		// Update positions based on velocities and accelerations
		updatePositions(point, dt)

		// Get new accelerations
		acceleration(point)

		// Update velocities
		updateVelocities(point, dt)

		// Update time
		t += dt
		k += 1

		if true && k%10 == 0 {
			kinEnergy, potEnergy = energies(point)
			totEnergy := kinEnergy + potEnergy
			dE := (totEnergy - totEnergy0) / totEnergy0
			fmt.Printf("\rt= %f Etot=%f Etot0=%f Ek=%f Ep=%f dE=%f",
				t, totEnergy, totEnergy0, kinEnergy, potEnergy, dE)
		}
	}
	fmt.Println()

	// Write results
	outFile, err := os.Create("MichaelJones.dat")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	for _, p := range point {
		fmt.Fprintf(outFile, "%d %10.6f %10.6f %10.6f %10.6f %10.6f %10.6f %10.6f\n",
			-1, p.m, p.r.x, p.r.y, p.r.z, p.v.x, p.v.y, p.v.z)
	}
}
