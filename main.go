package main

import (
	"fmt"
	"os"

	"github.com/codilime/floodgate/cmd"
)

func main() {
	if err := cmd.Execute(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
