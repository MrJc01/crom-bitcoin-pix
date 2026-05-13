package main

import (
	"fmt"
	"os"
)

// Version é injetada em build time via ldflags.
var Version = "dev"

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
