module github.com/canta2899/logo-ls

go 1.22.0

toolchain go1.23.2

require (
	github.com/mattn/go-colorable v0.1.14
	github.com/pborman/getopt/v2 v2.0.0-00010101000000-000000000000
	golang.org/x/sys v0.29.0
	golang.org/x/term v0.28.0
)

require github.com/mattn/go-isatty v0.0.20 // indirect

replace github.com/pborman/getopt/v2 => github.com/rkennedy/getopt/v2 v2.0.0-20231113214939-8e2a55a5c3b3
