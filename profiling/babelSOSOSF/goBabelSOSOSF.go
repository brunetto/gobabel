// Here there's a struct representing the cluster. It contains the slices of body
// structs representing the particles.
// A function compute and store distances.

package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/brunetto/goutils/debug"
	"github.com/davecheney/profile"
)

var Debug = false

type body struct {
	x, y, z, vx, vy, vz, ax, ay, az, a0x, a0y, a0z, m float64
}

type Cluster struct {
	N          int
	Bds        []body                      // bodies slice
	Rij        []struct{ x, y, z float64 } // distance between two particles
	RdotR      []float64
	KinEnergy  float64 // kinetic energy of the system
	PotEnergy  float64 // potential energy of the system
	TotEnergy  float64 // total energy of the system
	TotEnergy0 float64 // initial total energy of the system
	DE         float64 // relative energy error
}

func (cl *Cluster) ComputeDistances() {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	k := 0
	for i := 0; i < cl.N; i++ {
		for j := i + 1; j < cl.N; j++ {
			//Distance between the two stars
			cl.Rij[k].x = cl.Bds[i].x - cl.Bds[j].x
			cl.Rij[k].y = cl.Bds[i].y - cl.Bds[j].y
			cl.Rij[k].z = cl.Bds[i].z - cl.Bds[j].z
			cl.RdotR[k] = (cl.Rij[k].x * cl.Rij[k].x) +
				(cl.Rij[k].y * cl.Rij[k].y) +
				(cl.Rij[k].z * cl.Rij[k].z)
			k++
		}
	}
}

// Init create a new cluster and load particles from file
func (cl *Cluster) Init(inFileName string) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	var (
		minusOne int
		inFile   *os.File
		err      error
	)

	// Create new cluster
	cl.Bds = []body{} // the same as make([]*body, 0)

	// Load particles
	if inFile, err = os.Open(inFileName); err != nil {
		panic(err)
	}
	defer inFile.Close()

	for {
		bd := body{}
		if _, err := fmt.Fscanf(inFile, "%d %f %f %f %f %f %f %f\n",
			&minusOne, &(bd.m), &(bd.x), &(bd.y), &(bd.z), &(bd.vx), &(bd.vy), &(bd.vz)); err != nil {
			break
		}
		cl.Bds = append(cl.Bds, bd)
	}
	cl.N = len(cl.Bds)
	fmt.Println("Read ", cl.N, "lines")

	// Create distances slice
	cl.Rij = make([]struct{ x, y, z float64 }, cl.N*(cl.N-1)/2)
	cl.RdotR = make([]float64, cl.N*(cl.N-1)/2)
}

// Acceleration calculate the acceleration for each particle
func (cl *Cluster) Acceleration() {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	// Reset acceleration
	for i := 0; i < cl.N; i++ {
		cl.Bds[i].ax = 0
		cl.Bds[i].ay = 0
		cl.Bds[i].az = 0
	}

	// Compute distances
	cl.ComputeDistances()

	k := 0
	for i := 0; i < cl.N; i++ {
		for j := i + 1; j < cl.N; j++ {
			apre := 1.0 / math.Sqrt(cl.RdotR[k]*cl.RdotR[k]*cl.RdotR[k])

			//Update acceleration
			cl.Bds[i].ax -= cl.Bds[j].m * apre * cl.Rij[k].x
			cl.Bds[i].ay -= cl.Bds[j].m * apre * cl.Rij[k].y
			cl.Bds[i].az -= cl.Bds[j].m * apre * cl.Rij[k].z
			k++
		}
	}
}

// UpdatePositions updates each particle position
func (cl *Cluster) UpdatePositions(dt float64) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	for i := 0; i < cl.N; i++ {
		// Update the positions, based on the calculated accelerations and velocities
		cl.Bds[i].a0x = cl.Bds[i].ax
		cl.Bds[i].a0y = cl.Bds[i].ay
		cl.Bds[i].a0z = cl.Bds[i].az
		// For each axis (x/y/z)
		cl.Bds[i].x += dt*cl.Bds[i].vx + 0.5*dt*dt*cl.Bds[i].a0x
		cl.Bds[i].y += dt*cl.Bds[i].vy + 0.5*dt*dt*cl.Bds[i].a0y
		cl.Bds[i].z += dt*cl.Bds[i].vz + 0.5*dt*dt*cl.Bds[i].a0z
	}

}

