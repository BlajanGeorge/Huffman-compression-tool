package main

import (
	"log"
	"strings"
)

func ValidateArgsNumber(args []string) {
	if len(args) != 3 {
		log.Fatal("Wrong arguments number, format of the input should be [input, destination]")
	}
}

func ValidateFileName(args []string) {
	if len(strings.TrimSpace(args[1])) == 0 {
		log.Fatalf("Invalid input name")
	}
}

func ValidateDestination(args []string) {
	if len(strings.TrimSpace(args[2])) == 0 {
		log.Fatalf("Invalid destination name")
	}
}
