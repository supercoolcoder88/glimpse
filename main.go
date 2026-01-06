package main

import (
	"fmt"
	"glimpse/logs"
	"os"
)

func main() {
	err := logs.Read(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
