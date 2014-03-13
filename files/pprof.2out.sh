#!/bin/bash
set -x
# goBabelS

go tool pprof --callgrind goBabelS goBabelS1k.pprof > ./goBabelS1k.0
go tool pprof --svg goBabelS goBabelS1k.pprof > ./goBabelS1k.svg
go tool pprof --pdf goBabelS goBabelS1k.pprof > ./goBabelS1k.pdf

go tool pprof --callgrind goBabelS goBabelS16k.pprof > ./goBabelS16k.0
go tool pprof --svg goBabelS goBabelS16k.pprof > ./goBabelS16k.svg
go tool pprof --pdf goBabelS goBabelS16k.pprof > ./goBabelS16k.pdf

# goBabelSF

go tool pprof --callgrind goBabelSF goBabelSF1k.pprof > ./goBabelSF1k.0
go tool pprof --svg goBabelSF goBabelSF1k.pprof > ./goBabelSF1k.svg
go tool pprof --pdf goBabelSF goBabelSF1k.pprof > ./goBabelSF1k.pdf

go tool pprof --callgrind goBabelSF goBabelSF16k.pprof > ./goBabelSF16k.0
go tool pprof --svg goBabelSF goBabelSF16k.pprof > ./goBabelSF16k.svg
go tool pprof --pdf goBabelSF goBabelSF16k.pprof > ./goBabelSF16k.pdf

# goBabelSOS

go tool pprof --callgrind goBabelSOS goBabelSOS1k.pprof > ./goBabelSOS1k.0
go tool pprof --svg goBabelSOS goBabelSOS1k.pprof > ./goBabelSOS1k.svg
go tool pprof --pdf goBabelSOS goBabelSOS1k.pprof > ./goBabelSOS1k.pdf

go tool pprof --callgrind goBabelSOS goBabelSOS16k.pprof > ./goBabelSOS16k.0
go tool pprof --svg goBabelSOS goBabelSOS16k.pprof > ./goBabelSOS16k.svg
go tool pprof --pdf goBabelSOS goBabelSOS16k.pprof > ./goBabelSOS16k.pdf

# goBabelSOS

go tool pprof --callgrind goBabelSOSOSF goBabelSOSOSF1k.pprof > ./goBabelSOSOSF1k.0
go tool pprof --svg goBabelSOSOSF goBabelSOSOSF1k.pprof > ./goBabelSOSOSF1k.svg
go tool pprof --pdf goBabelSOSOSF goBabelSOSOSF1k.pprof > ./goBabelSOSOSF1k.pdf

go tool pprof --callgrind goBabelSOSOSF goBabelSOSOSF16k.pprof > ./goBabelSOSOSF16k.0
go tool pprof --svg goBabelSOSOSF goBabelSOSOSF16k.pprof > ./goBabelSOSOSF16k.svg
go tool pprof --pdf goBabelSOSOSF goBabelSOSOSF16k.pprof > ./goBabelSOSOSF16k.pdf

# goBabelSOSOSF16k

go tool pprof --callgrind MichaelJones MichaelJones1k.pprof > ./MichaelJones1k.0
go tool pprof --svg MichaelJones MichaelJones1k.pprof > ./MichaelJones1k.svg
go tool pprof --pdf MichaelJones MichaelJones1k.pprof > ./MichaelJones1k.pdf

go tool pprof --callgrind MichaelJones MichaelJones16k.pprof > ./MichaelJones16k.0
go tool pprof --svg MichaelJones MichaelJones16k.pprof > ./MichaelJones16k.svg
go tool pprof --pdf MichaelJones MichaelJones16k.pprof > ./MichaelJones16k.pdf

