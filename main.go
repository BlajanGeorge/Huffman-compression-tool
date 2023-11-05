package main

import (
	"os"
)

func main() {
	ValidateArgsNumber(os.Args)
	ValidateFileName(os.Args)
	ValidateDestination(os.Args)
	Compress(os.Args[1], os.Args[2])
}
