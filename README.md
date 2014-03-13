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

````bash
$ time ./goBabelS ../inputFilesSmall/input1k 
2014/03/13 11:06:08 profile: cpu profiling enabled, /tmp/profile201142823/cpu.pprof
Read  1024 lines
Starting: Etot0=-0.250000 Ek0=0.250000 Ep0=-0.500000
t= 1.000000 Etot=-0.189654 Etot0=-0.250000 Ek=0.220449 Ep=-0.410103 dE=-0.241383
2014/03/13 11:06:24 Wall time for  main.main :  16.744109844s

real    0m16.747s
user    0m16.669s
sys     0m0.064s

$ time ./goBabelS ../inputFilesBig/input16k 
2014/03/13 11:13:06 profile: cpu profiling enabled, /tmp/profile769946569/cpu.pprof
Read  16384 lines
Starting: Etot0=-0.250000 Ek0=0.250000 Ep0=-0.500000
t= 1.000000 Etot=-0.189432 Etot0=-0.250000 Ek=0.218587 Ep=-0.408020 dE=-0.242270
2014/03/13 12:24:37 Wall time for  main.main :  1h11m31.085543521s

real    71m31.090s
user    71m26.184s
sys     0m5.532s
````

goBabelSF
=========

Here distances are computed in a function and stored in a slice.

````bash
$ time ./goBabelSF ../inputFilesSmall/input1k
2014/03/13 11:11:08 profile: cpu profiling enabled, /tmp/profile919886927/cpu.pprof
Read  1024 lines
Starting: Etot0=-0.250000 Ek0=0.250000 Ep0=-0.500000
t= 1.000000 Etot=-0.069217 Etot0=-0.250000 Ek=0.291366 Ep=-0.360583 dE=-0.723132
2014/03/13 11:11:28 Wall time for  main.main :  20.268483801s

real    0m20.273s
user    0m20.193s
sys     0m0.076s

$ time ./goBabelSF ../inputFilesBig/input16k 
2014/03/13 11:13:17 profile: cpu profiling enabled, /tmp/profile243239193/cpu.pprof
Read  16384 lines
Starting: Etot0=-0.250000 Ek0=0.250000 Ep0=-0.500000
t= 1.000000 Etot=0.713182 Etot0=-0.250000 Ek=1.058857 Ep=-0.345675 dE=-3.8527273
2014/03/13 12:39:53 Wall time for  main.main :  1h26m35.96436336s

real    86m36.405s
user    86m27.384s
sys     0m9.013s
````

goBabelSOS
==========

Here particles are represented by a struct.
No function for compute distances.

````bash
$ time ./goBabelSOS ../inputFilesSmall/input1k
2014/03/13 11:11:11 profile: cpu profiling enabled, /tmp/profile034141521/cpu.pprof
Read  1024 lines
Starting: Etot0=-0.250000 Ek0=0.250000 Ep0=-0.500000
t= 1.000000 Etot=-0.189654 Etot0=-0.250000 Ek=0.220449 Ep=-0.410103 dE=-0.241383
2014/03/13 11:11:28 Wall time for  main.main :  16.775445217s

real    0m16.778s
user    0m16.717s
sys     0m0.064s

$ time ./goBabelSOS ../inputFilesBig/input16k 
2014/03/13 11:13:27 profile: cpu profiling enabled, /tmp/profile507866436/cpu.pprof
Read  16384 lines
Starting: Etot0=-0.250000 Ek0=0.250000 Ep0=-0.500000
t= 1.000000 Etot=-0.189432 Etot0=-0.250000 Ek=0.218587 Ep=-0.408020 dE=-0.242270
2014/03/13 12:24:29 Wall time for  main.main :  1h11m2.713812333s

real    71m2.717s
user    70m56.062s
sys     0m6.864s
````

goBabelSOSOSF
=============

Here there's a struct representing the cluster. It contains the slices of body
structs representing the particles.
A function compute and store distances.

````bash
$ time ./goBabelSOSOSF ../inputFilesSmall/input1k
2014/03/13 11:11:14 profile: cpu profiling enabled, /tmp/profile202456972/cpu.pprof
Read  1024 lines
Starting: Etot0=-0.250000 Ek0=0.250000 Ep0=-0.500000
t= 1.000000 Etot=-0.189654 Etot0=-0.250000 Ek=0.220449 Ep=-0.410103 dE=-0.241383
2014/03/13 11:11:39 Wall time for  main.main :  24.983169915s

real    0m24.987s
user    0m24.894s
sys     0m0.092s

$ time ./goBabelSOSOSF ../inputFilesBig/input16k 
2014/03/13 11:13:36 profile: cpu profiling enabled, /tmp/profile771707555/cpu.pprof
Read  16384 lines
Starting: Etot0=-0.250000 Ek0=0.250000 Ep0=-0.500000
t= 1.000000 Etot=-0.189432 Etot0=-0.250000 Ek=0.218587 Ep=-0.408020 dE=-0.242270
2014/03/13 13:00:26 Wall time for  main.main :  1h46m49.940454341s

real    106m50.363s
user    106m41.920s
sys     0m8.281s
````

MichaelJones
============

Provided by Michael Jones on the Golang-nuts mailing list
https://groups.google.com/d/msg/golang-nuts/lhAD3LnfO88/jg6e5OZfSYYJ

````bash
$ time ./MichaelJones ../inputFilesSmall/input1k 
2014/03/13 14:15:30 profile: cpu profiling enabled, /tmp/profile823321229/cpu.pprof
Read  1024 lines
Starting: Etot0=-0.250000 Ek0=0.250000 Ep0=-0.500000
t= 1.000000 Etot=-0.189654 Etot0=-0.250000 Ek=0.220449 Ep=-0.410103 dE=-0.241383

real    0m11.319s
user    0m11.241s
sys     0m0.088s

$ time ./MichaelJones ../inputFilesBig/input16k 
2014/03/13 14:16:25 profile: cpu profiling enabled, /tmp/profile043838095/cpu.pprof
Read  16384 lines
Starting: Etot0=-0.250000 Ek0=0.250000 Ep0=-0.500000
t= 1.000000 Etot=-0.189432 Etot0=-0.250000 Ek=0.218587 Ep=-0.408020 dE=-0.242270

real    47m35.187s                                                                                                                   
user    47m32.970s                                                                                                                   
sys     0m4.088s
````












