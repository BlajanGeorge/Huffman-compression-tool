package main

import (
	"log"
	"strings"
)

func ValidateArgsNumber(args []string) {
	if len(args) != 3 && len(args) != 4 {
		log.Fatal("Wrong arguments number, format of the input should be [-op (optional), input, destination]")
	}
}

func ValidateFileName(args []string, inputIndex int) {
	if len(strings.TrimSpace(args[inputIndex])) == 0 {
		log.Fatalf("Invalid input name")
	}
}

func ValidateDestination(args []string, dstIndex int) {
	if len(strings.TrimSpace(args[dstIndex])) == 0 {
		log.Fatalf("Invalid destination name")
	}
}

func ValidateOperationFlag(flag string) {
	if flag != "compression" && flag != "decompression" {
		log.Fatalf("Invalid operation flag")
	}
}
