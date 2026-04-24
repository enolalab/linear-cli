package main

import (
	"os"

	"github.com/enolalab/linear-cli/cmd"
)

func main() {
	code := cmd.Execute()
	os.Exit(code)
}
