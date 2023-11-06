package main

import (
	"flag"
	"os"
)

func main() {
	inputIndex, dstIndex := computeArgsIndex(os.Args)
	operationFlagPointer := flag.String("op", "compression", "Flag to state the operation to be done, default value is compression.")
	ValidateArgsNumber(os.Args)
	ValidateFileName(os.Args, inputIndex)
	ValidateDestination(os.Args, dstIndex)
	flag.Parse()
	ValidateOperationFlag(*operationFlagPointer)
	if *operationFlagPointer == "compression" {
		Compress(os.Args[inputIndex], os.Args[dstIndex])
	} else {
		Decompress(os.Args[inputIndex], os.Args[dstIndex])
	}
}

func computeArgsIndex(args []string) (int, int) {
	if len(args) == 3 {
		return 1, 2
	} else {
		return 2, 3
	}
}
