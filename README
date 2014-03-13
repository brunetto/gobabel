Go version of a n-body integrator for the site [http://www.nbabel.org](http://www.nbabel.org). 

Each folder has 

* the source
* the binary (compiled with go version go1.2.1 linux/amd64)
* the pprof file for 1k particles and 16k particles input
* the callgrind, the pdf, the svg output of go tool pprof

Input files are in the inputSmallFiles and inputBigFiles folders.
Timing.dat contains the timing I obtained on a Intel(R) Xeon(R) CPU E5-2620 @ 2.00GHz.

goBabelS
========

Basic version copied from the C version
that can be found at http://www.nbabel.org/codes/13.

goBabelSF
=========

Here distances are computed in a function and stored in a slice.

goBabelSOS
==========

Here particles are represented by a struct.
No function for compute distances.

goBabelSOSOSF
=============

Here there's a struct representing the cluster. It contains the slices of body
structs representing the particles.
A function compute and store distances.

MichaelJones
============

Provided by Michael Jones on the Golang-nuts mailing list
https://groups.google.com/d/msg/golang-nuts/lhAD3LnfO88/jg6e5OZfSYYJ














