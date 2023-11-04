package main

import (
	"os"
)

func main() {
	ValidateArgsNumber(os.Args)
	ValidateFileName(os.Args)
	ValidateDestination(os.Args)
	compress(os.Args[1])
}
