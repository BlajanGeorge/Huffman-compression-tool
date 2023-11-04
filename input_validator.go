package main

import (
	"log"
	"strings"
)

func ValidateArgsNumber(args []string) {
	if len(args) < 2 || len(args) > 3 {
		log.Fatal("Wrong arguments number, format of the input should be [fileName, destination]")
	}
}

func ValidateFileName(args []string) {
	if len(strings.TrimSpace(args[1])) == 0 {
		log.Fatalf("Invalid file name")
	}
}

func ValidateDestination(args []string) {
	if len(args) == 3 && len(strings.TrimSpace(args[2])) == 0 {
		log.Fatalf("Invalid destination")
	}
}
