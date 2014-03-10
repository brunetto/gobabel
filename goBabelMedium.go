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

type Cluster struct {
	N          int
	Bds        []*body                    // bodies slice
	Rij        *struct{ x, y, z float64 } // distance between two particles
	KinEnergy  float64                    // kinetic energy of the system
	PotEnergy  float64                    // potential energy of the system
	TotEnergy  float64                    // total energy of the system
	TotEnergy0 float64                    // initial total energy of the system
	DE         float64                    // relative energy error
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
	cl.Bds = make([]*body, 0)

	// Load particles
	if inFile, err = os.Open(inFileName); err != nil {
		panic(err)
	}
	defer inFile.Close()

	for {
		bd := &body{}
		if _, err := fmt.Fscanf(inFile, "%d %f %f %f %f %f %f %f\n",
			&minusOne, &(bd.m), &(bd.x), &(bd.y), &(bd.z), &(bd.vx), &(bd.vy), &(bd.vx)); err != nil {
			break
		}
		cl.Bds = append(cl.Bds, bd)
	}
	cl.N = len(cl.Bds)
	fmt.Println("Read ", cl.N, "lines")

	// Create distances slice
	cl.Rij = new(struct{ x, y, z float64 }) //= make([]struct{x, y, z float64}, cl.N * (cl.N-1) / 2)
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
	for i := 0; i < cl.N; i++ {
		for j := i + 1; j < cl.N; j++ {

			cl.Rij.x = cl.Bds[i].x - cl.Bds[j].x
			cl.Rij.y = cl.Bds[i].y - cl.Bds[j].y
			cl.Rij.z = cl.Bds[i].z - cl.Bds[j].z

			RdotR := (cl.Rij.x * cl.Rij.x) + (cl.Rij.y * cl.Rij.y) + (cl.Rij.z * cl.Rij.z)

			apre := 1.0 / math.Sqrt(RdotR*RdotR*RdotR)

			//Update acceleration
			cl.Bds[i].ax -= cl.Bds[j].m * apre * cl.Rij.x
			cl.Bds[i].ay -= cl.Bds[j].m * apre * cl.Rij.y
			cl.Bds[i].az -= cl.Bds[j].m * apre * cl.Rij.z
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

		// Update accelerations
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
		cl.KinEnergy += 0.5*cl.Bds[i].m*cl.Bds[i].vx*cl.Bds[i].vx + cl.Bds[i].vy*cl.Bds[i].vy + cl.Bds[i].vz*cl.Bds[i].vz
	}

	//Potential energy
	k := 0
	for i := 0; i < cl.N; i++ {
		for j := i + 1; j < cl.N; j++ {

			cl.Rij.x = cl.Bds[i].x - cl.Bds[j].x
			cl.Rij.y = cl.Bds[i].y - cl.Bds[j].y
			cl.Rij.z = cl.Bds[i].z - cl.Bds[j].z

			cl.PotEnergy -= (cl.Bds[i].m * cl.Bds[j].m) / math.Sqrt((cl.Rij.x*cl.Rij.x)+
				(cl.Rij.y*cl.Rij.y)+(cl.Rij.z*cl.Rij.z))
			k++
		}
	}
}

func main() {
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

	cl := new(Cluster)

	cl.Init(inFileName)

	// Compute initial energy of the system
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

	if outFile, err = os.Create("babelDumpMedium.dat"); err != nil {
		panic(err)
	}
	defer outFile.Close()

	for i := 0; i < cl.N; i++ {
		fmt.Fprintf(outFile, "%d %f %f %f %f %f %f %f\n",
			-1, cl.Bds[i].m, cl.Bds[i].x, cl.Bds[i].y, cl.Bds[i].z, cl.Bds[i].vx, cl.Bds[i].vy, cl.Bds[i].vx)
	}
}
