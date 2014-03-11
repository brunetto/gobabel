package main

import (
	"fmt"
	"math"
	"os"
)

var (
	n   int
	m   = []float64{} //Masses array
	rij = [3]float64{}
	r   = [][3]float64{} //Positions
	v   = [][3]float64{} //Velocities
	a   [][3]float64     //Accelerations
	a0  [][3]float64     //Prev accelerations
)

func acceleration() {
	// Reset acceleration
	for i := 0; i < n; i++ {
		a[i][0] = 0
		a[i][1] = 0
		a[i][2] = 0
	}

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			rij[0] = r[i][0] - r[j][0]
			rij[1] = r[i][1] - r[j][1]
			rij[2] = r[i][2] - r[j][2]

			RdotR := (rij[0] * rij[0]) + (rij[1] * rij[1]) + (rij[2] * rij[2])
			apre := 1.0 / math.Sqrt(RdotR*RdotR*RdotR)

			//Update acceleration
			a[i][0] -= m[j] * apre * rij[0]
			a[i][1] -= m[j] * apre * rij[1]
			a[i][2] -= m[j] * apre * rij[2]
		}
	}
}

// Update positions
func updatePositions(dt float64) {
	for i := 0; i < n; i++ {
		// Update the positions, based on the calculated accelerations and velocities
		a0[i][0] = a[i][0]
		a0[i][1] = a[i][1]
		a0[i][2] = a[i][2]
		// For each axis (x/y/z)
		r[i][0] += dt*v[i][0] + 0.5*dt*dt*a0[i][0]
		r[i][1] += dt*v[i][1] + 0.5*dt*dt*a0[i][1]
		r[i][2] += dt*v[i][2] + 0.5*dt*dt*a0[i][2]
	}

}

// Update velocities based on previous and new accelerations
func updateVelocities(dt float64) {
	//Update the velocities based on the previous and old accelerations
	for i := 0; i < n; i++ {
		v[i][0] += 0.5 * dt * (a0[i][0] + a[i][0])
		v[i][1] += 0.5 * dt * (a0[i][1] + a[i][1])
		v[i][2] += 0.5 * dt * (a0[i][2] + a[i][2])

		a0[i][0] = a[i][0]
		a0[i][1] = a[i][1]
		a0[i][2] = a[i][2]
	}
}

// Compute the energy of the system,
// contains an expensive O(N^2) part which can be moved to the acceleration part
// where this is already calculated
func energies() (EKin, EPot float64) {
	EKin = 0

	//Kinetic energy
	for i := 0; i < n; i++ {
		EKin += 0.5 * m[i] * (v[i][0]*v[i][0] + v[i][1]*v[i][1] + v[i][2]*v[i][2])
	}

	//Potential energy
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			//Distance between the two stars
			rij[0] = r[i][0] - r[j][0]
			rij[1] = r[i][1] - r[j][1]
			rij[2] = r[i][2] - r[j][2]

			EPot -= ((m[i] * m[i]) / math.Sqrt((rij[0]*rij[0])+(rij[1]*rij[1])+(rij[2]*rij[2])))
		}
	}
	return EKin, EPot
}

func main() {

	var (
		t                                               float64 = 0.0
		tend                                            float64 = 1.0
		dt                                              float64 = 1e-3
		k                                               int     = 0
		kinEnergy, potEnergy, totEnergy, totEnergy0, dE float64
		inFileName                                      string
		inFile                                          *os.File
		outFile                                         *os.File
		err                                             error
		minusOne                                        int
		tempm, x, y, z, vx, vy, vz                      float64
	)

	inFileName = os.Args[1]

	if inFile, err = os.Open(inFileName); err != nil {
		panic(err)
	}
	defer inFile.Close()

	for {

		if _, err := fmt.Fscanf(inFile, "%d %f %f %f %f %f %f %f\n",
			&minusOne, &(tempm), &(x), &(y), &(z), &(vx), &(vy), &(vz)); err != nil {
			break
		}

		m = append(m, tempm)
		r = append(r, [3]float64{x, y, z})
		v = append(v, [3]float64{vx, vy, vz})
	}
	n = len(m)
	fmt.Println("Read ", n, "lines")
	// Init accelerations
	a = make([][3]float64, n)
	a0 = make([][3]float64, n)
	// Compute initial energy of the system
	kinEnergy, potEnergy = energies()
	totEnergy0 = kinEnergy + potEnergy

	fmt.Printf("Starting: Etot0=%f Ek0=%f Ep0=%f\n", totEnergy0, kinEnergy, potEnergy)

	//Initialize the accelerations
	acceleration()
	//Start the main loop
	for {
		if t > tend {
			break
		}

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
			totEnergy = kinEnergy + potEnergy
			dE = (totEnergy - totEnergy0) / totEnergy0
			
			fmt.Printf("\rt= %f Etot=%f Etot0=%f Ek=%f Ep=%f dE=%f", t, totEnergy, totEnergy0, kinEnergy, potEnergy, dE)
		}
	}
	fmt.Println()

	// Write results
	if outFile, err = os.Create("babelDumpBase0.dat"); err != nil {
		panic(err)
	}
	defer outFile.Close()

	for i := 0; i < n; i++ {
		fmt.Fprintf(outFile, "%d %f %f %f %f %f %f %f\n",
			minusOne, m[i], r[i][0], r[i][1], r[i][2], v[i][0], v[i][1], v[i][2])
	}
}
