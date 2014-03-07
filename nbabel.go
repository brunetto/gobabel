package main

import (
	"math"
)

type vec struc {
	x float64
	y float64
	z float64
}

var (
	n int64 					// number of particles
	m = make([]float64, 0)	 	// particle masses
	r = make([]vec, 0)			// particle radii
	v = make([]vec, 0)			// particle velocities
	a = make([]vec, 0)			// particle accelerations
	a0  = make([]vec, 0)		// particle Prev accelerations
	rij = new{vec}				// distance between two particles
)

func acceleration () () {
	// Reset acceleration
	for idx := 0; idx<len(); idx++ {
		a[i].x = 0
		a[i].y = 0
		a[i].z = 0
	}
	
	// Calculate ??
	for i:=0; i<n; i++ {
		for j:=i+1; j<n; j++ {
			rij.x = r[i].x - r[j].x
			rij.y = r[i].y - r[j].y
			rij.z = r[i].z - r[j].z
			// ??
			RdotR = (rij.x * rij.x) + (rij.y * rij.y) + (rij.z * rij.z)
			apre  = 1.0 / math.Sqrt(RdotR * RdotR * RdotR)
		
			//Update acceleration
			a[i].x -= m[j] * apre * rij.x
			a[i].y -= m[j] * apre * rij.y
			a[i].z -= m[j] * apre * rij.z
		}
	}
}

// Update positions
func updatePositions (dt float64) () {
	for i:=0; i<n; i++ {
		// Update the positions, based on the calculated accelerations and velocities
		a0[i].x = a[i].x
		a0[i].y = a[i].y
		a0[i].z = a[i].z
		// For each axis (x/y/z)
		r[i].x += dt * v[i].x + 0.5 * dt * dt * a0[i].x
		r[i].y += dt * v[i].y + 0.5 * dt * dt * a0[i].y
		r[i].z += dt * v[i].z + 0.5 * dt * dt * a0[i].z
	}
	
}

// Update velocities based on previous and new accelerations
func updateVelocities (dt float64) () {
	//Update the velocities based on the previous and old accelerations
	for i:=0; i<n; i++ {
		v[i].x += 0.5 * dt * (a0[i].x + a[i].x)
		v[i].y += 0.5 * dt * (a0[i].y + a[i].y)
		v[i].z += 0.5 * dt * (a0[i].z + a[i].z)
		
		// Update accelerations 
		a0[i].x = a[i].x
		a0[i].y = a[i].y
		a0[i].z = a[i].z
	}  
} 


// Compute the energy of the system, 
// contains an expensive O(N^2) part which can be moved to the acceleration part
// where this is already calculated
void energies () (EKin, EPot float64) {
	EKin = 0
	
	//Kinetic energy
	for i:=0; i<n; i++ {
		EKin += 0.5 * m[i] * v[i].x * v[i].x + v[i].y * v[i].y + v[i].z * v[i].z
	}
	
	//Potential energy
	for i:=0; i<n; i++ {
		for j:=i+1; j<n; j++ {
			//Distance between the two stars
			rij.x = r[i].x - r[j].x
			rij.y = r[i].y - r[j].y
			rij.z = r[i].z - r[j].z

			EPot -= (m[i] * m[j]) / math.Sqrt((rij.x * rij.x) + (rij.y * rij.y) + (rij.z * rij.z))
		}
	}
	return EKin, EPot
}


func main () (int) {

	
	
}