// Update velocities based on previous and new accelerations
func (cl *Cluster) UpdateVelocities(dt float64) {
	if Debug {
		defer debug.TimeMe(time.Now())
	}
	//Update the velocities based on the previous and old accelerations
	for i := 0; i < cl.N; i++ {
		cl.Bds[i].vx += 0.5 * dt * (cl.Bds[i].a0x + cl.Bds[i].ax)
		cl.Bds[i].vy += 0.5 * dt * (cl.Bds[i].a0y + cl.Bds[i].ay)
		cl.Bds[i].vz += 0.5 * dt * (cl.Bds[i].a0z + cl.Bds[i].az)

		cl.Bds[i].a0x = cl.Bds[i].ax
		cl.Bds[i].a0y = cl.Bds[i].ay
		cl.Bds[i].a0z = cl.Bds[i].az
	}
}

// Compute the energy of the system,
// contains an expensive O(N^2) part which can be moved to the acceleration part
// where this is already calculated
func (cl *Cluster) Energies() {
	if Debug {
		defer debug.TimeMe(time.Now())
	}

	cl.KinEnergy = 0
	cl.PotEnergy = 0

	//Kinetic energy
	for i := 0; i < cl.N; i++ {
		cl.KinEnergy += 0.5 * cl.Bds[i].m * (cl.Bds[i].vx*cl.Bds[i].vx + cl.Bds[i].vy*cl.Bds[i].vy + cl.Bds[i].vz*cl.Bds[i].vz)
	}

	//Potential energy
	k := 0
	for i := 0; i < cl.N; i++ {
		for j := i + 1; j < cl.N; j++ {
			cl.PotEnergy -= ((cl.Bds[i].m * cl.Bds[j].m) / math.Sqrt(cl.RdotR[k]))
			k++
		}
	}
}

func main() {
	defer profile.Start(profile.CPUProfile).Stop()
	defer debug.TimeMe(time.Now())

	var (
		t          float64 = 0.0  // time
		tEnd       float64 = 1.0  // end of simulation
		dt         float64 = 1e-3 // timestep
		k          int     = 0
		inFileName string
		outFile    *os.File
		err        error
	)

	inFileName = os.Args[1]

	cl := &Cluster{}

	cl.Init(inFileName)

	// Compute initial energy of the system
	cl.ComputeDistances()
	cl.Energies()
	cl.TotEnergy0 = cl.KinEnergy + cl.PotEnergy

	fmt.Printf("Starting: Etot0=%f Ek0=%f Ep0=%f\n", cl.TotEnergy0, cl.KinEnergy, cl.PotEnergy)

	//Initialize the accelerations
	cl.Acceleration()
	//Start the main loop
	for {
		if t > tEnd {
			break
		}

		// Update positions based on velocities and accelerations
		cl.UpdatePositions(dt)

		// Get new accelerations
		cl.Acceleration()

		// Update velocities
		cl.UpdateVelocities(dt)

		// Update time
		t += dt
		k += 1

		if k%10 == 0 {
			cl.Energies()
			cl.TotEnergy = cl.KinEnergy + cl.PotEnergy
			cl.DE = (cl.TotEnergy - cl.TotEnergy0) / cl.TotEnergy0

			fmt.Printf("\rt= %f Etot=%f Etot0=%f Ek=%f Ep=%f dE=%f", t, cl.TotEnergy, cl.TotEnergy0, cl.KinEnergy, cl.PotEnergy, cl.DE)
		}
		if Debug {
			os.Exit(1)
		}
	}
	fmt.Println()

	// Write results

	if outFile, err = os.Create("babelDumpSOSOSF.dat"); err != nil {
		panic(err)
	}
	defer outFile.Close()

	for i := 0; i < cl.N; i++ {
		fmt.Fprintf(outFile, "%d %f %f %f %f %f %f %f\n",
			-1, cl.Bds[i].m, cl.Bds[i].x, cl.Bds[i].y, cl.Bds[i].z, cl.Bds[i].vx, cl.Bds[i].vy, cl.Bds[i].vz)
	}
}
