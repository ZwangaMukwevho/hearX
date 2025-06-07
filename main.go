// main.go
package main

import (
	"os"

	"hearx/pkg/cli"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "todo",
		Short: "Todo-service (serve + client)",
	}

	// global flags for client
	root.PersistentFlags().StringVar(&cli.Host, "host", "localhost", "gRPC host")
	root.PersistentFlags().StringVar(&cli.Port, "port", "50051", "gRPC port")

	// register sub-commands
	root.AddCommand(
		cli.ServeCmd(),
		cli.AddCmd(),
		cli.GetCmd(),
		cli.CompleteCmd(),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
