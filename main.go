// main.go
package main

import (
	"os"

	"hearx/pkg/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
